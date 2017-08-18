package coapserver

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"canopus"

	"xutils/xcrypto/xhash"
	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xtime"

	"dborm/zndx"
	"dborm/zndxview"

	"dev.project/BackEndCode/devcontrol/app"
	"dev.project/BackEndCode/devcontrol/coap"
	"dev.project/BackEndCode/devcontrol/coap/coapclient"
	"dev.project/BackEndCode/devcontrol/coap/coapview"
	"dev.project/BackEndCode/devcontrol/dal"
	"dev.project/BackEndCode/devcontrol/ioutil/dbutil"
)

const _OTA_BUFF_SIZE = 256

// 启动Coap服务
func StartCoapService() {
	server := canopus.NewCoapServer("5683")

	//设置请求URI
	// server.Post("/hello", onHelloHandler2)
	server.Post("/hello", onHello)
	server.Get("/ota", onOTAHandler2)
	server.Post("/test", onTestHandler)

	server.OnError(onErrorHandler)
	server.OnClose(onCloseHandler)

	server.Start()
}

// [/hello]资源处理
func onHello(request canopus.CoapRequest) (response canopus.CoapResponse) {
	defer xerr.CatchPanic()

	// 无需应答
	pResponseMsg := canopus.NewMessageOfType(canopus.MessageNonConfirmable, request.GetMessage().MessageID)
	payload := coapview.CoapResponseView{Msg: "excuting..."}
	pResponseMsg.SetStringPayload(payload.ToJson())
	response = canopus.NewResponse(pResponseMsg, nil)

	// 解析数据
	pRequestMsg := request.GetMessage()
	data := pRequestMsg.Payload.GetBytes()
	txt := fmt.Sprintf("RECV : \n\t%s", string(data))
	log.Println(txt)
	var request_info coapview.NodeUpinfo_Request
	err := json.Unmarshal(data, &request_info)
	xerr.ThrowPanic(err)

	// 更新节点地址
	tNow := time.Now()
	pDBMap := dbutil.GetDBMap()
	node := zndx.Node{
		Id:                 request_info.Node.Eui,
		RouteIp:            request.GetAddress().IP.String(),
		InRouteMappingPort: strconv.Itoa(request.GetAddress().Port),
		UploadTime:         xtime.TimeString(&tNow),
		IpType:             "ipv4",
		NodeCoapPort:       "5683"}
	err = node.SavePingInfo(pDBMap)
	xerr.ThrowPanic(err)

	// 节点日志
	data, _ = json.Marshal(&request_info)
	logg := zndx.NodeLog{
		LogTime:    tNow,
		LogTopic:   pRequestMsg.GetURIPath(),
		LogContent: string(data)}
	logg.Insert(pDBMap) // 插入日志，忽略错误

	// 节点上报信息
	upinfo_view := coapview.NodeUpinfoView{
		Uptime:             xtime.TimeString(&tNow),
		NodeAddr:           fmt.Sprintf("%s:%d", request.GetAddress().IP.String(), request.GetAddress().Port),
		NodeUpinfo_Request: request_info}
	err = upinfo_view.Save(pDBMap)
	xerr.ThrowPanic(err)

	// 累计时间及预警
	dev_list, err := zndx.Device_QueryByNodeID(request_info.NodeID(), dbutil.GetDBMap())
	xdebug.LogError(err)
	for _, dev := range dev_list {
		// 生成设备累计使用时间
		var pDevProp *zndx.DeviceProp
		prop_key, comments := "device.accumulate.time.seconds", "设备累计使用时间"
		if !zndx.DeviceProp_Exists(dev.Id, prop_key, pDBMap) {
			pDevProp = &zndx.DeviceProp{DeviceID: dev.Id, K: prop_key, V: "0", Comments: comments}
			err = pDevProp.Insert(pDBMap)
			xerr.ThrowPanic(err)
			continue
		}

		pDevProp, err = zndx.DeviceProp_Get(dev.Id, prop_key, pDBMap)
		xerr.ThrowPanic(err)
		acc, _ := strconv.Atoi(pDevProp.V)
		acc += app.GetConfig().CoapConfig.HeartbeatInterval
		pDevProp.V = fmt.Sprintf("%d", acc)
		_, err = pDevProp.Update(pDBMap)
		xerr.ThrowPanic(err)

		// 生成设备实时预警记录
		if pDeviceEx, err := zndxview.Query_DeviceExNodeView(dev.Id, pDBMap); err == nil {
			pDevAlert := zndx.DeviceAlert{DeviceId: dev.Id}
			err = pDevAlert.Delete_DeviceOvertime(dev.Id, pDBMap)
			xerr.ThrowPanic(err)
			time_used := int(pDeviceEx.UseTimeBefore) + acc
			if (pDeviceEx.MaxUseTime > 0) && (time_used > int(pDeviceEx.MaxUseTime)) {
				tspan_used_total := xtime.TimeSpan(time.Second * time.Duration(time_used))
				tspan_used_max := xtime.TimeSpan(time.Second * time.Duration(pDeviceEx.MaxUseTime))
				tspan_format := "[%0.2d小时%0.2d分钟%0.2d秒]"
				alert_desc := "累计使用时间" + tspan_used_total.ToString(tspan_format) + "超过了设定的最大使用时间" + tspan_used_max.ToString(tspan_format)
				pDevAlert.AlertType = "1" // 超时使用预警
				pDevAlert.AlertDescription = alert_desc
				pDevAlert.LastAlertTime = xtime.TimeString(&tNow)
				pDevAlert.Insert(pDBMap)
				xerr.ThrowPanic(err)
			}
		}
	}

	// 广播处理
	if app.GetConfig().BroadcastingSupport && upinfo_view.Node.IsBroadcast() {
		log.Println("\t 节点设备广播消息处理,,,,,,,,,")
		onHelloBroadcasting(request, &request_info)
		return response
	}

	// 更新设备供电的节点开关状态
	for _, s := range upinfo_view.Sw {
		dev := zndx.Device{NodeSwitchStatus: s.Stat, NodeSwitchStatusUpdateTime: xtime.TimeString(&tNow), PowerNodeId: upinfo_view.Node.Eui, PowerSwitchId: s.Id}
		err = dev.Update_NodeSwitchStatus(pDBMap)
		xerr.ThrowPanic(err)
	}

	// 更新接入设备的最后上报时间(2016-10-17日增加）
	dev := zndx.Device{JoinNodeUpdateTime: xtime.TimeString(&tNow), JoinNodeId: upinfo_view.Node.Eui}
	err = dev.Update_JoinNodeUpdateTime(pDBMap)
	xerr.ThrowPanic(err)

	// 更新设备各个状态
	var nodeSocketState = zndx.NodeSocketStatus{NodeId: upinfo_view.Node.Eui}
	err = nodeSocketState.DeleteByNodeID(pDBMap)
	xerr.ThrowPanic(err)
	for _, s := range upinfo_view.Dev { // i代表了设备连接节点的插口（JoinSocketId）
		//保存设备的各个状态值到数据库-------------------------------------------------------------
		var statusValueArray [10]string = s.Data.StringArray() // 将10个状态值放入数组(目前节点只支持10个状态值)
		for no, v := range statusValueArray {
			socket_index, _ := strconv.Atoi(s.Id)
			socket_id := fmt.Sprintf("%d", socket_index+1)
			var status = zndx.NodeSocketStatus{NodeId: upinfo_view.Node.Eui, SocketId: socket_id, SeqNo: no + 1, StatusValue: v, UpdateTime: xtime.TimeString(&tNow)}
			err = status.Insert(pDBMap)
			xerr.ThrowPanic(err)
		}
	}

	//报警处理-------------------------------------------------------------------------------
	for i, s := range upinfo_view.Dev { // i代表了设备连接节点的插口（JoinSocketId）
		joinSocketId := strconv.Itoa(i + 1)
		var statusValueArray [10]string = s.Data.StringArray() // 将10个状态值放入数组(目前节点只支持10个状态值)
		status_cmds, err := zndx.DeviceModelStatusCmd_Query_ByDevice(&zndx.Device{JoinNodeId: upinfo_view.Node.Eui, JoinSocketId: joinSocketId}, pDBMap)
		xdebug.LogError(err)
		for _, statusCmd := range status_cmds {
			//根据序号(SeqNo)找到对应的状态值
			statusValue := ""
			if statusCmd.SeqNo >= 1 {
				statusValue = statusValueArray[statusCmd.SeqNo-1]
			}
			//根据匹配串从开关状态值中匹配出状态值代码
			statusValueCode := ""
			if statusValue != "" {
				statusValueCode = _MatchVarValue(statusCmd.StatusValueMatchString, statusValue, "{val}") //1)解析并保存状态值
			}

			//预警处理(待调试)---------------------------------------------------------------------------------------------
			//1)按状态值条件进行预警
			dal.Refresh_DeviceAlert_ByStatusCode(upinfo_view.Node.Eui, joinSocketId, statusCmd, statusValueCode, pDBMap)
			//2)按具体的状态值预警
			dal.Refresh_DeviceAlert_ByStatusValue(upinfo_view.Node.Eui, joinSocketId, statusCmd, statusValueCode, pDBMap)
		}
	}

	payload.Msg = "ok"
	pResponseMsg.SetStringPayload(payload.ToJson())
	return response
}

