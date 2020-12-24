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

	"github.com/spf13/cobra"

	"github.com/hitzhangjie/gorpc-cli/config"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: config.LoadTranslation("versionCmdUsage", nil),
	Long:  config.LoadTranslation("versionCmdUsageLong", nil),
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]interface{}{
			"Hash": config.GORPCCliVersion,
		}
		trans := config.LoadTranslation("versionMsgFormat", data)
		fmt.Println(trans)
	},
}
