package logger

import (
	"fmt"
	"log"
	"os"
	"path"
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
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

type _FILE struct {
	dir      string
	filename string
	_suffix  int
	isCover  bool
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

func SetConsole(isConsole bool) {
	consoleAppender = isConsole
	// log.SetFlags(log.Ldate | log.Lmicroseconds)
}

func SetPrefix(s string) {
	prefixStr = s
}

func SetLevel(_level LEVEL) {
	logLevel = _level
}

func SetRollingFile(fileName string, maxNumber int32, maxSize int64, _unit UNIT) {
	fileDir := path.Dir(fileName)
	fileName = path.Base(fileName)
	maxFileCount = maxNumber
	maxFileSize = maxSize * int64(_unit)
	RollingFile = true
	dailyRolling = false
	dirMk(fileDir)
	logObj = &_FILE{dir: fileDir, filename: fileName, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()
	for i := 1; i <= int(maxNumber); i++ {
		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
			logObj._suffix = i
		} else {
			break
		}
	}
	if !logObj.isMustRename() {
		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	} else {
		logObj.rename()
	}
	go fileMonitor()
}

func SetRollingDaily(fileName string) {
	fileDir := path.Dir(fileName)
	fileName = path.Base(fileName)
	RollingFile = false
	dailyRolling = true
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	dirMk(fileDir)
	logObj = &_FILE{dir: fileDir, filename: fileName, _date: &t, isCover: false, mu: new(sync.RWMutex)}
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	if !logObj.isMustRename() {
		logObj.logfile, _ = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		logObj.lg = log.New(logObj.logfile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	} else {
		logObj.rename()
	}
}

func console(levelAndPrefix string, s ...interface{}) {
	if consoleAppender {
		_, file, line, _ := runtime.Caller(3)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.SetFlags(log.Ldate | log.Lmicroseconds)
		log.Println(file, strconv.Itoa(line), levelAndPrefix+trim(fmt.Sprint(s)))
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
	}
}

func Output(levelInt LEVEL, level string, v ...interface{}) {
	if dailyRolling {
		fileCheck()
	}
	defer catchError()
	if logObj != nil {
		logObj.mu.RLock()
		defer logObj.mu.RUnlock()
	}

	if logLevel <= levelInt {
		var levelAndPrefix string
		if prefixStr == "" {
			levelAndPrefix = level + " "
		} else {
			levelAndPrefix = level + " " + prefixStr + " "
		}
		if logObj != nil {
			logObj.lg.Output(3, levelAndPrefix+trim(fmt.Sprint(v))+"\n")
		}
		console(levelAndPrefix, v)
	}
}

//interface会在两端加了[]，去掉
func trim(s string) string {
	return strings.Trim(s, "[]")
}

func Debug(v ...interface{}) {
	Output(DEBUG, "DEBUG", v)
}
func Info(v ...interface{}) {
	Output(INFO, "INFO", v)
}
func Warn(v ...interface{}) {
	Output(WARN, "WARN", v)
}
func Error(v ...interface{}) {
	Output(ERROR, "ERROR", v)
}
func Fatal(v ...interface{}) {
	Output(FATAL, "FATAL", v)
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
			if fileSize(f.dir+"/"+f.filename) >= maxFileSize {
				return true
			}
		}
	}
	return false
}

func (f *_FILE) rename() {
	if dailyRolling {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
			f.lg = log.New(logObj.logfile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
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
		f.logfile.Close()
	}
	if isExist(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix))) {
		os.Remove(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix)))
	}
	os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+"."+strconv.Itoa(int(f._suffix)))
	f.logfile, _ = os.Create(f.dir + "/" + f.filename)
	f.lg = log.New(logObj.logfile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
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

func dirMk(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		fmt.Println("创建目录失败：", err.Error())
	}
}
