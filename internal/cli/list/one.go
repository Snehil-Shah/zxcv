// Package list implements the `zxcv list` command.
package list

import (
	"fmt"
	"sort"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// listOne prints all installed versions for one tool, one per line.
func listOne(name string) error {
	versions, err := installer.VersionsFor(name)
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		ui.Info(fmt.Sprintf("no versions of %s installed", name))
		return nil
	}
	sort.Strings(versions)
	for _, v := range versions {
		ui.Write(v)
	}
	return nil
}
