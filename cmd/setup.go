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
	"github.com/hitzhangjie/codeblocks/log"
	"github.com/spf13/cobra"

	"github.com/hitzhangjie/gorpc-cli/config"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:           "setup",
	Short:         "初始化设置 && 安装依赖工具",
	Long:          `初始化设置 && 安装依赖工具.`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("初始化设置 && 安装依赖工具")

		deps, err := config.LoadDependencies()
		if err != nil {
			return err
		}

		if len(deps) == 0 {
			log.Info("设置完成")
			return nil
		}

		err = config.CheckDependencies(deps)
		if err != nil {
			return err
		}

		log.Info("设置完成")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
