/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hitzhangjie/log"
	"github.com/spf13/cobra"
)

const (
	FuncRegisterPattern = "Register.*Service"
)

var (
	regexFuncRegisterPattern *regexp.Regexp
)

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// TODO use i18n to generate translations for English/Chinese
	visualizeCmd.Flags().StringP("projectdir", "d", "", "specify the directory of existed project")

	visualizeCmd.MarkFlagRequired("projectdir")

	regexFuncRegisterPattern, _ = regexp.Compile(FuncRegisterPattern)
}

// updateCmd represents the update command
var visualizeCmd = &cobra.Command{
	Use: "visualize",
	// TODO i18n settings
	Short: "visualize workflow",
	Long:  `visualize workflow`,
	RunE: func(cmd *cobra.Command, args []string) error {

		d, err := cmd.Flags().GetString("projectdir")
		if err != nil {
			return err
		}

		// 解析main.go
		fset, astFile, err := parseFile(filepath.Join(d, "main.go"))
		if err != nil {
			return err
		}

		// 检查main.main中注册了哪几个逻辑service
		services, err := registeredServices(fset, astFile)
		fmt.Printf("found registered services: %s\n", strings.Join(services, ", "))

		// 解析pb文件，检查service接口定义的method

		// 解析工程找到service定义，从接口方法开始展开，还原整个流程图
		fset, pkgs, err := parseDir(d)
		if err != nil {
			return err
		}

		for _, service := range services {
			methodSteps, _ := parseServiceMethods(fset, pkgs, service)
			for method, steps := range methodSteps {
				//_ = renderSteps(method, steps)
				//fmt.Println("--------------------------------------------")
				_ = renderStepsWithPlantUML(method, steps)
			}
		}

		return err
	},
}

func parseFile(file string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return fset, astFile, nil
}

func registeredServices(fset *token.FileSet, file *ast.File) ([]string, error) {

	services := []string{}

	ast.Inspect(file, func(n ast.Node) bool {

		// only main.main needed analyzing
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if fn.Name.Name != "main" {
			return true
		}

		// traverse all pb.Register${Service}Service statements,
		// it must be x := *ast.ExprStatement.(*ast.CallExpr).(*ast.SelectorExpr),
		//
		// see and test at: https://yuroyoro.github.io/goast-viewer/index.html.
		for _, stmt := range fn.Body.List {
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				continue
			}
			callExpr, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if !regexFuncRegisterPattern.MatchString(selectorExpr.Sel.Name) {
				continue
			}

			service := callExpr.Args[1].(*ast.UnaryExpr).X.(*ast.CompositeLit).Type.(*ast.Ident).Name
			services = append(services, service)
		}
		return true
	})

	return services, nil
}

func parseDir(dir string) (*token.FileSet, map[string]*ast.Package, error) {

	dirs, err := traverseDirs(dir)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	allPkgs := map[string]*ast.Package{}

	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
		if err != nil {
			return nil, nil, err
		}
		for k, v := range pkgs {
			allPkgs[k] = v
		}
	}

	return fset, allPkgs, nil
}

type Phase int

const (
	PhaseStart = iota
	PhaseExpand
	PhaseEnd
)

func parseServiceMethods(fset *token.FileSet, pkgs map[string]*ast.Package, service string) (map[string][]StatementX, error) {

	methodSteps := map[string][]StatementX{}

	// TODO remove hardcoded chan capcity 16
	funcNodesCh := make(chan *ast.FuncDecl, 16)

	for _, pkg := range pkgs {
		for fname, file := range pkg.Files {
			_ = fname

			ast.Inspect(file, func(n ast.Node) bool {

				fn, ok := isServiceMethod(n, service)
				if !ok {
					return true
				}

				// record the rpc methods ast node
				funcNodesCh <- fn

				return true
			})
		}
	}
	close(funcNodesCh)

	// analyze every service method
	for fn := range funcNodesCh {
		funcName, steps := inspectFn(fn, fset, pkgs, -1)
		methodSteps[funcName] = steps
	}

	return methodSteps, nil
}

