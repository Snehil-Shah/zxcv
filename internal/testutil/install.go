// Package testutil provides shared test helpers.
package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// MakeFakeInstall creates the install dir at installsDir/<name>/<version>/ with an
// executable bin/<binary> for each binary. Returns the install path.
func MakeFakeInstall(t *testing.T, name, version string, binaries ...string) string {
	t.Helper()
	installPath := filepath.Join(config.InstallsDir(), name, version)
	if err := os.MkdirAll(installPath, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, binary := range binaries {
		WriteFile(t, filepath.Join(installPath, "bin", binary), "x")
	}
	return installPath
}

// MakeFakeShims writes a shim stub for plugin at shimsDir/<binary> for each binary.
func MakeFakeShims(t *testing.T, plugin string, binaries ...string) {
	t.Helper()
	body := "#!/bin/sh\n# zxcv-plugin: " + plugin + "\n"
	for _, binary := range binaries {
		WriteFile(t, filepath.Join(config.ShimsDir(), binary), body)
	}
}
