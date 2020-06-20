package swagger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hitzhangjie/gorpc/descriptor"
	"github.com/hitzhangjie/gorpc/params"
	"github.com/hitzhangjie/gorpc/util/log"
	"github.com/hitzhangjie/gorpc/util/pb"

	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

// FieldDescriptorProto_Type_name 用于将 protobuf 库中的 field 类型定义转化为 string 的 map 表
var FieldDescriptorProto_Type_name = map[protobuf.FieldDescriptorProto_Type]string{
	1:  "double",
	2:  "float",
	3:  "int64",
	4:  "uint64",
	5:  "int32",
	6:  "fixed64",
	7:  "fixed32",
	8:  "bool",
	9:  "string",
	10: "group",
	11: "message",
	12: "bytes",
	13: "uint32",
	14: "enum",
	15: "sfixed32",
	16: "sfixed64",
	17: "sint32",
	18: "sint64",
}

// GenSwagger 提供外部结构，用于生成 swagger json
func GenSwagger(fd *descriptor.FileDescriptor, option *params.Option) error {
	// 1 先处理定义了 service 的 pb 文件，该文件中含有 rpc 方法出入参的数据模型定义
	serviceMessageMap, err := getServiceMessage(fd, option)
	if err != nil {
		return err
	}

	// 2 再处理 import 的 pb 文件，import 的 pb 文件中，也包含了一些数据模型的定义
	importMessageMap, err := getImportMessage(fd, option)
	if err != nil {
		return err
	}

	// 3 合并两种 pb 文件得到的数据模型
	msgMap := make(map[string]ModelStruct)
	for key, val := range serviceMessageMap {
		msgMap[key] = val
	}
	for key, val := range importMessageMap {
		msgMap[key] = val
	}

	pathsMap := make(map[string]map[string]*MethodStruct)

	// 4 组装各个方法的信息
	for _, service := range fd.Services {
		for _, rpc := range service.RPC {
			// 获取到每个 rpc 方法前面定义的描述（如果有）
			summary := rpc.LeadingComments
			// 校验之前在 option 中收集到的对于 rpc 方法的名称，如果为空，则将
			if len(rpc.SwaggerInfo.Title) != 0 {
				summary = rpc.SwaggerInfo.Title
			}
			// 组装 method 的信息
			method := &MethodStruct{
				Summary:     summary,
				OperationId: rpc.Name,
				Responses:   getResponseMsg(rpc, msgMap),
				Parameters:  getParameters(rpc, msgMap),
				Tags:        []string{service.Name},
				Description: rpc.SwaggerInfo.Description,
			}

			// 组装 path 消息
			pathsMap[rpc.FullyQualifiedCmd] = map[string]*MethodStruct{
				rpc.SwaggerInfo.Method: method,
			}
		}
	}

	// 5 获取文件信息组装 swagger json 文档头部信息
	filePath, err := filepath.Abs(fd.FilePath)
	if err != nil {
		return err
	}
	_, fileName := filepath.Split(filePath)
	title := strings.ReplaceAll(fileName, ".proto", "")
	infoMap := InfoStruct{
		Title:       title,
		Description: fmt.Sprintf("The api document of %s", fileName),
		Version:     "2.0",
	}

	// 6 组装整个 swagger json 信息
	swaggerJson := &SwaggerJson{
		Swagger:     "2.0",
		Info:        infoMap,
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Paths:       pathsMap,
		Definitions: msgMap,
	}

	// 7 格式化输出 json 文件，保证 json 输出的字符串不是直接在一行内输出
	jsonByte, err := json.MarshalIndent(swaggerJson, "", " ")
	if err != nil {
		return err
	}

	// 8 写入文件
	if err = ioutil.WriteFile("apidocs.swagger.json", jsonByte, 0666); err != nil {
		return err
	}

	log.Info("Generate the api document of ```%s``` success", fileName)

	return nil
}

// 组装 rpc 方法的出参信息
func getResponseMsg(rpc *descriptor.RPCDescriptor, msgMap map[string]ModelStruct) map[string]ResponseStruct {
	var ret map[string]ResponseStruct
	if rspMsg, ok := msgMap[rpc.ResponseType]; ok {
		ret = map[string]ResponseStruct{
			"200": {
				Description: rspMsg.Title,
				Schema:      SchemaStruct{Ref: fmt.Sprintf("#/definitions/%s", rpc.ResponseType)},
			},
		}
	}
	return ret
}

// 组装 rpc 方法的入参信息
func getParameters(rpc *descriptor.RPCDescriptor, msgMap map[string]ModelStruct) []*ParametersStruct {
	paraArr := make([]*ParametersStruct, 0)
	if reqMsg, ok := msgMap[rpc.RequestType]; ok {
		for name, property := range reqMsg.Properties {
			in := "query" // 默认是query 类型，如果为 message 类型的field，则为 body
			var schema *SchemaStruct
			if property.Type == "message" {
				in = "body"
				schema = &SchemaStruct{Ref: property.Ref}
			}
			para := &ParametersStruct{
				Name:        name,
				In:          in,
				Required:    false, // 默认设置为非必要字段
				Type:        property.Type,
				Schema:      schema,
				Format:      property.Format,
				Description: property.Description,
			}
			paraArr = append(paraArr, para)
		}
	}

	return paraArr
}

