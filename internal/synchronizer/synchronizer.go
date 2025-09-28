package synchronizer

import (
	"fmt"
	"snc/internal/config"
	"snc/internal/logger"
	"snc/internal/stream"
	"snc/internal/validate/dir"
)

type Synchronizer struct {
	cfg *config.Config
}

func NewSynchronizer(provider config.ConfigProvider) *Synchronizer {
	return &Synchronizer{cfg: provider.Config()}
}

func (s *Synchronizer) Sync() error {
	var hasErrors bool

	logger.Info("SYNC", "Starting synchronization process")
	logger.Debug("SYNC", "Configuration: Source=%s, Target=%s, DeleteMissing=%v",
		s.cfg.Source, s.cfg.Target, s.cfg.DeleteMissing)

	// Phase 1: Directory validation
	logger.Info("SYNC", "Phase 1: Validating directories")
	if err := dir.ValidateSyncDirs(s.cfg.Source, s.cfg.Target); err != nil {
		logger.Error("SYNC", "Directory validation failed: %v", err)
		hasErrors = true
	} else {
		logger.Success("SYNC", "Directory validation completed")
	}

	// Phase 2: File synchronization
	logger.Info("SYNC", "Phase 2: Synchronizing files")
	if err := stream.Sync(s.cfg); err != nil {
		logger.Error("SYNC", "File synchronization failed: %v", err)
		hasErrors = true
	} else {
		logger.Success("SYNC", "File synchronization completed")
	}

	// Phase 3: Delete missing files (if enabled)
	if s.cfg.DeleteMissing {
		logger.Info("SYNC", "Phase 3: Removing missing files")
		if err := stream.DeleteMissing(s.cfg.Source, s.cfg.Target); err != nil {
			logger.Error("SYNC", "Delete missing operation failed: %v", err)
			hasErrors = true
		} else {
			logger.Success("SYNC", "Delete missing operation completed")
		}
	} else {
		logger.Debug("SYNC", "Phase 3: Skipped (delete missing disabled)")
	}

	if hasErrors {
		logger.Warn("SYNC", "Synchronization completed with errors - check logs for details")
		return fmt.Errorf("sync completed with errors - check logs for details")
	}

	logger.Success("SYNC", "Synchronization completed successfully")
	return nil
}