// inspectFn analyze function code flow, like rpc call hierarchy, control flow, etc
//
// what should we visualize?
// - OOP communication, this depicts the dependencies btw different components
//
// 	 case1 : communication btw components by calling obj's method
// 	 ss := &student{}
// 	 ss.Name()
//
// 	 case2: communication btw components by calling pkg's exported function
// 	 pkg.Statement()
//
// - TODO control flow, if, for, switch, this depicts some important logic
// - TODO concurrency like go func(), serialization like wg.Wait()
func inspectFn(fn *ast.FuncDecl, fset *token.FileSet, pkgs map[string]*ast.Package, depth int) (string, []StatementX) {

	depth++

	steps := []StatementX{}
	service, _ := parseFuncRecvType(fn)
	funcName := fn.Name.Name

	if len(service) != 0 {
		funcName = fmt.Sprintf("%s.%s", service, fn.Name.Name)
	}

	for _, stmt := range fn.Body.List {

		//fmt.Printf("found statement: %+v\n", stmt)

		// TODO arguments
		var (
			pos       = fset.Position(stmt.Pos())
			args      = "..."
			statement = ""
			typ       = ""
			selName   = ""
			callExpr  *ast.CallExpr
		)

		switch stmt.(type) {
		case *ast.ExprStmt:
			exprStmt, ok := stmt.(*ast.ExprStmt)
			if !ok {
				fmt.Printf("fuck\n")
				continue
			}
			call, ok := exprStmt.X.(*ast.CallExpr)
			if !ok {
				continue
			}
			callExpr = call
		case *ast.AssignStmt:
			call, ok := stmt.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			callExpr = call
		default:
			continue
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		// x.Name represents receiver variable, or package name
		x := selectorExpr.X.(*ast.Ident)
		selName = selectorExpr.Sel.Name

		// x.Obj != nil, it's funcName
		// x.Obj == nil, it's package exported function
		if x.Obj != nil { // funcName
			rhs := selectorExpr.X.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs
			typ = rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Type.(*ast.Ident).Name
			if op := rhs[0].(*ast.UnaryExpr).Op.String(); len(op) != 0 {
				if op == "&" {
					op = "*"
				}
				typ = op + typ
			}
		}

		if len(typ) != 0 {
			statement = fmt.Sprintf("%s%s%s (%s)%s.%s(%s)",
				log.COLOR_GREEN, pos, log.COLOR_RESET, typ, x.Name, selName, args)
		} else {
			statement = fmt.Sprintf("%s%s%s %s.%s(%s)\n",
				log.COLOR_GREEN, pos, log.COLOR_RESET, x.Name, selName, args)
		}

		// TODO recursively expand the body at function
		findNode := findFuncNode(pkgs, typ, selName)

		if len(statement) != 0 {
			tmp := strings.TrimSpace(findNode.Doc.Text())
			idx := strings.IndexAny(tmp, " \t")
			comment := tmp[idx+1:]

			steps = append(steps, StatementX{
				Position:           pos,
				Statement:          statement,
				Comment:            comment,
				Caller:             funcName,
				CallHierarchyDepth: depth,
				X:                  x.Name,
				Typ:                typ,
				Seletor:            selName,
				Args:               []string{args},
			})
		}

		if findNode != nil {
			//fmt.Printf("found funcNode, %s.%s, %+v\n", typ, selName, findNode)
			// recursive expand function body
			_, nestedSteps := inspectFn(findNode, fset, pkgs, depth)
			steps = append(steps, nestedSteps...)
		} else {
			fmt.Printf("not found funcNode, %s.%s\n", typ, selName)
		}
	}

	return funcName, steps
}

func traverseDirs(dir string) ([]string, error) {
	dirs := []string{dir}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	return dirs, err
}

func renderSteps(method string, steps []StatementX) error {
	fmt.Printf("%s%s%s\v\r", log.COLOR_RED, method, log.COLOR_RESET)
	for phase, step := range steps {
		prefix := strings.Repeat("\t", step.CallHierarchyDepth)
		idx := strings.IndexAny(step.Comment, " \t")
		comment := step.Comment[idx+1:]
		fmt.Printf("\v\r%s|\v\r%s|%d-%s\v\r%s|\v\r%sv\v\r", prefix, prefix, phase, comment, prefix, prefix)
		fmt.Printf("\v\r%s%s", prefix, step.Statement)
	}
	return nil
}

func renderStepsWithPlantUML(method string, statements []StatementX) error {
	buf := &bytes.Buffer{}

	buf.WriteString("@startuml\n")

	entities := map[string]bool{}

	for _, statement := range statements {
		vals := strings.Split(statement.Caller, ".")
		callerEntity := vals[0]
		if _, ok := entities[callerEntity]; !ok {
			fmt.Fprintf(buf, "participant \"%s\"\n", callerEntity)
			entities[callerEntity] = true
		}
	}

	for _, statement := range statements {
		entity := ""
		if len(statement.Typ) != 0 {
			entity = strings.TrimPrefix(statement.Typ, "*")
		} else {
			entity = statement.X
		}
		operation := statement.Seletor

		if _, ok := entities[entity]; !ok {
			fmt.Fprintf(buf, "participant \"%s\"\n", entity)
			entities[entity] = true
		}

		v := strings.Split(statement.Caller, ".")
		callerEntity, callerAction := v[0], v[1]
		_ = callerAction
		fmt.Fprintf(buf, "\"%s\" -> \"%s\" : %s\n", callerEntity, entity, operation)
		fmt.Fprintf(buf, "activate \"%s\"\n", entity)
		fmt.Fprintf(buf, "note right\n")
		fmt.Fprintf(buf, "%s\n", statement.Comment)
		fmt.Fprintf(buf, "end note\n")
		fmt.Fprintf(buf, "deactivate \"%s\"\n", entity)
	}

	//fmt.Fprintf(buf, "deactivate \"%s\"\n", actor)
	buf.WriteString("@enduml\n")

	fmt.Printf("plantuml data: \n\n%s", buf.String())

	return nil
}

// isServiceMethod check with node `n` is a method definition node of service `service`
func isServiceMethod(n ast.Node, service string) (*ast.FuncDecl, bool) {

	// must be a func declaration
	fn, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil, false
	}

	// fn is function, rather than methods
	if fn.Recv == nil || len(fn.Recv.List) == 0 || fn.Recv.List[0] == nil || fn.Recv.List[0].Type == nil {
		return nil, false
	}

	// gorpc template make sure receiver type of methods of generated implemention of service interface
	// always conforms to form `(s *${service}) RPCMethod(ctx, req, rsp) error`.
	typ, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return nil, false
	}

	// filter out the methods whose receiver type has the same type as registered services
	ident, ok := typ.X.(*ast.Ident)
	if !ok || service != ident.Name {
		return nil, false
	}

	return fn, true
}

