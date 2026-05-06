// Package list implements the `zxcv list` command.
package list

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// group holds a tool and its installed versions.
type group struct {
	Name     string
	Versions []string
}

// listAll prints every installed tool with its versions, grouped tool-by-tool.
func listAll() error {
	installsDir := config.InstallsDir()
	entries, err := os.ReadDir(installsDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			ui.Info("no tools installed")
			return nil
		}
		return fmt.Errorf("read %s: %w", installsDir, err)
	}

	var groups []group
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		versions, err := installer.VersionsFor(e.Name())
		if err != nil {
			return err
		}
		if len(versions) > 0 {
			sort.Strings(versions)
			groups = append(groups, group{Name: e.Name(), Versions: versions})
		}
	}
	if len(groups) == 0 {
		ui.Info("no tools installed")
		return nil
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })

	for i, g := range groups {
		if i > 0 {
			ui.Write("")
		}
		ui.Write(g.Name)
		for _, v := range g.Versions {
			ui.Write(v)
		}
	}
	return nil
}