// 获取 service 定义的 pb 文件中的 model 定义
func getServiceMessage(nfd *descriptor.FileDescriptor, option *params.Option) (map[string]ModelStruct, error) {
	protodirs := option.Protodirs

	p, err := pb.LocateGoRPCProto()
	if err != nil {
		return nil, err
	}
	protodirs = append(protodirs, p)

	parser := protoparse.Parser{
		IncludeSourceCodeInfo: true,
		ImportPaths:           protodirs,
	}
	fds, err := parser.ParseFiles(option.Protofile)
	if err != nil {
		panic(err)
	}
	fd := fds[0]

	packageName := fd.GetPackage()

	msgMap := make(map[string]ModelStruct)

	for _, msg := range fd.GetMessageTypes() {
		key := fmt.Sprintf("%s.%s", packageName, msg.GetName())

		// 获取 message 的注释
		msgComment := msg.GetSourceInfo().GetLeadingComments()
		msgComment = strings.ReplaceAll(msgComment, "\n", "")

		msgMap[key] = ModelStruct{
			Type:       "object",
			Properties: getProperties(msg),
			Title:      msgComment,
		}
	}

	return msgMap, nil
}

// 获取 import proto 中定义的 model
func getImportMessage(nfd *descriptor.FileDescriptor, option *params.Option) (map[string]ModelStruct, error) {
	//pwd, err := os.Getwd()
	//if err != nil {
	//	return nil, err
	//}

	protodirs := option.Protodirs

	p, err := pb.LocateGoRPCProto()
	if err != nil {
		return nil, err
	}
	protodirs = append(protodirs, p)

	msgMap := make(map[string]ModelStruct)
	for fname, _ := range nfd.Pb2ImportPath {
		// 跳过google官方提供的pb文件，gorpc扩展文件，swagger 扩展文件
		if strings.HasPrefix(fname, "google/protobuf") || fname == "gorpc.proto" || fname == "swagger.proto" {
			continue
		}

		if fname == "validate.proto" {
			continue
		}

		//importPath := filepath.Join(pwd, fname)
		importPath := fname
		parser := protoparse.Parser{
			IncludeSourceCodeInfo: true,
			ImportPaths:           protodirs,
		}
		fds, err := parser.ParseFiles(importPath)
		if err != nil {
			panic(err)
		}
		fd := fds[0]

		packageName := fd.GetPackage()

		// 遍历每个 message 组装数据模型
		for _, msg := range fd.GetMessageTypes() {
			key := fmt.Sprintf("%s.%s", packageName, msg.GetName())

			// 获取 message 的注释，用于 title 描述
			msgComment := msg.GetSourceInfo().GetLeadingComments()
			msgComment = strings.ReplaceAll(msgComment, "\n", "")

			msgMap[key] = ModelStruct{
				Type:       "object",
				Properties: getProperties(msg),
				Title:      msgComment,
			}
		}
	}

	return msgMap, nil
}

// 获取 message field 的消息填充到 properties 中
func getProperties(msg *desc.MessageDescriptor) map[string]PropertiesStruct {
	// 获取 message 的 field 的消息填充到 properties 中
	propertiesMap := make(map[string]PropertiesStruct)
	for _, field := range msg.GetFields() {
		var inputComment string = ""
		// 获取 field 定义前面的注释
		fieldHeadComment := field.GetSourceInfo().GetLeadingComments()
		fieldHeadComment = strings.ReplaceAll(fieldHeadComment, "\n", "")

		// 获取 field 定义后面的注释
		fieldTailComment := field.GetSourceInfo().GetTrailingComments()
		fieldTailComment = strings.ReplaceAll(fieldTailComment, "\n", "")

		// field 定义后面的注释优先，定义为描述该 field 的注释。
		if len(fieldTailComment) != 0 {
			inputComment = fieldTailComment
		} else {
			inputComment = fieldHeadComment
		}

		// 如果是消息类型则引用
		if field.GetType() == protobuf.FieldDescriptorProto_TYPE_MESSAGE {
			propertiesMap[field.GetName()] = PropertiesStruct{
				Ref: fmt.Sprintf("#/definitions/%s", field.GetMessageType().GetFullyQualifiedName()),
			}
			continue
		}
		// 非消息类型的则直接组装 properties
		propertiesMap[field.GetName()] = PropertiesStruct{
			Type:        getTypeStr(field.GetType()),
			Format:      getFormatStr(field.GetType()),
			Description: inputComment,
		}
	}

	return propertiesMap
}

// 根据 protobuf 的类型返回具体的 field 类型
func getTypeStr(t protobuf.FieldDescriptorProto_Type) string {
	switch t {
	case protobuf.FieldDescriptorProto_TYPE_DOUBLE, protobuf.FieldDescriptorProto_TYPE_FLOAT,
		protobuf.FieldDescriptorProto_TYPE_INT64, protobuf.FieldDescriptorProto_TYPE_UINT64,
		protobuf.FieldDescriptorProto_TYPE_INT32, protobuf.FieldDescriptorProto_TYPE_UINT32,
		protobuf.FieldDescriptorProto_TYPE_FIXED64, protobuf.FieldDescriptorProto_TYPE_FIXED32,
		protobuf.FieldDescriptorProto_TYPE_SFIXED64, protobuf.FieldDescriptorProto_TYPE_SFIXED32,
		protobuf.FieldDescriptorProto_TYPE_SINT64, protobuf.FieldDescriptorProto_TYPE_SINT32:
		return "number"
	case protobuf.FieldDescriptorProto_TYPE_STRING, protobuf.FieldDescriptorProto_TYPE_BYTES:
		return "string"
	case protobuf.FieldDescriptorProto_TYPE_MESSAGE:
		return "message"
	}

	return "string"
}

// 根据 protobuf 的类型返回具体的 field 的格式
func getFormatStr(t protobuf.FieldDescriptorProto_Type) string {
	if val, ok := FieldDescriptorProto_Type_name[t]; ok {
		return val
	}

	return "string"
}
