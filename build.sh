#!/bin/sh

apk add --no-cache git
go get -v -d
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-w -s" -o release/minio-delayed-server
