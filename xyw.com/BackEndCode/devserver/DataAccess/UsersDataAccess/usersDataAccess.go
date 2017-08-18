package usersDataAccess

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"dev.project/BackEndCode/devserver/commons"
	"dev.project/BackEndCode/devserver/commons/xdebug"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/users"
	"dev.project/BackEndCode/devserver/viewmodel"

	"gopkg.in/gorp.v1"
)

var resp *commons.ResponseMsgSet

func GetWhereString(ws viewmodel.QueryUsersWhere) (where string) {

	if ws.Searhtxt != "" {
		searhint, _ := strconv.Atoi(ws.Searhtxt)
		if searhint != 0 {
			where = where + " and(us.Id=" + ws.Searhtxt
			where = where + " or st.Classesid=" + ws.Searhtxt
			where = where + " or st.Enrollmentyear=" + ws.Searhtxt
			where = where + " or us.Userphone='" + ws.Searhtxt + "'"
			where = where + " or us.Loginuser='" + ws.Searhtxt + "')"
		} else {
			where = where + " and(us.Loginuser='" + ws.Searhtxt + "' or us.Loginuser like '%" + ws.Searhtxt + "%'"
			where = where + " or us.Truename='" + ws.Searhtxt + "' or us.Truename like '%" + ws.Searhtxt + "%'"
			where = where + " or us.Nickname='" + ws.Searhtxt + "' or us.Nickname like '%" + ws.Searhtxt + "%'"
			where = where + " or us.Userphone='" + ws.Searhtxt + "' or st.Homeaddress like '%" + ws.Searhtxt + "%'"
			where = where + " or st.Nowaddress like '%" + ws.Searhtxt + "%')"
		}
	}
	if ws.Studentsid > 0 {
		where = where + " and us.Id=" + strconv.Itoa(ws.Studentsid)
	}
	if ws.Id > 0 {
		where = where + " and us.Id=" + strconv.Itoa(ws.Id)
	}
	if ws.Rolestype < 0 {
		rtstr := strconv.Itoa(ws.Rolestype)
		rtstr = strings.Replace(rtstr, "-", "", -1)
		where = where + " and us.Rolesid<>" + rtstr
		fmt.Println(rtstr)
	} else if ws.Rolestype > 0 {
		where = where + " and us.Rolesid=" + strconv.Itoa(ws.Rolestype)
	}
	if ws.Classesids != "" {
		where = where + " and st.Classesid in(" + ws.Classesids + ")"
	}
	if ws.Classesid > 0 {
		where = where + " and st.Classesid=" + strconv.Itoa(ws.Classesid)
	}
	if ws.Userstate > 0 {
		where = where + " and us.Userstate=" + strconv.Itoa(ws.Userstate)
	}
	if ws.Usersex > 0 {
		where = where + " and us.Usersex=" + strconv.Itoa(ws.Usersex)
	}
	if ws.Infostate > 0 {
		where = where + " and st.Infostate=" + strconv.Itoa(ws.Infostate)
	}
	if ws.Currentstate > 0 {
		where = where + " and st.Currentstate=" + strconv.Itoa(ws.Currentstate)
	}
	if ws.Collegeids != "" {
		where = where + " and tr.Collegeid in(" + ws.Collegeids + ")"
	}
	if ws.Collegeid > 0 {
		where = where + " and tr.Collegeid=" + strconv.Itoa(ws.Collegeid)
	}
	if ws.Majorids != "" {
		where = where + " and tr.Majorid in(" + ws.Majorids + ")"
	}
	if ws.Majorid > 0 {
		where = where + " and tr.Majorid=" + strconv.Itoa(ws.Majorid)
	}
	return where
}

