package actiondataControllers

import (
	"github.com/gin-gonic/gin"
)

func _AllowCrossDomain(c *gin.Context) {
	//设置参数，允许跨域调用
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	c.Writer.Header().Set("Content-Type", "application/x-www-form-urlencoded")
}
