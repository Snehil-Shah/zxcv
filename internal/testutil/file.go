// Package testutil provides shared test helpers.
package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// WriteFile creates a file at the given path with the given content.
func WriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}
