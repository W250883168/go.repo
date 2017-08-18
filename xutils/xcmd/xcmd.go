package xcmd

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// 执行cmd命令
func Exec(cmd string) error {
	command := exec.Command("cmd.exe", "/c ", cmd)
	return command.Run()
}

// 执行cmd命令(异步方式)
func AsyncExec(cmd string) {
	go Exec(cmd)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
