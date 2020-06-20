package en

import "github.com/nicksnyder/go-i18n/v2/i18n"

// root cmd
var rootCmdUsage = i18n.Message{
	ID:          "rootCmdUsage",
	Description: "usage of rootCmd",
	Other:       "gorpc is an efficient too to speedup development",
}

var rootCmdUsageLong = i18n.Message{
	ID:          "rootCmdUsageLong",
	Description: "usage of rootCmd",
	Other: `gorpc is an efficient too to speedup development.

for example:
- quickly generate project or rpcstub, based on pb
- send rpc request to test the target server
- update template to the newest version
- quickly open issue page to report an issue
- more ...`,
}

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

var updateCmdUsage = i18n.Message{
	ID:          "updateCmdUsage",
	Description: "usage of updateCmd",
	Other:       "update gorpc template to the newest version",
}

var updateCmdUsageLong = i18n.Message{
	ID:          "updateCmdUsageLong",
	Description: "usage of updateCmd",
	Other: `update gorpc template to the newest version.

by default, go get -u only update the binary, not the template.`,
}
