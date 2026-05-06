// Package toolversions_test tests the `toolversions` package.
package toolversions_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

func TestParseFile(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    toolversions.ToolVersions
	}{
		{
			name:    "basic",
			fixture: "basic.tool-versions",
			want: toolversions.ToolVersions{
				{Name: "nodejs", Version: "20.11.0"},
				{Name: "python", Version: "3.11.5"},
				{Name: "ruby", Version: "3.2.0"},
			},
		},
		{
			name:    "with comments",
			fixture: "with-comments.tool-versions",
			want: toolversions.ToolVersions{
				{Name: "nodejs", Version: "20.11.0"},
				{Name: "python", Version: "3.11.5"},
				{Name: "ruby", Version: "3.2.0"},
			},
		},
		{
			name:    "messy whitespace, blanks, orphan comment",
			fixture: "messy.tool-versions",
			want: toolversions.ToolVersions{
				{Name: "nodejs", Version: "20.11.0"},
				{Name: "python", Version: "3.11.5"},
				{Name: "go", Version: "1.21.5"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toolversions.ParseFile(filepath.Join("testdata", tt.fixture))
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
	_, err := toolversions.ParseFile(filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParseFile_Invalid(t *testing.T) {
	_, err := toolversions.ParseFile(filepath.Join("testdata", "invalid.tool-versions"))
	if err == nil {
		t.Error("expected error for invalid fixture, got nil")
	}
}

func TestToolVersionsLookup(t *testing.T) {
	tv := toolversions.ToolVersions{
		{Name: "nodejs", Version: "20.11.0"},
		{Name: "python", Version: "3.11.5"},
	}

	if got, ok := tv.Lookup("python"); !ok || got.Version != "3.11.5" {
		t.Errorf("Lookup(python) = (%+v, %v), want ({python 3.11.5}, true)", got, ok)
	}
	if _, ok := tv.Lookup("ruby"); ok {
		t.Errorf("Lookup(ruby) should not be found")
	}
}
