package httpview

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
)

// 摄像头参数
type CameraConfigView struct {
	CameraName  string
	CameraIP    string
	CameraPort  int
	CameraState int
	LoginUser   string
	LoginPass   string
}

// 教学电脑参数
type ComputerConfigView struct {
	ComputerName  string
	ComputerIP    string
	ComputerPort  int
	ComputerState int
	LoginUser     string
	LoginPass     string
}

// 视频录制请求
type VideoCaptureRequest struct {
	LessonID   int // 课堂ID
	UserID     int
	LocationID int
	Teacher    string
	Curriclum  string
	Chapter    string
	FileName   string
	Duration   int // 时长(ms)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}

	data, _ := json.Marshal(&VideoCaptureRequest{})
	log.Printf("%s\n", data)
}
