VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
GO=$(shell which go)
PWD=$(shell pwd)
ENV=GO111MODULE=on

BINARY = demo
PKG1 = "github.com/Zensey/go-archetype-project/cmd/demo"
PKGS = $(PKG1)
report = lint_report.txt

.DEFAULT_GOAL: $(BINARY)
all: get-deps $(BINARY)

get-deps:
	$(ENV) $(GO) get -u github.com/golang/lint/golint
	$(ENV) $(GO) get -u github.com/kisielk/errcheck
	$(ENV) $(GO) get -u honnef.co/go/tools/cmd/staticcheck
	$(ENV) $(GO) get -u honnef.co/go/tools/cmd/unused
	$(ENV) $(GO) get -u mvdan.cc/interfacer

test:
	$(ENV) $(GO) test $(PKG1) -v -run Main

lint:
	$(ENV) golint $(PKGS) &>> $(report)
	$(ENV) go tool vet ../../$(PKG1)/*.go  &>> $(report)
	$(ENV) go tool vet ../../$(PKG2)/*.go  &>> $(report)
	$(ENV) errcheck -ignore 'fmt:.*,encoding/binary:.*' -ignoretests $(PKGS)  &>> $(report) ||:
	$(ENV) errcheck $(PKGS)  &>> $(report) ||:
	$(ENV) staticcheck $(PKGS)  &>> $(report) ||:
	$(ENV) unused $(PKGS)  &>> $(report) ||:
	$(ENV) interfacer $(PKGS)  &>> $(report)

$(BINARY):
	$(ENV) $(GO) generate ./pkg/logger
	$(ENV) $(GO) build -v $(LDFLAGS) ./cmd/demo

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

docker-build:
	docker build -t go-archetype-project .
	docker run --rm --publish 8080:8080 -it go-archetype-project