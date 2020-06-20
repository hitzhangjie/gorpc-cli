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
	"os/exec"
	"runtime"

	"github.com/hitzhangjie/gorpc/config"
	"github.com/spf13/cobra"
)

// issueCmd represents the bug command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: config.LoadTranslation("issueCmdUsage", nil),
	Long:  config.LoadTranslation("issueCmdUsageLong", nil),
	RunE: func(cmd *cobra.Command, args []string) error {

		err := openIssuesInBrowser()
		if err != nil {
			return err
		}
		return nil
	},
}

func openIssuesInBrowser() error {

	fireAt := "https://github.com/hitzhangjie/gorpc/issues"

	switch runtime.GOOS {
	case "darwin":

		browser := "/Applications/Google Chrome.app"
		_, err := os.Lstat(browser)
		if os.IsNotExist(err) {
			browser = "/Applications/Safari.app"
		}

		cmd := exec.Command("/usr/bin/open", "-a", browser, fireAt)
		buf, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("open browser error: %v\n%s", err, string(buf))
			return err
		}

	default:
		fmt.Printf("请移步 [%s] 提交bug\n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(issueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
