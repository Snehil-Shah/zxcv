// Package plugin_test tests the `plugin` package.
package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/plugin"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestLink_NotADir(t *testing.T) {
	testutil.HermeticDataDir(t)
	file := filepath.Join(t.TempDir(), "f")
	testutil.WriteFile(t, file, "hi")
	err := plugin.New("foo").Link(context.Background(), file)
	if err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("expected 'not a directory' error, got %v", err)
	}
}

func TestLink_MissingBinInstall(t *testing.T) {
	testutil.HermeticDataDir(t)
	src := t.TempDir()
	err := plugin.New("foo").Link(context.Background(), src)
	if err == nil || !strings.Contains(err.Error(), "missing bin/install") {
		t.Errorf("expected 'missing bin/install' error, got %v", err)
	}
}

func TestLink_Success(t *testing.T) {
	testutil.HermeticDataDir(t)
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{"install": "#!/bin/sh\n"})
	p := plugin.New("foo")
	if err := p.Link(context.Background(), src); err != nil {
		t.Fatal(err)
	}
	target, err := os.Readlink(p.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if target != src {
		t.Errorf("symlink target = %q, want %q", target, src)
	}
}

func TestLink_FiresPostPluginAdd(t *testing.T) {
	testutil.HermeticDataDir(t)
	marker := filepath.Join(os.Getenv("HOME"), ".post-plugin-add-fired")
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{
		"install":         "#!/bin/sh\n",
		"post-plugin-add": "#!/bin/sh\ntouch \"" + marker + "\"\n",
	})
	if err := plugin.New("foo").Link(context.Background(), src); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("post-plugin-add did not fire: %v", err)
	}
}