/*
获取所有校区
根据校区代码查询某校区
查询所有校区
*/
func QueryUsersPG(ws viewmodel.QueryUsersWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from users as us left join students as st on us.Id=st.Id left join teacher as tr on tr.Id=us.Id where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "usersDataAccess|QueryUsersPG|获取用户数量")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.UsersInfoAll
		sql := "select us.Id as UsersId,Loginuser,Rolesid, IFNULL((SELECT roles.Rolesname FROM roles WHERE us.Rolesid = roles.Id), '') AS RoleName, Truename,Nickname,Userheadimg,Userphone,Userstate,Usersex,Usermac,Birthday,"
		sql = sql + " ifnull(st.Id,0) as StudentsId,ifnull(Enrollmentyear,0)as Enrollmentyear,ifnull(Homeaddress,'')as Homeaddress,ifnull(Nowaddress,'')as Nowaddress,"
		sql = sql + " ifnull(Classesid,0)as Classesid,ifnull(Infostate,0)as Infostate,ifnull(Currentstate,0)as Currentstate,ifnull(Attendance,0)as Attendance,ifnull(Needcoursenum,0)as Needcoursenum,ifnull(Alreadycoursenum,0)as Alreadycoursenum,ifnull(Absenteeism,0)as Absenteeism,"
		sql = sql + " ifnull(tr.Id,0) as TeacherId,ifnull(tr.Collegeid,0)as Collegeid, IFNULL((SELECT college.Collegename FROM college WHERE college.Id = tr.Collegeid), '') AS CollegeName, ifnull(tr.Majorid,0)as Majorid, IFNULL((SELECT major.Majorname FROM major WHERE major.Id = tr.Majorid),'') AS MajorName from users as us left join students as st on us.Id=st.Id left join teacher as tr on tr.Id=us.Id where 1=1"
		sql = sql + wheresql + core.GetLimitString(pg) + ";"
		fmt.Println("	<<<<<<<<<<<<<< \n" + sql)
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "usersDataAccess|QueryUsersPG|系统后台获取所有校区:")
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

//获取我的消息
func QueryIsNewMessage(lg core.BasicsToken, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from usermessage where 1=1"
	wheresql := " and Usersid=" + strconv.Itoa(lg.Usersid) + " and Readdate=Createdate"
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "usersDataAccess|QueryIsNewMessage|获取我的消息:")
	if sqlerr == nil {
		if countint > 0 {
			rd.Rcode = "1000"
			rd.Result = countint
		} else {
			rd.Result = 0
			rd.Rcode = "1099"
			rd.Reason = "未查找到数据"
		}
	} else {
		rd.Rcode = "1001"
		rd.Reason = "查询数据失败"
	}
	return rd
}

//获取我的消息
func QueryMyMessage(lg core.BasicsToken, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	//	sql := "select Id as Messageid,Title,Details,State,Createdate,Megtype,Readdate,GoUrl,GoParameter,MessageImg from usermessage where Usersid=? and Readdate=Createdate order by Id desc;"
	//	_, listerr := dbmap.Select(&list, sql, lg.Rolestype) //查询权限模块
	//	core.CheckErr(listerr, "usersDataAccess|GetMymess|获取我的消息:")
	//	return list
	countsql := "select count(*) from usermessage where 1=1"
	wheresql := " and Usersid=" + strconv.Itoa(lg.Usersid)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "usersDataAccess|QueryMyMessage|获取我的消息:")
	if sqlerr == nil {
		if countint > 0 {
			var list []viewmodel.Umessage
			sql := "select Id as Messageid,Title,Details,State,Createdate,Megtype,Readdate,GoUrl,GoParameter,MessageImg,MessageProfiles from usermessage where 1=1 "
			sql = sql + wheresql + " order by Createdate desc " + core.GetLimitString(pg) + ";"
			_, sqlerr := dbmap.Select(&list, sql)
			core.CheckErr(sqlerr, "usersDataAccess|QueryMyMessage|获取我的消息:")
			_, sqlerr = dbmap.Exec("update usermessage set Readdate='" + time.Now().Format("2006-01-02 15:04:05") + "' where 1=1 " + wheresql + ";")
			core.CheckErr(sqlerr, "usersDataAccess|QueryMyMessage|修改消息的读取时间:")
			rd.Rcode = "1000"
			pg.PageCount = int(countint)
			pg.PageData = list
			rd.Result = pg
		} else {
			rd.Rcode = "1099"
			rd.Reason = "未查找到数据"
		}
	} else {
		rd.Rcode = "1001"
		rd.Reason = "查询数据失败"
	}
	return rd
}

