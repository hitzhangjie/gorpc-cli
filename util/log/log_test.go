package log

import "testing"

func TestInfo(t *testing.T) {
	InitLogging(true)
	Info("my name is: %s", "zhangjie")
	Debug("my name is: %s", "zhangjie")
	Error("my name is: %s", "zhangjie")
}
