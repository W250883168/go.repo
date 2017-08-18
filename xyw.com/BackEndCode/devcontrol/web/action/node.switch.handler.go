package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gorp "gopkg.in/gorp.v1"

	"canopus"

	"dborm/zndx"
	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xtext"
	"xutils/xtime"

	"dev.project/BackEndCode/devcontrol/coap/coapclient"
	"dev.project/BackEndCode/devcontrol/dal"
	"dev.project/BackEndCode/devcontrol/ioutil/dbutil"
)

// 节点开关控制命令请求
type _NodeSwitch_ControlHandler_Request struct {
	UserID          string // 用户ID
	NodeID          string // 节点设备ID
	NodeSwitchIndex int    // 节点插口序号
	Params          string // 参数（可选）
}

// 教室内节点开关控制命令请求
type _NodeSwitch_RoomControlHandler_Request struct {
	UserID string // 用户ID
	RoomID string // 位置ID
	Params string // 参数（可选）
}

// 节点开关处理
func NodeSwitch_ControlHandler(cmd string, c *gin.Context) {
	defer xerr.CatchPanic()

	resp := xhttp.HttpResponse{Content: gin.H{"code": "1", "msg": "无效参数", "data": ""}}
	defer func() { c.JSON(http.StatusOK, resp.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _NodeSwitch_ControlHandler_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.NodeID)
	xtext.RequireNonBlank(request.UserID)

	// 查询设备
	dbmap := dbutil.GetDBMap()
	pNode, err := zndx.Node_Get(request.NodeID, dbmap)
	xerr.ThrowPanic(err)

	// 查询节点命令并发送
	node_cmds, err := dal.QueryList_NodeSwitchCmd_ByNode(pNode, cmd, dbmap)
	xerr.ThrowPanic(err)
	for _, nodecmd := range node_cmds {
		coapcmd := coapclient.CoapCommand{
			HostAddr:    fmt.Sprintf("%s:%s", nodecmd.RouterIP),
			Method:      canopus.Post,
			RequestURI:  fmt.Sprintf("/smart_switch/switch_%s", request.NodeSwitchIndex),
			QueryParams: map[string]string{"eui": request.NodeID},
			Payload:     nodecmd.Payload}
		reply := coapclient.Send(coapcmd)
		log.Println("Payload: ", reply.GetMessage().Payload.String())
	}

	resp.Content = gin.H{"code": "0", "msg": "执行成功", "data": "ok"}
}

// 节点开关控制数据处理
func onNodeSwitch_ControlHandler(cmd, payload string, pRequest *_NodeSwitch_ControlHandler_Request, dbmap *gorp.DbMap) {
	nowTime := xtime.NowString()

	// 生成日志（操作日志和详细日志）
	cmd_name := _CmdCode_ToName(cmd, "")
	opLog := zndx.DeviceOperateLog{0, nowTime, pRequest.UserID, cmd, "node", pRequest.NodeID, "", "", cmd, cmd_name, pRequest.Params}
	opLog.Insert(dbmap)
	devLog := zndx.DeviceOnOffLogView{nowTime, pRequest.UserID, "node", pRequest.NodeID, cmd, cmd_name, pRequest.Params}
	devLog.CreateLog(dbmap)
}
