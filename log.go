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
	if this == nil {
		*this = Logger{}
	}
	if this.hasPrefix == false {
		this.ResetPrefix()
	}
	this.prefixStr = fmt.Sprintf("%s %s", this.prefixStr, str)
	this.hasPrefix = true
}

func (this *Logger) SetPrefix(str string) {
	if this == nil {
		*this = Logger{}
	}
	this.prefixStr = str
	this.hasPrefix = true
}

func (this *Logger) ResetPrefix() {
	if this == nil {
		*this = Logger{}
	}
	this.prefixStr = GetPrefix()
	this.hasPrefix = true
}

func (this *Logger) SetTraceID(str string) {
	if this == nil {
		*this = Logger{}
	}
	this.traceID = fmt.Sprintf("traceid: %s", str)
}

func (this *Logger) getPrefix() string {
	if this == nil {
		*this = Logger{}
	}
	if this.hasPrefix == false {
		this.ResetPrefix()
	}

	if this.traceID == "" {
		return this.prefixStr
	}

	return this.traceID + " " + this.prefixStr
}

func (this *Logger) Debug(v ...interface{}) {
	Output(3, "DEBUG", this.getPrefix(), v...)
}

func (this *Logger) Debugf(format string, v ...interface{}) {
	Output(3, "DEBUG", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Info(v ...interface{}) {
	Output(3, "INFO", this.getPrefix(), v...)
}

func (this *Logger) Infof(format string, v ...interface{}) {
	Output(3, "INFO", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Warn(v ...interface{}) {
	Output(3, "WARN", this.getPrefix(), v...)
}

func (this *Logger) Warnf(format string, v ...interface{}) {
	Output(3, "WARN", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Error(v ...interface{}) {
	Output(3, "ERROR", this.getPrefix(), v...)
}

func (this *Logger) Errorf(format string, v ...interface{}) {
	Output(3, "ERROR", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Fatal(v ...interface{}) {
	Output(3, "FATAL", this.getPrefix(), v...)
}

func (this *Logger) Fatalf(format string, v ...interface{}) {
	Output(3, "FATAL", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Mix(v ...interface{}) {
	Output(3, "MIX", this.getPrefix(), v...)
}

func (this *Logger) Mixf(format string, v ...interface{}) {
	Output(3, "MIX", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Logger) Output(calldepth int, s string) error {
	Output(2+calldepth, "MIX", this.getPrefix(), s)
	return nil
}
