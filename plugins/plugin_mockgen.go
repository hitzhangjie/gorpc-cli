package plugins

import (
	"os/exec"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/log"
)

type MockgenPlugin struct {
}

func (m *MockgenPlugin) Run(fd *descriptor.FileDescriptor, opts *params.Option) error {

	if opts.Language != "go" {
		return nil
	}

	if _, err := exec.LookPath("mockgen"); err != nil {
		log.Error("please install mockgen in order to generate mockstub")
		return nil
	}

	cmd := exec.Command("go", "generate")
	if buf, err := cmd.CombinedOutput(); err != nil {
		log.Error("run mockgen error: %v,\n%s", err, string(buf))
		return err
	}

	return nil
}
