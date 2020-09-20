{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}
package main

import (
	_ "github.com/hitzhangjie/gorpc-config-tconf"
	_ "github.com/hitzhangjie/gorpc-filter/debuglog"
	_ "github.com/hitzhangjie/gorpc-filter/recovery"
	_ "github.com/hitzhangjie/gorpc-log-atta"
	_ "github.com/hitzhangjie/gorpc-metrics-m007"
	_ "github.com/hitzhangjie/gorpc-metrics-runtime"
	_ "github.com/hitzhangjie/gorpc-naming-polaris"
	_ "github.com/hitzhangjie/gorpc-opentracing-tjg"
	_ "github.com/hitzhangjie/gorpc-selector-cl5"
	_ "go.uber.org/automaxprocs"

	"github.com/hitzhangjie/gorpc-cli/log"

	gorpc "github.com/hitzhangjie/gorpc"
    {{ if ne $goPkgOption "" -}}
   	pb "{{$goPkgOption}}"
    {{- else -}}
    pb "{{$pkgName}}"
	{{- end }}
)

{{range $index, $service := .Services}}
{{- $svrName := $service.Name | camelcase | untitle -}}
type {{$svrName}}ServiceImpl struct {}
{{end}}

func main() {

	s := gorpc.NewServer()

    {{range $index, $service := .Services}}
    {{- $svrNameCamelCase := $service.Name | camelcase -}}
	pb.Register{{$svrNameCamelCase}}Service(s, &{{$svrNameCamelCase|untitle}}ServiceImpl{})
	{{end}}

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
