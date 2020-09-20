{{- $pkgName := .PackageName -}}
{{- $goPkgName := .PackageName -}}
{{- $importedPkg := .Pkg2ValidGoPkg -}}
{{- $imports := .Imports -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
  {{- $goPkgName = (splitList "/" $goPkgOption)|last|gopkg -}}
{{- end -}}
package main

import (
	"context"

{{ $service0 := index .Services 0 -}}
{{ $method := (index $service0.RPC 0) -}}
{{- $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}

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

{{ if or (eq (index (splitList "." $rpcReqType) 0) "pb") (eq (index (splitList "." $rpcRspType) 0) "pb") }}
{{ if ne $goPkgOption "" }}
	pb "{{ $goPkgOption }}"
{{- else }}
	pb "{{$pkgName -}}"
{{- end }}
{{- end }}
{{ range $imports }}
{{- $importPkg := . }}
{{/*- if or (hasprefix $importPkg $rpcReqType) (hasprefix $importPkg $rpcRspType) */}}
{{/*- if or (hassuffix (index (splitList "." $rpcReqType) 0) $importPkg) (hassuffix (index (splitList "." $rpcRspType) 0) $importPkg) */}}
{{- if or (ne (index (splitList "." $rpcReqType) 0) "pb") (ne (index (splitList "." $rpcRspType) 0) "pb") }}
    {{ $importPkg | gopkg_simple}} "{{ $importPkg }}"
{{ end -}}
{{/* end -*/}}
{{ end -}}
)

{{ $service := (index .Services .ServiceIndex) -}}
{{- $svrName := $service.Name | camelcase | untitle -}}
{{ range $index, $method := $service.RPC }}
{{- $rpcName := $method.Name | camelcase -}}
{{- $rpcReqType := $method.RequestType -}}
{{- $rpcRspType := $method.ResponseType -}}

{{- if or (eq (trimright "." $rpcReqType|gopkg) ($pkgName|gopkg)) (eq (trimright "." (gofulltype $rpcReqType $.FileDescriptor)|gopkg) ($goPkgOption|gopkg)) -}}
	{{- $rpcReqType = (printf "pb.%s" (splitList "." $rpcReqType|last|export)) -}}
{{- else -}}
	{{- $rpcReqType = (gofulltype $rpcReqType $.FileDescriptor) -}}
{{- end -}}

{{- if or (eq (trimright "." $rpcRspType|gopkg) ($pkgName|gopkg)) (eq (trimright "." (gofulltype $rpcRspType $.FileDescriptor)|gopkg) ($goPkgOption|gopkg)) -}}
	{{- $rpcRspType = (printf "pb.%s" (splitList "." $rpcRspType|last|export)) -}}
{{- else -}}
	{{- $rpcRspType = (gofulltype $rpcRspType $.FileDescriptor) -}}
{{- end -}}

// {{$rpcName}} ...
func (s *{{$svrName}}ServiceImpl) {{$rpcName}}(ctx context.Context, req *{{$rpcReqType}}, rsp *{{$rpcRspType}}) error {
	// implement business logic here ...
	// ...

	return nil
}

{{end}}
