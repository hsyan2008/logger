package logger

import (
	"fmt"
)

type Log struct {
	hasPrefix bool
	prefixStr string
}

func NewLog() *Log {
	return &Log{
		hasPrefix: true,
		prefixStr: GetPrefix(),
	}
}

func (this *Log) AppendPrefix(str string) {
	this.prefixStr = fmt.Sprintf("%s %s", this.getPrefix(), str)
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

func (this *Log) getPrefix() string {
	if this.hasPrefix == false {
		this.prefixStr = GetPrefix()
		this.hasPrefix = true
	}

	return this.prefixStr
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

func (this *Log) Output(calldepth int, s string) error {
	Output(2+calldepth, "MIX", this.getPrefix(), s)
	return nil
}
