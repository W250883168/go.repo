package xdebug

import "fmt"
import "log"
import "runtime"
import "runtime/debug"
import "xutils/xapp"

// 输出错误所在文件/函数/行数信息
func PrintErrFunc(ptr uintptr, file string, line int, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Panicln(err)
		}
	}()

	f := runtime.FuncForPC(ptr)
	fmt.Printf("\n  FILE: %s; LINE: %d; PACKAGE.FUNC: %+v \n  ERROR: %s\n\n", file, line, f.Name(), err.Error())
}

// 输出错误消息并打印出堆栈信息
func PrintStackTrace(err error) {
	if err != nil {
		log.Println("	<<<<<<<<<<<< EEROR: " + err.Error())
		debug.PrintStack()
		fmt.Println("")
	}
}

// 输出错误消息并打印出堆栈信息
func DebugError(err error) {
	if err != nil {
		log.Println("	<<<<<<<<<<<< EEROR: " + err.Error())
		debug.PrintStack()
		fmt.Println("")
	}
}

// 输出错误消息并打印出堆栈信息
func LogError(err error) {
	if err != nil {
		txt := fmt.Sprintf("	<<<<<<<<<< LOG: %s", err.Error())
		log.Println(txt)
		log.Println(string(debug.Stack()))
	}
}

func LogString(msg string) {
	log.Println(msg)
	log.Println(string(debug.Stack()))
}

// 输出错误消息
func LogErrorText(err error) {
	if err != nil {
		log.Printf("	<<<<<<<<<<<<< ERROR:	%s", err.Error())
	}
}

// 抛出异常
func PanicError(err error) {
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
}

// 恢复异常
func RecoverError() {
	if err := recover(); err != nil {
		if e, ok := err.(error); ok {
			LogError(e)
		}
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok && xapp.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
