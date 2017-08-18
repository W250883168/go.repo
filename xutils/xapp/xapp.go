package xapp

import (
	"fmt"
	"log"
	"runtime"
)

var gDebug bool = false

func IsDebugMode() bool {
	return gDebug
}

func SetRelaseMode() {
	gDebug = false
}

func SetDebugMode() {
	gDebug = true
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok && IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

}
