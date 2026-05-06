// Package runner is a process runner.
package runner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Run executes a script via `bash -c`.
func Run(ctx context.Context, script string, env map[string]string) (stdout, stderr []byte, err error) {
	cmd := exec.Command("bash", "-c", script)
	cmd.Env = mergeEnv(env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("start bash: %w", err)
	}

	// Kill the process group on context cancellation.
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		case <-done:
		}
	}()

	waitErr := cmd.Wait()
	close(done)

	return outBuf.Bytes(), errBuf.Bytes(), waitErr
}

// mergeEnv overlays extra onto `os.Environ()`, with extra winning on key collisions.
func mergeEnv(extra map[string]string) []string {
	base := os.Environ()
	if len(extra) == 0 {
		return base
	}
	out := make([]string, 0, len(base)+len(extra))
	for _, kv := range base {
		if i := strings.IndexByte(kv, '='); i >= 0 {
			if _, override := extra[kv[:i]]; override {
				continue
			}
		}
		out = append(out, kv)
	}
	for k, v := range extra {
		out = append(out, k+"="+v)
	}
	return out
}
