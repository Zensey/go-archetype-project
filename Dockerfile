FROM golang:1.15-alpine as stage0

RUN apk add --no-cache \
	build-base \
	git

ADD . /app/
WORKDIR /app
RUN make get-deps customer_svc

FROM alpine as release

COPY --from=stage0 /app/customer_svc /app/customer_svc
COPY --from=stage0 /app/config.yaml /app/config.yaml
COPY --from=stage0 /app/tpl /app/tpl

WORKDIR /app
CMD ["/app/customer_svc", "serve", "--config", "/app/config.yaml"]
