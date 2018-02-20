FROM golang:1.9.4-alpine3.7
RUN apk add --no-cache --virtual git && \
    go-wrapper download github.com/marcelcorso/sachet/cmd/... && \
    go-wrapper install github.com/marcelcorso/sachet/cmd/... && \
    rm -rf src pkg && \
    apk del git

EXPOSE 9876
ENTRYPOINT ["sachet"]
