package plugins

import (
	"os/exec"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/log"
)

type GoImportsPlugin struct {
}

func (m *GoImportsPlugin) Run(fd *descriptor.FileDescriptor, opts *params.Option) error {

	// run goimports to format your code
	goimports, err := exec.LookPath("goimports")
	if err != nil {
		log.Error("please install goimports in order to format code")
		return nil
	}

	cmd := exec.Command(goimports, "-w", ".")
	if buf, err := cmd.CombinedOutput(); err != nil {
		log.Error("run goimports error: %v,\n%s", err, string(buf))
		return err
	}

	return nil
}
