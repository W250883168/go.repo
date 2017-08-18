package deviceControllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"xutils/xerr"

	"dev.project/BackEndCode/devserver/DataAccess/deviceDataAccess"
	"dev.project/BackEndCode/devserver/model/core"
)

func GetDeviceCmd_ByRoom(c *gin.Context) {
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
	var request = struct {
		RoomID  int
		CmdCode string
	}{}
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = gResponseMsgs.FAIL.Text
	ret, err := deviceDataAccess.Query_DeviceCmd_ByRoom(request.RoomID, request.CmdCode, dbmap)
	xerr.ThrowPanic(err)

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = ret
}
