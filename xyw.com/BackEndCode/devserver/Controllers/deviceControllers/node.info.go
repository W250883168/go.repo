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

// 查询节点详细处理(关键字查询)
func Query_NodeDetail_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获异常

	// 响应数据
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 参数校验
	mapPara := request.Para.(map[string]interface{}) //获得通过断言实现类型转换
	keyword := mapPara["Keyword"].(string)

	// 校验权限
	//	TAG := "GetNode"
	//	dbmap := core.InitDb()
	//	defer dbmap.Db.Close()
	//	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
	//		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
	//		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
	//		return
	//	}

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	page := xhttp.PageInfo{PageIndex: request.Page.PageIndex, PageSize: request.Page.PageSize}
	list, err := zndxview.Query_NodeDetailView_ByKeyword(keyword, &page, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{Data: list}
}

// 查询节点详细处理(关键字查询)
func Query_NodeDetail_Info_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获异常

	// 响应数据
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	println(string(data))
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 参数校验
	mapPara := request.Para.(map[string]interface{}) //获得通过断言实现类型转换
	NodeId := mapPara["NodeId"].(string)

	// 校验权限
	//	TAG := "GetNode"
	//	dbmap := core.InitDb()
	//	defer dbmap.Db.Close()
	//	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
	//		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
	//		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
	//		return
	//	}

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	page := xhttp.PageInfo{PageIndex: request.Page.PageIndex, PageSize: request.Page.PageSize}
	list, err := zndxview.Query_NodeDetailView_ByID(NodeId, &page, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{Data: list}
}

// 获取节点下设备信息（详细）
func GetNode_DeviceDetailInfo_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request struct {
		UserID string
		NodeID string
	}
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 验证数据
	xtext.RequireNonBlank(request.NodeID)
	xtext.RequireNonBlank(request.UserID)

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	list, err := zndxview.QueryList_DeviceDetailView(request.NodeID, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &xhttp.HttpResponse{Content: &list}
}

// 获取节点基本信息
func Get_NodeBasicInfo_Handler(c *gin.Context) {
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
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request struct {
		Page    xhttp.PageInfo
		UserID  string
		Keyword string
	}
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
	list, err := zndxview.QueryList_NodeBasicView(request.Keyword, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &xhttp.HttpResponse{Content: &list}
}

// 获取节点或设备基本信息
func Get_NodeDevBasicInfo_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	// HTTP响应数据
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += "; ERR: " + err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request struct {
		Page    xhttp.PageInfo
		UserID  string
		Keyword string
	}
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
	list, err := zndxview.QueryList_NodeDevBasicView(request.Keyword, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &xhttp.HttpResponse{Content: &list}
}