//将数据写入到数据库中
func SendUserMessage(mgs *users.Usermessage, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(users.Usermessage{}, "usermessage").SetKeys(true, "Id")
	inserterr := dbmap.Insert(mgs)
	if inserterr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1001"
	}
	return rd
}

//根据用户的账号查询用户的相关信息
func QueryUsersInfo(Loginuser string, dbmap *gorp.DbMap) (us users.Users) {
	sql := "select * from users where Loginuser='" + Loginuser + "' limit 0,1;"
	fmt.Println(sql)
	selectoneerr := dbmap.SelectOne(&us, sql)
	core.CheckErr(selectoneerr, "usersDataAccess|QueryUsersInfo|根据账号查询用户信息错误："+Loginuser)
	return us
}

//获取学生信息
func GetStudentsinfo(lg viewmodel.GetStudentsinfo, bt core.BasicsToken, dbmap *gorp.DbMap) (ssarr viewmodel.Studentsinfo) {
	sql := "select st.Id as Studentsid,st.Enrollmentyear,st.Homeaddress,st.Nowaddress,st.Classesid,st.Infostate,st.Currentstate,us.Truename,us.Nickname,us.Userheadimg,us.Userphone,us.Usersex,us.Birthday,us.Rolesid,us.Userstate from students as st inner join users as us on st.Id=us.Id where st.Id=?;"
	errs2 := dbmap.SelectOne(&ssarr, sql, lg.Studentsid)
	core.CheckErr(errs2, "usersDataAccess|GetStudentsinfo|登录方法中根据账号密码查询错误:")
	type ListInt struct {
		Ct int
	}
	var ltint []ListInt
	sql = "select count(*)as Ct from students as st inner join curriculumsclasscentre as ccc on st.Classesid=ccc.Classesid"
	sql = sql + " inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid where st.Id=? union select count(*)as Ct from students as st "
	sql = sql + " inner join curriculumsclasscentre as ccc on st.Classesid=ccc.Classesid inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid"
	sql = sql + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id where st.Id=? and tr.state=2 union select count(*)as Ct from students as st"
	sql = sql + " inner join curriculumsclasscentre as ccc on st.Classesid=ccc.Classesid inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid"
	sql = sql + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id inner join pointtos as ps on (ps.Curriculumclassroomchaptercentreid=cccc.id and ps.Usersid=st.Id) where st.Id=? and tr.state=2 and ps.state=1;"
	_, errs3 := dbmap.Select(&ltint, sql, lg.Studentsid, lg.Studentsid, lg.Studentsid)
	if len(ltint) == 3 {
		ssarr.Needcoursenum = ltint[0].Ct
		ssarr.Alreadycoursenum = ltint[1].Ct
		ssarr.Absenteeism = ltint[2].Ct
		ssarr.Attendance = float32(ssarr.Absenteeism) / float32(ssarr.Alreadycoursenum)
		//更新底层数据
		dbmap.Exec("update students set Needcoursenum=?,Alreadycoursenum=?,Absenteeism=?,Attendance=? where Id=?;", ssarr.Needcoursenum, ssarr.Alreadycoursenum, ssarr.Absenteeism, ssarr.Attendance, lg.Studentsid)
	}
	core.CheckErr(errs3, "usersDataAccess|GetStudentsinfo|登录方法中根据账号密码查询错误:")
	if errs2 == nil && errs3 == nil {
		return ssarr
	} else {
		return viewmodel.Studentsinfo{}
	}
}

