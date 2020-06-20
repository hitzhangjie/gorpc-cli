// +build experimental

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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/hitzhangjie/gorpc/descriptor"
	"github.com/hitzhangjie/gorpc/params"
	"github.com/hitzhangjie/gorpc/parser"
	"github.com/hitzhangjie/gorpc/tpl"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update your project or rpc stub",
	Long:  `update your project or rpc stub`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// 解析参数
		err := cmd.ParseFlags(args)
		if err != nil {
			return fmt.Errorf("parse flags error: %v", err)
		}

		// 加载选项
		option, err := loadUpdateOptions(cmd.Flags())
		if err != nil {
			return fmt.Errorf("load options error: %v", err)
		}

		// 执行更新
		fd, err := parser.ParseProtoFile(option)
		if err != nil {
			return fmt.Errorf("parse protofile:%s error:%v", option.Protofile, err)
		}
		fd.FilePath = option.ProtofileAbs
		fd.Dump()

		err = update(fd, option)
		if err != nil {
			return fmt.Errorf("update project error: %v", err)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().StringArray("protodir", []string{"."}, "include path of the target protofile")
	updateCmd.Flags().StringP("protofile", "p", "", "protofile used as IDL of target service")
	updateCmd.Flags().String("protocol", "gorpc", "protocol to use, gorpc, http, etc")
	updateCmd.Flags().BoolP("verbose", "v", false, "show verbose logging info")
	updateCmd.Flags().String("assetdir", "", "path of project template")
	updateCmd.Flags().Bool("alias", false, "enable alias mode of rpc name")
	updateCmd.Flags().Bool("rpconly", false, "generate rpc stub only")
	updateCmd.Flags().String("lang", "go", "programming language, including go, java, python")
	updateCmd.Flags().StringP("mod", "m", "", "go module, default: ${pb.package}")
	updateCmd.Flags().StringP("output", "o", "", "output directory")
	updateCmd.Flags().BoolP("force", "f", false, "enable overwritten existed code forcibly")
}

func loadUpdateOptions(flagSet *pflag.FlagSet) (*params.Option, error) {
	return loadCreateOption(flagSet)
}

func update(fd *descriptor.FileDescriptor, option *params.Option) error {

	// 代码生成
	outputdir := filepath.Join(os.TempDir(), fd.PackageName)

	err := tpl.GenerateFiles(fd, option, outputdir)

	if err != nil {
		return fmt.Errorf("generate files error:%v", err)
	}

	os.RemoveAll(outputdir)

	return nil
}
