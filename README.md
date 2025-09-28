# SNC

A fast and reliable file synchronization tool written in Go. SNC efficiently synchronizes files between source and target directories with support for different update detection strategies and optional cleanup of missing files.

## Features

- **Fast synchronization** with configurable update detection methods
- **Two update strategies**:
  - `modtime`: Fast detection using file modification time and size (default)
  - `sha256`: Reliable detection using SHA256 checksums
- **Optional cleanup** of files that exist in target but not in source
- **Comprehensive logging** with configurable log levels
- **Error handling** with detailed error reporting
- **Directory validation** before synchronization

## Installation

### Prerequisites

- Go 1.24.4 or later

### Build from source

```bash
# Clone the repository
git clone <repository-url>
cd snc

# Build the binary
make build

# Or build manually
go build -o snc ./cmd/src
```

## Usage

```bash
snc [OPTIONS] <source> <target>
```

### Options

- `--delete-missing`: Delete files from target that do not exist in source (default: false)
- `--log-level LEVEL`: Set logging level - error, warn, info, debug (default: info)
- `--update-method METHOD`: Method for detecting file updates - modtime, sha256 (default: modtime)

### Arguments

- `source`: Source directory path
- `target`: Target directory path

## Examples

### Basic synchronization

```bash
# Sync files from source to target directory
./snc /path/to/source /path/to/target

# With verbose logging
./snc --log-level debug /path/to/source /path/to/target
```

### Synchronization with cleanup

```bash
# Sync and remove files that don't exist in source
./snc --delete-missing /path/to/source /path/to/target
```

### Using SHA256 for reliable detection

```bash
# Use SHA256 checksums for update detection (slower but more reliable)
./snc --update-method sha256 /path/to/source /path/to/target
```

### Complete example with all options

```bash
# Full sync with cleanup, debug logging, and SHA256 detection
./snc --delete-missing --log-level debug --update-method sha256 /home/user/documents /backup/documents
```

## Development

### Running tests

```bash
# Run all tests
make test

# Or run tests manually
go test -v ./...
```

### Building

```bash
# Format, test, and build
make build

# Clean build artifacts
make clean
```

### Project structure

```
snc/
├── cmd/src/main.go          # Main application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── errors/              # Error handling and types
│   ├── logger/              # Logging utilities
│   ├── stream/              # File synchronization logic
│   ├── synchronizer/        # Main synchronization orchestrator
│   └── validate/dir/        # Directory validation
├── go.mod                   # Go module definition
└── Makefile                 # Build automation
```

## Update Strategies

### ModTime Strategy (Default)

- **Speed**: Very fast
- **Reliability**: Good for most cases
- **Use case**: General file synchronization
- **Detection**: File size and modification time

### SHA256 Strategy

- **Speed**: Slower (reads entire file content)
- **Reliability**: Highly reliable
- **Use case**: Critical data synchronization
- **Detection**: SHA256 checksum comparison
