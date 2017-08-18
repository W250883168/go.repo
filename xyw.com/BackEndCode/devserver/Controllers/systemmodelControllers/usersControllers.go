package systemmodelControllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"xutils/xtext"

	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/commons"
	"dev.project/BackEndCode/devserver/commons/xdebug"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/users"
	"dev.project/BackEndCode/devserver/viewmodel"

	"github.com/gin-gonic/gin"
)

// 查询用户人员
func GetUsersList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryUsersWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetUsersList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetUsersList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetUsersList", dbmap)
		if rd.Rcode == "1000" {
			obj.Rolestype = -3
			rd = usersDataAccess.QueryUsersPG(obj, page, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//查询用户人员
func GetStudentsList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryUsersWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetStudentsList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetStudentsList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetStudentsList", dbmap)
		if rd.Rcode == "1000" {
			obj.Rolestype = 3
			rd = usersDataAccess.Query_StudentInfoDetail_ByPage(obj, page, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//查询用户人员
func GetTeacherList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryUsersWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetTeacherList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetTeacherList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetTeacherList", dbmap)
		if rd.Rcode == "1000" {
			obj.Rolestype = 2
			rd = usersDataAccess.QueryUsersPG(obj, page, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加用户人员
func AddUsers(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddUsers|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddUsers", dbmap)
		if rd.Rcode == "1000" {
			obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
			obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
			us := users.Users{Loginuser: obj.Loginuser, Loginpwd: obj.Loginpwd, Rolesid: obj.Rolesid, Truename: obj.Truename, Nickname: obj.Nickname,
				Userheadimg: obj.Userheadimg, Userphone: obj.Userphone, Userstate: 0, Usersex: obj.Usersex, Usermac: "", Birthday: obj.Birthday}
			tr := users.Teacher{Collegeid: 0, Majorid: 0}
			st := users.Students{Enrollmentyear: 0, Homeaddress: "", Nowaddress: "", Classesid: 0, Infostate: 0, Currentstate: 0, Attendance: 0,
				Needcoursenum: 0, Alreadycoursenum: 0, Absenteeism: 0}
			adderr := usersDataAccess.AddUsersInfoAll(us, st, tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|AddUsers|添加用户数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Result = us.Id
			} else {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加用户人员
func AddStudents(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddStudents|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		//		lgs, _ := c.Get("users")
		//		lg := lgs.(core.BasicsToken)
		//		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddStudents", dbmap)
		//		if rd.Rcode == "1000" {
		obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
		obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
		us := users.Users{Loginuser: obj.Loginuser, Loginpwd: obj.Loginpwd, Rolesid: obj.Rolesid, Truename: obj.Truename, Nickname: obj.Nickname,
			Userheadimg: obj.Userheadimg, Userphone: obj.Userphone, Userstate: 0, Usersex: obj.Usersex, Usermac: "", Birthday: obj.Birthday}
		tr := users.Teacher{Collegeid: 0, Majorid: 0}
		st := users.Students{Enrollmentyear: obj.Enrollmentyear, Homeaddress: obj.Homeaddress, Nowaddress: obj.Nowaddress, Classesid: obj.Classesid, Infostate: 1,
			Currentstate: 1, Attendance: 0, Needcoursenum: 0, Alreadycoursenum: 0, Absenteeism: 0}
		adderr := usersDataAccess.AddUsersInfoAll(us, st, tr, dbmap)
		core.CheckErr(adderr, "systemmodelControllers|AddStudents|添加学生数据失败")
		if adderr == nil {
			rd.Rcode = "1000"
			rd.Result = us.Id
		} else {
			rd.Rcode = "1001"
			rd.Reason = "添加失败:" + adderr.Error()
		}
		//		} else {
		//			rd.Rcode = "1105"
		//			rd.Reason = "此功能未授权,不可用"
		//		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加用户人员
func AddTeacher(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddTeacher|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddTeacher", dbmap)
		if rd.Rcode == "1000" {

			obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
			obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
			us := users.Users{Loginuser: obj.Loginuser, Loginpwd: obj.Loginpwd, Rolesid: obj.Rolesid, Truename: obj.Truename, Nickname: obj.Nickname,
				Userheadimg: obj.Userheadimg, Userphone: obj.Userphone, Userstate: 0, Usersex: obj.Usersex, Usermac: "", Birthday: obj.Birthday}
			tr := users.Teacher{Collegeid: obj.Collegeid, Majorid: obj.Majorid}
			st := users.Students{Enrollmentyear: 0, Homeaddress: "", Nowaddress: "", Classesid: 0, Infostate: 0, Currentstate: 0, Attendance: 0,
				Needcoursenum: 0, Alreadycoursenum: 0, Absenteeism: 0}
			adderr := usersDataAccess.AddUsersInfoAll(us, st, tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|AddTeacher|添加教师数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Result = us.Id
			} else {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

// 修改用户人员
func ChangeUsers(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	// fmt.Println("	<<<<<<<<<<<<\n", strdata)
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeUsers|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeUsers", dbmap)
		if rd.Rcode == "1000" {
			obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
			obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
			us := users.Users{Id: obj.Id, Loginuser: obj.Loginuser, Loginpwd: obj.Loginpwd, Rolesid: obj.Rolesid, Truename: obj.Truename, Nickname: obj.Nickname,
				Userheadimg: obj.Userheadimg, Userphone: obj.Userphone, Userstate: 0, Usersex: obj.Usersex, Usermac: "", Birthday: obj.Birthday}
			tr := users.Teacher{Id: 0, Collegeid: 0, Majorid: 0}
			st := users.Students{Enrollmentyear: 0, Homeaddress: "", Nowaddress: "", Classesid: 0, Infostate: 0, Currentstate: 0, Attendance: 0,
				Needcoursenum: 0, Alreadycoursenum: 0, Absenteeism: 0}
			adderr := usersDataAccess.UpdateUsersInfoAll(&us, &st, &tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|ChangeUsers|添加校区数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

// 修改学生人员
func ChangeStudents(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeStudents|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		//		lgs, _ := c.Get("users")
		//		lg := lgs.(core.BasicsToken)
		//		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeStudents", dbmap)
		//		if rd.Rcode == "1000" {
		obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
		obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
		us := users.Users{Id: obj.Id, Loginuser: obj.Loginuser, Loginpwd: obj.Loginpwd, Rolesid: obj.Rolesid, Truename: obj.Truename, Nickname: obj.Nickname,
			Userheadimg: obj.Userheadimg, Userphone: obj.Userphone, Userstate: 0, Usersex: obj.Usersex, Usermac: "", Birthday: obj.Birthday}
		tr := users.Teacher{Collegeid: 0, Majorid: 0}
		st := users.Students{Id: obj.Id, Enrollmentyear: obj.Enrollmentyear, Homeaddress: obj.Homeaddress, Nowaddress: obj.Nowaddress, Classesid: obj.Classesid, Infostate: 1,
			Currentstate: 1}
		fmt.Println("us: ", us)
		fmt.Println("st: ", st)
		fmt.Println("tr: ", tr)
		adderr := usersDataAccess.UpdateUsersInfoAll(&us, &st, &tr, dbmap)
		core.CheckErr(adderr, "systemmodelControllers|ChangeStudents|添加校区数据失败")
		if adderr == nil {
			rd.Rcode = "1000"
			rd.Reason = "修改成功"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "修改失败:" + adderr.Error()
		}
		//		} else {
		//			rd.Rcode = "1105"
		//			rd.Reason = "此功能未授权,不可用"
		//		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改用户人员
func ChangeTeacher(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeTeacher|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeTeacher", dbmap)
		if rd.Rcode == "1000" {
			obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
			obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
			us := users.Users{Id: obj.Id, Loginuser: obj.Loginuser, Loginpwd: obj.Loginpwd, Rolesid: obj.Rolesid, Truename: obj.Truename, Nickname: obj.Nickname,
				Userheadimg: obj.Userheadimg, Userphone: obj.Userphone, Userstate: 0, Usersex: obj.Usersex, Usermac: "", Birthday: obj.Birthday}
			tr := users.Teacher{Id: obj.Id, Collegeid: obj.Collegeid, Majorid: obj.Majorid}
			st := users.Students{Enrollmentyear: 0, Homeaddress: "", Nowaddress: "", Classesid: 0, Infostate: 0, Currentstate: 0, Attendance: 0,
				Needcoursenum: 0, Alreadycoursenum: 0, Absenteeism: 0}
			adderr := usersDataAccess.UpdateUsersInfoAll(&us, &st, &tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|ChangeTeacher|添加校区数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除用户人员
func DelUsers(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelUsers|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelUsers", dbmap)
		if rd.Rcode == "1000" {
			us := users.Users{Id: obj.Id}
			tr := users.Teacher{Id: obj.TeacherId}
			st := users.Students{Id: obj.StudentsId}
			adderr := usersDataAccess.DelUsersInfoAll(&us, &st, &tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelUsers|删除校区数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除学生人员
func DelStudents(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelStudents|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelStudents", dbmap)
		if rd.Rcode == "1000" {
			us := users.Users{Id: obj.Id, Rolesid: 3}
			tr := users.Teacher{Id: obj.Id}
			st := users.Students{Id: obj.Id}
			adderr := usersDataAccess.DelUsersInfoAll(&us, &st, &tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelStudents|删除校区数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除教师人员
func DelTeacher(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelTeacher|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelTeacher", dbmap)
		if rd.Rcode == "1000" {
			us := users.Users{Id: obj.Id, Rolesid: 2}
			tr := users.Teacher{Id: obj.Id}
			st := users.Students{Id: obj.Id}
			adderr := usersDataAccess.DelUsersInfoAll(&us, &st, &tr, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelTeacher|删除教师人员")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

/****************/
// 修改用户信息
func ModifyUserInfo(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.UsersInfoAll

	resp := commons.ResponseMsgSet_Instance()
	rd.Rcode = strconv.Itoa(resp.DATA_MALFORMED.Code)
	rd.Reason = resp.DATA_MALFORMED.Text

	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	if err := json.Unmarshal(data, &obj); err != nil {
		core.CheckErr(err, "systemmodelControllers|ChangeUsers|模块数据转换失败")
		xdebug.PrintStackTrace(err)
		c.JSON(200, rd)
		return
	}

	TAG := "ChangeUsers"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	lgs, _ := c.Get("users")
	lg := lgs.(core.BasicsToken)
	rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, TAG, dbmap)
	if rd.Rcode == strconv.Itoa(resp.SUCCESS.Code) {
		obj.Loginpwd = xtext.SpaceTrim(obj.Loginpwd)
		obj.Loginuser = xtext.SpaceTrim(obj.Loginuser)
		us := users.Users{
			Id:          obj.Id,
			Loginuser:   obj.Loginuser,
			Loginpwd:    obj.Loginpwd,
			Rolesid:     obj.Rolesid,
			Truename:    obj.Truename,
			Nickname:    obj.Nickname,
			Userheadimg: obj.Userheadimg,
			Userphone:   obj.Userphone,
			Userstate:   0,
			Usersex:     obj.Usersex,
			Usermac:     "",
			Birthday:    obj.Birthday}

		rd.Rcode = strconv.Itoa(resp.SUCCESS.Code)
		rd.Reason = resp.SUCCESS.Text
		if sqlerr := usersDataAccess.UpdateUser_ExceptPassword(&us, dbmap); sqlerr != nil {
			rd.Rcode = strconv.Itoa(resp.FAIL.Code)
			rd.Reason = resp.FAIL.Text
			xdebug.PrintStackTrace(sqlerr)
		}
	} else {
		rd.Rcode = strconv.Itoa(resp.AUTH_LIMITED.Code)
		rd.Reason = resp.AUTH_LIMITED.Text
	}

	c.JSON(200, rd)
}

// 查询用户人员
func GetUsersList2(c *gin.Context) {
	var rd core.Returndata
	var where = users.QueryWhere4User{PageIndex: 0, PageSize: 100, RoleID: -1}
	resp := commons.ResponseMsgSet_Instance()
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	if err := json.Unmarshal(data, &where); err != nil {
		rd.Rcode = strconv.Itoa(resp.DATA_MALFORMED.Code)
		rd.Reason = resp.DATA_MALFORMED.Text
		c.JSON(200, rd)

		xdebug.PrintStackTrace(err)
		return
	}
	TAG := "GetUsersList"
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	lgs, _ := c.Get("users")
	lg := lgs.(core.BasicsToken)
	rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, TAG, dbmap)
	if rd.Rcode == strconv.Itoa(resp.SUCCESS.Code) {
		rd = usersDataAccess.QueryUsersPG2(where, dbmap)
	} else {
		rd.Rcode = strconv.Itoa(resp.AUTH_LIMITED.Code)
		rd.Reason = resp.AUTH_LIMITED.Text
	}

	c.JSON(200, rd)
}
