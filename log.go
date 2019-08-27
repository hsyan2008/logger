package logger

import (
	"fmt"
)

type Log struct {
	hasPrefix bool
	prefixStr string
	traceID   string
}

func NewLog() *Log {
	return &Log{
		hasPrefix: true,
		prefixStr: GetPrefix(),
	}
}

func (this *Log) AppendPrefix(str string) {
	if this.hasPrefix == false {
		this.ResetPrefix()
	}
	this.prefixStr = fmt.Sprintf("%s %s", this.prefixStr, str)
	this.hasPrefix = true
}

func (this *Log) SetPrefix(str string) {
	this.prefixStr = str
	this.hasPrefix = true
}

func (this *Log) ResetPrefix() {
	this.prefixStr = GetPrefix()
	this.hasPrefix = true
}

func (this *Log) SetTraceID(str string) {
	this.traceID = fmt.Sprintf("traceid: %s", str)
}

func (this *Log) getPrefix() string {
	if this.hasPrefix == false {
		this.ResetPrefix()
	}

	if this.traceID == "" {
		return this.prefixStr
	}

	return this.traceID + " " + this.prefixStr
}

func (this *Log) Debug(v ...interface{}) {
	Output(3, "DEBUG", this.getPrefix(), v...)
}

func (this *Log) Debugf(format string, v ...interface{}) {
	Output(3, "DEBUG", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Log) Info(v ...interface{}) {
	Output(3, "INFO", this.getPrefix(), v...)
}

func (this *Log) Infof(format string, v ...interface{}) {
	Output(3, "INFO", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Log) Warn(v ...interface{}) {
	Output(3, "WARN", this.getPrefix(), v...)
}

func (this *Log) Warnf(format string, v ...interface{}) {
	Output(3, "WARN", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Log) Error(v ...interface{}) {
	Output(3, "ERROR", this.getPrefix(), v...)
}

func (this *Log) Errorf(format string, v ...interface{}) {
	Output(3, "ERROR", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Log) Fatal(v ...interface{}) {
	Output(3, "FATAL", this.getPrefix(), v...)
}

func (this *Log) Fatalf(format string, v ...interface{}) {
	Output(3, "FATAL", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Log) Mix(v ...interface{}) {
	Output(3, "MIX", this.getPrefix(), v...)
}

func (this *Log) Mixf(format string, v ...interface{}) {
	Output(3, "MIX", this.getPrefix(), fmt.Sprintf(format, v...))
}

func (this *Log) Output(calldepth int, s string) error {
	Output(2+calldepth, "MIX", this.getPrefix(), s)
	return nil
}
