package deviceControllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"dborm/zndx"
	"xutils/xerr"
	"xutils/xhttp"

	core "dev.project/BackEndCode/devserver/model/core"
	devmodel "dev.project/BackEndCode/devserver/model/deviceModel"
)

// 获取所有设备故障信息
func GetAllFaultInfoList4App(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获错误

	var err error
	var rd = core.Returndata{Rcode: gResponseMsgs.DATA_MALFORMED.CodeText(), Reason: gResponseMsgs.DATA_MALFORMED.Text}
	defer func() {
		if err != nil {
			rd.Reason += "; ERR: " + err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	//获得查询参数
	var requestData devmodel.RequestData
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//获得参数map
	params := requestData.Para.(map[string]interface{})
	siteType := params["SiteType"].(string)       // 参数：SiteType
	siteId := params["SiteId"].(string)           // 参数：SiteId
	modelId := params["ModelId"].(string)         // 参数：ModelId
	keyWord := params["KeyWord"].(string)         // 参数：KeyWord
	devState := int(params["DevState"].(float64)) // 参数：DevState
	// log.Println(requestData)

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
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	dFault := zndx.DeviceFault{InputUserId: requestData.Auth.Usersid, Status: strconv.Itoa(devState)}
	pPage := &xhttp.PageInfo{PageIndex: requestData.Page.PageIndex, PageSize: requestData.Page.PageSize}
	list, err := dFault.Query_DeviceFault(siteType, siteId, modelId, keyWord, dbmap, pPage)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{devmodel.PageData{
		PageIndex:   pPage.PageIndex,
		PageSize:    pPage.PageSize,
		RecordCount: pPage.RowTotal,
		PageCount:   pPage.PageTotal()}, list}
}
