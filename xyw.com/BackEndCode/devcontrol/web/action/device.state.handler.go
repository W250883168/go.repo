package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"

	"dborm/zndx"
	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xtext"

	"dev.project/BackEndCode/devcontrol/app"
	"dev.project/BackEndCode/devcontrol/coap"
	"dev.project/BackEndCode/devcontrol/coap/coapview"
	"dev.project/BackEndCode/devcontrol/dal"
	"dev.project/BackEndCode/devcontrol/ioutil/dbutil"
	"dev.project/BackEndCode/devcontrol/model"
	"dev.project/BackEndCode/devcontrol/view/httpview"
)

// 设备状态查询请求（按设备ID查询）
type _Device_StateQuery_Request struct {
	UserID    string   // 用户ID
	DeviceIDs []string // 设备ID
	Params    string   // 参数（可选）
}

// 设备状态查询请求（按房间ID查询）
type _RoomDevice_StateQuery_Request struct {
	UserID string // 用户ID
	RoomID string // 房间ID
	Params string // 参数（可选）
}

// 设备状态
type DeviceStatus struct {
	StatusCode      string
	StatusName      string
	StatusValueCode string
	StatusValueName string
}

// 最后发送内容
type LastSendContent struct {
	CmdCode string
	Value   string
}

// 设备数据
type DeviceData struct {
	DeviceId        string
	DeviceName      string
	DeviceImg       string //关闭时显示的图片
	DeviceImg2      string //开启时显示的图片
	DevicePage      string
	NodeSwitch      string // on/off/offline(无应答时为offline,离线的意思)
	DeviceSwitch    string // on/off(无应答时为off)
	UseTimeBefore   int64  //上系统前已使用时间 单位：秒
	UseTimeAfter    int64  //上系统后累计使用时间 单位：秒
	IsCanUse        string //0-不可用 1-可用
	IsHaveAlert     string //0-没有预警消息 1-有预警消息
	DeviceStatus    []DeviceStatus
	LastSendContent []LastSendContent

	PowerMeterStat []coapview.PowerMeterView // 计量统计
}

// 设备状态结果
type _DevStatResult struct {
	Code string
	Name string
	Data []DeviceData
}

