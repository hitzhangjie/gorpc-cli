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
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
	"github.com/hitzhangjie/gorpc-cli/util/debug"
	"github.com/phayes/freeport"
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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// 获取要分析的微服务目录
		pkg, _ := cmd.Flags().GetString("package")
		if len(pkg) == 0 {
			return errors.New("package empty")
		}

		if pkg == "." {
			pkg, err = os.Getwd()
			if err != nil {
				return err
			}
		}
		fmt.Printf("----------------------------------\n")
		fmt.Printf("debug `%s`\n", pkg)

		// 准备进行go代码编译
		bin, err := buildDebugVersion(pkg)
		if err != nil {
			return err
		}
		fmt.Printf("build `%s`\n", bin)

		port, err := freeport.GetFreePort()
		if err != nil {
			return err
		}
		fmt.Printf("allocate port: %d\n", port)

		// start debugger in headerless mode
		addr := fmt.Sprintf("127.0.0.1:%d", port)
		pid, err := debugInHeadlessMode(bin, addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "start debugger error: %v\n", err)
			return
		}
		fmt.Printf("----------------------------------\n")
		fmt.Printf("start debugger ok, tracee: %d\n", pid)

		// 等debugger server起来
		//for !detectAddrInUse(addr) {
		//	time.Sleep(time.Millisecond*100)
		//}
		time.Sleep(time.Second * 2)

		// 初始化rpc client，并开始设置好断点
		rpc := rpc2.NewClient(addr)
		if rpc == nil {
			return errors.New("start debugger client error")
		}
		rpc.SetReturnValuesLoadConfig(&apiLoadConfig)

		// set breakpoints at statements of RPC
		bp, err := rpc.CreateBreakpoint(&api.Breakpoint{
			//File: "/Users/zhangjie/gorpc101/gorpc-cli/testcase/testcase.debug/main.go",
			//Line: 60,
			FunctionName: "main.(*Client).Invoke",
			//Cond:          "",
		})
		if err != nil {
			return err
		}
		fmt.Printf("create breakpoint ok, pos: %#x, %s\n", bp.Addr, bp.FunctionName)
		fmt.Printf("create breakpoint ok, pos: %s:%d\n", bp.File, bp.Line)

		// state, err := rpc.StepInstruction()
		// BUG: thread blocked in function prologue, using rpc.Step() instead
		state, err := rpc.Step()
		if err != nil {
			return err
		}
		fmt.Printf("stepin into the function: pc: %#x\n", state.CurrentThread.PC)

		// TODO 这里应该通过FDE计算得到返回值，然后将其设置到rip，然后直接stepin
		// 现在是图省事，走了stepout，函数体还是被完整执行了，不想完整执行，比如：
		// 希望绕过真正的网络io的ctx超时部分，这部分可能比较耗时间.
		state, err = rpc.StepOut()
		if err != nil {
			return err
		}

		tab, err := debug.BuildLineTable(bin)
		if err != nil {
			return fmt.Errorf("build lntab error: %v", err)
		}
		f, l, fn := tab.PCToLine(state.CurrentThread.PC)

		// note: main.go:30 -> main.go:60 -> main.go:30
		fmt.Printf("stepout the function: pos: %#x %s\n", state.CurrentThread.PC, fn.Name)
		fmt.Printf("stepout the function: pos: %s:%d\n", f, l)

		scope := api.EvalScope{
			GoroutineID:  state.SelectedGoroutine.ID,
			Frame:        0,
			DeferredCall: 0,
		}
		//err = printLocalVariables(rpc, scope)
		//if err != nil {
		//	return err
		//}
		rpc.Next()

		// TODO rsp可能为nil，没想到什么好办法，delve不支持
		// 可以用修改桩代码的方式，始终返回一个非指针类型的结构体
		//err0 := rpc.SetVariable(scope, "rsp", "&helloRsp{}")
		//err1 := rpc.SetVariable(scope, "rsp.code", "1024")
		//err2 := rpc.SetVariable(scope, "rsp.msg", `"hello, world"`)
		//err3 := rpc.SetVariable(scope, "err", "nil")
		// can not call function with nil ReturnInfoLoadConfig
		//_, err2 := rpc.Call(state.SelectedGoroutine.ID, `rsp.msg="cool"`, false)
		//_, err3 := rpc.Call(state.SelectedGoroutine.ID, `err=errors.New("xxxxx")`, false)
		//v, err0 := rpc.EvalVariable(scope, "rsp", apiLoadConfig)
		//fmt.Printf("var %+v, error: %v\n", v, err0)

		fmt.Printf("----------------------------------\n")
		fmt.Printf("check the values before mocked\n")
		// 显示rsp.code
		v, err := rpc.EvalVariable(scope, "rsp.code", apiLoadConfig)
		if err != nil {
			return err
		}
		fmt.Printf("var rsp.code = %s, error: %v\n", v.Value, err)

		// 显示rsp.msg
		v, err = rpc.EvalVariable(scope, "rsp.msg", apiLoadConfig)
		if err != nil {
			return err
		}
		fmt.Printf("var rsp.msg = %s, error: %v\n", v.Value, err)

		// 使用输入的值mock这里的变量值
		fmt.Printf("----------------------------------\n")
		fmt.Printf("mock the return values using input\n")
		err = rpc.SetVariable(scope, "rsp.code", "9999")
		if err != nil {
			return err
		}
		v, err = rpc.EvalVariable(scope, "rsp.code", apiLoadConfig)
		if err != nil {
			return err
		}
		fmt.Printf("var rsp.code = %s\n", v.Value)

		return nil
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
	debugCmd.Flags().StringP("package", "p", ".", "package to debug") // TODO add i18n translations
}

func buildDebugVersion(pkg string) (string, error) {
	bin := filepath.Join(filepath.Dir(pkg), "_debug_"+filepath.Base(pkg))

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	err = os.Chdir(pkg)
	if err != nil {
		return "", err
	}
	defer os.Chdir(wd)

	// note: don't use exec.Command("go", "build", `-gcflags="all=-N -l"`, "-o", bin)
	// only shell consumes the double quotes "all and -l", so drop this in exec.Command.
	// see: https://github.com/golang/go/issues/42482.
	cmd := exec.Command("go", "build", "-gcflags=all=-N -l", "-o", bin)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error: %v, details: %s", err, string(buf))
	}
	return bin, nil
}

func debugInHeadlessMode(bin string, addr string) (pid int, err error) {
	cmd := exec.Command("dlv", "exec", bin,
		"--headless", "--api-version=2", "--log", fmt.Sprintf("--listen=%s", addr))
	err = cmd.Start()
	if err != nil {
		return -1, err
	}
	go cmd.Wait()
	return cmd.Process.Pid, nil
}

// TODO not working as expected
func detectAddrInUse(addr string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

func printLocalVariables(rpc *rpc2.RPCClient, scope api.EvalScope) error {
	vars, err := rpc.ListLocalVariables(scope, api.LoadConfig{
		FollowPointers:     true,
		MaxVariableRecurse: 8,
		MaxStringLen:       1024,
		MaxArrayValues:     1024,
		MaxStructFields:    64,
	})
	if err != nil {
		return fmt.Errorf("local variables error: %v", err)
	}
	for i, v := range vars {
		fmt.Printf("%d-%#x %s %s = %s\n", i, v.Addr, v.Name, v.RealType, v.Value)
	}
	return nil
}

var apiLoadConfig = api.LoadConfig{
	FollowPointers:     true,
	MaxVariableRecurse: 8,
	MaxStringLen:       256,
	MaxArrayValues:     64,
	MaxStructFields:    16,
}