//添加用户
func AddUsers(us *users.Users, st *users.Students, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(users.Students{}, "students")
	inerr = dbmap.Insert(us)
	if inerr == nil && us.Rolesid == 3 { //判断是否是学生
		st.Id = us.Id
		inerr = dbmap.Insert(st)
		core.CheckErr(inerr, "usersDataAccess|AddUsers|添加学生详细信息:")
		_, inerr = dbmap.Exec("update classes set Classesnum=Classesnum+1 where Id=?;", st.Classesid) //更新班级人数
		core.CheckErr(inerr, "usersDataAccess|AddUsers|更新班级人数:")
		majorid, inerr := dbmap.SelectInt("select Majorid from classes where Id=?", st.Classesid) //获取专业ID
		core.CheckErr(inerr, "usersDataAccess|AddUsers|获取专业ID:")
		_, inerr = dbmap.Exec("update major set Majornum=Majornum+1 where Id=?;", majorid) //更新专业人数
		core.CheckErr(inerr, "usersDataAccess|AddUsers|更新专业人数:")
		collegeid, inerr := dbmap.SelectInt("select Collegeid from major where Id=?", majorid) //获取学院ID
		core.CheckErr(inerr, "usersDataAccess|AddUsers|获取学院ID:")
		_, inerr = dbmap.Exec("update college set Collegenum=Collegenum+1 where Id=?;", collegeid) //更新学院人数
		core.CheckErr(inerr, "usersDataAccess|AddUsers|更新学院人数:")
	} else {
		core.CheckErr(inerr, "usersDataAccess|AddUsers|添加用户:")
		//fmt.Println(inerr)
	}
	return inerr
}

//添加用户
func AddUsersInfoAll(us users.Users, st users.Students, tr users.Teacher, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(users.Students{}, "students")
	dbmap.AddTableWithName(users.Teacher{}, "teacher")
	inerr = dbmap.Insert(&us)
	if inerr == nil && us.Rolesid == 3 { //判断是否是学生
		st.Id = us.Id
		inerr = dbmap.Insert(&st)
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|添加学生详细信息:")
		_, inerr = dbmap.Exec("update classes set Classesnum=Classesnum+1 where Id=?;", st.Classesid) //更新班级人数
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|更新班级人数:")
		majorid, inerr := dbmap.SelectInt("select Majorid from classes where Id=?", st.Classesid) //获取专业ID
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|获取专业ID:")
		_, inerr = dbmap.Exec("update major set Majornum=Majornum+1 where Id=?;", majorid) //更新专业人数
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|更新专业人数:")
		collegeid, inerr := dbmap.SelectInt("select Collegeid from major where Id=?", majorid) //获取学院ID
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|获取学院ID:")
		_, inerr = dbmap.Exec("update college set Collegenum=Collegenum+1 where Id=?;", collegeid) //更新学院人数
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|更新学院人数:")
	} else if inerr == nil && us.Rolesid == 2 {
		tr.Id = us.Id
		inerr = dbmap.Insert(&tr)
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|添加教师详细信息:")
	} else {
		core.CheckErr(inerr, "usersDataAccess|AddUsersInfoAll|添加用户:")
	}
	return inerr
}

