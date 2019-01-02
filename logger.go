package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	_VER string = "1.0.2"
)

type LEVEL int32

var logLevel LEVEL = 1
var maxFileSize int64
var maxFileCount int32
var dailyRolling bool = true
var consoleAppender bool = true
var RollingFile bool = false
var logObj *_FILE
var prefixStr = ""
var logGoID = false

var logConsole = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

const DATEFORMAT = "2006-01-02"

type UNIT int64

const (
	_       = iota
	KB UNIT = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type _FILE struct {
	filePath string //完整路径
	dir      string //目录
	filename string //文件名
	_suffix  int
	isCover  bool
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

func SetConsole(isConsole bool) {
	consoleAppender = isConsole
}

func SetPrefix(s string) {
	prefixStr = s
}

func GetPrefix() string {
	if logGoID {
		return strings.TrimSpace("GoID:" + GoroutineID() + " " + prefixStr)
	} else {
		return prefixStr
	}
}

func SetLogGoID(b bool) {
	logGoID = b
}

var l = len("goroutine ")

func GoroutineID() string {
	var buf [32]byte
	n := runtime.Stack(buf[:], false)

	b := bytes.NewBuffer(buf[l:n])
	s, _ := b.ReadString(' ')

	return strings.TrimSpace(s)
}

func SetLevel(_level LEVEL) {
	logLevel = _level
}
func Level() LEVEL {
	return logLevel
}
func SetLevelStr(level string) {
	logLevel = getLogLevel(level)
}

func SetRollingFile(fileName string, maxNumber int32, maxSize int64, unit string) {
	fileDir := filepath.Dir(fileName)
	fileName = filepath.Base(fileName)
	_unit := getLogUnit(unit)
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	RollingFile = true
	dailyRolling = false
	mkdir(fileDir)
	logObj = &_FILE{
		filePath: filepath.Join(fileDir, fileName),
		dir:      fileDir,
		filename: fileName,
		isCover:  false,
		mu:       new(sync.RWMutex),
	}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()
	for i := 1; i <= int(maxNumber); i++ {
		if isExist(logObj.filePath + "." + strconv.Itoa(i)) {
			logObj._suffix = i
		} else {
			break
		}
	}
	if !logObj.isMustRename() {
		logObj.initFile()
	} else {
		logObj.rename()
	}
	go fileMonitor()
}

func SetRollingDaily(fileName string) {
	fileDir := filepath.Dir(fileName)
	fileName = filepath.Base(fileName)
	RollingFile = false
	dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	mkdir(fileDir)
	logObj = &_FILE{
		filePath: filepath.Join(fileDir, fileName),
		dir:      fileDir,
		filename: fileName,
		_date:    &t,
		isCover:  false,
		mu:       new(sync.RWMutex),
	}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if !logObj.isMustRename() {
		logObj.initFile()
	} else {
		logObj.rename()
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func Output(calldepth int, level string, v ...interface{}) {
	if logLevel > getLogLevel(level) {
		return
	}

	if dailyRolling {
		fileCheck()
	}
	defer catchError()

	//print to file
	if logObj != nil {
		_ = logObj.lg.Output(calldepth, getLevelAndPrefix(level, false)+fmt.Sprintln(v...))
	}
	//print to console
	if consoleAppender {
		_ = logConsole.Output(calldepth, getLevelAndPrefix(level, true)+fmt.Sprintln(v...))
	}
}

func getLevelAndPrefix(level string, isConsole bool) (levelAndPrefix string) {
	if isConsole {
		level = getColor(level)
	}
	prefix := GetPrefix()
	if len(prefix) == 0 {
		levelAndPrefix = level + " "
	} else {
		levelAndPrefix = level + " " + prefix + " "
	}

	return levelAndPrefix
}

func getColor(level string) string {
	if runtime.GOOS == "windows" {
		return level
	}
	var color string
	switch level {
	case "DEBUG":
		color = "[32;1m"
	case "INFO":
		color = "[33;1m"
	case "WARN":
		color = "[35;1m"
	case "ERROR":
		color = "[31;1m"
	case "FATAL":
		color = "[31;7m"
	}
	return color + level + "[37;0m"
}

func Debug(v ...interface{}) {
	Output(3, "DEBUG", v...)
}
func Debugf(f string, v ...interface{}) {
	Output(3, "DEBUG", fmt.Sprintf(f, v...))
}
func Info(v ...interface{}) {
	Output(3, "INFO", v...)
}
func Infof(f string, v ...interface{}) {
	Output(3, "INFO", fmt.Sprintf(f, v...))
}
func Warn(v ...interface{}) {
	Output(3, "WARN", v...)
}
func Warnf(f string, v ...interface{}) {
	Output(3, "WARN", fmt.Sprintf(f, v...))
}
func Error(v ...interface{}) {
	Output(3, "ERROR", v...)
}
func Errorf(f string, v ...interface{}) {
	Output(3, "ERROR", fmt.Sprintf(f, v...))
}
func Fatal(v ...interface{}) {
	Output(3, "FATAL", v...)
}
func Fatalf(f string, v ...interface{}) {
	Output(3, "FATAL", fmt.Sprintf(f, v...))
}

var checkMustRenameTime int64

func (f *_FILE) isMustRename() bool {
	//3秒检查一次，不然太频繁
	if checkMustRenameTime != 0 && time.Now().Unix()-checkMustRenameTime < 3 {
		return false
	}
	checkMustRenameTime = time.Now().Unix()
	if dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	} else {
		if maxFileCount > 1 {
			if fileSize(f.filePath) >= maxFileSize {
				return true
			}
		}
	}
	return false
}

func (f *_FILE) rename() {
	if dailyRolling {
		fn := f.filePath + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				_ = f.logfile.Close()
			}
			err := os.Rename(f.filePath, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.initFile()
		}
	} else {
		f.coverNextOne()
	}
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(maxFileCount) + 1)
}

func (f *_FILE) coverNextOne() {
	f._suffix = f.nextSuffix()
	if f.logfile != nil {
		_ = f.logfile.Close()
	}
	if isExist(f.filePath + "." + strconv.Itoa(int(f._suffix))) {
		_ = os.Remove(f.filePath + "." + strconv.Itoa(int(f._suffix)))
	}
	_ = os.Rename(f.filePath, f.filePath+"."+strconv.Itoa(int(f._suffix)))
	f.initFile()
}

func (f *_FILE) initFile() {
	f.logfile, _ = os.OpenFile(f.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	f.lg = log.New(f.logfile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			fileCheck()
		}
	}
}

func fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer logObj.mu.Unlock()
		logObj.rename()
	}
}

func mkdir(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic("create dir faild：" + err.Error())
	}
}
func getLogLevel(l string) LEVEL {
	switch strings.ToUpper(l) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	case "OFF":
		return OFF
	default:
		return DEBUG
	}
}

func getLogUnit(u string) UNIT {
	switch strings.ToUpper(u) {
	case "K", "KB":
		return KB
	case "M", "MB":
		return MB
	case "G", "GB":
		return GB
	case "T", "TB":
		return TB
	default:
		return KB
	}
}
