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

var createCmdFlagProtodir = i18n.Message{
	ID:          "createCmdFlagProtodir",
	Description: "usage of flag --protodir",
	Other:       `include path of the target protofile`,
}

var createCmdFlagProtofile = i18n.Message{
	ID:          "createCmdFlagProtofile",
	Description: "usage of flag --protofile",
	Other:       `protofile used as IDL of target service`,
}

var createCmdFlagProtocol = i18n.Message{
	ID:          "createCmdFlagProtocol",
	Description: "usage of flag --protocol",
	Other:       `protocol to use, gorpc, http, etc`,
}

var createCmdFlagVerbose = i18n.Message{
	ID:          "createCmdFlagVerbose",
	Description: "usage of flag --verbose",
	Other:       `show verbose logging info`,
}

var createCmdFlagAssetdir = i18n.Message{
	ID:          "createCmdFlagAssetdir",
	Description: "usage of flag --assetdir",
	Other:       `path of project template`,
}

var createCmdFlagRpcOnly = i18n.Message{
	ID:          "createCmdFlagRpcOnly",
	Description: "usage of flag --rpconly",
	Other:       `generate rpc stub only`,
}

var createCmdFlagLang = i18n.Message{
	ID:          "createCmdFlagLang",
	Description: "usage of flag --lang",
	Other:       `programming language, including go, java, python`,
}

var createCmdFlagMod = i18n.Message{
	ID:          "createCmdFlagMod",
	Description: "usage of flag --mod",
	Other:       `go module, default: ${pb.package}`,
}

var createCmdFlagOutput = i18n.Message{
	ID:          "createCmdFlagOutput",
	Description: "usage of flag --output",
	Other:       `output directory`,
}

var createCmdFlagForce = i18n.Message{
	ID:          "createCmdFlagForce",
	Description: "usage of flag --force",
	Other:       `enable overwritten existed code forcibly`,
}

var createCmdFlagPlugins = i18n.Message{
	ID:          "createCmdFlagPlugins",
	Description: "gorpc create --plugins, enable plugins",
	Other:       "enabled plugins list, joined by '+', support goimports, mock, swagger",
}

// issue cmd
var bugCmdUsage = i18n.Message{
	ID:          "bugCmdUsage",
	Description: "usage of issueCmd",
	Other:       "report an issue",
}

var bugCmdUsageLong = i18n.Message{
	ID:          "bugCmdUsageLong",
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
or input file, just like the curl gostyle. Besides, target 
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

// update cmd
var updateCmdUsage = i18n.Message{
	ID:          "updateCmdUsage",
	Description: "usage of updateCmd",
	Other:       "quickly update gorpc template or rpcstub, based on pb",
}

var updateCmdUsageLong = i18n.Message{
	ID:          "updateCmdUsageLong",
	Description: "usage of updateCmd",
	Other: `quickly update gorpc template or rpcstub, based on pb.

so far, this feature hasn't beed fully implemented, coming soon.`,
}

