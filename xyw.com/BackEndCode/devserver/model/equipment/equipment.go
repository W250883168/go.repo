package equipment

import (
	"dev.project/BackEndCode/devserver/model/core"
	"fmt"
)

type ClassroomComputerConfig struct {
	Id                int
	Classroomid       int    //教室Id
	Computermac       string //电脑mac
	Computerip        string //教室ip
	ComputerState     int    //电脑状态
	Computerupdaccept int    //电脑upd接收端口
	Computerupdsend   int    //电脑upd发送端口
	Computerwebport   int    //web请求端口
	Computerweburl    string //web请求url
	Computerremarks   string //备注
	OpenDate          string //打开时间
	UpdateDate        string //更新时间
}

type CommandSendlog struct {
	Id           int
	Classroomid  int    //教室Id
	CmdIp        string //目标Ip
	CmdPort      int    //目标端口
	CmdStr       string //命令内容
	CmdType      string //命令类型
	CmdUsersId   int    //命令发送人ID
	CmdUsersName string //命令发送人名称
	CmdDate      string //命令发送时间
	CmdState     int    //命令发送状态[0:未发送,1:发送中,2:已发送,3:未回应,4:已回应,5:已结束,6:异常终止]
	CmdError     string //错误信息
	CmdMac       string //客户端的mac地址
}

type CommandRecord struct { //命令记录表
	Id         int
	CmdStr     string
	CmdType    int
	CmdDescibe string
	CmdPara    string
	CmdCode    string
}

type VideoConfig struct { //视频录像机配置表
	Id              int
	CameraIp        string
	CameraPort      int
	CameraState     int
	Classroomid     int    //教室Id
	CameraLoginUser string //登录账号
	CameraLoginPass string //登录密码
}

func init() {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(ClassroomComputerConfig{}, "classroomcomputerconfig").SetKeys(true, "Id")
	dbmap.AddTableWithName(CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
	dbmap.AddTableWithName(CommandRecord{}, "commandrecord").SetKeys(true, "Id")
	dbmap.AddTableWithName(VideoConfig{}, "videoconfig").SetKeys(true, "Id")
	err := dbmap.CreateTablesIfNotExists()
	core.CheckErr(err, "Create tables failed")
	fmt.Println("lives模块开始初始化")
}
