package plugins

import (
	"fmt"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
	"github.com/hitzhangjie/gorpc-cli/util/swagger"
)

type SwaggerPlugin struct {
}

func (s *SwaggerPlugin) Name() string {
	return "swagger"
}

func (s *SwaggerPlugin) Run(fd *descriptor.FileDescriptor, opts *params.Option) error {

	if err := swagger.GenSwagger(fd, opts); err != nil {
		return fmt.Errorf("create swagger api document error: %v", err)
	}

	return nil
}