// 广播处理
func onHelloBroadcasting(request canopus.CoapRequest, pUpinfo *coapview.NodeUpinfo_Request) {
	if isServerAddr := (pUpinfo.Node.S_addr == app.GetConfig().ThisHostAddr); isServerAddr {
		// 停止广播
		go func() {
			log.Println("\t		 Stop Broadcasting,,,,,,,")
			cmd := coapclient.CoapCommand{
				HostAddr:    request.GetAddress().String(),
				Method:      canopus.Post,
				RequestURI:  "/eui_broadcast",
				QueryParams: map[string]string{"eui": pUpinfo.Node.Eui},
				Payload:     "STOP"}
			coapclient.Send(cmd)
		}()
	} else {
		go func() {
			// 设置服务器IP地址
			log.Println("\t		 Set ServerIP")
			cmd := coapclient.CoapCommand{
				HostAddr:    request.GetAddress().String(),
				Method:      canopus.Post,
				RequestURI:  "/infor",
				QueryParams: map[string]string{"eui": pUpinfo.Node.Eui, "param": "s_addr"},
				Payload:     app.GetConfig().ThisHostAddr}
			coapclient.Send(cmd)
			time.Sleep(time.Millisecond * 10) // 两次命令间隔10毫秒

			// Save Changes
			log.Println("\t		Save Changes")
			save_cmd := coapclient.CoapCommand{
				HostAddr:    request.GetAddress().String(),
				Method:      canopus.Post,
				RequestURI:  "/infor",
				QueryParams: map[string]string{"eui": pUpinfo.Node.Eui, "param": "save"},
				Payload:     "1"}
			coapclient.Send(save_cmd)
		}()
	}
}

