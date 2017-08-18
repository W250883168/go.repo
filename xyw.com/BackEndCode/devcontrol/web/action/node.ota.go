package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"canopus"

	"xutils/xerr"

	"dev.project/BackEndCode/devcontrol/coap"
	"dev.project/BackEndCode/devcontrol/coap/coapclient"
)

// 节点OTA请求处理
func Node_OTAHandler(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var code, msg, data string = "0", "执行成功", "ok"
	defer func() {
		if err != nil {
			code, msg, data = "1001", err.Error(), "fail"
		}
		c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg, "data": data})
	}()

	// 解析参数
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))
	var request = struct{ Eui, Version string }{}
	err = json.Unmarshal(body, &request)
	log.Printf("%+v", request)
	xerr.ThrowPanic(err)

	// 重置进度
	coap.GetOTAContext().PutValue(request.Eui, "0.0")

	var payload_struct = struct{ version string }{version: request.Version}
	payload_buff, _ := json.Marshal(&payload_struct)
	log.Println(string(payload_buff))

	// 发送命令
	addr := "192.168.0.177:20005"
	payload := "START OTA"
	coapcmd := coapclient.CoapCommand{
		HostAddr:    addr,
		Method:      canopus.Post,
		RequestURI:  "/node_ota",
		QueryParams: map[string]string{"eui": request.Eui},
		Payload:     payload}
	reply, _, err := coapclient.Send2(coapcmd)
	if reply != nil {
		log.Println("Payload: ", reply.GetMessage().Payload.String())
	}
}

// 节点OTA请求处理
func Node_OTAHandler2(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var code, msg, data string = "0", "执行成功", "ok"
	defer func() {
		if err != nil {
			code, msg, data = "1001", err.Error(), "fail"
		}
		c.JSON(http.StatusOK, gin.H{"code": code, "msg": msg, "data": data})
	}()

	// 解析参数
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Println(string(body))
	var request = struct{ Eui, Version string }{}
	err = json.Unmarshal(body, &request)
	log.Printf("%+v", request)
	xerr.ThrowPanic(err)

	// 重置进度
	coap.GetOTAContext().PutValue(request.Eui, "0.0")

	var payload_struct = struct{ version string }{version: request.Version}
	payload_buff, _ := json.Marshal(&payload_struct)
	log.Println(string(payload_buff))

	// 发送命令
	addr := "192.168.0.22:20005"
	payload := "START OTA;" + request.Version
	coapcmd := coapclient.CoapCommand{
		HostAddr:    addr,
		Method:      canopus.Post,
		RequestURI:  "/node_ota",
		QueryParams: map[string]string{"eui": request.Eui},
		Payload:     payload}
	reply, _, err := coapclient.Send2(coapcmd)
	if reply != nil {
		log.Println("Payload: ", reply.GetMessage().Payload.String())
	}
}

// 节点OTA进度查询
func Node_OTAProgress(eui string, c *gin.Context) {
	defer xerr.CatchPanic()

	progress := coap.GetOTAContext().GetValue(eui)
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "执行成功", "data": progress})
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
