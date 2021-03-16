VERSION=$(shell git describe --tags --always)
LDFLAGS=-tags netgo -ldflags \
    "-X github.com/Zensey/go-archetype-project/pkg/driver/config.Version=$(VERSION) \
    -X github.com/Zensey/go-archetype-project/pkg/driver/config.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` \
    -X github.com/Zensey/go-archetype-project/pkg/driver/config.Commit=`git rev-parse HEAD`"


GO=$(shell which go)
PWD=$(shell pwd)
ENV=GO111MODULE=on
GOBIN=GOBIN=$(PWD)/.bin/

export GO111MODULE := on
export PATH := .bin:${PATH}
export PWD := $(shell pwd)

BINARY = customer_svc

.DEFAULT_GOAL: $(BINARY)
all: get-deps $(BINARY)

get-deps:
	$(ENV) $(GOBIN) $(GO) get -u github.com/markbates/pkger/cmd/pkger
	$(ENV) $(GOBIN) $(GO) get -u github.com/gobuffalo/packr/v2/packr2

.PHONY: pack
pack:
	mkdir -p generated
	.bin/pkger -o generated -exclude .git -exclude .github -exclude .idea -exclude .bin -exclude cmd
	.bin/packr2

$(BINARY): pack
	$(ENV) $(GO) build -o $(BINARY) -v $(LDFLAGS)

test:
	$(ENV) $(GO) test -count=1 github.com/Zensey/go-archetype-project/pkg/customer -v

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

docker-build:
	docker build -t customers .
	docker run --rm --publish 3010:3010 --env-file .env -it customers

start-db-local:
	docker-compose start

start:
	env $(cat .env.serve | xargs) ./$(BINARY) serve --config config.yaml