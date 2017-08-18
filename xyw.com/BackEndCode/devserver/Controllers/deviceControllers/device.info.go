package deviceControllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xtext"

	"dborm/zndxview"

	"dev.project/BackEndCode/devserver/model/core"
	devmodel "dev.project/BackEndCode/devserver/model/deviceModel"
)

// 获取设备基本信息请求
type _Get_DeviceBasicInfo_Request struct {
	Page xhttp.PageInfo

	UserID  string
	Keyword string
}

// 查询设备基本信息列表
func Get_DeviceBasicInfo_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += "; ERR: " + err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request _Get_DeviceBasicInfo_Request
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 验证数据
	// xtext.RequireNonBlank(request.NodeID)
	xtext.RequireNonBlank(request.UserID)

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	list, err := zndxview.QueryList_DeviceBasicView(request.Keyword, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &xhttp.HttpResponse{Content: &list}
}

func Query_DeviceDetail_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += "; ERR: " + err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 获取请求数据
	para := request.Para.(map[string]interface{})
	keyword := para["Keyword"].(string)

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	list, err := zndxview.QueryList_DeviceDetailView0_ByKeyword(keyword, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{Data: list}
}

func QueryValid_DeviceDetail_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += "; ERR: " + err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 获取请求数据
	para := request.Para.(map[string]interface{})
	keyword := para["Keyword"].(string)

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	list, err := zndxview.QueryList_ValidDeviceDetailView0_ByKeyword(keyword, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{Data: list}
}

func Query_DeviceDetail_Info_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += "; ERR: " + err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 获取请求数据
	para := request.Para.(map[string]interface{})
	DeviceId := para["DeviceId"].(string)

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	list, err := zndxview.QueryList_DeviceDetailView0_ById(DeviceId, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{Data: list}
}
