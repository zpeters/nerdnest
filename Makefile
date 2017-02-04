#
#  Makefile for Go
#
SHELL=/usr/bin/env bash
VERSION=$(shell git describe --tags --always)

default: build

build:
	go build

clean:
	rm -f nerdnest
	rm -f bin/nerdtest-*

cross:
	echo "Building darwin-amd64..."
	GOOS="darwin" GOARCH="amd64" go build -o bin/nerdnest-mac-amd64-${VERSION}

	echo "Building windows-amd64..."
	GOOS="windows" GOARCH="amd64" go build -o bin/nerdnest-64-${VERSION}.exe

	echo "Building linux-arm..."
	GOOS="linux" GOARCH="arm" go  build -o bin/nerdnest-linux-arm-${VERSION}

	echo "Building linux-amd64..."
	GOOS="linux" GOARCH="amd64" go build -o bin/nerdnest-linux-amd64-${VERSION}
