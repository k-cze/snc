package config

import (
	"flag"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := &Config{
		Source:        "/source",
		Target:        "/target",
		DeleteMissing: true,
		LogLevel:      "debug",
		UpdateMethod:  "sha256",
	}

	// Test Config struct fields
	if cfg.Source != "/source" {
		t.Errorf("Expected Source '/source', got '%s'", cfg.Source)
	}
	if cfg.Target != "/target" {
		t.Errorf("Expected Target '/target', got '%s'", cfg.Target)
	}
	if !cfg.DeleteMissing {
		t.Error("Expected DeleteMissing to be true")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("Expected LogLevel 'debug', got '%s'", cfg.LogLevel)
	}
	if cfg.UpdateMethod != "sha256" {
		t.Errorf("Expected UpdateMethod 'sha256', got '%s'", cfg.UpdateMethod)
	}
}

func TestFlagConfig(t *testing.T) {
	cfg := &Config{
		Source:        "/test/source",
		Target:        "/test/target",
		DeleteMissing: false,
		LogLevel:      "info",
		UpdateMethod:  "modtime",
	}

	flagConfig := &FlagConfig{cfg: cfg}

	// Test ConfigProvider interface
	retrievedConfig := flagConfig.Config()
	if retrievedConfig != cfg {
		t.Error("Expected retrieved config to be the same as original")
	}

	// Test field access
	if retrievedConfig.Source != "/test/source" {
		t.Errorf("Expected Source '/test/source', got '%s'", retrievedConfig.Source)
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedConfig *Config
		expectError    bool
	}{
		{
			name: "valid arguments with defaults",
			args: []string{"--delete-missing", "--log-level", "debug", "--update-method", "sha256", "/source", "/target"},
			expectedConfig: &Config{
				Source:        "/source",
				Target:        "/target",
				DeleteMissing: true,
				LogLevel:      "debug",
				UpdateMethod:  "sha256",
			},
			expectError: false,
		},
		{
			name: "valid arguments with minimal flags",
			args: []string{"/source", "/target"},
			expectedConfig: &Config{
				Source:        "/source",
				Target:        "/target",
				DeleteMissing: false,
				LogLevel:      "info",
				UpdateMethod:  "modtime",
			},
			expectError: false,
		},
		{
			name:        "missing source argument",
			args:        []string{"/target"},
			expectError: true,
		},
		{
			name:        "missing target argument",
			args:        []string{"/source"},
			expectError: true,
		},
		{
			name:        "no arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "too many arguments",
			args:        []string{"/source", "/target", "/extra"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine to avoid conflicts between tests
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set up test arguments
			oldArgs := os.Args
			os.Args = append([]string{os.Args[0]}, tt.args...)
			defer func() {
				os.Args = oldArgs
			}()

			flagConfig, err := ParseFlags()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if flagConfig == nil {
				t.Fatal("Expected flagConfig to be non-nil")
			}

			config := flagConfig.Config()
			if config == nil {
				t.Fatal("Expected config to be non-nil")
			}

			// Compare expected vs actual config
			if config.Source != tt.expectedConfig.Source {
				t.Errorf("Expected Source '%s', got '%s'", tt.expectedConfig.Source, config.Source)
			}
			if config.Target != tt.expectedConfig.Target {
				t.Errorf("Expected Target '%s', got '%s'", tt.expectedConfig.Target, config.Target)
			}
			if config.DeleteMissing != tt.expectedConfig.DeleteMissing {
				t.Errorf("Expected DeleteMissing %v, got %v", tt.expectedConfig.DeleteMissing, config.DeleteMissing)
			}
			if config.LogLevel != tt.expectedConfig.LogLevel {
				t.Errorf("Expected LogLevel '%s', got '%s'", tt.expectedConfig.LogLevel, config.LogLevel)
			}
			if config.UpdateMethod != tt.expectedConfig.UpdateMethod {
				t.Errorf("Expected UpdateMethod '%s', got '%s'", tt.expectedConfig.UpdateMethod, config.UpdateMethod)
			}
		})
	}
}

func TestParseFlagsWithInvalidUpdateMethod(t *testing.T) {
	// Reset flag.CommandLine to avoid conflicts
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Set up test arguments with invalid update method
	oldArgs := os.Args
	os.Args = []string{os.Args[0], "--update-method", "invalid", "/source", "/target"}
	defer func() {
		os.Args = oldArgs
	}()

	flagConfig, err := ParseFlags()
	if err != nil {
		t.Fatalf("Unexpected error during flag parsing: %v", err)
	}

	// The flag parsing should succeed, but the update method validation
	// happens later in the stream package
	config := flagConfig.Config()
	if config.UpdateMethod != "invalid" {
		t.Errorf("Expected UpdateMethod 'invalid', got '%s'", config.UpdateMethod)
	}
}
