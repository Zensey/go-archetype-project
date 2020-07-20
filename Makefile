VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags "-X main.version=$(VERSION)"
GO=$(shell which go)
PWD=$(shell pwd)
ENV=GO111MODULE=on

BINARY = demo
PKG1 = "github.com/Zensey/go-archetype-project/cmd/demo"
report = lint_report.txt

.DEFAULT_GOAL: $(BINARY)
all: get-deps $(BINARY)

get-deps:
	$(ENV) $(GO) get -u ./...

test:
	$(ENV) $(GO) test $(PKG1) -v -run Main

lint:
	$(ENV) go vet ../../../$(PKG1)/*.go

$(BINARY):
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
	docker build -t go-archetype-project .
	docker-compose up

# gen-sw:
# 	docker run --rm --user `id -u`:`id -g` -v ${PWD}:/local jimschubert/swagger-codegen-cli generate \
# 	    -l go \
# 	    -i /local/swagger.yaml \
# 	    -o /local/out

swagger-ui-dev:
	docker run --rm -p 80:8080 -e URL=http://localhost:8080/files/swagger.yaml swaggerapi/swagger-ui

run-dev:
	rm demo ; make demo ; DB_ADDR=localhost:5432 DB_NAME=db DB_USER=db DB_PASSWORD=xxx ERPLY_USERNAME=jnashicq@gmail.com ERPLY_PASSWORD=demo1234 ERPLY_CLIENTCODE=113746 ./demo

