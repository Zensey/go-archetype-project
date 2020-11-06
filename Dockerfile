FROM golang:1.14-alpine

RUN apk add --no-cache \
	build-base git

ADD . /app/
WORKDIR /app
RUN make fiatconv
CMD ["/app/fiatconv"]
