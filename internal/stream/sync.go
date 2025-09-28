package stream

import (
	"io"
	"os"
	"path/filepath"
	"snc/internal/config"
	"snc/internal/errors"
	"snc/internal/logger"
	"time"
)

// Sync performs file synchronization using the specified configuration
func Sync(cfg *config.Config) error {
	logger.Info("STREAM", "Starting file synchronization from %s to %s", cfg.Source, cfg.Target)
	logger.Info("STREAM", "Using update method: %s", cfg.UpdateMethod)

	// Create update strategy
	updateStrategy, err := NewUpdateStrategy(cfg.UpdateMethod)
	if err != nil {
		logger.Error("STREAM", "Failed to create update strategy: %v", err)
		return errors.NewSyncError(errors.ErrSyncFailed, "update strategy creation", err)
	}

	var fileCount, copiedCount, skippedCount, errorCount int

	err = filepath.WalkDir(cfg.Source, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			logger.Error("STREAM", "Error accessing %s: %v", path, err)
			errorCount++
			return nil // continue walking
		}

		if d.IsDir() {
			logger.Debug("STREAM", "Skipping directory: %s", path)
			return nil
		}

		fileCount++
		logger.Debug("STREAM", "Processing file: %s", path)

		// Process the file
		if err := processFileWithStrategy(cfg.Source, cfg.Target, path, d, updateStrategy); err != nil {
			logger.Error("STREAM", "Failed to process file %s: %v", path, err)
			errorCount++
		} else {
			copiedCount++
		}

		return nil
	})

	if err != nil {
		logger.Error("STREAM", "Directory walk failed: %v", err)
		return errors.NewSyncError(errors.ErrSyncFailed, "sync operation", err)
	}

	logger.Info("STREAM", "Synchronization completed: %d files processed, %d copied, %d skipped, %d errors",
		fileCount, copiedCount, skippedCount, errorCount)

	return nil
}

// processFileWithStrategy handles a single file during synchronization using the specified update strategy
func processFileWithStrategy(srcRoot, dstRoot, srcPath string, d os.DirEntry, strategy UpdateStrategy) error {
	// Calculate relative path
	rel, relErr := filepath.Rel(srcRoot, srcPath)
	if relErr != nil {
		logger.Error("STREAM", "Cannot compute relative path for %s: %v", srcPath, relErr)
		return errors.NewRelativePathError(srcPath, relErr)
	}

	dstPath := filepath.Join(dstRoot, rel)
	logger.Debug("STREAM", "Processing: %s -> %s", srcPath, dstPath)

	// Check if destination file exists
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		// File doesn't exist, copy it
		logger.Progress("STREAM", "COPY", "New file: %s", rel)
		return copyFile(srcPath, dstPath)
	} else if err != nil {
		// Error accessing destination file
		logger.Error("STREAM", "Cannot access destination file %s: %v", dstPath, err)
		return errors.NewFileStatError(dstPath, err)
	}

	// File exists, check if update is needed using the strategy
	needsUpdate, err := strategy.NeedsUpdate(srcPath, dstPath)
	if err != nil {
		logger.Error("STREAM", "Failed to check if file needs update %s: %v", srcPath, err)
		return err
	}

	if needsUpdate {
		logger.Progress("STREAM", "UPDATE", "Modified file: %s", rel)
		return copyFile(srcPath, dstPath)
	} else {
		logger.Debug("STREAM", "Skipping unchanged file: %s", rel)
		return nil
	}
}

func copyFile(src, dst string) error {
	logger.Debug("STREAM", "Starting copy: %s -> %s", src, dst)

	// ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		logger.Error("STREAM", "Cannot create parent directory for %s: %v", dst, err)
		return errors.NewSyncError(errors.ErrCannotCreateParentDir, dst, err)
	}

	// Open source file
	in, err := os.Open(src)
	if err != nil {
		logger.Error("STREAM", "Cannot open source file %s: %v", src, err)
		return errors.NewFileError(errors.ErrCannotOpenFile, src, err)
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil {
			logger.Warn("STREAM", "Failed to close source file %s: %v", src, closeErr)
		}
	}()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		logger.Error("STREAM", "Cannot create destination file %s: %v", dst, err)
		return errors.NewFileError(errors.ErrCannotCreateFile, dst, err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			logger.Warn("STREAM", "Failed to close destination file %s: %v", dst, closeErr)
		}
	}()

	// Copy file contents
	bytesCopied, err := io.Copy(out, in)
	if err != nil {
		logger.Error("STREAM", "File copy failed from %s to %s: %v", src, dst, err)
		return errors.NewSyncError(errors.ErrFileCopyFailed.WithSourcePath(src).WithTargetPath(dst), "copy operation", err)
	}

	// Preserve file modtime
	if srcInfo, statErr := in.Stat(); statErr == nil {
		if chtimesErr := os.Chtimes(dst, time.Now(), srcInfo.ModTime()); chtimesErr != nil {
			logger.Warn("STREAM", "Failed to preserve modtime for %s: %v", dst, chtimesErr)
		}
	} else {
		logger.Warn("STREAM", "Failed to stat source file %s for modtime: %v", src, statErr)
	}

	logger.Success("STREAM", "Copied %s -> %s (%d bytes)", src, dst, bytesCopied)
	return nil
}
