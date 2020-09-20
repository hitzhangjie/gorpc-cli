package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/fs"
)

// getOutputDir return outputdir
//
// default outputdir is equal to basename of protofile, if option `-o` specified, use its value as outputdir
func getOutputDir(options *params.Option) (string, error) {

	// 当前路径下存在go.mod，并且未指定-mod/-o，则直接使用当前路径
	// -mod不影响输出目录，-mod指定时也不参考，参考-o或者pb文件名
	if len(options.GoModEx) != 0 &&
		options.GoModEx == options.GoMod && // -mod 二者相同，相当于没指定
		len(options.OutputDir) == 0 { // -o
		return os.Getwd()
	}

	// 指定了-o的情况下，直接使用
	if len(options.OutputDir) != 0 {
		return options.OutputDir, nil
	}

	// 其他情况下使用pbfile的basename
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}

	b := fs.BaseNameWithoutExt(options.Protofile)
	return filepath.Join(d, b), nil
}

// changeProtofileDir 不同语言希望pb存放位置不同，hardcode
func changeProtofileDir(pbDir, language string) error {
	var err error
	switch language {
	case "java":
		var pbGenDir string
		err = filepath.Walk(pbDir, func(fpath string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				proPath := strings.Replace(pbDir, "/stub", "", -1)
				return fs.Copy(fpath, filepath.Join(proPath, info.Name())) // java路径调整 suggest from youngwwang
			} else {
				if !strings.HasSuffix(fpath, "com") {
					return nil
				}
				pbGenDir = fpath
				return nil
			}
		})
		// 移出pb生成的.java文件，删除pb生成的/com/github/gorpc/greeter/...文件夹
		if pbGenDir != "" {
			err = os.RemoveAll(pbGenDir)
		}
	default:
	}
	return err
}
