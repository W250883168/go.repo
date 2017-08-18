package viewmodel

import (
	"TimingService/Model"
	"log"
	core "xutils/xcore"
)

type PostSumbitTask struct { //任务提交对象
	Auth  core.BasicsToken        //身份认证
	Para  taskModel.TimedTask     //定时任务数据
	Event taskModel.EventSetTable //定义的事情信息
}
type TTScanDevice struct {
	TTUrl string
	Para  interface{}
}

//type SendSwitch struct {
//	UId     string
//	Id      string
//	Type    string
//	CmdCode string
//	Para    string
//}
type ViewTaskList struct {
	TaskId          int
	TaskState       int    //:定时任务的状态[0:未启动,1:已启动,2:已结束]
	TaskIsOpen      int    //任务的开启状态[0:未开启，1:已开启]
	TaskType        int    //:任务类型[0:单次任务,1:循环任务,2:多次任务]
	TaskExecNum     int    //:执行次数[多次任务时:填写执行的次数，执行一次-1，循环任务时：默认-1，单次任务时：1]
	MakeUsersId     int    //:制定任务人Id
	MakeDate        string //:任务制定时间
	TaskName        string //:任务的名称
	TimePoint       string //:任务定的时间触发的[示例：每天的下午2点 ,存储特定的时间格式数据]
	RepeatType      string //:重复类型[执行一次、每天、工作日、周末、自定义]
	RepeatValue     string //:重复类型的值
	EventSetTableId int    //事件定义Id
	ClassRoomId     int    //教室Id
	BuildingId      int    //楼栋Id[必填]
	FloorsId        int    //楼层Id
	CampusId        int    //校区Id[必填]
	Classroomsname  string //教室Id
	BuildingName    string //楼栋Id[必填]
	Floorname       string //楼层Id
	Campusname      string //校区名称
	EventName       string //事件名称
	EventContent    string //事件执行的内容
	MakeUsersname   string //制定人名称
}

type QueryTaskWhere struct { //任务查询条件对象
	Auth core.BasicsToken //身份认证
	Page core.PageData    //分页信息
	Para interface{}      //查询参数
}

func init() {
	log.Println("viewmodel模块开始初始化")
}
