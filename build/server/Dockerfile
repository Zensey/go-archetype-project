FROM golang:1.21-alpine AS builder

RUN apk add --no-cache make git build-base
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=1 go test -race -v ./...
RUN CGO_ENABLED=0 make server

FROM scratch
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/assets/quotes.yml /app/assets/quotes.yml
CMD ["/app/server"]
