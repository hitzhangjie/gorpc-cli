package cmds

import (
	"sync"
)

var (
	cmds = map[string]Commander{}
	mux  sync.Mutex
)

// SubCmds return cmds registered subcmds.
func SubCmd(subcmd string) (Commander, bool) {
	mux.Lock()
	defer mux.Unlock()

	cmd, ok := cmds[subcmd]
	return cmd, ok
}

func RegisterSubCmd(name string, commander Commander) {
	mux.Lock()
	cmds[name] = commander
	mux.Unlock()
}

func init() {
	RegisterSubCmd("create", newCreateCmd())
	RegisterSubCmd("update", newUpdateCmd())
	RegisterSubCmd("help", newHelpCmd())
}
