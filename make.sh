#!/bin/bash -e

#step0: update the template version
echo "VERSION=${rev}" > install/VERSION

#step1: update the GoRPCVersion
rev=$(git rev-parse --short HEAD)
cat > config/version.go <<EOF
package config

var GORPCCliVersion string = "${rev}"
EOF

#step2: extract the message
goi18n extract -format json i18n/zh/message_zh.go
mv active.en.json install/active.zh.json
goi18n extract -format json i18n/en/message_en.go
mv active.en.json install/

#step3: compress the templates
rm bindata
tar cvfz install.tgz install
go run util/bindata.go -file install.tgz
rm install.tgz

