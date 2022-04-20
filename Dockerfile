FROM golang:1.18 AS builder

WORKDIR /build

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod vendor -o sachet github.com/messagebird/sachet/cmd/sachet

FROM alpine:3.15

COPY --from=builder /build/sachet /usr/local/bin
COPY --chown=nobody examples/config.yaml /etc/sachet/config.yaml
RUN apk update && \
    apk add --no-cache ca-certificates

USER nobody
EXPOSE 9876
ENTRYPOINT ["sachet"]
CMD ["-config", "/etc/sachet/config.yaml"]
