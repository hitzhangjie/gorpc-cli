package tpl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hitzhangjie/gorpc/descriptor"
	"github.com/hitzhangjie/gorpc/params"
	"github.com/hitzhangjie/gorpc/util/fs"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/hitzhangjie/gorpc/util/log"
)

// GenerateFiles 处理go模板文件，并输出到目录outputdir中
func GenerateFiles(fd *descriptor.FileDescriptor, option *params.Option, outputdir string) (err error) {

	// 准备输出目录
	if err := fs.PrepareOutputdir(outputdir); err != nil {
		return fmt.Errorf("GenerateFiles prepareOutputdir:%v", err)
	}

	// 遍历模板文件进行处理
	f := func(path string, info os.FileInfo, err error) error {
		return processTemplateFile(path, info, err, &mixedOptions{fd, option, outputdir /*serviceIdx,*/})
	}

	err = filepath.Walk(option.Assetdir, f)
	if err != nil {
		return fmt.Errorf("GenerateFiles filepath.Walk:%v", err)
	}

	return nil
}

func GenerateFile(fd *descriptor.FileDescriptor, infile, outfile string, option *params.Option, serviceIndex ...int) (err error) {

	assetdir := option.Assetdir
	if !filepath.IsAbs(assetdir) {
		return errors.New("assetdir must be specified an absolute path")
	}

	// stat template
	tplFilePath := infile
	if _, err = os.Stat(tplFilePath); err != nil {
		log.Error("%v", err)
		return err
	}

	// create output file
	fout, err := os.Create(outfile)
	if err != nil {
		log.Error("%v", err)
		return err
	}
	defer fout.Close()

	// template execute and populate the output file
	var tplInstance *template.Template

	baseName := filepath.Base(infile)

	if funcMap == nil {
		tplInstance, err = template.New(baseName).ParseFiles(tplFilePath)
	} else {
		tplInstance, err = template.New(baseName).Funcs(funcMap).ParseFiles(tplFilePath)
	}
	if err != nil {
		log.Error("%v", err)
		return err
	}

	// 将需要的descriptor信息、命令行控制参数信息、其他分文件需要的serviceIndex信息传入
	err = tplInstance.Execute(fout, struct {
		*descriptor.FileDescriptor
		*params.Option
		ServiceIndex int
	}{
		fd,
		option,
		func() int {
			if len(serviceIndex) != 0 {
				return serviceIndex[0]
			}
			return 999999
		}(),
	})
	if err != nil {
		log.Error("%v", err)
		return err
	}

	return nil
}

// processTemplateFile 处理模板文件
func processTemplateFile(entry string, info os.FileInfo, err error, options *mixedOptions) error {

	// if incoming error encounter, return at once
	if err != nil {
		return err
	}

	fd := options.FileDescriptor
	option := options.Option
	outputdir := options.OutputDir
	//serviceIdx := options.ServiceIdx

	// keep same files/folders hierarchy in the outputdir/assetdir
	relativePath := strings.TrimPrefix(entry, option.Assetdir)
	if len(relativePath) == 0 {
		return nil
	}
	relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))

	log.Debug("file entry srcPath:%s", entry)
	// 如果server stub需要分文件，则指定rpc_server_stub模板文件名
	if relativePath == option.GoRPCCfg.RPCServerStub {
		outPath := filepath.Join(outputdir, relativePath)
		dir := filepath.Dir(outPath)
		for idx, sd := range fd.Services {
			base := strcase.ToSnake(sd.Name) + "." + option.GoRPCCfg.LangFileExt
			switch option.GoRPCCfg.Language {
			case "java":
				base = strcase.ToCamel(sd.Name) + "." + option.GoRPCCfg.Language
			}
			outPath = filepath.Join(dir, base)
			if err := GenerateFile(fd, entry, outPath, option, idx); err != nil {
				return err
			}
			continue
		}
		return nil
	}
	// 如果server stub需要分文件，则指定rpc_server_impl_stub模板文件名
	if relativePath == option.GoRPCCfg.RPCServerImplStub {
		outPath := filepath.Join(outputdir, relativePath)
		dir := filepath.Dir(outPath)
		for idx, sd := range fd.Services {
			base := strings.ToLower(sd.Name) + "." + option.GoRPCCfg.LangFileExt
			switch option.GoRPCCfg.Language {
			case "java":
				base = strcase.ToCamel(sd.Name) + "Impl." + option.GoRPCCfg.LangFileExt
			}
			outPath = filepath.Join(dir, base)
			if err := GenerateFile(fd, entry, outPath, option, idx); err != nil {
				return err
			}
			continue
		}
		return nil
	}
	// {{service}}.go 测试文件{{service}}_test.go
	if relativePath == option.GoRPCCfg.RPCServerTestStub {
		outPath := filepath.Join(outputdir, relativePath)
		dir := filepath.Dir(outPath)
		for idx, sd := range fd.Services {
			base := strcase.ToSnake(sd.Name) + "_test." + option.GoRPCCfg.LangFileExt
			outPath = filepath.Join(dir, base)
			if err := GenerateFile(fd, entry, outPath, option, idx); err != nil {
				return err
			}
			log.Debug("entry destPath: %s", outPath)
			continue
		}
		return nil
	}

	outPath := filepath.Join(outputdir, relativePath)
	log.Debug("file entry destPath: %s", outPath)

	// if `entry` is directory, create the same entry in `outputdir`
	if info.IsDir() {
		return os.MkdirAll(outPath, os.ModePerm)
	}
	outPath = strings.TrimSuffix(outPath, option.GoRPCCfg.TplFileExt)

	return GenerateFile(fd, entry, outPath, option)
}

// mixedOptions 将众多选项聚集在一起，简化方法签名
type mixedOptions struct {
	*descriptor.FileDescriptor
	*params.Option
	OutputDir string
	//ServiceIdx int
}
