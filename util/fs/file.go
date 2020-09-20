package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hitzhangjie/gorpc-cli/util/lang"
)

// BaseNameWithoutExt return basename without extension of `filename`,
// in which `filename` may contains directory.
func BaseNameWithoutExt(filename string) string {
	return lang.TrimRight(".", filepath.Base(filename))
}

// LocateFile 定位protofile目录路径
//
// 要想保证protofile能够被搜索到，protodirs需要提供protofile的父路径，不能是父路径的父路径
func LocateFile(protofile string, protodirs []string) (string, error) {

	// `-protodir not specified` or `-protodir=.`
	if len(protodirs) == 0 || (len(protodirs) == 1 && (protodirs)[0] == ".") {
		abs, _ := filepath.Abs(".")
		return filepath.Join(abs, protofile), nil
	}

	// $protodir/$protofile
	dirs := UniqFilePath(protodirs)

	//查找protofile的绝对路径
	fpaths := []string{}
	for _, dir := range dirs {
		fp := filepath.Join(dir, protofile)
		fin, err := os.Stat(fp)
		if err == nil && !fin.IsDir() {
			fpaths = append(fpaths, fp)
		}
	}
	if len(fpaths) == 0 {
		return "", fmt.Errorf("%s not found in dirs: %v", protofile, protodirs)
	} else if len(fpaths) > 1 {
		return "", fmt.Errorf("%s found duplicate ones: %v", protofile, fpaths)
	}

	// `-protofile=abc/d.proto`, works like `-protodir=abc -protofile=d.proto`ma
	absPath, err := filepath.Abs(fpaths[0])
	if err != nil {
		return "", err
	}
	return absPath, nil
}

// UniqFilePath 文件系统路径去重
func UniqFilePath(dirs []string) []string {

	set := map[string]struct{}{}
	for _, p := range dirs {
		abs, _ := filepath.Abs(p)
		set[abs] = struct{}{}
	}

	uniq := []string{}
	for dir := range set {
		uniq = append(uniq, dir)
	}

	return uniq
}

// PrepareOutputdir create outputdir if it doesn't exist,
// return error if `outputdir` existed while it is not a directory,
// return error if any other error occurs.
func PrepareOutputdir(outputdir string) error {

	var err error

	fin, err := os.Lstat(outputdir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return os.MkdirAll(outputdir, os.ModePerm)
	}

	if !fin.IsDir() {
		return fmt.Errorf("target %s already existed", outputdir)
	}

	return nil
}
