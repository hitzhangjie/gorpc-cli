package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// dependency 依赖工具
type dependency struct {
	Name    string // 依赖工具名称
	Version string // 工具最小版本
	Cmd     string // 获取版本命令
}

// loadDependencies 加载工具版本信息
func loadDependencies() ([]*dependency, error) {
	deps := []*dependency{
		{
			Name:    "protoc",
			Version: "v3.6.0",
			Cmd:     "protoc --version | awk '{print $2}'",
		}, {
			Name:    "protoc-gen-go",
			Version: "",
			Cmd:     "",
		},
	}
	return deps, nil
}

// checkDependencies 检查工具版本
func checkDependencies(deps []*dependency) error {

	for _, dep := range deps {
		// 检查依赖是否存在
		_, err := exec.LookPath(dep.Name)
		if err != nil {
			return fmt.Errorf("%s not found, %v", dep.Name, err)
		}

		// 不需要检查版本就跳过
		if len(dep.Cmd) == 0 || len(dep.Version) == 0 {
			continue
		}

		// 执行自定义命令获取版本信息
		cmd := exec.Command("sh", "-c", dep.Cmd)

		buf, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s load version, %v, \n\t%s\n", dep.Name, err, string(buf))
		}

		err = checkVersion(string(buf), dep.Version)
		if err != nil {
			return fmt.Errorf("%s mismatch, %v", dep.Name, err)
		}
	}

	return nil
}

// checkVersion 检查版本是否满足要求
func checkVersion(version, required string) error {

	if len(version) != 0 && version[0] == 'v' || version[0] == 'V' {
		version = version[1:]
	}

	if len(required) != 0 && required[0] == 'v' || required[0] == 'V' {
		required = required[1:]
	}

	m1, n1, r1 := versions(version)
	m2, n2, r2 := versions(required)

	if !(m1 >= m2 && n1 >= n2 && r1 >= r2) {
		return fmt.Errorf("require version: %s", required)
	}

	return nil
}

// versions 获取主、副、修订版本
func versions(ver string) (major, minor, revision int) {

	var err error

	vv := strings.Split(ver, ".")

	if len(vv) >= 1 {
		major, err = strconv.Atoi(vv[0])
		if err != nil {
			return
		}
	}

	if len(vv) >= 2 {
		minor, err = strconv.Atoi(vv[1])
		if err != nil {
			return
		}
	}

	if len(vv) >= 3 {
		revision, err = strconv.Atoi(vv[2])
		if err != nil {
			return
		}
	}

	return
}
