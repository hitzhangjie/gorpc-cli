package params

import (
	"github.com/hitzhangjie/gorpc-cli/config"
)

type Option struct {
	// pb option
	Protodirs    []string // pb import路径
	Protofile    string   // protofile文件
	ProtofileAbs string   // protofile绝对路径
	AliasOn      bool     // 解析MethodOption或者注释中//@alias=别名，用来代替pb文件中rpc

	// template option
	Assetdir string // 服务模板路径
	Language string // 开发语言，如go，java，cpp等
	Protocol string // 协议类型
	//HttpOn   bool   // 生成http相关代码，使用-protocol=http代替-httpon
	RpcOnly bool // 只生成rpc相关代码，而非完整工程
	// gorpc.json
	GoRPCCfg *config.LanguageCfg

	// gomod option
	GoMod   string // 当前工程指定的gomod
	GoModEx string // go.mod中提取出的module

	// logging option
	Verbose bool // 输出verbose日志信息

	// logging option
	OutputDir string // 项目输出路径
	Force     bool   // 强制写入

	// swagger option
	SwaggerOn bool // 解析 MethodOption 的swagger
}

var option = &Option{}

func NewOption() *Option {
	return &Option{
		Protodirs: []string{},
	}
}
