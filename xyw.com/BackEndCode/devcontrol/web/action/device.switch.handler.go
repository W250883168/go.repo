package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gorp "gopkg.in/gorp.v1"

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

// 设备开关命令请求
type _Device_SwitchHandler_Request struct {
	UserID   string // 用户ID
	DeviceID string // 设备ID
	Params   string // 参数（可选）
}

// 房间设备一键开关命令请求
type _RoomDevice_SwitchHandler_Request struct {
	UserID string // 用户ID
	RoomID string // 位置ID
	Params string // 参数（可选）
}

type _FloorDevice_SwitchHandler_Request struct {
	UserID  string // 用户ID
	FloorID string // 位置ID
	Params  string // 参数（可选）
}

// 楼层设备一键开关处理(By FloorID)
func FloorDevice_SwitchHandler(c *gin.Context, cmd, payload string) {
	defer xerr.CatchPanic()

	resp_content := gin.H{"code": "1", "msg": "无效参数", "data": ""}
	resp := xhttp.HttpResponse{Content: &resp_content}
	defer func() { c.JSON(http.StatusOK, resp.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _FloorDevice_SwitchHandler_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.FloorID)
	xtext.RequireNonBlank(request.UserID)

	// 数据查询
	resp_content["data"] = "数据查询错误"
	dbmap := dbutil.GetDBMap()

	//生成日志（操作日志和详细日志）
	nowTime := xtime.NowString()
	cmd_name := _CmdCode_ToName(cmd, "")
	opLog := zndx.DeviceOperateLog{0, nowTime, request.UserID, cmd, "floor", request.FloorID, "", "", cmd, cmd_name, request.Params}
	opLog.Insert(dbmap)
	devLog := zndx.DeviceOnOffLogView{nowTime, request.UserID, "floor", request.FloorID, cmd, cmd_name, request.Params}
	devLog.CreateLog(dbmap)

	dev_list, err := zndx.Device_QueryByFloorID(request.FloorID, dbmap)
	xerr.ThrowPanic(err)
	onRoomDevice_SwitchHandler(dev_list, cmd, dbmap)

	// 发送命令
	for _, dev := range dev_list {
		node_cmd, err := dal.Query_NodeSwitchCmd_ByDevice(&dev, cmd, dbmap)
		xdebug.LogError(err)
		if err == nil {
			coapcmd := coapclient.CoapCommand{
				HostAddr:    fmt.Sprintf("%s:%s", node_cmd.RouterIP, node_cmd.InRouterMappingPort),
				Method:      canopus.Post,
				RequestURI:  fmt.Sprintf("/smart_switch/switch_%s", dev.PowerSwitchId),
				QueryParams: map[string]string{"eui": dev.PowerNodeId},
				Payload:     payload}
			go func() {
				if reply := coapclient.Send(coapcmd); reply != nil {
					log.Printf("RevcPayload: %+v\n", reply.GetMessage())
				}
			}()
		}
		// 每次命令发送间隔100ms
		time.Sleep(time.Millisecond * 100)
	}

	resp.Content = gin.H{"code": "0", "msg": "执行成功", "data": "ok"}
}

// 节点开关处理(By DeviceID)
func Device_SwitchHandler(c *gin.Context, cmd, payload string) {
	defer xerr.CatchPanic()

	resp_content := gin.H{"code": "1", "msg": "无效参数", "data": ""}
	resp := xhttp.HttpResponse{Content: &resp_content}
	defer func() { c.JSON(http.StatusOK, resp.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _Device_SwitchHandler_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.DeviceID)
	xtext.RequireNonBlank(request.UserID)

	// 数据查询
	resp_content["data"] = "数据查询错误"
	dbmap := dbutil.GetDBMap()
	device_id := request.DeviceID
	pDevice, err := zndx.Device_Get(device_id, dbmap)
	xerr.ThrowPanic(err)

	// 保存数据
	onDevice_SwitchHandler(pDevice, cmd, &request, dbmap)

	// 查询设备命令并发送
	log.Println("-----------查询设备命令并发送")
	dev_cmds, err := dal.Query_DevicePowerCmdView_ByDeviceID(request.DeviceID, cmd, dbmap)
	xerr.ThrowPanic(err)
	for _, dev_cmd := range dev_cmds {
		coapcmd := coapclient.CoapCommand{
			HostAddr:    fmt.Sprintf("%s:%s", dev_cmd.RouterIP, dev_cmd.InRouteMappingPort),
			Method:      canopus.Post,
			RequestURI:  dev_cmd.RequestURI,
			QueryParams: map[string]string{"eui": pDevice.JoinNodeId},
			Payload:     dev_cmd.Payload}
		if reply := coapclient.Send(coapcmd); reply != nil {
			log.Printf("RevcPayload: %+v\n", reply.GetMessage())
		}
	}

	// 查询节点命令并发送
	log.Println("-----------查询设备命令并发送")
	node_cmds, err := dal.QueryList_NodeSwitchCmd_ByNode(&zndx.Node{Id: pDevice.PowerNodeId}, cmd, dbmap)
	xerr.ThrowPanic(err)
	for _, node_cmd := range node_cmds {
		coapcmd := coapclient.CoapCommand{
			HostAddr:    fmt.Sprintf("%s:%s", node_cmd.RouterIP, node_cmd.InRouterMappingPort),
			Method:      canopus.Post,
			RequestURI:  fmt.Sprintf("/smart_switch/switch_%s", pDevice.PowerSwitchId),
			QueryParams: map[string]string{"eui": pDevice.PowerNodeId},
			Payload:     payload}
		if reply := coapclient.Send(coapcmd); reply != nil {
			log.Printf("RevcPayload: %+v\n", reply.GetMessage())
		}
	}

	resp.Content = gin.H{"code": "0", "msg": "执行成功", "data": "ok"}
}

// 房间节点开关处理(By RoomID)
func RoomDevice_SwitchHandler(c *gin.Context, cmd, payload string) {
	defer xerr.CatchPanic()

	resp_content := gin.H{"code": "1", "msg": "无效参数", "data": ""}
	resp := xhttp.HttpResponse{Content: &resp_content}
	defer func() { c.JSON(http.StatusOK, resp.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _RoomDevice_SwitchHandler_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.RoomID)
	xtext.RequireNonBlank(request.UserID)

	// 数据查询
	resp_content["data"] = "数据查询错误"
	dbmap := dbutil.GetDBMap()

	//生成日志（操作日志和详细日志）
	nowTime := xtime.NowString()
	cmd_name := _CmdCode_ToName(cmd, "")
	opLog := zndx.DeviceOperateLog{0, nowTime, request.UserID, cmd, "classroom", request.RoomID, "", "", cmd, cmd_name, request.Params}
	opLog.Insert(dbmap)
	devLog := zndx.DeviceOnOffLogView{nowTime, request.UserID, "classroom", request.RoomID, cmd, cmd_name, request.Params}
	devLog.CreateLog(dbmap)

	dev_list, err := zndx.Device_QueryByRoomID(request.RoomID, dbmap)
	xerr.ThrowPanic(err)
	onRoomDevice_SwitchHandler(dev_list, cmd, dbmap)

	// 发送命令
	for _, dev := range dev_list {
		node_cmd, err := dal.Query_NodeSwitchCmd_ByDevice(&dev, cmd, dbmap)
		xdebug.DebugError(err)
		if err == nil {
			coapcmd := coapclient.CoapCommand{
				HostAddr:    fmt.Sprintf("%s:%s", node_cmd.RouterIP, node_cmd.InRouterMappingPort),
				Method:      canopus.Post,
				RequestURI:  fmt.Sprintf("/smart_switch/switch_%s", dev.PowerSwitchId),
				QueryParams: map[string]string{"eui": dev.PowerNodeId},
				Payload:     payload}
			go func() {
				if reply := coapclient.Send(coapcmd); reply != nil {
					log.Printf("RevcPayload: %+v\n", reply.GetMessage())
				}
			}()

			time.Sleep(time.Millisecond * 100)
		}
	}

	resp.Content = gin.H{"code": "0", "msg": "执行成功", "data": "ok"}
}

// 教室设备开关处理
func onRoomDevice_SwitchHandler(dev_list []zndx.Device, cmd string, dbmap *gorp.DbMap) {
	defer xerr.CatchPanic() // 捕获异常
	nowTime := xtime.NowString()

	// 处理设备的开/关
	for _, dev := range dev_list {
		//1) 更新设备自身状态
		device_id := dev.Id
		dev.DeviceSelfStatus = cmd
		err := dev.Update_DeviceSelfStatus(dbmap)
		xerr.ThrowPanic(err)

		//2) 更新设备开关记录
		exist := zndx.DeviceUseLog_Exists_DeviceOpenLog(device_id, dbmap)
		if cmd == "on" && !exist { // 开机操作 && 无开机记录
			err = zndx.DeviceUseLog_Insert_DeviceOpenLog(device_id, xtime.NowString(), dbmap)
			xerr.ThrowPanic(err)
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
}

// 设备开关数据处理
func onDevice_SwitchHandler(pDevice *zndx.Device, cmd string, pRequest *_Device_SwitchHandler_Request, dbmap *gorp.DbMap) {
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
		if (err == nil) && (pDeviceEx.IsAlert) && (pDeviceEx.MaxUseTime > 0) {
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

func init() {
	log.Print(fmt.Sprint(""))
}
