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

	switch option {
	case "go_package":
		val = opts.GetGoPackage()
	case "java_package":
		val = opts.GetJavaPackage()
	default:
		err = errors.New("fileoption not supported")
		return
	}

	if len(val) == 0 {
		err = fmt.Errorf("%s not defined", option)
		return
	}

	return
}
