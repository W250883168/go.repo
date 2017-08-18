package actiondataDataAccess

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"dev.project/BackEndCode/devserver/model/actiondata"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/curriculum"
	"dev.project/BackEndCode/devserver/viewmodel"

	"gopkg.in/gorp.v1"
)

/*
传入课程班级章节中间表ID:Curriculumclassroomchaptercentreid
*/
func QueryClassPointtos(lg viewmodel.GetPointtos, dbmap *gorp.DbMap) (list []viewmodel.PointtosUsers) {
	if lg.Curriculumclassroomchaptercentreid > 0 || lg.Curriculumclassroomchaptercentreids != "" { //判断传入进来的数值是否正确
		sql2 := "select ps.Usersid,us.Truename,us.Userheadimg,ps.state from curriculumclassroomchaptercentre as cccc inner join pointtos as ps on cccc.Id=ps.Curriculumclassroomchaptercentreid inner join users as us on us.Id=ps.Usersid where 1=1"
		wherestr := ""
		var ptlist []actiondata.Pointtos
		if lg.Curriculumclassroomchaptercentreid > 0 {
			wherestr = wherestr + " and cccc.Id=" + strconv.Itoa(lg.Curriculumclassroomchaptercentreid) + ";"
		}
		if lg.Curriculumclassroomchaptercentreids != "" {
			wherestr = wherestr + " and cccc.Id in(" + lg.Curriculumclassroomchaptercentreids + ");"
		}
		sql1 := fmt.Sprintf(`select st.Id as Usersid,cccc.Id as Curriculumclassroomchaptercentreid from Curriculumclassroomchaptercentre cccc inner join curriculumsclasscentre ccc on cccc.Curriculumsclasscentreid=ccc.Id
				inner join students st on st.Classesid=ccc.Classesid where 1=1 %s
				and st.Id not in (select ps.Usersid from curriculumclassroomchaptercentre as cccc inner join pointtos as ps on cccc.Id=ps.Curriculumclassroomchaptercentreid where 1=1 %s);
				;`, wherestr, wherestr)
		dbmap.Select(&ptlist, sql1)
		if len(ptlist) > 0 { //判断返回数据的长度
			for _, v := range ptlist {
				AddPointtos(&v, dbmap) //循环添加数据
			}
		}
		_, listerr := dbmap.Select(&list, (sql2 + wherestr)) //添加完成后在进行查询
		core.CheckErr(listerr, "actiondataDataAccess|QueryClassPointtos|查询班级点到的相关信息:")
		return list
	} else {
		return nil
	}
}

/*
教师点击移动中控获取教室Id
*/
func QueryClassClassroomId(lg core.BasicsToken, dbmap *gorp.DbMap) (rd core.Returndata) {
	if lg.Usersid > 0 { //判断传入进来的数值是否正确
		timestr := time.Now().Format("2006-01-02 15:04:05")
		sql2 := "select tr.Classroomid from curriculumclassroomchaptercentre as cccc inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid where (Enddate<=? or Begindate<=?) and Usersid=? order by Enddate desc limit 0,1;"
		fmt.Println(sql2, timestr, timestr, lg.Usersid)
		idint, iderr := dbmap.SelectInt(sql2, timestr, timestr, lg.Usersid)
		core.CheckErr(iderr, "actiondataDataAccess|QueryClassClassroomId|教师点击移动中控获取教室Id:")
		if iderr != nil {
			rd.Rcode = "1001"
			rd.Reason = "数据查询错误"
		} else {
			if idint == 0 {
				rd.Rcode = "1099"
				rd.Reason = "未找到相关记录"
			} else {
				rd.Rcode = "1000"
				rd.Result = idint
			}
		}
	} else {
		rd.Rcode = "2001"
		rd.Reason = "数据提交不正确"
	}
	return rd
}

/*
在无Id传入的情况下获取当前老师的状况下的课程章节点到
*/
func GetQueryCurriculumClassroomChapterCentreId(lg core.BasicsToken, dbmap *gorp.DbMap) (rd core.Returndata) {
	if lg.Usersid > 0 { //判断传入进来的数值是否正确
		timestr := time.Now().Format("2006-01-02 15:04:05")
		sql2 := "select Id from curriculumclassroomchaptercentre where (Enddate<=? or Begindate<=?) and Usersid=? order by Enddate desc limit 0,1;"
		idint, iderr := dbmap.SelectInt(sql2, timestr, timestr, lg.Usersid)
		core.CheckErr(iderr, "actiondataDataAccess|GetQueryCurriculumClassroomChapterCentreId|在无Id传入的情况下获取当前老师的状况下的课程章节点到:")
		if iderr != nil {
			rd.Rcode = "1001"
			rd.Reason = "数据查询错误"
		} else {
			if idint == 0 {
				rd.Rcode = "1099"
				rd.Reason = "未找到相关记录"
			} else {
				rd.Rcode = "1000"
				rd.Result = idint
			}
		}
	} else {
		rd.Rcode = "2001"
		rd.Reason = "数据提交不正确"
	}
	return rd
}

