{{- $svrNameCamelCase := (index .Services .ServiceIndex).Name | camelcase -}}
{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}

{{- with .FileOptions.go_package -}}
{{- $goPkgOption = . -}}
{{- end -}}

{{- $serviceIndex := .ServiceIndex -}}
{{ $service := (index .Services .ServiceIndex) -}}

package main

import (
	"context"
    "testing"

	gorpc "github.com/hitzhangjie/gorpc"
	_ "github.com/hitzhangjie/gorpc/http"
	_ "github.com/hitzhangjie/gorpc-selector-cl5"

    "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

    {{ if ne $goPkgOption "" -}}
	pb "{{$goPkgOption}}"
    {{- else -}}
	pb "{{$pkgName}}"
    {{- end }}
    {{ range .Imports }}
    {{ if ne . $pkgName }}
	{{. | gopkg_simple}} "{{- . -}}"
	{{ end }}
    {{- end }}
)

{{$svrName := $service.Name | camelcase | untitle}}
var {{$svrName}}Service = &{{$svrName}}ServiceImpl{}

{{ if ne $goPkgOption "" -}}
//go:generate mockgen -destination=stub/{{$goPkgOption}}/{{$svrName|snakecase}}_mock.go -package={{$goPkgOption|gopkg_simple}} -self_package={{$goPkgOption}} {{$goPkgOption}} {{$svrName|title}}ClientProxy
{{- else -}}
//go:generate mockgen -destination=stub/{{$pkgName}}/{{$svrName|snakecase}}_mock.go -package={{$pkgName|gopkg_simple}} -self_package={{$pkgName}} {{$pkgName}} {{$svrName|title}}ClientProxy
{{- end }}

{{range (index .Services .ServiceIndex).RPC}}
{{- $rpcName := .Name | camelcase -}}
{{- $rpcReqType := .RequestType -}}
{{- $rpcRspType := .ResponseType -}}

{{- if or (eq (trimright "." $rpcReqType) $pkgName) (eq (trimright "." (gofulltype $rpcReqType $.FileDescriptor)) $goPkgOption) -}}
	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
{{- else -}}
	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
{{- end -}}

{{- if or (eq (trimright "." $rpcRspType) $pkgName) (eq (trimright "." (gofulltype $rpcRspType $.FileDescriptor)) $goPkgOption) -}}
	{{- $rpcRspType = (printf "pb.%s" (splitList "." $rpcRspType|last|export)) -}}
{{- else -}}
	{{- $rpcRspType = (gofulltype $rpcRspType $.FileDescriptor) -}}
{{- end -}}

func Test_{{$svrNameCamelCase}}_{{$rpcName}}(t *testing.T) {

	// 开始写mock逻辑
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

   	{{$svrNameCamelCase|untitle}}ClientProxy := pb.NewMock{{$svrNameCamelCase}}ClientProxy(ctrl)

	// 预期行为
	m := {{$svrNameCamelCase|untitle}}ClientProxy.EXPECT().{{$rpcName}}(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	m.DoAndReturn(func(ctx context.Context, req interface{}, opts ...interface{}) (interface{}, error) {

		r, ok := req.(*{{$rpcReqType}})
		if !ok {
			panic("invalid request")
		}

        rsp := &{{$rpcRspType}}{}
		err := {{$svrName}}Service.{{$rpcName}}(gorpc.BackgroundContext(), r, rsp)
		return rsp, err
	})

	// 开始写单元测试逻辑
	req := &{{$rpcReqType}}{}

	rsp, err := {{$svrNameCamelCase|untitle}}ClientProxy.{{$rpcName}}(gorpc.BackgroundContext(), req)

    // 输出入参和返回 (检查t.Logf输出，运行 `go test -v`)
    t.Logf("{{$svrNameCamelCase}}_{{$rpcName}} req: %v", req)
    t.Logf("{{$svrNameCamelCase}}_{{$rpcName}} rsp: %v, err: %v", rsp, err)

	assert.Nil(t, err)
}

{{end}}
