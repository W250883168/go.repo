package timingControllers

import (
	"TimingService/DataAccess"
	"TimingService/Model"
	"TimingService/Viewmodel"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	//	"strings"
	//	"time"
	core "xutils/xcore"
	"xutils/xerr"
	//	"xutils/xhttp"
	"xutils/xtext"
	//	"xutils/xtime"

	"github.com/gin-gonic/gin"
)

//添加定时任务
func SaveTimedTask(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: "1002", Reason: "数据验证失败,请检测是否填写任务名称和选择任务定时类型"}
	defer func() { c.JSON(200, rd) }()

	//获得查询参数
	var requestData viewmodel.PostSumbitTask
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(requestData.Para.TaskName)   //检测名称不能为空
	xtext.RequireNonBlank(requestData.Para.RepeatType) //任务定时类型不能为空

	//	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(taskModel.TimedTask{}, "TimedTask").SetKeys(true, "TaskId")

	if requestData.Para.EventSetTableId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择响应的事件"
		return
	} else {
		//设置执行代码
		qesto := timingDataAccess.QueryEventSetTableOne(requestData.Para.EventSetTableId, dbmap)
		ec := qesto.EventContent
		ec = strings.Replace(ec, "@UserID", strconv.Itoa(requestData.Auth.Usersid), -1)
		ec = strings.Replace(ec, "@RoomID", strconv.Itoa(requestData.Para.ClassRoomId), -1)
		ec = strings.Replace(ec, "@FloorId", strconv.Itoa(requestData.Para.FloorsId), -1)
		requestData.Para.TaskContent = ec
	}
	if requestData.Para.CampusId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的校区位置"
		return
	}
	if requestData.Para.BuildingId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的校区楼栋位置"
		return
	}
	if requestData.Para.FloorsId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的校区楼栋楼层位置"
		return
	}
	if requestData.Para.TimePoint == "" {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的具体时间，例12:00,14:25,18:30"
		return
	}
	switch requestData.Para.RepeatType {
	case "每天":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	case "自定义":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
		if requestData.Para.RepeatValue == "" {
			rd.Rcode = "1002"
			rd.Reason = "请填写自定义值，例星期一、星期二、星期三、星期四、星期五、星期六、星期天"
			return
		}
		requestData.Para.TaskExecNum = -1
	case "只执行一次":
		requestData.Para.TaskType = 0
		requestData.Para.TaskExecNum = 1
	case "工作日":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	case "法定工作日":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	case "法定节假日":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	default:
		rd.Rcode = "1002"
		rd.Reason = "提交数据出错，未找到此定时任务类型"
		return
	}
	requestData.Para.MakeUsersId = requestData.Auth.Usersid
	requestData.Para.MakeDate = time.Now().Format("2006-01-02 15:04:05") //xtime.FormatString()
	//	requestData.Para.TaskIsOpen = 1
	//开启事务
	trans, err := dbmap.Begin()
	xerr.ThrowPanic(err)
	//	tNow := xtime.NowString()
	defer func() {
		if err != nil {
			trans.Rollback()
		}
	}()
	isUnique := timingDataAccess.QueryUniqueTimedTask(requestData.Para, dbmap) //检查唯一约束
	if isUnique > 0 {
		rd.Rcode = "1001"
		rd.Reason = "添加失败:判断到此为重复场景"
		return
	}
	//	requestData.Para.TaskContent = `{"TTUrl":"http://192.168.0.201:8090/device/node/control/switch/off/floor","Para":{"UserID":"9","FloorID":"1","Params":""}}`
	err = timingDataAccess.AddTimedTask(&requestData.Para, trans)
	xerr.ThrowPanic(err)

	//提交事务
	err = trans.Commit()
	xerr.ThrowPanic(err)

	// OK
	rd.Rcode = "1000"
	rd.Reason = ""
}

