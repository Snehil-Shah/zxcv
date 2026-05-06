// Package plugin manages asdf plugins.
package plugin

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// IsLatestString reports whether v is the literal "latest" or "latest:<prefix>".
func IsLatestString(v string) bool {
	return v == "latest" || strings.HasPrefix(v, "latest:")
}

// Latest returns the latest stable version matching spec, where spec is "latest" or "latest:<prefix>".
func (p *Plugin) Latest(ctx context.Context, spec string) (string, error) {
	prefix := strings.TrimPrefix(strings.TrimPrefix(spec, "latest"), ":")
	stdout, _, err := p.RunCallback(ctx, "latest-stable", []string{prefix}, nil)
	if err == nil {
		v := strings.TrimSpace(string(stdout))
		if v == "" {
			return "", fmt.Errorf("latest-stable returned no version for prefix %q", prefix)
		}
		return v, nil
	}
	if !errors.Is(err, ErrCallbackNotImplemented) {
		return "", fmt.Errorf("latest-stable: %w", err)
	}

	// Fallback: derive from bin/list-all output.
	versions, err := p.ListAll(ctx)
	if err != nil {
		return "", err
	}
	return pickLatest(versions, prefix)
}

// pickLatest returns the last non-prerelease version matching prefix.
func pickLatest(versions []string, prefix string) (string, error) {
	var match string
	for _, v := range versions {
		if prefix != "" && !strings.HasPrefix(v, prefix) {
			continue
		}
		if isPrerelease(v) {
			continue
		}
		match = v
	}
	if match == "" {
		return "", fmt.Errorf("no stable version matching prefix %q", prefix)
	}
	return match, nil
}

// isPrerelease reports whether v looks like a SemVer-style prerelease.
// HACK: obviously unreliable
func isPrerelease(v string) bool {
	lower := strings.ToLower(v)
	for _, marker := range []string{"-alpha", "-beta", "-rc", "-pre", "-dev"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}
