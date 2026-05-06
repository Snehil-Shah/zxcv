// Package plugin_test tests the `plugin` package.
package plugin_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestRunCallback_NotFound(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "foo", map[string]string{"install": "#!/bin/sh\n"})
	_, _, err := p.RunCallback(context.Background(), "nonexistent", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "not implemented") {
		t.Errorf("expected 'not implemented' error, got %v", err)
	}
}

func TestRunCallback_Success(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "foo", map[string]string{
		"install":  "#!/bin/sh\n",
		"list-all": "#!/bin/sh\necho \"$TEST_VAR 1.0 2.0\"\n",
	})
	stdout, _, err := p.RunCallback(context.Background(), "list-all", nil, map[string]string{
		"TEST_VAR": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(string(stdout)), "foo 1.0 2.0"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}
