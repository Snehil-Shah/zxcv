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

func TestGenerate_DefaultsToBin(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{"install": "#!/bin/sh\nexit 0\n"})
	installPath := testutil.MakeFakeInstall(t, "mytool", "1.0.0", "foo", "bar")
	shimsDir := config.ShimsDir()

	if err := binary.Generate(context.Background(), p, "1.0.0", installPath); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"foo", "bar"} {
		path := filepath.Join(shimsDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected shim at %s: %v", path, err)
			continue
		}
		body := string(data)
		if !strings.Contains(body, "# zxcv-plugin: mytool") {
			t.Errorf("shim %s: missing plugin comment, got:\n%s", name, body)
		}
		if !strings.Contains(body, `exec zxcv exec "mytool" "`+name+`"`) {
			t.Errorf("shim %s: missing exec line, got:\n%s", name, body)
		}
	}
}

func TestGenerate_SkipsNonExecutable(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{"install": "#!/bin/sh\nexit 0\n"})
	installPath := testutil.MakeFakeInstall(t, "mytool", "1.0.0", "real")
	// Non-executable file (e.g., template, README) — should NOT get a shim.
	if err := os.WriteFile(filepath.Join(installPath, "bin", "README"), []byte("docs"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := binary.Generate(context.Background(), p, "1.0.0", installPath); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(config.ShimsDir(), "real")); err != nil {
		t.Errorf("expected shim for real binary: %v", err)
	}
	if _, err := os.Stat(filepath.Join(config.ShimsDir(), "README")); err == nil {
		t.Errorf("non-executable README should not have a shim")
	}
}

func TestGenerate_HonorsListBinPaths(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{
		"install":        "#!/bin/sh\nexit 0\n",
		"list-bin-paths": "#!/bin/sh\necho \"bin tools\"\n",
	})
	installPath := testutil.MakeFakeInstall(t, "mytool", "1.0.0", "foo")
	shimsDir := config.ShimsDir()
	testutil.WriteFile(t, filepath.Join(installPath, "tools", "bar"), "x")

	if err := binary.Generate(context.Background(), p, "1.0.0", installPath); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"foo", "bar"} {
		if _, err := os.Stat(filepath.Join(shimsDir, name)); err != nil {
			t.Errorf("missing shim for %s: %v", name, err)
		}
	}
}
