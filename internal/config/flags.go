package config

import (
	"flag"
	"fmt"
	"os"
)

// FlagConfig implements ConfigProvider using CLI flags
type FlagConfig struct {
	cfg *Config
}

// Config returns the parsed config
func (f *FlagConfig) Config() *Config {
	return f.cfg
}

// ParseFlags parses CLI flags and returns a FlagConfig
func ParseFlags() (*FlagConfig, error) {
	usage := func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [--delete-missing] [--log-level LEVEL] <source> <target>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Usage = usage

	deleteMissing := flag.Bool("delete-missing", false, "Delete files from target that do not exist in source")
	logLevel := flag.String("log-level", "info", "Set logging level (error, warn, info, debug)")
	updateMethod := flag.String("update-method", "modtime", "Method for detecting file updates (modtime, sha256)")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid arguments: source and target paths are required")
	}

	cfg := &Config{
		Source:        args[0],
		Target:        args[1],
		DeleteMissing: *deleteMissing,
		LogLevel:      *logLevel,
		UpdateMethod:  *updateMethod,
	}

	return &FlagConfig{cfg: cfg}, nil
}
