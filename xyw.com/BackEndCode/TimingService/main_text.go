package main

import (
	app "TimingService/Action"
	"fmt"
	"log"
	"runtime"
	xconfig "xutils/xconfig"
	core "xutils/xcore"

	"github.com/gin-gonic/gin"
)

func main() {
	var cf xconfig.Config
	cf.InitConfig("./config.ini")
	core.CheckDB()
	r := gin.Default()
	app.LoadAction(r)
	app.LoadWebFile(r)
	r.Run(":" + cf.Read("server", "serverprot"))

}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}
