package action

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xtime"

	// "vodx/app"
	"vodx/ioutil/dbutil"
	"vodx/mqclient"
	"vodx/mqclient/mqview"
	"vodx/orm/dborm"
	"vodx/web/view/httpview"
)

// 开始视频录制
func BeginVideo(pRequest *httpview.VideoCaptureRequest) {
	// 打开MQ连接
	pClient, err := mqclient.OpenDefault() // mqclient.Connect(app.GetConfig().MQConnString)
	xerr.ThrowPanic(err)
	defer pClient.Close()

	// 查询摄像机信息
	pDBMap := dbutil.GetDBMap()
	var camera_list []dborm.CameraInfo
	camera_list, err = dborm.CameraInfo_QueryByLocation(pRequest.LocationID, pDBMap)
	xerr.ThrowPanic(err)
	for _, camera := range camera_list {
		// 发送
		msg := mqview.MQMessage{
			MessageID: time.Now().Format(xtime.FORMAT_yyyyMMddHHmmssfff),
			Message: struct {
				dborm.CameraInfo
				Action string // run/stop/pause
			}{CameraInfo: camera, Action: "run"}}
		mqclient.SendMessage("ffmpeg.cmd.capture.camera", msg, pClient)

		// 日志
		buff, _ := json.Marshal(msg)
		log.Println(string(buff))
		log := dborm.OperateLog{LogDate: time.Now(), LogContent: string(buff)}
		err = log.Insert(pDBMap)
		xdebug.LogError(err)
	}

	// 查询教学电脑信息
	var pc_list []dborm.ComputerInfo
	pc_list, err = dborm.ComputerInfo_QueryByLocation(pRequest.LocationID, pDBMap)
	xerr.ThrowPanic(err)
	for _, pc := range pc_list {
		// 发送
		msg := mqview.MQMessage{
			MessageID: time.Now().Format(xtime.FORMAT_yyyyMMddHHmmssfff),
			Message: struct {
				dborm.ComputerInfo
				Action string // run/stop/pause
			}{ComputerInfo: pc, Action: "run"}}
		mqclient.SendMessage("ffmpeg.cmd.capture.screen", msg, pClient)

		// 日志
		buff, _ := json.Marshal(msg)
		log.Println(string(buff))
		log := dborm.OperateLog{LogDate: time.Now(), LogContent: string(buff)}
		err = log.Insert(pDBMap)
		xdebug.LogError(err)
	}

}

// 停止视频录制
func EndVideo(pRequest *httpview.VideoCaptureRequest) {
	pDBMap := dbutil.GetDBMap()

	var err error
	var list []dborm.CameraInfo
	list, err = dborm.CameraInfo_QueryByLocation(pRequest.LocationID, pDBMap)
	xerr.ThrowPanic(err)

	pClient, err := mqclient.OpenDefault() // mqclient.Connect(app.GetConfig().MQConnString)
	xerr.ThrowPanic(err)
	defer pClient.Close()

	for _, info := range list {
		// 发送
		msg := mqview.MQMessage{
			MessageID: time.Now().Format(xtime.FORMAT_yyyyMMddHHmmssfff),
			Message: struct {
				dborm.CameraInfo
				Action string // run/stop/pause
			}{CameraInfo: info, Action: "stop"}}
		mqclient.SendMessage("ffmpeg.cmd.capture.camera", msg, pClient)

		// 日志
		buff, _ := json.Marshal(msg)
		log.Println(string(buff))
		log := dborm.OperateLog{LogDate: time.Now(), LogContent: string(buff)}
		err = log.Insert(pDBMap)
		xdebug.LogError(err)
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}
}
