package xdebug

import "runtime"
import "runtime/debug"
import "fmt"
import "log"

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

func PrintStackTrace(err error) {
	if err != nil {
		fmt.Println("	<<<<<<<<<<<< EEROR: " + err.Error())
		debug.PrintStack()
		fmt.Println("")
	}
}

func HandlePanic() {
	go func() {
		if err := recover(); err != nil {
			log.Panicln(err)
		}
	}()
}

func DebugError(err error) {
	if err != nil {
		fmt.Println("	<<<<<<<<<<<< EEROR: " + err.Error())
		debug.PrintStack()
		fmt.Println("")
	}
}

func LogError(err error) {
	if err != nil {
		txt := fmt.Sprintf("	<<<<<<<<<< LOG: %s", err.Error())
		log.Println(txt)
		log.Println(string(debug.Stack()))
	}
}

// 抛出异常
func HandleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

// 恢复异常
func DoRecover() {
	if err := recover(); err != nil {
		if err, ok := err.(error); ok {
			DebugError(err)
		}
	}
}

// 捕获异常
func CatchError() (err error) {
	if catch := recover(); catch != nil {
		if err2, ok := err.(error); ok {
			fmt.Printf("CatchERR: %s\n", err.Error())
			err = err2
		}
	}

	return err
}
