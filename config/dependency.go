package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hitzhangjie/codeblocks/log"
)

// Dependency describes a Dependency
type Dependency struct {
	Executable string `json:"executable"`  // name of Dependency
	VersionMin string `json:"version_min"` // minimum version, a.b.c
	VersionCmd string `json:"version_cmd"` // cmd to get version
	InstallCmd string `json:"install_cmd"` // cmd to install
	Fallback   string `json:"fallback"`    // fallback message
}

func (d *Dependency) Version() (string, error) {
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
func (d *Dependency) Installed() (bool, error) {
	_, err := exec.LookPath(d.Executable)
	if err != nil {
		return false, fmt.Errorf("not found in PATH")
	}
	return true, nil
}

func (d *Dependency) Install() error {
	if len(d.InstallCmd) == 0 {
		return fmt.Errorf("install cmd empty, tips: %s", d.Fallback)
	}

	cmd := exec.Command("sh", "-c", d.InstallCmd)
	buf, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("install cmd error: %v, \n%s", err, string(buf))
	}
	return nil
}

func (d *Dependency) CheckVersion() (passed bool, err error) {
	// skip checking if cmd/version not specified
	if len(d.VersionMin) == 0 {
		return true, nil
	}

	v, err := d.Version()
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

// LoadDependencies load dependencies and version requirements
func LoadDependencies() ([]Dependency, error) {

	d, err := LocateTemplatePath()
	if err != nil {
		return nil, err
	}

	f := filepath.Join(d, "dependencies.json")
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	deps := []Dependency{}
	err = json.Unmarshal(b, &deps)
	if err != nil {
		return nil, err
	}

	return deps, nil
}

// CheckDependencies check if dependencies meet the version requirements
func CheckDependencies(deps []Dependency) error {

	for _, dep := range deps {
		// check installed or not
		_, err := dep.Installed()
		if err != nil {
			log.Info("check %s, not installed, try installing", dep.Executable)
			err = dep.Install()
			if err != nil {
				log.Error("check %s, not installed, try installing failed: %v", dep.Executable, err)
				return fmt.Errorf("check %s, not installed: %v", dep.Executable, err)
			}
			log.Info("check %s, not installed, try installing done", dep.Executable)
		}

		// check if `dep's version` meet the requirement of `versionCmd`
		ok, err := dep.CheckVersion()
		if err != nil {
			log.Error("check %s, check version error: %v", dep.Executable, err)
			return fmt.Errorf("check %s, check error: %v", dep.Executable, err)
		}
		if !ok {
			log.Error("check %s, check version too old", dep.Executable)
			return fmt.Errorf("check %s, version too old", dep.Executable)
		}
		log.Info("check %s, passed", dep.Executable)
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
