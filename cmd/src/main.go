package main

import (
	"os"
	"snc/internal/config"
	"snc/internal/logger"
	"snc/internal/synchronizer"
)

func main() {
	cfgProvider, err := config.ParseFlags()
	if err != nil {
		logger.Error("MAIN", "Failed to parse configuration: %v", err)
		os.Exit(2)
	}

	// Set log level from config if available
	if cfgProvider.Config().LogLevel != "" {
		logger.SetLevelFromString(cfgProvider.Config().LogLevel)
	}

	logger.Info("MAIN", "Starting file synchronization tool")
	logger.Info("MAIN", "Source: %s, Target: %s, Delete missing: %v",
		cfgProvider.Config().Source,
		cfgProvider.Config().Target,
		cfgProvider.Config().DeleteMissing)

	sn := synchronizer.NewSynchronizer(cfgProvider)
	if err := sn.Sync(); err != nil {
		logger.Error("MAIN", "Sync completed with errors: %v", err)
		os.Exit(1)
	}

	logger.Success("MAIN", "Synchronization completed successfully")
}
