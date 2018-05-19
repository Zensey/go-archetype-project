FROM golang:latest
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
ADD . /app/
WORKDIR /app
RUN make demo
RUN make test
CMD ["/app/build/_workspace/bin/demo"]
