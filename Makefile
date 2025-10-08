.PHONY: build clean gosec lint staticcheck govulncheck test

BINARY_NAME=zfssnap
BUILD_DIR=bin
CMD_DIR=./cmd/zfssnap

build: build-freebsd-amd64 build-freebsd-arm64 build-linux-amd64 build-linux-arm64

build-freebsd-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=freebsd GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 $(CMD_DIR)

build-freebsd-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=freebsd GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-arm64 $(CMD_DIR)

build-linux-amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)

build-linux-arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)

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
