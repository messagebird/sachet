FROM golang:alpine

RUN apk add --no-cache --virtual git && \
    go-wrapper download github.com/marcelcorso/sachet && \
    go-wrapper install github.com/marcelcorso/sachet && \
    rm -rf src pkg && \
    apk del git

COPY example-config.yaml /etc/sachet/config.yaml

EXPOSE 9876
ENTRYPOINT ["sachet"]
CMD ["-config", "/etc/sachet/config.yaml"]