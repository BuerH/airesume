BINARY := airesume
PKG := ./cmd/airesume
DIST := dist
VERSION ?= dev
GOFLAGS ?= -buildvcs=false

.PHONY: all build test check fmt clean release-snapshot

all: check build

build:
	go build $(GOFLAGS) -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/$(BINARY) $(PKG)

test:
	go test $(GOFLAGS) ./...

fmt:
	gofmt -w ./cmd ./internal

check:
	gofmt -w ./cmd ./internal
	go test $(GOFLAGS) ./...

clean:
	rm -rf bin $(DIST)

release-snapshot:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(BINARY)-linux-amd64 $(PKG)
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(BINARY)-linux-arm64 $(PKG)
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(BINARY)-darwin-amd64 $(PKG)
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(BINARY)-darwin-arm64 $(PKG)
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o $(DIST)/$(BINARY)-windows-amd64.exe $(PKG)
	cd $(DIST) && sha256sum * > checksums.txt
