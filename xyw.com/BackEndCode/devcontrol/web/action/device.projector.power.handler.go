package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"

	"canopus"

	"dborm/zndx"
	"dborm/zndxview"
	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xtext"
	"xutils/xtime"

	"dev.project/BackEndCode/devcontrol/coap/coapclient"
	"dev.project/BackEndCode/devcontrol/dal"
	"dev.project/BackEndCode/devcontrol/ioutil/dbutil"
)

// EPSON投影仪请求
type _EpsonProjector_Request struct {
	UserID   string // 用户ID
	DeviceID string // 设备ID
	Params   string // 参数（可选）
}

// 投影仪请求处理
func Device_ProjectorEpson_Handler(cmd, payload, uri string, c *gin.Context) {
	defer xerr.CatchPanic()

	var err error
	resp_content := gin.H{"code": "1", "msg": "无效参数", "data": ""}
	resp := xhttp.HttpResponse{Content: resp_content}
	defer func() {
		if err != nil {
			resp_content["data"] = err.Error()
		}
		c.JSON(http.StatusOK, resp.Content)
	}()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _EpsonProjector_Request
	err = json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.DeviceID)
	xtext.RequireNonBlank(request.UserID)

	// 查询设备
	dbmap := dbutil.GetDBMap()
	pDevice, err := zndx.Device_Get(request.DeviceID, dbmap)
	xerr.ThrowPanic(err)

	// 保存数据
	// onProjectorPower(cmd, payload, &request, dbmap)
	onProjector_SwitchHandler(pDevice, cmd, &request, dbmap)

	if cmd == "off" {
		// 查询设备命令并发送命令
		dev_cmds, err := dal.Query_DevicePowerCmdView_ByDeviceID(request.DeviceID, cmd, dbmap)
		xerr.ThrowPanic(err)
		//		log.Println("-----------查询设备命令并发送命令")
		log.Println("--------------// 查询设备命令并发送命令", dev_cmds)
		for _, dev_cmd := range dev_cmds {
			if !strings.HasPrefix(dev_cmd.RequestURI, "/") {
				dev_cmd.RequestURI = "/" + dev_cmd.RequestURI
			}
			coapcmd := coapclient.CoapCommand{
				HostAddr:    fmt.Sprintf("%s:%s", dev_cmd.RouterIP, dev_cmd.InRouteMappingPort),
				Method:      canopus.Post,
				RequestURI:  dev_cmd.RequestURI,
				QueryParams: map[string]string{"eui": pDevice.JoinNodeId},
				Payload:     "PWR OFF"}
			if reply := coapclient.Send(coapcmd); reply != nil {
				log.Printf("RevcPayload: %+v\n", reply.GetMessage())
			}

			time.Sleep(time.Second) // 等待设备关闭动作完成
		}

		// 查询节点命令并发送命令(投影仪设备/EPSON)
		node_cmds, err := dal.QueryList_NodeSwitchCmd_ByNode(&zndx.Node{Id: pDevice.JoinNodeId}, cmd, dbmap)
		xerr.ThrowPanic(err)
		log.Println("---------------查询节点命令并发送命令")
		for _, nodecmd := range node_cmds {
			if !strings.HasPrefix(nodecmd.RequestURI, "/") {
				nodecmd.RequestURI = "/" + nodecmd.RequestURI
			}
			coapcmd := coapclient.CoapCommand{
				HostAddr:    fmt.Sprintf("%s:%s", nodecmd.RouterIP, nodecmd.InRouterMappingPort),
				Method:      canopus.Post,
				RequestURI:  nodecmd.RequestURI,
				QueryParams: map[string]string{"eui": pDevice.PowerNodeId},
				Payload:     "SWITCH OFF"}
			if reply := coapclient.Send(coapcmd); reply != nil {
				log.Printf("RevcPayload: %+v\n", reply.GetMessage())
			}
		}
	}

	if cmd == "on" {
		// 查询节点命令并发送命令(投影仪设备/EPSON)
		log.Println("// 查询节点命令并发送命令(投影仪设备/EPSON)")
		node_cmds, err := dal.QueryList_NodeSwitchCmd_ByNode(&zndx.Node{Id: pDevice.JoinNodeId}, cmd, dbmap)
		xerr.ThrowPanic(err)
		for _, nodecmd := range node_cmds {
			if !strings.HasPrefix(nodecmd.RequestURI, "/") {
				nodecmd.RequestURI = "/" + nodecmd.RequestURI
			}
			coapcmd := coapclient.CoapCommand{
				HostAddr:    fmt.Sprintf("%s:%s", nodecmd.RouterIP, nodecmd.InRouterMappingPort),
				Method:      canopus.Post,
				RequestURI:  nodecmd.RequestURI,
				QueryParams: map[string]string{"eui": pDevice.PowerNodeId},
				Payload:     "SWITCH ON"}
			if reply := coapclient.Send(coapcmd); reply != nil {
				log.Printf("RevcPayload: %+v\n", reply.GetMessage())
			}

			time.Sleep(time.Second) // 等待设备上电动作完成
		}

		// 查询设备命令并发送命令
		dev_cmds, err := dal.Query_DevicePowerCmdView_ByDeviceID(request.DeviceID, cmd, dbmap)
		xerr.ThrowPanic(err)
		log.Println("--------------// 查询设备命令并发送命令", dev_cmds)
		for _, dev_cmd := range dev_cmds {
			if !strings.HasPrefix(dev_cmd.RequestURI, "/") {
				dev_cmd.RequestURI = "/" + dev_cmd.RequestURI
			}
			coapcmd := coapclient.CoapCommand{
				HostAddr:    fmt.Sprintf("%s:%s", dev_cmd.RouterIP, dev_cmd.InRouteMappingPort),
				Method:      canopus.Post,
				RequestURI:  dev_cmd.RequestURI,
				QueryParams: map[string]string{"eui": pDevice.JoinNodeId},
				Payload:     "PWR ON"}
			if reply := coapclient.Send(coapcmd); reply != nil {
				log.Printf("RevcPayload: %+v\n", reply.GetMessage())
			}
		}
	}

	resp.Content = gin.H{"code": "0", "msg": "执行成功", "data": "ok"}
}

