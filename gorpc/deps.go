package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hitzhangjie/gorpc-cli/config"
)

// dependency describes a dependency
type dependency struct {
	Executable string `json:"executable"`  // name of dependency
	VersionMin string `json:"version_min"` // minimum version, a.b.c
	VersionCmd string `json:"version_cmd"` // cmd to get version
	InstallCmd string `json:"install_cmd"` // cmd to install
	Fallback   string `json:"fallback"`    // fallback message
}

func (d *dependency) version() (string, error) {
	if len(d.VersionCmd) == 0 {
		return "", errors.New("install cmd empty")
	}

	// run specified cmd to get version
	cmd := exec.Command("sh", "-c", d.VersionCmd)

	buf, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// check installed or not
func (d *dependency) installed() (bool, error) {
	_, err := exec.LookPath(d.Executable)
	if err != nil {
		return false, fmt.Errorf("not found in PATH")
	}
	return true, nil
}

func (d *dependency) checkVersion() (passed bool, err error) {
	// skip checking if cmd/version not specified
	if len(d.VersionMin) == 0 {
		return true, nil
	}

	v, err := d.version()
	if err != nil {
		return false, err
	}

	version := v
	required := d.VersionMin

	if len(version) != 0 && version[0] == 'v' || version[0] == 'V' {
		version = version[1:]
	}

	if len(required) != 0 && required[0] == 'v' || required[0] == 'V' {
		required = required[1:]
	}

	m1, n1, r1 := versions(version)
	m2, n2, r2 := versions(required)

	if !(m1 >= m2 && n1 >= n2 && r1 >= r2) {
		return false, fmt.Errorf("version too old")
	}
	return true, nil
}

// loadDependencies load dependencies and version requirements
func loadDependencies() ([]dependency, error) {

	d, err := config.LocateTemplatePath()
	if err != nil {
		return nil, err
	}

	f := filepath.Join(d, "dependencies.json")
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	deps := []dependency{}
	err = json.Unmarshal(b, &deps)
	if err != nil {
		return nil, err
	}

	return deps, nil
}

// checkDependencies check if dependencies meet the version requirements
func checkDependencies(deps []dependency) error {

	for _, dep := range deps {
		// check installed or not
		_, err := dep.installed()
		if err != nil {
			return fmt.Errorf("check %s, not installed: %v", dep.Executable, err)
		}

		// check if `dep's version` meet the requirement of `versionCmd`
		ok, err := dep.checkVersion()
		if err != nil {
			return fmt.Errorf("check %s, check error: %v", dep.Executable, err)
		}
		if !ok {
			return fmt.Errorf("check %s, version too old", dep.Executable)
		}
	}

	return nil
}

// versions extract the major, minor and revision (patching) version
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
