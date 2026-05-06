// Package registry maintains the asdf plugin registry.
package registry

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/go-git/go-git/v5"

	"github.com/Snehil-Shah/zxcv/internal/config"
)

// Refresh force-refreshes the index, clones if missing, pulls otherwise.
func Refresh(ctx context.Context) error {
	dir := config.RegistryDir()
	_, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return clone(ctx, dir)
	}
	if err != nil {
		return fmt.Errorf("stat %s: %w", dir, err)
	}
	if pullErr := pull(ctx, dir); pullErr != nil {
		if rmErr := os.RemoveAll(dir); rmErr != nil {
			return fmt.Errorf("recover from %v: %w", pullErr, rmErr)
		}
		return clone(ctx, dir)
	}
	return nil
}

// clone performs a shallow clone of the plugin index to dir.
func clone(ctx context.Context, dir string) error {
	_, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:   defaultURL,
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("clone %s: %w", defaultURL, err)
	}
	return nil
}

// pull resets the worktree (in case of go-git's filemode quirks) then pulls.
func pull(ctx context.Context, dir string) error {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("open plugin index: %w", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("worktree: %w", err)
	}
	if err := wt.Reset(&git.ResetOptions{Mode: git.HardReset}); err != nil {
		return fmt.Errorf("reset worktree: %w", err)
	}
	err = wt.PullContext(ctx, &git.PullOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("pull plugin index: %w", err)
	}
	return nil
}
