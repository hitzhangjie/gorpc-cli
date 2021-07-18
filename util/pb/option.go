package pb

import (
	"errors"
	"fmt"

	"github.com/jhump/protoreflect/desc"
)

func GetFileOption(fd *desc.FileDescriptor, option string) (val string, err error) {
	if fd == nil {
		err = errors.New("fd *desc.FileDescriptor is nil")
		return
	}

	opts := fd.GetFileOptions()
	if opts == nil {
		err = errors.New("no fileoption defined")
		return
	}

	if option != "go_package" {
		err = errors.New("fileoption not supported")
		return
	}
	val = opts.GetGoPackage()

	if len(val) == 0 {
		err = fmt.Errorf("%s not defined", option)
		return
	}
	return
}
