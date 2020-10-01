package extension

//go:generate protoc -I gorpc --go_out=gorpc gorpc.proto
//go:generate cp -f gorpc/gorpc.proto ../install/
