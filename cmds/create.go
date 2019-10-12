package cmds

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hitzhangjie/go-rpc-cmdline/config"
	"github.com/hitzhangjie/go-rpc-cmdline/params"
	"github.com/hitzhangjie/go-rpc-cmdline/parser"
	"github.com/hitzhangjie/go-rpc-cmdline/parser/gomod"
	"github.com/hitzhangjie/go-rpc-cmdline/tpl"
	"github.com/hitzhangjie/go-rpc-cmdline/util/fs"
	"github.com/hitzhangjie/go-rpc-cmdline/util/log"
	"github.com/hitzhangjie/go-rpc-cmdline/util/pb"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// CreateCmd generate project or generate rpcstub(-rpconly)
//
// 1. create project:
// 		gorpc create -protodir=<dir1> -protodir=<dir2> -protofile=greeter.proto -protocol=gorpc -v
//		or
//		gorpc create -protofile=greeter.proto
// 2. create rpcstub:
// 		gorpc create -protofile=greeter.proto -rpconly
type CreateCmd struct {
	Cmd
	*params.Option
}

// newCreateCmd build CreateCmd
func newCreateCmd() *CreateCmd {

	cmd := Cmd{
		usageLine: `gorpc create`,
		descShort: `
how to create project:
	gorpc create -protodir=. -protofile=*.proto -protocol=gorpc -alias
	gorpc create -protofile=*.proto -protocol=gorpc`,
		descLong: `
gorpc create:
	-protodir, search path for protofile, default: "."
	-protofile, protofile to handle
	-protocol, protocol to use, including: gorpc, nrpc, ilive, sso, default: gorpc 
	-lang, language including: go, java, cpp, default: go
	-alias, enable alias mode, //@alias=${rpcName}, default: false
	-rpconly, generate rpc stub only, default: false"`,
	}

	return &CreateCmd{cmd, params.NewOption()}
}

// Run execute the CreateCmd logic
func (c *CreateCmd) Run(args ...string) (err error) {

	// `gorpc create`, parse the arguments
	c.initFlagSet()
	c.parseFlagSet(args)

	// `-protofile=abc/d.proto`, works like `-protodir=abc -protofile=d.proto`
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

	// pass `-assetdir` to gorpc to use customized template dir, instead of the one specified in gorpc.json
	if len(c.Assetdir) == 0 {
		c.Assetdir = c.GoRPCConfig.AssetDir
	}

	// init logging level
	log.InitLogging(c.Verbose)

	// if `-rpconly` specified, then only rpc stub need generated
	if c.RpcOnly {
		return c.generateRPCStub()
	}
	return c.create()
}

// initFlagSet build flagset of CreateCmd
func (c *CreateCmd) initFlagSet() {

	fs := flag.NewFlagSet("createcmd", flag.ContinueOnError)

	fs.Var(&params.RepeatedOption{}, "protodir", "search path of protofile")
	fs.String("protofile", "any.proto", "protofile to handle")
	fs.String("protocol", "gorpc", "protocol to use, gorpc, chick or swan")
	fs.Bool("v", false, "verbose mode")
	fs.String("assetdir", "", "search path of project template")
	fs.Bool("alias", false, "rpcname alias mode")
	fs.Bool("rpconly", false, "generate rpc stub only")
	fs.String("lang", "go", "language, including go, java, cpp, etc")

	c.flagSet = fs
}

// parseFlagSet parse flagset of CreateCmd
func (c *CreateCmd) parseFlagSet(args []string) {

	c.flagSet.Parse(args)

	params.LookupFlag(c.flagSet, "protodir", &c.Protodirs)
	params.LookupFlag(c.flagSet, "protofile", &c.Protofile)
	params.LookupFlag(c.flagSet, "lang", &c.Language)
	params.LookupFlag(c.flagSet, "protocol", &c.Protocol)
	params.LookupFlag(c.flagSet, "alias", &c.AliasOn)
	params.LookupFlag(c.flagSet, "rpconly", &c.RpcOnly)
	params.LookupFlag(c.flagSet, "assetdir", &c.Assetdir)
	params.LookupFlag(c.flagSet, "v", &c.Verbose)
}

