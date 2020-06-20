package pb

import (
	"runtime"
	"testing"
)

func TestLocateGoRPCProto(t *testing.T) {
	fp, err := LocateGoRPCProto()
	if err != nil {
		t.Fatalf("locate gorpc.proto error = %v", err)
	}
	t.Logf("locate gorpc.proto in = %s", fp)
}

func TestHomeWindows(t *testing.T) {
	os := runtime.GOOS
	if os == "windows" {
		path, err := homeWindows()
		if path == "" {
			t.Fatal(err)
		}
	} else {
		t.Logf("os is not windows")
	}
}

func Test_homeWindows(t *testing.T) {
	d, err := homeWindows()
	if runtime.GOOS == "windows" && err != nil {
		t.Error(err)
	}
	if runtime.GOOS == "windows" && len(d) == 0 {
		t.Errorf("get home directory error on Windows")
	}
}