/*
提交课程班级章节中间表ID:Curriculumclassroomchaptercentreid,学生ID,状态[0:未到,1:已到]
*/
func UpdateStudentsPointtos(lg viewmodel.GetPointtos, Studentsid int, State int, dbmap *gorp.DbMap) (rd core.Returndata) {
	if (lg.Curriculumclassroomchaptercentreids != "") && Studentsid > 0 { //判断传入进来的数值是否正确
		var exeerror error
		timestr := time.Now().Format("2006-01-02 15:04:05")
		isokerrsql := "select count(*) from teachingrecord as tr inner join curriculumclassroomchaptercentre as cccc on cccc.Id=tr.Curriculumclassroomchaptercentreid where tr.Curriculumclassroomchaptercentreid in(" + lg.Curriculumclassroomchaptercentreids + ") and tr.State=0 and cccc.Begindate>='" + timestr + "';"
		countint, exeerror := dbmap.SelectInt(isokerrsql)
		if countint == 0 {
			sql := "update pointtos set State=?,Ismodify=1 where Curriculumclassroomchaptercentreid in(" + lg.Curriculumclassroomchaptercentreids + ") and Usersid=?;"
			fmt.Println(sql)
			_, exeerror = dbmap.Exec(sql, State, Studentsid)
			UpdateToclassrate(lg.Curriculumclassroomchaptercentreids, dbmap)
			core.CheckErr(exeerror, "actiondataDataAccess|UpdateStudentsPointtos|提交课程班级章节中间表ID:")
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "不能对还未上课的课程进行点到修改"
		}
	} else {
		rd.Rcode = "1001"
		rd.Reason = "提交数据格式不正确"
	}
	return rd
}

/*
更新出勤到课率
*/
func UpdateToclassrate(Curriculumclassroomchaptercentreids string, dbmap2 *gorp.DbMap) {
	var qdarr2 []curriculum.Curriculumclassroomchaptercentre
	//查询准备结束上课的教室
	dbmap2.Select(&qdarr2, "select * from curriculumclassroomchaptercentre where Id in("+Curriculumclassroomchaptercentreids+");")
	for p := 0; p < len(qdarr2); p++ {
		//更新课程班级章节id的到课率和到课人数
		count2, _ := dbmap2.SelectInt("select count(*) from pointtos where Curriculumclassroomchaptercentreid=? and state>=1;", qdarr2[p].Id)
		Toclassrate := float32(float32(count2) / float32(qdarr2[p].Plannumber)) //到课率
		dbmap2.Exec("update curriculumclassroomchaptercentre set Actualnumber=?,Toclassrate=? where Id=?;", int(count2), Toclassrate, qdarr2[p].Id)
		//		CheckErr(roomerr3, "update curriculumclassroomchaptercentre set Classroomstate='0' where Id=?;")
		//		//获取课程下的章节平均到课率
		avgToclassrate, _ := dbmap2.SelectFloat("select avg(Toclassrate) from curriculumclassroomchaptercentre as cccc inner join teachingrecord as td on cccc.Id=td.Curriculumclassroomchaptercentreid where cccc.Curriculumsclasscentreid=? and td.state>=1;", qdarr2[p].Id)
		count, _ := dbmap2.SelectInt("select count(*) from curriculumclassroomchaptercentre as cccct left join curriculumsclasscentre as ccc on cccct.Curriculumsclasscentreid=ccc.Id inner join teachingrecord as trd on trd.Curriculumclassroomchaptercentreid=cccct.Id where cccct.Curriculumsclasscentreid=? and trd.state=2;", qdarr2[p].Curriculumsclasscentreid)
		livecount, _ := dbmap2.SelectInt("select count(*) from curriculumclassroomchaptercentre as cccct left join curriculumsclasscentre as ccc on cccct.Curriculumsclasscentreid=ccc.Id inner join teachingrecord as trd on trd.Curriculumclassroomchaptercentreid=cccct.Id where cccct.Curriculumsclasscentreid=? and trd.state=2 and(cccct.Islive=1 or cccct.Isondomian=1);", qdarr2[p].Curriculumsclasscentreid)
		dbmap2.Exec("update curriculumsclasscentre set Newchapter=?,Newlivechapter=?,Averageclassrate=? where Id=?", int(count), int(livecount), float32(avgToclassrate), qdarr2[p].Curriculumsclasscentreid)
	}
}