func (c *CreateCmd) create() error {

	// locate the protofile `-protofile` in search path `-protodir`, returns the abs path
	dir, err := pb.LocateProtoFile(&c.Protodirs, c.Protofile)
	if err != nil {
		return err
	}
	log.Info("Found protofile:%s in following dir:%v", c.Protofile, dir)

	// 解析pb
	fd, err := parser.ParseProtoFile(c.Option)
	if err != nil {
		return fmt.Errorf("Parse protofile:%s error:%v", c.Protofile, err)
	}

	// 解析gomod
	mod, err := gomod.LoadGoMod()
	if err == nil && len(mod) != 0 {
		c.GoMod = mod
	}
	dump(fd)

	// 代码生成
	// - 准备输出目录
	outputdir, err := getOutputDir(fd, c.Option)
	if err != nil {
		return err
	}
	// - 生成代码
	protofileAbsPath := path.Join(dir, c.Protofile)

	err = tpl.GenerateFiles(fd, protofileAbsPath, outputdir, c.Option)
	if err != nil {
		os.RemoveAll(outputdir)
		return err
	}
	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	pbOutDir := path.Join(outputdir, "rpc")
	if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, pbOutDir, fd.Dependencies); err != nil {
		return fmt.Errorf("GenerateFiles: %v", err)
	}
	// - copy *.proto to outpoutdir/rpc/
	basename := path.Base(protofileAbsPath)
	if err := fs.Copy(protofileAbsPath, path.Join(pbOutDir, basename)); err != nil {
		return err
	}
	// - move outputdir/rpc to outputdir/dir($gopkgdir)
	src := path.Join(outputdir, "rpc")
	fileOption := fmt.Sprintf("%s_package", c.Option.GoRPCConfig.Language)
	gopkgdir := fd.PackageName
	if fo := fd.FileOptions[fileOption]; fo != nil {
		if v := fd.FileOptions[fileOption].(string); len(v) != 0 {
			gopkgdir = v
		}
	}

	// - 将outputdir/rpc移动到outputdir/$gopkgdir/
	dest := path.Join(outputdir, gopkgdir)
	if err := os.MkdirAll(path.Dir(dest), os.ModePerm); err != nil {
		return err
	}
	if err := fs.Move(src, dest); err != nil {
		return err
	}
	// - 将stub文件重命名，加上前缀service
	sd := fd.Services[0]
	err = filepath.Walk(dest, func(fpath string, info os.FileInfo, err error) error {
		dir, fname := filepath.Split(fpath)

		for _, f := range c.Option.GoRPCConfig.RPCClientStub {
			if strings.HasSuffix(fname, c.Option.GoRPCConfig.Language) &&
				strings.HasSuffix(f, fname+c.Option.GoRPCConfig.TplFileExt) {
				fs.Move(fpath, path.Join(dir, sd.Name+"."+fname))
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 格式化操作
	err = filepath.Walk(outputdir, func(fpath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(fpath, ".go") && !info.IsDir() {
			err := gofmt(fpath)
			if err != nil {
				log.Error("Warn: gofmt file:%s error:%v", fpath, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Info("Generate project %s```%s```%s success", log.COLOR_RED, sd.Name, log.COLOR_GREEN)
	return nil
}

func (c *CreateCmd) generateRPCStub() error {

	// 检查pb中的导入路径
	fpaths, err := pb.LocateProtoFile(&c.Protodirs, c.Protofile)
	if err != nil {
		return err
	}
	log.Info("Found protofile:%s in following dir:%v", c.Protofile, fpaths)

	// 解析pb
	fd, err := parser.ParseProtoFile(c.Option)
	if err != nil {
		return fmt.Errorf("Parse protofile:%s error:%v", c.Protofile, err)
	}

	// 解析gomod
	mod, err := gomod.LoadGoMod()
	if err == nil && len(mod) != 0 {
		c.GoMod = mod
	}
	dump(fd)

	// 代码生成
	// - 准备输出目录
	outputdir, err := os.Getwd()
	if err != nil {
		return err
	}
	// - 生成代码，只处理clientstub
	for _, f := range c.Option.GoRPCConfig.RPCClientStub {
		in := path.Join(c.Assetdir, f)
		log.Debug("handle:%s", in)
		out := path.Join(outputdir, strings.TrimSuffix(path.Base(in), c.GoRPCConfig.TplFileExt))
		if err := tpl.GenerateFile(fd, in, out, c.Option); err != nil {
			return err
		}
	}
	// 将stub文件gorpc.go重命名
	// fixme, handle .gorpc.go
	sd := fd.Services[0]
	err = filepath.Walk(outputdir, func(fpath string, info os.FileInfo, err error) error {
		if fname := path.Base(fpath); fname == "gorpc.go" {
			fs.Move(fpath, path.Join(path.Dir(fpath), sd.Name+".gorpc.go"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	//if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, outputdir, fd.Dependencies); err != nil {
	if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, outputdir, fd.ImportPathMappings); err != nil {
		return fmt.Errorf("GenerateFiles: %v", err)
	}
	log.Info("Generate rpc stub success")
	return nil
}

func dump(fd *parser.FileDescriptor) {
	log.Debug("************************** FileDescriptor ***********************")
	buf, _ := json.MarshalIndent(fd, "", "  ")
	log.Debug("\n%s", string(buf))
	log.Debug("*****************************************************************")
}

// purgeNonRpcStub 清理非rpcstub文件
//
// todo 不同语言将${lang}_package转换成文件系统路径时，处理方式不同，先硬编码几个
func purgeNonRpcStub(fd *parser.FileDescriptor, outputdir string, option *params.Option) error {

	// 先简单粗暴搞一发
	pkg, err := pkgFileOption(fd, option.Language)
	if err != nil {
		return err
	}

	// copy outputdir/rpc/ to tmpdir/$pkg/
	src := path.Join(outputdir, pkg)
	tmpdir := path.Join(os.TempDir(), fd.PackageName)
	if err := fs.Copy(src, tmpdir); err != nil {
		return err
	}

	// delete any file or directory under outputdir
	err = fs.DeleteFilesUnderDir(outputdir)
	if err != nil {
		return err
	}

	// copy any file under tmpdir/$pkg/ into outputdir
	err = fs.CopyFileUnderDir(tmpdir, outputdir)
	if err != nil {
		return err
	}

	// remove tmpdir
	if err := os.RemoveAll(tmpdir); err != nil {
		return err
	}

	return nil
}

// pkgFileOption 获取pb文件中的FileOption:${lang}_package
func pkgFileOption(fd *parser.FileDescriptor, lang string) (string, error) {
	pkg := fd.PackageName
	if strings.ToLower(lang) != "cpp" {
		fo := fmt.Sprintf("%s_package", lang)
		opt, ok := fd.FileOptions[fo]
		if !ok || opt == nil {
			return "", fmt.Errorf("invalid FileOption:%s", fo)
		}
		if v, ok := opt.(string); !ok || len(pkg) == 0 {
			return "", fmt.Errorf("invalid FileOption:%s", fo)
		} else {
			pkg = v
		}
	}
	return pkg, nil
}
