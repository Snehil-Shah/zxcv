// Package tooldefinitions_test tests the `tooldefinitions` package.
package tooldefinitions_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/tooldefinitions"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    tooldefinitions.ToolDefinitions
	}{
		{
			name:    "basic",
			fixture: "basic.tool-definitions",
			want: tooldefinitions.ToolDefinitions{
				{Name: "internal-cli", URL: "https://github.com/ourorg/asdf-internal-cli.git"},
				{Name: "legacy-plugin", URL: "git@github.com:legacy-org/asdf-legacy.git"},
				{Name: "my-plugin", Path: "/Users/me/work/asdf-my-plugin"},
				{Name: "our-cli", Path: filepath.Join("testdata", "tooling", "asdf-our-cli")},
			},
		},
		{
			name:    "with comments",
			fixture: "with-comments.tool-definitions",
			want: tooldefinitions.ToolDefinitions{
				{Name: "internal-cli", URL: "https://github.com/ourorg/asdf-internal-cli.git"},
				{Name: "my-plugin", Path: "/Users/me/work/asdf-my-plugin"},
			},
		},
		{
			name:    "messy whitespace",
			fixture: "messy.tool-definitions",
			want: tooldefinitions.ToolDefinitions{
				{Name: "internal-cli", URL: "https://github.com/ourorg/asdf-internal-cli.git"},
				{Name: "our-cli", Path: filepath.Join("testdata", "tooling", "asdf-our-cli")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tooldefinitions.ParseFile(filepath.Join("testdata", tt.fixture))
			if err != nil {
				t.Fatalf("ParseFile(%s): %v", tt.fixture, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := tooldefinitions.ParseFile(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseFile_Invalid(t *testing.T) {
	_, err := tooldefinitions.ParseFile(filepath.Join("testdata", "invalid.tool-definitions"))
	if err == nil {
		t.Error("expected error for invalid fixture, got nil")
	}
}

func TestToolDefinitionsLookup(t *testing.T) {
	d := tooldefinitions.ToolDefinitions{
		{Name: "internal-cli", URL: "https://example.com/p.git"},
		{Name: "our-cli", Path: "/tmp/our-cli"},
	}

	got, ok := d.Lookup("our-cli")
	if !ok || got.Path != "/tmp/our-cli" {
		t.Errorf("Lookup(our-cli) = (%+v, %v), want matching our-cli", got, ok)
	}
	if _, ok := d.Lookup("missing"); ok {
		t.Error("Lookup(missing) ok=true, want false")
	}
}
