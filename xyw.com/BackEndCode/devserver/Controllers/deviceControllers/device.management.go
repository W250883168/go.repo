package deviceControllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"xutils/xerr"
	//	"xutils/xhttp"
	"xutils/xtext"

	"dborm/zndx"
	"dev.project/BackEndCode/devserver/model/core"
)

// 设备解绑请求
type _DeviceUnbind_Request struct {
	UserID   string
	DeviceID string
	NodeID   string
}

// 设备解绑处理
func Do_DeviceUnbind_Handler(c *gin.Context) {
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
	var request _DeviceUnbind_Request
	data, _ := ioutil.ReadAll(c.Request.Body)
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 验证数据
	xtext.RequireNonBlank(request.UserID)
	xtext.RequireNonBlank(request.NodeID)
	xtext.RequireNonBlank(request.DeviceID)

	// 查询数据
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	pDevice, err := zndx.Device_Get(request.DeviceID, dbmap)
	xerr.ThrowPanic(err)

	// 解绑电源节点
	if pDevice.PowerNodeId == request.NodeID {
		pDevice.PowerNodeId = ""
		pDevice.PowerSwitchId = ""
	}
	// 解绑通讯节点
	if pDevice.JoinNodeId == request.NodeID {
		pDevice.JoinMethod = ""
		pDevice.JoinNodeId = ""
		pDevice.JoinSocketId = ""
	}
	// 保存数据
	_, err = pDevice.Update(dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = "ok"
}

func init() {
	log.Print("")
}
