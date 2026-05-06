// Package plugin_test tests the `plugin` package.
package plugin_test

import (
	"context"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestLatest_ViaLatestStable(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{
		"install":       "#!/bin/sh\n",
		"latest-stable": "#!/bin/sh\necho \"1.2.3\"\n",
	})
	got, err := p.Latest(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.2.3" {
		t.Errorf("got %q, want %q", got, "1.2.3")
	}
}

func TestLatest_PrefixPassedThrough(t *testing.T) {
	testutil.HermeticDataDir(t)
	// latest-stable echoes its prefix arg back so we can verify it's piped through.
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{
		"install":       "#!/bin/sh\n",
		"latest-stable": "#!/bin/sh\necho \"$1.99\"\n",
	})
	got, err := p.Latest(context.Background(), "20")
	if err != nil {
		t.Fatal(err)
	}
	if got != "20.99" {
		t.Errorf("got %q, want %q", got, "20.99")
	}
}

func TestLatest_FallbackToListAll(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{
		"install":  "#!/bin/sh\n",
		"list-all": "#!/bin/sh\necho \"1.0.0 1.0.1 1.1.0-rc1 1.1.0\"\n",
	})
	got, err := p.Latest(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.1.0" {
		t.Errorf("got %q, want %q", got, "1.1.0")
	}
}

func TestLatest_FallbackPrefix(t *testing.T) {
	testutil.HermeticDataDir(t)
	p := testutil.MakeFakePlugin(t, "mytool", map[string]string{
		"install":  "#!/bin/sh\n",
		"list-all": "#!/bin/sh\necho \"1.0.0 2.0.0 2.1.0 3.0.0\"\n",
	})
	got, err := p.Latest(context.Background(), "2")
	if err != nil {
		t.Fatal(err)
	}
	if got != "2.1.0" {
		t.Errorf("got %q, want %q", got, "2.1.0")
	}
}