/*
更改上课状态
传入数据有 教室Id，教师Id，上课状态
*/
func ChangeClassState(Classroomid int, Usersid int, State int, Ccccid int, Ccccids string, dbmap *gorp.DbMap) (rd core.Returndata) {
	if Classroomid > 0 && Usersid > 0 { //
		datestr := core.Timeaction(time.Now().Format("2006-01-02 15:04:05"))
		sql1 := "select cccc.Id as Ccccid,tr.Classroomid from curriculumclassroomchaptercentre as cccc inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid where tr.Classroomid=? and cccc.Usersid=? "
		if Ccccid > 0 {
			sql1 = sql1 + " and tr.State<=" + strconv.Itoa(State) + " and cccc.Id=" + strconv.Itoa(Ccccid) + " ;" //limit 0,1
		} else if Ccccids != "" {
			sql1 = sql1 + " and tr.State<=" + strconv.Itoa(State) + " and cccc.Id in (" + Ccccids + ") ;"
		} else {
			sql1 = sql1 + " and cccc.Begindate<='" + datestr + "' and cccc.Enddate>='" + datestr + "'" + " ;" //limit 0,1
		}
		type CcccIds struct {
			Ccccid      int
			Classroomid int
		}
		var cccclist []CcccIds
		_, Iderr := dbmap.Select(&cccclist, sql1, Classroomid, Usersid)
		//		Id, Iderr := dbmap.SelectInt(sql1, Classroomid, Usersid)
		core.CheckErr(Iderr, "actiondataDataAccess|ChangeClassState|更改上课状态:")
		if len(cccclist) > 0 && Iderr == nil {
			for _, v := range cccclist {
				valstate := State
				_, trerr1 := dbmap.Exec("update teachingrecord set State=? where curriculumclassroomchaptercentreid=?;", strconv.Itoa((valstate + 1)), v.Ccccid)
				core.CheckErr(trerr1, "actiondataDataAccess|ChangeClassState|更改上课状态|update teachingrecord set State=? where Id=?;")
				if State == 0 {
					valstate = 1
				} else if State == 1 {
					valstate = 0
				}
				_, roomerr1 := dbmap.Exec("update classrooms set Classroomstate=? where Id=?;", strconv.Itoa(valstate), v.Classroomid)
				core.CheckErr(roomerr1, "actiondataDataAccess|ChangeClassState|更改上课状态|update classrooms set Classroomstate=? where Id=?;")
			}
			rd.Result = cccclist
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1003"
			rd.Reason = "系统后台未找到相关数据"
		}
	} else {
		fmt.Println("更改上课状态:", Classroomid, Usersid)
		rd.Rcode = "1002"
		rd.Reason = "数据提交不准确"
	}
	return rd
}

/*
批量提交课程班级章节中间表ID:Curriculumclassroomchaptercentreid,学生ID,状态[0:未到,1:已到]
*/
func UpdateListStudentsPointtos(updata viewmodel.PostUpdatePointtosData, dbmap *gorp.DbMap) (rd core.Returndata) {
	if len(updata.CcccIds) > 0 && len(updata.CcccIds) == len(updata.StudentsIds) && len(updata.States) == len(updata.CcccIds) {
		var CcccIdstr []string
		for _, v := range updata.CcccIds {
			CcccIdstr = append(CcccIdstr, strconv.Itoa(v))
		}
		CcccIdstrs := strings.Join(CcccIdstr, ",")
		timestr := time.Now().Format("2006-01-02 15:04:05")
		isokerrsql := "select count(*) from teachingrecord as tr inner join curriculumclassroomchaptercentre as cccc on cccc.Id=tr.Curriculumclassroomchaptercentreid where tr.Curriculumclassroomchaptercentreid in(" + CcccIdstrs + ") and tr.State=0 and cccc.Begindate>='" + timestr + "';"
		countint, _ := dbmap.SelectInt(isokerrsql)
		if countint == 0 {
			for k := 0; k < len(updata.CcccIds); k++ {
				if updata.CcccIds[k] > 0 && updata.StudentsIds[k] > 0 { //判断传入进来的数值是否正确
					_, exeerror := dbmap.Exec("update pointtos set State=?,Ismodify=1 where Curriculumclassroomchaptercentreid=? and Usersid=?;", updata.States[k], updata.CcccIds[k], updata.StudentsIds[k])
					core.CheckErr(exeerror, "actiondataDataAccess|UpdateListStudentsPointtos|批量提交课程班级章节中间表:")
					if exeerror != nil {
						rd.Rcode = rd.Rcode + "|1001"
						rd.Reason = rd.Reason + "|数据[" + strconv.Itoa(updata.CcccIds[k]) + "][" + strconv.Itoa(updata.StudentsIds[k]) + "]操作执行失败"
					}
				} else {
					rd.Rcode = rd.Rcode + "|1002"
					rd.Reason = rd.Reason + "|数据[" + strconv.Itoa(updata.CcccIds[k]) + "][" + strconv.Itoa(updata.StudentsIds[k]) + "]提交不正确"
				}
			}
		} else {
			rd.Rcode = "1001"
			rd.Reason = "不能对还未上课的课程进行点到修改"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交不正确"
	}
	if rd.Rcode == "" {
		rd.Rcode = "1000"
	}
	return rd
}

/*
出勤数据--实时出勤
*/
func QueryAttendancelist(lg viewmodel.PostQueryCurriculums, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) { //(list []viewmodel.QueryPeoples) {
	var timestr = time.Now().Format("2006-01-02 15:04:05")
	var list []viewmodel.QueryPeoples
	sql1 := "select count(pts.Id) as Seatsnumbers,ifnull(sum(pts.State),0)as Sumnumbers,crs.Id as Classroomid,bd.Id as BuildingId,bd.Buildingname as BuildingName,fls.Id as FloorId,fls.Floorname as FloorName,fls.FloorsImage,crs.Classroomsname as ClassroomName,crs.Classroomstate as ClassroomState,crs.Classroomstype,crs.Classroomicon from "
	sql1 = sql1 + " classrooms as crs inner join floors as fls on crs.Floorsid=fls.Id inner join building as bd on fls.Buildingid=bd.Id inner join campus as cps on bd.Campusid=cps.Id "
	sql1 = sql1 + " left join teachingrecord as tr on tr.Classroomid=crs.Id left join curriculumclassroomchaptercentre as cccc on (cccc.Id=tr.Curriculumclassroomchaptercentreid and cccc.Begindate<='" + timestr + "' and cccc.Enddate>='" + timestr + "')"
	sql1 = sql1 + " left join pointtos as pts on (pts.Curriculumclassroomchaptercentreid=cccc.Id) where 1=1"
	if lg.Campusids != "" {
		sql1 = sql1 + " and cps.Id in(" + lg.Campusids + ")"
	}
	if lg.Buildingids != "" {
		sql1 = sql1 + " and bd.Id in(" + lg.Buildingids + ")"
	}
	if lg.Floorsids != "" {
		sql1 = sql1 + " and fls.Id in(" + lg.Floorsids + ")"
	}
	sql1 = sql1 + " group by bd.Id,fls.Id,crs.Id,crs.Classroomsname,crs.Classroomstate,crs.Classroomstype,crs.Classroomicon;"
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "actiondataDataAccess|QueryAttendancelist|根据校区Id、楼栋id、楼层id、查询教室出勤数据")
	rd.Rcode = "1000"
	rd.Result = list
	return rd
}

