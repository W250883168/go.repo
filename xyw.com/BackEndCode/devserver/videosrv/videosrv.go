package videosrv

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"strconv"

	"dev.project/BackEndCode/devserver/commons/xdebug"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/equipment"
)

// 服务器设置
const gSERVER_IP string = "127.0.0.1"
const gSERVER_PORT int = 1616
const gURL_PATH string = "/vod"
const gNET_PROTOCOL string = "udp4"

//数据缓冲区
const gBUFFSIZE = 1024 * 2

var gBuff = make([]byte, gBUFFSIZE)
var pConn *net.UDPConn
var bStopped bool = true

func StartService() {
	//本地监听udp地址
	addr := gSERVER_IP + ":" + strconv.Itoa(gSERVER_PORT)
	udpLocalAddr, err := net.ResolveUDPAddr(gNET_PROTOCOL, addr)
	xdebug.HandleError(err)

	//本地监听udp连接
	updConn, err := net.ListenUDP(gNET_PROTOCOL, udpLocalAddr)
	xdebug.HandleError(err)

	//接收消息
	pConn = updConn
	bStopped = false
	go HandleMessage2(pConn)
}

func StopService() {
	if !IsServiceStopped() {
		pConn.Close()
	}

	bStopped = true
}

func IsServiceStopped() bool {
	return bStopped
}

//发送消息
func SendCommand(cmd equipment.CommandSendlog) (ok bool) {
	defer xdebug.DoRecover()
	addr := cmd.CmdIp + ":" + strconv.Itoa(cmd.CmdPort) // udp地址
	udpAddr, err := net.ResolveUDPAddr(gNET_PROTOCOL, addr)
	xdebug.HandleError(err)

	conn, err := net.DialUDP(gNET_PROTOCOL, nil, udpAddr)
	defer conn.Close()
	xdebug.HandleError(err)

	data, _ := json.Marshal(&cmd)
	var msg = CmdMessage{
		CmdType:  reflect.TypeOf(equipment.CommandSendlog{}).Name(),
		JsonText: string(data)}

	data, _ = json.Marshal(&msg)
	_, err = conn.Write(data)
	xdebug.HandleError(err)

	ok = true
	return ok
}

func SendCmdMessage(msg CmdMessage, addr string) (ok bool) {
	defer xdebug.DoRecover()
	udpAddr, err := net.ResolveUDPAddr(gNET_PROTOCOL, addr)
	xdebug.HandleError(err)

	conn, err := net.DialUDP(gNET_PROTOCOL, nil, udpAddr)
	defer conn.Close()
	xdebug.HandleError(err)

	data, _ := json.Marshal(msg)
	fmt.Println(addr)
	fmt.Println(string(data))

	_, err = conn.Write(data)
	xdebug.HandleError(err)

	ok = true
	return ok
}

//接收消息
func HandleMessage2(conn *net.UDPConn) {
	for !IsServiceStopped() {
		n, _, err := conn.ReadFromUDP(gBuff)
		xdebug.DebugError(err)
		if n > 0 {
			fmt.Println("REVC: " + string(gBuff[:n]))
			var msg CmdMessage
			if err := json.Unmarshal(gBuff[:n], &msg); err != nil {
				xdebug.DebugError(err)
				continue
			}

			switch msg.CmdType {
			case "CommandSendlog":
				handle_CommandSendlog(msg.JsonText)
			case "VideoCaptureCommand":
				handle_VideoCaptureCommand(msg.JsonText)
			default:
			}
		}
	}
}

func handle_VideoCaptureCommand(jonsText string) {

}

func handle_CommandSendlog(jsonText string) {
	defer xdebug.DoRecover()
	var cmdlog equipment.CommandSendlog
	err := json.Unmarshal([]byte(jsonText), &cmdlog)
	xdebug.HandleError(err)

	if cmdlog.CmdType == "0" { // FFmpegServer命令
		var vConfig = equipment.VideoConfig{
			CameraIp:    cmdlog.CmdIp,
			CameraPort:  cmdlog.CmdPort,
			CameraState: cmdlog.CmdState,
			Classroomid: cmdlog.Classroomid}

		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		dbmap.AddTableWithName(equipment.VideoConfig{}, "videoconfig").SetKeys(true, "Id")
		sql := `UPDATE videoconfig SET CameraIp = ?, CameraPort = ?, CameraState = ? WHERE (Classroomid = ?);`
		_, err := dbmap.Exec(sql, vConfig.CameraIp, vConfig.CameraPort, vConfig.CameraState, vConfig.Classroomid)
		xdebug.HandleError(err)
	}

	if err == nil {
		cmdlog.CmdStr = "ok"
		SendCommand(cmdlog)
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
