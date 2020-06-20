package zh

import "github.com/nicksnyder/go-i18n/v2/i18n"

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

	'gorpc create' 有两种模式:
	- 生成一个完整的服务工程
	- 生成被调服务的rpcstub，需指定'--rpconly'选项.`,
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
	Other:       "更新gorpc代码模板到最新版本",
}

var updateCmdUsageLong = i18n.Message{
	ID:          "updateCmdUsageLong",
	Description: "usage of updateCmd",
	Other: `更新gorpc代码模板到最新版本. 

默认地, 通过"go get -u"更新，只是更新了gorpc命令，但没有更新代码模板`,
}
