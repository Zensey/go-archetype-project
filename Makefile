VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
SHELL=/bin/bash
GO=$(shell which go)
ENV=$(shell pwd)/build/env.sh
GOBIN=$(shell pwd)/build/_workspace/bin
PWD=$(shell pwd)

BINARY1 = api
BINARY2 = worker1
BINARY3 = worker2
PKG1 = "github.com/Zensey/go-archetype-project/cmd/api"
PKG2 = "github.com/Zensey/go-archetype-project/pkg/logger"
PKG3 = "github.com/Zensey/go-archetype-project/pkg/utils"
PKGS = $(PKG1) $(PKG2)
report = lint_report.txt

.DEFAULT_GOAL: all
all: get-deps $(BINARY1) $(BINARY2) $(BINARY3) test

get-deps:
	$(ENV) $(GO) get -u github.com/golang/dep/cmd/dep
	#$(ENV) $(GO) get -u -a golang.org/x/tools/cmd/stringer
	#$(ENV) $(GO) get -u github.com/golang/lint/golint
	#$(ENV) $(GO) get -u github.com/kisielk/errcheck
	#$(ENV) $(GO) get -u honnef.co/go/tools/cmd/staticcheck
	#$(ENV) $(GO) get -u honnef.co/go/tools/cmd/unused
	#$(ENV) $(GO) get -u mvdan.cc/interfacer

	$(ENV) $(GOBIN)/dep ensure -v

test:
	$(ENV) $(GO) test $(PKG3) -v

lint:
	$(ENV) golint $(PKGS)  &>> $(report)
	$(ENV) go tool vet ../../$(PKG1)/*.go  &>> $(report)
	$(ENV) go tool vet ../../$(PKG2)/*.go  &>> $(report)
	$(ENV) errcheck -ignore 'fmt:.*,encoding/binary:.*' -ignoretests $(PKGS)  &>> $(report) ||:
	$(ENV) errcheck $(PKGS)  &>> $(report) ||:
	$(ENV) staticcheck $(PKGS)  &>> $(report) ||:
	$(ENV) unused $(PKGS)  &>> $(report) ||:
	$(ENV) interfacer $(PKGS)  &>> $(report)

$(BINARY1):
	$(ENV) $(GO) install -v $(LDFLAGS) ./cmd/$(BINARY1)
$(BINARY2):
	$(ENV) $(GO) install -v $(LDFLAGS) ./cmd/$(BINARY2)
$(BINARY3):
	$(ENV) $(GO) install -v $(LDFLAGS) ./cmd/$(BINARY3)

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*
	rm -fr vendor

strip-$(BINARY): $(BINARY)
	strip -s $(GOBIN)/$(BINARY)

docker-build:
	docker build -t go-archetype-project .
	docker run -h docker.local --rm --publish 8888:8888 --publish 5432:5432 -it go-archetype-project
