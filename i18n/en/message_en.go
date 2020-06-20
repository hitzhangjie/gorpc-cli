package en

import "github.com/nicksnyder/go-i18n/v2/i18n"

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
