package errors

import (
	"fmt"
)

// Error types for different categories
var (
	// Directory-related errors
	ErrNotADirectory          = NewError("path is not a directory")
	ErrDirectoryNotAccessible = NewError("path is not accessible")
	ErrCannotCreateDirectory  = NewError("cannot create directory")
	ErrSourceDirValidation    = NewError("source directory validation failed")
	ErrTargetDirValidation    = NewError("target directory validation failed")

	// File-related errors
	ErrFileNotAccessible = NewError("file is not accessible")
	ErrCannotOpenFile    = NewError("cannot open file")
	ErrCannotCreateFile  = NewError("cannot create file")
	ErrCannotReadFile    = NewError("cannot read file")
	ErrCannotWriteFile   = NewError("cannot write file")
	ErrCannotCloseFile   = NewError("cannot close file")
	ErrFileCopyFailed    = NewError("file copy failed")
	ErrFileNotFound      = NewError("file not found")
	ErrCannotDeleteFile  = NewError("cannot delete file")

	// Sync-related errors
	ErrSyncFailed                = NewError("sync operation failed")
	ErrCannotComputeRelativePath = NewError("cannot compute relative path")
	ErrCannotCreateParentDir     = NewError("cannot create parent directory")
	ErrCannotStatFile            = NewError("cannot get file information")
)

// Error represents a custom error with context
type Error struct {
	message string
	context map[string]interface{}
}

// NewError creates a new error with a message
func NewError(message string) *Error {
	return &Error{
		message: message,
		context: make(map[string]interface{}),
	}
}

// Error implements the error interface
func (e *Error) Error() string {
	if len(e.context) == 0 {
		return e.message
	}

	contextStr := ""
	for key, value := range e.context {
		if contextStr != "" {
			contextStr += ", "
		}
		contextStr += fmt.Sprintf("%s: %v", key, value)
	}

	return fmt.Sprintf("%s (%s)", e.message, contextStr)
}

// WithContext adds context to the error
func (e *Error) WithContext(key string, value interface{}) *Error {
	e.context[key] = value
	return e
}

// WithPath adds a path context to the error
func (e *Error) WithPath(path string) *Error {
	return e.WithContext("path", path)
}

// WithSourcePath adds a source path context to the error
func (e *Error) WithSourcePath(path string) *Error {
	return e.WithContext("source", path)
}

// WithTargetPath adds a target path context to the error
func (e *Error) WithTargetPath(path string) *Error {
	return e.WithContext("target", path)
}

// Wrap wraps an existing error with additional context
func (e *Error) Wrap(err error) *Error {
	return e.WithContext("wrapped_error", err.Error())
}

// Helper functions for common error patterns

// NewDirectoryError creates a directory-related error with path context
func NewDirectoryError(baseErr *Error, path string, cause error) error {
	return fmt.Errorf("%s: %w: %v", path, baseErr, cause)
}

// NewFileError creates a file-related error with path context
func NewFileError(baseErr *Error, path string, cause error) error {
	return fmt.Errorf("%s: %w: %v", path, baseErr, cause)
}

// NewSyncError creates a sync-related error with context
func NewSyncError(baseErr *Error, context string, cause error) error {
	return fmt.Errorf("%s: %w: %v", context, baseErr, cause)
}

// NewValidationError creates a validation error with context
func NewValidationError(baseErr *Error, context string, cause error) error {
	return fmt.Errorf("%s: %w: %v", context, baseErr, cause)
}

// NewFileAccessError creates a formatted error message for file access issues
func NewFileAccessError(path string, cause error) error {
	return fmt.Errorf("error accessing %s: %w", path, NewFileError(ErrFileNotAccessible, path, cause))
}

// NewRelativePathError creates a formatted error message for relative path computation issues
func NewRelativePathError(path string, cause error) error {
	return fmt.Errorf("cannot compute relative path for %s: %w", path, NewSyncError(ErrCannotComputeRelativePath, path, cause))
}

// NewFileDeleteError creates a formatted error message for file deletion issues
func NewFileDeleteError(path string, cause error) error {
	return fmt.Errorf("failed to delete %s: %w", path, NewFileError(ErrCannotDeleteFile, path, cause))
}

// NewFileStatError creates a formatted error message for file stat issues
func NewFileStatError(path string, cause error) error {
	return fmt.Errorf("cannot stat %s: %w", path, NewFileError(ErrCannotStatFile, path, cause))
}

// NewDirectoryCreateError creates a formatted error message for directory creation issues
func NewDirectoryCreateError(path string, cause error) error {
	return fmt.Errorf("cannot create parent directory for %s: %w", path, NewSyncError(ErrCannotCreateParentDir, path, cause))
}

// NewFileOpenError creates a formatted error message for file opening issues
func NewFileOpenError(path string, cause error) error {
	return fmt.Errorf("cannot open source %s: %w", path, NewFileError(ErrCannotOpenFile, path, cause))
}

// NewFileCreateError creates a formatted error message for file creation issues
func NewFileCreateError(path string, cause error) error {
	return fmt.Errorf("cannot create target %s: %w", path, NewFileError(ErrCannotCreateFile, path, cause))
}

// NewFileCopyError creates a formatted error message for file copy issues
func NewFileCopyError(src, dst string, cause error) error {
	return fmt.Errorf("copy failed from %s to %s: %w", src, dst, NewSyncError(ErrFileCopyFailed.WithSourcePath(src).WithTargetPath(dst), "copy operation", cause))
}

// NewFileCloseError creates a formatted error message for file closing issues
func NewFileCloseError(path string, cause error) error {
	return fmt.Errorf("closing target %s failed: %w", path, NewFileError(ErrCannotCloseFile, path, cause))
}
