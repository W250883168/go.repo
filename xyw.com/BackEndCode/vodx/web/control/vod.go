package control

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"

	"xutils/xerr"

	"vodx/web/action"
	"vodx/web/view/httpview"
)

func BeginVideo(c *gin.Context) {
	defer xerr.CatchPanic()

	// 数据响应
	var err error
	var response = struct{ Content interface{} }{}
	defer func() {
		if err != nil {
			response.Content = err.Error()
		}

		c.JSON(http.StatusOK, response)
	}()

	// 解析数据
	buff, _ := ioutil.ReadAll(c.Request.Body)
	var request httpview.VideoCaptureRequest
	err = json.Unmarshal(buff, &request)
	xerr.ThrowPanic(err)

	// 执行动作
	log.Printf("%+v\n", request)
	action.BeginVideo(&request)
	response.Content = "ok"
}

func EndVideo(c *gin.Context) {
	defer xerr.CatchPanic()

	// 数据响应
	var err error
	var response = struct{ Content interface{} }{}
	defer func() {
		if err != nil {
			response.Content = err.Error()
		}

		c.JSON(http.StatusOK, response)
	}()

	// 解析数据
	buff, _ := ioutil.ReadAll(c.Request.Body)
	var request httpview.VideoCaptureRequest
	err = json.Unmarshal(buff, &request)
	xerr.ThrowPanic(err)

	// 执行动作
	log.Printf("%+v\n", request)
	action.EndVideo(&request)
	response.Content = "ok"
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}
}
