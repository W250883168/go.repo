package taskModel

import (
	"log"
)

type TimedTask struct { //定时任务
	TaskId          int
	TaskState       int    //:定时任务的状态[0:未启动,1:已启动,2:已结束]
	TaskIsOpen      int    //任务的开启状态[0:未开启，1:已开启]
	TaskType        int    //:任务类型[0:单次任务,1:循环任务,2:多次任务]
	TaskExecNum     int    //:执行次数[多次任务时:填写执行的次数，执行一次-1，循环任务时：默认-1，单次任务时：1]
	MakeUsersId     int    //:制定任务人Id
	MakeDate        string //:任务制定时间
	ExecBeginDate   string //:任务开始的时间[精确到时分秒]
	ExecEndDate     string //:任务结束的时间[精确到时分秒]
	TaskName        string //:任务的名称
	TaskContent     string //:任务的内容[执行的代码]
	TimeLong        int    //:任务定的时长[示例：定10分钟后执行]
	TimePoint       string //:任务定的时间触发的[示例：每天的下午2点 ,存储特定的时间格式数据]
	RepeatType      string //:重复类型[执行一次、每天、工作日、周末、自定义]
	RepeatValue     string //:重复类型的值
	EventSetTableId int    //事件定义Id
	ClassRoomId     int    //教室Id
	BuildingId      int    //楼栋Id[必填]
	FloorsId        int    //楼层Id
	CampusId        int    //校区Id[必填]
}
type EventSetTable struct { //事件设置表
	EventSetTableId int    //事件定义Id
	EventName       string //事件名称
	EventContent    string //事件执行的内容
	ClassRoomId     int    //教室Id
	BuildingId      int    //楼栋Id[必填]
	FloorsId        int    //楼层Id
	CampusId        int    //校区Id[必填]
	NodeId          string //节点Id
	DeviceId        string //设备Id
	CmdId           string //命令Id
}

func init() {
	log.Println("taskModel模块开始初始化")
}
