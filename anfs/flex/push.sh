#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o driver github.com/andyleap/anfs/anfs
docker build -t andyleap/anfs-flex-install:latest .
docker push andyleap/anfs-flex-install:latest

