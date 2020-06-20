/*
Copyright © 2020 zhijiezhang

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
package main

import (
	"os"
	"os/exec"

	"github.com/hitzhangjie/gorpc-cli/cmd"
	"github.com/hitzhangjie/gorpc-cli/util/log"
)

func main() {

	// 检查protoc有没有安装
	if _, err := exec.LookPath("protoc"); err != nil {
		log.Error("Please install protoc ... error:\n\t==> %v", err)
		os.Exit(1)
	}

	// 检查protoc-gen-go有没有安装
	if _, err := exec.LookPath("protoc-gen-go"); err != nil {
		log.Error("Please install protoc-gen-go ... error:\n\t==> %v", err)
		os.Exit(1)
	}

	// 执行命令
	cmd.Execute()
}
