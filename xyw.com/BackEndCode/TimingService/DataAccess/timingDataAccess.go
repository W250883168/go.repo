package timingDataAccess

import (
	"TimingService/Model"
	"TimingService/Viewmodel"
	"fmt"
	"log"
	"runtime"
	core "xutils/xcore"

	"gopkg.in/gorp.v1"
)

//添加定时任务
func AddTimedTask(obj *taskModel.TimedTask, dbmap *gorp.Transaction) (inerr error) {
	//	dbmap.AddTableWithName(taskModel.TimedTask{}, "TimedTask").SetKeys(true, "TaskId")
	fmt.Printf("%+v", obj)
	inerr = dbmap.Insert(obj)
	core.CheckErr(inerr, "timingDataAccess|AddTimedTask|添加定时任务:")
	return inerr
}

//判断任务是否重复
func QueryUniqueTimedTask(obj taskModel.TimedTask, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(*) from TimedTask where RepeatType=? and RepeatValue=? and EventSetTableId=? and FloorsId=? and ClassRoomId=?;", obj.RepeatType, obj.RepeatValue, obj.EventSetTableId, obj.FloorsId, obj.ClassRoomId)
	return num
}

//修改定时任务
func UpdateTimedTask(obj *taskModel.TimedTask, dbmap *gorp.Transaction) (inerr error) {
	//	dbmap.AddTableWithName(taskModel.TimedTask{}, "TimedTask").SetKeys(true, "TaskId")
	_, inerr = dbmap.Update(obj)
	log.Printf("%+v\n", obj)
	core.CheckErr(inerr, "timingDataAccess|UpdateTimedTask|修改定时任务:")
	return inerr
}

//删除定时任务
func DeleteTimedTask(obj *taskModel.TimedTask, dbmap *gorp.Transaction) (inerr error) {
	_, inerr = dbmap.Exec("delete from TimedTask where TaskId in(?)", obj.TaskId)
	return inerr
}

/*
查询具体的数据信息
*/
func QueryTimedTaskinfo(ws string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var returnobj viewmodel.ViewTaskList
	sql1 := `select tt.TaskId,tt.TaskState,tt.TaskIsOpen,tt.TaskType,tt.TaskExecNum,tt.TaskName,tt.TimePoint,tt.RepeatType,
				tt.RepeatValue,tt.EventSetTableId,tt.ClassRoomId,tt.BuildingId,tt.FloorsId,tt.CampusId,est.EventName,est.EventContent,
				cps.Campusname,bd.Buildingname,fs.Floorname,ifnull(cr.Classroomsname,'')Classroomsname,us.Truename as MakeUsersname,tt.MakeDate
				from TimedTask tt left join EventSetTable est on tt.EventSetTableId=est.EventSetTableId 
				left join campus cps on tt.CampusId=cps.Id left join Building bd on tt.BuildingId=bd.Id 
				left join floors fs on tt.FloorsId=fs.Id left join classrooms cr on tt.ClassRoomId=cr.Id 
				left join users us on tt.MakeUsersId=us.Id where 1=1 `
	sql1 = sql1 + ws + ";"
	sqlerr := dbmap.SelectOne(&returnobj, sql1)
	core.CheckErr(sqlerr, "timingDataAccess|QueryTimedTaskinfo|查询具体的数据信息:")
	if sqlerr == nil && returnobj.TaskId > 0 {
		rd.Rcode = "1000"
		rd.Result = returnobj
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未找到数据"
	}
	return rd
}

/*
打开或关闭场景任务
*/
func OnOrOffTimedTask(IsOpen int, TaskId int, dbmap *gorp.DbMap) (rd core.Returndata) {
	_, sqlerr := dbmap.Exec("update TimedTask set TaskIsOpen=?,TaskState=0,ExecBeginDate='',ExecEndDate='' where TaskId=?;", IsOpen, TaskId)
	if sqlerr == nil {
		rd.Rcode = "1000"
		rd.Result = "操作成功"
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未找到数据"
	}
	return rd
}

