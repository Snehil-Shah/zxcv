// Package installer_test tests the `installer` package.
package installer_test

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestUninstallAll_WipesEverything(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", nil)
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")
	testutil.MakeFakeShims(t, "alpha", "alpha")
	testutil.MakeFakePlugin(t, "beta", nil)
	testutil.MakeFakeInstall(t, "beta", "2.0.0", "beta")
	testutil.MakeFakeShims(t, "beta", "beta")

	if err := installer.UninstallAll(context.Background()); err != nil {
		t.Fatalf("UninstallAll: %v", err)
	}
	for _, dir := range []string{config.InstallsDir(), config.PluginsDir(), config.ShimsDir()} {
		if _, err := os.Stat(dir); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s should be removed: %v", dir, err)
		}
	}
}

func TestUninstallTool_RemovesEverythingForName(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", nil)
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")
	testutil.MakeFakeInstall(t, "alpha", "2.0.0", "alpha")
	testutil.MakeFakeShims(t, "alpha", "alpha")
	testutil.MakeFakePlugin(t, "beta", nil)
	testutil.MakeFakeInstall(t, "beta", "9.0.0", "beta")
	testutil.MakeFakeShims(t, "beta", "beta")

	if err := installer.UninstallTool(context.Background(), "alpha"); err != nil {
		t.Fatalf("UninstallTool: %v", err)
	}

	for _, p := range []string{
		filepath.Join(config.InstallsDir(), "alpha"),
		filepath.Join(config.PluginsDir(), "alpha"),
		filepath.Join(config.ShimsDir(), "alpha"),
	} {
		if _, err := os.Stat(p); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s should be removed: %v", p, err)
		}
	}
	// beta untouched.
	for _, p := range []string{
		filepath.Join(config.InstallsDir(), "beta"),
		filepath.Join(config.PluginsDir(), "beta"),
		filepath.Join(config.ShimsDir(), "beta"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s should still exist: %v", p, err)
		}
	}
}

func TestUninstallTool_WipesPartialsAndDownloads(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", nil)
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")
	testutil.MakeFakeShims(t, "alpha", "alpha")
	// Stale partial dir from an interrupted previous install.
	if err := os.MkdirAll(filepath.Join(config.InstallsDir(), "alpha", "2.0.0.partial"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Stale download dir from an interrupted previous install.
	if err := os.MkdirAll(filepath.Join(config.DownloadsDir(), "alpha", "2.0.0"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := installer.UninstallTool(context.Background(), "alpha"); err != nil {
		t.Fatal(err)
	}

	for _, p := range []string{
		filepath.Join(config.InstallsDir(), "alpha"),
		filepath.Join(config.PluginsDir(), "alpha"),
		filepath.Join(config.DownloadsDir(), "alpha"),
		filepath.Join(config.ShimsDir(), "alpha"),
	} {
		if _, err := os.Stat(p); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s should be removed: %v", p, err)
		}
	}
}

func TestUninstallTool_NotInstalled(t *testing.T) {
	testutil.HermeticDataDir(t)
	err := installer.UninstallTool(context.Background(), "ghost")
	if err == nil || !strings.Contains(err.Error(), "not installed") {
		t.Errorf("expected 'not installed' error, got %v", err)
	}
}

func TestUninstallVersion_LastVersion_RemovesToolAndPruneShim(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", nil)
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")
	testutil.MakeFakeShims(t, "alpha", "alpha")

	if err := installer.UninstallVersion(context.Background(), "alpha", "1.0.0"); err != nil {
		t.Fatalf("UninstallVersion: %v", err)
	}

	for _, p := range []string{
		filepath.Join(config.InstallsDir(), "alpha"),
		filepath.Join(config.PluginsDir(), "alpha"),
		filepath.Join(config.ShimsDir(), "alpha"),
	} {
		if _, err := os.Stat(p); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("%s should be removed: %v", p, err)
		}
	}
}

func TestUninstallVersion_NonLast_KeepsToolAndShim(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", nil)
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")
	testutil.MakeFakeInstall(t, "alpha", "2.0.0", "alpha")
	testutil.MakeFakeShims(t, "alpha", "alpha")

	if err := installer.UninstallVersion(context.Background(), "alpha", "1.0.0"); err != nil {
		t.Fatalf("UninstallVersion: %v", err)
	}

	if _, err := os.Stat(filepath.Join(config.InstallsDir(), "alpha", "1.0.0")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("v1.0.0 should be removed: %v", err)
	}
	for _, p := range []string{
		filepath.Join(config.InstallsDir(), "alpha", "2.0.0"),
		filepath.Join(config.PluginsDir(), "alpha"),
		filepath.Join(config.ShimsDir(), "alpha"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s should still exist: %v", p, err)
		}
	}
}

func TestUninstallVersion_NotInstalled(t *testing.T) {
	testutil.HermeticDataDir(t)
	err := installer.UninstallVersion(context.Background(), "ghost", "1.0.0")
	if err == nil || !strings.Contains(err.Error(), "not installed") {
		t.Errorf("expected 'not installed' error, got %v", err)
	}
}

func TestUninstallVersion_FiresUninstallCallback(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", map[string]string{
		"uninstall": "#!/bin/sh\necho \"$ASDF_INSTALL_VERSION\" > \"$HOME/.uninstall-fired\"\n",
	})
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")

	if err := installer.UninstallVersion(context.Background(), "alpha", "1.0.0"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".uninstall-fired"))
	if err != nil {
		t.Fatalf("uninstall callback did not fire: %v", err)
	}
	if got := strings.TrimSpace(string(data)); got != "1.0.0" {
		t.Errorf("callback got version %q, want %q", got, "1.0.0")
	}
}

func TestUninstallVersion_FiresPrePluginRemove_OnLastVersion(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", map[string]string{
		"pre-plugin-remove": "#!/bin/sh\ntouch \"$HOME/.pre-plugin-remove-fired\"\n",
	})
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")

	if err := installer.UninstallVersion(context.Background(), "alpha", "1.0.0"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".pre-plugin-remove-fired")); err != nil {
		t.Errorf("pre-plugin-remove did not fire on last version: %v", err)
	}
}

func TestUninstallVersion_NonLast_SkipsPrePluginRemove(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "alpha", map[string]string{
		"pre-plugin-remove": "#!/bin/sh\ntouch \"$HOME/.pre-plugin-remove-fired\"\n",
	})
	testutil.MakeFakeInstall(t, "alpha", "1.0.0", "alpha")
	testutil.MakeFakeInstall(t, "alpha", "2.0.0", "alpha")

	if err := installer.UninstallVersion(context.Background(), "alpha", "1.0.0"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".pre-plugin-remove-fired")); err == nil {
		t.Errorf("pre-plugin-remove should NOT fire when other versions remain")
	}
}

