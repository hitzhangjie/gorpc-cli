package zh

import "github.com/nicksnyder/go-i18n/v2/i18n"

// root cmd
var rootCmdUsage = i18n.Message{
	ID:          "rootCmdUsage",
	Description: "usage of rootCmd",
	Other:       "gorpc 是一个效率工具，方便gorpc服务的开发",
}

var rootCmdUsageLong = i18n.Message{
	ID:          "rootCmdUsageLong",
	Description: "usage of rootCmd",
	Other: `gorpc 是一个效率工具，方便gorpc服务的开发.

例如:
- 根据pb快速生成服务工程，或者rpcstub
- 发送rpc请求到目的服务用以测试服务接口
- 更新代码模板到最新版本
- 快速打开项目issue页面报告一个issue
- 更多 ...`,
}

// create cmd
var createCmdUsage = i18n.Message{
	ID:          "createCmdUsage",
	Description: "usage of createCmd",
	Other:       "指定pb文件快速创建工程或rpcstub",
}

var createCmdUsageLong = i18n.Message{
	ID:          "createCmdUsageLong",
	Description: "usage of createCmd",
	Other: `指定pb文件快速创建工程或rpcstub，

gorpc create 有两种模式:
- 生成一个完整的服务工程
- 生成被调服务的rpcstub，需指定'--rpconly'选项.`,
}

var createCmdFlagProtodir = i18n.Message{
	ID:          "createCmdFlagProtodir",
	Description: "usage of flag --protodir",
	Other:       `指定pb文件的搜索路径，可指定多个`,
}

var createCmdFlagProtofile = i18n.Message{
	ID:          "createCmdFlagProtofile",
	Description: "usage of flag --protofile",
	Other:       `指定pb文件作为IDL来指导代码生成`,
}

var createCmdFlagProtocol = i18n.Message{
	ID:          "createCmdFlagProtocol",
	Description: "usage of flag --protocol",
	Other:       `指定使用的协议, 如：gorpc, http等`,
}

var createCmdFlagVerbose = i18n.Message{
	ID:          "createCmdFlagVerbose",
	Description: "usage of flag --verbose",
	Other:       `显示详细的运行日志`,
}

var createCmdFlagAssetdir = i18n.Message{
	ID:          "createCmdFlagAssetdir",
	Description: "usage of flag --assetdir",
	Other:       `指定代码模板的路径`,
}

var createCmdFlagRpcOnly = i18n.Message{
	ID:          "createCmdFlagRpcOnly",
	Description: "usage of flag --rpconly",
	Other:       `仅生成rpcstub而非完整工程`,
}

var createCmdFlagLang = i18n.Message{
	ID:          "createCmdFlagLang",
	Description: "usage of flag --lang",
	Other:       `指定使用的编程语言, 如：go, java, python等`,
}

var createCmdFlagMod = i18n.Message{
	ID:          "createCmdFlagMod",
	Description: "usage of flag --mod",
	Other:       `指定go module, 默认是: ${pb.package}`,
}

var createCmdFlagOutput = i18n.Message{
	ID:          "createCmdFlagOutput",
	Description: "usage of flag --output",
	Other:       `指定生成代码的输出路径`,
}

var createCmdFlagForce = i18n.Message{
	ID:          "createCmdFlagForce",
	Description: "usage of flag --force",
	Other:       `强制覆盖已存在的代码`,
}

var createCmdFlagSwagger = i18n.Message{
	ID:          "createCmdFlagSwagger",
	Description: "usage of flag --swagger",
	Other:       `生成swagger api文档`,
}

var createCmdFlagMock = i18n.Message{
	ID: "createCmdFlagMock",
	Description: "usage of flag --mock",
	Other: `mockgen生成接口测试代码`,
}

// issue cmd
var issueCmdUsage = i18n.Message{
	ID:          "issueCmdUsage",
	Description: "usage of issueCmd",
	Other:       "反馈一个issue",
}

var issueCmdUsageLong = i18n.Message{
	ID:          "issueCmdUsageLong",
	Description: "usage of issueCmd",
	Other:       `浏览器打开issue页面，反馈一个issue.`,
}

// rpc cmd
var rpcCmdUsage = i18n.Message{
	ID:          "rpcCmdUsage",
	Description: "usage of rpcCmd",
	Other:       "发送rpc请求给指定服务",
}

var rpcCmdUsageLong = i18n.Message{
	ID:          "rpcCmdUsageLong",
	Description: "usage of rpcCmd",
	Other: `发送rpc请求给指定服务.

可以通过stdin输入协议请求头、请求体，也可以通过输入文件输入，就像curl
一样。当然，也可以指定目的服务地址、网络类型、超时时间等.`,
}

// version cmd
var versionCmdUsage = i18n.Message{
	ID:          "versionCmdUsage",
	Description: "usage of versionCmd",
	Other:       "显示gorpc命令的版本 (commit hash)",
}

var versionCmdUsageLong = i18n.Message{
	ID:          "versionCmdUsageLong",
	Description: "usage of versionCmd",
	Other:       `显示gorpc命令的版本 (commit hash)`,
}

var versionMsgFormat = i18n.Message{
	ID:          "versionMsgFormat",
	Description: "version msg format",
	Other:       "gorpc 命令版本：{{.Hash}}",
}

// update cmd
var updateCmdUsage = i18n.Message{
	ID:          "updateCmdUsage",
	Description: "usage of updateCmd",
	Other:       "指定pb文件更新工程或者rpcstub",
}

var updateCmdUsageLong = i18n.Message{
	ID:          "updateCmdUsageLong",
	Description: "usage of updateCmd",
	Other:       `指定pb文件快速更新工程或者rpcstub，当前未完全实现`,
}
