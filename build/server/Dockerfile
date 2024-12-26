FROM golang:1.21-alpine AS builder

RUN apk add --no-cache make git
ADD . /app/
WORKDIR /app
RUN make server

FROM scratch
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/quotes.yml /app/quotes.yml
CMD ["/app/server"]
