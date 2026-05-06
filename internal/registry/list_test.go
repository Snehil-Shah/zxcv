// Package registry_test tests the `registry` package.
package registry_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestListAll_NoRegistry(t *testing.T) {
	testutil.HermeticDataDir(t)

	got, err := registry.ListAll()
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestListAll_SkipsDirectories(t *testing.T) {
	testutil.HermeticDataDir(t)
	writePluginEntry(t, "real", "repository = x\n")
	if err := os.MkdirAll(filepath.Join(config.RegistryDir(), "plugins", "subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := registry.ListAll()
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	want := []string{"real"}
	if len(got) != 1 || got[0] != want[0] {
		t.Errorf("got %v, want %v", got, want)
	}
}
