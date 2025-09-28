package errors

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	// Test basic error creation
	err := NewError("test error")
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}

	// Test error with context
	err = err.WithContext("key1", "value1")
	expected := "test error (key1: value1)"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test error with multiple context
	err = err.WithContext("key2", "value2")
	// Order is not guaranteed in map iteration, so we check both possible orders
	errorStr := err.Error()
	if errorStr != "test error (key1: value1, key2: value2)" &&
		errorStr != "test error (key2: value2, key1: value1)" {
		t.Errorf("Expected error with both contexts, got '%s'", errorStr)
	}
}

func TestErrorWithPath(t *testing.T) {
	err := NewError("file error")
	err = err.WithPath("/test/path")

	errorStr := err.Error()
	if errorStr != "file error (path: /test/path)" {
		t.Errorf("Expected error with path context, got '%s'", errorStr)
	}
}

func TestErrorWithSourceAndTargetPaths(t *testing.T) {
	err := NewError("sync error")
	err = err.WithSourcePath("/source")
	err = err.WithTargetPath("/target")

	errorStr := err.Error()
	// Check that both source and target are present
	if errorStr != "sync error (source: /source, target: /target)" &&
		errorStr != "sync error (target: /target, source: /source)" {
		t.Errorf("Expected error with source and target paths, got '%s'", errorStr)
	}
}

func TestErrorWrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := NewError("wrapper error")
	err = err.Wrap(originalErr)

	errorStr := err.Error()
	if errorStr != "wrapper error (wrapped_error: original error)" {
		t.Errorf("Expected wrapped error, got '%s'", errorStr)
	}
}

func TestHelperFunctions(t *testing.T) {
	originalErr := errors.New("original error")

	// Test NewDirectoryError
	dirErr := NewDirectoryError(ErrNotADirectory, "/test/dir", originalErr)
	if dirErr == nil {
		t.Error("Expected non-nil directory error")
	}

	// Test NewFileError
	fileErr := NewFileError(ErrCannotOpenFile, "/test/file", originalErr)
	if fileErr == nil {
		t.Error("Expected non-nil file error")
	}

	// Test NewSyncError
	syncErr := NewSyncError(ErrSyncFailed, "test operation", originalErr)
	if syncErr == nil {
		t.Error("Expected non-nil sync error")
	}

	// Test NewValidationError
	validationErr := NewValidationError(ErrSourceDirValidation, "test validation", originalErr)
	if validationErr == nil {
		t.Error("Expected non-nil validation error")
	}

	// Test NewFileAccessError
	accessErr := NewFileAccessError("/test/file", originalErr)
	if accessErr == nil {
		t.Error("Expected non-nil file access error")
	}

	// Test NewRelativePathError
	relativeErr := NewRelativePathError("/test/path", originalErr)
	if relativeErr == nil {
		t.Error("Expected non-nil relative path error")
	}

	// Test NewFileDeleteError
	deleteErr := NewFileDeleteError("/test/file", originalErr)
	if deleteErr == nil {
		t.Error("Expected non-nil file delete error")
	}

	// Test NewFileStatError
	statErr := NewFileStatError("/test/file", originalErr)
	if statErr == nil {
		t.Error("Expected non-nil file stat error")
	}

	// Test NewDirectoryCreateError
	createErr := NewDirectoryCreateError("/test/dir", originalErr)
	if createErr == nil {
		t.Error("Expected non-nil directory create error")
	}

	// Test NewFileOpenError
	openErr := NewFileOpenError("/test/file", originalErr)
	if openErr == nil {
		t.Error("Expected non-nil file open error")
	}

	// Test NewFileCreateError
	createFileErr := NewFileCreateError("/test/file", originalErr)
	if createFileErr == nil {
		t.Error("Expected non-nil file create error")
	}

	// Test NewFileCopyError
	copyErr := NewFileCopyError("/source", "/target", originalErr)
	if copyErr == nil {
		t.Error("Expected non-nil file copy error")
	}

	// Test NewFileCloseError
	closeErr := NewFileCloseError("/test/file", originalErr)
	if closeErr == nil {
		t.Error("Expected non-nil file close error")
	}
}

func TestErrorTypes(t *testing.T) {
	// Test that all error types are defined
	errorTypes := []*Error{
		ErrNotADirectory,
		ErrDirectoryNotAccessible,
		ErrCannotCreateDirectory,
		ErrSourceDirValidation,
		ErrTargetDirValidation,
		ErrFileNotAccessible,
		ErrCannotOpenFile,
		ErrCannotCreateFile,
		ErrCannotReadFile,
		ErrCannotWriteFile,
		ErrCannotCloseFile,
		ErrFileCopyFailed,
		ErrFileNotFound,
		ErrCannotDeleteFile,
		ErrSyncFailed,
		ErrCannotComputeRelativePath,
		ErrCannotCreateParentDir,
		ErrCannotStatFile,
	}

	for _, err := range errorTypes {
		if err == nil {
			t.Error("Expected error type to be non-nil")
		}
		if err.Error() == "" {
			t.Error("Expected error type to have non-empty message")
		}
	}
}
