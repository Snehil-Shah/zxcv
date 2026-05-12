// Package main is the entry point for `zxcv`.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Snehil-Shah/zxcv/internal/cli"
	"github.com/Snehil-Shah/zxcv/internal/ui"
)

// version is set at build time via goreleaser ldflags.
var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // graceful shutdowns
	defer stop()
	if err := cli.New(version).Run(ctx, os.Args); err != nil {
		ui.Error("zxcv: " + err.Error())
		os.Exit(1)
	}
}
