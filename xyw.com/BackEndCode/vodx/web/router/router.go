package router

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"

	"vodx/web/control"
)

func Load(pEngine *gin.Engine) {
	_VoidHandler := func(c *gin.Context) {}
	video := pEngine.Group("/video", _VoidHandler)
	vod := video.Group("/vod", _VoidHandler)
	vod.POST("/beginvideo", _AllowCrossDomain, control.BeginVideo)
	vod.POST("/endvideo", _AllowCrossDomain, control.EndVideo)

	vod.POST("/foox", func(c *gin.Context) {

	})
}

func _AllowCrossDomain(c *gin.Context) {
	//设置参数，允许跨域调用
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	c.Writer.Header().Set("Content-Type", "application/x-www-form-urlencoded")
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}
}
