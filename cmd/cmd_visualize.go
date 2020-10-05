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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hitzhangjie/log"
	"github.com/spf13/cobra"
)

const (
	FuncMain            = "main"
	FuncRegisterService = "RegisterHelloSvrService"
	FuncRegisterPattern = "Register.*Service"
)

var (
	regex *regexp.Regexp
)

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// TODO use i18n to generate translations for English/Chinese
	visualizeCmd.Flags().StringP("projectdir", "d", "", "specify the directory of existed project")

	visualizeCmd.MarkFlagRequired("projectdir")

	regex, _ = regexp.Compile(FuncRegisterPattern)
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
		//fmt.Printf("found registered services: %s\n", strings.Join(services, ", "))

		// 解析pb文件，检查service接口定义的method

		// 解析工程找到service定义，从接口方法开始展开，还原整个流程图
		fset, pkgs, err := parseDir(d)
		if err != nil {
			return err
		}

		for _, service := range services {
			methodSteps, _ := parseServiceMethods(fset, pkgs, service)
			for method, steps := range methodSteps {
				_ = renderSteps(method, steps)
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
			if !regex.MatchString(selectorExpr.Sel.Name) {
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

func parseServiceMethods(fset *token.FileSet, pkgs map[string]*ast.Package, service string) (map[string][]string, error) {

	methodSteps := map[string][]string{}

	for _, pkg := range pkgs {
		for fname, file := range pkg.Files {
			_ = fname

			ast.Inspect(file, func(n ast.Node) bool {
				fn, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}

				// function, rather than methods
				if fn.Recv == nil || len(fn.Recv.List) == 0 || fn.Recv.List[0] == nil || fn.Recv.List[0].Type == nil {
					return true
				}

				typ, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
				if !ok {
					return true
				}

				ident, ok := typ.X.(*ast.Ident)
				if !ok || service != ident.Name {
					return true
				}

				method := fmt.Sprintf("%s.%s", service, fn.Name)
				steps := []string{}

				// TODO what should we visualize?
				// - OOP communication, this depicts the dependencies btw different components
				// - control flow, if, for, switch, this depicts some important logic

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

					// TODO OOP communication

					// case1 : communication btw components by calling obj's method
					// ss := &student{}
					// ss.Name()
					//
					// case2: communication btw components by calling pkg's exported function
					// pkg.Func()
					x := selectorExpr.X.(*ast.Ident)
					xName := x.Name
					selName := selectorExpr.Sel.Name

					// TODO arguments
					args := "..."
					step := ""
					pos := fset.Position(stmt.Pos())

					if x.Obj != nil { // method
						rhs := selectorExpr.X.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs
						typ := rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit).Type.(*ast.Ident).Name
						if op := rhs[0].(*ast.UnaryExpr).Op.String(); len(op) != 0 {
							if op == "&" {
								op = "*"
							}
							typ = op + typ
						}
						step = fmt.Sprintf("%s%s%s (%s)%s.%s(%s)",
							log.COLOR_GREEN, pos, log.COLOR_RESET, typ, xName, selName, args)
					} else { // package exported function
						step = fmt.Sprintf("%s%s%s %s.%s(%s)\n",
							log.COLOR_GREEN, pos, log.COLOR_RESET, xName, selName, args)
					}

					if len(step) != 0 {
						steps = append(steps, step)
					}
				}

				methodSteps[method] = steps

				return true
			})
		}
	}

	return methodSteps, nil
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

func renderSteps(method string, steps []string) error {
	fmt.Printf("%s*%s%s", log.COLOR_RED, method, log.COLOR_RESET)
	for _, step := range steps {
		// 递归的使用\v\b进行绘制
		// 串行的使用\v\r进行绘制
		fmt.Printf("\v\r|\v\r|\v\rv\v\r")
		fmt.Printf("%s\v\r", step)
	}
	fmt.Println()
	return nil
}
