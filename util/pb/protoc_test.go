package pb

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var wd string

func TestMain(m *testing.M) {
	d, err := os.Getwd()
	if err != nil {
		os.Exit(1)
	}
	wd = filepath.Join(d, "testcase")
	os.Exit(m.Run())
}

func TestProtoc(t *testing.T) {

	languages := []string{"go", "java"}
	outputdir := wd

	type args struct {
		protodirs    []string
		protofile    string
		pbpkgMapping map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{"1-service.proto", args{[]string{wd}, "service.proto", nil}},
		{"2-helloworld-import-gorpc.proto", args{[]string{wd, "/etc/gorpc", "/Users/zhangjie/.gorpc"}, "helloworld-import-gorpc.proto", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, lang := range languages {
				err := Protoc(tt.args.protodirs, tt.args.protofile, lang, outputdir, tt.args.pbpkgMapping)
				if err != nil {
					t.Errorf("Protoc() error = %v", err)
				}
			}
		})
	}

	// clean
	err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".pb.go") || strings.HasSuffix(path, ".java") {
			os.RemoveAll(path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
