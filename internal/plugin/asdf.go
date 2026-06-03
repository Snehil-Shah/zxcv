// Package plugin manages asdf plugins.
package plugin

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

//go:embed scripts/asdf.sh
var asdfShim []byte

// ensureAsdfShim writes the asdf shim to the libexec dir.
func ensureAsdfShim() (string, error) {
	dir := config.LibexecDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", dir, err)
	}
	path := filepath.Join(dir, "asdf")
	if err := os.WriteFile(path, asdfShim, 0o755); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return dir, nil
}

// withAsdfShimOnPath returns a copy of env with the asdf shim dir prepended to PATH.
func withAsdfShimOnPath(env map[string]string) (map[string]string, error) {
	libexec, err := ensureAsdfShim()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(env)+1)
	for k, v := range env {
		out[k] = v
	}
	existing := out["PATH"]
	if existing == "" {
		existing = os.Getenv("PATH")
	}
	out["PATH"] = libexec + string(os.PathListSeparator) + existing
	return out, nil
}
