package cmd

import (
	"os"
	"path/filepath"

	"github.com/hitzhangjie/gorpc-cli/params"

	"github.com/hitzhangjie/fs"
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
