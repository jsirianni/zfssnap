.PHONY: build clean gosec lint staticcheck govulncheck test

build: build-freebsd-amd64 build-freebsd-arm64 build-linux-amd64 build-linux-arm64

build-freebsd-amd64:
	@sh scripts/build-binary.sh freebsd amd64

build-freebsd-arm64:
	@sh scripts/build-binary.sh freebsd arm64

build-linux-amd64:
	@sh scripts/build-binary.sh linux amd64

build-linux-arm64:
	@sh scripts/build-binary.sh linux arm64

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
	@sh scripts/clean.sh
