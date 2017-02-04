#!/usr/bin/env bash
VERSION=`git describe --tags --always`

echo "Building darwin-amd64..."
GOOS="darwin" GOARCH="amd64" go build -o bin/nerdnest-mac-amd64-${VERSION}

echo "Building windows-amd64..."
GOOS="windows" GOARCH="amd64" go build -o bin/nerdnest-64-${VERSION}.exe

echo "Building linux-arm..."
GOOS="linux" GOARCH="arm" go  build -o bin/nerdnest-linux-arm-${VERSION}

echo "Building linux-amd64..."
GOOS="linux" GOARCH="amd64" go build -o bin/nerdnest-linux-amd64-${VERSION}