// 更新用户
func UpdateUsersInfoAll(us *users.Users, st *users.Students, tr *users.Teacher, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(users.Students{}, "students").SetKeys(true, "Id")
	dbmap.AddTableWithName(users.Teacher{}, "teacher").SetKeys(true, "Id")
	_, inerr = dbmap.Update(us)
	if inerr == nil && us.Rolesid == 3 { //判断是否是学生
		var oldst users.Students
		inerr = dbmap.SelectOne(&oldst, "select * from students where Id=?;", us.Id)
		core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|查询修改前的信息:")
		if oldst.Classesid != st.Classesid {
			_, inerr = dbmap.Exec("update classes set Classesnum=Classesnum-1 where Id=?;", oldst.Classesid) //更新班级人数
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|更新班级人数:")
			majorid, inerr := dbmap.SelectInt("select Majorid from classes where Id=?", oldst.Classesid) //获取专业ID
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|获取专业ID:")
			_, inerr = dbmap.Exec("update major set Majornum=Majornum-1 where Id=?;", majorid) //更新专业人数
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|更新专业人数:")
			collegeid, inerr := dbmap.SelectInt("select Collegeid from major where Id=?", majorid) //获取学院ID
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|获取学院ID:")
			_, inerr = dbmap.Exec("update college set Collegenum=Collegenum-1 where Id=?;", collegeid) //更新学院人数
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|更新学院人数:")
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|添加学生详细信息:")
			_, inerr = dbmap.Exec("update classes set Classesnum=Classesnum+1 where Id=?;", st.Classesid) //更新班级人数
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|更新班级人数:")
			majorid, inerr = dbmap.SelectInt("select Majorid from classes where Id=?", st.Classesid) //获取专业ID
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|获取专业ID:")
			_, inerr = dbmap.Exec("update major set Majornum=Majornum+1 where Id=?;", majorid) //更新专业人数
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|更新专业人数:")
			collegeid, inerr = dbmap.SelectInt("select Collegeid from major where Id=?", majorid) //获取学院ID
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|获取学院ID:")
			_, inerr = dbmap.Exec("update college set Collegenum=Collegenum+1 where Id=?;", collegeid) //更新学院人数
			core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|更新学院人数:")
		}
		//		updatesql := "update students set Enrollmentyear=?,Homeaddress=?,Nowaddress=?,Classesid=? where Id=?"
		//		_, inerr = dbmap.Exec(updatesql, st.Enrollmentyear, st.Homeaddress, st.Nowaddress, st.Classesid, st.Id)
		st.Id = oldst.Id
		_, inerr = dbmap.Update(st)
		core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|修改学生信息:")
	} else if inerr == nil && us.Rolesid == 2 {
		fmt.Println("1111")
		_, inerr = dbmap.Update(tr)
		core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|添加教师详细信息:")
	} else {
		fmt.Println("2222")
		core.CheckErr(inerr, "usersDataAccess|UpdateUsersInfoAll|添加用户:")
	}

	return inerr
}

//删除
func DelUsersInfoAll(us *users.Users, st *users.Students, tr *users.Teacher, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(users.Students{}, "students").SetKeys(true, "Id")
	dbmap.AddTableWithName(users.Teacher{}, "teacher").SetKeys(true, "Id")
	_, inerr = dbmap.Delete(us)
	if inerr == nil && us.Rolesid == 3 { //判断是否是学生
		var oldst users.Students
		inerr = dbmap.SelectOne(&oldst, "select * from students where Id=?;", us.Id)
		core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|查询修改前的信息:")
		if oldst.Classesid != st.Classesid {
			_, inerr = dbmap.Exec("update classes set Classesnum=Classesnum-1 where Id=?;", oldst.Classesid) //更新班级人数
			core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|更新班级人数:")
			majorid, inerr := dbmap.SelectInt("select Majorid from classes where Id=?", oldst.Classesid) //获取专业ID
			core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|获取专业ID:")
			_, inerr = dbmap.Exec("update major set Majornum=Majornum-1 where Id=?;", majorid) //更新专业人数
			core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|更新专业人数:")
			collegeid, inerr := dbmap.SelectInt("select Collegeid from major where Id=?", majorid) //获取学院ID
			core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|获取学院ID:")
			_, inerr = dbmap.Exec("update college set Collegenum=Collegenum-1 where Id=?;", collegeid) //更新学院人数
			core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|更新学院人数:")
		}
		_, inerr = dbmap.Delete(st)
	} else if inerr == nil && us.Rolesid == 2 {
		_, inerr = dbmap.Delete(tr)
		core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|删除教师详细信息")
	} else {
		core.CheckErr(inerr, "usersDataAccess|DelUsersInfoAll|删除用户")
	}
	return inerr
}

