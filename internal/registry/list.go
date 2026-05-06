// Package registry maintains the asdf plugin registry.
package registry

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// ListAll returns every plugin name in the upstream registry, unsorted.
// Skips directories so only proper plugin entries (files) are included.
func ListAll() ([]string, error) {
	dir := filepath.Join(config.RegistryDir(), upstreamPluginsSubdir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		names = append(names, e.Name())
	}
	return names, nil
}