/*
获取课表数据
根据角色不同查看的数据也不同
学生：只看班级课程数据
教师：只看安排自己将要上的课程数据
管理者：看所有的数据
请求数据
[开始时间、结束时间、Token令牌、用户ID、角色ID]
响应数据
[课程名称，课程ID，章节名称，章节ID，开始时间，结束时间，所在校区，所在楼栋，所在楼层，所在教室,课程状态]
*/
func QueryCurriculumsTable(lg viewmodel.PostQueryCurriculums, bt core.BasicsToken, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) { //(list []viewmodel.GetCurriculumslist) {
	var list []viewmodel.GetCurriculumslist
	if bt.Rolestype > 0 { //判断账号是否正确
		sql := "select cccc.Id as Curriculumclassroomchaptercentreid,cc.Curriculumname,cc.Id as Curriculumsid,ct.Id as Chaptersid,ct.Chaptername,cccc.Begindate,cccc.Enddate,cs.Classesname,crs.Classroomsname,fls.Floorname,bd.Buildingname,cps.Campusname,tr.state,mr.Majorname,cg.Collegename,us.Nickname,cccc.Usersid as TeacherId"
		sql = sql + ",cccc.Plannumber,cccc.Actualnumber,cccc.Toclassrate"
		sqlfrom := " from curriculums as cc inner join chapters as ct on cc.Id=ct.Curriculumsid"
		sqlfrom = sqlfrom + " inner join curriculumsclasscentre as ccc on cc.Id=ccc.Curriculumsid inner join curriculumclassroomchaptercentre as cccc on (cccc.Curriculumsclasscentreid=ccc.Id and ct.Id=cccc.Chaptersid)"
		sqlfrom = sqlfrom + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id inner join classrooms as crs on crs.Id= tr.Classroomid"
		sqlfrom = sqlfrom + " inner join floors as fls on fls.Id=crs.Floorsid inner join building as bd on bd.Id=fls.Buildingid inner join campus as cps on cps.Id=bd.Campusid inner join classes as cs on ccc.Classesid=cs.Id"
		sqlfrom = sqlfrom + " inner join major as mr on mr.Id=cs.Majorid inner join college as cg on cg.Id=mr.Collegeid inner join users as us on us.Id=cccc.Usersid where 1=1 "
		wheresql := ""
		ordersql := " order by cccc.Begindate"
		sql = sql + sqlfrom
		switch bt.Rolestype {
		case 1: //管理员
			break
		case 2: //老师
			wheresql = wheresql + " and cccc.Usersid=" + strconv.Itoa(bt.Usersid)
			break
		case 3: //学生
			wheresql = wheresql + " and ccc.Classesid in (select Classesid from students where Id=" + strconv.Itoa(bt.Usersid) + ")"
		}
		if lg.State > -1 { //判断课程状态
			wheresql = wheresql + " and tr.state=" + strconv.Itoa(lg.State)
		}
		if lg.Begindate != "" {
			wheresql = wheresql + " and Begindate>='" + lg.Begindate + "'"
		}
		if lg.Enddate != "" {
			wheresql = wheresql + " and Enddate<='" + lg.Enddate + "'"
		}
		if lg.Teacherids != "" {
			wheresql = wheresql + " and cccc.Usersid in(" + lg.Teacherids + ")"
		}
		if lg.Collegeids != "" {
			wheresql = wheresql + " and Collegeid in(" + lg.Collegeids + ")"
		}
		if lg.Majorids != "" {
			wheresql = wheresql + " and Majorid in(" + lg.Majorids + ")"
		}
		if lg.Classesids != "" {
			wheresql = wheresql + " and ccc.Classesid in(" + lg.Classesids + ")"
		}
		if lg.Curriculumsids != "" {
			wheresql = wheresql + " and ct.Curriculumsid in(" + lg.Curriculumsids + ")"
		}
		if lg.Searhtxt != "" {
			wheresql = wheresql + " and cc.Curriculumname like '%" + lg.Searhtxt + "%'"
		}
		if lg.Campusid > 0 {
			wheresql = wheresql + " and bd.Campusid=" + strconv.Itoa(lg.Campusid)
		}
		if lg.Buildingid > 0 {
			wheresql = wheresql + " and fls.Buildingid=" + strconv.Itoa(lg.Buildingid)
		}
		if lg.Floorsid > 0 {
			wheresql = wheresql + " and crs.Floorsid=" + strconv.Itoa(lg.Floorsid)
		}
		if lg.Classroomid > 0 {
			wheresql = wheresql + " and crs.Id=" + strconv.Itoa(lg.Classroomid)
		}
		if pg.PageIndex > 0 {
			sqlconut := "select count(*) "
			sqlconut = sqlconut + sqlfrom
			count, errs2 := dbmap.SelectInt(sqlconut + wheresql)
			core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsTable|根据楼栋查看教室的出勤情况|"+sqlconut+wheresql)
			if count > 0 && errs2 == nil {
				pg.PageCount = int(count)
				_, errs2 = dbmap.Select(&list, sql+wheresql+ordersql+core.GetLimitString(pg)) //查询权限模块
				pg.PageData = list
				rd.Result = pg
				rd.Rcode = "1000"
			} else {
				rd.Rcode = "1099"
				rd.Reason = "未搜索到数据"
				rd.Result = pg
				core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsTable|根据楼栋查看教室的出勤情况|"+sql+wheresql+ordersql+core.GetLimitString(pg))
			}
		} else {
			_, errs2 := dbmap.Select(&list, sql+wheresql+ordersql+core.GetLimitString(pg)) //查询权限模块
			core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsTable|获取课表数据|"+sql+wheresql+ordersql+core.GetLimitString(pg))
			rd.Rcode = "1000"
			rd.Result = list
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "查询失败"
	}
	return rd
}

/*
查询历史出勤记录
*/
func QueryHistoryAttendance(lg viewmodel.PostQueryCurriculums, bt core.BasicsToken, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) { //(list []viewmodel.GetCurriculumslist) {
	var list []viewmodel.GetCurriculumslist
	if bt.Rolestype > 0 { //判断账号是否正确
		sql := "select cccc.Id as Curriculumclassroomchaptercentreid,cc.Curriculumname,cc.Id as Curriculumsid,ct.Id as Chaptersid,ct.Chaptername,cccc.Begindate,cccc.Enddate,cs.Classesname,crs.Classroomsname,fls.Floorname,bd.Buildingname,cps.Campusname,tr.state,mr.Majorname,cg.Collegename,us.Nickname"
		sql = sql + ",cccc.Plannumber,cccc.Actualnumber,cccc.Toclassrate"
		sqlfrom := " from curriculums as cc inner join chapters as ct on cc.Id=ct.Curriculumsid"
		sqlfrom = sqlfrom + " inner join curriculumsclasscentre as ccc on cc.Id=ccc.Curriculumsid inner join curriculumclassroomchaptercentre as cccc on (cccc.Curriculumsclasscentreid=ccc.Id and ct.Id=cccc.Chaptersid)"
		sqlfrom = sqlfrom + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id inner join classrooms as crs on crs.Id= tr.Classroomid"
		sqlfrom = sqlfrom + " inner join floors as fls on fls.Id=crs.Floorsid inner join building as bd on bd.Id=fls.Buildingid inner join campus as cps on cps.Id=bd.Campusid inner join classes as cs on ccc.Classesid=cs.Id"
		sqlfrom = sqlfrom + " inner join major as mr on mr.Id=cs.Majorid inner join college as cg on cg.Id=mr.Collegeid inner join users as us on us.Id=cccc.Usersid where 1=1 "
		wheresql := ""
		ordersql := " order by cccc.Begindate"
		sql = sql + sqlfrom
		switch bt.Rolestype {
		case 1: //管理员
			break
		case 2: //老师
			wheresql = wheresql + " and cccc.Usersid=" + strconv.Itoa(bt.Usersid)
			break
		case 3: //学生
			wheresql = wheresql + " and ccc.Classesid in (select Classesid from students where Id=" + strconv.Itoa(bt.Usersid) + ")"
		}
		if lg.State > -1 { //判断课程状态
			wheresql = wheresql + " and tr.state=" + strconv.Itoa(lg.State)
		}
		if lg.Begindate != "" {
			wheresql = wheresql + " and Begindate>='" + lg.Begindate + "'"
		}
		if lg.Enddate != "" {
			wheresql = wheresql + " and Enddate<='" + lg.Enddate + "'"
		}
		if lg.Teacherids != "" {
			wheresql = wheresql + " and cccc.Usersid in(" + lg.Teacherids + ")"
		}
		if lg.Collegeids != "" {
			wheresql = wheresql + " and Collegeid in(" + lg.Collegeids + ")"
		}
		if lg.Majorids != "" {
			wheresql = wheresql + " and Majorid in(" + lg.Majorids + ")"
		}
		if lg.Classesids != "" {
			wheresql = wheresql + " and ccc.Classesid in(" + lg.Classesids + ")"
		}
		if lg.Curriculumsids != "" {
			wheresql = wheresql + " and ct.Curriculumsid in(" + lg.Curriculumsids + ")"
		}
		if lg.Searhtxt != "" {
			wheresql = wheresql + " and (cc.Curriculumname like '%" + lg.Searhtxt + "%' or mr.Majorname like '%" + lg.Searhtxt + "%' or cg.Collegename like '%" + lg.Searhtxt + "%' or us.Nickname like '%" + lg.Searhtxt + "%'"
			wheresql = wheresql + " or ct.Chaptername like '%" + lg.Searhtxt + "%' or cs.Classesname like '%" + lg.Searhtxt + "%' or crs.Classroomsname like '%" + lg.Searhtxt + "%' or bd.Buildingname like '%" + lg.Searhtxt + "%')"
		}
		if lg.Campusid > 0 {
			wheresql = wheresql + " and bd.Campusid=" + strconv.Itoa(lg.Campusid)
		}
		if lg.Buildingid > 0 {
			wheresql = wheresql + " and fls.Buildingid=" + strconv.Itoa(lg.Buildingid)
		}
		if lg.Floorsid > 0 {
			wheresql = wheresql + " and crs.Floorsid=" + strconv.Itoa(lg.Floorsid)
		}
		if lg.Classroomid > 0 {
			wheresql = wheresql + " and crs.Id=" + strconv.Itoa(lg.Classroomid)
		}
		if pg.PageIndex > 0 {
			sqlconut := "select count(*) "
			sqlconut = sqlconut + sqlfrom
			count, errs2 := dbmap.SelectInt(sqlconut + wheresql)
			core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsTable|根据楼栋查看教室的出勤情况|"+sqlconut+wheresql)
			if count > 0 && errs2 == nil {
				pg.PageCount = int(count)
				_, errs2 = dbmap.Select(&list, sql+wheresql+ordersql+core.GetLimitString(pg)) //查询权限模块
				pg.PageData = list
				rd.Result = pg
				rd.Rcode = "1000"
				core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsTable|根据楼栋查看教室的出勤情况|"+sql+wheresql+ordersql+core.GetLimitString(pg))
			} else {
				rd.Rcode = "1099"
				rd.Reason = "查询失败"
				rd.Result = pg
			}
		} else {
			fmt.Println(sql + wheresql + ordersql + core.GetLimitString(pg))
			_, errs2 := dbmap.Select(&list, sql+wheresql+ordersql+core.GetLimitString(pg)) //查询权限模块
			core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsTable|获取课表数据|"+sql+wheresql+ordersql+core.GetLimitString(pg))
			rd.Rcode = "1000"
			rd.Result = list
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "查询失败"
	}
	return rd
}
func QueryCurriculumsinfo(lg viewmodel.PostQueryCurriculums, dbmap *gorp.DbMap) (rd core.Returndata) { //(list []viewmodel.GetCurriculumslist) {
	var obj viewmodel.GetCurriculumslist
	sql := "select cccc.Id as Curriculumclassroomchaptercentreid,cc.Curriculumname,cc.Id as Curriculumsid,ct.Id as Chaptersid,ct.Chaptername,cccc.Begindate,cccc.Enddate,cs.Classesname,crs.Classroomsname,fls.Floorname,bd.Buildingname,cps.Campusname,tr.state,mr.Majorname,cg.Collegename,us.Nickname"
	sql = sql + ",cccc.Plannumber,cccc.Actualnumber,cccc.Toclassrate"
	sqlfrom := " from curriculums as cc inner join chapters as ct on cc.Id=ct.Curriculumsid"
	sqlfrom = sqlfrom + " inner join curriculumsclasscentre as ccc on cc.Id=ccc.Curriculumsid inner join curriculumclassroomchaptercentre as cccc on (cccc.Curriculumsclasscentreid=ccc.Id and ct.Id=cccc.Chaptersid)"
	sqlfrom = sqlfrom + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id inner join classrooms as crs on crs.Id= tr.Classroomid"
	sqlfrom = sqlfrom + " inner join floors as fls on fls.Id=crs.Floorsid inner join building as bd on bd.Id=fls.Buildingid inner join campus as cps on cps.Id=bd.Campusid inner join classes as cs on ccc.Classesid=cs.Id"
	sqlfrom = sqlfrom + " inner join major as mr on mr.Id=cs.Majorid inner join college as cg on cg.Id=mr.Collegeid inner join users as us on us.Id=ccc.Usersid where 1=1 "
	wheresql := " and cccc.Id=" + strconv.Itoa(lg.Curriculumclassroomchaptercentreid)
	ordersql := " order by cccc.Begindate"
	sql = sql + sqlfrom
	fmt.Println(sql + wheresql + ordersql + ";")
	errs2 := dbmap.SelectOne(&obj, sql+wheresql+ordersql+";") //查询权限模块
	core.CheckErr(errs2, "actiondataDataAccess|QueryCurriculumsinfo|获取课表数据|"+sql+wheresql+ordersql)
	rd.Rcode = "1000"
	rd.Result = obj
	return rd
}

/*
设置或者取消教室收藏记录
*/
func SetOrCancelClassroomCollection(lg viewmodel.PostCollection, bt core.BasicsToken, dbmap *gorp.DbMap) bool {
	var isok error
	var iscount int64
	fmt.Println(lg.State)
	if lg.State == 1 { //取消[删除数据]
		_, isok = dbmap.Exec("delete from Classroomcollection where Usersid=? and Classroomid=?;", bt.Usersid, lg.Classroomid)
		core.CheckErr(isok, "actiondataDataAccess|SetOrCancelClassroomCollection|设置或者取消教室收藏记录|删除数据")
	} else { //添加
		iscount, isok = dbmap.SelectInt("select count(*) from Classroomcollection where Usersid=? and Classroomid=?;", bt.Usersid, lg.Classroomid)
		core.CheckErr(isok, "actiondataDataAccess|SetOrCancelClassroomCollection|设置或者取消教室收藏记录|添加数据")
		if iscount == 0 {
			cc := actiondata.Classroomcollection{Classroomid: lg.Classroomid, Usersid: bt.Usersid, Createdate: time.Now().Format("2006-01-02 15:04:05")}
			dbmap.AddTableWithName(actiondata.Classroomcollection{}, "classroomcollection").SetKeys(true, "Id")
			isok = dbmap.Insert(&cc)
			core.CheckErr(isok, "actiondataDataAccess|SetOrCancelClassroomCollection|设置或者取消教室收藏记录|更新教室收藏数")
		}
	}
	iscount, isok = dbmap.SelectInt("select count(*) from Classroomcollection where Classroomid=?;", lg.Classroomid)
	core.CheckErr(isok, "actiondataDataAccess|SetOrCancelClassroomCollection|设置或者取消教室收藏记录|更新教室收藏数")
	if isok == nil {
		dbmap.Exec("update classrooms set Collectionnumbers=? where Id=?;", iscount, lg.Classroomid)
		return true
	} else {
		return false
	}
}

/*
查询我关注的课程数据
*/
func QueryMyAttentionRecord(lg viewmodel.QueryAttentionRecordWhere, bt core.BasicsToken, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.QueryAttentionRecordList) {
	sql1 := "select ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,ifnull(ccc.Classesid,0)as Classesid,cs.Classesname"
	sql1 = sql1 + " from attentionrecord as atr inner join curriculumsclasscentre as ccc on (atr.Classesid=ccc.Classesid and atr.Curriculumsid=ccc.Curriculumsid) inner join curriculums as cc on cc.Id=ccc.Curriculumsid "
	sql1 = sql1 + " left join students as sts on sts.Id=atr.Usersid inner join classes as cs on cs.Id=ccc.Classesid where 1=1"
	sql1 = sql1 + GetQueryAttentionRecordWhere(bt, lg)
	if bt.Rolestype == 2 {
		sql1 = strings.Replace(sql1, "ccc.Usersid", "atr.Usersid", 1)
		//sql1 = sql1 + " and atr.Usersid=" + strconv.Itoa(lg.Usersid)
	}
	sql1 = sql1 + core.GetLimitString(pg)
	fmt.Println(sql1)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "actiondataDataAccess|QueryMyAttentionRecord|查询我关注的课程数据")
	return list
}

/*
查询我所在的班级关注的课程
*/
func QueryMyClassAttentionRecord(lg viewmodel.QueryAttentionRecordWhere, bt core.BasicsToken, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.QueryAttentionRecordList) {
	sql1 := "select ccc.Id as Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,sts.Classesid,cs.Classesname"
	sql1 = sql1 + " from students as sts inner join curriculumsclasscentre as ccc on sts.Classesid=ccc.Classesid inner join curriculums as cc on cc.Id=ccc.Curriculumsid inner join classes as cs on cs.Id=ccc.Classesid inner join curriculumclassroomchaptercentre cccc on cccc.Curriculumsclasscentreid=ccc.Id where 1=1"
	sql1 = sql1 + GetQueryAttentionRecordWheres(bt, lg) + " group by ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,sts.Classesid,cs.Classesname " + core.GetLimitString(pg)
	fmt.Println(sql1)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "actiondataDataAccess|QueryMyClassAttentionRecord|查询我所在的班级关注的课程")
	return list
}

