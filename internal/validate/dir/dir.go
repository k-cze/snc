package dir

import (
	"os"
	"snc/internal/errors"
)

type configOption func(*config)

type config struct {
	allowCreate bool
	perm        os.FileMode
}

func newConfig(opts ...configOption) *config {
	// defaults
	cfg := &config{
		allowCreate: false,
		perm:        0755,
	}

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func withAllowCreate() configOption {
	return func(cfg *config) {
		cfg.allowCreate = true
	}
}

func validateDir(path string, opts ...configOption) error {
	cfg := newConfig(opts...)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) && cfg.allowCreate {
			if mkErr := os.MkdirAll(path, cfg.perm); mkErr != nil {
				return errors.NewDirectoryError(errors.ErrCannotCreateDirectory, path, mkErr)
			}
			return nil
		}
		return errors.NewDirectoryError(errors.ErrDirectoryNotAccessible, path, err)
	}

	if !info.IsDir() {
		return errors.NewDirectoryError(errors.ErrNotADirectory, path, nil)
	}
	return nil
}

// ValidateSyncDirs validates the source and target directories.
// Source must exist; target is created if missing.
func ValidateSyncDirs(src, dst string) error {
	if err := validateDir(src); err != nil {
		return errors.NewValidationError(errors.ErrSourceDirValidation, "source directory", err)
	}
	if err := validateDir(dst, withAllowCreate()); err != nil {
		return errors.NewValidationError(errors.ErrTargetDirValidation, "target directory", err)
	}
	return nil
}
