// +build windows

package logger

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			Warn(err)
		}
	}()
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
	}
}
