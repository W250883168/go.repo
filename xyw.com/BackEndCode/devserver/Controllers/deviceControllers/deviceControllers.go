package deviceControllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xnumeric"
	"xutils/xtext"
	"xutils/xtime"

	"dborm/zndx"
	"dborm/zndxview"

	devdao "dev.project/BackEndCode/devserver/DataAccess/deviceDataAccess"
	userdao "dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	core "dev.project/BackEndCode/devserver/model/core"
	devmodel "dev.project/BackEndCode/devserver/model/deviceModel"

	"dev.project/BackEndCode/devserver/commons"
)

// 获取设备使用日志
func GetDeviceUseLogList(c *gin.Context) {
	defer xerr.CatchPanic()

	responses := commons.ResponseMsgSet_Instance()
	var rd = core.Returndata{Rcode: responses.DATA_MALFORMED.CodeText(), Reason: responses.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	// 解析数据
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	// 获得查询参数
	mapPara := requestData.Para.(map[string]interface{}) // 获得通过断言实现类型转换
	obj_DeviceId := mapPara["DeviceId"]
	deviceId := obj_DeviceId.(string) // 类型转换，转换失败引发panic

	// 校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "GetDeviceUseLogList"
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = responses.AUTH_LIMITED.CodeText()
		rd.Reason = responses.AUTH_LIMITED.Text
		return
	}

	//查询数据
	// rd = devdao.GetDeviceUseLogList(requestData, deviceId, dbmap)
	rd.Rcode = responses.FAIL.CodeText()
	rd.Reason = responses.FAIL.Text
	devlog := zndx.DeviceUseLog{DeviceId: deviceId}
	pageinfo := xhttp.PageInfo{PageIndex: requestData.Page.PageIndex, PageSize: requestData.Page.PageSize}
	list, err := devlog.GetByDeviceID(&pageinfo, dbmap)
	xerr.ThrowPanic(err)

	// 响应数据
	rd.Rcode = responses.SUCCESS.CodeText()
	rd.Reason = responses.SUCCESS.Text
	rd.Result = &devmodel.ResultData{devmodel.PageData{
		PageIndex:   pageinfo.PageIndex,
		PageSize:    pageinfo.PageSize,
		RecordCount: pageinfo.RowTotal,
		PageCount:   pageinfo.PageTotal(),
	}, list}
}

