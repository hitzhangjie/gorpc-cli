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
	"os"

	"github.com/hitzhangjie/gorpc-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: config.LoadTranslation("updateCmdUsage", nil),
	Long:  config.LoadTranslation("updateCmdUsageLong", nil),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 当前用户对应的模板候选目录列表
		paths, err := config.TemplateSearchPaths()
		if err != nil {
			panic(err)
		}

		// 确定一个有效的模板路径，如果未安装则安装模板
		installTo, err := config.TemplateInstallPath(paths)

		// 已经安装模板，则先删除模板
		if err == nil {
			if err = os.RemoveAll(installTo); err != nil {
				return err
			}
		}

		// 重新安装模板
		installTo = paths[0]
		if err = config.InstallTemplate(installTo); err != nil {
			return err
		}
		return nil
	},
}
