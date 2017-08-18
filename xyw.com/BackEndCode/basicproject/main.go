package main

import (
	_ "net/http/pprof"

	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"

	"xutils/xerr"

	con "xutils/xconfig"
	core "xutils/xcore"

	app "basicproject/action"
)

func main() {
	var cf con.Config
	cf.InitConfig("./config.ini")

	core.CheckDB()

	//性能监测，通过“http://host:port/debug/pprof/”访问
	go func() {
		defer xerr.CatchPanic()

		addr := fmt.Sprintf(":%s", cf.Read("server", "pprof.port"))
		log.Printf("\n	性能监测:  http://host%s/debug/pprof	\n", addr)
		log.Println(http.ListenAndServe(addr, nil))
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	app.LoadAction(r)
	app.LoadSystemAction(r)
	app.LoadWebModel(r)
	//	app.LoadDeviceAction(r)
	r.Run(":" + cf.Read("server", "serverprot"))
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if pFile := core.GetLogFile(); pFile != nil {
		log.Printf("<<<<<<\t	LOGFILE:  %s\n", pFile.Name())
		os.Stdout.WriteString("test.....................\n")
		log.SetOutput(os.Stdout)
	}
}