//获取设备操作日志
func GetDeviceOperateLogList(c *gin.Context) {
	defer xerr.CatchPanic()

	responses := commons.ResponseMsgSet_Instance()
	var rd = core.Returndata{Rcode: responses.DATA_MALFORMED.CodeText(), Reason: responses.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//校验权限
	TAG := "GetDeviceOperateLogList"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = responses.AUTH_LIMITED.CodeText()
		rd.Reason = responses.AUTH_LIMITED.Text
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换
	//参数：DeviceId
	tDeviceId := mapPara["DeviceId"]
	deviceId := tDeviceId.(string)
	xtext.RequireNonBlank(deviceId)

	//查询数据
	rd = devdao.GetDeviceOperateLogList2(requestData, deviceId, dbmap)
}

//获取设备预警信息
func GetDeviceAlertInfoList(c *gin.Context) {
	//设置参数，允许跨域调用

	var rd core.Returndata
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		log.Println(rd.Reason + ":" + err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	TAG := "GetDeviceAlertInfoList"
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tDeviceId := mapPara["DeviceId"]
	if tDeviceId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		log.Println(TAG, rd.Reason)
		return
	}
	deviceId, ok := tDeviceId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceAlertInfoList(requestData, deviceId, dbmap)
}

//获取设备故障信息
func GetDeviceFaultInfoList(c *gin.Context) {
	//设置参数，允许跨域调用

	TAG := "GetDeviceFaultInfoList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tDeviceId := mapPara["DeviceId"]
	if tDeviceId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	deviceId, ok := tDeviceId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceFaultInfoList(requestData, deviceId, dbmap)

	c.JSON(200, rd)
}

//获取教室状态信息
func GetClassroomStatusList(c *gin.Context) {
	log.Printf("###-----------------%s\n", c.Request.RequestURI)

	TAG := "GetClassroomStatusList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：BuildingId
	tId := mapPara["BuildingIds"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据GetClassroomStatusData
	rd = devdao.GetClassroomStatusList(requestData, id, dbmap)

	//--------------------------------------------------------

	c.JSON(200, rd)

}

//获取所有设备操作日志
func GetAllOperateLogList(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	//获得查询参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验权限
	TAG := "GetAllOperateLogList"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, request.Auth.Rolestype, request.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//获得参数map
	mapPara := request.Para.(map[string]interface{}) //获得通过断言实现类型转换
	fromTime := mapPara["FromTime"].(string)         //参数：FromTime
	toTime := mapPara["ToTime"].(string)             //参数：ToTime
	keyWord := mapPara["KeyWord"].(string)           //参数：KeyWord

	//查询数据
	rd = devdao.GetAllOperateLogList(request, fromTime, toTime, keyWord, dbmap)
}

//获取所有设备预警消息
func GetAllAlertInfoList(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//校验权限
	TAG := "GetAllAlertInfoList"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//获得参数
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换
	siteType := mapPara["SiteType"].(string)             // 参数：SiteType
	siteId := mapPara["SiteId"].(string)                 // 参数：SiteId
	modelId := mapPara["ModelId"].(string)               // 参数：ModelId
	keyWord := mapPara["KeyWord"].(string)               // 参数：KeyWord

	//查询数据
	rd = devdao.GetAllAlertInfoList2(requestData, siteType, siteId, modelId, keyWord, dbmap)
}

//获取所有设备故障信息
func GetAllFaultInfoList(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd core.Returndata
	rd.Rcode = gResponseMsgs.DATA_MALFORMED.CodeText()
	rd.Reason = gResponseMsgs.DATA_MALFORMED.Text + "fdssssssssssssss"
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换
	//参数：SiteType
	tSiteType := mapPara["SiteType"]
	siteType := tSiteType.(string)
	//参数：SiteId
	tSiteId := mapPara["SiteId"]
	siteId := tSiteId.(string)
	//参数：ModelId
	tModelId := mapPara["ModelId"]
	modelId := tModelId.(string)
	//参数：KeyWord
	tKeyWord := mapPara["KeyWord"]
	keyWord := tKeyWord.(string)

	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "GetAllFaultInfoList"
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//查询数据
	rd = devdao.GetAllFaultInfoList2(requestData, siteType, siteId, modelId, keyWord, dbmap)
}

//获取设备数量
func GetDeviceQty(c *gin.Context) {
	//设置参数，允许跨域调用

	TAG := "GetDeviceQty"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：SiteType
	tSiteType := mapPara["SiteType"]
	if tSiteType == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	siteType, ok := tSiteType.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//参数：SiteId
	tSiteId := mapPara["SiteId"]
	if tSiteId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	siteId, ok := tSiteId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//参数：ModelId
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	modelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceQty(requestData, siteType, siteId, modelId, dbmap)

	c.JSON(200, rd)
}

//设备分析-按设备型号统计使用时间
func GetUseTimeByModel(c *gin.Context) {
	defer xerr.CatchPanic()

	// HTTP响应
	responses := commons.ResponseMsgSet_Instance()
	var rd core.Returndata
	rd.Rcode = responses.DATA_MALFORMED.CodeText()
	rd.Reason = responses.DATA_MALFORMED.Text
	defer func() { c.JSON(http.StatusOK, rd) }()

	// 解析网络数据
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	// 获得查询参数
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换
	//参数：FromTime
	tFromTime := mapPara["FromTime"]
	fromTime := tFromTime.(string) // 类型转换，类型错误则引发panic， 下同
	xtext.RequireNonBlank(fromTime)
	//参数：ToTime
	tToTime := mapPara["ToTime"]
	toTime := tToTime.(string)
	xtext.RequireNonBlank(toTime)
	//参数：SiteType
	tSiteType := mapPara["SiteType"]
	siteType := tSiteType.(string)
	// xtext.RequireNonBlank(siteType)
	//参数：SiteId
	tSiteId := mapPara["SiteId"]
	siteId := tSiteId.(string)
	// xtext.RequireNonBlank(siteId)
	//参数：ModelId
	tModelId := mapPara["ModelId"]
	modelId := tModelId.(string)
	// xtext.RequireNonBlank(modelId)

	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "GetUseTimeByModel"
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = responses.AUTH_LIMITED.CodeText()
		rd.Reason = responses.AUTH_LIMITED.Text
		return
	}

	//查询数据
	rd = devdao.GetUseTimeByModel2(requestData, fromTime, toTime, siteType, siteId, modelId, dbmap)
}

//设备分析-按设备位置统计使用时间
func GetUseTimeBySite(c *gin.Context) {
	//设置参数，允许跨域调用

	TAG := "GetUseTimeBySite"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：FromTime
	tFromTime := mapPara["FromTime"]
	if tFromTime == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	fromTime, ok := tFromTime.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if fromTime == "" {
		rd.Rcode = "1002"
		rd.Reason = "参数FromTo不能为空，请传入值！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//参数：ToTime
	tToTime := mapPara["ToTime"]
	if tFromTime == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	toTime, ok := tToTime.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if toTime == "" {
		rd.Rcode = "1002"
		rd.Reason = "参数ToTime不能为空，请传入值！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//参数：SiteType
	tSiteType := mapPara["SiteType"]
	if tSiteType == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	siteType, ok := tSiteType.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//参数：SiteId
	tSiteId := mapPara["SiteId"]
	if tSiteId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	siteId, ok := tSiteId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//参数：ModelId
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	modelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetUseTimeBySite(requestData, fromTime, toTime, siteType, siteId, modelId, dbmap)

	c.JSON(200, rd)
}

//获取设备型号树
func GetDeviceModelTree(c *gin.Context) {
	//设置参数，允许跨域调用

	TAG := "GetDeviceModelTree"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelTree(requestData, dbmap)

	c.JSON(200, rd)
}

//故障管理-获取故障记录
func GetFault(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//校验权限
	TAG := "GetFault"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换
	id := mapPara["Id"].(string)                         //参数：Id
	xtext.RequireNonBlank(id)

	//查询数据
	rd = devdao.GetFault(requestData, id, dbmap)
}

//故障管理-获取教室设备
func GetClassroomDevice(c *gin.Context) {

	TAG := "GetClassroomDevice"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：ClassroomId
	tClassroomId := mapPara["ClassroomId"]
	if tClassroomId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	classroomId, ok := tClassroomId.(string)
	if !ok {
		rd.Rcode = "1003"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetClassroomDevice(requestData, classroomId, dbmap)

	c.JSON(200, rd)
}

//故障管理-获取设备对应型号的所有故障分类
func GetDeviceAllFaultType(c *gin.Context) {
	TAG := "GetDeviceAllFaultType"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tDeviceId := mapPara["DeviceId"]
	if tDeviceId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	deviceId, ok := tDeviceId.(string)
	if !ok {
		rd.Rcode = "1003"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceAllFaultType(requestData, deviceId, dbmap)

	c.JSON(200, rd)
}

//故障管理-获取设备对应型号的所有故障现象词条
func GetDevicFaultWord(c *gin.Context) {
	TAG := "GetDevicFaultWord"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tDeviceId := mapPara["DeviceId"]
	if tDeviceId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	deviceId, ok := tDeviceId.(string)
	if !ok {
		rd.Rcode = "1003"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.DeviceSiteName(requestData, deviceId, dbmap)

	c.JSON(200, rd)
}

//故障管理-故障登记
func RegisterFault(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var request devmodel.RequestRegisterFaultData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(request.Para.Id)
	xtext.RequireNonBlank(request.Para.OT)
	if request.Para.OT == "submit" {
		xtext.RequireNonBlank(request.Para.DeviceId)
		xtext.RequireNonBlank(request.Para.FaultSummary)
		xtext.RequireNonBlank(request.Para.HappenTime)
		request.Para.DataValidate()
	}

	//校验权限
	TAG := "RegisterFault"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, request.Auth.Rolestype, request.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//查询故障表信息
	rd.Rcode = gResponseMsgs.EXE_CMD_FAIL.CodeText()
	rd.Reason = gResponseMsgs.EXE_CMD_FAIL.Text
	ft, err := devdao.QueryFaultTableInfo(request.Para.Id, dbmap)
	xerr.ThrowPanic(err)

	//如果故障已经提交，则不能再更改或添加
	if ft.Id != "" && ft.Status != "0" {
		rd.Rcode = "1002"
		rd.Reason = "故障已经被提交，不能再修改"
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	xerr.ThrowPanic(err)
	tNow := xtime.NowString()
	defer func() {
		if err != nil {
			trans.Rollback()
		}
	}()

	if ft.Id == "" { //故障添加
		err = devdao.RegisterFault_Add(tNow, request, ft, trans)
		xerr.ThrowPanic(err)
	} else { //故障编辑
		if ft.Status == "0" { //草稿状态时，才能被编辑
			err = devdao.RegisterFault_Edit(request, ft, trans)
			xerr.ThrowPanic(err)
		}
	}

	//故障提交
	if request.Para.OT == "submit" {
		err = devdao.RegisterFault_Submit(tNow, request.Para.DeviceId, request.Para.Id, request.Para.IsCanUse, request.Auth.Usersid, trans)
		xerr.ThrowPanic(err)
	}

	//提交事务
	err = trans.Commit()
	xerr.ThrowPanic(err)

	// OK
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
}

//故障管理-故障提交(按故障id单独提交)
func SubmitFault(c *gin.Context) {
	TAG := "SubmitFault"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestFaultIdData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Id == "" {
		rd.Rcode = "1003"
		rd.Reason = "故障Id不能为空！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询故障表信息
	ft, err := devdao.QueryFaultTableInfo(requestData.Para.Id, dbmap)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "查询故障记录出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//如果故障已经提交，则不能再更改或添加
	if ft.Id != "" && ft.Status != "0" {
		rd.Rcode = "1002"
		rd.Reason = "故障已经被提交，不能再修改"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获取当前时间
	t := xtime.NowString()

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//故障提交
	err = devdao.RegisterFault_Submit(t, ft.DeviceId, requestData.Para.Id, ft.IsCanUse, requestData.Auth.Usersid, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "提交故障时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "提交事务时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//故障管理-故障删除
func DeleteFault(c *gin.Context) {
	TAG := "DeleteFault"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestFaultIdData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Id == "" {
		rd.Rcode = "1003"
		rd.Reason = "故障Id不能为空！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询故障表信息
	_, err = devdao.QueryFaultTableInfo(requestData.Para.Id, dbmap)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "查询故障记录出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//故障删除
	err = devdao.DeleteFault(requestData.Para.Id, requestData.Auth.Usersid, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除故障时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除故障时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//故障管理-故障受理
func AcceptanceFault(c *gin.Context) {
	TAG := "AcceptanceFault"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestAcceptanceFaultData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Id == "" {
		rd.Rcode = "1003"
		rd.Reason = "故障Id不能为空！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.RepairPerson == "" {
		rd.Rcode = "1003"
		rd.Reason = "维修人不能为空！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.RepairPersonTel == "" {
		rd.Rcode = "1003"
		rd.Reason = "维修人电话不能为空！"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询故障表信息
	ft, err := devdao.QueryFaultTableInfo(requestData.Para.Id, dbmap)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "查询故障记录出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//如果故障状态不为待受理状态，则不处理
	if ft.Status == "0" {
		rd.Rcode = "1002"
		rd.Reason = "当前故障还未提交，不能受理"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	if ft.Status >= "2" {
		rd.Rcode = "1002"
		rd.Reason = "当前故障已经受理，无需再受理"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//故障受理
	t := xtime.NowString()
	err = devdao.AcceptanceFault(t, requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "提交故障时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "提交事务时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//故障管理-维修登记
func RegisterRepair(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获异常

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	//获得查询参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestRegisterRepairData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(request.Para.Id)
	xtext.RequireNonBlank(request.Para.OT)
	if request.Para.OT != "save" && request.Para.OT != "submit" {
		panic(errors.New("操作类型的值只能是save(暂存)或submit(提交)"))
	}

	if request.Para.OT == "submit" {
		xtext.RequireNonBlank(request.Para.RepairPerson)
		xtext.RequireNonBlank(request.Para.RepairFinishTime)

		if request.Para.RepairIsCanUse != "0" && request.Para.RepairIsCanUse != "1" {
			panic(errors.New("设备是否可以使用的值只能是0(不可使用）和1(可以使用)！"))
		}
		if request.Para.RepairResult != "1" && request.Para.RepairResult != "2" {
			panic(errors.New("维修结果的值只能是1(未修复)和2(已修复)"))
		}
	}

	//校验权限
	TAG := "RegisterRepair"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, request.Auth.Rolestype, request.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//查询故障表信息
	ft, err := devdao.QueryFaultTableInfo(request.Para.Id, dbmap)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "查询故障记录出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//检查故障状态，如果不是维修中状态则不处理
	if ft.Status != "2" {
		rd.Rcode = "1002"
		rd.Reason = "当前故障的状态不是‘维修中’，不能暂存或提交"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	xerr.ThrowPanic(err)

	//暂存维修记录
	tNow := xtime.NowString()
	err = devdao.RegisterRepair_Edit(tNow, request, ft, trans)
	xerr.ThrowPanic(err)

	//提交维修记录
	if request.Para.OT == "submit" {
		err = devdao.RegisterRepair_Submit(tNow, request, ft, trans)
		xerr.ThrowPanic(err)
	}

	//提交事务
	err = trans.Commit()
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = "ok"
}

//节点配置--------------------------------------------------------

//节点配置-获取节点型号列表
func GetNodeModelList(c *gin.Context) {
	TAG := "GetNodeModelList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetNodeModelList(requestData, keyWord, dbmap)

	c.JSON(200, rd)
}

//节点配置-获取节点型号
func GetNodeModel(c *gin.Context) {
	TAG := "GetNodeModel"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetNodeModel(requestData, id, dbmap)

	c.JSON(200, rd)
}

//节点配置-保存节点型号
func SaveNodeModel(c *gin.Context) {
	TAG := "SaveNodeModel"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestNodeModelData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Id == "" {
		rd.Rcode = "1003"
		rd.Reason = "节点型号Id不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.Name == "" {
		rd.Rcode = "1003"
		rd.Reason = "节点型号名称不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveNodeModel(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//节点配置-删除节点型号
func DeleteNodeModel(c *gin.Context) {
	TAG := "DeleteNodeModel"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	if id == "" {
		rd.Rcode = "1002"
		rd.Reason = "错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteNodeModel(id, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

func OnDeletingNodeModel(c *gin.Context) {
	defer xerr.CatchPanic()

	// 响应
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//获得参数
	param := request.Para.(map[string]interface{})
	modelID := param["Id"].(string)
	xtext.RequireNonBlank(modelID)

	//校验权限
	TAG := "DeleteNodeModel"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	exist := zndx.Node_Exists_ByModelID(modelID, dbmap) ||
		zndx.NodeModelCmd_Exists_ByModelID(modelID, dbmap)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = exist
}

//节点配置-获取节点型号命令列表
func GetNodeModelCMDList(c *gin.Context) {
	TAG := "GetNodeModelCMDList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误1"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误2"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误3"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误4"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetNodeModelCMDList(requestData, keyWord, ModelId, dbmap)

	c.JSON(200, rd)
}

//节点配置-获取节点型号命令
func GetNodeModelCMD(c *gin.Context) {
	TAG := "GetNodeModelCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：DeviceId
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误1"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误2:"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetNodeModelCMD(requestData, int(id), dbmap)

	c.JSON(200, rd)
}

//节点配置-保存节点型号命令
func SaveNodeModelCMD(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//获得查询参数
	var requestData devmodel.RequestNodeModelCMDData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(requestData.Para.CmdCode)
	xtext.RequireNonBlank(requestData.Para.ModelId)
	xtext.RequireNonBlank(requestData.Para.CmdName)
	xtext.RequireNonBlank(requestData.Para.RequestType)
	xtext.RequireNonBlank(requestData.Para.RequestURI)
	xtext.RequireNonBlank(requestData.Para.URIQuery)

	// 校验权限
	TAG := "SaveNodeModelCMD"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	// 保存数据
	cmd := zndx.NodeModelCmd{
		Id:             requestData.Para.Id,
		ModelId:        requestData.Para.ModelId,
		CmdCode:        requestData.Para.CmdCode,
		CmdName:        requestData.Para.CmdName,
		RequestType:    requestData.Para.RequestType,
		URIQuery:       requestData.Para.URIQuery,
		CmdDescription: requestData.Para.CmdDescription,
		RequestURI:     requestData.Para.RequestURI,
		Payload:        requestData.Para.Payload,
		CloseCmdFlag:   requestData.Para.CloseCmdFlag,
		OpenCmdFlag:    requestData.Para.OpenCmdFlag}
	err = cmd.Save(dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text

}

//节点配置-删除节点型号命令
func DeleteNodeModelCMD(c *gin.Context) {
	TAG := "DeleteNodeModelCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteNodeModelCMD(int(id), trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

func OnDeletingNodeModelCmd(c *gin.Context) {
	defer xerr.CatchPanic()

	// 响应
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	//解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		return
	}

	//获得参数
	param := request.Para.(map[string]interface{})
	id := param["Id"].(float64) //参数：Id

	//校验权限
	TAG := "DeleteNodeModelCMD"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	pCmd := &zndx.NodeModelCmd{Id: int(id)}
	exist := pCmd.Referenced_ByDevice(dbmap)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = exist
}

//节点配置-获取节点列表
func GetNodeList(c *gin.Context) {
	TAG := "GetNodeList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：路由器Ip，节点名称，节点型号Id
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误1"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误2"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tNodeId := mapPara["NodeId"]
	if tNodeId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误3"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	NodeId, ok := tNodeId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误4"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tCampusids := mapPara["Campusids"]
	if tCampusids == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误5"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Campusids, ok := tCampusids.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误6"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tBuildingids := mapPara["Buildingids"]
	if tBuildingids == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误7"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Buildingids, ok := tBuildingids.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误8"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tFloorsids := mapPara["Floorsids"]
	if tFloorsids == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误9"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Floorsids, ok := tFloorsids.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误10"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tClassRoomIds := mapPara["ClassRoomIds"]
	if tClassRoomIds == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误11"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ClassRoomIds, ok := tClassRoomIds.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误12"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tIsNoSave := mapPara["IsNoSave"]
	if tClassRoomIds == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误11"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	IsNoSave, ok := tIsNoSave.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误12"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	//查询数据
	rd = devdao.GetNodeList(requestData, keyWord, NodeId, Campusids, Buildingids, Floorsids, ClassRoomIds, IsNoSave, dbmap)

	c.JSON(200, rd)
}

//节点配置-获取节点
func GetNode(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获异常

	// 响应数据
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += fmt.Sprintf(" ERR: %s", err.Error())
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	// 参数校验
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换
	node_id := mapPara["NodeId"].(string)

	// 校验权限
	TAG := "GetNode"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, requestData.Auth.Rolestype, requestData.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	// 查询数据
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	list, err := zndxview.QueryList_NodeDetailView_ByNodeID(node_id, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.FOUND_NODATA.CodeText()
	rd.Reason = gResponseMsgs.FOUND_NODATA.Text
	if len(list) > 0 {
		rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
		rd.Reason = gResponseMsgs.SUCCESS.Text
		rd.Result = &devmodel.ResultData{Data: list[0]}
	}
}

//节点配置-保存节点
func SaveNode(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	//获得查询参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestNodeData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(request.Para.Id)
	// xtext.RequireNonBlank(request.Para.IpType)
	// xtext.RequireNonBlank(request.Para.NodeCoapPort)
	// xtext.RequireNonBlank(request.Para.InRouteMappingPort)
	// xtext.RequireNonBlank(request.Para.RouteIp)
	// xtext.RequireNonBlank(request.Para.UploadTime)
	xtext.RequireNonBlank(request.Para.ModelId)
	xnumeric.RequireBetweenInt(request.Para.ClassRoomId, 1, math.MaxInt32)

	//校验权限
	TAG := "SaveNode"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, request.Auth.Rolestype, request.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//保存数据
	rd.Rcode = gResponseMsgs.EXE_CMD_FAIL.CodeText()
	rd.Reason = gResponseMsgs.EXE_CMD_FAIL.Text
	node := zndx.Node{
		Id:                 request.Para.Id,
		Name:               request.Para.Name,
		ModelId:            request.Para.ModelId,
		ClassRoomId:        request.Para.ClassRoomId,
		IpType:             request.Para.IpType,
		NodeCoapPort:       request.Para.NodeCoapPort,
		InRouteMappingPort: request.Para.InRouteMappingPort,
		RouteIp:            request.Para.RouteIp,
		UploadTime:         request.Para.UploadTime}
	err = node.Save(dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
}

//节点配置-删除节点
func DeleteNode(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	buff, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(buff, &request)
	xerr.ThrowPanic(err)

	//获得参数
	param := request.Para.(map[string]interface{})
	NodeId := param["NodeId"].(string) //参数：NodeId

	//校验权限
	TAG := "DeleteNode"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	rd.Rcode = gResponseMsgs.EXE_CMD_FAIL.CodeText()
	rd.Reason = gResponseMsgs.EXE_CMD_FAIL.Text
	pNode := &zndx.Node{Id: NodeId}
	_, err = pNode.Delete(dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = "ok"
}

func OnDeletingNode(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	buff, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(buff, &request)
	xerr.ThrowPanic(err)

	//获得参数
	param := request.Para.(map[string]interface{})
	nodeID := param["NodeId"].(string) //参数：NodeId

	//校验权限
	TAG := "DeleteNode"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	exist := zndx.Device_Exists_ByNodeID(nodeID, dbmap) ||
		zndx.NodeSocketStatus_Exists_ByNodeID(nodeID, dbmap) ||
		zndx.EventSetTable_Exists_ByNodeID(nodeID, dbmap)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = exist
}

//设备配置-获取设备型号列表
func GetDeviceModelList(c *gin.Context) {
	TAG := "GetDeviceModelList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：路由器Ip，节点名称，节点型号Id
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelList(requestData, keyWord, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-获取设备型号
func GetDeviceModel(c *gin.Context) {
	TAG := "GetDeviceModel"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModel(requestData, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-保存设备型号
func SaveDeviceModel(c *gin.Context) {
	TAG := "SaveDeviceModel"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestDeviceModelData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Id == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号Id不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.Type < 1 {
		rd.Rcode = "1003"
		rd.Reason = "设备型号类型不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.Name == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号名称不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveDeviceModel(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-删除设备型号
func DeleteDeviceModel(c *gin.Context) {
	TAG := "deleteDeviceModel"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：ModelId
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteDeviceModel(ModelId, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

func OnDeletingDeviceModel(c *gin.Context) {
	defer xerr.CatchPanic()

	// 数据响应
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//获得参数
	param := request.Para.(map[string]interface{})
	ModelId := param["ModelId"].(string) //参数：ModelId
	xtext.RequireNonBlank(ModelId)

	//校验权限
	TAG := "deleteDeviceModel"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		return
	}

	exist := zndx.DeviceModel_Exists_ByParentID(ModelId, dbmap) ||
		zndx.Device_Exists_ByModelID(ModelId, dbmap) ||
		zndx.DeviceModelControlCmd_Exists_ByModelID(ModelId, dbmap) ||
		zndx.DeviceModelStatusCmd_Exists_ByModelID(ModelId, dbmap) ||
		zndx.DeviceModelStatusValueCode_Exists_ByModelID(ModelId, dbmap) ||
		zndx.DeviceModelFaultType_Exists_ByModelID(ModelId, dbmap) ||
		zndx.DeviceModelFaultWord_Exists_ByModelID(ModelId, dbmap)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = exist
}

//设备配置-获取设备型号列表
func GetDeviceModelStatusCMDList(c *gin.Context) {
	TAG := "GetDeviceModelStatusCMDList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：路由器Ip，节点名称，节点型号Id
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelStatusCMDList(requestData, keyWord, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-获取设备型号
func GetDeviceModelStatusCMD(c *gin.Context) {
	TAG := "GetDeviceModelStatusCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelStatusCMD(requestData, int(Id), dbmap)

	c.JSON(200, rd)
}

//设备配置-保存设备型号状态命令
func SaveDeviceModelStatusCMD(c *gin.Context) {
	TAG := "SaveDeviceModelStatusCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestDeviceModelStatusCMDData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.StatusCode == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态命令编码不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.StatusName == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态命令名称不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.StatusValueMatchString == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态命令格式不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.Payload == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态命令定义不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.ModelId == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.SwitchStatusFlag == "" {
		requestData.Para.SwitchStatusFlag = "0"
	}
	if requestData.Para.SelectValueFlag == "" {
		requestData.Para.SelectValueFlag = "0"
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveDeviceModelStatusCMD(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-删除设备型号
func DeleteDeviceModelStatusCMD(c *gin.Context) {
	TAG := "DeleteDeviceModelStatusCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：ModelId
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	//参数：StatusCode
	tStatusCode := mapPara["StatusCode"]
	if tStatusCode == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	StatusCode, ok := tStatusCode.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteDeviceModelStatusCMD(ModelId, StatusCode, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

// OnDeletingDeviceModelStatusCmd
func OnDeletingDeviceModelStatusCmd(c *gin.Context) {
	defer xerr.CatchPanic()

	//数据响应
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "DeleteDeviceModelStatusCMD"
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	//获得参数
	mapPara := request.Para.(map[string]interface{})
	model_id := mapPara["ModelId"].(string) //参数：ModelId
	xtext.RequireNonBlank(model_id)
	status_code := mapPara["StatusCode"].(string) //参数：StatusCode
	xtext.RequireNonBlank(status_code)

	pCmd := &zndx.DeviceModelStatusCmd{ModelId: model_id, StatusCode: status_code}
	exist := zndx.DeviceModelStatusValueCode_Exists_ByStatusCmd(pCmd, dbmap)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = exist
}

//设备配置-获取设备型号状态编码列表
func GetDeviceModelStatusValueCodeList(c *gin.Context) {
	TAG := "GetDeviceModelStatusValueCodeList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tStatusCode := mapPara["StatusCode"]
	if tStatusCode == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	StatusCode, ok := tStatusCode.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelStatusValueCodeList(requestData, keyWord, StatusCode, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-获取设备型号状态编码
func GetDeviceModelStatusValueCode(c *gin.Context) {
	TAG := "GetDeviceModelStatusValueCode"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelStatusValueCode(requestData, int(Id), dbmap)

	c.JSON(200, rd)
}

//设备配置-保存设备型号状态编码
func SaveDeviceModelStatusValueCode(c *gin.Context) {
	TAG := "SaveDeviceModelStatusValueCode"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestDeviceModelStatusValueCodeData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.StatusCode == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.StatusValueCode == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态命令编码值不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.StatusValueName == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号状态命令名称值不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.ModelId == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveDeviceModelStatusValueCode(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-删除设备型号状态编码
func DeleteDeviceModelStatusValueCode(c *gin.Context) {
	TAG := "DeleteDeviceModelStatusValueCode"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteDeviceModelStatusValueCode(int(Id), trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-获取设备型号控制命令列表
func GetDeviceModelControlCMDList(c *gin.Context) {
	TAG := "GetDeviceModelControlCMDList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelControlCMDList(requestData, keyWord, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-获取设备型号控制命令详细
func GetDeviceModelControlCMD(c *gin.Context) {
	TAG := "GetDeviceModelControlCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelControlCMD(requestData, int(Id), dbmap)

	c.JSON(200, rd)
}

//设备配置-保存设备型号控制命令
func SaveDeviceModelControlCMD(c *gin.Context) {
	TAG := "SaveDeviceModelControlCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestDeviceModelControlCMDData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.ModelId == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号Id不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.CmdCode == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号控制命令代码不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.CmdName == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号控制命令名称不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.RequestURI == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号请求地址不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.URIQuery == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号请求地址参数不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.RequestType == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号请求类型不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.Payload == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号控制命令不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveDeviceModelControlCMD(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-删除设备型号控制命令
func DeleteDeviceModelControlCMD(c *gin.Context) {
	TAG := "DeleteDeviceModelControlCMD"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteDeviceModelControlCMD(int(Id), trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

func OnDeletingDeviceModelControlCmd(c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//获得参数
	param := request.Para.(map[string]interface{})
	Id := param["Id"].(float64) //参数：Id

	//校验权限
	TAG := "DeleteDeviceModelControlCMD"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		return
	}

	cmd := zndx.DeviceModelControlCmd{Id: int(Id)}
	referenced := cmd.Referenced_ByDevice(dbmap)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = referenced
}

//设备配置-获取设备列表
func GetDeviceList(c *gin.Context) {
	TAG := "GetDeviceList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tBuildingid := mapPara["Buildingid"]
	if tBuildingid == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Buildingid, ok := tBuildingid.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tFloorsid := mapPara["Floorsid"]
	if tFloorsid == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Floorsid, ok := tFloorsid.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tCampusid := mapPara["Campusid"]
	if tCampusid == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Campusid, ok := tCampusid.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tClassroomId := mapPara["ClassroomId"]
	if tClassroomId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ClassroomId, ok := tClassroomId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	//查询数据
	rd = devdao.GetDeviceList(requestData, keyWord, ModelId, int(Buildingid), int(Floorsid), int(Campusid), int(ClassroomId), dbmap)

	c.JSON(200, rd)
}

//设备配置-获取设备
func GetDevice(c *gin.Context) {
	TAG := "GetDevice"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDevice(requestData, Id, dbmap)

	c.JSON(200, rd)
}

//设备配置-保存设备
func SaveDevice(c *gin.Context) {
	defer xerr.CatchPanic() //捕获异常

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	// 解析数据
	var request devmodel.RequestDeviceData
	data, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println(string(data))
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(request.Para.Id)
	xtext.RequireNonBlank(request.Para.Name)
	xtext.RequireNonBlank(request.Para.ModelId)
	if request.Para.JoinMethod == "node" {
		xtext.RequireNonBlank(request.Para.PowerNodeId)
		xtext.RequireNonBlank(request.Para.JoinNodeId)
		xtext.RequireEqual(request.Para.PowerNodeId, request.Para.JoinNodeId)
	}

	//校验权限
	TAG := "SaveDevice"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if ok := doAuthValidate(TAG, request.Auth.Rolestype, request.Auth.Usersid, dbmap); !ok {
		rd.Rcode = gResponseMsgs.AUTH_LIMITED.CodeText()
		rd.Reason = gResponseMsgs.AUTH_LIMITED.Text
		return
	}

	//开启事务
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	trans, err := dbmap.Begin()
	xerr.ThrowPanic(err)

	// 节点验证
	pPowerNode := &zndx.Node{Id: request.Para.PowerNodeId} // 电源节点
	pJoinNode := &zndx.Node{Id: request.Para.JoinNodeId}   // 接入节点
	if !zndx.Node_Exists(pPowerNode.Id, dbmap) && !zndx.Node_Exists(pJoinNode.Id, dbmap) {
		rd.Rcode = gResponseMsgs.EXE_CMD_FAIL.CodeText()
		rd.Reason = gResponseMsgs.EXE_CMD_FAIL.Text + " 节点验证失败"
		return
	}

	//保存数据
	err = devdao.SaveDevice(request, trans)
	xerr.ThrowPanic(err)

	//提交事务
	err = trans.Commit()
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = "ok"
}

//设备配置-删除设备
func DeleteDevice(c *gin.Context) {
	var rd core.Returndata
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		return
	}

	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "DeleteDevice"
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	dev := zndx.Device{Id: Id}
	_, err = dev.Delete(dbmap)
	xdebug.LogError(err)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
	}

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
}

func OnDeletingDevice(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		c.JSON(http.StatusOK, rd)
	}()

	//获得查询参数
	buff, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err := json.Unmarshal(buff, &request)
	xerr.ThrowPanic(err)

	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "DeleteDevice"
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	//获得参数map
	kv := request.Para.(map[string]interface{})
	devID := kv["DeviceID"].(string) //参数：DeviceID

	// 查询设备被依赖否
	exist := zndx.DeviceUseLog_Exists_ByDeviceID(devID, dbmap) ||
		zndx.DeviceOperateLog_Exists_ByDeviceID(devID, dbmap) ||
		zndx.DeviceDetailLog_Exists_ByDeviceID(devID, dbmap) ||
		zndx.DeviceFault_Exists_ByDeviceID(devID, dbmap) ||
		zndx.DeviceAlert_Exists_ByDeviceID(devID, dbmap) ||
		zndx.DeviceLastSendContent_Exists_ByDeviceID(devID, dbmap) ||
		zndx.PJLink_Exists_ByDeviceID(devID, dbmap) ||
		zndx.EventSetTable_Exists_ByDeviceID(devID, dbmap)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = exist
}

//设备配置-获取故障分类列表
func GetDeviceModelFaultTypeList(c *gin.Context) {
	TAG := "GetDeviceModelFaultTypeList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	//查询数据
	rd = devdao.GetDeviceModelFaultTypeList(requestData, keyWord, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-获取故障分类
func GetDeviceModelFaultType(c *gin.Context) {
	TAG := "GetDeviceModelFaultType"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelFaultType(requestData, Id, dbmap)

	c.JSON(200, rd)
}

//设备配置-保存故障分类
func SaveDeviceModelFaultType(c *gin.Context) {
	TAG := "SaveDeviceModelFaultType"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestDeviceModelFaultTypeData
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Id == "" {
		rd.Rcode = "1003"
		rd.Reason = "故障分类Id不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.Name == "" {
		rd.Rcode = "1003"
		rd.Reason = "故障分类名称不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.ModelId == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveDeviceModelFaultType(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-删除故障分类
func DeleteDeviceModelFaultType(c *gin.Context) {
	TAG := "DeleteDeviceModelFaultType"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteDeviceModelFaultType(Id, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//OnDeletingDeviceModelFaultType
func OnDeletingDeviceModelFaultType(c *gin.Context) {
	defer xerr.CatchPanic()

	//数据响应
	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() { c.JSON(http.StatusOK, rd) }()

	//解析数据
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request devmodel.RequestData
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	TAG := "DeleteDeviceModelFaultType"
	rd = userdao.CheckVaild(request.Auth.Rolestype, request.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		return
	}

	//获得参数
	mapPara := request.Para.(map[string]interface{})
	id := mapPara["Id"].(string) //参数：Id
	xtext.RequireNonBlank(id)

	pFaultType := &zndx.DeviceModelFaultType{Id: id}
	referenced := pFaultType.Referenced_ByDeviceFault(dbmap)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = referenced
}

//设备配置-获取故障现象常用词条列表
func GetDeviceModelFaultWordList(c *gin.Context) {
	TAG := "GetDeviceModelFaultWordList"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	tModelId := mapPara["ModelId"]
	if tModelId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	ModelId, ok := tModelId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	//查询数据
	rd = devdao.GetDeviceModelFaultWordList(requestData, keyWord, ModelId, dbmap)

	c.JSON(200, rd)
}

//设备配置-获取故障现象常用词条
func GetDeviceModelFaultWord(c *gin.Context) {
	TAG := "GetDeviceModelFaultWord"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//查询数据
	rd = devdao.GetDeviceModelFaultWord(requestData, int(Id), dbmap)

	c.JSON(200, rd)
}

//设备配置-保存故障现象常用词条
func SaveDeviceModelFaultWord(c *gin.Context) {
	TAG := "SaveDeviceModelFaultWord"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestDeviceModelFaultWordData
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}
	//校验数据
	if requestData.Para.Name == "" {
		rd.Rcode = "1003"
		rd.Reason = "故障现象常用词条名称不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	if requestData.Para.ModelId == "" {
		rd.Rcode = "1003"
		rd.Reason = "设备型号不能为空!"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//保存数据
	err = devdao.SaveDeviceModelFaultWord(requestData, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "保存数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

//设备配置-删除故障现象常用词条
func DeleteDeviceModelFaultWord(c *gin.Context) {
	TAG := "DeleteDeviceModelFaultWord"
	var rd core.Returndata

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason+":"+err.Error())
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//校验权限
	rd = userdao.CheckVaild(requestData.Auth.Rolestype, requestData.Auth.Usersid, TAG, dbmap)
	if rd.Rcode != "1000" { //权限未通过验证
		rd.Rcode = "1002"
		rd.Reason = "未通过权限验证"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//参数：Id
	tId := mapPara["Id"]
	if tId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}
	Id, ok := tId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		return
	}

	//删除数据
	err = devdao.DeleteDeviceModelFaultWord(int(Id), trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除数据时出错：" + err.Error()
		c.JSON(200, rd)
		log.Println(TAG, rd.Reason)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}
