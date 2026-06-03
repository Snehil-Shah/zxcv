// Package plugin_test tests the `plugin` package.
package plugin_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestAsdfShim_Version(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "test", map[string]string{
		"check": "#!/bin/sh\nasdf version\n",
	})
	stdout, _, err := p.RunCallback(context.Background(), "check", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(stdout)); got != "v0.19.0" {
		t.Errorf("stdout = %q, want v0.19.0", got)
	}
}

func TestAsdfShim_Reshim(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "test", map[string]string{
		"check": "#!/bin/sh\nasdf reshim foo 1.0\n",
	})
	if _, _, err := p.RunCallback(context.Background(), "check", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAsdfShim_Unknown(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "test", map[string]string{
		"check": "#!/bin/sh\nasdf nonsense\n",
	})
	_, stderr, err := p.RunCallback(context.Background(), "check", nil, nil)
	if err == nil {
		t.Fatal("expected non-nil error from unknown subcommand")
	}
	if !strings.Contains(string(stderr), "not supported") {
		t.Errorf("stderr = %q, want substring 'not supported'", stderr)
	}
}

func TestAsdfShim_Current(t *testing.T) {
	testutil.HermeticDataDir(t)
	stubDir := t.TempDir()
	testutil.WriteFile(t, filepath.Join(stubDir, "zxcv"), "#!/bin/sh\necho \"9.9.9\"\n")
	t.Setenv("PATH", stubDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	p := testutil.MakeFakePlugin(t, "test", map[string]string{
		"check": "#!/bin/sh\nasdf current --no-header elixir\n",
	})
	stdout, _, err := p.RunCallback(context.Background(), "check", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(stdout)); got != "elixir 9.9.9" {
		t.Errorf("stdout = %q, want %q", got, "elixir 9.9.9")
	}
}

func TestAsdfShim_CurrentRequiresTool(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "test", map[string]string{
		"check": "#!/bin/sh\nasdf current --no-header\n",
	})
	_, stderr, err := p.RunCallback(context.Background(), "check", nil, nil)
	if err == nil {
		t.Fatal("expected non-nil error when no tool arg given")
	}
	if !strings.Contains(string(stderr), "requires a tool") {
		t.Errorf("stderr = %q, want substring 'requires a tool'", stderr)
	}
}
