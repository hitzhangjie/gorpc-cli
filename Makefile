user	:=	$(shell whoami)
rev 	:= 	$(shell git rev-parse --short HEAD)

# GOBIN > GOPATH > INSTALLDIR
GOBIN	:=	$(shell echo ${GOBIN} | cut -d':' -f1)
GOPATH	:=	$(shell echo $(GOPATH) | cut -d':' -f1)
BIN		:= 	""

# check GOBIN
ifneq ($(GOBIN),)
	BIN=$(GOBIN)
else
# check GOPATH
ifneq ($(GOPATH),)
	BIN=$(GOPATH)/bin
else
# check INSTALL
#ifneq ($(INSTALL),)
#	BIN=$(INSTALL)
#endif
endif
endif

all: *.go
	@GO111MODULE=on go build -o gorpc -ldflags="-X github.com/hitzhangjie/gorpc/config.GORPCCliVersion=$(rev)"

experimental: *.go
	@GO111MODULE=on go build -o gorpc -tags=experimental -ldflags="-X github.com/hitzhangjie/gorpc/config.GORPCCliVersion=$(rev)"

debug: *.go
	@GO111MODULE=on go build -o gorpc -gcflags="all=-N -l" -ldflags="-X github.com/hitzhangjie/gorpc/config.GORPCCliVersion=$(rev)"

.PHONY: clean
.PHONY: install
.PHONY: uninstall

install:
ifeq ($(user),root)
#root, install for all user
	@cp ./gorpc /usr/bin
	@[ -d /etc/gorpc ] || mkdir /etc/gorpc
	@cp -rf ./install/* /etc/gorpc/
	@go get -u github.com/golang/mock/mockgen
else
#!root, install for current user
	$(shell if [ -z $(BIN) ]; then read -p "Please select installdir: " REPLY; mkdir -p $${REPLY}; cp ./gorpc $${REPLY}/; else mkdir -p $(BIN); cp ./gorpc $(BIN); fi)
	@[ -d ~/.gorpc ] || mkdir ~/.gorpc
	@cp -rf ./install/* ~/.gorpc/
	@which mockgen &> /dev/null || go get github.com/golang/mock/mockgen
endif
	@echo "install finished"

uninstall:
ifeq ($(user),root)
#root, install for all user
	@rm -f /usr/bin/gorpc &> /dev/null
	@rm -rf /etc/gorpc &>/dev/null
else
#!root, install for current user
	$(shell for i in `which -a gorpc | grep -v '/usr/bin/gorpc' 2>/dev/null | sort | uniq`; do read -p "Press to remove $${i} (y/n): " REPLY; if [ $${REPLY} = "y" ]; then rm -f $${i}; fi; done)
	@rm -rf ~/.gorpc &>/dev/null
endif
	@echo "uninstall finished"

clean:
	@rm -f ./gorpc
	@echo "clean finished"

fmt:
	@gofmt -s -w .
	@goimports -w -local github.com .

static:
	@rm bindata
	@tar cvfz install.tgz install
	@go run util/bindata.go -file install.tgz
	@rm install.tgz
