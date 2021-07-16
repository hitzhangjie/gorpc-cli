package plugins

import (
	"os/exec"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/codeblocks/log"
)

type GoImportsPlugin struct {
}

func (m *GoImportsPlugin) Name() string {
	return "goimports"
}

func (m *GoImportsPlugin) Run(fd *descriptor.FileDescriptor, opts *params.Option) error {

	// check if language is go
	if opts.Language != "go" {
		return nil
	}

	// run goimports to format your code
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
