#!/bin/bash
CGO_ENABLED=1 GOOS=linux CC=x86_64-unknown-linux-gnu-gcc \
    CXX=x86_64-unknown-linux-gnu-g++ GOARCH=amd64 go build -ldflags \
    "-X 'git.uozi.org/uozi/cosy-example-api/settings.buildTime=$(date +%s)'" -o store -v main.go
