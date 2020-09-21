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
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"git.code.oa.com/go-neat/tencent/polaris/selector"
	"git.code.oa.com/trpc-go/trpc-go/client"
	"git.code.oa.com/trpc-go/trpc-go/codec"

	"github.com/golang/protobuf/proto"
	"github.com/hitzhangjie/gorpc-cli/config"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rpcCmd represents the rpc command
var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: config.LoadTranslation("rpcCmdUsage", nil),
	Long:  config.LoadTranslation("rpcCmdUsageLong", nil),
	RunE: func(cmd *cobra.Command, args []string) error {

		// 检查参数
		err := cmd.ParseFlags(args)
		if err != nil {
			return fmt.Errorf("parse flags error: %v", err)
		}

		// 加载rpc相关选项
		opts, err := loadRPCOptions(cmd.Flags())
		if err != nil {
			return fmt.Errorf("load rpc options error: %v", err)
		}

		// 准备发起rpc请求
		ctx := gorpc.BackgroundContext()
		msg := gorpc.Message(ctx)

		msg.WithClientRPCName(string(opts.reqhead.GetFunc()))
		msg.WithCalleeServiceName(string(opts.reqhead.GetCallee()))

		callopts := []client.Option{
			client.WithProtocol("gorpc"),
			client.WithReqHead(&opts.reqhead),
			client.WithRspHead(&opts.rsphead),
			client.WithTimeout(opts.timeout),
			client.WithTarget(opts.target),
			client.WithNamespace(opts.namespace),
		}

		var pb bool
		if _, ok := opts.request.(proto.Message); ok {
			callopts = append(callopts, client.WithSerializationType(codec.SerializationTypePB))
			pb = true
		} else {
			pb = false
			callopts = append(callopts, client.WithSerializationType(codec.SerializationTypeJSON))
			callopts = append(callopts, client.WithCurrentSerializationType(codec.SerializationTypeNoop))
		}

		var infinite bool
		if opts.times == 0 {
			infinite = true
		}
		fmt.Println()

		for {
			err := client.DefaultClient.Invoke(ctx, opts.request, opts.response, callopts...)
			if pb {
				fmt.Printf("req pb body:%s\nrsp pb body:%s\nreq head:%s\nrsp head:%s\nnode:%s\nerr:%v\n\n",
					opts.request, opts.response, &opts.reqhead, &opts.rsphead, node, err)
			} else {
				fmt.Printf("req json body:%s\nrsp json body:%s\nreq head:%s\nrsp head:%s\nnode:%s\nerr:%v\n\n",
					opts.request.(*codec.Body).Data, opts.response.(*codec.Body).Data, &opts.reqhead, &opts.rsphead, node, err)
			}

			if infinite {
				time.Sleep(opts.interval)
				continue
			}

			opts.times--
			if opts.times == 0 {
				break
			}
			time.Sleep(opts.interval)
		}

		return nil

	},
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rpcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rpcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rpcCmd.Flags().String("head", "", "gorpc request head json data")
	rpcCmd.Flags().String("body", "", "*required, rpc request body json data")
	rpcCmd.Flags().String("headfile", "", "read gorpc request head json data from head file path")
	rpcCmd.Flags().String("bodyfile", "", "read gorpc request body json data from body file path")

	rpcCmd.Flags().String("target", "ip://127.0.0.1:8000", "server target address")
	rpcCmd.Flags().Bool("production", false, "server namespace is production")
	rpcCmd.Flags().Duration("timeout", time.Second, "rpc request timeout")
	rpcCmd.Flags().Duration("interval", time.Second, "rpc request interval")
	rpcCmd.Flags().Uint64("times", 1, "request loop times, 0: infinite loop")

	rpcCmd.Flags().String("protofile", "", "back server proto file path: helloworld.proto")
	rpcCmd.Flags().String("reqbody", "", "rpc request body message name: gorpc.test.helloworld.HelloRequest")
	rpcCmd.Flags().String("rspbody", "", "rpc response body message name: gorpc.test.helloworld.HelloReply")
	rpcCmd.Flags().String("func", "", "rpc method name: /gorpc.test.helloworld.Greeter/SayHello, required")
	rpcCmd.Flags().String("callee", "", "back server service name: gorpc.test.helloworld.Greeter")

	rpcCmd.MarkFlagRequired("protofile")
	rpcCmd.MarkFlagRequired("reqbody")
	rpcCmd.MarkFlagRequired("rspbody")
	rpcCmd.MarkFlagRequired("callee")
	rpcCmd.MarkFlagRequired("callee")
}