func CheckVaild(RolesId int, UsersId int, funcode string, dbmap *gorp.DbMap) (rd core.Returndata) {
	//	if UsersId == 29 || UsersId == 30 || UsersId == 26 || UsersId == 27 || UsersId == 81 || UsersId == 84 {
	//		rd.Rcode = "1000"
	//	} else {
	sql := "select count(*) from roles as rs inner join rolemodulecenter as rmc on rmc.Rolesid=rs.Id inner join rolemodulefunctioncenter as rmfc on rmfc.Rolemodulecenterid=rmc.Id"
	sql = sql + " inner join systemmodulefunctions as smf on smf.Id=rmfc.Systemmodulefunctionsid inner join users as us on us.Rolesid=rs.Id where rs.Id in(?) and smf.Functioncode=? and us.Id=?;"
	//fmt.Println(sql, RolesId, funcode, UsersId)
	countnum, counterr := dbmap.SelectInt(sql, RolesId, funcode, UsersId)
	core.CheckErr(counterr, "usersDataAccess|CheckVaild|验证用户是否有权限:")
	if counterr != nil {
		rd.Rcode = "1003"
		rd.Reason = "系统请求丢失，请联系管理员"
	} else {
		if countnum > 0 {
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1004"
			rd.Reason = "无此操作权限"
		}
	}
	//	}
	return rd
}

// 添加日志
func AddLoginlog(log *users.LoginLog, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(users.LoginLog{}, "loginlog").SetKeys(true, "Id")
	inerr = dbmap.Insert(log)
	core.CheckErr(inerr, "basicsetDataAccess|AddLoginlog|添加登录日志")
	return inerr
}

// 学生信息详细(查询分页)
func Query_StudentInfoDetail_ByPage(ws viewmodel.QueryUsersWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from users as us left join students as st on us.Id=st.Id left join teacher as tr on tr.Id=us.Id where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "usersDataAccess|QueryUsersPG|获取用户数量")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.StudentInfoDetail
		// 学生信息详细查询 SQL
		sql := `SELECT 	us.Id AS UserID,
						IFNULL(us.Loginuser, '') AS LoginUser, 
						IFNULL(us.Rolesid, 0) AS RoleID,
						IFNULL((SELECT roles.Rolesname FROM roles WHERE us.Rolesid = roles.Id), '') AS RoleName,
						IFNULL(us.Loginpwd, '') AS LoginPwd,
						IFNULL(us.Truename, '') AS TrueName,
						IFNULL(us.Nickname, '') AS NickName,
						IFNULL(us.Userheadimg, '') AS UserHeadImg,
						IFNULL(us.Userphone, '') AS UserPhone,
						IFNULL(us.Userstate, 0) AS UserState,
						IFNULL(us.Usersex, 0) AS UserSex,
						IFNULL(us.Usermac, '') AS Usermac,
						IFNULL(us.Birthday, '')AS Birthday,
					
						IFNULL(st.Id, 0) AS StudentID,
						IFNULL(st.Enrollmentyear, 0) AS EnrollmentYear,
						IFNULL(st.Homeaddress, '') AS HomeAddress,
						IFNULL(st.Nowaddress, '') AS NowAddress,
						IFNULL(st.Classesid, 0) AS ClasseID,
						IFNULL((SELECT classes.Classesname FROM classes WHERE classes.Id = st.Classesid),'') AS ClassName,
						IFNULL(st.Infostate, 0) AS InfoState,
						IFNULL(st.Currentstate, 0) AS Currentstate,
						IFNULL(st.Attendance, 0) AS Attendance,
						IFNULL(st.Needcoursenum, 0) AS NeedCourseNum,
						IFNULL(st.Alreadycoursenum, 0) AS AlreadyCourseNum,
						IFNULL(st.Absenteeism, 0) AS Absenteeism	
					FROM users AS us 
							LEFT JOIN students AS st ON us.Id = st.Id							
					WHERE 1 = 1  `
		sql = sql + wheresql + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "usersDataAccess|Query_StudentInfoDetail_ByPage|系统后台查询学生详细")
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

