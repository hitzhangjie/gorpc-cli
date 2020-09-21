package pb

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/gorpc-cli/descriptor"
	"github.com/hitzhangjie/gorpc-cli/util/lang"
	"github.com/hitzhangjie/gorpc-cli/util/log"
)

// Protoc process `protofile` to generate *.pb.go or *.java, which is specified by `language`
//
// Please note protoc whose version is less than 3.6.0 has bug:
//
// File does not reside within any path specified using --proto_path (or -I).
// You must specify a --proto_path which encompasses this file.
// Note that the proto_path must be an exact prefix of the .proto file names
// -- protoc is too dumb to figure out when two paths (e.g. absolute and relative)
// are equivalent (it's harder than you think).
//
// 当指定如下选项时，protoc仍然无法处理，这其实是个明显的bug，protoc v3.6.0+及以上版本都可以正常执行：
// protoc --proto_path=/root/test --go_out=paths=source_relative:/root/test/greeter/rpc greeter.proto
// or
// protoc --proto_path=. --go_out=paths=source_relative:/root/test/greeter/rpc greeter.proto
//
// 下面存在一些"排除--proto_path为当前路径的操作"，纯粹是为了兼容老的protoc处理相对路径、绝对路径的bug
func Protoc(fd *descriptor.FileDescriptor, protodirs []string, protofile, language, outputdir string, pbpkgMapping map[string]string) error {

	_, baseName := filepath.Split(protofile)

	// locate gorpc.proto
	_, dep := pbpkgMapping["gorpc.proto"]
	if dep {
		p, err := LocateGoRPCProto()
		if err != nil {
			return err
		}
		protodirs = append(protodirs, p)
	}

	args := []string{}

	// make --proto_path
	argsProtoPath, err := makeProtoPath(protodirs)
	if err != nil {
		return err
	}

	// make --go_out
	argsGoOut := makeGoOut(fd, pbpkgMapping, protofile, language, outputdir)

	args = append(args, argsProtoPath...)
	args = append(args, argsGoOut, baseName)
	log.Debug("protoc %s", strings.Join(args, " "))

	cmd := exec.Command("protoc", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Run command: `%s`, error: %s", strings.Join(cmd.Args, " "), string(output))
	}

	return nil
}

func makeGoOut(fd *descriptor.FileDescriptor, pbpkgMapping map[string]string, protofile string, language string, outputdir string) string {

	//protofileValidGoPkg := pbpkgMapping[protofile]

	var out string
	var pbpkg string

	if len(pbpkgMapping) != 0 {
		for k, v := range pbpkgMapping {

			// 1. google官方库就交给protoc、protoc-gen-go自行处理好了
			// 2. 其他import的pb，如与protofile validGoPkg相同，那么protoreflect/jhump解析出pb对应的package为空
			//    BUG: https://github.com/hitzhangjie/gorpc-cli/issues/96 解决这里的循环依赖问题！
			if strings.HasPrefix(k, "google/protobuf") || len(v) == 0 {
				continue
			}

			// 如果指定go_package会简化很多问题，google/protobuf建议为pb文件添加go_package这个fileoption，并将在后续
			// 新版protoc-gen-go中将其作为一个强制的约束。
			// 经测试，上述ISSUE在为pb文件指定go_package的情况下可以完美解决，现在为了更好的兼容性做下处理，针对没有指定
			// go_package这个fileoption的情况也保证生成代码的正确性。
			hasGoPackageOpt := true

			rfd := fd.RawFileDescriptor()
			deps := rfd.GetDependencies()
			for _, dep := range deps {
				depName := dep.GetFile().GetName()
				if k != depName {
					continue
				}
				fmt.Println("k==", k, "depName==", depName)

				_, err := GetFileOption(dep.GetFile(), "go_package")
				if err != nil {
					hasGoPackageOpt = false
				}
			}

			if hasGoPackageOpt {
				pbpkg += ",M" + k + "=" + v
			} else {
				pbpkg += ",M" + k + "=" + lang.PBValidGoPackage(v)
			}
		}
	}

	switch language {
	case "go", "python":
		out = fmt.Sprintf("--%s_out=paths=source_relative%s:%s", language, pbpkg, outputdir)
	case "java":
		out = fmt.Sprintf("--%s_out=%s", language, outputdir)
	default:
		os.MkdirAll(outputdir, os.ModePerm)
		if len(pbpkg) != 0 {
			pbpkg += ":"
		}
		out = fmt.Sprintf("--%s_out=%s%s", language, pbpkg, outputdir)
	}

	return out
}

func makeProtoPath(protodirs []string) ([]string, error) {

	args := []string{}

	// check if protoc is earlier version, to help make some descisions for compatibility
	old, err := isOldProtocVersion()
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	existed := map[string]struct{}{}
	for _, protodir := range protodirs {

		if _, ok := existed[protodir]; ok {
			continue
		}

		if protodir == wd {
			if old {
				continue
			}
			args = append(args, fmt.Sprintf("--proto_path=%s", protodir))
		} else {
			args = append(args, fmt.Sprintf("--proto_path=%s", protodir))
		}

		existed[protodir] = struct{}{}
	}
	return args, nil
}