// 设备状态查询处理（按设备ID查询）
func Device_StateQuery_Handler(ptype string, c *gin.Context) {
	defer xerr.CatchPanic()

	reply := httpview.HttpResponseView{HttpCode: http.StatusOK, Content: gin.H{"code": "1", "msg": "无效参数", "data": ""}}
	defer func() { c.JSON(reply.HttpCode, reply.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _Device_StateQuery_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.UserID)

	// 查询数据
	dbmap := dbutil.GetDBMap()
	dev_data_list := []DeviceData{} // 定义存储返回结果的数据集
	for _, device_id := range request.DeviceIDs {
		dev_infos := dal.Query_DeviceBaseInfo(device_id, ptype, dbmap) // 查询设备基本信息
		for _, d := range dev_infos {
			var dev_data = DeviceData{
				DeviceId:        d.DeviceId,
				DeviceName:      d.DeviceName,
				DeviceImg:       d.DeviceImg,
				DeviceImg2:      d.DeviceImg2,
				DevicePage:      d.DevicePage,
				DeviceSwitch:    d.DeviceSelfStatus, //直接取设备自身状态（对设备执行开或关操作时，已经将状态写到了该字段）
				UseTimeBefore:   d.UseTimeBefore,
				UseTimeAfter:    d.UseTimeAfter,
				IsCanUse:        d.IsCanUse,
				IsHaveAlert:     d.IsHaveAlert,
				DeviceStatus:    []DeviceStatus{},
				NodeSwitch:      "offline", // 默认离线
				LastSendContent: []LastSendContent{}}

			// 获取设备最后一次发送的值(针对设备时才查询，如果是针对教室时不查询
			lastSendContents, _ := zndx.DeviceLastSendContent_Query_ByDevice(device_id, dbmap)
			for _, v := range lastSendContents {
				content := LastSendContent{v.CmdCode, v.LastSendContent}
				dev_data.LastSendContent = append(dev_data.LastSendContent, content)
			}

			// 1)首先检查给设备供电的节点是否在线
			var online1 = xtext.IsBlank(d.PowerNodeId) // 如果没有节点给设备供电，则默认为在线
			if !online1 && xtext.IsNotBlank(d.NodeSwitchStatusUpdateTime) {
				online1 = _IsOnline(d.NodeSwitchStatusUpdateTime, coap.OfflineTime) // 用更新时间判断节点是否在线
			}

			//2)其次检查设备接入节点（RS232资源、红外资源等）是否在线
			var online2 = xtext.IsBlank(d.JoinNodeId)
			if !online2 && xtext.IsNotBlank(d.JoinNodeUpdateTime) {
				online2 = _IsOnline(d.JoinNodeUpdateTime, coap.OfflineTime) //用更新时间判断节点是否在线
			}

			//3)最后计算最终设备是否在线
			if online1 && online2 {
				dev_data.NodeSwitch = "inline"
			}

			// 计算设备的工作参数状态值
			dev_data.DeviceStatus = _CalcDeviceStatusData(d, dbmap)
			dev_data_list = append(dev_data_list, dev_data)
		}
	}

	reply.Content = _DevStatResult{Code: "0", Name: "操作成功", Data: dev_data_list}
}

// 设备状态查询处理（按设备ID查询）
func Device_StateQuery_Handler2(ptype string, c *gin.Context) {
	defer xerr.CatchPanic()

	reply := httpview.HttpResponseView{HttpCode: http.StatusOK, Content: gin.H{"code": "1", "msg": "无效参数", "data": ""}}
	defer func() { c.JSON(reply.HttpCode, reply.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _Device_StateQuery_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.UserID)

	// 查询数据
	dbmap := dbutil.GetDBMap()
	dev_data_list := []DeviceData{} // 定义存储返回结果的数据集
	for _, device_id := range request.DeviceIDs {
		dev_infos := dal.Query_DeviceBaseInfo(device_id, ptype, dbmap) // 查询设备基本信息
		for _, d := range dev_infos {
			var dev_data = DeviceData{
				DeviceId:        d.DeviceId,
				DeviceName:      d.DeviceName,
				DeviceImg:       d.DeviceImg,
				DeviceImg2:      d.DeviceImg2,
				DevicePage:      d.DevicePage,
				DeviceSwitch:    d.DeviceSelfStatus, //直接取设备自身状态（对设备执行开或关操作时，已经将状态写到了该字段）
				UseTimeBefore:   d.UseTimeBefore,
				UseTimeAfter:    d.UseTimeAfter,
				IsCanUse:        d.IsCanUse,
				IsHaveAlert:     d.IsHaveAlert,
				NodeSwitch:      "offline", // 默认离线
				DeviceStatus:    []DeviceStatus{},
				LastSendContent: []LastSendContent{},
				PowerMeterStat:  []coapview.PowerMeterView{}}

			if pUpinfo, _ := zndx.NodeUpinfo_Get(d.PowerNodeId, dbmap); pUpinfo != nil {
				// 节点开关状态
				offTimeout := time.Duration(app.GetConfig().OffTimeout) * time.Second
				if pUpinfo.Onlined(offTimeout) {
					dev_data.NodeSwitch = "inline"
				}

				// 计量
				view_list, _ := coapview.NodeUpinfoView_GetListFrom(pUpinfo)
				for _, upinfo := range view_list {
					dev_data.PowerMeterStat = append(dev_data.PowerMeterStat, upinfo.Pe)
				}

				// 设备的工作参数状态值
				templates, err := zndx.DeviceStatusTemplate_QueryByModelID("device.model", dbmap)
				xdebug.LogError(err)
				if length := len(view_list); (err == nil) && (length > 0) {
					upinfo_view := view_list[length-1]
					status_arr := upinfo_view.Dev[0].Data.StringArray()
					for i, temp := range templates {
						dev_status := DeviceStatus{StatusCode: temp.StatusCode, StatusName: temp.StatusName, StatusValueCode: status_arr[i]}
						dev_data.DeviceStatus = append(dev_data.DeviceStatus, dev_status)
					}
				}
			}

			// 获取设备最后一次发送的值(针对设备时才查询，如果是针对教室时不查询
			lastSendContents, _ := zndx.DeviceLastSendContent_Query_ByDevice(device_id, dbmap)
			for _, v := range lastSendContents {
				dev_data.LastSendContent = append(dev_data.LastSendContent, LastSendContent{v.CmdCode, v.LastSendContent})
			}

			dev_data_list = append(dev_data_list, dev_data)
		}
	}

	reply.Content = _DevStatResult{Code: "0", Name: "操作成功", Data: dev_data_list}
}

// 设备状态查询处理（按房间ID查询）
func RoomDevice_StateQuery_Handler(ptype string, c *gin.Context) {
	defer xerr.CatchPanic()

	reply := httpview.HttpResponseView{HttpCode: http.StatusOK, Content: gin.H{"code": "1", "msg": "无效参数", "data": ""}}
	defer func() { c.JSON(reply.HttpCode, reply.Content) }()

	// 解析参数
	data, _ := ioutil.ReadAll(c.Request.Body)
	var request _RoomDevice_StateQuery_Request
	err := json.Unmarshal(data, &request)
	xerr.ThrowPanic(err)

	// 校验参数
	xtext.RequireNonBlank(request.UserID)
	xtext.RequireNonBlank(request.RoomID)

	// 查询数据
	dbmap := dbutil.GetDBMap()
	dev_data_list := []DeviceData{}                                     // 定义存储返回结果的数据集
	dev_infos := dal.Query_DeviceBaseInfo(request.RoomID, ptype, dbmap) // 查询设备基本信息
	for _, d := range dev_infos {
		var dev_data = DeviceData{
			DeviceId:        d.DeviceId,
			DeviceName:      d.DeviceName,
			DeviceImg:       d.DeviceImg,
			DeviceImg2:      d.DeviceImg2,
			DevicePage:      d.DevicePage,
			DeviceSwitch:    "off", // 默认离线
			UseTimeBefore:   d.UseTimeBefore,
			UseTimeAfter:    d.UseTimeAfter,
			IsCanUse:        d.IsCanUse,
			IsHaveAlert:     d.IsHaveAlert,
			DeviceStatus:    []DeviceStatus{},
			NodeSwitch:      "offline", // 默认离线
			LastSendContent: []LastSendContent{}}

		// 直接取设备自身状态（对设备执行开或关操作时，已经将状态写到了该字段）
		if xtext.IsNotBlank(d.DeviceSelfStatus) {
			dev_data.DeviceSwitch = d.DeviceSelfStatus
		}

		// 1)首先检查给设备供电的节点是否在线
		var online1 = xtext.IsBlank(d.PowerNodeId) // 如果没有节点给设备供电，则默认为在线
		if !online1 && xtext.IsNotBlank(d.NodeSwitchStatusUpdateTime) {
			online1 = _IsOnline(d.NodeSwitchStatusUpdateTime, coap.OfflineTime) // 用更新时间判断节点是否在线
		}

		//2)其次检查设备接入节点（RS232资源、红外资源等）是否在线
		var online2 = xtext.IsBlank(d.JoinNodeId)
		if !online2 && xtext.IsNotBlank(d.JoinNodeUpdateTime) {
			online2 = _IsOnline(d.JoinNodeUpdateTime, coap.OfflineTime) //用更新时间判断节点是否在线
		}

		//3)最后计算最终设备是否在线
		if online1 && online2 {
			dev_data.NodeSwitch = "inline"
		}

		// 计算设备的工作参数状态值
		dev_data.DeviceStatus = _CalcDeviceStatusData(d, dbmap)
		dev_data_list = append(dev_data_list, dev_data)
	}

	reply.Content = _DevStatResult{Code: "0", Name: "操作成功", Data: dev_data_list}
}

//计算设备的当前工作状态（通过node接入的设备）
func _CalcDeviceStatusData(d model.DeviceBasciInfoView, dbmap *gorp.DbMap) (list []DeviceStatus) {
	defer xerr.CatchPanic()
	list = []DeviceStatus{}

	//查询设备状态配置记录
	stat_cmds, err := zndx.DeviceModelStatusCmd_Query_ByDevice(&zndx.Device{JoinNodeId: d.JoinNodeId, JoinSocketId: d.JoinSocketId}, dbmap)
	xerr.ThrowPanic(err)

	//查询设备状态上报记录（节点每10秒上报一次）
	//定义变量存储要返回前端的设备状态态（一个设备可以有多个状态，所以用数组形式）
	pStatus := &zndx.NodeSocketStatus{NodeId: d.JoinNodeId, SocketId: d.JoinSocketId}
	nodeSockets, err := pStatus.Query_ByNodeInfo(dbmap)
	xerr.ThrowPanic(err)

	//对设备配置的状态命令进行循环,依次取出各个状态的值-------------------------------
	for i, stat_cmd := range stat_cmds {
		//保存设备状态值
		var ds = DeviceStatus{StatusCode: stat_cmd.StatusCode, StatusName: stat_cmd.StatusName}

		//从上报的数据中找到当前状态对应的状态值
		var nodeSocket zndx.NodeSocketStatus
		if len(nodeSockets) > i {
			nodeSocket = nodeSockets[i]       //取出状态值
			if nodeSocket.StatusValue != "" { //如果状态值不等于空
				if _IsOnline(nodeSocket.UpdateTime, coap.OfflineTime) { //如果是最新上报的数据(在线）
					//获取状态值代码
					ds.StatusValueCode = _MatchVarValue(stat_cmd.StatusValueMatchString, nodeSocket.StatusValue, "{val}") //1)解析并保存状态值

					//获取状态值名称
					ds.StatusValueName = ds.StatusValueCode
					if stat_cmd.SelectValueFlag == "1" && ds.StatusValueCode != "" {
						pValueCode := &zndx.DeviceModelStatusValueCode{ModelId: stat_cmd.ModelId, StatusCode: stat_cmd.StatusCode, StatusValueCode: ds.StatusValueCode}
						ds.StatusValueName, _ = pValueCode.Query_StatusValueName(dbmap)
					}
				}
			}
		}

		list = append(list, ds)
	}

	return list
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
