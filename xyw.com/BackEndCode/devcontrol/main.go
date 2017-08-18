package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"dev.project/BackEndCode/devcontrol/app"
	"dev.project/BackEndCode/devcontrol/coap/coapserver"
	"dev.project/BackEndCode/devcontrol/ioutil/dbutil"
	webcontrol "dev.project/BackEndCode/devcontrol/web/control"
)

func main() {
	defer _DumpErr()
	_AppInit()

	// 启动CoAP服务
	go coapserver.StartCoapService()

	gin.SetMode(gin.ReleaseMode)
	pEngine := gin.Default()
	webcontrol.RouteConfig(pEngine)
	addr := fmt.Sprintf(":%d", app.GetConfig().GinHTTPPort)
	pEngine.Run(addr)
}

func _AppInit() {
	// 检查数据库
	dbutil.CheckDB()

	//性能监测: http://host:port/debug/pprof
	go func() {
		addr := fmt.Sprintf(":%d", app.GetConfig().PprofHTTPPort)
		log.Println(http.ListenAndServe(addr, nil))
	}()
}

func _DumpErr() {
	if p := recover(); p != nil {
		log.Println("ERR: ", p)
		log.Println(string(debug.Stack()))
		panic(p)
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	if pLogFile := app.GetLogFile(); pLogFile != nil {
		log.Printf("<<<<<<<\t 日志文件:  %s\n", pLogFile.Name())
		// os.Stdout = pLogFile
		// log.SetOutput(pLogFile)
	}

	buff, _ := json.Marshal(app.GetConfig())
	log.Printf("<<<<<<<<<< APP配置: \n\t%s\n", string(buff))
	log.Printf(`性能监测:  http://%s:%d/debug/pprof`, app.GetConfig().ThisHostAddr, app.GetConfig().PprofHTTPPort)
}
