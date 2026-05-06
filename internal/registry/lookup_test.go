// Package registry_test tests the `registry` package.
package registry_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestLookup_Found(t *testing.T) {
	testutil.HermeticDataDir(t)
	writePluginEntry(t, "node", "repository = https://github.com/asdf-vm/asdf-nodejs.git\n")

	url, found, err := registry.Lookup("node")
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("found = false, want true")
	}
	if url != "https://github.com/asdf-vm/asdf-nodejs.git" {
		t.Errorf("url = %q", url)
	}
}

func TestLookup_NotFound(t *testing.T) {
	testutil.HermeticDataDir(t)

	url, found, err := registry.Lookup("ghost")
	if err != nil {
		t.Fatalf("expected no error for missing plugin, got: %v", err)
	}
	if found {
		t.Errorf("found = true, want false")
	}
	if url != "" {
		t.Errorf("url = %q, want empty", url)
	}
}

func TestLookup_MissingRepositoryKey(t *testing.T) {
	testutil.HermeticDataDir(t)
	writePluginEntry(t, "broken", "other = thing\n")

	_, _, err := registry.Lookup("broken")
	if err == nil || !strings.Contains(err.Error(), "missing 'repository' key") {
		t.Errorf("expected missing-repository error, got: %v", err)
	}
}

// writePluginEntry simulates an entry in the cloned upstream registry.
func writePluginEntry(t *testing.T, name, body string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(config.RegistryDir(), "plugins", name), body)
}