// Deprecated
func onHelloHandler2(request canopus.CoapRequest) (response canopus.CoapResponse) {
	defer xerr.CatchPanic()
	// log.Println("###*****************", request.GetMessage().GetURIPath())

	msg := canopus.NewMessageOfType(canopus.MessageNonConfirmable, request.GetMessage().MessageID) //不需应答，所以不能使用canopus.MessageAcknowledgment
	msg.SetStringPayload("ok")
	response = canopus.NewResponse(msg, nil)

	// 解析数据
	data := request.GetMessage().Payload.GetBytes()
	log.Printf("RECV : \n\t%s", string(data))
	var info coapview.DevPayloadView
	err := json.Unmarshal(data, &info)
	xerr.ThrowPanic(err)

	// 数据检查
	log.Printf("eui=%s; msgtype=%d; msgid=%d\n\n", info.Node.Eui, request.GetMessage().MessageType, request.GetMessage().MessageID)
	_ValidateMessage(request.GetMessage())
	info.Node.Validate()

	// 节点设备广播消息处理
	//	if info.Node.IsBroadcast() {
	//		log.Println("节点设备广播消息处理,,,,,,,,,")
	//		_HandleBroadcast(request, &info)
	//		return response
	//	}

	// Redis 缓存
	//	var client *redis.Client
	//	if client, err = redisutil.DefaultClient(); err == nil {
	//		input_info := coapview.NodeStatInfo{DevPayloadView: info, NodeAddr: request.GetAddress().String(), When: xtime.NowString()}
	//		defer func() {
	//			if err == nil {
	//				k := input_info.DevPayloadView.Node.Eui
	//				v := input_info.JsonText()
	//				expire := time.Duration(coapclient.OffTimeout()) * time.Second
	//				client.Set(k, v, expire)
	//			}
	//		}()

	//		if txt := client.Get(info.Node.Eui).Val(); xtext.IsNotBlank(txt) {
	//			var exist_info coapview.NodeStatInfo
	//			json.Unmarshal([]byte(txt), &exist_info)
	//			if exist_info.NodeAddr == request.GetAddress().String() {
	//				return response
	//			}
	//		}
	//	}

	dbmap := dbutil.GetDBMap()
	tNow := xtime.NowString()
	node := zndx.Node{
		Id:                 info.Node.Eui,
		InRouteMappingPort: strconv.Itoa(request.GetAddress().Port),
		RouteIp:            request.GetAddress().IP.String(),
		UploadTime:         tNow,
		IpType:             "ipv4",
		NodeCoapPort:       "5683"}
	// log.Println(node)
	err = node.SavePingInfo(dbmap)
	xerr.ThrowPanic(err)

	// 更新设备供电的节点开关状态
	for _, s := range info.Sw {
		dev := zndx.Device{NodeSwitchStatus: s.Stat, NodeSwitchStatusUpdateTime: tNow, PowerNodeId: info.Node.Eui, PowerSwitchId: s.Id}
		err = dev.Update_NodeSwitchStatus(dbmap)
		xerr.ThrowPanic(err)
	}

	// 更新接入设备的最后上报时间(2016-10-17日增加）
	dev := zndx.Device{JoinNodeUpdateTime: tNow, JoinNodeId: info.Node.Eui}
	err = dev.Update_JoinNodeUpdateTime(dbmap)
	xerr.ThrowPanic(err)

	// 更新设备各个状态
	var nodeSocketState = zndx.NodeSocketStatus{NodeId: info.Node.Eui}
	err = nodeSocketState.DeleteByNodeID(dbmap)
	xerr.ThrowPanic(err)

	for i, s := range info.Dev { // i代表了设备连接节点的插口（JoinSocketId）
		//保存设备的各个状态值到数据库-------------------------------------------------------------
		var statusValueArray [10]string = s.Data.StringArray() // 将10个状态值放入数组(目前,节点只支持10个状态值)
		for no, v := range statusValueArray {
			socket_index, _ := strconv.Atoi(s.Id)
			socket_id := fmt.Sprintf("%d", socket_index+1)
			var status = zndx.NodeSocketStatus{NodeId: info.Node.Eui, SocketId: socket_id, SeqNo: no + 1, StatusValue: v, UpdateTime: tNow}
			err = status.Insert(dbmap)
			xerr.ThrowPanic(err)
		}

		//报警处理-------------------------------------------------------------------------------
		joinSocketId := strconv.Itoa(i + 1)
		status_cmds, err := zndx.DeviceModelStatusCmd_Query_ByDevice(&zndx.Device{JoinNodeId: info.Node.Eui, JoinSocketId: joinSocketId}, dbmap)
		xdebug.LogError(err)
		for _, statusCmd := range status_cmds {
			//根据序号(SeqNo)找到对应的状态值
			statusValue := ""
			if statusCmd.SeqNo >= 1 {
				statusValue = statusValueArray[statusCmd.SeqNo-1]
			}
			//根据匹配串从开关状态值中匹配出状态值代码
			statusValueCode := ""
			if statusValue != "" {
				statusValueCode = _MatchVarValue(statusCmd.StatusValueMatchString, statusValue, "{val}") //1)解析并保存状态值
			}

			//预警处理(待调试)---------------------------------------------------------------------------------------------
			//1)按状态值条件进行预警
			dal.Refresh_DeviceAlert_ByStatusCode(info.Node.Eui, joinSocketId, statusCmd, statusValueCode, dbmap)
			//2)按具体的状态值预警
			dal.Refresh_DeviceAlert_ByStatusValue(info.Node.Eui, joinSocketId, statusCmd, statusValueCode, dbmap)
		}
	}

	return response
}

