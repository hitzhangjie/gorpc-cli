package descriptor

import (
	"encoding/json"

	"github.com/hitzhangjie/gorpc/util/log"
	"github.com/jhump/protoreflect/desc"
)

// FileDescriptor 文件作用域相关的描述信息
type FileDescriptor struct {
	FilePath    string                 // filepath (abs)
	PackageName string                 // pb包名称
	AppName     string                 // app名称
	ServerName  string                 // server名称
	Imports     []string               // pb文件可能import其他pb文件，rpc请求、响应中若有引用，记录类型对应的导入包名
	FileOptions map[string]interface{} // fileoptions
	Services    []*ServiceDescriptor   // 支持多service

	Pb2ValidGoPkg map[string]string // k=pb文件名, v=protoc处理后package名
	Pb2ImportPath map[string]string // k=pb文件名，v=go代码中对应importpath

	//BUG: https://github.com/hitzhangjie/gorpc/issues/186
	//Deprecated
	Pkg2ValidGoPkg map[string]string // k=pb文件package directive, v=protoc处理后package名

	//BUG: https://github.com/hitzhangjie/gorpc/issues/186
	//Deprecated
	Pkg2ImportPath map[string]string // k=pb文件package directive, v=go代码中对应importpath

	RpcMessageType map[string]string // k=pb定义的pkg.typ，v=有效的go中的pkg.typ

	fd *desc.FileDescriptor // 原始descriptor
}

// RawFileDescriptor 返回原始的desc.FileDescriptor
func (fd *FileDescriptor) RawFileDescriptor() *desc.FileDescriptor {
	return fd.fd
}

// SetRawFileDescriptor 设置原始的desc.FileDescriptor
func (fd *FileDescriptor) SetRawFileDescriptor(rawFd *desc.FileDescriptor) {
	fd.fd = rawFd
}

func (fd *FileDescriptor) Dump() {
	log.Debug("************************** FileDescriptor ***********************")
	buf, _ := json.MarshalIndent(fd, "", "  ")
	log.Debug("\n%s", string(buf))
	log.Debug("*****************************************************************")
}

// ServiceDescriptor service作用域相关的描述信息
type ServiceDescriptor struct {
	Name string           // 服务名称
	RPC  []*RPCDescriptor // rpc接口定义
}

// RPCDescriptor rpc作用域相关的描述信息
//
// RequestType由于涉及到
type RPCDescriptor struct {
	Name              string            // RPC方法名
	Cmd               string            // RPC命令字
	FullyQualifiedCmd string            // 完整的RPC命令字，用于ServiceDesc、client请求时命令字
	RequestType       string            // RPC请求消息类型，包含package，比如package_a.TypeA
	ResponseType      string            // RPC响应消息类型，包含package，比如package_b.TypeB
	LeadingComments   string            // RPC前置注释信息
	TrailingComments  string            // RPC后置注释信息
	SwaggerInfo       SwaggerDescriptor // 用于生成 swagger 文档的信息
}

// SwaggerDescriptor swagger api 文档生成所需的描述信息
type SwaggerDescriptor struct {
	Title       string // RPC 方法名
	Method      string // http 协议的 method，如果该方法支持 http 协议
	Description string // 方法的描述
}