//判断条件获取
func GetQueryAttentionRecordWheres(bt core.BasicsToken, lg viewmodel.QueryAttentionRecordWhere) (where string) {
	if lg.Islive > 0 {
		where = where + " and ccc.Islive=" + strconv.Itoa(lg.Islive)
	}
	if lg.Isondemand > 0 {
		where = where + " and cccc.Isondomian=" + strconv.Itoa(lg.Isondemand)
	}
	if bt.Usersid > 0 {
		if bt.Rolestype == 3 {
			where = where + " and sts.Id=" + strconv.Itoa(bt.Usersid)
		} else if bt.Rolestype == 2 {
			where = where + " and ccc.Usersid=" + strconv.Itoa(bt.Usersid)
		}
	}
	if lg.Subjectcode != "" {
		where = where + " and left(cc.Subjectcode,LENGTH('" + lg.Subjectcode + "'))='" + lg.Subjectcode + "'"
	}
	return where
}

//判断条件获取
func GetQueryAttentionRecordWhere(bt core.BasicsToken, lg viewmodel.QueryAttentionRecordWhere) (where string) {
	if lg.Islive > 0 {
		where = where + " and ccc.Islive=" + strconv.Itoa(lg.Islive)
	}
	if lg.Isondemand > 0 {
		where = where + " and ccc.Isondemand=" + strconv.Itoa(lg.Isondemand)
	}
	if bt.Usersid > 0 {
		if bt.Rolestype == 3 {
			where = where + " and sts.Id=" + strconv.Itoa(bt.Usersid)
		} else if bt.Rolestype == 2 {
			where = where + " and ccc.Usersid=" + strconv.Itoa(bt.Usersid)
		}
	}
	if lg.Subjectcode != "" {
		where = where + " and left(cc.Subjectcode,LENGTH('" + lg.Subjectcode + "'))='" + lg.Subjectcode + "'"
	}
	return where
}

