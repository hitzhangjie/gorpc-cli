package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMove(t *testing.T) {
	type args struct {
		src string
		dst string
	}

	//backup
	src := filepath.Join(wd, "testcase/move")
	dst := filepath.Join(wd, "testcase/move.bak")
	err := Copy(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(src)
		os.Rename(dst, src)
	}()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// testcases for move a file
		{"case1.1-src-notexist-dst-notexist", args{"notexist", "move/notexist"}, true},
		{"case1.2-src-notexist-dst-exist", args{"notexist", "move/d"}, true},
		{"case2.1-src-file-dst-notexist-dir(dst)isfolder", args{"move/a", "move/d/a"}, false},
		{"case2.2-src-file-dst-notexist-dir(dst)isnotfolder", args{"move/b", "move/nf/a"}, true},
		{"case2.3-src-file-dst-notexist-dir(dst)notexist", args{"move/b", "move/mf/a"}, true},
		{"case3.1-src-file-dst-folder-dst/basename(src)-exist", args{"move/b", "move/d"}, false},
		{"case3.2-src-file-dst-folder-dst/basename(src)-notexist", args{"move/c", "move/d"}, false},
		{"case4.1-src-file-dst-file", args{"move/fd", "move/fd"}, false},
		// testcases for move a directory
		{"case5.1-src-folder-dst-notexist-dir(dst)folder", args{"move/d", "move/e/d"}, false},
		{"case5.2-src-folder-dst-notexist-dir(dst)file", args{"move/e", "move/nf/e"}, true},
		{"case5.3-src-folder-dst-notexist-dir(dst)notexist", args{"move/e", "move/notexist/e"}, true},
		{"case6.1-src-folder-dst-file", args{"move/e", "move/nf"}, true},
		{"case7.1-src-folder-dst-folder-dst/basename(src)existed+empty", args{"move/z", "move/x/y/"}, false},
		{"case7.2-src-folder-dst-folder-dst/basename(src)existed_notempty", args{"move/z1", "move/x/"}, true},
		{"case7.3-src-folder-dst-folder-dst/basename(src)notexist", args{"move/z1", "move/q"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := filepath.Join(wd, "testcase", tt.args.src)
			dst := filepath.Join(wd, "testcase", tt.args.dst)
			if err := Move(src, dst); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