//修改定时任务
func ChangeTimedTask(c *gin.Context) {
	defer xerr.CatchPanic()

	var rd = core.Returndata{Rcode: "1002", Reason: "数据验证失败,请检测是否填写任务名称和选择任务定时类型"}
	defer func() { c.JSON(200, rd) }()

	//获得查询参数
	var requestData viewmodel.PostSumbitTask
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	xerr.ThrowPanic(err)

	//校验数据
	xtext.RequireNonBlank(requestData.Para.TaskName)   //检测名称不能为空
	xtext.RequireNonBlank(requestData.Para.RepeatType) //任务定时类型不能为空

	//	//校验权限
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(taskModel.TimedTask{}, "TimedTask").SetKeys(true, "TaskId")

	if requestData.Para.TaskId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的需要提交此任务的Id"
		return
	}
	if requestData.Para.EventSetTableId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择响应的事件"
		return
	} else {
		//设置执行代码
		qesto := timingDataAccess.QueryEventSetTableOne(requestData.Para.EventSetTableId, dbmap)
		ec := qesto.EventContent
		ec = strings.Replace(ec, "@UserID", strconv.Itoa(requestData.Auth.Usersid), -1)
		ec = strings.Replace(ec, "@RoomID", strconv.Itoa(requestData.Para.ClassRoomId), -1)
		ec = strings.Replace(ec, "@FloorId", strconv.Itoa(requestData.Para.FloorsId), -1)
		requestData.Para.TaskContent = ec
	}
	if requestData.Para.CampusId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的校区位置"
		return
	}
	if requestData.Para.BuildingId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的校区楼栋位置"
		return
	}
	if requestData.Para.FloorsId <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的校区楼栋楼层位置"
		return
	}
	if requestData.Para.TimePoint == "" {
		rd.Rcode = "1002"
		rd.Reason = "请选择执行的具体时间，例12:00,14:25,18:30"
		return
	}
	switch requestData.Para.RepeatType {
	case "每天":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	case "自定义":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
		if requestData.Para.RepeatValue == "" {
			rd.Rcode = "1002"
			rd.Reason = "请填写自定义值，例星期一、星期二、星期三、星期四、星期五、星期六、星期天"
			return
		}
		requestData.Para.TaskExecNum = -1
	case "只执行一次":
		requestData.Para.TaskType = 0
		requestData.Para.TaskExecNum = 1
	case "工作日":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	case "法定工作日":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	case "法定节假日":
		requestData.Para.TaskType = 2
		requestData.Para.TaskExecNum = -1
	default:
		rd.Rcode = "1002"
		rd.Reason = "提交数据出错，未找到此定时任务类型"
		return
	}
	requestData.Para.MakeUsersId = requestData.Auth.Usersid
	requestData.Para.MakeDate = time.Now().Format("2006-01-02 15:04:05") //xtime.FormatString()
	//	requestData.Para.TaskIsOpen = 1
	//开启事务
	trans, err := dbmap.Begin()
	xerr.ThrowPanic(err)
	//	tNow := xtime.NowString()
	defer func() {
		if err != nil {
			trans.Rollback()
		}
	}()
	isUnique := timingDataAccess.QueryUniqueTimedTask(requestData.Para, dbmap) //检查唯一约束
	if isUnique != 1 {
		rd.Rcode = "1001"
		rd.Reason = "修改失败:判断到此为重复场景"
		return
	}
	err = timingDataAccess.UpdateTimedTask(&requestData.Para, trans)
	xerr.ThrowPanic(err)

	//提交事务
	err = trans.Commit()
	xerr.ThrowPanic(err)

	// OK
	rd.Rcode = "1000"
	rd.Reason = ""
}

