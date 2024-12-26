FROM golang:1.21-alpine AS builder

RUN apk add --no-cache make git
ADD . /app/
WORKDIR /app
RUN make client

FROM scratch
COPY --from=builder /app/client /app/client
CMD ["/app/client"]
