FROM golang:alpine AS build

RUN apk update && \
    apk add --no-cache git openssl ca-certificates && \
    go get github.com/messagebird/sachet/cmd/...

FROM alpine
COPY --from=build /go/bin/sachet /usr/local/bin
COPY --chown=nobody examples/config.yaml /etc/sachet/config.yaml
RUN apk update && \
    apk add --no-cache ca-certificates

USER nobody
EXPOSE 9876
ENTRYPOINT ["sachet"]
CMD ["-config", "/etc/sachet/config.yaml"]
