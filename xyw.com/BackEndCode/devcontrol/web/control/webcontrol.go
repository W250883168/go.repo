package webcontrol

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"dev.project/BackEndCode/devcontrol/web/action"
)

// 路由配置
func RouteConfig(r *gin.Engine) {
	r.GET("/", _AllowCrossDomain, func(c *gin.Context) {
		c.String(http.StatusOK, "CoapServer已启动")
	})

	// 设备开关(On)处理(通过节点控制)
	uri := "/device/node/control/switch/on/device"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload := "on", "SWITCH ON"
		action.Device_SwitchHandler(c, cmd, payload)
	})

	// 设备开关(Off)处理(通过节点控制)
	uri = "/device/node/control/switch/off/device"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload := "off", "SWITCH OFF"
		action.Device_SwitchHandler(c, cmd, payload)
	})

	// 切换器设备命令处理(视和讯/vnriver)
	uri = "/device/multiplexer/control/vnriver"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd := "toggle"
		action.Device_MultiplexerVnriver_Handler(cmd, c)
	})

	// 投影仪设备(Epson)命令处理(Power On)
	uri = "/device/projector/control/power/on/epson"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload, request_uri := "on", "PWR ON", "/rs232/d0"
		action.Device_ProjectorEpson_Handler(cmd, payload, request_uri, c)
	})

	// 投影仪设备(Epson)命令处理(Power Off)
	uri = "/device/projector/control/power/off/epson"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload, request_uri := "off", "PWR OFF", "/rs232/d0"
		action.Device_ProjectorEpson_Handler(cmd, payload, request_uri, c)
	})

	// 设备状态查询处理(按设备查询)
	uri = "/device/node/state/device"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		ptype := "device"
		action.Device_StateQuery_Handler(ptype, c)
	})

	// 设备状态查询处理(按位置查询)
	uri = "/device/node/state/room"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		ptype := "classroom"
		action.RoomDevice_StateQuery_Handler(ptype, c)
	})

	// 一键关闭房间设备
	uri = "/device/node/control/switch/off/room"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload := "off", "SWITCH OFF"
		action.RoomDevice_SwitchHandler(c, cmd, payload)
	})

	// 一键开启房间设备
	uri = "/device/node/control/switch/on/room"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload := "on", "SWITCH ON"
		action.RoomDevice_SwitchHandler(c, cmd, payload)
	})

	// 一键开启楼层设备
	uri = "/device/node/control/switch/on/floor"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload := "on", "SWITCH ON"
		action.FloorDevice_SwitchHandler(c, cmd, payload)
	})

	// 一键关闭楼层设备
	uri = "/device/node/control/switch/off/floor"
	r.POST(uri, _AllowCrossDomain, func(c *gin.Context) {
		cmd, payload := "off", "SWITCH OFF"
		action.FloorDevice_SwitchHandler(c, cmd, payload)
	})

	// OTA测试
	r.POST("/node_ota", _AllowCrossDomain, func(c *gin.Context) {
		// eui := "00124b000cd52302"
		// ver := "yyy"
		action.Node_OTAHandler2(c)
	})

	// OTA进度
	r.POST("/node_ota_progress", _AllowCrossDomain, func(c *gin.Context) {
		eui := "00124b000cd52302"
		action.Node_OTAProgress(eui, c)
	})
}

// 设置参数，允许跨域调用
func _AllowCrossDomain(c *gin.Context) {
	log.Printf("	###*****************URI=%s\n\n", c.Request.RequestURI)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	c.Writer.Header().Set("Content-Type", "application/x-www-form-urlencoded")
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
