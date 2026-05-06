// Package registry_test tests the `registry` package.
package registry_test

import (
	"reflect"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/registry"
	"github.com/Snehil-Shah/zxcv/internal/testutil"
)

func TestSearch_NoRegistry(t *testing.T) {
	testutil.HermeticDataDir(t)

	got, err := registry.Search("anything")
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestSearch_SubstringCaseInsensitive(t *testing.T) {
	testutil.HermeticDataDir(t)
	writePluginEntry(t, "node", "repository = x\n")
	writePluginEntry(t, "nodejs", "repository = x\n")
	writePluginEntry(t, "python", "repository = x\n")

	got, err := registry.Search("NODE")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"node", "nodejs"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSearch_Sorted(t *testing.T) {
	testutil.HermeticDataDir(t)
	for _, name := range []string{"zig", "alpha", "mango"} {
		writePluginEntry(t, name, "repository = x\n")
	}

	got, err := registry.Search("")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"alpha", "mango", "zig"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
