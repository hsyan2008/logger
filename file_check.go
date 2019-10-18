// +build !windows

package logger

import "syscall"

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			Warn(err)
		}
	}()
	//防止多进程的并发操作
	if logObj != nil && logObj.logfile != nil {
		err := syscall.Flock(int(logObj.logfile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err != nil {
			return
		}
		defer syscall.Flock(int(logObj.logfile.Fd()), syscall.LOCK_UN)
	}
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
	}
}
