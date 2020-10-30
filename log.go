package logger

import (
	"fmt"
)

type Logger struct {
	hasPrefix bool
	prefixStr string
	traceID   string
	calldepth int
}

func NewLogger() *Logger {
	return &Logger{
		hasPrefix: true,
		prefixStr: GetPrefix(),
		calldepth: 2,
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

func (this *Logger) SetCallDepth(calldepth int) {
	this.calldepth = calldepth
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

func (this *Logger) GetPrefix() string {
	return this.prefixStr
}

func (this *Logger) getFullPrefix() string {
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
	Output(this.calldepth+1, "DEBUG", this.getFullPrefix(), v...)
}

func (this *Logger) Debugf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > DEBUG {
		return
	}
	Output(this.calldepth+1, "DEBUG", this.getFullPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Info(v ...interface{}) {
	if logLevel == OFF || logLevel > INFO {
		return
	}
	Output(this.calldepth+1, "INFO", this.getFullPrefix(), v...)
}

func (this *Logger) Infof(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > INFO {
		return
	}
	Output(this.calldepth+1, "INFO", this.getFullPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Warn(v ...interface{}) {
	if logLevel == OFF || logLevel > WARN {
		return
	}
	Output(this.calldepth+1, "WARN", this.getFullPrefix(), v...)
}

func (this *Logger) Warnf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > WARN {
		return
	}
	Output(this.calldepth+1, "WARN", this.getFullPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Error(v ...interface{}) {
	if logLevel == OFF || logLevel > ERROR {
		return
	}
	Output(this.calldepth+1, "ERROR", this.getFullPrefix(), v...)
}

func (this *Logger) Errorf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > ERROR {
		return
	}
	Output(this.calldepth+1, "ERROR", this.getFullPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Fatal(v ...interface{}) {
	if logLevel == OFF || logLevel > FATAL {
		return
	}
	Output(this.calldepth+1, "FATAL", this.getFullPrefix(), v...)
}

func (this *Logger) Fatalf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > FATAL {
		return
	}
	Output(this.calldepth+1, "FATAL", this.getFullPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Mix(v ...interface{}) {
	if logLevel == OFF || logLevel > MIX {
		return
	}
	Output(this.calldepth+1, "MIX", this.getFullPrefix(), v...)
}

func (this *Logger) Mixf(format string, v ...interface{}) {
	if logLevel == OFF || logLevel > MIX {
		return
	}
	Output(this.calldepth+1, "MIX", this.getFullPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Output(calldepth int, s string) error {
	Output(this.calldepth+calldepth, "MIX", this.getFullPrefix(), s)
	return nil
}
