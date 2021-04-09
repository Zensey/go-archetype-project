GO=$(shell which go)
PWD=$(shell pwd)
ENV=GO111MODULE=on
BINARY = game


.DEFAULT_GOAL: $(BINARY)
all: $(BINARY)

$(BINARY):
	$(ENV) $(GO) build -o $(BINARY)

clean:
	rm $(BINARY)
