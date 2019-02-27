VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
SHELL=/bin/bash
GO=$(shell which go)
ENV=$(shell pwd)/build/env.sh
GOBIN=$(shell pwd)/build/_workspace/bin
PWD=$(shell pwd)

BINARY = demo
PKG1 = "github.com/Zensey/go-archetype-project/cmd/demo"
PKG2 = "github.com/Zensey/go-archetype-project/pkg/logger"
PKGS = $(PKG1) $(PKG2)
report = lint_report.txt

.DEFAULT_GOAL: $(BINARY)
all: get-deps $(BINARY)

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
	GOCACHE=off $(ENV) $(GO) test $(PKG1) -v

lint:
	$(ENV) golint $(PKGS)  &>> $(report)
	$(ENV) go tool vet ../../$(PKG1)/*.go  &>> $(report)
	$(ENV) go tool vet ../../$(PKG2)/*.go  &>> $(report)
	$(ENV) errcheck -ignore 'fmt:.*,encoding/binary:.*' -ignoretests $(PKGS)  &>> $(report) ||:
	$(ENV) errcheck $(PKGS)  &>> $(report) ||:
	$(ENV) staticcheck $(PKGS)  &>> $(report) ||:
	$(ENV) unused $(PKGS)  &>> $(report) ||:
	$(ENV) interfacer $(PKGS)  &>> $(report)

$(BINARY):
	$(ENV) $(GO) generate "github.com/Zensey/go-archetype-project/pkg/logger"
	$(ENV) $(GO) install -v $(LDFLAGS) ./cmd/$(BINARY)

clean:
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

strip-$(BINARY): $(BINARY)
	strip -s $(GOBIN)/$(BINARY)

docker-build:
	docker build -t go-archetype-project .
	docker run --rm --publish 8080:8080 -it go-archetype-project