type rpcOptions struct {
	head     string
	body     string
	headfile string
	bodyfile string

	target     string
	production bool
	namespace  string
	timeout    time.Duration
	interval   time.Duration
	times      uint64

	protofile string
	reqbody   string
	rspbody   string

	xfunc  string
	callee string

	reqhead gorpc.RequestProtocol
	rsphead gorpc.ResponseProtocol

	request  interface{}
	response interface{}
}

func loadRPCOptions(flagSet *pflag.FlagSet) (*rpcOptions, error) {

	opts := rpcOptions{}

	selector.RegisterDefault()

	opts.protofile, _ = flagSet.GetString("protofile")

	opts.head, _ = flagSet.GetString("head")
	opts.body, _ = flagSet.GetString("body")
	opts.headfile, _ = flagSet.GetString("headfile")
	opts.bodyfile, _ = flagSet.GetString("bodyfile")

	opts.production, _ = flagSet.GetBool("production")
	opts.target, _ = flagSet.GetString("target")
	opts.timeout, _ = flagSet.GetDuration("timeout")
	opts.times, _ = flagSet.GetUint64("times")
	opts.interval, _ = flagSet.GetDuration("interval")

	opts.reqbody, _ = flagSet.GetString("reqbody")
	opts.rspbody, _ = flagSet.GetString("rspbody")
	opts.xfunc, _ = flagSet.GetString("func")
	opts.callee, _ = flagSet.GetString("callee")

	if opts.headfile != "" {
		data, err := ioutil.ReadFile(opts.headfile)
		if err != nil {
			return nil, err
		}
		opts.head = strings.Trim(string(data), " \n\t")
	}
	if opts.head != "" {
		err := codec.Unmarshal(codec.SerializationTypeJSON, []byte(opts.head), &opts.reqhead)
		if err != nil {
			return nil, err
		}
	}
	if opts.xfunc != "" {
		opts.reqhead.Func = []byte(opts.xfunc)
	}
	if len(opts.reqhead.GetFunc()) == 0 {
		return nil, errors.New("func empty")
	}

	if opts.callee != "" {
		opts.reqhead.Callee = []byte(opts.callee)
	}

	if len(opts.reqhead.GetCallee()) == 0 {
		s := strings.Split(string(opts.reqhead.GetFunc()), "/")
		if len(s) > 1 {
			opts.reqhead.Callee = []byte(s[1])
		}
	}

	if opts.bodyfile != "" {
		data, err := ioutil.ReadFile(opts.bodyfile)
		if err != nil {
			return nil, err
		}
		opts.body = strings.Trim(string(data), " \n\t")
	}
	if opts.body == "" {
		return nil, errors.New("request body data empty")
	}

	if opts.production {
		opts.namespace = "Production"
	} else {
		opts.namespace = "Development"
	}

	if opts.protofile != "" {
		p := protoparse.Parser{}
		descs, err := p.ParseFiles(opts.protofile)
		if err != nil {
			panic(err)
		}
		if len(descs) != 1 {
			panic("descs length invalid")
		}
		desc := descs[0]
		reqMsgDesc := desc.FindMessage(opts.reqbody)
		if reqMsgDesc == nil {
			panic("req body name not exist")
		}
		rspMsgDesc := desc.FindMessage(opts.rspbody)
		if rspMsgDesc == nil {
			panic("rsp body name not exist")
		}
		msg := dynamic.NewMessage(reqMsgDesc)
		err = msg.UnmarshalJSON([]byte(opts.body))
		if err != nil {
			panic(err)
		}
		opts.request = msg
		opts.response = dynamic.NewMessage(rspMsgDesc)
	} else {
		opts.request = &codec.Body{Data: []byte(opts.body)}
		opts.response = &codec.Body{}
	}

	return &opts, nil
}
