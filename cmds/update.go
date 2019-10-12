package cmds

import (
	"flag"
	"fmt"
	"github.com/hitzhangjie/go-rpc-cmdline/config"
	"github.com/hitzhangjie/go-rpc-cmdline/params"
	"github.com/hitzhangjie/go-rpc-cmdline/parser"
	"github.com/hitzhangjie/go-rpc-cmdline/parser/gomod"
	"github.com/hitzhangjie/go-rpc-cmdline/tpl"
	"github.com/hitzhangjie/go-rpc-cmdline/util/log"
	"github.com/hitzhangjie/go-rpc-cmdline/util/pb"
	"os"
	"path"
	"path/filepath"
)

type UpdateCmd struct {
	Cmd
	*params.Option
}

func newUpdateCmd() *UpdateCmd {

	cmd := Cmd{
		usageLine: `gorpc update`,
		descShort: `
how to update project:
	gorpc update -protodir=. -protofile=*.proto -protocol=whisper
	gorpc update -protofile=*.proto -protocol=whisper`,

		descLong: `
gorpc update:
	-protodir, search path for protofile
	-protofile, protofile to handle
	-lang, language including: go, java, cpp, etc
	-protocol, protocol to use, gorpc, chick or swan`,
	}

	return &UpdateCmd{cmd, params.NewOption()}
}

func (c *UpdateCmd) Run(args ...string) (err error) {

	c.initFlagSet()
	c.parseFlagSet(args)

	// `-protofile=abc/d.proto`, works like `-protodir=abc -protofile=d.proto`ma
	p, err := filepath.Abs(c.Protofile)
	if err != nil {
		panic(err)
	}
	c.Protofile = filepath.Base(p)
	c.Protodirs = append(c.Protodirs, filepath.Dir(p))

	// load language config in gorpc.json
	c.GoRPCConfig, err = config.GetLanguageCfg(c.Language)
	if err != nil {
		return err
	}

	// using assetdir in gorpc.json
	if len(c.Assetdir) == 0 {
		c.Assetdir = c.GoRPCConfig.AssetDir
	}

	// init logging level
	log.InitLogging(c.Verbose)

	return c.update()
}

func (c *UpdateCmd) initFlagSet() {

	fs := flag.NewFlagSet("updatecmd", flag.ContinueOnError)

	fs.Var(&params.RepeatedOption{}, "protodir", "search path of protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "gorpc", "protocol to use, gorpc, chick or swan")
	//fs.Bool("httpon", false, "enable http mode")
	fs.Bool("v", false, "verbose mode")
	fs.String("assetdir", "", "search path of project template")
	fs.Bool("alias", false, "rpcname alias mode")
	//fs.Bool("rpconly", false, "generate rpc stub only")
	fs.String("lang", "go", "language, including go, java, cpp, etc")

	c.flagSet = fs
}

func (c *UpdateCmd) parseFlagSet(args []string) {

	c.flagSet.Parse(args)

	params.LookupFlag(c.flagSet, "protodir", &c.Protodirs)
	params.LookupFlag(c.flagSet, "protofile", &c.Protofile)
	params.LookupFlag(c.flagSet, "lang", &c.Language)
	params.LookupFlag(c.flagSet, "protocol", &c.Protocol)
	params.LookupFlag(c.flagSet, "alias", &c.AliasOn)
	params.LookupFlag(c.flagSet, "assetdir", &c.Assetdir)
	params.LookupFlag(c.flagSet, "v", &c.Verbose)
}

func (c *UpdateCmd) update() error {

	// 检查pb中的导入路径
	dir, err := pb.LocateProtoFile(&c.Protodirs, c.Protofile)
	if err != nil {
		return err
	}
	log.Info("Found protofile:%s in following dir:%v", c.Protofile, dir)

	// 解析pb
	fd, err := parser.ParseProtoFile(c.Option)
	if err != nil {
		return fmt.Errorf("parse protofile:%s error:%v", c.Protofile, err)
	}

	// 解析gomod
	mod, err := gomod.LoadGoMod()
	if err == nil && len(mod) != 0 {
		c.GoMod = mod
	}

	// 代码生成
	fp := path.Join(dir, c.Protofile)

	outputdir := path.Join(os.TempDir(), fd.PackageName)

	err = tpl.GenerateFiles(fd, fp, outputdir, c.Option)

	if err != nil {
		return fmt.Errorf("generate files error:%v", err)
	}

	os.RemoveAll(outputdir)

	log.Info("Generate files success")

	return nil
}