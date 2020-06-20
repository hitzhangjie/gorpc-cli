package fs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var wd string

func TestMain(m *testing.M) {
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wd = d
	os.Exit(m.Run())
}

func TestBaseFileNameWithoutExt(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"1-gorpc.proto", args{"gorpc.proto"}, "gorpc"},
		{"2-hello.world.proto", args{"hello.world.proto"}, "hello.world"},
		{"3-gorpc.app.server.go", args{"gorpc.app.server.go"}, "gorpc.app.server"},
		{"4-github.com/group/repo/gorpc.app.proto", args{"github.com/group/repo/gorpc.app.proto"}, "gorpc.app"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseNameWithoutExt(tt.args.filename); got != tt.want {
				t.Errorf("BaseNameWithoutExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocateProtofileDir(t *testing.T) {

	tests := []struct {
		name     string
		filename string
		search   string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{"1-good.dat", "good.dat", filepath.Join(wd, "testcase/a/b/"), false},
		{"2-bad.dat", "bad.dat", filepath.Join(wd, "testcase/a/b/c/"), false},
		{"3-hello.dat", "hello.dat", filepath.Join(wd, "testcase/a/b/c/d/"), false},
		{"4-notexist.dat", "notexit.dat", filepath.Join(wd, "testcase/a/b/c/d/"), true},
		{"5-good.dat", "notexit.dat", filepath.Join(wd, "testcase/"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LocateFile(tt.filename, []string{tt.search})
			if (err != nil) != tt.wantErr {
				t.Errorf("LocateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUniqFilePath(t *testing.T) {
	type args struct {
		dirs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{"testcase-1", args{[]string{"/a", "/b", "/a"}}, []string{"/a", "/b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UniqFilePath(tt.args.dirs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrepareOutputdir(t *testing.T) {
	type args struct {
		outputdir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"testcase-1", args{filepath.Join(wd, "testcase/a/")}, false},          // target dir already existed, return nil
		{"testcase-2", args{filepath.Join(wd, "testcase/a/b/good.dat")}, true}, // target existed but not dir, return error
		{"testcase-3", args{filepath.Join(wd, "testcase/fff")}, false},         // target dir not existed, create it, return nil
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PrepareOutputdir(tt.args.outputdir)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareOutputdir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "testcase-3" {
				os.RemoveAll(tt.args.outputdir)
			}
		})
	}
}
