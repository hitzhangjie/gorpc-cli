package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/hitzhangjie/gorpc/bindata"
	"github.com/hitzhangjie/gorpc/util/compress"
	"github.com/hitzhangjie/gorpc/util/fs"
)

// LanguageCfg 开发语言相关的配置信息，如对应的模板工程目录、模板工程中的serverstub文件、clientstub文件
type LanguageCfg struct {
	Language          string   `json:"language"` // required: 语言名称，如go、java
	LangFileExt       string   `json:"lang_file_ext"`
	AssetDir          string   `json:"asset_dir"`            // required: 语言对应的工程目录
	TplFileExt        string   `json:"tpl_file_ext"`         // required: 工程中模板文件的后缀名，如.tpl
	RPCServerStub     string   `json:"rpc_server_stub"`      // optional: 工程中对应的rpc server stub文件名（按service.method分文件生成时有用)
	RPCServerImplStub string   `json:"rpc_server_impl_stub"` // optional: 工程中对应的rpc server impl stub文件名（按service.method分文件生成时有用)
	RPCServerTestStub string   `json:"rpc_server_test_stub"` // optional: 工程中对应的rpc server stub测试文件名（按service.method分文件生成时有用)
	RPCClientStub     []string `json:"rpc_client_stub"`      // required: 工程中对应的rpc client stub文件列表
}

// configs 所有语言的配置信息，汇总在此
var configs = map[string]*LanguageCfg{}

func init() {

	// 当前用户对应的模板候选目录列表
	paths, err := templateSearchPath()
	if err != nil {
		panic(err)
	}

	// 确定一个有效的模板路径，如果未安装则安装模板
	installTo, err := templateInstallPath(paths)
	if err == ErrTemplateNotFound {
		installTo = paths[0]
		err = installTemplate(installTo)
		if err != nil {
			panic(err)
		}
	}

	// 加载配置文件
	fin, err := os.Open(filepath.Join(installTo, "gorpc.json"))
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadAll(fin)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(dat, &configs)
	if err != nil {
		panic(err)
	}

	for k, v := range configs {
		if err := validate(k, v); err != nil {
			panic(err)
		}
	}
}

func installTemplate(installTo string) error {

	tmpDir := filepath.Join(os.TempDir(), "gorpc")

	_ = os.RemoveAll(installTo)
	_ = os.RemoveAll(tmpDir)

	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = compress.Untar(tmpDir, bytes.NewBuffer(bindata.InstallTgzBytes))
	if err != nil {
		return err
	}

	err = fs.Move(filepath.Join(tmpDir, "install"), installTo)
	if err != nil {
		return err
	}

	return nil
}

// GetLanguageCfg 加载开发语言对应的配置信息
func GetLanguageCfg(lang string) (*LanguageCfg, error) {
	cfg, ok := configs[lang]
	if !ok {
		return nil, fmt.Errorf("language:%s not supported, check config 'gorpc.json'", lang)
	}
	return cfg, nil
}

var ErrTemplateNotFound = errors.New("template not found")

// templateSearchPath 获取gorpc安装路径
// root安装到/etc/gorpc，非root用户安装到$HOME/.gorpc
func templateSearchPath() (dirs []string, err error) {

	u, err := user.Current()
	if err != nil {
		return
	}

	candidateDirs := []string{filepath.Join(u.HomeDir, ".gorpc"), "/etc/gorpc"}
	if u.Username == "root" {
		candidateDirs = []string{"/etc/gorpc"}
	}

	return candidateDirs, nil
}

// templateInstallPath 确定一个有效的模板路径
func templateInstallPath(dirs []string) (dir string, err error) {
	for _, d := range dirs {
		if fin, err := os.Lstat(d); err == nil && fin.IsDir() {
			return d, nil
		}
	}
	return "", ErrTemplateNotFound
}

func validate(lang string, cfg *LanguageCfg) error {

	dirs, err := templateSearchPath()
	if err != nil {
		return err
	}

	dir, err := templateInstallPath(dirs)
	if err != nil {
		return err
	}

	if len(lang) == 0 {
		return errors.New("invalid language, check config 'gorpc.json'")
	}
	cfg.Language = lang
	if cfg.LangFileExt == "" {
		cfg.LangFileExt = lang
	}
	// asset dir
	if len(cfg.AssetDir) == 0 {
		return errors.New("invalid asset_dir, check config 'gorpc.json'")
	}
	if !path.IsAbs(cfg.AssetDir) {
		cfg.AssetDir = filepath.Join(dir, cfg.AssetDir)
	}
	// tpl_file_ext
	if len(cfg.TplFileExt) == 0 {
		return errors.New("invalid tpl_file_ext, check config 'gorpc.json'")
	}
	// rpc_server_stub, 分文件用，不设置也ok

	// rpc_client_stub，-rpconly用
	if len(cfg.RPCClientStub) == 0 {
		return errors.New("invalid rpc_client_stub, check config 'gorpc.json'")
	}
	return nil
}
