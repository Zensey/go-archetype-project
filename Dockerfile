FROM golang:latest
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
ADD . /app/
WORKDIR /app
RUN make get-deps
RUN make test
RUN make demo
CMD ["/app/build/_workspace/bin/demo"]
