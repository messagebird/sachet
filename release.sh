#!/usr/bin/env bash
# a hack to generate releases like other prometheus projects
# use like this: 
#       VERSION=1.0.1 ./release.sh

set -e 

rm -rf "bin/sachet-$VERSION.linux-amd64"
mkdir "bin/sachet-$VERSION.linux-amd64"
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "bin/sachet-$VERSION.linux-amd64/sachet" github.com/messagebird/sachet/cmd/sachet
cd bin
tar -zcvf "sachet-$VERSION.linux-amd64.tar.gz" "sachet-$VERSION.linux-amd64"

GITHUB_USER=messagebird
GITHUB_REPO=sachet

# go get -u github.com/aktau/github-release
# dont forget to set your token like
# export GITHUB_TOKEN=blabla
git tag -a $VERSION -m "version $VERSION"

github-release release \
    --user $GITHUB_USER \
    --repo $GITHUB_REPO \
    --tag $VERSION \
    --name $VERSION \
    --description "version $VERSION!" 

github-release upload \
    --user $GITHUB_USER \
    --repo $GITHUB_REPO \
    --tag $VERSION \
    --name "sachet-$VERSION.linux-amd64.tar.gz" \
    --file "sachet-$VERSION.linux-amd64.tar.gz"


