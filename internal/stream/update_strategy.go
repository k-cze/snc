package stream

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// UpdateStrategy defines the interface for different file update detection methods
// Different strategies offer different trade-offs between speed and reliability
type UpdateStrategy interface {
	NeedsUpdate(srcPath, dstPath string) (bool, error)
	Name() string
}

// ModTimeStrategy uses file modification time and size for update detection
//
// Pros:
//   - Very fast (no file content reading required)
//   - Low CPU and I/O usage
//   - Good for most common use cases
//
// Cons:
//   - Less reliable (files can have same modtime/size but different content)
//   - May miss updates if file timestamps are not preserved
//   - Not suitable for files that are frequently modified with same timestamp
//
// This is the default strategy for backward compatibility and performance
type ModTimeStrategy struct{}

func (m *ModTimeStrategy) Name() string {
	return "modtime"
}

func (m *ModTimeStrategy) NeedsUpdate(srcPath, dstPath string) (bool, error) {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return false, fmt.Errorf("cannot stat source file %s: %w", srcPath, err)
	}

	dstInfo, err := os.Stat(dstPath)
	if err != nil {
		return false, fmt.Errorf("cannot stat destination file %s: %w", dstPath, err)
	}

	// Simple check: size or modtime differs
	if srcInfo.Size() != dstInfo.Size() {
		return true, nil
	}
	if !srcInfo.ModTime().Equal(dstInfo.ModTime()) {
		return true, nil
	}
	return false, nil
}

// SHA256Strategy uses SHA256 checksums for update detection
//
// Pros:
//   - Highly reliable (guaranteed to detect any content changes)
//   - Cryptographically secure hash function
//   - Works regardless of file timestamps or metadata
//   - Suitable for critical data synchronization
//
// Cons:
//   - Slower than modtime strategy (requires reading entire file content)
//   - Higher CPU usage for large files
//   - Higher I/O usage (must read both source and destination files)
//
// Recommended for critical data or when file timestamps cannot be trusted
type SHA256Strategy struct{}

func (s *SHA256Strategy) Name() string {
	return "sha256"
}

func (s *SHA256Strategy) NeedsUpdate(srcPath, dstPath string) (bool, error) {
	srcHash, err := calculateSHA256(srcPath)
	if err != nil {
		return false, fmt.Errorf("cannot calculate SHA256 for source file %s: %w", srcPath, err)
	}

	dstHash, err := calculateSHA256(dstPath)
	if err != nil {
		return false, fmt.Errorf("cannot calculate SHA256 for destination file %s: %w", dstPath, err)
	}

	return srcHash != dstHash, nil
}

// calculateSHA256 calculates the SHA256 hash of a file
func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// NewUpdateStrategy creates an UpdateStrategy based on the method name
//
// Supported methods:
//   - "modtime": Fast but less reliable (default)
//   - "sha256":  Slower but highly reliable
//
// The modtime strategy is recommended for most use cases due to its speed,
// while sha256 is recommended for critical data synchronization where
// reliability is more important than performance.
func NewUpdateStrategy(method string) (UpdateStrategy, error) {
	switch method {
	case "modtime":
		return &ModTimeStrategy{}, nil
	case "sha256":
		return &SHA256Strategy{}, nil
	default:
		return nil, fmt.Errorf("unsupported update method: %s (supported: modtime, sha256)", method)
	}
}
