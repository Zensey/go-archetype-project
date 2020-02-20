VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
GO=$(shell which go)
PWD=$(shell pwd)
ENV=GO111MODULE=on

BINARY = demo
PKG1 = "github.com/Zensey/go-archetype-project/cmd/demo"
PKG2 = "github.com/Zensey/go-archetype-project/pkg/logger"
report = lint_report.txt

.DEFAULT_GOAL: $(BINARY)
all: get-deps $(BINARY)

get-deps:
	$(ENV) $(GO) get -u -a golang.org/x/tools/cmd/stringer
	$(ENV) $(GO) get -u github.com/go-pg/pg/v9
	$(ENV) $(GO) get -u ./...

test:
	$(ENV) $(GO) test $(PKG1) -v -run Main

lint:
	$(ENV) go vet ../../../$(PKG1)/*.go
	$(ENV) go vet ../../../$(PKG2)/*.go

$(BINARY):
	#$(ENV) $(GO) generate ./pkg/logger
	$(ENV) $(GO) build -v $(LDFLAGS) ./cmd/demo

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

docker-db-reset:
	docker-compose down; docker volume rm go-archetype-project_postgresql go-archetype-project_postgresql_data

docker-db-shell:
	docker-compose exec db psql -U db

docker-build:
	#docker build -t go-archetype-project .
	docker-compose up
