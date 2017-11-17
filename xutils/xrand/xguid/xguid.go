package xguid

import (
	"fmt"
	"log"
	"runtime"

	"github.com/beevik/guid"

	"go.repo/xutils/xapp"
)

func foo() {
	log.Println(guid.NewString())
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok && xapp.IsDebugMode() {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

}
