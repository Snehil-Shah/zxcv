// Package testutil provides shared test helpers.
package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/plugin"
)

// MakeFakePlugin creates a fake plugin under pluginsDir/<name>/ with the given scripts under bin/.
func MakeFakePlugin(t *testing.T, name string, scripts map[string]string) *plugin.Plugin {
	t.Helper()
	MakeFakePluginAt(t, filepath.Join(config.PluginsDir(), name), scripts)
	return plugin.New(name)
}

// MakeFakePluginAt creates a fake plugin-like at the given path with the given scripts under bin/.
func MakeFakePluginAt(t *testing.T, path string, scripts map[string]string) {
	t.Helper()
	bin := filepath.Join(path, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range scripts {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(body), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}
