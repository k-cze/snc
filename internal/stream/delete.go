package stream

import (
	"os"
	"path/filepath"
	"snc/internal/logger"
)

// DeleteMissing removes files from dst that do not exist in src
func DeleteMissing(srcRoot, dstRoot string) error {
	logger.Info("DELETE", "Starting cleanup of missing files from %s", dstRoot)

	var fileCount, deletedCount, errorCount int

	err := filepath.WalkDir(dstRoot, func(dstPath string, d os.DirEntry, err error) error {
		if err != nil {
			logger.Error("DELETE", "Error accessing %s: %v", dstPath, err)
			errorCount++
			return nil
		}

		if d.IsDir() {
			logger.Debug("DELETE", "Skipping directory: %s", dstPath)
			return nil
		}

		fileCount++
		logger.Debug("DELETE", "Checking file: %s", dstPath)

		// compute relative path to dst root
		rel, relErr := filepath.Rel(dstRoot, dstPath)
		if relErr != nil {
			logger.Error("DELETE", "Cannot compute relative path for %s: %v", dstPath, relErr)
			errorCount++
			return nil
		}

		srcPath := filepath.Join(srcRoot, rel)

		// check if file exists in source
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			// File doesn't exist in source, delete it
			if err := os.Remove(dstPath); err != nil {
				logger.Error("DELETE", "Failed to delete missing file %s: %v", dstPath, err)
				errorCount++
			} else {
				logger.Progress("DELETE", "REMOVE", "Deleted missing file: %s", rel)
				deletedCount++
			}
		} else if err != nil {
			// Log error accessing source file but continue
			logger.Error("DELETE", "Error accessing source file %s: %v", srcPath, err)
			errorCount++
		} else {
			logger.Debug("DELETE", "File exists in source, keeping: %s", rel)
		}

		return nil
	})

	if err != nil {
		logger.Error("DELETE", "Directory walk failed: %v", err)
		return err
	}

	logger.Info("DELETE", "Cleanup completed: %d files checked, %d deleted, %d errors",
		fileCount, deletedCount, errorCount)

	return nil
}
