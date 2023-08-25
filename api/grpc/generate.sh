#!/usr/bin/env sh

rm -rf ./go && mkdir go

docker build -t proto-builder .
docker run --rm \
        -v "$(pwd)"/go:/gen/go \
        --user "$(id -u):$(id -u)" \
        proto-builder
