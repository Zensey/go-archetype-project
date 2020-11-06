VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
GO=$(shell which go)
PWD=$(shell pwd)
ENV=GO111MODULE=on

BINARY = fiatconv
PKG1 = "github.com/Zensey/go-archetype-project/cmd/fiatconv"
PKGS = $(PKG1)
report = lint_report.txt

.DEFAULT_GOAL: $(BINARY)
all: get-deps $(BINARY)

get-deps:

lint:
	$(ENV) golint $(PKGS) &>> $(report)
	$(ENV) go tool vet ../../$(PKG1)/*.go  &>> $(report)
	$(ENV) errcheck -ignore 'fmt:.*,encoding/binary:.*' -ignoretests $(PKGS)  &>> $(report) ||:
	$(ENV) errcheck $(PKGS)  &>> $(report) ||:
	$(ENV) staticcheck $(PKGS)  &>> $(report) ||:
	$(ENV) unused $(PKGS)  &>> $(report) ||:
	$(ENV) interfacer $(PKGS)  &>> $(report)

$(BINARY):
	$(ENV) $(GO) build -v $(LDFLAGS) ./cmd/fiatconv

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

docker-build:
	docker build -t app-fiatconv-build .
	docker container create --name temp app-fiatconv-build
	docker container cp temp:/app/fiatconv ./
	docker container rm temp
