// Package main generates the auto-generated CLI reference markdown.
//
// Writes the root command to docs/usage.md and one file
// per subcommand under docs/reference/.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	docs "github.com/urfave/cli-docs/v3"
	cli "github.com/urfave/cli/v3"

	zcli "github.com/Snehil-Shah/zxcv/internal/cli"
)

// NOTE: This file hardcodes project's structural conventions.

const (
	usagePath = "docs/usage.md"
	refDir    = "docs/reference"

	// Banner prepended to every generated file. HTML comments don't render in mkdocs/GitHub.
	autogenBanner = "<!-- This file is auto-generated. Do not edit manually. -->\n\n"
)

// cli-docs/v3 emits a `# SYNOPSIS` section in awkward style, strip it.
var synopsisRe = regexp.MustCompile(`(?s)\n# SYNOPSIS\n.*?\n# `)

// Matches every Markdown heading.
var headingRe = regexp.MustCompile(`(?m)^(#+) `)

func main() {
	root := zcli.New("dev")
	if err := os.MkdirAll(refDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir:", err)
		os.Exit(1)
	}

	// Root command lives outside the reference folder as a top-level usage page.
	if err := write(usagePath, root); err != nil {
		fmt.Fprintln(os.Stderr, "usage:", err)
		os.Exit(1)
	}

	// One page per subcommand.
	for _, c := range root.Commands {
		if c.Hidden {
			continue
		}
		if err := write(filepath.Join(refDir, c.Name+".md"), c); err != nil {
			fmt.Fprintln(os.Stderr, c.Name+":", err)
			os.Exit(1)
		}
	}
}

// write writes a markdown reference file for the given command at the given path.
func write(path string, cmd *cli.Command) error {
	md, err := docs.ToMarkdown(cmd)
	if err != nil {
		return err
	}
	// Normalize headings, the generated ones are kinda awkward.
	md = synopsisRe.ReplaceAllString(md, "\n# ")
	md = strings.Replace(md, "# NAME\n\n", "", 1)
	md = headingRe.ReplaceAllString(md, "#$1 ")
	md = autogenBanner + "# " + cmd.Name + "\n\n" + md
	return os.WriteFile(path, []byte(md), 0o644)
}
