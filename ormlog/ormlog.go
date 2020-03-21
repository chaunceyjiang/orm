package ormlog

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// 错误 红色
	errorLog = log.New(os.Stderr, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	// 警告 黄色
	warnLog = log.New(os.Stderr, "\033[33m[warning]\033[0m ", log.LstdFlags|log.Lshortfile)
	// info 绿色
	infoLog = log.New(os.Stderr, "\033[32m[info]\033[0m ", log.LstdFlags|log.Lshortfile)

	debugLog = log.New(os.Stderr, "[debug] ", log.LstdFlags|log.Lshortfile)

	loggers = []*log.Logger{debugLog, infoLog, warnLog, errorLog}
	mu      sync.Mutex
)

var (
	// Error Println
	Error = errorLog.Println
	// ErrorF Printf
	ErrorF   = errorLog.Printf
	Warning  = warnLog.Println
	WarningF = warnLog.Printf
	Info     = infoLog.Println
	InfoF    = infoLog.Printf
	Debug    = debugLog.Println
	DebugF   = debugLog.Printf
)

const (
	DebugLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	if level<0{
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if DebugLevel < level{
		// 小于设置的日志级别丢弃日志信息
		debugLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level{
		// 小于设置的日志级别丢弃日志信息
		infoLog.SetOutput(ioutil.Discard)
	}
	if WarningLevel < level{
		// 小于设置的日志级别丢弃日志信息
		warnLog.SetOutput(ioutil.Discard)
	}
	if ErrorLevel < level{
		// 小于设置的日志级别丢弃日志信息
		errorLog.SetOutput(ioutil.Discard)
	}

}
