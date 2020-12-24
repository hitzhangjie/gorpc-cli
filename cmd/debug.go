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
)

// debugCmd represents the debug command
//
// **更好的开发测试支持**：
//
// 1.自动编译程序（带-gcflags参数），生成调试符号、禁用内联优化，方便后面调试；
// 2.goast分析服务接口实现（包括请求、响应类型分析）；
// 3.goast分析函数中外部依赖，如db、rpc调用，可以是任意函数调用；
//   - step1：通过配置文件的方式给出一些要分析的规则列表，如常见的分析db、rpc调用；
//            db、redis等操作参数类型相对简单，rpc调用pb结合proto反射很容易可以搞定；
//   - step2：支持任意函数调用；
//            参数类型不确定，需要遍历工程goast分析函数出入参类型，不难可能有点繁琐；
// 4.启动进程，启动debugger server，并将上述接口函数入口点、db、redis、rpc调用等全部添加breakpoints；
// 5.弹出界面，显示服务接口列表，用户选择并填写req、rsp，点击发送请求；
//   在一次测试中，req作为测试输入，rsp用来做测试最后的断言；
// 6.服务受到请求之后，开始执行对应的接口方法
//   接口方法代码执行时，执行路径上的所有像db、redis、rpc调用（或用户指定函数调用），弹出界面询问输入、输出；
//   - 通过debugger server api发起stepin操作，进入函数第一条指令
//   - 通过FDE计算，找到函数退出地址，修改rip直接退出函数（代替真正的执行）
//   - 根据用户输入，更新返回值，如更新rsp、err
// 7.继续往下执行，类似mock测试中涉及到外部依赖需要mock的地方都按照步骤6中的方式来处理
// 8.一条测试执行结束，整个过程中用户输入的内容：
//   - 自动作为测试用例沉淀到用例配置文件中；
//   - 自动生成mock函数、mock代码，不需要每位开发手动编写；
//   - 自动生成测试函数，加载该配置文件，可以持续跑测试；
//
// **有两个明显的好处**：
//
// - 开发阶段，开发测试更方便了：
//   - 一个用例需要mock哪些代码，可能需要多次阅读代码，这次mock掉了db，还漏了rpc，不用关心了，用到哪个mock哪个；
//   - 工具自动化生成mock的代码，也不用关心mock怎么写了，每次写类似代码，太繁琐太耗时；
// - 代码质量，开发真正的关心其测试覆盖率，省掉写接口测试代码、mock逻辑的耗时；
//
// tdd模式，鼓励先规划测试用例，再投入开发，红绿红绿……
// 该issue提及的方式，与tdd也不冲突，只需要一个加载配置文件的helper函数辅助一下即可。
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("debug called")
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// debugCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// debugCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
