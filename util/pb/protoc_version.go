package pb

import (
	"os/exec"
	"strconv"
	"strings"
)

func protocVersion() (version string, err error) {
	cmd := exec.Command("protoc", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	version = strings.TrimPrefix(strings.TrimSpace(string(output)), "libprotoc ")
	return
}

const (
	MAJOR_VERSION    = 3
	MINOR_VERSION    = 6
	REVISION_VERSION = 0
)

func isOldProtocVersion() (old bool, err error) {

	version, err := protocVersion()
	if err != nil {
		return
	}
	return oldVersion(version)
}

func oldVersion(version string) (old bool, err error) {

	vv := strings.Split(version, ".")
	major, err := strconv.Atoi(vv[0])
	if err != nil {
		return
	}
	minor, err := strconv.Atoi(vv[1])
	if err != nil {
		return
	}
	revision, err := strconv.Atoi(vv[2])
	if err != nil {
		return
	}

	if major < MAJOR_VERSION ||
		(major == MAJOR_VERSION && minor < MINOR_VERSION) ||
		(major == MAJOR_VERSION && minor == MINOR_VERSION && revision < REVISION_VERSION) {
		old = true
		return
	}

	return
}
