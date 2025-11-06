.PHONY: generate build run client test test-rest

PROJECT_NAME := $(shell basename $(git rev-parse --show-toplevel))

# Detect current OS and architecture using Go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Output directory
BUILD_DIR := bin

# Detect .exe extension if Windows target
EXT :=
ifeq ($(GOOS),windows)
	EXT := .exe
endif
# Binary name pattern (e.g., bin/myproject-linux-amd64)
BINARY := $(BUILD_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH)$(EXT)


generate:
	buf generate

build: generate
	@echo "Building $(PROJECT_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY) ./...
	@echo "✅ Built: $(BINARY)"

build-linux:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-windows:
	@$(MAKE) build GOOS=windows GOARCH=amd64

build-arm64:
	@$(MAKE) build GOOS=linux GOARCH=arm64

	# Clean rule
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)

run:
	go run ./cmd/library/

test-rest:
	@echo ":rocket: Running REST enpoints tests"


deps:
# 	go get github.com/swaggo/http-swagger/v2
	go mod tidy
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/flipp-oss/protoc-gen-avro@latest


