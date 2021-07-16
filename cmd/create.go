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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/gorpc-cli/config"
	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/parser"
	"github.com/hitzhangjie/gorpc-cli/plugins"
	"github.com/hitzhangjie/gorpc-cli/tpl"
	"github.com/hitzhangjie/gorpc-cli/util/lang"
	"github.com/hitzhangjie/gorpc-cli/util/pb"

	"github.com/hitzhangjie/codeblocks/format"
	"github.com/hitzhangjie/codeblocks/fs"
	"github.com/hitzhangjie/codeblocks/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	createOption     *params.Option
	createSuccess    bool
	createOutputDir  string
	createDescriptor *descriptor.FileDescriptor
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringArray("protodir", []string{"."}, config.LoadTranslation("createCmdFlagProtodir", nil))
	createCmd.Flags().StringP("protofile", "p", "", config.LoadTranslation("createCmdFlagProtofile", nil))
	createCmd.Flags().String("protocol", "gorpc", config.LoadTranslation("createCmdFlagProtocol", nil))
	createCmd.Flags().BoolP("verbose", "v", false, config.LoadTranslation("createCmdFlagVerbose", nil))
	createCmd.Flags().String("assetdir", "", config.LoadTranslation("createCmdFlagAssetdir", nil))
	createCmd.Flags().Bool("rpconly", false, config.LoadTranslation("createCmdFlagRpcOnly", nil))
	createCmd.Flags().String("lang", "go", config.LoadTranslation("createCmdFlagLang", nil))
	createCmd.Flags().StringP("mod", "m", "", config.LoadTranslation("createCmdFlagMod", nil))
	createCmd.Flags().StringP("output", "o", "", config.LoadTranslation("createCmdFlagOutput", nil))
	createCmd.Flags().BoolP("force", "f", false, config.LoadTranslation("createCmdFlagForce", nil))

	// plugins
	createCmd.Flags().String("plugins", "goimports", config.LoadTranslation("createCmdFlagPlugins", nil))

	createCmd.MarkFlagRequired("protofile")
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: config.LoadTranslation("createCmdUsage", nil),
	Long:  config.LoadTranslation("createCmdUsageLong", nil),
	PreRunE: func(cmd *cobra.Command, args []string) error {

		// 检查命令行参数
		option, err := loadCreateOption(cmd.Flags())
		if err != nil {
			return fmt.Errorf("check flags error: %v", err)
		}
		createOption = option

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		// 初始化日志级别
		if createOption.Verbose {
			log.SetFlags(log.LVerbose)
		}
		log.Info("ready to process protofile: %s", createOption.ProtofileAbs)

		// 解析pb
		fd, err := parser.ParseProtoFile(createOption)
		if err != nil {
			return fmt.Errorf("Parse protofile: %s error: %v", createOption.Protofile, err)
		}
		fd.FilePath = createOption.ProtofileAbs
		fd.Dump()

		// 创建工程
		var outputdir string
		if !createOption.RpcOnly {
			outputdir, err = create(fd, createOption)
		} else {
			outputdir, err = generateRPCStub(fd, createOption)
		}

		if err != nil {
			if !createOption.RpcOnly {
				return fmt.Errorf("create gorpc project error: %v", err)
			}
			return fmt.Errorf("create gorpc stub error: %v", err)
		}

		createOption = createOption
		createOutputDir = outputdir
		createSuccess = true
		createDescriptor = fd

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {

		value, _ := cmd.Flags().GetString("plugins")
		enabled := map[string]bool{}
		for _, v := range strings.Split(value, "+") {
			enabled[v] = true
		}

		if !createSuccess {
			return nil
		}

		err := os.Chdir(createOutputDir)
		if err != nil {
			return err
		}

		for _, p := range plugins.Plugins {
			if !enabled[p.Name()] {
				continue
			}
			err := p.Run(createDescriptor, createOption)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func loadCreateOption(flagSet *pflag.FlagSet) (*params.Option, error) {

	option := loadCreateFlagsetToOption(flagSet)

	// 检查pb文件是否合法
	if len(option.Protofile) == 0 {
		return nil, errors.New("invalid protofile")
	}

	// 定位pb文件
	target, err := fs.LocateFile(option.Protofile, option.Protodirs)
	if err != nil {
		return nil, err
	}
	option.Protofile = filepath.Base(target)
	option.ProtofileAbs = target
	option.Protodirs = append(option.Protodirs, filepath.Dir(target))

	// 加载gorpc.json中定义的语言相关的配置
	option.GoRPCCfg, err = config.GetLanguageCfg(option.Language)
	if err != nil {
		return nil, fmt.Errorf("load config via gorpc.json error: %v", err)
	}
	if len(option.Assetdir) == 0 {
		option.Assetdir = option.GoRPCCfg.AssetDir
	}

	// 判断gomod
	// - 优先使用-mod指定的moduleName
	// - 没有指定-mod选项的话，再考虑加载本地go.mod，兼容老的操作逻辑
	// - 如果本地也没有指定go.mod，再考虑pb中的package（模板实现的）
	if len(option.GoMod) == 0 {
		mod, err := lang.LoadGoMod()
		if err == nil && len(mod) != 0 {
			option.GoModEx = mod
			option.GoMod = mod
		}
	}

	return option, nil
}

func loadCreateFlagsetToOption(flagSet *pflag.FlagSet) *params.Option {

	option := &params.Option{}

	option.Protodirs, _ = flagSet.GetStringArray("protodir")
	option.Protofile, _ = flagSet.GetString("protofile")
	option.Language, _ = flagSet.GetString("lang")
	option.Protocol, _ = flagSet.GetString("protocol")
	option.RpcOnly, _ = flagSet.GetBool("rpconly")
	option.Assetdir, _ = flagSet.GetString("assetdir")
	option.Verbose, _ = flagSet.GetBool("verbose")
	option.GoMod, _ = flagSet.GetString("mod")
	option.OutputDir, _ = flagSet.GetString("output")
	option.Force, _ = flagSet.GetBool("force")

	return option
}

// create 代码生成，生成完整的工程
func create(fd *descriptor.FileDescriptor, option *params.Option) (outputdir string, err error) {

	// - 准备输出目录
	outputdir, err = getOutputDir(option)
	if err != nil {
		return
	}

	if !isSafeOutputDir(outputdir) && !option.Force {
		err = fmt.Errorf("reject overwrite existed code: %s,\nuse --force/-f to make it if you want", outputdir)
		return
	}

	defer func() {
		if err != nil {
			removeDirAsNeeded(outputdir)
		}
	}()

	// - 生成代码
	err = tpl.GenerateFiles(fd, option, outputdir)
	if err != nil {
		return
	}

	// create rpcstub
	stubDir := filepath.Join(outputdir, "stub")
	if _, err = os.Lstat(stubDir); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(stubDir, os.ModePerm); err != nil {
			return
		}
	}
	stub := filepath.Join(outputdir, "stub")

	// - move outputdir/rpc to outputdir/stub/dir($gopkgdir)
	fileOption := fmt.Sprintf("%s_package", option.GoRPCCfg.Language)
	pbPackage, err := parser.GetPbPackage(fd, fileOption)
	if err != nil {
		return
	}

	if fileOption == "java_package" {
		pathLast := path.Join(strings.Split(pbPackage, ".")...)
		pbPackage = path.Join("client/src/main/java", strings.ToLower(pathLast))
	} else if fileOption == "python_package" {
		pbPackage = strings.Replace(pbPackage, ".", "_", -1)
	}

	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	pbOutDir := filepath.Join(stub, pbPackage)
	err = os.MkdirAll(pbOutDir, os.ModePerm)
	if err != nil {
		return
	}

	pb2pkg := fd.Pb2ImportPath

	// 处理-protofile指定的pb文件
	err = pb.Protoc(fd, option.Protodirs, option.Protofile, option.Language, pbOutDir, pb2pkg)
	if err != nil {
		err = fmt.Errorf("GenerateFiles: %v", err)
		return
	}

	// - copy *.proto to outpoutdir/rpc/
	basename := filepath.Base(fd.FilePath)
	err = fs.Copy(fd.FilePath, filepath.Join(pbOutDir, basename))
	if err != nil {
		return
	}

	// - 处理${protofile}依赖的其他pb文件
	//BUG: 目录组织问题，不再按照pb相对路径关系进行组织，全部按照stub/package进行组织
	//err = handleDependencies(fd, option, pbPackage, pbOutDir)
	err = handleDependencies(fd, option, pbPackage, stub)
	if err != nil {
		return
	}

	// - 将outputdir/rpc移动到outputdir/$gopkgdir/
	src := filepath.Join(outputdir, "rpc")
	defer os.RemoveAll(src)
	dest := path.Join(stub, pbPackage)

	err = filepath.Walk(src, func(fpath string, info os.FileInfo, err error) (e error) {

		if fpath == src {
			return nil
		}

		if fname := filepath.Base(fpath); fname == "gorpc.go" {
			// - 将stub文件gorpc.go重命名，
			fname = fs.BaseNameWithoutExt(basename)
			return fs.Move(fpath, filepath.Join(dest, fname+".gorpc.go"))
		} else {
			return fs.Move(fpath, filepath.Join(dest, filepath.Base(fpath)))
		}
	})
	if err != nil {
		return
	}

	// 格式化操作
	if err = format.FormatDir(outputdir); err != nil {
		return
	}

	log.Info("generate project %s```%s```%s success", log.COLOR_RED, basename, log.COLOR_GREEN)
	return
}

func generateRPCStub(fd *descriptor.FileDescriptor, option *params.Option) (outputdir string, err error) {

	// 代码生成
	// - 准备输出目录
	outputdir, err = os.Getwd()
	if err != nil {
		return
	}

	if filepath.IsAbs(option.OutputDir) {
		outputdir = option.OutputDir
	} else {
		outputdir = filepath.Join(outputdir, option.OutputDir)
	}
	err = os.MkdirAll(outputdir, os.ModePerm)
	if err != nil {
		return
	}

	// - 生成代码，只处理clientstub
	generated := map[string]struct{}{}
	for _, f := range option.GoRPCCfg.RPCClientStub {
		in := filepath.Join(option.Assetdir, f)
		log.Debug("handle:%s", in)
		out := filepath.Join(outputdir, strings.TrimSuffix(filepath.Base(in), option.GoRPCCfg.TplFileExt))
		if err = tpl.GenerateFile(fd, in, out, option); err != nil {
			return
		}

		if err = format.Format(out); err != nil {
			return
		}

		generated[out] = struct{}{}
	}

	// 将stub文件gorpc.go重命名
	basename := fs.BaseNameWithoutExt(fd.FilePath)

	for fpath := range generated {
		if fname := filepath.Base(fpath); fname != "gorpc.go" {
			continue
		}
		dst := filepath.Join(path.Dir(fpath), basename+".gorpc.go")
		if err = fs.Move(fpath, dst); err != nil {
			return
		}
		break
	}

	//将所有package相同的依赖过滤出来
	var protofiles []string
	protofiles = append(protofiles, option.Protofile)
	for fname, pkg := range fd.Pb2ValidGoPkg {
		if pkg == fd.PackageName {
			protofiles = append(protofiles, fname)
		}
	}
	// - generate *.pb.go or *.java or *.pb.h/*.pb.cc under outputdir/rpc/
	if err = pb.Protoc(fd, option.Protodirs, option.Protofile, option.Language, outputdir, fd.Pb2ImportPath); err != nil {
		//if err = pb.Protoc(c.Option.Protodirs, c.Option.Protofile, c.Option.Language, outputdir, fd.Pb2ValidGoPkg); err != nil {
		err = fmt.Errorf("GenerateFiles: %v", err)
		return
	}
	log.Info("generate rpc stub success")
	return
}

func isSafeOutputDir(dir string) bool {

	_, err := os.Lstat(dir)

	// 目录不存在，说明不存在写覆盖的情况
	if err != nil {
		if os.IsNotExist(err) {
			return true
		}
		return false
	}

	// 存在的话，检测下目录下是否存在源文件，存在就有覆盖风险
	err = filepath.Walk(dir, func(p string, inf os.FileInfo, err error) error {
		if strings.HasSuffix(p, ".go") {
			return fmt.Errorf("go code detected: %s", p)
		}
		return nil
	})

	if err != nil {
		return false
	}
	return true
}

func removeDirAsNeeded(path string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if dir == path {
		return
	}
	os.RemoveAll(path)
}

// handleDependencies 处理-protofile指定的pb文件中import的其他pb文件，包括protoc处理，与拷贝pb文件
//
// 准备protoc生成pb文件对应的*.pb.go，需要注意的是，避免生成循环依赖的代码
//
// jhump/protoreflect解析结果，如果是与-protofile相同的pkgname，那么importpath为"",
//
// protoc --go_out=M$pb=$pkgname，这里需要做兼容处理:
// 	1. 避免传递$pkgname为空, 否则protoc会生成这样的代码：
//    ```go
//    package $pkgname
//    import (
//        "."
//    )
//    ```
// 	2. 避免传递与-protofile相同pkgname的情况，不然会导致循环依赖:
//    ```go
//    package $pkgname
//    import (
//        "$pkgname"
//    )
//    ```
func handleDependencies(fd *descriptor.FileDescriptor, option *params.Option, pbPackage string, outputDir string) (err error) {
	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return err
	}

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return err
	}

	includeDirs := []string{}
	for fname := range fd.Pb2ImportPath {
		dir, _ := filepath.Split(fname)
		includeDirs = append(includeDirs, dir)
	}

	// 逐一处理依赖的其他pb文件
	for fname, importPath := range fd.Pb2ImportPath {

		// 如果是${protofile}跳过不处理
		if fname == fd.FilePath {
			continue
		}

		// 跳过google官方提供的pb文件，gorpc扩展文件，swagger 扩展文件
		if strings.HasPrefix(fname, "google/protobuf") || fname == "gorpc.proto" || fname == "swagger.proto" {
			continue
		}

		pbOutDir := filepath.Join(outputDir, importPath)
		if option.Language == "java" {
			pbOutDir = filepath.Join(outputDir, pbPackage)
		}
		if err := os.MkdirAll(pbOutDir, os.ModePerm); err != nil {
			return err
		}

		// 继承上一级的目录,避免出现目录找不到的问题
		searchPath := option.Protodirs

		parentDirs := []string{wd}
		parentDirs = append(parentDirs, option.Protodirs...)

		for _, pDir := range parentDirs {
			for _, incDir := range includeDirs {

				includeDir := filepath.Join(pDir, incDir)
				includeDir = filepath.Clean(includeDir)

				if fin, err := os.Lstat(includeDir); err != nil {
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					if !fin.IsDir() {
						return fmt.Errorf("import path: %s, not directory", includeDir)
					}
					searchPath = append(searchPath, includeDir)
				}
			}
		}

		if err := pb.Protoc(fd, searchPath, fname, option.Language, pbOutDir, fd.Pb2ImportPath); err != nil {
			return fmt.Errorf("GenerateFiles: %v", err)
		}

		// 拷贝pb文件
		p, err := fs.LocateFile(fname, option.Protodirs)
		if err != nil {
			return err
		}

		_, baseName := filepath.Split(fname)
		src := p
		dst := filepath.Join(pbOutDir, baseName)

		log.Debug("copy file %s to %s", src, dst)
		if err := fs.Copy(src, dst); err != nil {
			return err
		}

		// 初始化gomod
		//
		// 避免重复初始化go.mod
		fp := filepath.Join(pbOutDir, "go.mod")
		fin, err := os.Stat(fp)
		if err == nil && !fin.IsDir() {
			continue
		}

		// TODO 移动到createCmd.PostRun

		// 执行go mod init, 与pbPackage相同也不用初始化
		if option.Language != "go" {
			continue
		}

		if len(importPath) != 0 && importPath != pbPackage {
			os.Chdir(pbOutDir)

			cmd := exec.Command("go", "mod", "init", importPath)
			if buf, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("process %s, initialize go.mod in stub/%s error: %v", fname, importPath, string(buf))
			}
			log.Debug("process %s, initialize go.mod success in xxxout/%s: go mod init %s", fname, importPath, importPath)
		}
	}

	if err = os.Chdir(wd); err != nil {
		return err
	}

	return nil
}
