// Package testutil provides shared test helpers.
package testutil

import (
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// HermeticDataDir mocks our data dir to point to a temp dir.
func HermeticDataDir(t *testing.T) string {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	return config.DataDir()
}