//删除定时任务
func DelTimedTask(c *gin.Context) {
	var rd core.Returndata
	rd.Rcode = "1000"
	//获得查询参数
	var requestData viewmodel.PostSumbitTask
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		return
	}
	//校验数据
	if requestData.Para.TaskId <= 0 {
		rd.Rcode = "1003"
		rd.Reason = "请提交任务Id"
		c.JSON(200, rd)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//开启事务
	trans, err := dbmap.Begin()
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "开启事务时错误：" + err.Error()
		c.JSON(200, rd)
		return
	}

	//故障删除
	err = timingDataAccess.DeleteTimedTask(&requestData.Para, trans)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "删除故障时出错：" + err.Error()
		c.JSON(200, rd)
		trans.Rollback()
		return
	}

	//提交事务
	err = trans.Commit()
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "删除故障时出错：" + err.Error()
		c.JSON(200, rd)
		trans.Rollback()
		return
	}

	c.JSON(200, rd)
}

func GetTimedTaskinfo(c *gin.Context) {

	var rd core.Returndata
	//获得查询参数
	var requestData viewmodel.QueryTaskWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//TaskId
	tKeyWord := mapPara["TaskId"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	keyWord, ok := tKeyWord.(float64)
	if !ok || int(keyWord) <= 0 {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	wheresql := " and tt.TaskId=" + strconv.Itoa(int(keyWord))

	//查询数据
	rd = timingDataAccess.QueryTimedTaskinfo(wheresql, dbmap)
	c.JSON(200, rd)
}

func PostOnOrOffTimedTask(c *gin.Context) {
	var rd core.Returndata
	//获得查询参数
	var requestData viewmodel.QueryTaskWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//TaskId
	tTaskId := mapPara["TaskId"]
	if tTaskId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	TaskId, ok := tTaskId.(float64)
	if !ok || int(TaskId) < 0 {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	//IsOpen
	tIsOpen := mapPara["IsOpen"]
	if tIsOpen == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	IsOpen, ok := tIsOpen.(float64)
	if !ok || int(IsOpen) < 0 {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}

	//查询数据
	rd = timingDataAccess.OnOrOffTimedTask(int(IsOpen), int(TaskId), dbmap)
	c.JSON(200, rd)
}

func GetTimedTasklist(c *gin.Context) {
	var rd core.Returndata
	//获得查询参数
	var requestData viewmodel.QueryTaskWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
		c.JSON(200, rd)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//KeyWord
	tKeyWord := mapPara["KeyWord"]
	if tKeyWord == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	keyWord, ok := tKeyWord.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	//CampusId
	tCampusId := mapPara["CampusId"]
	if tCampusId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	CampusId, ok := tCampusId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	//BuildingId
	tBuildingId := mapPara["BuildingId"]
	if tCampusId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	BuildingId, ok := tBuildingId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	//FloorsId
	tFloorsId := mapPara["FloorsId"]
	if tCampusId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	FloorsId, ok := tFloorsId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	//ClassRoomId
	tClassRoomId := mapPara["ClassRoomId"]
	if tCampusId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	ClassRoomId, ok := tClassRoomId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	//IsME
	tIsME := mapPara["IsME"]
	if tCampusId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	IsME, ok := tIsME.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误"
		c.JSON(200, rd)
		return
	}
	wheresql := ""
	if keyWord != "" {
		wheresql = wheresql + " and tt.TaskName like '%" + keyWord + "%'"
	}
	if CampusId > 0 {
		wheresql = wheresql + " and tt.CampusId=" + strconv.Itoa(int(CampusId))
	}
	if BuildingId > 0 {
		wheresql = wheresql + " and tt.BuildingId=" + strconv.Itoa(int(BuildingId))
	}
	if FloorsId > 0 {
		wheresql = wheresql + " and tt.FloorsId=" + strconv.Itoa(int(FloorsId))
	}
	if ClassRoomId > 0 {
		wheresql = wheresql + " and tt.ClassRoomId=" + strconv.Itoa(int(ClassRoomId))
	}
	if IsME > 0 {
		wheresql = wheresql + " and tt.MakeUsersId=" + strconv.Itoa(requestData.Auth.Usersid) //(int(ClassRoomId))
	}

	//查询数据
	rd = timingDataAccess.QueryTimedTasklist(wheresql, requestData.Page, dbmap)
	c.JSON(200, rd)
}

func GetEventSetTablelist(c *gin.Context) {
	var rd core.Returndata
	//获得查询参数
	var requestData viewmodel.QueryTaskWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &requestData)
	if err != nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误0"
		c.JSON(200, rd)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()

	//获得参数map
	mapPara := requestData.Para.(map[string]interface{}) //获得通过断言实现类型转换

	//KeyWord
	tFloorsId := mapPara["FloorsId"]
	if tFloorsId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误1"
		c.JSON(200, rd)
		return
	}
	FloorsId, ok := tFloorsId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误2"
		c.JSON(200, rd)
		return
	}
	//ClassRoomId
	tClassRoomId := mapPara["ClassRoomId"]
	if tClassRoomId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误3"
		c.JSON(200, rd)
		return
	}
	ClassRoomId, ok := tClassRoomId.(float64)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误4"
		c.JSON(200, rd)
		return
	}
	//NodeId
	tNodeId := mapPara["NodeId"]
	if tNodeId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误5"
		c.JSON(200, rd)
		return
	}
	NodeId, ok := tNodeId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误6"
		c.JSON(200, rd)
		return
	}
	//DeviceId
	tDeviceId := mapPara["DeviceId"]
	if tDeviceId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误7"
		c.JSON(200, rd)
		return
	}
	DeviceId, ok := tDeviceId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误8"
		c.JSON(200, rd)
		return
	}
	//CmdId
	tCmdId := mapPara["CmdId"]
	if tCmdId == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误9"
		c.JSON(200, rd)
		return
	}
	CmdId, ok := tCmdId.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误10"
		c.JSON(200, rd)
		return
	}
	//IsFloors
	tIsFloors := mapPara["IsFloors"]
	if tIsFloors == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误11"
		c.JSON(200, rd)
		return
	}
	IsFloors, ok := tIsFloors.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误12"
		c.JSON(200, rd)
		return
	}
	//IsClassRoom
	tIsClassRoom := mapPara["IsClassRoom"]
	if tIsClassRoom == nil {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误13"
		c.JSON(200, rd)
		return
	}
	IsClassRoom, ok := tIsClassRoom.(string)
	if !ok {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式（参数）错误14"
		c.JSON(200, rd)
		return
	}
	wheresql := ""
	if IsFloors == "1" {
		//			wheresql = wheresql + " and FloorsId>0"
	}
	if IsClassRoom == "1" {
		//			wheresql = wheresql + " and ClassRoomId>0"
	}
	fmt.Println("int(ClassRoomId):", int(ClassRoomId), "ClassRoomId:", ClassRoomId)
	fmt.Println("int(FloorsId):", int(FloorsId), "FloorsId:", FloorsId)
	if int(ClassRoomId) > 0 {
		wheresql = wheresql + " and(ClassRoomId=" + strconv.Itoa(int(ClassRoomId)) + " or ClassRoomId=-1)"
	} else if int(FloorsId) > 0 {
		wheresql = wheresql + " and(FloorsId=" + strconv.Itoa(int(FloorsId)) + " or FloorsId=-1)"
	}
	if NodeId != "" {
		wheresql = wheresql + " and NodeId='" + NodeId + "'"
	}
	if DeviceId != "" {
		wheresql = wheresql + " and DeviceId='" + DeviceId + "'"
	}
	if CmdId != "" {
		wheresql = wheresql + " and CmdId='" + CmdId + "'"
	}
	fmt.Println("wheresql:", wheresql)
	//查询数据
	rd = timingDataAccess.QueryEventSetTablelist(wheresql, requestData.Page, dbmap)
	c.JSON(200, rd)
}
