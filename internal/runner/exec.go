// Package runner is a process runner.
package runner

import "syscall"

// Exec replaces the current process with binPath via syscall.Exec.
func Exec(binPath string, args []string, extras map[string]string) error {
	argv := append([]string{binPath}, args...)
	return syscall.Exec(binPath, argv, mergeEnv(extras))
}

// ExecBash replaces the current process with bash running script via syscall.Exec.
// argv0 sets bash's $0; args become $1.. inside the script.
func ExecBash(argv0, script string, args []string, extras map[string]string) error {
	argv := []string{"bash", "-c", script, argv0}
	argv = append(argv, args...)
	return syscall.Exec("/bin/bash", argv, mergeEnv(extras))
}
