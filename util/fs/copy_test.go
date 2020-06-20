package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1-copy/fa", args{"copy/fa", "target/fa"}, false},
		{"2-copy/fa", args{"copy/fa", "target/copy/fa"}, false},
		{"3-copy/d", args{"copy/d", "target/d"}, false},
		{"4-copy/d", args{"copy/d", "target/e/d"}, false},
	}
	tmpdir := filepath.Join(wd, "testcase/target")
	err := os.MkdirAll(tmpdir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := filepath.Join(wd, "testcase", tt.args.src)
			dest := filepath.Join(wd, "testcase", tt.args.dest)

			if err := Copy(src, dest); (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
