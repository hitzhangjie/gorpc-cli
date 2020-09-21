#!/bin/bash -e

# step-2: check and install required dependencies
which git &> /dev/null || (echo "git not found...install" && exit)
which gofmt &> /dev/null || (echo "gofmt not found...exit" && exit)
which goimports &> /dev/null || (echo "goimports not found...install" && go get golang.org/x/tools/cmd/goimports)
which goi18n &> /dev/null || (echo "goi18n not found...install" && go get -u github.com/nicksnyder/go-i18n/v2/goi18n)
which bindata &> /dev/null || (echo "bindata not found...install" && go get -u github.com/hitzhangjie/bindata)

# step-1: reformat the code
#gofmt -s -w .
#goimports -w -local github.com .

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
rm -rf bindata
tar cvfz install.tgz install
bindata -file install.tgz
rm install.tgz
