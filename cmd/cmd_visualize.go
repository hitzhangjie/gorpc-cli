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
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	FuncMain            = "main"
	FuncRegisterService = "RegisterHelloSvrService"
	FuncRegisterPattern = "Register.*Service"
)

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// TODO use i18n to generate translations for English/Chinese
	visualizeCmd.Flags().StringP("projectdir", "d", "", "specify the directory of existed project")

	visualizeCmd.MarkFlagRequired("projectdir")
}

// updateCmd represents the update command
var visualizeCmd = &cobra.Command{
	Use: "visualize",
	// TODO i18n settings
	Short: "visualize workflow",
	Long:  `visualize workflow`,
	RunE: func(cmd *cobra.Command, args []string) error {

		d, _ := cmd.Flags().GetString("projectdir")

		return parse(d)
	},
}

func parse(dir string) error {

	dirs := []string{dir}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	allPkgs := map[string]*ast.Package{}

	for _, dir := range dirs {
		pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for k, v := range pkgs {
			allPkgs[k] = v
		}
	}

	for _, pkg := range allPkgs {
		for fname, file := range pkg.Files {
			fmt.Printf("------------- %s --------------\n", fname)
			ast.Inspect(file, func(n ast.Node) bool {
				// perform analysis here
				// ...

				fn, ok := n.(*ast.FuncDecl)
				if !ok {
					return true
				}
				switch fn.Name.Name {
				case FuncMain:
					buf := bytes.Buffer{}
					err := format.Node(&buf, fset, fn)
					if err != nil {
						panic(err)
					}
					fmt.Printf("%s\n", buf.String())
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
						if selectorExpr.Sel.Name == FuncRegisterService {
							fmt.Println("func:", FuncRegisterService, "found")
							service := callExpr.Args[1].(*ast.UnaryExpr).X.(*ast.CompositeLit).Type.(*ast.Ident).Name
							fmt.Println("service:", service)
						}
					}

				case FuncRegisterService:
					buf := bytes.Buffer{}
					err := format.Node(&buf, fset, fn)
					if err != nil {
						panic(err)
					}
					fmt.Printf("%s\n", buf.String())
				default:
					return true
				}

				return true
			})
		}
	}

	return nil
}
