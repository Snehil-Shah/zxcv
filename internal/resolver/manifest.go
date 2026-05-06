// Package resolver resolves the applicable manifest files.
package resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

// GlobalManifest returns the path to the global `.tool-versions` under $HOME.
func GlobalManifest() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home dir: %w", err)
	}
	return filepath.Join(home, manifestName), nil
}

// ApplicableManifest returns the path to the nearest `.tool-versions` walking up from dir.
// If no manifest is found, returns `dir/.tool-versions` (the path to create).
func ApplicableManifest(dir string) string {
	d := dir
	for {
		candidate := filepath.Join(d, manifestName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}
	return filepath.Join(dir, manifestName)
}