// Deprecated
func _HandleBroadcast(request canopus.CoapRequest, info *coapview.DevPayloadView) {
	if isServerAddr := (info.Node.S_addr == app.GetConfig().ThisHostAddr); isServerAddr {
		// Stop Broadcasting
		go func() {
			log.Println("\t		 Stop Broadcasting,,,,,,,")
			cmd := coapclient.CoapCommand{
				HostAddr:    request.GetAddress().String(),
				Method:      canopus.Post,
				RequestURI:  "/eui_broadcast",
				QueryParams: map[string]string{"eui": info.Node.Eui},
				Payload:     "STOP"}
			coapclient.Send(cmd)
		}()
	} else {
		go func() {
			// Set ServerIP
			log.Println("\t		 Set ServerIP")
			cmd := coapclient.CoapCommand{
				HostAddr:    request.GetAddress().String(),
				Method:      canopus.Post,
				RequestURI:  "/infor",
				QueryParams: map[string]string{"eui": info.Node.Eui, "param": "s_addr"},
				Payload:     app.GetConfig().ThisHostAddr}
			coapclient.Send(cmd)

			time.Sleep(time.Millisecond * 10) // 两次命令间隔10毫秒

			// Save Changes
			log.Println("\t		Save Changes")
			save_cmd := coapclient.CoapCommand{
				HostAddr:    request.GetAddress().String(),
				Method:      canopus.Post,
				RequestURI:  "/infor",
				QueryParams: map[string]string{"eui": info.Node.Eui, "param": "save"},
				Payload:     "1"}
			coapclient.Send(save_cmd)
		}()
	}
}

