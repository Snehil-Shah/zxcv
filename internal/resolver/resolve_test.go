// Package resolver_test tests the `resolver` package.
package resolver_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/resolver"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestResolve(t *testing.T) {
	testutil.HermeticDataDir(t)

	// Tree:
	//   root/.tool-versions          python 3.10.0, ruby 3.0.0
	//   root/project/.tool-versions  nodejs 20.11.0, ruby 3.2.0  (overrides ruby)
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, ".tool-versions"), "python 3.10.0\nruby 3.0.0\n")

	project := filepath.Join(root, "project")
	if err := os.Mkdir(project, 0o755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(project, ".tool-versions"), "nodejs 20.11.0\nruby 3.2.0\n")

	resolutions, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	got := map[string]resolver.Resolution{}
	for _, r := range resolutions {
		got[r.Name] = r
	}

	cases := []struct {
		tool    string
		version string
		source  string
	}{
		{"nodejs", "20.11.0", filepath.Join(project, ".tool-versions")},
		{"ruby", "3.2.0", filepath.Join(project, ".tool-versions")},
		{"python", "3.10.0", filepath.Join(root, ".tool-versions")},
	}
	for _, c := range cases {
		r, ok := got[c.tool]
		if !ok {
			t.Errorf("missing resolution for %q", c.tool)
			continue
		}
		if r.Version != c.version {
			t.Errorf("%s: version = %q, want %q", c.tool, r.Version, c.version)
		}
		if r.Source != c.source {
			t.Errorf("%s: source = %q, want %q", c.tool, r.Source, c.source)
		}
	}
}

func TestResolve_WithDefinitions(t *testing.T) {
	testutil.HermeticDataDir(t)

	// Tree:
	//   root/.tool-versions          mytool 1.0.0
	//   root/.tool-definitions       mytool https://github.com/me/asdf-mytool.git
	//   root/project/.tool-versions  other 2.0.0   (no sibling .tool-definitions)
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, ".tool-versions"), "mytool 1.0.0\n")
	testutil.WriteFile(t, filepath.Join(root, ".tool-definitions"), "mytool https://github.com/me/asdf-mytool.git\n")

	project := filepath.Join(root, "project")
	if err := os.Mkdir(project, 0o755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(project, ".tool-versions"), "other 2.0.0\n")

	resolutions, err := resolver.Resolve(project)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	got := map[string]resolver.Resolution{}
	for _, r := range resolutions {
		got[r.Name] = r
	}

	mytool, ok := got["mytool"]
	if !ok {
		t.Fatal("missing resolution for mytool")
	}
	if mytool.Definition == nil {
		t.Error("mytool: Definition is nil, want non-nil")
	} else if mytool.Definition.URL != "https://github.com/me/asdf-mytool.git" {
		t.Errorf("mytool: Definition.URL = %q, want git URL", mytool.Definition.URL)
	}

	other, ok := got["other"]
	if !ok {
		t.Fatal("missing resolution for other")
	}
	if other.Definition != nil {
		t.Errorf("other: Definition = %+v, want nil (no sibling .tool-definitions)", other.Definition)
	}
}

func TestResolveVersion_FromCwd(t *testing.T) {
	testutil.HermeticDataDir(t)
	cwd := t.TempDir()
	testutil.WriteFile(t, filepath.Join(cwd, ".tool-versions"), "nodejs 20.11.0\nruby 3.2.0\n")

	got, err := resolver.ResolveVersion(cwd, "ruby")
	if err != nil {
		t.Fatal(err)
	}
	if got != "3.2.0" {
		t.Errorf("got %q, want %q", got, "3.2.0")
	}
}

func TestResolveVersion_FromParent(t *testing.T) {
	testutil.HermeticDataDir(t)
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, ".tool-versions"), "python 3.10.0\n")
	project := filepath.Join(root, "project")
	if err := os.Mkdir(project, 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := resolver.ResolveVersion(project, "python")
	if err != nil {
		t.Fatal(err)
	}
	if got != "3.10.0" {
		t.Errorf("got %q, want %q", got, "3.10.0")
	}
}

func TestResolveVersion_NearerWins(t *testing.T) {
	testutil.HermeticDataDir(t)
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, ".tool-versions"), "ruby 3.0.0\n")
	project := filepath.Join(root, "project")
	if err := os.Mkdir(project, 0o755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(project, ".tool-versions"), "ruby 3.2.0\n")

	got, err := resolver.ResolveVersion(project, "ruby")
	if err != nil {
		t.Fatal(err)
	}
	if got != "3.2.0" {
		t.Errorf("nearer .tool-versions should win: got %q, want %q", got, "3.2.0")
	}
}

func TestResolveVersion_FromHome(t *testing.T) {
	home := testutil.HermeticDataDir(t)
	homeDir := filepath.Dir(home)
	testutil.WriteFile(t, filepath.Join(homeDir, ".tool-versions"), "node 18.0.0\n")

	cwd := t.TempDir()
	got, err := resolver.ResolveVersion(cwd, "node")
	if err != nil {
		t.Fatal(err)
	}
	if got != "18.0.0" {
		t.Errorf("got %q, want %q", got, "18.0.0")
	}
}

func TestResolveVersion_NotFound(t *testing.T) {
	testutil.HermeticDataDir(t)
	cwd := t.TempDir()
	testutil.WriteFile(t, filepath.Join(cwd, ".tool-versions"), "nodejs 20.0.0\n")

	_, err := resolver.ResolveVersion(cwd, "ghost")
	if err == nil || !strings.Contains(err.Error(), "no version of ghost set") {
		t.Errorf("expected 'no version of ghost set' error, got %v", err)
	}
}
