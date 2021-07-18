package plugins

import (
	"os/exec"

	"github.com/hitzhangjie/codeblocks/log"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
)

// GoImportsPlugin goimports to format your code
type GoImportsPlugin struct {
}

// Name returns plugin's name
func (m *GoImportsPlugin) Name() string {
	return "goimports"
}

// Run run goimports to format your code
func (m *GoImportsPlugin) Run(fd *descriptor.FileDescriptor, opts *params.Option) error {
	goimports, err := exec.LookPath("goimports")
	if err != nil {
		log.Error("please install goimports in order to format code")
		return nil
	}

	// goimports -w -local github.com .
	cmd := exec.Command(goimports, "-w", "-local", "github.com", ".")
	if buf, err := cmd.CombinedOutput(); err != nil {
		log.Error("run goimports error: %v,\n%s", err, string(buf))
		return err
	}

	return nil
}
