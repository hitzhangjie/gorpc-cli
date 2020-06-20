package zh

import "github.com/nicksnyder/go-i18n/v2/i18n"

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
