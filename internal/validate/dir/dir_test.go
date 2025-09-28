package dir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateDir(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		path        string
		opts        []configOption
		expectError bool
		setup       func() // function to set up test state
	}{
		{
			name:        "existing directory",
			path:        tempDir,
			opts:        []configOption{},
			expectError: false,
			setup:       func() {},
		},
		{
			name:        "non-existent directory without allowCreate",
			path:        filepath.Join(tempDir, "nonexistent"),
			opts:        []configOption{},
			expectError: true,
			setup:       func() {},
		},
		{
			name:        "non-existent directory with allowCreate",
			path:        filepath.Join(tempDir, "newdir"),
			opts:        []configOption{withAllowCreate()},
			expectError: false,
			setup:       func() {},
		},
		{
			name:        "file instead of directory",
			path:        filepath.Join(tempDir, "file.txt"),
			opts:        []configOption{},
			expectError: true,
			setup: func() {
				// Create a file
				file, err := os.Create(filepath.Join(tempDir, "file.txt"))
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file.Close()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := validateDir(tt.path, tt.opts...)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateSyncDirs(t *testing.T) {
	// Create temporary test directories
	tempDir, err := os.MkdirTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "source")
	dstDir := filepath.Join(tempDir, "destination")

	tests := []struct {
		name        string
		src         string
		dst         string
		expectError bool
		setup       func()
	}{
		{
			name:        "valid source and destination",
			src:         srcDir,
			dst:         dstDir,
			expectError: false,
			setup: func() {
				os.MkdirAll(srcDir, 0755)
			},
		},
		{
			name:        "non-existent source",
			src:         filepath.Join(tempDir, "nonexistent_src"),
			dst:         dstDir,
			expectError: true,
			setup: func() {
				os.MkdirAll(dstDir, 0755)
			},
		},
		{
			name:        "source is a file",
			src:         filepath.Join(tempDir, "source_file.txt"),
			dst:         dstDir,
			expectError: true,
			setup: func() {
				// Create a file instead of directory
				file, err := os.Create(filepath.Join(tempDir, "source_file.txt"))
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				file.Close()
				os.MkdirAll(dstDir, 0755)
			},
		},
		{
			name:        "destination will be created",
			src:         srcDir,
			dst:         filepath.Join(tempDir, "new_destination"),
			expectError: false,
			setup: func() {
				os.MkdirAll(srcDir, 0755)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up previous test
			os.RemoveAll(tempDir)
			os.MkdirAll(tempDir, 0755)

			tt.setup()

			err := ValidateSyncDirs(tt.src, tt.dst)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				// Check that it's a validation error
				if !isValidationError(err) {
					t.Errorf("Expected validation error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestConfigOptions(t *testing.T) {
	// Test withAllowCreate option
	cfg := newConfig(withAllowCreate())
	if !cfg.allowCreate {
		t.Error("Expected allowCreate to be true")
	}
	if cfg.perm != 0755 {
		t.Errorf("Expected default perm 0755, got %o", cfg.perm)
	}

	// Test default config
	cfg = newConfig()
	if cfg.allowCreate {
		t.Error("Expected allowCreate to be false by default")
	}
	if cfg.perm != 0755 {
		t.Errorf("Expected default perm 0755, got %o", cfg.perm)
	}
}

func TestConfigWithMultipleOptions(t *testing.T) {
	// Test with multiple options (though we only have one currently)
	cfg := newConfig(withAllowCreate())
	if !cfg.allowCreate {
		t.Error("Expected allowCreate to be true")
	}
}

// Helper function to check if an error is a validation error
func isValidationError(err error) bool {
	// Check if the error message contains validation-related text
	// This is a simple check - in a real scenario you might want to use error types
	return err != nil && (contains(err.Error(), "validation") ||
		contains(err.Error(), "source directory") ||
		contains(err.Error(), "target directory"))
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
