# Project name (binary name)
BINARY_NAME=snc

# Go command
GO=go

# Build flags (optional: add -ldflags for version info if needed)
BUILD_FLAGS=

.PHONY: all build clean run

all: build

## Format, test and build the binary into the project root
build:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Running tests..."
	$(GO) test ./...
	@echo "Building binary..."
	$(GO) build $(BUILD_FLAGS) -o $(BINARY_NAME) ./cmd/src


## Remove the built binary
clean:
	rm -f $(BINARY_NAME)

## Run unit tests with verbose output
test:
	$(GO) test -v ./...