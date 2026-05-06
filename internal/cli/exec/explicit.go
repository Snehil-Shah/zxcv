// Package exec implements the `zxcv exec` command.
package exec

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Snehil-Shah/zxcv/internal/binary"
)

// runExplicit execs the under the stated plugin version, bypasses cwd resolutions.
func runExplicit(ctx context.Context, pluginAtVer, binName string, args []string) error {
	plugin, version, ok := strings.Cut(pluginAtVer, "@")
	if !ok || plugin == "" || version == "" {
		return fmt.Errorf("invalid form %q, expected <tool>@<version>", pluginAtVer)
	}
	if strings.HasPrefix(version, "latest") {
		return errors.New("'latest' not supported in exec, specify a concrete version")
	}
	r, err := binary.FindAt(ctx, plugin, binName, version)
	if err != nil {
		return err
	}
	return execBinary(r, args)
}