/*
查询信息列表
*/
func QueryTimedTasklist(ws string, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := `select count(*) from TimedTask tt left join EventSetTable est on tt.EventSetTableId=est.EventSetTableId 
				left join campus cps on tt.CampusId=cps.Id left join Building bd on tt.BuildingId=bd.Id 
				left join floors fs on tt.FloorsId=fs.Id left join classrooms cr on tt.ClassRoomId=cr.Id where 1=1`
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + ws)
	core.CheckErr(sqlerr, "timingDataAccess|QueryTimedTasklist|查询信息列表:")
	var list []viewmodel.ViewTaskList
	if sqlerr == nil && countint > 0 {
		sql1 := `select tt.TaskId,tt.TaskState,tt.TaskIsOpen,tt.TaskType,tt.TaskExecNum,tt.TaskName,tt.TimePoint,tt.RepeatType,
				tt.RepeatValue,tt.EventSetTableId,tt.ClassRoomId,tt.BuildingId,tt.FloorsId,tt.CampusId,est.EventName,est.EventContent,
				cps.Campusname,bd.Buildingname,fs.Floorname,ifnull(cr.Classroomsname,'')Classroomsname,us.Truename as MakeUsersname,tt.MakeDate
				from TimedTask tt left join EventSetTable est on tt.EventSetTableId=est.EventSetTableId 
				left join campus cps on tt.CampusId=cps.Id left join Building bd on tt.BuildingId=bd.Id 
				left join floors fs on tt.FloorsId=fs.Id left join classrooms cr on tt.ClassRoomId=cr.Id
				left join users us on tt.MakeUsersId=us.Id where 1=1 `
		sql1 = sql1 + ws + " order by case when tt.ExecBeginDate>= now() then 0 else 1 end,abs(TIMESTAMPDIFF(SECOND,now(),tt.ExecBeginDate)) asc" + core.GetLimitString(pg) + ""
		_, sqlerr := dbmap.Select(&list, sql1)
		core.CheckErr(sqlerr, "timingDataAccess|QueryTimedTasklist|查询信息列表:")
		rd.Rcode = "1000"
		pg.PageCount = int(countint)
		pg.PageData = list
		rd.Result = pg
	} else {
		rd.Rcode = "1000"
		pg.PageCount = int(countint)
		pg.PageData = list
		rd.Result = pg
		//		rd.Rcode = "1002"
		//		rd.Reason = "未找到数据"
	}
	return rd
}

//添加定时任务事件
func AddEventSetTable(obj *taskModel.TimedTask, dbmap *gorp.Transaction) (inerr error) {
	//	dbmap.AddTableWithName(taskModel.TimedTask{}, "EventSetTable").SetKeys(true, "EventSetTableId")
	inerr = dbmap.Insert(obj)
	core.CheckErr(inerr, "timingDataAccess|AddEventSetTable|添加定时任务事件:")
	return inerr
}

//修改定时任务事件
func UpdateEventSetTable(obj *taskModel.EventSetTable, dbmap *gorp.Transaction) (inerr error) {
	//	dbmap.AddTableWithName(taskModel.EventSetTable{}, "EventSetTable").SetKeys(true, "EventSetTableId")
	_, inerr = dbmap.Update(obj)
	core.CheckErr(inerr, "timingDataAccess|UpdateEventSetTable|修改定时任务事件:")
	return inerr
}

//删除定时任务事件
func DeleteEventSetTable(obj *taskModel.EventSetTable, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//根据定时任务的Id来删除
		_, inerr = tran.Exec("delete from EventSetTable where EventSetTableId in(?)", obj.EventSetTableId)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

/*
查询事件信息列表
*/
func QueryEventSetTablelist(ws string, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := `select count(*) from EventSetTable where 1=1`
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + ws)
	core.CheckErr(sqlerr, "timingDataAccess|QueryEventSetTablelist|查询事件信息列表")
	if sqlerr == nil && countint > 0 {
		var list []taskModel.EventSetTable
		sql1 := `select * from EventSetTable where 1=1 `
		sql1 = sql1 + ws + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql1)
		core.CheckErr(sqlerr, "timingDataAccess|QueryEventSetTablelist|查询事件信息列表")
		rd.Rcode = "1000"
		pg.PageCount = int(countint)
		pg.PageData = list
		rd.Result = pg
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未找到数据"
	}
	return rd
}

/*
查询事件信息列表
*/
func QueryEventSetTableOne(Id int, dbmap *gorp.DbMap) (list taskModel.EventSetTable) {
	sql1 := `select * from EventSetTable where 1=1 and EventSetTableId=?`
	sqlerr := dbmap.SelectOne(&list, sql1, Id)
	core.CheckErr(sqlerr, "timingDataAccess|QueryEventSetTablelist|查询事件信息列表")
	return list
}
func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
