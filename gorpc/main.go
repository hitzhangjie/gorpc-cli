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
	_ "github.com/hitzhangjie/gorpc-cli/bindata"
	"github.com/hitzhangjie/gorpc-cli/cmd"
)

func main() {

	deps, err := loadDependencies()
	if err != nil {
		panic(err)
	}

	err = checkDependencies(deps)
	if err != nil {
		panic(err)
	}

	cmd.Execute()
}
