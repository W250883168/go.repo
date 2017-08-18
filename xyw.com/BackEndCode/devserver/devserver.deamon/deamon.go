package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"xutils/xerr"
)

func main() {
	fmt.Println("Hello,,,,")
	fmt.Println(os.Args)
	fmt.Println(os.Getwd())
	os.Chdir("..")
	fmt.Println(os.Getwd())

	cmd := "devserver.exe"
	for {
		do(cmd)
	}

}

func do(cmd string) {
	defer xerr.CatchPanic()

	log.Println("服务启动时间： ", time.Now())
	pCommand := exec.Command(cmd)
	pCommand.Start()
	pCommand.Wait()
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}
