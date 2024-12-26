get-deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u honnef.co/go/tools/cmd/unused
	go get -u mvdan.cc/interfacer

client:
	go build ./cmd/client

server:
	go build ./cmd/server