func isTargetMethod(n ast.Node, recvType, funcName string) (*ast.FuncDecl, bool) {

	// must be a func declaration
	fn, ok := n.(*ast.FuncDecl)
	if !ok {
		return nil, false
	}

	if len(recvType) != 0 {
		// fn is function, rather than methods
		if fn.Recv == nil || len(fn.Recv.List) == 0 || fn.Recv.List[0] == nil || fn.Recv.List[0].Type == nil {
			return nil, false
		}

		// gorpc template make sure receiver type of methods of generated implemention of service interface
		// always conforms to form `(s *${service}) RPCMethod(ctx, req, rsp) error`.
		typ, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return nil, false
		}

		// filter out the methods whose receiver type has the same type as registered services
		ident, ok := typ.X.(*ast.Ident)
		if !ok || recvType != ident.Name {
			return nil, false
		}
	}

	// filter out the methods whose name not matches
	if funcName != fn.Name.Name {
		return nil, false
	}

	return fn, true
}

func findFuncNode(pkgs map[string]*ast.Package, recvType, funcName string) *ast.FuncDecl {

	// TODO how to seperate *Student and Student
	recvType = strings.TrimPrefix(recvType, "*")

	var findNode *ast.FuncDecl

NEXT:
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := isTargetMethod(n, recvType, funcName)
				if !ok {
					return true
				}
				findNode = fn
				return false
			})
			if findNode != nil {
				break NEXT
			}
		}
	}

	return findNode
}

func parseFuncRecvType(fn *ast.FuncDecl) (string, error) {

	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return "", errors.New("invalid receiver type")
	}

	typ, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
	if !ok {
		return "", errors.New("not *ast.StarExpr")
	}

	ident, ok := typ.X.(*ast.Ident)
	if !ok {
		panic("not *ast.Ident")
	}

	// TODO seperate value and pointer?
	// `(s Student) xxx()` or `(s *Student) xxx()`
	return ident.Name, nil
}

type StatementX struct {
	Position           token.Position
	Statement          string
	Comment            string
	Caller             string
	CallHierarchyDepth int
	X                  string   // receiver variable, or package name
	Typ                string   // receiver type
	Seletor            string   // member function or package exported function
	Args               []string // function args
}
