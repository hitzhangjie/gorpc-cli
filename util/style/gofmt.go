package style

import (
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/log"
)

// Format 原地格式化代码
func Format(fpath string, opt *params.Option) error {
	switch opt.Language {
	case "go":
		return GoFmt(fpath)
	default:
		// 不支持
		return nil
	}
}

// GoFmt 原地格式化go代码
func GoFmt(fpath string) error {

	in, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}

	out, err := format.Source(in)
	if err != nil {
		log.Error("%v", err)
		return err
	}

	err = ioutil.WriteFile(fpath, out, 0644)
	if err != nil {
		return err
	}

	return nil
}

// GoFmtDir 原地格式化go代码目录
func GoFmtDir(dir string) error {
	err := filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(fpath, ".go") && !info.IsDir() {
			err := GoFmt(fpath)
			if err != nil {
				log.Error("Warn: style file:%s error:%v", fpath, err)
			}
		}
		return nil
	})
	return err
}
