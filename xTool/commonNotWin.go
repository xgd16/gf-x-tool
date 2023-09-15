//go:build !windows

package xTool

import "syscall"

// SafeExit 安全触发退出
func SafeExit() {
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
