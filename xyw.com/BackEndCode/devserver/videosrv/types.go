package videosrv

type FFmpegCmd_CameraConfig struct {
	ID          int    // ID
	ClassroomID int    // 教室ID
	CameraIP    string // 摄像机IP
	CameraPort  int    // 摄像机端口
	CmdType     string // 命令类型
	UserID      int    // 命令发送人ID
	UserName    string // 命令发送人名称
	CmdDate     string // 命令发送时间
	CmdState    int    // 命令发送状态 [0:未发送,1:发送中,2:已发送,3:未回应,4:已回应]
	CmdError    string // 错误信息
}

// 用户认证
type UserAuth struct {
	UserID   int
	RoleType int
	Token    string
}

// 摄像机参数
type CameraArgs struct {
	CameraIP   string // 摄像机IP
	CameraPort int    // 摄像机端口
	CameraUser string // 摄像机用户名
	CameraPass string // 摄像机用户密码
}

// 目标参数
type TargetArgs struct {
	TargetIP   string // 目标IP
	TargetPort int    // 目标端口
	TargetUser string // 目标用户名
	TargetPass string // 目标用户密码
}

// 课堂参数
type CurriculumArgs struct {
	CurriculumName     string // 课程ID
	TeacherName        string // 教师ID
	CharpterName       string // 章节ID
	CurriculumDuration int    // 课堂时长(秒)
}

// 视频参数
type FFmpegVideoArgs struct {
	VideoFile     string // 视频文件名
	VideoDuration int    // 录制时长
}

type CmdMessage struct {
	CmdID    string
	CmdType  string
	JsonText string
}

// 视频录制命令
type VideoCaptureCommand struct {
	TargetArgs      // 目标参数
	CurriculumArgs  // 课程视频参数
	FFmpegVideoArgs // 视频录制参数

	CmdID       string // 命令ID
	CmdType     int    // 命令类型(1:BeginVideo, 2:StopVideo, 0:PauseVideo)
	CmdState    int    // 命令发送状态 [0:未发送,1:发送中,2:已发送,3:未回应,4:已回应]
	ClassroomID int    // 教室ID

	//	{"TargetIP": "",
	//    "TargetPort": 0,
	//    "TargetUser": "",
	//    "TargetPass": "",
	//    "CurriculumName": "",
	//    "TeacherName": "",
	//    "CharpterName": "",
	//    "CurriculumDuration": 0,
	//    "VideoFile": "",
	//    "VideoDuration": 0,
	//    "CmdID": "3333333",
	//    "CmdType": 0,
	//    "CmdState": 0,
	//    "ClassroomID": 0}

}

// 视频录制请求
type VideoCaptureRequest struct {
	UserAuth                         // 用户认证
	CurriculumArgs                   // 课堂参数
	ClassroomID                  int // 教室ID
	CurriculumClassroomChapterID int // 课程教室章节ID

	//	{"UserID": 0,
	//    "RoleType": 0,
	//    "Token": "",
	//    "CurriculumName": "",
	//    "TeacherName": "",
	//    "CharpterName": "",
	//    "CurriculumDuration": 0,
	//    "ClassroomID": 0,
	//    "CurriculumClassroomChapterID": 0}
}
