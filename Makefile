GOPATH=$(shell pwd)
GOROOT=$(shell echo ~)/opt/go
GO=$(shell which go)

VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
SHELL=/bin/bash

BINARY = go.archetype
MAIN_PKG = zensey.go.archetype
report = lint_report.txt

.DEFAULT_GOAL: $(BINARY)

all: $(BINARY)

get-deps:
	export GOPATH=$(GOPATH)
	#export GOROOT=$(GOROOT)

	$(GO) get -u -a "golang.org/x/tools/cmd/stringer"

	$(GO) get -u github.com/golang/lint/golint
	$(GO) get -u github.com/kisielk/errcheck
	$(GO) get -u honnef.co/go/tools/cmd/staticcheck
	$(GO) get -u honnef.co/go/tools/cmd/unused
	$(GO) get -u mvdan.cc/interfacer

lint:
	export GOPATH=$(GOPATH)

	bin/golint src/$(MAIN_PKG) &>> $(report)
	go tool vet src/$(MAIN_PKG)/*.go &>> $(report)
	bin/errcheck -ignore 'fmt:.*,encoding/binary:.*' -ignoretests cloudcc/ &>> $(report) ||:
	bin/errcheck $(MAIN_PKG)/ &>> $(report) ||:
	bin/staticcheck $(MAIN_PKG) &>> $(report) ||:
	bin/unused $(MAIN_PKG) &>> $(report) ||:
	bin/interfacer $(MAIN_PKG)

$(BINARY):
	export GOPATH=$(GOPATH)

	$(GO) generate $(MAIN_PKG)
	$(GO) build -o $(BINARY) $(MAIN_PKG)

clean:
	rm -rf $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

