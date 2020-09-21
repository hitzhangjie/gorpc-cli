package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// dependency 依赖工具
type dependency struct {
	Name    string
	Version string
	Cmd     string
}

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

// Pipeline strings together the given exec.Cm
// to the Unix pipeline.  Each command's standard output is connected to the
// standard input of the next command, and the output of the final command in
// the pipeline is returned, along with the collected standard error of all
// commands and the first error found (if any).
//
// To provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}
