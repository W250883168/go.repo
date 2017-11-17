package xlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"

	xutil "go.repo/xutils/xapp"
	"go.repo/xutils/xfile"
)

var logFile string
var pFileInfo *os.FileInfo
var pLog *log.Logger

func GetLogger() *log.Logger {
	return pLog
}

func loginit() {
	logFlags := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	prefix := "XLog: "
	logFile = time.Now().Format("20060102_150405.log")

	if file, ok := xfile.CreateFile(logfile); ok {
		pLog = log.New(file, prefix, logFlags)
		// log.SetOutput(file)
		if fInfo, err := file.Stat(); err == nil {
			pFileInfo = &fInfo
		}
	}
}

func LogFilePath() string {
	return logFile
}

func LogFileSize() (size int64) {
	if pFileInfo != nil {
		size = pFileInfo.Size()
	}

	return size
}

func SetOutput(w io.Writer) {
	GetLogger().SetOutput(w)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(1); ok && xutil.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	loginit()
}
