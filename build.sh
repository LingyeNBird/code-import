#!/usr/bin/env bash
#set -x

VERSION=$(git rev-parse HEAD)
BuildTime=$(date '+%Y-%m-%d %H:%M:%S')
#LDFAGS="-X 'ccrctl/cmd.BuildTime=${BuildTime}' -X 'ccrctl/cmd.Version=${VERSION}' "
#CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build  -mod=vendor -ldflags "-X 'main.BuildTime=$(date '+%Y-%m-%d %H:%M:%S')' -X 'main.Version=${VERSION}' -o ccrctl_${VERSION}_linux_amd64 ./

CGO_ENABLED=0  GOOS=linux  GOARCH=amd64 go build -mod=vendor -ldflags "-X 'ccrctl/cmd.BuildTime=$BuildTime' -X 'ccrctl/cmd.Version=${VERSION}' " -o  ccrctl_linux_amd64 ./
CGO_ENABLED=0  GOOS=darwin  GOARCH=amd64 go build -mod=vendor -ldflags "-X 'ccrctl/cmd.BuildTime=$BuildTime' -X 'ccrctl/cmd.Version=${VERSION}' " -o  ccrctl_darwin_amd64 ./
