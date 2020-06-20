package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestCreateCmd(t *testing.T) {
	pwd, _ := os.Getwd()
	wd := filepath.Dir(pwd)
	testcase := filepath.Join(wd, "testcase/testcase.create")

	type args struct {
		dir    string
		pbfile string
		outdir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"1-1pb-without-import",
			args{"1-1pb-without-import", "helloworld.proto", "1-helloworld"},
			false,
		}, {
			"2-Npb-same-pkg",
			args{"2-Npb-same-pkg", "helloworld.proto", "2-Npb-same-pkg"},
			false,
		}, {
			"3-Npb-diff-pkgs",
			args{"3-Npb-diff-pkgs", "helloworld.proto", "3-Npb-diff-pkgs"},
			false,
		}, {
			"4-Npb-same-pkg-diff-dirs",
			args{"4-Npb-same-pkg-diff-dirs", "helloworld.proto", "4-Npb-same-pkg-diff-dirs"},
			false,
		}, {
			"4-Npb-same-pkg-diff-dirs-more",
			args{"4-Npb-same-pkg-diff-dirs-more", "helloworld.proto", "4-Npb-same-pkg-diff-dirs-more"},
			false,
		}, {
			"6-Npb-same-pkgdirective-diff-gopkgopts",
			args{"6-Npb-same-pkgdirective-diff-gopkgopts", "feeds_read.proto", "6-Npb-same-pkgdirective-diff-gopkgopts"},
			false,
		}, {
			"6.1-other-scene-biz_service",
			args{"6-other-scene/biz_service", "biz_init_rpc.proto", "6-other-scene/biz-service"},
			false,
		}, {
			"6.2-other-scene-hello_service",
			args{"6-other-scene/hello_service", "hello.proto", "6-other-scene/hello_service"},
			false,
		}, {
			"6.3-other-scene-rec_service",
			args{"6-other-scene/rec_service", "rec_interface_service.proto", "6-other-scene/rec_service"},
			false,
		}, {
			"6.4-other-scene-test_service",
			args{"6-other-scene/test_service", "test_service_rpc.proto", "6-other-scene/test_service"},
			false,
		},
	}

	out := filepath.Join(testcase, "generated")
	os.RemoveAll(out)

	for _, tt := range tests {

		any := filepath.Join(testcase, tt.args.dir)
		err := os.Chdir(any)
		if err != nil {
			t.Fatal(err)
		}

		dirs := []string{}
		err = filepath.Walk(any, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		})
		if err != nil {
			t.Fatal("walk testcase error")
		}

		opts := []string{}
		for _, d := range dirs {
			opts = append(opts, "--protodir", d)
		}

		opts = append(opts, "--protofile", tt.args.pbfile)
		out := filepath.Join(testcase, "generated", tt.args.outdir)
		opts = append(opts, "-o", out)

		resetFlags(createCmd)
		if err := createCmd.RunE(createCmd, opts); (err != nil) != tt.wantErr {
			t.Errorf("TEST CreateCmd/%s Run() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Value.Type() == "stringSlice" {
			// XXX: unfortunately, flag.Value.Set() appends to original
			// slice, not resets it, so we retrieve pointer to the slice here
			// and set it to new empty slice manually
			value := reflect.ValueOf(flag.Value).Elem().FieldByName("value")
			ptr := (*[]string)(unsafe.Pointer(value.Pointer()))
			*ptr = make([]string, 0)
		}
		if flag.Value.Type() == "stringArray" {
			// XXX: unfortunately, flag.Value.Set() appends to original
			// slice, not resets it, so we retrieve pointer to the slice here
			// and set it to new empty slice manually
			value := reflect.ValueOf(flag.Value).Elem().FieldByName("value")
			ptr := (*[]string)(unsafe.Pointer(value.Pointer()))
			*ptr = make([]string, 0)
		}
		flag.Value.Set(flag.DefValue)
	})
	for _, cmd := range cmd.Commands() {
		resetFlags(cmd)
	}
}
