#!/usr/bin/env bash
#set -x
#VERSION=$(git rev-parse HEAD)
#BuildTime=$(date '+%Y-%m-%d')
App="cnb-code-import"
TAG=${App}:latest

docker build -t ${TAG} .






