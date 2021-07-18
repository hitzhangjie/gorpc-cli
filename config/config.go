package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/codeblocks/tar"

	"github.com/hitzhangjie/gorpc-cli/bindata"
)

// TemplateCfg 开发语言相关的配置信息，如对应的模板工程目录、模板工程中的serverstub文件、clientstub文件
type TemplateCfg struct {
	AssetDir          string   `json:"asset_dir"`            // required: 语言对应的模板工程目录，如asset_go
	LangFileExt       string   `json:"lang_file_ext"`        // required: 文件扩展名，如.go
	TplFileExt        string   `json:"tpl_file_ext"`         // required: 工程中模板文件的后缀名，如.tpl
	RPCServerStub     string   `json:"rpc_server_stub"`      // optional: 工程中对应的rpc server stub文件名（按service.method分文件生成时有用)
	RPCServerImplStub string   `json:"rpc_server_impl_stub"` // optional: 工程中对应的rpc server impl stub文件名（按service.method分文件生成时有用)
	RPCServerTestStub string   `json:"rpc_server_test_stub"` // optional: 工程中对应的rpc server stub测试文件名（按service.method分文件生成时有用)
	RPCClientStub     []string `json:"rpc_client_stub"`      // required: 工程中对应的rpc client stub文件列表
}

var templateCfg *TemplateCfg

func init() {
	// 当前用户对应的模板候选目录列表
	paths, err := TemplateSearchPaths()
	if err != nil {
		panic(err)
	}

	needReinstall := false

	// 确定一个有效的模板路径，如果未安装则安装模板
	installTo, err := TemplateInstallPath(paths)
	if err != nil {
		needReinstall = true
	} else {
		fp := filepath.Join(installTo, "VERSION")
		if !hasSameTplVersion(fp) {
			needReinstall = true
		}
	}

	if needReinstall {
		installTo = paths[0]
		err = InstallTemplate(installTo)
		if err != nil {
			panic(err)
		}
	}

	// 加载配置文件
	initializeConfig(installTo)

	// 加载i18n配置
	initializeI18NMessages(installTo)
}

// LocateTemplatePath locate where templates are installed
func LocateTemplatePath() (string, error) {
	paths, err := TemplateSearchPaths()
	if err != nil {
		panic(err)
	}

	return TemplateInstallPath(paths)
}

func hasSameTplVersion(fp string) bool {

	buf, err := ioutil.ReadFile(fp)
	if err != nil {
		panic(err)
	}
	vals := strings.Split(string(buf), "=")
	if len(vals) != 2 {
		panic("invalid VERSION file")
	}
	if vals[1] != GORPCCliVersion {
		return false
	}
	return true
}

func initializeConfig(installTo string) {
	dat, err := os.ReadFile(filepath.Join(installTo, "gorpc.json"))
	if err != nil {
		panic(err)
	}

	cfg := TemplateCfg{}
	err = json.Unmarshal(dat, &cfg)
	if err != nil {
		panic(err)
	}

	if !path.IsAbs(cfg.AssetDir) {
		cfg.AssetDir = filepath.Join(installTo, cfg.AssetDir)
	}

	if err := validate(&cfg); err != nil {
		panic(err)
	}
	templateCfg = &cfg
}

func InstallTemplate(installTo string) error {
	_ = os.RemoveAll(installTo)

	return tar.Untar(installTo, bytes.NewBuffer(bindata.AssetsGo))
}

// GetTemplateCfg 加载开发语言对应的配置信息
func GetTemplateCfg() (*TemplateCfg, error) {
	return templateCfg, nil
}

var ErrTemplateNotFound = errors.New("template not found")

// TemplateSearchPaths 获取gorpc安装路径
// root安装到/etc/gorpc，非root用户安装到$HOME/.gorpc
func TemplateSearchPaths() (dirs []string, err error) {

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

// TemplateInstallPath 确定一个有效的模板路径
func TemplateInstallPath(dirs []string) (dir string, err error) {
	for _, d := range dirs {
		if fin, err := os.Lstat(d); err == nil && fin.IsDir() {
			return d, nil
		}
	}
	return "", ErrTemplateNotFound
}

func validate(cfg *TemplateCfg) error {
	// lang_file_ext
	if cfg.LangFileExt == "" {
		return errors.New("invalid lang_file_ext, check config 'gorpc.json'")
	}
	// asset dir
	if len(cfg.AssetDir) == 0 {
		return errors.New("invalid asset_dir, check config 'gorpc.json'")
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
