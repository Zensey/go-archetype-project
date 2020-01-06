FROM golang:1.13-alpine

RUN apk add --no-cache build-base
RUN apk add --no-cache git

ADD . /app/
WORKDIR /app

RUN make get-deps
RUN make demo
CMD ["/app/demo"]
