// Package runner_test tests the `runner` package.
package runner_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/runner"
)

func TestRun_Stdout(t *testing.T) {
	stdout, _, err := runner.Run(context.Background(), "echo hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(string(stdout)), "hello"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRun_Stderr(t *testing.T) {
	_, stderr, err := runner.Run(context.Background(), "echo oops >&2", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(string(stderr)), "oops"; got != want {
		t.Errorf("stderr = %q, want %q", got, want)
	}
}

func TestRun_ExitNonZero(t *testing.T) {
	_, _, err := runner.Run(context.Background(), "exit 7", nil)
	if err == nil {
		t.Fatal("expected error on non-zero exit, got nil")
	}
}

func TestRun_EnvOverride(t *testing.T) {
	stdout, _, err := runner.Run(context.Background(), `echo "$TEST"`, map[string]string{"TEST": "hi"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(string(stdout)), "hi"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}
