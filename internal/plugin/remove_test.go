// Package plugin_test tests the `plugin` package.
package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestRemove_Success(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "foo", nil)
	if err := p.Remove(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(p.Dir); !os.IsNotExist(err) {
		t.Errorf("expected plugin dir to be gone, got err=%v", err)
	}
}

func TestRemove_FiresPrePluginRemove(t *testing.T) {
	testutil.HermeticDataDir(t)
	marker := filepath.Join(os.Getenv("HOME"), ".pre-plugin-remove-fired")
	p := testutil.MakeFakePlugin(t, "foo", map[string]string{
		"pre-plugin-remove": "#!/bin/sh\ntouch \"" + marker + "\"\n",
	})
	if err := p.Remove(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("pre-plugin-remove did not fire: %v", err)
	}
}
