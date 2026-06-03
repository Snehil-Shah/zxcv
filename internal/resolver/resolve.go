// Package resolver resolves the applicable manifest files.
package resolver

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Snehil-Shah/zxcv/internal/config"
	"github.com/Snehil-Shah/zxcv/internal/tooldefinitions"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

// Resolve walks from startDir up to the filesystem root to collect all tool versions.
func Resolve(startDir string) ([]Resolution, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return nil, err
	}

	var resolutions []Resolution
	seen := make(map[string]struct{})

	// record adds tools from a manifest to the resolutions list, skipping any already seen.
	record := func(path string, tools toolversions.ToolVersions, defs tooldefinitions.ToolDefinitions) {
		for _, t := range tools {
			if _, ok := seen[t.Name]; ok {
				continue
			}
			seen[t.Name] = struct{}{}
			var def *tooldefinitions.Definition
			if d, found := defs.Lookup(t.Name); found {
				def = &d
			}
			resolutions = append(resolutions, Resolution{Tool: t, Source: path, Definition: def})
		}
	}

	// ingest records a dir.
	ingest := func(d string) error {
		tools, err := loadVersions(d)
		if err != nil {
			return err
		}
		defs, err := loadDefinitions(d)
		if err != nil {
			return err
		}
		record(filepath.Join(d, manifestName), tools, defs)
		return nil
	}

	// Trip to the root!
	for {
		if err := ingest(dir); err != nil {
			return nil, err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if err := ingest(config.GlobalDir()); err != nil {
		return nil, err
	}

	return resolutions, nil
}

// ResolveVersion returns the active version for tool name (based on manifests).
func ResolveVersion(startDir, name string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}
	for {
		tools, err := loadVersions(dir)
		if err != nil {
			return "", err
		}
		if t, ok := tools.Lookup(name); ok {
			return t.Version, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	tools, err := loadVersions(config.GlobalDir())
	if err != nil {
		return "", err
	}
	if t, ok := tools.Lookup(name); ok {
		return t.Version, nil
	}
	return "", fmt.Errorf("no version of %s set", name)
}

// loadVersions parses the `.tool-versions` in dir.
func loadVersions(dir string) (toolversions.ToolVersions, error) {
	tools, err := toolversions.ParseFile(filepath.Join(dir, manifestName))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return tools, nil
}

// loadDefinitions parses the `.tool-definitions` in dir.
func loadDefinitions(dir string) (tooldefinitions.ToolDefinitions, error) {
	defs, err := tooldefinitions.ParseFile(filepath.Join(dir, definitionsName))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return defs, nil
}
