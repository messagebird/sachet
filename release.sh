#!/usr/bin/env bash
# a hack to generate releases like other prometheus projects
# use like this: 
#       VERSION=1.0.1 ./release.sh

set -e 

# github user and repo
USER=messagebird
REPO=sachet

rm -rf "bin/$REPO-$VERSION.linux-amd64"
mkdir "bin/$REPO-$VERSION.linux-amd64"
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "bin/sachet-$VERSION.linux-amd64/$REPO" github.com/$USER/$REPO/cmd/$REPO
cd bin
tar -zcvf "$REPO-$VERSION.linux-amd64.tar.gz" "$REPO-$VERSION.linux-amd64"

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
    --name "$REPO-$VERSION.linux-amd64.tar.gz" \
    --file "$REPO-$VERSION.linux-amd64.tar.gz"