// 投影仪开关处理（数据）
func onProjector_SwitchHandler(pDevice *zndx.Device, cmd string, pRequest *_EpsonProjector_Request, dbmap *gorp.DbMap) {
	defer xerr.CatchPanic() // 捕获异常
	device_id := pDevice.Id
	nowTime := xtime.NowString()

	//生成日志（操作日志和详细日志）
	cmd_name := _CmdCode_ToName(cmd, "")
	opLog := zndx.DeviceOperateLog{0, nowTime, pRequest.UserID, cmd, "device", device_id, "", "", cmd, cmd_name, pRequest.Params}
	opLog.Insert(dbmap)
	devLog := zndx.DeviceOnOffLogView{nowTime, pRequest.UserID, "device", device_id, cmd, cmd_name, pRequest.Params}
	devLog.CreateLog(dbmap)

	// 处理设备的开/关
	//1) 更新设备自身状态
	pDevice.DeviceSelfStatus = cmd
	err := pDevice.Update_DeviceSelfStatus(dbmap)
	xerr.ThrowPanic(err)

	//2) 更新设备开关记录
	exist := zndx.DeviceUseLog_Exists_DeviceOpenLog(device_id, dbmap)
	if cmd == "on" && !exist { // 开机操作 && 无开机记录
		zndx.DeviceUseLog_Insert_DeviceOpenLog(device_id, xtime.NowString(), dbmap)
	}

	if cmd == "off" && exist { // 关机操作 && 有开机记录
		devLog, err := zndx.DeviceUseLog_Query_DeviceOpenLog(device_id, dbmap)
		xerr.ThrowPanic(err)
		// 计算累计使用时间
		powerOnTime, err := time.Parse(xtime.FormatString(), devLog.OnTime)
		powerOffTime, err := time.Parse(xtime.FormatString(), nowTime)
		useTime := int64(powerOffTime.Sub(powerOnTime).Seconds()) //使用时间（秒）

		// 写入关闭时间、计算本次使用时间
		devLog.OffTime = nowTime
		devLog.UseTime = useTime
		_, err = devLog.Update(dbmap)
		xerr.ThrowPanic(err)

		//更新设备累计使用时间
		dev := zndx.Device{Id: device_id, UseTimeAfter: useTime}
		err = dev.Update_DeviceUseTime(dbmap)
		xerr.ThrowPanic(err)

		// 处理预警(按使用时间预警)
		pDeviceEx, err := zndxview.Query_DeviceExNodeView(device_id, dbmap)
		if (err == nil) && (pDeviceEx.IsAlert) {
			pDevAlert := zndx.DeviceAlert{DeviceId: device_id}
			err = pDevAlert.Delete_DeviceOvertime(device_id, dbmap)
			xerr.ThrowPanic(err)

			//如果累计使用时间大于最大使用时间，则创建新的使用日期预警
			totalUseTime := pDeviceEx.UseTimeBefore + pDeviceEx.UseTimeAfter + useTime
			if totalUseTime > pDeviceEx.MaxUseTime {
				total_use_tspan := xtime.TimeSpan(time.Second * time.Duration(totalUseTime))
				max_use_tspan := xtime.TimeSpan(time.Second * time.Duration(pDeviceEx.MaxUseTime))
				tspan_format := "[%0.2d小时%0.2d分钟%0.2d秒]"
				alert_desc := "累计使用时间" + total_use_tspan.ToString(tspan_format) + "超过了设定的最大使用时间" + max_use_tspan.ToString(tspan_format)

				pDevAlert.AlertType = "1" // 超时使用预警
				pDevAlert.AlertDescription = alert_desc
				pDevAlert.LastAlertTime = nowTime
				pDevAlert.Insert(dbmap)
				xerr.ThrowPanic(err)
			}
		}
	}
}

// 投影仪电源请求处理（数据）
func onProjectorPower(cmd, payload string, pRequest *_EpsonProjector_Request, dbmap *gorp.DbMap) {
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
		// ObjectName
		// UseWhoseCmd
		CmdCode: cmd,
		CmdName: _CmdCode_ToName(cmd, ""),
		Para:    pRequest.Params}
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
