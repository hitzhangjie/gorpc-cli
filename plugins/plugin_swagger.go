package plugins

import (
	"fmt"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/swagger"
)

type SwaggerPlugin struct {
}

func (s *SwaggerPlugin) Run(fd *descriptor.FileDescriptor, opts *params.Option) error {

	if opts.Language != "go" {
		return nil
	}

	if opts.SwaggerOn {
		if err := swagger.GenSwagger(fd, opts); err != nil {
			return fmt.Errorf("create swagger api document error: %v", err)
		}
	}

	return nil
}