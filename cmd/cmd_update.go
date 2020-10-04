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

	"github.com/hitzhangjie/gorpc-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)

	// TODO use i18n to generate translations for English/Chinese
	updateCmd.Flags().StringP("protofile", "p", "", "specify the protofile to process")
	updateCmd.Flags().StringP("projectdir", "d", "", "specify the directory of existed project")

	updateCmd.MarkFlagRequired("protofile")
	updateCmd.MarkFlagRequired("projectdir")
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: config.LoadTranslation("updateCmdUsage", nil),
	Long:  config.LoadTranslation("updateCmdUsageLong", nil),
	RunE: func(cmd *cobra.Command, args []string) error {

		p, _ := cmd.Flags().GetString("protofile")
		d, _ := cmd.Flags().GetString("projectdir")

		fmt.Println("protofile:", p)
		fmt.Println("projectdir:", d)

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

				buf := bytes.Buffer{}
				err := format.Node(&buf, fset, fn)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", buf.String())
				return true
			})
		}
	}

	return nil
}
