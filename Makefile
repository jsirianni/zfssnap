.PHONY: build clean gosec lint staticcheck govulncheck test

BINARY_NAME=zfssnap
BUILD_DIR=bin
CMD_DIR=./cmd/zfssnap

VERSION=$(shell cat VERSION)
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "\
	-X github.com/jsirianni/zfssnap/internal/version.semver=$(VERSION) \
	-X github.com/jsirianni/zfssnap/internal/version.commitHash=$(COMMIT_HASH) \
	-X github.com/jsirianni/zfssnap/internal/version.buildTime=$(BUILD_TIME)"

build: build-freebsd-amd64 build-freebsd-arm64 build-linux-amd64 build-linux-arm64

build-freebsd-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=freebsd GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 $(CMD_DIR)

build-freebsd-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=freebsd GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-arm64 $(CMD_DIR)

build-linux-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)

build-linux-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)

test:
	go test -v -race ./...

gosec:
	gosec ./...

lint:
	revive -config revive.toml ./...

staticcheck:
	staticcheck ./...

govulncheck:
	govulncheck ./...

clean:
	rm -rf $(BUILD_DIR)
