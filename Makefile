.PHONY: build test lint vet fmt clean

BINARY := algod-monitor
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/algod-monitor

test:
	go test -v -race ./...

lint: vet fmt
	@echo "Lint passed."

vet:
	go vet ./...

fmt:
	@test -z "$$(gofmt -l .)" || (echo "Run gofmt:"; gofmt -l .; exit 1)

clean:
	rm -rf $(BUILD_DIR)
