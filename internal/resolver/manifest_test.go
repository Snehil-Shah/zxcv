// Package resolver_test tests the `resolver` package.
package resolver_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/resolver"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestGlobalManifest(t *testing.T) {
	home := testutil.HermeticDataDir(t)
	homeDir := filepath.Dir(home)

	got := resolver.GlobalManifest()
	want := filepath.Join(homeDir, ".tool-versions")
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplicableManifest_FromCwd(t *testing.T) {
	testutil.HermeticDataDir(t)
	cwd := t.TempDir()
	manifest := filepath.Join(cwd, ".tool-versions")
	testutil.WriteFile(t, manifest, "nodejs 20.0.0\n")

	if got := resolver.ApplicableManifest(cwd); got != manifest {
		t.Errorf("got %q, want %q", got, manifest)
	}
}

func TestApplicableManifest_FromParent(t *testing.T) {
	testutil.HermeticDataDir(t)
	root := t.TempDir()
	parentManifest := filepath.Join(root, ".tool-versions")
	testutil.WriteFile(t, parentManifest, "python 3.11.0\n")

	project := filepath.Join(root, "project")
	if err := os.Mkdir(project, 0o755); err != nil {
		t.Fatal(err)
	}

	if got := resolver.ApplicableManifest(project); got != parentManifest {
		t.Errorf("got %q, want %q (should walk up)", got, parentManifest)
	}
}

func TestApplicableManifest_NoneFound(t *testing.T) {
	testutil.HermeticDataDir(t)
	cwd := t.TempDir()
	want := filepath.Join(cwd, ".tool-versions")

	if got := resolver.ApplicableManifest(cwd); got != want {
		t.Errorf("got %q, want %q (should default to cwd)", got, want)
	}
}
