// Package registry maintains the asdf plugin registry.
package registry

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// Search returns plugin names matching query (case-insensitive substring), sorted alphabetically.
func Search(query string) ([]string, error) {
	dir := filepath.Join(config.RegistryDir(), upstreamPluginsSubdir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}
	q := strings.ToLower(query)
	var matches []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.Contains(strings.ToLower(name), q) {
			continue
		}
		matches = append(matches, name)
	}
	sort.Strings(matches)
	return matches, nil
}
