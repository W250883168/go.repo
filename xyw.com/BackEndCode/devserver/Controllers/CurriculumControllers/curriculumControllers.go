package curriculumControllers

import (
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	//	"dev.project/BackEndCode/devserver/model/actiondata"
	"dev.project/BackEndCode/devserver/DataAccess/curriculumDataAccess"
	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/viewmodel"

	"github.com/gin-gonic/gin"
)

/*
获取某个课程下某个班级每个学生的平均到课率
*/
func GetEveryoneAverageclassrate(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.GetAverageclassrate
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "curriculumControllers|GetEveryoneAverageclassrate|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "curriculumControllers|GetEveryoneAverageclassrate|令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "geteveryoneaverageclassrate", dbmap)
			if rd.Rcode == "1000" {
				res := curriculumDataAccess.GetEveryoneAverageclassrate(lg, bt, dbmap)
				rd.Result = res
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
提交数据：教师Id，班级Id
获取教师班级下所有课程详细的到课数据
*/
func GetCurriculumChaptersInfo(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.GetAverageclassrate
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "curriculumControllers|GetCurriculumChaptersInfo|获取教师班级下所有课程详细的到课数据")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "GetCurriculumChaptersInfo", dbmap)
			if rd.Rcode == "1000" {
				res := curriculumDataAccess.QueryCurriculumChaptersInfo(lg, bt, dbmap)
				rd.Result = res
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
获取某个课程下某个班级每个学生的平均到课率
*/
func GetStudentsClassesAvg(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.GetStudentsinfo
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "curriculumControllers|GetStudentsClassesAvg|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "curriculumControllers|GetStudentsClassesAvg|令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "getstudentsclassesavg", dbmap)
			if rd.Rcode == "1000" {
				res := curriculumDataAccess.GetStudentsClassesAvg(lg, dbmap)
				rd.Result = res
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
查询课程、学科分类
*/
func GetSubjectclassList(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.GetStudentsinfo
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "curriculumControllers|GetSubjectclassList|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "curriculumControllers|GetSubjectclassList|令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			res := curriculumDataAccess.GetSubjectclassList(lg, dbmap)
			rd.Result = res
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
获取筛选条件数据
*/
func GetFilterData(c *gin.Context) {
	var rd core.Returndata
	var wh viewmodel.QueryFilterWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &wh)
	core.CheckErr(errs1, "curriculumControllers|GetFilterData|登录数据json格式转换错误"+string(data))
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "curriculumControllers|GetFilterData|令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			arrint := make([]interface{}, 2)
			arrint0 := curriculumDataAccess.GetFilterDataCurriculum(bt, wh, dbmap)
			arrint1 := curriculumDataAccess.GetFilterDataClasses(bt, wh, dbmap)
			arrint[0] = arrint0
			arrint[1] = arrint1
			rd.Result = arrint
			rd.Rcode = "1000"
			rd.Reason = "成功"
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
获取班级的到课平均率
*/
func GetClassAveragerate(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.GetAverageclassrate
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "curriculumControllers|GetClassAveragerate|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "curriculumControllers|GetClassAveragerate|令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "getclassaveragerate", dbmap)
			if rd.Rcode == "1000" {
				res := curriculumDataAccess.GetClassAveragerate(lg, bt, dbmap)
				rd.Result = res
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
管理者查看各种出勤统计分析
*/
func AttendanceAnalysis(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.AttendanceAnalysisWhere
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "curriculumControllers|AttendanceAnalysis|登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "curriculumControllers|AttendanceAnalysis|令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "AttendanceAnalysis", dbmap)
			if rd.Rcode == "1000" {
				res := curriculumDataAccess.QueryAttendanceAnalysisList(lg, bt, dbmap)
				rd.Result = res
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
