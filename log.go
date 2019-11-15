package logger

import (
	"fmt"
)

type Logger struct {
	hasPrefix bool
	prefixStr string
	traceID   string
}

func NewLogger() *Logger {
	return &Logger{
		hasPrefix: true,
		prefixStr: GetPrefix(),
	}
}

func (this *Logger) AppendPrefix(str string) {
	if str == "" {
		return
	}

	if this.hasPrefix == false {
		this.ResetPrefix()
	}

	if this.prefixStr == "" {
		this.prefixStr = str
	} else {
		this.prefixStr = this.prefixStr + " " + str
	}
	this.hasPrefix = true
}

func (this *Logger) SetPrefix(str string) {
	this.prefixStr = str
	this.hasPrefix = true
}

func (this *Logger) ResetPrefix() {
	this.prefixStr = GetPrefix()
	this.hasPrefix = true
}

func (this *Logger) SetTraceID(str string) {
	this.traceID = str
}

func (this *Logger) GetTraceID() string {
	return this.traceID
}

func (this *Logger) getPrefix() string {
	if this.hasPrefix == false {
		this.ResetPrefix()
	}

	if this.traceID == "" {
		return this.prefixStr
	}

	if this.prefixStr == "" {
		return "trace_id:" + this.traceID
	}

	return "trace_id:" + this.traceID + " " + this.prefixStr
}

func (this *Logger) Debug(v ...interface{}) {
	if logLevel == OFF || logLevel > DEBUG {
		return
	}
	Output(3, "DEBUG", this.getPrefix(), v...)
}

func (this *Logger) Debugf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > DEBUG {
		return
	}
	Output(3, "DEBUG", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Info(v ...interface{}) {
	if logLevel == OFF || logLevel > INFO {
		return
	}
	Output(3, "INFO", this.getPrefix(), v...)
}

func (this *Logger) Infof(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > INFO {
		return
	}
	Output(3, "INFO", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Warn(v ...interface{}) {
	if logLevel == OFF || logLevel > WARN {
		return
	}
	Output(3, "WARN", this.getPrefix(), v...)
}

func (this *Logger) Warnf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > WARN {
		return
	}
	Output(3, "WARN", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Error(v ...interface{}) {
	if logLevel == OFF || logLevel > ERROR {
		return
	}
	Output(3, "ERROR", this.getPrefix(), v...)
}

func (this *Logger) Errorf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > ERROR {
		return
	}
	Output(3, "ERROR", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Fatal(v ...interface{}) {
	if logLevel == OFF || logLevel > FATAL {
		return
	}
	Output(3, "FATAL", this.getPrefix(), v...)
}

func (this *Logger) Fatalf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > FATAL {
		return
	}
	Output(3, "FATAL", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Mix(v ...interface{}) {
	if logLevel == OFF || logLevel > MIX {
		return
	}
	Output(3, "MIX", this.getPrefix(), v...)
}

func (this *Logger) Mixf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > MIX {
		return
	}
	Output(3, "MIX", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Output(calldepth int, s string) error {
	Output(2+calldepth, "MIX", this.getPrefix(), s)
	return nil
}
