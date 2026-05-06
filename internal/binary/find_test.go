// Package binary_test tests the `binary` package.
package binary_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/binary"
	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestFind_Success(t *testing.T) {
	testutil.HermeticDataDir(t)
	installPath := testutil.MakeFakeInstall(t, "mytool", "1.0.0", "foo")

	cwd := t.TempDir()
	testutil.WriteFile(t, filepath.Join(cwd, ".tool-versions"), "mytool 1.0.0\n")

	got, err := binary.Find(context.Background(), "mytool", "foo", cwd)
	if err != nil {
		t.Fatal(err)
	}
	wantBin := filepath.Join(installPath, "bin", "foo")
	if got.BinPath != wantBin {
		t.Errorf("BinPath = %q, want %q", got.BinPath, wantBin)
	}
	if got.Plugin == nil || got.Plugin.Name != "mytool" {
		t.Errorf("Plugin.Name = %q, want %q", got.Plugin.Name, "mytool")
	}
	if got.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", got.Version, "1.0.0")
	}
	if got.Plugin == nil || got.Plugin.Dir != filepath.Join(config.PluginsDir(), "mytool") {
		t.Errorf("Plugin.Dir = %q", got.Plugin.Dir)
	}
}

func TestFind_NoVersionSet(t *testing.T) {
	testutil.HermeticDataDir(t)
	_, err := binary.Find(context.Background(), "mytool", "foo", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "no version of mytool set") {
		t.Errorf("expected 'no version' error, got %v", err)
	}
}

func TestFind_NotInstalled(t *testing.T) {
	testutil.HermeticDataDir(t)
	cwd := t.TempDir()
	testutil.WriteFile(t, filepath.Join(cwd, ".tool-versions"), "mytool 1.0.0\n")

	_, err := binary.Find(context.Background(), "mytool", "foo", cwd)
	if err == nil || !strings.Contains(err.Error(), "not installed") {
		t.Errorf("expected 'not installed' error, got %v", err)
	}
}

func TestFind_SystemVersion(t *testing.T) {
	testutil.HermeticDataDir(t)
	sysDir := t.TempDir()
	sysFoo := filepath.Join(sysDir, "foo")
	testutil.WriteFile(t, sysFoo, "x")
	t.Setenv("PATH", sysDir)

	cwd := t.TempDir()
	testutil.WriteFile(t, filepath.Join(cwd, ".tool-versions"), "mytool system\n")

	got, err := binary.Find(context.Background(), "mytool", "foo", cwd)
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != sysFoo {
		t.Errorf("BinPath = %q, want %q", got.BinPath, sysFoo)
	}
	if got.Plugin != nil || got.Version != "" || got.InstallPath != "" {
		t.Errorf("expected empty plugin metadata for system version, got %+v", got)
	}
}

func TestFindAt_Found(t *testing.T) {
	testutil.HermeticDataDir(t)
	installPath := testutil.MakeFakeInstall(t, "mytool", "1.0.0", "foo")

	got, err := binary.FindAt(context.Background(), "mytool", "foo", "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != filepath.Join(installPath, "bin", "foo") {
		t.Errorf("BinPath = %q", got.BinPath)
	}
	if got.Plugin == nil || got.Plugin.Name != "mytool" || got.Version != "1.0.0" {
		t.Errorf("got %+v", got)
	}
}

func TestFindAt_SystemVersion(t *testing.T) {
	testutil.HermeticDataDir(t)
	sysDir := t.TempDir()
	sysFoo := filepath.Join(sysDir, "foo")
	testutil.WriteFile(t, sysFoo, "x")
	t.Setenv("PATH", sysDir)

	got, err := binary.FindAt(context.Background(), "mytool", "foo", "system")
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != sysFoo {
		t.Errorf("BinPath = %q, want %q", got.BinPath, sysFoo)
	}
}

func TestFindAt_NotInstalled(t *testing.T) {
	testutil.HermeticDataDir(t)
	_, err := binary.FindAt(context.Background(), "mytool", "foo", "1.0.0")
	if err == nil || !strings.Contains(err.Error(), "not installed") {
		t.Errorf("expected 'not installed' error, got %v", err)
	}
}

