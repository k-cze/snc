package stream

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestModTimeStrategy(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")

	// Create test files
	createTestFile(t, srcFile, "test content")
	createTestFile(t, dstFile, "test content")

	strategy := &ModTimeStrategy{}

	// Test 1: Same content, same modtime - should not need update
	needsUpdate, err := strategy.NeedsUpdate(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if needsUpdate {
		t.Error("Expected no update needed for identical files")
	}

	// Test 2: Different content, same modtime - should need update
	createTestFile(t, dstFile, "different content")
	needsUpdate, err = strategy.NeedsUpdate(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !needsUpdate {
		t.Error("Expected update needed for different content")
	}

	// Test 3: Same content, different modtime - should need update
	createTestFile(t, dstFile, "test content")
	// Touch the source file to change modtime
	os.Chtimes(srcFile, time.Now(), time.Now().Add(time.Hour))
	needsUpdate, err = strategy.NeedsUpdate(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !needsUpdate {
		t.Error("Expected update needed for different modtime")
	}

	// Test 4: Non-existent source file
	needsUpdate, err = strategy.NeedsUpdate("nonexistent.txt", dstFile)
	if err == nil {
		t.Error("Expected error for non-existent source file")
	}

	// Test 5: Non-existent destination file
	needsUpdate, err = strategy.NeedsUpdate(srcFile, "nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent destination file")
	}

	// Test 6: Name method
	if strategy.Name() != "modtime" {
		t.Errorf("Expected name 'modtime', got '%s'", strategy.Name())
	}
}

func TestSHA256Strategy(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "destination.txt")

	// Create test files
	createTestFile(t, srcFile, "test content")
	createTestFile(t, dstFile, "test content")

	strategy := &SHA256Strategy{}

	// Test 1: Same content - should not need update
	needsUpdate, err := strategy.NeedsUpdate(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if needsUpdate {
		t.Error("Expected no update needed for identical files")
	}

	// Test 2: Different content - should need update
	createTestFile(t, dstFile, "different content")
	needsUpdate, err = strategy.NeedsUpdate(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !needsUpdate {
		t.Error("Expected update needed for different content")
	}

	// Test 3: Same content, different modtime - should not need update
	createTestFile(t, dstFile, "test content")
	os.Chtimes(srcFile, time.Now(), time.Now().Add(time.Hour))
	needsUpdate, err = strategy.NeedsUpdate(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if needsUpdate {
		t.Error("Expected no update needed for same content despite different modtime")
	}

	// Test 4: Non-existent source file
	needsUpdate, err = strategy.NeedsUpdate("nonexistent.txt", dstFile)
	if err == nil {
		t.Error("Expected error for non-existent source file")
	}

	// Test 5: Non-existent destination file
	needsUpdate, err = strategy.NeedsUpdate(srcFile, "nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent destination file")
	}

	// Test 6: Name method
	if strategy.Name() != "sha256" {
		t.Errorf("Expected name 'sha256', got '%s'", strategy.Name())
	}
}

func TestNewUpdateStrategy(t *testing.T) {
	tests := []struct {
		method    string
		wantName  string
		wantError bool
	}{
		{"modtime", "modtime", false},
		{"sha256", "sha256", false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			strategy, err := NewUpdateStrategy(tt.method)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if strategy.Name() != tt.wantName {
				t.Errorf("Expected name '%s', got '%s'", tt.wantName, strategy.Name())
			}
		})
	}
}

func TestCalculateSHA256(t *testing.T) {
	// Create temporary test file
	tempFile, err := os.CreateTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write test content
	testContent := "test content for sha256"
	_, err = tempFile.WriteString(testContent)
	if err != nil {
		t.Fatalf("Failed to write test content: %v", err)
	}
	tempFile.Close()

	// Calculate hash
	hash, err := calculateSHA256(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to calculate SHA256: %v", err)
	}

	// Hash should not be empty
	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	// Hash should be consistent
	hash2, err := calculateSHA256(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to calculate SHA256 second time: %v", err)
	}

	if hash != hash2 {
		t.Error("Hash should be consistent for same content")
	}

	// Test with non-existent file
	_, err = calculateSHA256("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// Helper function to create test files
func createTestFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", path, err)
	}
}
