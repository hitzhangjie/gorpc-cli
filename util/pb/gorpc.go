package pb

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// LocateGoRPCProto 定位gorpc.proto路径
func LocateGoRPCProto() (string, error) {

	var (
		search      []string
		errNotFound = errors.New("cannot locate gorpc.proto.")
	)

	search = append(search, filepath.Join(os.Getenv("HOME"), ".gorpc"))

	switch runtime.GOOS {
	case "linux", "darwin":
		search = append(search, "/etc/gorpc")
	case "windows":
		userHomePath, err := homeWindows()
		if err != nil {
			return "", errors.New(errNotFound.Error() + err.Error())
		}
		search = append(search, filepath.Join(userHomePath, ".gorpc"))
	}

	for _, p := range search {
		ap, err := filepath.Abs(p)
		if err != nil {
			continue
		}

		fp := filepath.Join(ap, "gorpc.proto")
		fin, err := os.Stat(fp)
		if err != nil || fin.IsDir() {
			continue
		}

		return ap, nil
	}

	return "", errNotFound
}

// homeWindows 返回HOME路径
func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("Window OS HOMEDRIVE, HOMEPATH, and USERPROFILE are blank.")
	}
	return home, nil
}
