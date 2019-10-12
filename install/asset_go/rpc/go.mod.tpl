{{- with .FileOptions.go_package -}}
module {{.}}
{{- else -}}
module {{$.PackageName}}
{{- end }}

go 1.12

