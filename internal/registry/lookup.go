// Package registry maintains the asdf plugin registry.
package registry

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// Lookup returns the plugin git URL registered for name.
func Lookup(name string) (string, bool, error) {
	path := filepath.Join(config.RegistryDir(), upstreamPluginsSubdir, name)

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("stat %s: %w", path, err)
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return "", false, fmt.Errorf("parse %s: %w", path, err)
	}
	url := cfg.Section("").Key("repository").String()
	if url == "" {
		return "", false, fmt.Errorf("%s: missing 'repository' key", path)
	}
	return url, true, nil
}
