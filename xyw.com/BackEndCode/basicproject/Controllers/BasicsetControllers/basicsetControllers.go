package basicsetControllers

import (
	"encoding/json"
	"time"
	//	"fmt"
	"basicproject/viewmodel"
	"io/ioutil"
	"strconv"
	"strings"
	core "xutils/xcore"

	"basicproject/DataAccess/basicsetDataAccess"
	//	"basicproject/DataAccess/curriculumDataAccess"
	"basicproject/DataAccess/usersDataAccess"

	"github.com/gin-gonic/gin"
)

//获取全部
func GetAll(c *gin.Context) {
	var rd core.Returndata
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	arrit := make([]interface{}, 7)
	arrit0 := basicsetDataAccess.QueryCampus(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit1 := basicsetDataAccess.QueryBuilding(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit2 := basicsetDataAccess.QueryFloors(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit3 := basicsetDataAccess.QueryCollege(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit4 := basicsetDataAccess.QueryMajor(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit5 := basicsetDataAccess.QueryClasses(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit6 := basicsetDataAccess.QueryTeachers(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	arrit[0] = arrit0
	arrit[1] = arrit1
	arrit[2] = arrit2
	arrit[3] = arrit3
	arrit[4] = arrit4
	arrit[5] = arrit5
	arrit[6] = arrit6

	rd.Result = arrit
	rd.Rcode = "1000"
	rd.Reason = "成功"
	c.JSON(200, rd)
}

//获取全部校区
func GetCampus(c *gin.Context) {
	var rd core.Returndata
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	res := basicsetDataAccess.QueryCampus(viewmodel.QueryBasicsetWhere{}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
	rd.Result = res
	rd.Rcode = "1000"
	rd.Reason = "成功"
	c.JSON(200, rd)
}

//获取楼栋
func GetBuilding(c *gin.Context) {
	var rd core.Returndata
	Campusid := c.Query("campusid") //获取校区Id
	Campusidint, _ := strconv.Atoi(Campusid)
	if Campusidint > 0 {
		ws := viewmodel.QueryBasicsetWhere{Campusid: Campusidint}
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		list := basicsetDataAccess.QueryBuilding(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
		if list != nil {
			rd.Result = list
			rd.Rcode = "1000"
			rd.Reason = "成功"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "系统没有预置任何校区数据"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未提交正确的数据"
	}
	c.JSON(200, rd)
}

//根据校区和专业查询教师
func QueryTeachers(c *gin.Context) {
	var rd core.Returndata
	//	Collegeid int
	//	Majorid   int
	collegeid := c.Query("collegeid") //获取学院Id
	collegeidint, _ := strconv.Atoi(collegeid)
	majorid := c.Query("majorid") //获取专业、科系Id
	majoridint, _ := strconv.Atoi(majorid)
	if collegeidint > 0 {
		ws := viewmodel.QueryBasicsetWhere{Collegeid: collegeidint, Majorid: majoridint}
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		list := basicsetDataAccess.QueryTeachers(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
		if list != nil {
			rd.Result = list
			rd.Rcode = "1000"
			rd.Reason = "成功"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "系统没有预置任何校区数据"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未提交正确的数据"
	}
	c.JSON(200, rd)
}

//获取楼层和教室
func GetFloorsandrooms(c *gin.Context) {
	var rd core.Returndata
	Buildingid := c.Query("buildingid") //获取楼栋ID
	Buildingidint, _ := strconv.Atoi(Buildingid)
	if Buildingidint == 0 {
		rd.Rcode = "1002"
		rd.Reason = "未收到请求数据"
	} else {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		ws := viewmodel.QueryBasicsetWhere{Buildingid: Buildingidint}
		floorslist := basicsetDataAccess.QueryFloors(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap)
		if len(floorslist) > 0 {
			for k := 0; k < len(floorslist); k++ {
				floorslist[k].Rooms = basicsetDataAccess.QueryClassrooms(viewmodel.QueryBasicsetWhere{Floorsid: floorslist[k].Floorsid}, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap) //rooms
			}
			rd.Result = floorslist
			rd.Rcode = "1000"
			rd.Reason = "成功"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "系统没有预置任何校区数据"
		}
	}
	c.JSON(200, rd)
}

//获取教室当前状态的详细情况
func GetClassRoomInfo(c *gin.Context) {
	var rd core.Returndata
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	c.Writer.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	classroomidstr := c.Query("id") //教室Id
	classroomid, _ := strconv.Atoi(classroomidstr)
	if classroomid == 0 {
		rd.Rcode = "1002"
		rd.Reason = "未收到请求数据"
	} else {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		ws := viewmodel.QueryBasicsetWhere{Classroomid: classroomid}
		//需要返回的结果有,教室的名称，教室的编号，教室所在的位置，教室当前的状态，教室内所上的课程，授课老师
		classroominfo1 := basicsetDataAccess.QueryClassroomsInfo(ws, dbmap)
		//		if classroominfo1.Classroomstate == 1 {
		datestr := time.Now().Format("2006-01-02 15:04:05")
		datestr = core.Timeaction(datestr)
		//		classroominfo2 := curriculumDataAccess.QueryCurriculumByClassRoom(ws.Classroomid, datestr, dbmap)
		//		if len(classroominfo2) > 0 {
		//			classroominfo1.Curriculumname = classroominfo2[0].Curriculumname
		//			classroominfo1.Nickname = classroominfo2[0].Nickname
		//			classroominfo1.Chaptername = classroominfo2[0].Chaptername
		//			classroominfo1.Classesid = classroominfo2[0].Classesid
		//			classroominfo1.Classesname = classroominfo2[0].Classesname
		//			classroominfo1.Curriculumclassroomchaptercentreid = classroominfo2[0].Curriculumclassroomchaptercentreid
		//			classroominfo1.Qccci = classroominfo2
		//		}
		//		}
		rd.Result = classroominfo1
		rd.Rcode = "1000"
		rd.Reason = "成功"
	}
	c.JSON(200, rd)
}

/*
根据校区Id、楼栋id、楼层id、用户Id等相关信息查询
*/
func QueryClassroom(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryClassroom
	data, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "basicsetControllers|QueryClassroom|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "basicsetControllers|QueryClassroom|令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = basicsetDataAccess.QueryClassroominfo(viewmodel.QueryBasicsetWhere{Campusid: lg.Campusid, Buildingid: lg.Buildingid, Floorsid: lg.Floorsid, Usersid: bt.Usersid}, core.PageData{PageIndex: lg.Pageindex, PageSize: 10}, dbmap)
			//			rd.Rcode = "1000"
			//			rd.Reason = ""
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

/*
查询统计校区，楼栋等人数
*/
func GetQueryPeoples(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryBasicsetWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "basicsetControllers|GetQueryPeoples|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "basicsetControllers|GetQueryPeoples|令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, lg.Usersid, "GetQueryPeoples", dbmap)
			if rd.Rcode == "1000" {
				lg.Floorsids = strings.Trim(lg.Floorsids, "|")
				lg.Floorsids = strings.Trim(lg.Floorsids, " ")
				lg.Buildingids = strings.Trim(lg.Buildingids, "|")
				lg.Buildingids = strings.Trim(lg.Buildingids, " ")
				lg.Campusids = strings.Trim(lg.Campusids, "|")
				lg.Campusids = strings.Trim(lg.Campusids, " ")
				rd.Result = basicsetDataAccess.GetSelectPeoples(lg, dbmap)
				rd.Rcode = "1000"
				rd.Reason = ""
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
			}
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

/*
根据教室Id查询教室内所有人员列表信息
查询条件有？
*/
func GetClassRoomPeopleInfo(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryClassRoomPeopleInfo
	data, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "basicsetControllers|GetClassRoomPeopleInfo|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "basicsetControllers|GetClassRoomPeopleInfo|令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, lg.Usersid, "GetClassRoomPeopleInfo", dbmap)
			if rd.Rcode == "1000" {
				rd.Result = basicsetDataAccess.QueryClassRoomPeopleInfo(lg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
			}
			rd.Rcode = "1000"
			rd.Reason = ""
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

/*
根据教室Id查询教室内所有人员列表信息汇总
查询条件有？
*/
func GetClassRoomPeopleCount(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryClassRoomPeopleInfo
	data, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "basicsetControllers|GetClassRoomPeopleCount|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "basicsetControllers|GetClassRoomPeopleCount|令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, lg.Usersid, "GetClassRoomPeopleCount", dbmap)
			if rd.Rcode == "1000" {
				rd.Result = basicsetDataAccess.QueryClassRoomPeopleCount(lg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
			}
			rd.Rcode = "1000"
			rd.Reason = ""
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

/*
根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计
*/
func GetStreamPeoplesAnalysis(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryStreamPeoplesWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	//fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "basicsetControllers|GetStreamPeoplesAnalysis|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "basicsetControllers|GetStreamPeoplesAnalysis|令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "GetStreamPeoplesAnalysis", dbmap)
			if rd.Rcode == "1000" {
				rd = basicsetDataAccess.QueryStreamPeoplesAnalysis(lg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
			}
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}
