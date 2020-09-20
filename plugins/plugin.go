package plugins

import (
	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/params"
)

var (
	Plugins = []Plugin{
		&SwaggerPlugin{},
		&MockgenPlugin{},
		&GoImportsPlugin{},
	}
)

// Plugin 插件接口
type Plugin interface {
	Run(fd *descriptor.FileDescriptor, opts *params.Option) error
}
