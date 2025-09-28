package stream

import (
	"os"
	"path/filepath"
	"snc/internal/config"
	"testing"
)

func TestSync(t *testing.T) {
	// Create temporary test directories
	tempDir, err := os.MkdirTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "source")
	dstDir := filepath.Join(tempDir, "destination")

	// Create source directory and files
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "subdir/file3.txt"}
	for _, file := range testFiles {
		filePath := filepath.Join(srcDir, file)
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			t.Fatalf("Failed to create subdir: %v", err)
		}
		err = os.WriteFile(filePath, []byte("content for "+file), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	tests := []struct {
		name          string
		config        *config.Config
		expectError   bool
		expectedFiles []string
	}{
		{
			name: "sync with modtime strategy",
			config: &config.Config{
				Source:       srcDir,
				Target:       dstDir,
				UpdateMethod: "modtime",
			},
			expectError:   false,
			expectedFiles: []string{"file1.txt", "file2.txt", "subdir/file3.txt"},
		},
		{
			name: "sync with sha256 strategy",
			config: &config.Config{
				Source:       srcDir,
				Target:       dstDir,
				UpdateMethod: "sha256",
			},
			expectError:   false,
			expectedFiles: []string{"file1.txt", "file2.txt", "subdir/file3.txt"},
		},
		{
			name: "sync with invalid strategy",
			config: &config.Config{
				Source:       srcDir,
				Target:       dstDir,
				UpdateMethod: "invalid",
			},
			expectError: true,
		},
		{
			name: "sync with non-existent source",
			config: &config.Config{
				Source:       "/non/existent/source",
				Target:       dstDir,
				UpdateMethod: "modtime",
			},
			expectError: false, // The current implementation doesn't return error for non-existent source
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up destination directory
			os.RemoveAll(dstDir)

			err := Sync(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check that expected files were copied
			for _, file := range tt.expectedFiles {
				dstFile := filepath.Join(dstDir, file)
				if _, err := os.Stat(dstFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s to exist in destination", file)
				}
			}
		})
	}
}

func TestProcessFileWithStrategy(t *testing.T) {
	// Create temporary test directories
	tempDir, err := os.MkdirTemp("", "sync_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	srcDir := filepath.Join(tempDir, "source")
	dstDir := filepath.Join(tempDir, "destination")

	// Create source directory
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}

	// Create test file
	srcFile := filepath.Join(srcDir, "test.txt")
	err = os.WriteFile(srcFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Get file info
	fileInfo, err := os.Stat(srcFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	// Create DirEntry mock
	dirEntry := &mockDirEntry{fileInfo: fileInfo}

	tests := []struct {
		name        string
		strategy    UpdateStrategy
		setupDst    func() // function to set up destination
		expectError bool
	}{
		{
			name:     "new file with modtime strategy",
			strategy: &ModTimeStrategy{},
			setupDst: func() {
				// No destination file
			},
			expectError: false,
		},
		{
			name:     "new file with sha256 strategy",
			strategy: &SHA256Strategy{},
			setupDst: func() {
				// No destination file
			},
			expectError: false,
		},
		{
			name:     "existing identical file with modtime strategy",
			strategy: &ModTimeStrategy{},
			setupDst: func() {
				// Create identical destination file
				err := os.WriteFile(filepath.Join(dstDir, "test.txt"), []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create destination file: %v", err)
				}
			},
			expectError: false,
		},
		{
			name:     "existing identical file with sha256 strategy",
			strategy: &SHA256Strategy{},
			setupDst: func() {
				// Create identical destination file
				err := os.WriteFile(filepath.Join(dstDir, "test.txt"), []byte("test content"), 0644)
				if err != nil {
					t.Fatalf("Failed to create destination file: %v", err)
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up destination directory
			os.RemoveAll(dstDir)
			os.MkdirAll(dstDir, 0755)

			tt.setupDst()

			err := processFileWithStrategy(srcDir, dstDir, srcFile, dirEntry, tt.strategy)

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

// Mock DirEntry for testing
type mockDirEntry struct {
	fileInfo os.FileInfo
}

func (m *mockDirEntry) Name() string {
	return m.fileInfo.Name()
}

func (m *mockDirEntry) IsDir() bool {
	return m.fileInfo.IsDir()
}

func (m *mockDirEntry) Type() os.FileMode {
	return m.fileInfo.Mode()
}

func (m *mockDirEntry) Info() (os.FileInfo, error) {
	return m.fileInfo, nil
}
