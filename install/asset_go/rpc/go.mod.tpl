{{- $pkgName := .PackageName -}}
{{- $goPkgOption := "" -}}
{{- with .FileOptions.go_package -}}
  {{- $goPkgOption = . -}}
{{- end -}}

{{- if ne $goPkgOption "" -}}
module {{$goPkgOption}}
{{- else -}}
module {{$pkgName}}
{{- end }}

go 1.12
