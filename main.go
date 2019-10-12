package main

import (
	"github.com/hitzhangjie/go-rpc-cmdline/cmds"
	"github.com/hitzhangjie/go-rpc-cmdline/log"
	"os"
)

func main() {

	help, ok := cmds.SubCmd("help")
	if !ok {
		panic("gorpc <help> subcmd not registered")
	}

	if l := len(os.Args); l == 1 {
		help.Run()
		return
	}

	cmd, ok := cmds.SubCmd(os.Args[1])
	if !ok || cmd == nil {
		help.Run()
		return
	}

	if err := cmd.Run(os.Args[2:]...); err != nil {
		log.Error("Run command:%v error:\n\t\t%v", os.Args, err)
	}
}