/*
关注或取消关注
*/
func SetAttentionRecord(lg *actiondata.Attentionrecord, dbmap *gorp.DbMap) (inerrs error) {
	var inerr error
	if lg.Curriculumsid > 0 && lg.Classesid > 0 && lg.Usersid > 0 {
		dbmap.AddTableWithName(actiondata.Attentionrecord{}, "attentionrecord").SetKeys(true, "Id")
		var at actiondata.Attentionrecord
		inerr = dbmap.SelectOne(&at, "select * from attentionrecord where Usersid=? and Curriculumsid=? and Classesid=?;", lg.Usersid, lg.Curriculumsid, lg.Classesid)
		core.CheckErr(inerr, "actiondataDataAccess|SetAttentionRecord|关注或取消关注|查询数据是否存在")
		if at.Id > 0 {
			if lg.State == 0 { //取消
				_, inerr = dbmap.Exec("delete from attentionrecord where Usersid=? and Curriculumsid=? and Classesid=?;", lg.Usersid, lg.Curriculumsid, lg.Classesid)
				core.CheckErr(inerr, "actiondataDataAccess|SetAttentionRecord|关注或取消关注|取消")
				if inerr == nil {
					_, inerr = dbmap.Exec("update curriculumsclasscentre set FollowSumnum=FollowSumnum-1 where Curriculumsid=? and Classesid=?;", lg.Curriculumsid, lg.Classesid)
					core.CheckErr(inerr, "actiondataDataAccess|SetAttentionRecord|关注或取消关注|修改")
				}
			} else {
				lg.State = 2
				_, inerr = dbmap.Update(&at)
				core.CheckErr(inerr, "actiondataDataAccess|SetAttentionRecord|关注或取消关注|修改关注状态")
			}
		} else {
			inerr = dbmap.Insert(lg)
			core.CheckErr(inerr, "actiondataDataAccess|SetAttentionRecord|关注或取消关注|添加关注状态")
			if inerr == nil {
				_, inerr = dbmap.Exec("update curriculumsclasscentre set FollowSumnum=FollowSumnum+1 where Curriculumsid=? and Classesid=?;", lg.Curriculumsid, lg.Classesid)
			}
		}
	}
	fmt.Println(inerr)
	return inerr
}

