# build details
DEFAULT_GOOS = mac
GOARCH = amd64
BUILD_DIR = bin
BINARY = nerdnest


# version stuff
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# commands
all: dep vet fmt build

cross: dep vet fmt mac linux windows

dep:
	dep ensure

build: ${DEFAULT_GOOS}

fmt:
	goimports -w ./cmd/
	goimports -w ./internal/

vet:
	go vet ./cmd/...
	go vet ./internal/...

clean:
	rm -f nerdnest
	rm -f bin/*

test:
	go test ./cmd/...
	go test ./internal/...

### build for different platforms
linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-linux-${GOARCH}-${VERSION} ./cmd/...
mac:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-darwin-${GOARCH}-${VERSION} ./cmd/...
windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-windows-${GOARCH}-${VERSION} ./cmd/...


### phone to speed up performance
.PHONY: dep build fmt vet clean
