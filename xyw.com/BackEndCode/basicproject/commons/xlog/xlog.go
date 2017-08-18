package xlog

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"basicproject/commons/xfile"
)

var pLog *log.Logger

func GetLogger() *log.Logger {
	if pLog == nil {
		loginit()
	}

	return pLog
}

func loginit() {
	logFlags := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	prefix := "XLog: "
	logfile := time.Now().Format("20060102_150405.log")

	if file, ok := xfile.CreateFile(logfile); ok {
		pLog = log.New(file, prefix, logFlags)
		// log.SetOutput(file)
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(1); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
		// fmt.Println(fun.FileLine(ptr))
	}

	loginit()
}