/*
判断是否已关注
*/
func IsAttentionRecordOk(lg *actiondata.Attentionrecord, dbmap *gorp.DbMap) (b bool) {
	b = false
	if lg.Curriculumsid > 0 && lg.Classesid > 0 && lg.Usersid > 0 {
		dbmap.AddTableWithName(actiondata.Attentionrecord{}, "attentionrecord").SetKeys(true, "Id")
		//var at actiondata.Attentionrecord
		lgerr := dbmap.SelectOne(lg, "select * from attentionrecord where Usersid=? and Curriculumsid=? and Classesid=?;", lg.Usersid, lg.Curriculumsid, lg.Classesid)
		core.CheckErr(lgerr, "actiondataDataAccess|IsAttentionRecordOk|判断是否已关注")
		if lg.Id > 0 {
			b = true
		}
	}
	return b
}

//查询用户所有的教室收藏记录
func QueryUserCollection(Usersid int, dbmap *gorp.DbMap) (list []actiondata.Classroomcollection) {
	sql1 := "select * from classroomcollection where Usersid="
	_, sserr1 := dbmap.Select(&list, sql1+strconv.Itoa(Usersid)+";")
	core.CheckErr(sserr1, "actiondataDataAccess|QueryUserCollection|查询用户所有的教室收藏记录")
	return list
}

//添加预置点到的数据
func AddPointtos(cr *actiondata.Pointtos, dbmap *gorp.DbMap) (inerr error) {
	if cr.Curriculumclassroomchaptercentreid > 0 && cr.Usersid > 0 {
		dbmap.AddTableWithName(actiondata.Pointtos{}, "pointtos").SetKeys(true, "Id")
		inerr = dbmap.Insert(cr)
		core.CheckErr(inerr, "actiondataDataAccess|AddPointtos|添加预置点到的数据")
	}
	return inerr
}

//添加上课的教室对应关系数据
func AddTeachingrecord(cr *actiondata.Teachingrecord, dbmap *gorp.DbMap) (inerr error) {
	if cr.Curriculumclassroomchaptercentreid > 0 && cr.Classroomid > 0 {
		dbmap.AddTableWithName(actiondata.Teachingrecord{}, "teachingrecord").SetKeys(true, "Id")
		inerr = dbmap.Insert(cr)
		core.CheckErr(inerr, "actiondataDataAccess|AddTeachingrecord|添加上课的教室对应关系数据")
	}
	return inerr
}