func onCloseHandler(cs canopus.CoapServer) {
	log.Println("CoapServer OnClose ")
}

func onErrorHandler(err error) {
	log.Println("CoapServer OnError: ", err.Error())
}

// ota升级处理
func onOTAHandler2(req canopus.CoapRequest) (resp canopus.CoapResponse) {
	defer xerr.CatchPanic()
	log.Printf("收到ota请求：%+v	\nPayloadBuff: %s", req.GetMessage(), req.GetMessage().Payload.String())
	msg := canopus.NewMessageOfType(canopus.MessageNonConfirmable, req.GetMessage().MessageID)
	resp = canopus.NewResponse(msg, nil)

	// 解析数据
	pRequestMsg := req.GetMessage()
	payload_buff := pRequestMsg.Payload.GetBytes()
	log.Println("Payload.DATA: ", string(payload_buff))
	payload_struct := struct{ Offset, Eui, Version string }{}
	err := json.Unmarshal(payload_buff, &payload_struct)
	xerr.ThrowPanic(err)
	log.Printf("%+v", payload_struct)

	// 打开文件
	bin_file := payload_struct.Version
	offset, _ := strconv.Atoi(payload_struct.Offset)
	f, err := os.Open(bin_file)
	xerr.ThrowPanic(err)
	defer f.Close()
	finfo, err := f.Stat()
	xerr.ThrowPanic(err)

	// 检查是否超出文件大小
	if int64(offset) >= finfo.Size() {
		msg.SetStringPayload("eof")
		log.Println("Offset Over Filesize,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,")
		return resp
	}

	// 读数据
	buff := make([]byte, _OTA_BUFF_SIZE)
	log.Printf("RequestFileOffset:  %d", offset)
	reads, err := f.ReadAt(buff, int64(offset))
	if (reads < _OTA_BUFF_SIZE) && (err.Error() == io.EOF.Error()) {
		log.Printf("Last Reads:= %d", reads)
	} else {
		xerr.ThrowPanic(err)
	}

	// 返回数据
	buff_str := string(buff[:reads])
	msg.SetStringPayload(buff_str)

	// 进度
	data_transferred := int(offset) + int(reads)
	progress := fmt.Sprintf("%.1f", 100*float64(data_transferred)/float64(finfo.Size()))
	coap.GetOTAContext().PutValue(payload_struct.Eui, progress)
	log.Printf("%X", buff)
	log.Println("MD5/Checksum:  ", xhash.ToMD5(buff_str))
	log.Printf("Progressing... %s%%", progress)
	return resp
}

