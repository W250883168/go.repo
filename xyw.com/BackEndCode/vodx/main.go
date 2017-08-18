package main

import (
	"vodx/app"

	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"

	"vodx/web/router"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	pEngine := gin.Default()
	router.Load(pEngine)
	pEngine.Run(fmt.Sprintf(":%d", app.GetConfig().HttpPort))
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}
}
