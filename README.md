# gorpc

***gorpc*** is an efficient tool to help developers :
- `gorpc create`, generate service template or rpc stub
- `gorpc update`, update service template or rpc stub
- `gorpc issue`, fire your browser at issue page
- `gorpc version`, show the version
- `gorpc rpc`, launch rpc test to your work

gorpc is built upon cobra, it is easy to extend new abilities.

## Using "Google Protobuf" as IDL

***Google Protobuf*** is developed by Google, it's a self-descriptive message format.

Using protobuf as IDL (Interface Descriptor Language) is very common, Google also 
provides a protobuf compiler called `protoc`.

Before, I wrote an article to introduce how `protoc` and `proto-gen-go` coordinate to work, 
and some internals of protocol buffer. If you're interested, please read my it:
[Protoc及其插件工作原理分析(精华版)](https://hitzhangjie.github.io/blog/2017-05-23-protoc%E5%8F%8A%E6%8F%92%E4%BB%B6%E5%B7%A5%E4%BD%9C%E5%8E%9F%E7%90%86%E5%88%86%E6%9E%90%E7%B2%BE%E5%8D%8E%E7%89%88/)

## How does protoc and protoc plugins coordinate to work ?

We usually execute command like `protoc --cpp_out/--java_out` to generate cpp 
header/source, java files. 

While for many other languages, `protoc` doesn't implement it, such as, go programming 
language. If we want to generate `*.pb.go` like `*.pb.h/*.pb.cc` or `*.pb.java`, we 
should implement a plugin for `go`, we already have a tool, that is `protoc-gen-go`.

Take protoc and protoc-gen-go as an example, let's see how they coordinate to work.

Just now, we know protobuf is a self-descriptive message format, when file `*.proto` 
parsed by `protoc`, a `FileDesciptorProto` object will be built, it contains nearly
everything about the `*.proto` we written. If you know little about internals of `protoc`
or `protobuf` itself, please refer to my article metioned above.

when run command `protoc --go_out *.proto`, protoc will read your protofile and parse it,
after that, it build a FileDescriptorProto object, then it will serialize it and search
executable named `protoc-gen-go` in your `PATH` shell env variable. If found, it will
fork a childprocess to run `protoc-gen-go`, and parentprocess `protoc` will create a 
pipe btw itself and childprocess to communicate. `protoc` will send a `CodeGenerateRequest`
to the childprocess `protoc-gen-go` via pipe. This `CodeGenerateRequest` contains 
serialized `FileDescriptorProto`, then `protoc-gen-go` read from pipe and extract it.
`protoc-gen-go` will be responsible for generate source code by `g.P("..")`. This generated
source code info will be responded to `protoc`, `protoc` process will create file and 
write file content (source code).

This is the way `protoc` and `protoc-gen-go` works.

## Why we choose go templates to generate code ?

Usually, we can use protoc-gen-go as an starting point, we can add some files to generate 
besiding *.pb.go, for example, some default configuration files, or other go code.

But, writing a new protoc plugin like protoc-gen-go is really not a good idea for 
generating source code, especially you want to generate language for many more languages, 
or you want to add some flags, etc.

This manner, writing a new protoc plugin, increases the difficulty in maintenance and 
extensibility. If you have written some before, you'll know what I am saying.
- use `g.P(), g.In(), g.Out()` to generate or format code;

Though we could use go tempalates instead of this, like `hitzhangjie/protoc-gen-gorpc`, 
passing flags is another big problem, which seriously limit the funtionality, and make 
it hard to use.

We could parse the *.proto file once, then using template technology to generate files. 
If we want to support new project template, we just change or add project template, needless
to change the code. And, we could add command, subcommand, flags easily by cobra to extend
the tool's functionalities.

## How to use gorpc ?

run `gorpc` or `gorpc help` to show the help message, you can run `gorpc help create` to see
more details relevant to `gorpc create`.

```bash
gorpc 是一个效率工具，方便gorpc服务的开发.

例如: 
- 指定pb文件，快速生成完整的工程，或者生成对应的rpcstub
- 对目标服务发起rpc测试请求 

尝试用gorpc框架+gorpc工具来编写你的下一个gorpc服务吧 !

Usage:
  gorpc [command]

Available Commands:
  create      指定pb文件快速创建工程或rpcstub
  help        Help about any command
  issue       反馈一个issue
  version     显示gorpc命令的版本(commit hash)

Flags:
  -h, --help   help for gorpc
```

## Contribution

Welcome the contribution from you !
