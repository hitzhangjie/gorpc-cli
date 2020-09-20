package extension

//go:generate protoc -I ../install/ --go_out=gorpc gorpc.proto
//go:generate protoc -I ../install/ --go_out=swagger swagger.proto
