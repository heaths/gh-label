#!/bin/bash

# Based on https://github.com/vilmibm/gh-user-status/blob/cead3abf46ffb5fd3c178a0ba6f2c69c3dbabf7e/build.sh

mkdir -p builds
GOOS=darwin GOARCH=amd64 go build -o builds/darwin-x86_64
GOOS=linux GOARCH=386 go build -o builds/linux-i386
GOOS=linux GOARCH=amd64 go build -o builds/linux-x86_64
GOOS=windows GOARCH=386 go build -o builds/windows-i386
GOOS=windows GOARCH=amd64 go build -o builds/windows-x86_64