func TestFindAt_HonorsListBinPaths(t *testing.T) {
	testutil.HermeticDataDir(t)
	// Plugin declares non-default bin path via list-bin-paths.
	testutil.MakeFakePlugin(t, "weirdo", map[string]string{
		"install":        "#!/bin/sh\nexit 0\n",
		"list-bin-paths": "#!/bin/sh\necho \"bin tools\"\n",
	})
	// Binary lives in tools/, not bin/.
	installPath := testutil.MakeFakeInstall(t, "weirdo", "1.0.0")
	binPath := filepath.Join(installPath, "tools", "bar")
	testutil.WriteFile(t, binPath, "x")

	got, err := binary.FindAt(context.Background(), "weirdo", "bar", "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != binPath {
		t.Errorf("BinPath = %q, want %q", got.BinPath, binPath)
	}
}

func TestFindAt_HonorsPluginShimsDir(t *testing.T) {
	testutil.HermeticDataDir(t)
	// Plugin ships a pre-built shim at <plugin>/shims/myhelper, no install dir entry.
	testutil.MakeFakePlugin(t, "weirdo", map[string]string{
		"install": "#!/bin/sh\nexit 0\n",
	})
	helperPath := filepath.Join(config.PluginsDir(), "weirdo", "shims", "myhelper")
	testutil.WriteFile(t, helperPath, "#!/bin/sh\necho hi\n")
	// Install dir exists but is empty.
	testutil.MakeFakeInstall(t, "weirdo", "1.0.0")

	got, err := binary.FindAt(context.Background(), "weirdo", "myhelper", "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != helperPath {
		t.Errorf("BinPath = %q, want %q", got.BinPath, helperPath)
	}
}

func TestFindAt_HonorsExecPath(t *testing.T) {
	testutil.HermeticDataDir(t)
	// exec-path receives ($1=installPath, $2=binary, $3=relativePath) and prints a new relative path.
	testutil.MakeFakePlugin(t, "weirdo", map[string]string{
		"install":   "#!/bin/sh\nexit 0\n",
		"exec-path": "#!/bin/sh\necho \"alt/$2\"\n",
	})
	// Binary in default bin/, but exec-path will rewrite to alt/bar.
	installPath := testutil.MakeFakeInstall(t, "weirdo", "1.0.0", "bar")
	want := filepath.Join(installPath, "alt", "bar")
	testutil.WriteFile(t, want, "x")

	got, err := binary.FindAt(context.Background(), "weirdo", "bar", "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != want {
		t.Errorf("BinPath = %q, want %q", got.BinPath, want)
	}
}

func TestFindAt_SkipsNonExecutable(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakePlugin(t, "weirdo", map[string]string{
		"install": "#!/bin/sh\nexit 0\n",
		// NOTE: We order it such that the "bin" path is traversed first so we can confirm we skip the stub
		"list-bin-paths": "#!/bin/sh\necho \"bin tools\"\n",
	})

	installPath := testutil.MakeFakeInstall(t, "weirdo", "1.0.0")
	// Non-executable file in bin/<binary> — should be skipped.
	stub := filepath.Join(installPath, "bin", "node")
	if err := os.MkdirAll(filepath.Dir(stub), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(stub, []byte("template"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Real executable in tools/<binary>.
	real := filepath.Join(installPath, "tools", "node")
	testutil.WriteFile(t, real, "x")

	got, err := binary.FindAt(context.Background(), "weirdo", "node", "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if got.BinPath != real {
		t.Errorf("BinPath = %q, want %q (should skip non-exec stub)", got.BinPath, real)
	}
}

func TestPluginForShim_Found(t *testing.T) {
	testutil.HermeticDataDir(t)
	testutil.MakeFakeShims(t, "nodejs", "node")

	got, err := binary.PluginForShim("node")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "nodejs" {
		t.Errorf("Name = %q, want %q", got.Name, "nodejs")
	}
}

func TestPluginForShim_NoShim(t *testing.T) {
	testutil.HermeticDataDir(t)
	_, err := binary.PluginForShim("ghost")
	if err == nil || !strings.Contains(err.Error(), "no shim") {
		t.Errorf("expected 'no shim' error, got %v", err)
	}
}
