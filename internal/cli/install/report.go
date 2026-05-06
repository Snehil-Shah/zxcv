// Package install implements the `zxcv install` command.
package install

import (
	"fmt"
	"strings"

	"github.com/Snehil-Shah/zxcv/internal/installer"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// labelsFor builds a label slice in the same order as targets, used by ui.Progress.
func labelsFor(targets []installer.Target) []string {
	labels := make([]string, len(targets))
	for i, t := range targets {
		labels[i] = t.Name + "@" + t.Version
	}
	return labels
}

// progressCallback returns an OnComplete that updates prog as each target finishes.
func progressCallback(prog *ui.Progress) func(int, installer.Result) {
	return func(idx int, r installer.Result) {
		if r.Err == nil {
			prog.Done(idx)
		} else {
			prog.Failed(idx)
		}
	}
}

// finalize prints the summary line, dumps stderr for failed targets, and returns an aggregate error.
func finalize(results []installer.Result) error {
	var failed []installer.Result
	for _, r := range results {
		if r.Err != nil {
			failed = append(failed, r)
		}
	}
	succeeded := len(results) - len(failed)
	summary := fmt.Sprintf("%d/%d tools installed", succeeded, len(results))

	ui.Info("")
	if len(failed) == 0 {
		ui.Info(ui.Green("✓") + " " + summary)
		return nil
	}
	ui.Info(ui.Red("✗") + " " + summary + fmt.Sprintf(" (%d failed)", len(failed)))

	for _, r := range failed {
		ui.Error("")
		ui.Error(fmt.Sprintf("%s@%s: %s", r.Target.Name, r.Target.Version, r.Err))
		body := r.Stderr
		if len(body) == 0 {
			body = r.Stdout
		}
		if trimmed := strings.TrimRight(string(body), "\n"); trimmed != "" {
			ui.Info(ui.Dim(trimmed))
		}
	}
	ui.Error("")
	return fmt.Errorf("%d/%d installs failed", len(failed), len(results))
}
