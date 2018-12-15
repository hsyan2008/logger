// +build !windows

package logger

import "syscall"

func init() {
	_ = syscall.Umask(011)
}
