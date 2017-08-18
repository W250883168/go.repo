package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	gorp "gopkg.in/gorp.v1"

	"canopus"

	"dborm/zndx"
	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xtext"
	"xutils/xtime"

	"dev.project/BackEndCode/devcontrol/coap/coapclient"
	"dev.project/BackEndCode/devcontrol/dal"
	"dev.project/BackEndCode/devcontrol/ioutil/dbutil"
)

// 切换器请求
type _VnriverMultiplexer_Request struct {
	UserID   string // 用户ID
	DeviceID string // 设备ID
	InPort   int    // 输入端口(int)
	OutPort  int    // 输出端口(int)
	Params   string // 参数（可选）
}

// 视和讯切换器处理
func Device_MultiplexerVnriver_Handler(cmd string, c *gin.Context) {
	defer xerr.CatchPanic()

	resp := xhttp.HttpResponse{Content: gin.H{"code": "1", "msg": "无效参数", "data": ""}}
	defer func() { c.JSON(http.StatusOK, resp.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _VnriverMultiplexer_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.DeviceID)
	xtext.RequireNonBlank(request.UserID)

	// 查询设备
	dbmap := dbutil.GetDBMap()
	device_id := request.DeviceID
	pDevice, err := zndx.Device_Get(device_id, dbmap)
	xerr.ThrowPanic(err)

	// 保存数据
	payload_data := []byte{0xfe, 0xfe, 00, 0x31, byte(request.InPort), byte(request.OutPort), 0xaa, 0xaa} // 切换器串口命令
	log.Printf("Payload_Data: %x\n", payload_data)
	onVnriverMultiplexer(cmd, fmt.Sprintf("%x", payload_data), &request, dbmap)

	// 查询设备命令
	_, err = dal.Query_DevicePowerCmdView_ByDeviceID(device_id, cmd, dbmap)
	xerr.ThrowPanic(err)

	// 查询节点命令并发送
	node_cmds, err := dal.QueryList_NodeSwitchCmd_ByNode(&zndx.Node{Id: pDevice.JoinNodeId}, cmd, dbmap)
	xerr.ThrowPanic(err)
	for _, node_cmd := range node_cmds {
		coapcmd := coapclient.CoapCommand{
			HostAddr:    fmt.Sprintf("%s:%s", node_cmd.RouterIP, node_cmd.InRouterMappingPort),
			Method:      canopus.Post,
			RequestURI:  fmt.Sprintf("/rs232/rs232_%s", pDevice.JoinSocketId),
			QueryParams: map[string]string{"eui": pDevice.JoinNodeId},
			Payload:     string(payload_data)}
		if reply := coapclient.Send(coapcmd); reply != nil {
			log.Printf("RevcPayload: %+v\n", reply.GetMessage())
		}
	}

	resp.Content = gin.H{"code": "0", "msg": "执行成功", "data": "ok"}
}

// 切换器数据处理
func onVnriverMultiplexer(cmd, payload string, pRequest *_VnriverMultiplexer_Request, dbmap *gorp.DbMap) {
	device_id := pRequest.DeviceID
	nowTime := xtime.NowString()

	// 保存设备命令最后发送的内容
	pSendContent := &zndx.DeviceLastSendContent{DeviceId: device_id, CmdCode: cmd, LastSendContent: payload, SendTime: nowTime}
	pSendContent.Save(dbmap)

	//生成操作日志
	op_log := zndx.DeviceOperateLog{
		OperateTime:   nowTime,
		OperateUserId: pRequest.UserID,
		OperateType:   "other",
		OperateObject: "device",
		ObjectId:      pRequest.DeviceID,
		CmdCode:       cmd,
		CmdName:       _CmdCode_ToName(cmd, ""),
		Para:          pRequest.Params}
	err := op_log.Insert(dbmap)
	xdebug.LogError(err)

	//生成设备详细日志
	user_id, _ := strconv.Atoi(pRequest.UserID)
	detail_log := zndx.DeviceDetailLog{
		OperateUserId: user_id,
		OperateTime:   nowTime,
		DeviceId:      pRequest.DeviceID,
		CmdCode:       cmd,
		CmdName:       _CmdCode_ToName(cmd, ""),
		Para:          pRequest.Params}
	err = detail_log.Insert(dbmap)
	xdebug.LogError(err)
}
