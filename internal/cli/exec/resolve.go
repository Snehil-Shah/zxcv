// Package exec implements the `zxcv exec` command.
package exec

import (
	"context"
	"fmt"
	"os"

	"github.com/Snehil-Shah/zxcv/internal/binary"
)

// runResolved resolves plugin's binary at cwd's resolved version and execs it.
func runResolved(ctx context.Context, plugin, binName string, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd: %w", err)
	}
	r, err := binary.Find(ctx, plugin, binName, cwd)
	if err != nil {
		return err
	}
	return execBinary(r, args)
}
