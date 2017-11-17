package logutil

import (
	"log"
	"os"
	"runtime"
	"time"

	xutil "go.repo/xutils/xapp"
	"go.repo/xutils/xfile"
)

const MaxFileSize = 10 * 1024 * 1024 // 10M
var gLogFile string
var gFile *os.File

func LogFilePath() string {
	return gLogFile
}

func Printf(format string, v ...interface{}) {
	doCheckFileSize()
	log.Printf(format, v...)
}

func Println(v ...interface{}) {
	doCheckFileSize()
	log.Println(v...)
}

func doCheckFileSize() {
	if gFile != nil {
		fi, err := gFile.Stat()
		if err == nil && fi.Size() > MaxFileSize {
			loginit()
		}
	}
}

func loginit() {
	gLogFile = time.Now().Format("20060102_150405.log")
	if file, ok := xfile.CreateFile(gLogFile); ok {
		gFile = file
		log.SetOutput(file)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("LogUtil: ")
}

func init() {
	if ptr, _, line, ok := runtime.Caller(1); ok && xutil.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	loginit()
}
