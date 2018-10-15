FROM golang:1.12.9-alpine3.9

RUN apk add --no-cache \
	build-base

ADD . /app/
WORKDIR /app
RUN make test
RUN make demo
CMD ["/app/demo"]
