package pb

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

// LocateSECVProto 查找GOPATH目录下的SECV目录，用于在protoc.go中拼接环境变量
func LocateSECVProto() (string, error) {

	var (
		search      []string
		errNotFound = errors.New("cannot locate validate.proto. please check your `/etc/profile` setting or GOPATH")
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

		fp := filepath.Join(ap, "validate.proto")
		fin, err := os.Stat(fp)
		if err != nil || fin.IsDir() {
			continue
		}

		return ap, nil
	}

	return "", errNotFound
}
