package en

import "github.com/nicksnyder/go-i18n/v2/i18n"

// create cmd
var createCmdUsage = i18n.Message{
	ID:          "createCmdUsage",
	Description: "usage of createCmd",
	Other:       "quickly create project or rpcstub, based on pb",
}

var createCmdUsageLong = i18n.Message{
	ID:          "createCmdUsageLong",
	Description: "usage of createCmd",
	Other: `quickly create project or rpcstub, based on pb.

'gorpc create' works in 2 modes:
- generate a complete project
- generate only the rpcstub, with option '--rpconly'.`,
}

// issue cmd
var issueCmdUsage = i18n.Message{
	ID:          "issueCmdUsage",
	Description: "usage of issueCmd",
	Other:       "report an issue",
}

var issueCmdUsageLong = i18n.Message{
	ID:          "issueCmdUsageLong",
	Description: "usage of issueCmd",
	Other:       `fire your browser at the issue page to report an issue.`,
}

// rpc cmd
var rpcCmdUsage = i18n.Message{
	ID:          "rpcCmdUsage",
	Description: "usage of rpcCmd",
	Other:       "send rpc request to server",
}

var rpcCmdUsageLong = i18n.Message{
	ID:          "rpcCmdUsageLong",
	Description: "usage of rpcCmd",
	Other: `send rpc request to server.

the request header, request body can be specified via stdin 
or input file, just like the curl style. Besides, target 
server address, network, timeout can be specified, too.`,
}

// version cmd
var versionCmdUsage = i18n.Message{
	ID:          "versionCmdUsage",
	Description: "usage of versionCmd",
	Other:       "show gorpc version (commit hash)",
}

var versionCmdUsageLong = i18n.Message{
	ID:          "versionCmdUsageLong",
	Description: "usage of versionCmd",
	Other:       `show gorpc version (commit hash)`,
}

var versionMsgFormat = i18n.Message{
	ID:          "versionMsgFormat",
	Description: "version msg format",
	Other:       "gorpc version: {{.Hash}}",
}
