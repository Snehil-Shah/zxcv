// Package installer_test tests the `installer` package.
package installer_test

import (
	"context"
	_ "embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
	"github.com/Snehil-Shah/zxcv/internal/tooldefinitions"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

//go:embed testdata/install_success.sh
var scriptInstallSuccess string

//go:embed testdata/install_fail_partial.sh
var scriptInstallFailPartial string

//go:embed testdata/install_marker.sh
var scriptInstallMarker string

func TestInstall_Success(t *testing.T) {
	testutil.HermeticDataDir(t)
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{"install": scriptInstallSuccess})
	target := installer.Target{
		Tool:       toolversions.Tool{Name: "mytool", Version: "1.0.0"},
		Definition: &tooldefinitions.Definition{Name: "mytool", Path: src},
	}
	res := installer.New().Install(context.Background(), []installer.Target{target})[0]
	if res.Err != nil {
		t.Fatalf("Install: %v\nstdout: %s\nstderr: %s", res.Err, res.Stdout, res.Stderr)
	}
	bin := filepath.Join(config.InstallsDir(), "mytool", "1.0.0", "bin", "mytool")
	if _, err := os.Stat(bin); err != nil {
		t.Errorf("expected installed binary at %s, got %v", bin, err)
	}
	if _, err := os.Stat(filepath.Join(config.ShimsDir(), "mytool")); err != nil {
		t.Errorf("expected shim for mytool: %v", err)
	}
}

func TestInstall_FailureCleansPartial(t *testing.T) {
	testutil.HermeticDataDir(t)
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{"install": scriptInstallFailPartial})
	target := installer.Target{
		Tool:       toolversions.Tool{Name: "mytool", Version: "1.0.0"},
		Definition: &tooldefinitions.Definition{Name: "mytool", Path: src},
	}
	res := installer.New().Install(context.Background(), []installer.Target{target})[0]
	if res.Err == nil {
		t.Fatal("expected error, got nil")
	}
	installed := filepath.Join(config.InstallsDir(), "mytool", "1.0.0")
	if _, err := os.Stat(installed); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("install dir should not exist: %v", err)
	}
	if _, err := os.Stat(installed + ".partial"); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("partial dir should be cleaned up: %v", err)
	}
}

func TestInstall_Parallel(t *testing.T) {
	testutil.HermeticDataDir(t)
	names := []string{"alpha", "beta", "gamma"}
	targets := make([]installer.Target, len(names))
	for i, name := range names {
		src := t.TempDir()
		testutil.MakeFakePluginAt(t, src, map[string]string{"install": scriptInstallMarker})
		targets[i] = installer.Target{
			Tool:       toolversions.Tool{Name: name, Version: "1.0.0"},
			Definition: &tooldefinitions.Definition{Name: name, Path: src},
		}
	}
	results := installer.New().Install(context.Background(), targets)
	if len(results) != len(names) {
		t.Fatalf("got %d results, want %d", len(results), len(names))
	}
	for i, r := range results {
		if r.Err != nil {
			t.Errorf("result %d: %v\nstderr: %s", i, r.Err, r.Stderr)
		}
	}
	for _, name := range names {
		path := filepath.Join(config.InstallsDir(), name, "1.0.0", "bin", "v")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing install for %s: %v", name, err)
		}
	}
}

func TestInstall_FiresPostPluginAdd(t *testing.T) {
	testutil.HermeticDataDir(t)
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{
		"install":         "#!/bin/sh\nmkdir -p \"$ASDF_INSTALL_PATH\"\n",
		"post-plugin-add": "#!/bin/sh\ntouch \"$HOME/.post-plugin-add-fired\"\n",
	})
	target := installer.Target{
		Tool:       toolversions.Tool{Name: "mytool", Version: "1.0.0"},
		Definition: &tooldefinitions.Definition{Name: "mytool", Path: src},
	}
	res := installer.New().Install(context.Background(), []installer.Target{target})[0]
	if res.Err != nil {
		t.Fatalf("Install: %v\nstderr: %s", res.Err, res.Stderr)
	}
	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".post-plugin-add-fired")); err != nil {
		t.Errorf("post-plugin-add did not fire: %v", err)
	}
}

func TestInstall_ExpandsLatest(t *testing.T) {
	testutil.HermeticDataDir(t)
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{
		"install":       "#!/bin/sh\nmkdir -p \"$ASDF_INSTALL_PATH\"\n",
		"latest-stable": "#!/bin/sh\necho \"5.0.0\"\n",
	})
	target := installer.Target{
		Tool:       toolversions.Tool{Name: "mytool", Version: "latest"},
		Definition: &tooldefinitions.Definition{Name: "mytool", Path: src},
	}
	res := installer.New().Install(context.Background(), []installer.Target{target})[0]
	if res.Err != nil {
		t.Fatalf("Install: %v\nstderr: %s", res.Err, res.Stderr)
	}
	if res.Target.Version != "5.0.0" {
		t.Errorf("Result.Target.Version = %q, want %q (latest should resolve)", res.Target.Version, "5.0.0")
	}
	if _, err := os.Stat(filepath.Join(config.InstallsDir(), "mytool", "5.0.0")); err != nil {
		t.Errorf("expected install at concrete version, got %v", err)
	}
}

func TestInstall_SkipsSystemVersion(t *testing.T) {
	testutil.HermeticDataDir(t)
	// Plugin's install script would exit non-zero if tried to be installed.
	src := t.TempDir()
	testutil.MakeFakePluginAt(t, src, map[string]string{"install": "#!/bin/sh\nexit 1\n"})
	target := installer.Target{
		Tool:       toolversions.Tool{Name: "mytool", Version: "system"},
		Definition: &tooldefinitions.Definition{Name: "mytool", Path: src},
	}

	res := installer.New().Install(context.Background(), []installer.Target{target})[0]
	if res.Err != nil {
		t.Fatalf("expected no-op for system version, got error: %v", res.Err)
	}
	if _, err := os.Stat(filepath.Join(config.InstallsDir(), "mytool", "system")); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("system version should not create an install dir, got err=%v", err)
	}
}