/******更新用户信息*********/
// 更新用户信息(密码除外)
func UpdateUser_ExceptPassword(us *users.Users, dbmap *gorp.DbMap) (err error) {
	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
	sql := `UPDATE users 
			SET  Loginuser = ?, 
				 Rolesid = ?,
				 Truename = ?,
				 Nickname = ?,
				 Userheadimg = ?,
				 Userphone = ?,
				 Userstate = ?,
				 Usersex = ?,
				 Usermac = ?,
				 Birthday = ?,
				 ThirdPartyId = ?,
				 Os = ?
			WHERE (Id = ?);`
	if _, err := dbmap.Exec(sql, us.Loginuser, us.Rolesid, us.Truename, us.Nickname, us.Userheadimg, us.Userphone, us.Userstate, us.Usersex, us.Usermac, us.Birthday, us.ThirdPartyId, us.Os); err != nil {
		xdebug.PrintStackTrace(err)
	}

	return err
}

// 更新用户密码
func UpdateUserPassword(user *users.Users, dbmap *gorp.DbMap) (err error) {
	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
	sql := `UPDATE users SET Loginpwd = ? WHERE (Id = ?);`
	if _, err := dbmap.Exec(sql, user.Loginpwd, user.Id); err != nil {
		xdebug.PrintStackTrace(err)
	}

	return err
}

func QueryUsersPG2(where users.QueryWhere4User, dbmap *gorp.DbMap) (reply core.Returndata) {
	var list []viewmodel.UsersInfoAll
	var pg = core.PageData{PageIndex: where.PageIndex, PageSize: where.PageSize, PageData: list}
	reply.Result = pg
	var sqlerr error
	countsql := `SELECT COUNT(*) FROM users AS TUser ` + where.WhereString()
	count, sqlerr := dbmap.SelectInt(countsql)
	if sqlerr != nil {
		reply.Rcode = strconv.Itoa(resp.FAIL.Code)
		reply.Reason = resp.FAIL.Text

		xdebug.PrintStackTrace(sqlerr)
		return reply
	}
	// 未找到数据
	if count <= 0 {
		reply.Rcode = strconv.Itoa(resp.FOUND_NODATA.Code)
		reply.Reason = resp.FOUND_NODATA.Text
		return reply
	}

	sql := `SELECT 	TUser.Id AS UsersId,
					IFNULL(TUser.Loginuser, '') AS Loginuser, 
					IFNULL(TUser.Rolesid, 0) AS Rolesid,
					IFNULL((SELECT roles.Rolesname FROM roles WHERE TUser.Rolesid = roles.Id), '') AS RoleName,
					IFNULL(TUser.Loginpwd, '') AS Loginpwd,
					IFNULL(TUser.Truename, '') AS Truename,
					IFNULL(TUser.Nickname, '') AS Nickname,
					IFNULL(TUser.Userheadimg, '') AS Userheadimg,
					IFNULL(TUser.Userphone, '') AS Userphone,
					IFNULL(TUser.Userstate, 0) AS Userstate,
					IFNULL(TUser.Usersex, 0) AS Usersex,
					IFNULL(TUser.Usermac, '') AS Usermac,
					IFNULL(TUser.Birthday, '') AS Birthday									
			FROM users AS TUser `
	sql += where.WhereString() + where.LimitString() + " ;"
	if _, sqlerr := dbmap.Select(&list, sql); sqlerr != nil {
		reply.Rcode = strconv.Itoa(resp.FAIL.Code)
		reply.Reason = resp.FAIL.Text

		xdebug.PrintStackTrace(sqlerr)
		return
	}

	pg.PageCount = int(count)
	pg.PageData = list
	reply.Result = pg
	reply.Rcode = strconv.Itoa(resp.SUCCESS.Code)
	reply.Reason = resp.SUCCESS.Text
	return reply
}

func init() {
	// 初始化响应消息集合
	resp = commons.ResponseMsgSet_Instance()
}