// ota升级
func onOTAHandler(req canopus.CoapRequest) (resp canopus.CoapResponse) {
	defer xerr.CatchPanic()
	msg := canopus.NewMessageOfType(canopus.MessageNonConfirmable, req.GetMessage().MessageID)
	resp = canopus.NewResponse(msg, nil)

	// 解析数据
	pRequestMsg := req.GetMessage()
	payload_buff := pRequestMsg.Payload.GetBytes()
	le_offset := binary.LittleEndian.Uint32(payload_buff)
	log.Println(fmt.Sprintf("RECV %d Bytes, Offset=%d", len(payload_buff), le_offset))

	// 读文件
	node_eui := "00124b000cd52302"
	bin_file := "ota-image.bin"
	f, err := os.Open(bin_file)
	xerr.ThrowPanic(err)
	defer f.Close()
	finfo, err := f.Stat()
	xerr.ThrowPanic(err)

	offset := le_offset
	if int64(offset) >= finfo.Size() {
		msg.SetStringPayload("eof")
		log.Println("Offset Over Filesize,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,")
		return resp
	}

	buff := make([]byte, _OTA_BUFF_SIZE)
	log.Printf("RequestFileOffset:  %d", offset)
	reads, err := f.ReadAt(buff, int64(offset))
	if (reads < _OTA_BUFF_SIZE) && (err.Error() == io.EOF.Error()) {
		log.Printf("Last Reads:= %d", reads)
	} else {
		xerr.ThrowPanic(err)
	}

	buff_str := string(buff[:reads])
	msg.SetStringPayload(buff_str)

	data_transferred := int(offset) + int(reads)
	log.Printf("%X", buff)
	checksum := xhash.ToMD5(buff_str)
	log.Println("MD5/Checksum:  ", checksum)
	progress := fmt.Sprintf("%3.1f", 100*float64(data_transferred)/float64(finfo.Size()))
	log.Printf("Progressing... %s%%", progress)
	coap.GetOTAContext().PutValue(node_eui, progress)

	log.Println(coap.GetOTAContext().KValue)
	return resp
}

func onTestHandler(req canopus.CoapRequest) (resp canopus.CoapResponse) {
	defer xerr.CatchPanic()

	// 解析数据
	pRequestMsg := req.GetMessage()
	bytes_received := pRequestMsg.Payload.GetBytes()
	payload_txt := string(bytes_received)
	log.Println(fmt.Sprintf("RECV : \n\t%s", payload_txt))

	msg := canopus.NewMessageOfType(canopus.MessageNonConfirmable, req.GetMessage().MessageID)
	msg.SetStringPayload(payload_txt)
	resp = canopus.NewResponse(msg, nil)
	return resp
}

// 不处理消息类型为0/CON（需要应答）的消息
func _ValidateMessage(pMessage *canopus.Message) {
	if pMessage.MessageType == canopus.MessageConfirmable {
		panic(errors.New("不处理消息类型为0/CON（需要应答）的消息"))
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

}
