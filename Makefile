get-deps:
	go get -u github.com/golang/lint/golint
	go install github.com/kisielk/errcheck@latest
	go get -u honnef.co/go/tools

txparser:
	go build ./cmd/txparser
