package curriculumDataAccess

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"dev.project/BackEndCode/devserver/model/actiondata"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/curriculum"
	"dev.project/BackEndCode/devserver/model/live"
	"dev.project/BackEndCode/devserver/viewmodel"

	"gopkg.in/gorp.v1"
)

func GetWhereString(ws viewmodel.QueryCurriculumWhere) (where string) {

	if ws.Subjectcode != "" {
		where = where + " and sc.Subjectcode='" + ws.Subjectcode + "'"
	}
	if ws.Subjectcode != "" {
		where = where + " and (sc.Subjectcode like '%" + ws.Subjectcode + "%' or sc.Subjectname like '%" + ws.Subjectname + "%')"
	}
	if ws.Curriculumname != "" {
		where = where + " and Curriculumname like '%" + ws.Curriculumname + "%'"
	}
	if ws.Seacrchtxt != "" {
		where = where + " and(Curriculumname like '%" + ws.Seacrchtxt + "%' or Curriculumsdetails like '%" + ws.Seacrchtxt + "%' or Subjectcode like '%" + ws.Seacrchtxt + "%')"
	}
	if ws.Curriculumnature != "" {
		where = where + " and Curriculumnature='" + ws.Curriculumnature + "'"
	}
	if ws.Curriculumstype != "" {
		where = where + " and Curriculumstype='" + ws.Curriculumstype + "'"
	}
	if ws.Curriculumsid > 0 {
		where = where + " and Curriculumsid=" + strconv.Itoa(ws.Curriculumsid) + ""
	}
	if ws.Chaptername != "" {
		where = where + " and(Chaptername='" + ws.Chaptername + "' or Chaptername like '%" + ws.Chaptername + "%')"
	}
	if ws.TeacherId > 0 {
		where = where + " and UsersId=" + strconv.Itoa(ws.TeacherId) + ""
	}
	if ws.Classesid > 0 {
		where = where + " and Classesid=" + strconv.Itoa(ws.Classesid) + ""
	}
	if ws.Begindate != "" {
		where = where + " and Begindate<='" + ws.Begindate + "'"
	}
	if ws.Enddate != "" {
		where = where + " and Enddate>='" + ws.Enddate + "'"
	}
	return where
}
func GetWhereString2(ws viewmodel.QueryCurriculumWhere) (where string) {

	if ws.Subjectcode != "" {
		where = where + " and sc.Subjectcode='" + ws.Subjectcode + "'"
	}
	if ws.Seacrchtxt != "" {
		where = where + " and(cc.Curriculumname like '%" + ws.Seacrchtxt + "%' or sc.Subjectcode like '%" + ws.Seacrchtxt + "%' or sc.Subjectname like '%" + ws.Seacrchtxt + "%' or cs.Classesname like '%" + ws.Seacrchtxt + "%')"
	}
	if ws.Curriculumsid > 0 {
		where = where + " and ccc.Curriculumsid=" + strconv.Itoa(ws.Curriculumsid) + ""
	}
	if ws.TeacherId > 0 {
		where = where + " and ccc.UsersId=" + strconv.Itoa(ws.TeacherId) + ""
	}
	if ws.Classesid > 0 {
		where = where + " and ccc.Classesid=" + strconv.Itoa(ws.Classesid) + ""
	}
	return where
}
func GetWhereString3(ws viewmodel.QueryCurriculumWhere) (where string) {

	if ws.Seacrchtxt != "" {
		where = where + " and(ct.Chaptername like '%" + ws.Seacrchtxt + "%')"
	}
	if ws.CurriculumsclasscentreId > 0 {
		where = where + " and cccc.CurriculumsclasscentreId=" + strconv.Itoa(ws.CurriculumsclasscentreId) + ""
	}
	if ws.TeacherId > 0 {
		where = where + " and cccc.UsersId=" + strconv.Itoa(ws.TeacherId) + ""
	}
	if ws.Begindate != "" {
		where = where + " and cccc.Begindate<='" + ws.Begindate + "'"
	}
	if ws.Enddate != "" {
		where = where + " and cccc.Enddate>='" + ws.Enddate + "'"
	}
	return where
}

/*
获取某教师下某班级下所有已上课程的到课率
*/
func QueryCurriculumChaptersInfo(lg viewmodel.GetAverageclassrate, bt core.BasicsToken, dbmap *gorp.DbMap) (list []viewmodel.CurriculumChapters) {
	//	sql := "select ccc.Id as Curriculumsclasscentreid,ccc.Curriculumsid,ccc.Newchapter,ccc.Averageclassrate,ccs.Curriculumname from curriculumsclasscentre as ccc inner join curriculums as ccs on ccc.Curriculumsid=ccs.Id"
	//	//	sql = sql + " where ccc.Classesid=? and ccc.Usersid=?;"
	//	sql = sql + " where ccc.Usersid=?;"
	//	fmt.Println(sql)
	//	//	_, numerr := dbmap.Select(&list, sql, lg.Classesid, lg.Teacherid)
	//	_, numerr := dbmap.Select(&list, sql, lg.Teacherid)
	//	core.CheckErr(numerr, "curriculumDataAccess|QueryCurriculumChaptersInfo|获取某教师下某班级下所有已上课程的到课率")
	//	sql = "select cccc.Id as Curriculumclassroomchaptercentreid,cccc.Chaptersid,ct.Chaptername,cccc.Plannumber,cccc.Actualnumber,cccc.Toclassrate,cccc.Begindate,cccc.Enddate,tr.State"
	//	sql = sql + " from curriculumclassroomchaptercentre as cccc inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid"
	//	sql = sql + " inner join Chapters as ct on cccc.Chaptersid=ct.Id where cccc.Usersid=? and tr.State=? and cccc.Curriculumsclasscentreid=?;"
	//	if len(list) > 0 {
	//		for i, v := range list {
	//			var itemlist []viewmodel.CurriculumChaptersInfo
	//			fmt.Println(sql)
	//			_, selecterr := dbmap.Select(&itemlist, sql, lg.Teacherid, 2, v.Curriculumsclasscentreid)
	//			core.CheckErr(selecterr, "curriculumDataAccess|QueryCurriculumChaptersInfo|获取某教师下某班级下所有已上课程的到课率:")
	//			list[i].Infos = itemlist
	//		}
	//	}
	//	return list
	sql := "select cccc.Curriculumsclasscentreid,ccm.Curriculumname,ccm.Id as Curriculumsid,count(cccc.Chaptersid) as Newchapter,avg(cccc.Toclassrate) as Averageclassrate from curriculumclassroomchaptercentre as cccc inner join Chapters as cts on cccc.Chaptersid=cts.Id inner join curriculums as ccm on cts.Curriculumsid=ccm.Id inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid"
	sql = sql + " where cccc.Usersid=? and tr.State=2 group by ccm.Id;"
	fmt.Println(sql)
	_, numerr := dbmap.Select(&list, sql, lg.Teacherid)
	core.CheckErr(numerr, "curriculumDataAccess|QueryCurriculumChaptersInfo|获取某教师下某班级下所有已上课程的到课率")
	sql = "select cccc.Id as Curriculumclassroomchaptercentreid,cccc.Chaptersid,cts.Chaptername,cccc.Plannumber,cccc.Actualnumber,cccc.Toclassrate,cccc.Begindate,cccc.Enddate,tr.State,cs.Id as ClassesId,cs.Classesname"
	sql = sql + " from curriculumclassroomchaptercentre as cccc inner join Chapters as cts on cccc.Chaptersid=cts.Id inner join curriculums as ccm on cts.Curriculumsid=ccm.Id inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid"
	sql = sql + " inner join curriculumsclasscentre as ccc on cccc.Curriculumsclasscentreid=ccc.Id inner join classes as cs on cs.Id=ccc.Classesid where cccc.Usersid=? and tr.State=? and ccm.Id=? group by cccc.Id,cs.Id;"
	if len(list) > 0 {
		for i, v := range list {
			var itemlist []viewmodel.CurriculumChaptersInfo
			fmt.Println(sql)
			_, selecterr := dbmap.Select(&itemlist, sql, lg.Teacherid, 2, v.Curriculumsid)
			core.CheckErr(selecterr, "curriculumDataAccess|QueryCurriculumChaptersInfo|获取某教师下某班级下所有已上课程的到课率:")
			list[i].Infos = itemlist
		}
	}
	return list
}

/*
获取某个课程下某个班级每个学生的平均到课率
*/
func GetEveryoneAverageclassrate(lg viewmodel.GetAverageclassrate, bt core.BasicsToken, dbmap *gorp.DbMap) (list []viewmodel.ResponAverageclassrate) {
	sql := "select count(tr.id) as Ct from curriculumsclasscentre as ccc inner join curriculumclassroomchaptercentre as cccc on cccc.Curriculumsclasscentreid=ccc.Id "
	sql = sql + "inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id where tr.state=2 "
	wheresql := ""
	ordersql := " group by ccc.Id,us.Id,us.Truename,us.Userheadimg"
	if lg.Curriculumsid > 0 { //所有课程
		wheresql = wheresql + " and ccc.Curriculumsid=" + strconv.Itoa(lg.Curriculumsid)
	}
	if lg.Classesid > 0 { //所有班级
		wheresql = wheresql + " and ccc.Classesid=" + strconv.Itoa(lg.Classesid)
	}
	if lg.Begindate != "" { //开始时间
		wheresql = wheresql + " and cccc.Enddate>='" + lg.Begindate + "'"
	}
	if lg.Enddate != "" { //结束时间
		wheresql = wheresql + " and cccc.Enddate<='" + lg.Enddate + "'"
	}
	if lg.Pattern == 3 { //模式
		ordersql = ordersql + " limit 0,7;"
	}
	switch bt.Rolestype {
	case 1:
		break
	case 2:
		wheresql = wheresql + " and ccc.Usersid=" + strconv.Itoa(bt.Usersid)
		break
	default:
		return nil
		break
	}
	countnum, numerr := dbmap.SelectInt(sql + wheresql + ";")
	core.CheckErr(numerr, "curriculumDataAccess|GetEveryoneAverageclassrate|获取某个课程下某个班级每个学生的平均到课率")
	sql = "select ccc.Id as Cccid,us.Id as Studentsid,us.Truename,us.Userheadimg,count(ps.id)as Absenteeismnum from curriculumsclasscentre as ccc inner join classes as cs on ccc.Classesid=cs.Id"
	sql = sql + " inner join students as sts on cs.Id=sts.Classesid inner join users as us on sts.Id=us.Id"
	sql = sql + " inner join curriculumclassroomchaptercentre as cccc on cccc.Curriculumsclasscentreid=ccc.Id"
	sql = sql + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id"
	sql = sql + " left join pointtos as ps on (ps.Curriculumclassroomchaptercentreid=cccc.Id and ps.Usersid=sts.Id and ps.State=0)where tr.state=2 "
	_, cuclisterr := dbmap.Select(&list, sql+wheresql+ordersql)
	core.CheckErr(cuclisterr, "curriculumDataAccess|GetEveryoneAverageclassrate|获取某个课程下某个班级每个学生的平均到课率")
	if len(list) > 0 {
		for k := 0; k < len(list); k++ {
			list[k].Sumnum = int(countnum)
			if list[k].Absenteeismnum == 0 {
				list[k].Averageclassrate = 1
			} else {
				list[k].Averageclassrate = float32((float32(countnum) - float32(list[k].Absenteeismnum)) / float32(countnum))
			}
		}
	}
	return list
}

/*
学生每个课程的出勤统计
*/
func GetStudentsClassesAvg(lg viewmodel.GetStudentsinfo, dbmap *gorp.DbMap) (list []viewmodel.ResponStudentsClassesAvg) {
	sql := `select Usersid,Curriculumsid,Curriculumname,Curriculumicon,Classesnum,Sumcounts,Absentnum,((Classesnum-Absentnum)/Classesnum)AvgAttendance from(select sts.Id Usersid,ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,
	IFNULL(cccc1.Classesnum,0) Classesnum,IFNULL(cccc.Sumcounts,0) Sumcounts,IFNULL(cccc2.Absentnum,0) Absentnum
	from students sts inner join curriculumsclasscentre ccc on sts.Classesid=ccc.Classesid
	inner join curriculums cc on cc.Id=ccc.Curriculumsid inner join (select cccc.Curriculumsclasscentreid,
	count(*) Sumcounts from curriculumsclasscentre ccc inner join curriculumclassroomchaptercentre cccc on ccc.Id=cccc.Curriculumsclasscentreid
	group by cccc.Curriculumsclasscentreid) cccc on cccc.Curriculumsclasscentreid=ccc.Id
	inner join (select cccc1.Curriculumsclasscentreid,count(*) Classesnum from curriculumclassroomchaptercentre cccc1 inner join teachingrecord tr on (tr.Curriculumclassroomchaptercentreid=cccc1.Id and tr.state=2)
	group by cccc1.Curriculumsclasscentreid) cccc1 on cccc1.Curriculumsclasscentreid=ccc.Id
	left join (select cccc2.Curriculumsclasscentreid,ps.Usersid,count(ps.Id) Absentnum from curriculumclassroomchaptercentre cccc2
	inner join pointtos ps on (ps.Curriculumclassroomchaptercentreid=cccc2.Id and ps.state=0)
	inner join teachingrecord tr on(tr.Curriculumclassroomchaptercentreid=cccc2.Id and tr.state=2)
	group by cccc2.Curriculumsclasscentreid,ps.Usersid) cccc2 on (cccc2.Curriculumsclasscentreid=ccc.Id and cccc2.Usersid=sts.Id)
	where sts.Id=? group by sts.Id,ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cccc1.Classesnum,cccc.Sumcounts,cccc2.Absentnum) tb;`
	_, cuclisterr := dbmap.Select(&list, sql, lg.Studentsid)
	core.CheckErr(cuclisterr, "curriculumDataAccess|GetStudentsClassesAvg|学生每个课程的出勤统计")
	return list
}

/*
查询课程分类
*/
func GetSubjectclassList(lg viewmodel.GetStudentsinfo, dbmap *gorp.DbMap) (list []viewmodel.ViewSubjectclass) {
	sql := "select sc.Subjectcode,sc.Subjectname,count(cc.Id) as Countnum from subjectclass as sc left join curriculums as cc on (sc.Subjectcode=cc.Subjectcode or sc.Subjectcode=left(cc.Subjectcode,LENGTH(sc.Subjectcode))) where sc.Superiorsubjectcode=? group by sc.Subjectcode,sc.Subjectname;"
	_, cuclisterr := dbmap.Select(&list, sql, lg.SubjectclassCode)
	core.CheckErr(cuclisterr, "curriculumDataAccess|GetSubjectclassList|查询课程分类")
	return list
}

/*
查询课程分类
*/
func GetSubjectclassListPG(lg viewmodel.QueryCurriculumWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from subjectclass where 1=1"
	wheresql := GetWhereString(lg)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "curriculumDataAccess|GetSubjectclassListPG|系统后台获取课程分类:")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.ViewSubjectclass
		sql := "select * from subjectclass sc where 1=1 "
		sql = sql + wheresql + " group by sc.Subjectcode,sc.Subjectname " + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "curriculumDataAccess|GetSubjectclassListPG|系统后台获取课程分类:")
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
获取筛选条件课程数据1
*/
func GetFilterDataCurriculum(lg core.BasicsToken, wh viewmodel.QueryFilterWhere, dbmap *gorp.DbMap) (list []viewmodel.RespFilterCurriculum) {

	sql1 := "select ccs.Id as Curriculumsid,ccs.Curriculumname,ccs.Curriculumicon,ccs.Curriculumnature,ccs.Curriculumstype from curriculums as ccs inner join curriculumsclasscentre as ccc on ccs.Id=ccc.Curriculumsid" //获取所有的课程
	wheresql := " where 1=1"
	switch lg.Rolestype {
	case 1:
		wheresql = wheresql + ""
		break
	case 2:
		wheresql = wheresql + " and ccc.Usersid=" + strconv.Itoa(lg.Usersid)
		break
	default:
		return nil
		break
	}
	if wh.TeacherIds != "" && wh.TeacherIds != "0" {
		wheresql = " where 1=1 and ccc.Usersid in(" + wh.TeacherIds + ")"
	}
	_, selecterr := dbmap.Select(&list, sql1+wheresql+" group by ccs.Id,ccs.Curriculumname,ccs.Curriculumicon,ccs.Curriculumnature,ccs.Curriculumstype;")
	//	fmt.Println(sql1 + wheresql + " group by ccs.Id,ccs.Curriculumname,ccs.Curriculumicon,ccs.Curriculumnature,ccs.Curriculumstype;")
	core.CheckErr(selecterr, "curriculumDataAccess|GetFilterDataCurriculum|获取筛选条件课程数据1")
	return list
}

/*
获取筛选条件班级数据2
*/
func GetFilterDataClasses(lg core.BasicsToken, wh viewmodel.QueryFilterWhere, dbmap *gorp.DbMap) (list []viewmodel.RespFilterClass) {
	wheresql := " where 1=1"
	switch lg.Rolestype {
	case 1:
		wheresql = wheresql + ""
		break
	case 2:
		wheresql = wheresql + " and ccc.Usersid=" + strconv.Itoa(lg.Usersid)
		break
	default:
		return nil
		break
	}
	if wh.TeacherIds != "" && wh.TeacherIds != "0" {
		wheresql = " where 1=1 and ccc.Usersid in(" + wh.TeacherIds + ")"
	}
	sql1 := "select cs.Id as Classesid,cs.Classesname,cs.Classesnum,cs.Classesicon,cs.Classstate from curriculumsclasscentre as ccc inner join classes as cs on ccc.Classesid=cs.Id"
	_, selecterr := dbmap.Select(&list, sql1+wheresql+" group by cs.Id,cs.Classesname,cs.Classesnum,cs.Classesicon,cs.Classstate;")
	core.CheckErr(selecterr, "curriculumDataAccess|GetFilterDataClasses|获取筛选条件班级数据2")
	return list
}

/*
获取班级的到课平均率
*/
func GetClassAveragerate(lg viewmodel.GetAverageclassrate, bt core.BasicsToken, dbmap *gorp.DbMap) (list []viewmodel.ResponAverageclassrate) {
	sqlpattern := "select "
	sql := " cs.Id as Classesid,cs.Classesname,cs.Classesnum,cs.Classesicon,cs.Classstate,avg(cccc.Toclassrate) as Averageclassrate,sum(cccc.Plannumber) as Plannumber,sum(cccc.Actualnumber)as Actualnumber"
	sql = sql + " from curriculumsclasscentre as ccc inner join classes as cs on ccc.Classesid=cs.Id inner join curriculumclassroomchaptercentre as cccc on cccc.Curriculumsclasscentreid=ccc.Id inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id"
	wheresql := " where tr.state=2"
	ordersql := " group by Ddate,cs.Id,cs.Classesname,cs.Classesnum,cs.Classesicon,cs.Classstate"
	if lg.Curriculumsid > 0 { //所有课程
		wheresql = wheresql + " and ccc.Curriculumsid=" + strconv.Itoa(lg.Curriculumsid)
	}
	if lg.Majorid > 0 { //按照专业查询
		wheresql = wheresql + " and cs.Majorid=" + strconv.Itoa(lg.Majorid)
	}
	if lg.Classesid > 0 { //按照班级Id
		wheresql = wheresql + " and ccc.Classesid=" + strconv.Itoa(lg.Classesid)
	}
	if lg.Begindate != "" { //开始时间
		wheresql = wheresql + " and cccc.Enddate>='" + lg.Begindate + "'"
	}
	if lg.Enddate != "" { //结束时间
		wheresql = wheresql + " and cccc.Enddate<='" + lg.Enddate + "'"
	}
	if lg.Pattern == 1 { //本学期
		sqlpattern = sqlpattern + "DATE_FORMAT(cccc.Enddate,'%Y%m%d') as Ddate,"
		ordersql = ordersql + ";"
	}
	if lg.Pattern == 2 { //近一个月
		sqlpattern = sqlpattern + "DATE_FORMAT(cccc.Enddate,'%Y%m%d') as Ddate,"
		ordersql = ordersql + ";"
	}
	if lg.Pattern == 3 { //近7天
		sqlpattern = sqlpattern + "DATE_FORMAT(cccc.Enddate,'%Y%m%d') as Ddate,"
		ordersql = ordersql + " limit 0,7;"
	}
	switch bt.Rolestype {
	case 1:
		break
	case 2:
		wheresql = wheresql + " and ccc.Usersid=" + strconv.Itoa(bt.Usersid)
		break
	default:
		return nil
		break
	}
	fmt.Println(sqlpattern + sql + wheresql + ordersql)
	_, selecterr := dbmap.Select(&list, sqlpattern+sql+wheresql+ordersql)
	core.CheckErr(selecterr, "curriculumDataAccess|GetClassAveragerate|获取班级的到课平均率")
	return list
}

//添加学科数据
func AddSubjectclass(sjc *curriculum.Subjectclass, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Subjectclass{}, "subjectclass").SetKeys(true, "Id")
	inerr = dbmap.Insert(sjc)
	core.CheckErr(inerr, "curriculumDataAccess|AddSubjectclass|添加学科数据")
	return inerr
}

//查询学科代码是否唯一
func QueryUniqueSubjectcode(Subjectcode, Subjectname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Subjectcode) from  subjectclass WHERE (Subjectcode=? or Subjectname=?);", Subjectcode, Subjectname)
	return num
}

// 修改学科数据
func UpdateSubjectclass(cps *curriculum.Subjectclass, dbmap *gorp.DbMap) (err error) {
	//	dbmap.AddTableWithName(curriculum.Subjectclass{}, "subjectclass").SetKeys(true, "Id")
	//	_, inerr = dbmap.Update(cps)

	sql := `UPDATE subjectclass SET Subjectname=? WHERE (Subjectcode=?);`
	if _, err = dbmap.Exec(sql, cps.Subjectname, cps.Subjectcode); err != nil {
		core.CheckErr(err, "curriculumDataAccess|UpdateSubjectclass|修改学科数据:")
	}

	return err
}

//删除学科数据
func DeleteSubjectclass(cps *curriculum.Subjectclass, dbmap *gorp.DbMap) (rd core.Returndata) {
	//	dbmap.AddTableWithName(curriculum.Subjectclass{}, "subjectclass").SetKeys(true, "Id")
	//	_, inerr = dbmap.Delete(cps)
	countint, inerr := dbmap.SelectInt("select count(*) from subjectclass where Superiorsubjectcode='" + cps.Subjectcode + "'")
	core.CheckErr(inerr, "curriculumDataAccess|DeleteSubjectclass|删除学科数据:")
	if countint == 0 {
		_, inerr = dbmap.Exec("delete from subjectclass where Subjectcode='" + cps.Subjectcode + "'")
		core.CheckErr(inerr, "curriculumDataAccess|DeleteSubjectclass|删除学科数据:")
		if inerr == nil {
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "执行失败"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "有下级学科数据，不能删除"
	}
	return rd
}

/*
查询课程
*/
func GetCurriculumsListPG(lg viewmodel.QueryCurriculumWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from curriculums as cc inner join subjectclass as sc on cc.Subjectcode=sc.Subjectcode where 1=1"
	wheresql := GetWhereString(lg)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumsListPG|系统后台获取课程:")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.ViewCurriculums //[]curriculum.Curriculums []viewmodel.ViewSubjectclass
		sql := "select cc.Id as CurriculumsId,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumnature,cc.Curriculumstype,cc.Curriculumsdetails,cc.Chaptercount,cc.Averageclassrate,cc.Subjectcode,sc.Subjectname from curriculums as cc inner join subjectclass as sc on cc.Subjectcode=sc.Subjectcode where 1=1 "
		sql = sql + wheresql + core.GetLimitString(pg) + ";"
		fmt.Println(sql)
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumsListPG|系统后台获取课程:")
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
查询课程
*/
func GetCurriculumsInfo(lg viewmodel.QueryCurriculumWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	var info viewmodel.ViewCurriculums
	sql := "select cc.Id as CurriculumsId,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumnature,cc.Curriculumstype,cc.Curriculumsdetails,cc.Chaptercount,cc.Averageclassrate,cc.Subjectcode,sc.Subjectname "
	sql = sql + "from curriculums as cc inner join subjectclass as sc on cc.Subjectcode=sc.Subjectcode where 1=1 and cc.Id=" + strconv.Itoa(lg.Curriculumsid) + ";"
	fmt.Println(sql)
	sqlerr := dbmap.SelectOne(&info, sql)
	core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumsInfo|系统后台获取课程:")
	if sqlerr == nil {
		if info.CurriculumsId > 0 {
			rd.Rcode = "1000"
			pg.PageCount = 1
			pg.PageData = info
			rd.Result = pg
		} else {
			rd.Rcode = "1099"
			rd.Reason = "未找到数据"
		}
	} else {
		rd.Rcode = "1001"
		rd.Reason = "数据查询失败"
	}
	return rd
}

//添加课程数据
func AddCurriculums(cc *curriculum.Curriculums, dbmap *gorp.DbMap) (inerr error) {
	ifnull := cc
	ifnull.Id = 0
	cc.Id = 0
	dbmap.AddTableWithName(curriculum.Curriculums{}, "curriculums").SetKeys(true, "Id")
	fmt.Println("添加课程数据:", ifnull)
	inerr = dbmap.SelectOne(cc, "select * from curriculums where Curriculumname=? and Subjectcode=? and Curriculumstype=?", ifnull.Curriculumname, ifnull.Subjectcode, ifnull.Curriculumstype)
	if inerr == nil {
		if cc.Id == 0 {
			inerr = dbmap.Insert(cc)
		}
	} else {
		if cc.Id == 0 {
			inerr = dbmap.Insert(cc)
		}
	}
	core.CheckErr(inerr, "curriculumDataAccess|AddCurriculums|添加课程数据")
	return inerr
}

//修改课程数据
func UpdateCurriculums(cps *curriculum.Curriculums, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Curriculums{}, "curriculums").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "curriculumDataAccess|UpdateCurriculums|修改课程数据:")
	return inerr
}

//删除课程数据
func DeleteCurriculums(cps *curriculum.Curriculums, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Curriculums{}, "curriculums").SetKeys(true, "Id")
	_, inerr = dbmap.Delete(cps)
	core.CheckErr(inerr, "curriculumDataAccess|DeleteCurriculums|删除课程数据:")
	return inerr
}

/*
查询课程章节
*/
func GetChaptersListPG(lg viewmodel.QueryCurriculumWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from chapters where 1=1"
	wheresql := GetWhereString(lg)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "curriculumDataAccess|GetChaptersListPG|系统后台获取课程章节:")
	if sqlerr == nil && countint > 0 {
		var list []curriculum.Chapters
		sql := "select * from chapters where 1=1 "
		sql = sql + wheresql + " order by ChaptersIndex " + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "curriculumDataAccess|GetChaptersListPG|系统后台获取课程章节:")
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

///*
//查询课程章节,执行sql语句
//*/
//func QueryChaptersList(sql string, dbmap *gorp.DbMap) (list []curriculum.Curriculums) {
//	_, sqlerr := dbmap.Select(&list, sql)
//	core.CheckErr(sqlerr, "curriculumDataAccess|GetChaptersListPG|系统后台获取课程章节:")
//	return list
//}

//添加课程章节数据
func AddChapters(ct *curriculum.Chapters, dbmap *gorp.DbMap) (inerr error) {
	ifnull := ct
	ifnull.Id = 0
	ct.Id = 0
	dbmap.AddTableWithName(curriculum.Chapters{}, "chapters").SetKeys(true, "Id")
	inerr = dbmap.SelectOne(ct, "select * from chapters where Curriculumsid=? and Chaptername=?", ifnull.Curriculumsid, ifnull.Chaptername)
	if ct.Id == 0 {
		inerr = dbmap.Insert(ct)
	}
	go UpdateChaptersNum(ct.Curriculumsid)
	core.CheckErr(inerr, "curriculumDataAccess|AddChapters|添加课程章节数据")
	return inerr
}

//更改主课程的章节数
func UpdateChaptersNum(Curriculumsid int) {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	sqlcount := "select count(*) from chapters where Curriculumsid=?"
	countint, counterr := dbmap.SelectInt(sqlcount, Curriculumsid)
	if countint == 0 {
		countint = 1
	}
	if counterr == nil {
		dbmap.Exec("update curriculums set Chaptercount=? where Id=?", countint, Curriculumsid) //更新课程下的章节数
	}
}

//修改课程章节数据
func UpdateChapters(cps *curriculum.Chapters, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Chapters{}, "chapters").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	go UpdateChaptersNum(cps.Curriculumsid)
	core.CheckErr(inerr, "curriculumDataAccess|UpdateChapters|修改课程章节数据:")
	return inerr
}

//删除课程章节数据
func DeleteChapters(cps *curriculum.Chapters, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Chapters{}, "chapters").SetKeys(true, "Id")
	_, inerr = dbmap.Delete(cps)
	go UpdateChaptersNum(cps.Curriculumsid)
	core.CheckErr(inerr, "curriculumDataAccess|DeleteChapters|删除课程章节数据:")
	return inerr
}

/*
查询课程班级中间配置数据
*/
func GetCurriculumsClassCentreListPG(lg viewmodel.QueryCurriculumWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from curriculumsclasscentre as ccc inner join curriculums as cc on ccc.Curriculumsid=cc.Id inner join users as us on ccc.Usersid=us.Id "
	countsql = countsql + " inner join subjectclass as sc on cc.Subjectcode=sc.Subjectcode inner join classes as cs on ccc.Classesid=cs.Id where 1=1"
	wheresql := GetWhereString2(lg)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	fmt.Println(countsql + wheresql)
	core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumsClassCentreListPG|系统后台获取课程章节总数:")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.GetCurriculumsClassCentreList
		sql := "select ccc.Id as curriculumsclasscentreId,ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumnature,cc.Curriculumstype,us.Truename,ccc.Usersid as TeacherId,cc.Chaptercount,cc.Subjectcode,sc.Subjectname,ccc.Classesid,ccc.Isondemand,ccc.Islive,cs.Classesname"
		sql = sql + " from curriculumsclasscentre as ccc inner join curriculums as cc on ccc.Curriculumsid=cc.Id inner join users as us on ccc.Usersid=us.Id "
		sql = sql + " inner join subjectclass as sc on cc.Subjectcode=sc.Subjectcode inner join classes as cs on ccc.Classesid=cs.Id where 1=1 "
		sql = sql + wheresql + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumsClassCentreListPG|系统后台获取课程章节列表:")
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

//添加课程班级关联表数据
func AddCurriculumsClassCentre(ct *curriculum.Curriculumsclasscentre, dbmap *gorp.DbMap) (inerr error) {
	ifnull := ct
	ifnull.Id = 0
	ct.Id = 0
	dbmap.AddTableWithName(curriculum.Curriculumsclasscentre{}, "curriculumsclasscentre").SetKeys(true, "Id")
	inerr = dbmap.SelectOne(ct, "select * from curriculumsclasscentre where Curriculumsid=? and Classesid=? and Usersid=?", ifnull.Curriculumsid, ifnull.Classesid, ifnull.Usersid)
	if inerr == nil {
		if ct.Id == 0 {
			inerr = dbmap.Insert(ct)
		}
	} else {
		if ct.Id == 0 {
			inerr = dbmap.Insert(ct)
		}
	}
	core.CheckErr(inerr, "curriculumDataAccess|AddCurriculumsClassCentre|添加课程班级关联表数据")
	return inerr
}

//添加课程班级关联表数据
func UpdateCurriculumsClassCentre(ct *curriculum.Curriculumsclasscentre, dbmap *gorp.DbMap) (rd core.Returndata) {
	execsql := "update curriculumsclasscentre set Usersid=" + strconv.Itoa(ct.Usersid) + ",Islive=" + strconv.Itoa(ct.Islive) + ",Isondemand=" + strconv.Itoa(ct.Isondemand) + " where Id=" + strconv.Itoa(ct.Id) + ";"
	_, sqlerr := dbmap.Exec(execsql)
	core.CheckErr(sqlerr, "课程主计划修改失败:")
	if sqlerr == nil {
		execsql = "update curriculumclassroomchaptercentre ccc,teachingrecord tr set Usersid=" + strconv.Itoa(ct.Usersid) + " where tr.State=0 and tr.Curriculumclassroomchaptercentreid=ccc.Id and ccc.Curriculumsclasscentreid=" + strconv.Itoa(ct.Id) + ";"
		_, sqlerr = dbmap.Exec(execsql)
		if sqlerr == nil {
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1001"
			rd.Reason = "课程子计划修改失败"
		}
	} else {
		rd.Rcode = "1001"
		rd.Reason = "课程主计划修改失败"
	}
	//	dbmap.AddTableWithName(curriculum.Curriculumsclasscentre{}, "curriculumsclasscentre").SetKeys(true, "Id")
	//	_, inerr = dbmap.Update(ct)
	//	core.CheckErr(inerr, "curriculumDataAccess|UpdateCurriculumsClassCentre|修改课程班级关联表数据:")
	//	return inerr
	return rd
}

//添加课程班级关联表数据
func DeleteCurriculumsClassCentre(ct *curriculum.Curriculumsclasscentre, dbmap *gorp.DbMap) (rd core.Returndata) { //(inerr error) {
	//	dbmap.AddTableWithName(curriculum.Curriculumsclasscentre{}, "curriculumsclasscentre").SetKeys(true, "Id")
	//	_, inerr = dbmap.Exec("delete from curriculumclassroomchaptercentre where Curriculumsclasscentreid=" + strconv.Itoa(ct.Id))
	//	core.CheckErr(inerr, "curriculumDataAccess|DeleteCurriculumsClassCentre|删除课程班级关联子表数据:")
	//	_, inerr = dbmap.Delete(ct)
	//	core.CheckErr(inerr, "curriculumDataAccess|DeleteCurriculumsClassCentre|删除课程班级关联表数据:")
	//	return inerr
	_, errexec := dbmap.Exec("call Delkebd(?);", ct.Id)
	if errexec == nil {
		//		commons.ResponseMsg
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1001"
	}
	return rd
}

/*
查询课程班级中间详细配置数据
*/
func GetCurriculumclassroomchaptercentreListPG(lg viewmodel.QueryCurriculumWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from Curriculumclassroomchaptercentre cccc inner join Chapters ct on cccc.Chaptersid=ct.Id inner join curriculumsclasscentre ccc on cccc.curriculumsclasscentreid=ccc.Id inner join classes as cs on ccc.Classesid=cs.Id"
	countsql = countsql + " inner join teachingrecord tr on tr.CurriculumclassroomchaptercentreId=cccc.Id left join classrooms cr on tr.Classroomid=cr.Id "
	countsql = countsql + " left join floors fs on cr.Floorsid=fs.id left join building bd on fs.Buildingid=bd.Id left join Campus cps on bd.Campusid=cps.Id inner join Users us on cccc.Usersid=us.Id where 1=1"
	wheresql := GetWhereString3(lg)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumclassroomchaptercentreListPG|查询课程班级中间详细配置数据1:")
	fmt.Println(countsql + wheresql)
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.GetCurriculumclassroomchaptercentreList //curriculum.Curriculumclassroomchaptercentre
		sql := "select cccc.Id CurriculumclassroomchaptercentreId,ct.Chaptername,cs.Classesname,ccc.Classesid,cccc.Begindate,cccc.Enddate,cccc.Islive,cccc.Isondomian,ifnull(cr.Classroomsname,'')Classroomsname,ifnull(bd.Buildingname,'')Buildingname,ifnull(cps.Campusname,'')Campusname,ifnull(tr.State,0)State,tr.Classroomid,cccc.Usersid TeacherId,cccc.Chaptersid,us.Truename"
		sql = sql + " from Curriculumclassroomchaptercentre cccc inner join Chapters ct on cccc.Chaptersid=ct.Id inner join curriculumsclasscentre ccc on cccc.curriculumsclasscentreid=ccc.Id inner join classes as cs on ccc.Classesid=cs.Id"
		sql = sql + " inner join teachingrecord tr on tr.CurriculumclassroomchaptercentreId=cccc.Id left join classrooms cr on tr.Classroomid=cr.Id "
		sql = sql + " left join floors fs on cr.Floorsid=fs.id left join building bd on fs.Buildingid=bd.Id left join Campus cps on bd.Campusid=cps.Id inner join Users us on cccc.Usersid=us.Id where 1=1 "
		sql = sql + wheresql + "\n ORDER BY ct.ChaptersIndex, ct.Id \n" + core.GetLimitString(pg) + ""
		// log.Println("<<<<<<<<:\n", sql)
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "curriculumDataAccess|GetCurriculumclassroomchaptercentreListPG|查询课程班级中间详细配置数据2:")
		// fmt.Println(sql)
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
检查排课数据的有效性
*/
func QueryCheckUpData(lg viewmodel.GetCurriculumclassroomchaptercentreList, dbmap *gorp.DbMap) (list []viewmodel.GetCurriculumclassroomchaptercentreList) {
	sql := fmt.Sprintf(`select us.Truename from curriculumclassroomchaptercentre cccc inner join teachingrecord tr 
						on cccc.Id=tr.Curriculumclassroomchaptercentreid inner join users us on cccc.Usersid=us.Id where cccc.UsersId!= %d  
						and tr.Classroomid= %d  and cccc.Begindate>=' %s ' and cccc.Enddate<=' %s ';`, lg.TeacherId, lg.Classroomid, lg.Begindate, lg.Enddate)
	dbmap.Select(&list, sql)
	return list
}

//添加课程班级章节中间详细配置数据
func AddCurriculumclassroomchaptercentre(ct *curriculum.Curriculumclassroomchaptercentre, Classroomid int, dbmap *gorp.DbMap) (inerr error) {
	ifnull := ct
	ifnull.Id = 0
	ct.Id = 0
	dbmap.AddTableWithName(curriculum.Curriculumclassroomchaptercentre{}, "curriculumclassroomchaptercentre").SetKeys(true, "Id")
	inerr = dbmap.SelectOne(ct, "select * from curriculumclassroomchaptercentre where Curriculumsclasscentreid=? and Chaptersid=? and Usersid=?", ifnull.Curriculumsclasscentreid, ifnull.Chaptersid, ifnull.Usersid)
	if inerr == nil {
		if ct.Id == 0 {
			inerr = dbmap.Insert(ct)
		}
	} else {
		if ct.Id == 0 {
			inerr = dbmap.Insert(ct)
		}
	}
	//视频数据处理
	lv := live.Lives{}
	dbmap.AddTableWithName(live.Lives{}, "lives").SetKeys(true, "Id")
	if ct.Islive > 0 || ct.Isondomian > 0 { //查询是否有计划如果有则修改时间和老师，如果没有则添加相关数据
		lv.Curriculumclassroomchaptercentreid = ct.Id
		lv.Begindate = ct.Begindate
		lv.Whenlong = 45
		dbmap.Insert(&lv)
	}
	//点到数据处理 begin
	dbmap.AddTableWithName(actiondata.Pointtos{}, "pointtos").SetKeys(true, "Id")
	var ptarr []actiondata.Pointtos //获取需要添加进去的学生数据
	_, inerr = dbmap.Select(&ptarr, "select st.Id Usersid,cccc.Id Curriculumclassroomchaptercentreid,0 State,0 Ismodify from Curriculumclassroomchaptercentre cccc inner join curriculumsclasscentre ccc on cccc.Curriculumsclasscentreid=ccc.Id inner join students st on ccc.Classesid=st.Classesid where cccc.Id=?;", ct.Id)
	//迭代循环添加
	for _, v := range ptarr {
		dbmap.Insert(&v)
	}
	//点到数据处理 end
	//上课教室 begin
	tr := actiondata.Teachingrecord{Classroomid: Classroomid, Curriculumclassroomchaptercentreid: ct.Id, State: 0}
	dbmap.AddTableWithName(actiondata.Teachingrecord{}, "teachingrecord").SetKeys(true, "Id")
	dbmap.Insert(&tr)
	//上课教室 end

	core.CheckErr(inerr, "curriculumDataAccess|AddCurriculumclassroomchaptercentre|添加课程班级章节关联表数据")
	return inerr
}

//修改课程班级章节中间详细配置数据
func UpdateCurriculumclassroomchaptercentre(ct *curriculum.Curriculumclassroomchaptercentre, Classroomid int, dbmap *gorp.DbMap) (rd core.Returndata) {
	//	dbmap.AddTableWithName(curriculum.Curriculumclassroomchaptercentre{}, "Curriculumclassroomchaptercentre").SetKeys(true, "Id")
	//	_, inerr = dbmap.Update(ct)
	//	core.CheckErr(inerr, "curriculumDataAccess|UpdateCurriculumclassroomchaptercentre|修改课程班级章节中间详细配置数据:")
	//	return inerr
	_, sqlerr := dbmap.Exec("update Curriculumclassroomchaptercentre set Usersid=?,Islive=?,Isondomian=?,Begindate=?,Enddate=? where Id=?;", ct.Usersid, ct.Islive, ct.Isondomian, ct.Begindate, ct.Enddate, ct.Id)
	if sqlerr == nil {
		lv := live.Lives{}
		dbmap.AddTableWithName(live.Lives{}, "lives").SetKeys(true, "Id")
		dbmap.SelectOne(&lv, "select * from lives where Curriculumclassroomchaptercentreid=? limit 0,1;", ct.Id)
		if ct.Islive > 0 || ct.Isondomian > 0 { //查询是否有计划如果有则修改时间和老师，如果没有则添加相关数据
			if lv.Id > 0 {
				lv.Begindate = ct.Begindate
				dbmap.Update(&lv)
			} else {
				lv.Curriculumclassroomchaptercentreid = ct.Id
				lv.Begindate = ct.Begindate
				lv.Ischeckcomment = 0
				lv.IsRelease = 0
				lv.Whenlong = 45
				dbmap.Insert(&lv)
			}
		} else if ct.Islive <= 0 && ct.Isondomian <= 0 { //直接干掉数据
			dbmap.Delete(&lv)
		}
		if Classroomid > 0 { //如果大于0代表修改了上课的教室
			dbmap.Exec("update teachingrecord set Classroomid=? where Curriculumclassroomchaptercentreid=?", Classroomid, ct.Id)
		}
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1001"
		rd.Reason = "课程子计划修改失败"
	}
	return rd
}

//删除课程班级章节中间详细配置数据
func DeleteCurriculumclassroomchaptercentre(ct *curriculum.Curriculumclassroomchaptercentre, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Curriculumclassroomchaptercentre{}, "curriculumclassroomchaptercentre").SetKeys(true, "Id")
	_, inerr = dbmap.Exec("delete from curriculumclassroomchaptercentre where Curriculumsclasscentreid=" + strconv.Itoa(ct.Id))
	core.CheckErr(inerr, "curriculumDataAccess|DeleteCurriculumclassroomchaptercentre|删除课程班级章节中间详细配置数据:")
	_, inerr = dbmap.Delete(ct)
	core.CheckErr(inerr, "curriculumDataAccess|DeleteCurriculumclassroomchaptercentre|删除课程班级章节中间详细配置数据:")
	return inerr
}

func UpdateLivePath(Usersid int, Classroomid int, path string, virtualpath string, typeint int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var uvfpc viewmodel.UpVideoFileCollect
	timestr := time.Now().Format("2006-01-02 15:04:05")
	sql := "select cccc.Islive,cccc.Isondomian,tr.Curriculumclassroomchaptercentreid from curriculumclassroomchaptercentre as cccc inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid"
	sql = sql + " where tr.Classroomid=? and to_days(cccc.Enddate)=to_days(now())"
	if Usersid > 0 {
		sql = sql + " and cccc.UsersId=" + strconv.Itoa(Usersid)
	}
	if typeint == 1 {
		sql = sql + " and cccc.Begindate<='" + timestr + "'"
		sql = sql + " order by cccc.Begindate desc limit 0,1;"
	} else {
		sql = sql + " and cccc.Enddate<='" + timestr + "'"
		sql = sql + " order by cccc.Enddate desc limit 0,1;"
	}
	fmt.Println(sql, Classroomid)
	dbmap.SelectOne(&uvfpc, sql, Classroomid)
	Enclosurenamestr := "电脑桌面视频文件"
	if typeint == 1 {
		Enclosurenamestr = "教室录像视频文件"
	}
	//判断获取当前教室是否有录播的课程
	if uvfpc.Islive > 0 || uvfpc.Isondomian > 0 { //如果有则更新播放文件
		var lv live.Lives
		//如果已经有文件了，则将文件资源存入到课程资源中，并显示未公开的状态
		sqllv := "select * from Lives where Curriculumclassroomchaptercentreid=?;"
		fmt.Println(sqllv)
		dbmap.SelectOne(&lv, sqllv, uvfpc.Curriculumclassroomchaptercentreid)
		if lv.Id == 0 { //如果没有找到数据，则默认添加一条新的数据
			rd.Rcode = "1006"
			rd.Result = "未找到本次章节课程的信息"
		} else {
			if lv.Livepath1 == "" || lv.Livepath2 == "" {
				dbmap.AddTableWithName(live.Lives{}, "lives").SetKeys(true, "Id")
				fmt.Println("保存视频")
				if typeint == 1 {
					fmt.Println("保存线路1视频")
					lv.Livepath2 = virtualpath
				} else {
					fmt.Println("保存线路2视频")
					lv.Livepath1 = virtualpath
				}
				_, uperr := dbmap.Update(&lv)
				fmt.Println(uperr)
				core.CheckErr(uperr, "curriculumDataAccess|UpdateLivePath|保存视频录制的文件路径")
				rd.Rcode = "1000"
			} else {
				fmt.Println("存入资源1")
				ecrs := curriculum.Enclosure{Curriculumclassroomchaptercentreid: uvfpc.Curriculumclassroomchaptercentreid,
					Enclosurename: Enclosurenamestr, Enclosuretype: "视频文件", Enclosuresize: 0, EnclosureVirtualPath: lv.Livepath1, Enclosurepath: path,
					Createdate: time.Now().Format("2006-01-02 15:05:03"), IsPublish: 1}
				AddEnclosure(&ecrs, dbmap)
				rd.Rcode = "1000"
			}
		}
	} else { //如果没有 则将文件更新到课程资源中
		fmt.Println("存入资源2")
		ecrs := curriculum.Enclosure{Curriculumclassroomchaptercentreid: uvfpc.Curriculumclassroomchaptercentreid,
			Enclosurename: Enclosurenamestr, Enclosuretype: "视频文件", Enclosuresize: 0, EnclosureVirtualPath: virtualpath, Enclosurepath: path,
			Createdate: time.Now().Format("2006-01-02 15:05:03"), IsPublish: 1}
		AddEnclosure(&ecrs, dbmap)
		rd.Rcode = "1000"
	}
	return rd
}

//添加课程班级章节附件关联表数据
func AddEnclosure(ct *curriculum.Enclosure, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(curriculum.Enclosure{}, "enclosure").SetKeys(true, "Id")
	inerr = dbmap.Insert(ct)
	core.CheckErr(inerr, "curriculumDataAccess|AddEnclosure|添加课程班级章节附件关联表数据")
	return inerr
}

//根据教室ID查询此教室是否在上课，在上什么课,可能存在多个班级在上课
func QueryCurriculumByClassRoom(classroomid int, datestr string, dbmap *gorp.DbMap) (list []viewmodel.QueryClassroomCurriculumInfo) {
	sql1 := "select cc.Curriculumname,ct.Chaptername,us.Nickname,ccc.Classesid,cs.Classesname,tr.Curriculumclassroomchaptercentreid from teachingrecord as tr inner join curriculumclassroomchaptercentre as cccc on tr.Curriculumclassroomchaptercentreid=cccc.Id "
	sql1 = sql1 + " inner join curriculumsclasscentre as ccc on cccc.Curriculumsclasscentreid=ccc.Id inner join Classes as cs on cs.Id=ccc.Classesid inner join Chapters as ct on ct.Id=cccc.Chaptersid inner join curriculums as cc on ct.Curriculumsid=cc.Id"
	sql1 = sql1 + " inner join users as us on us.Id=cccc.Usersid where tr.Classroomid=? and cccc.Begindate<=? and cccc.Enddate>=?;"
	fmt.Println(sql1, classroomid, datestr, datestr)
	_, serr := dbmap.Select(&list, sql1, classroomid, datestr, datestr)
	core.CheckErr(serr, "curriculumDataAccess|QueryCurriculumByClassRoom|根据教室ID查询此教室是否在上课，在上什么课")
	return list

}

/*
管理者查看各种出勤统计分析
*/
func QueryAttendanceAnalysisList(lg viewmodel.AttendanceAnalysisWhere, bt core.BasicsToken, dbmap *gorp.DbMap) (list []viewmodel.QueryAttendanceAnalysis) {
	//sql := "select avg(cccc.Toclassrate) as Toclassrate,ces.Enrollmentyear,mj.Id as Majorid,mj.Majorname,cg.Id as Collegeid,cg.Collegename,cc.Id as Curriculumsid,cc.Curriculumname,ccc.Classesid"
	sql := "select "
	ordersql := " group by ces.Enrollmentyear,cg.Id,mj.Id,cc.Id,ccc.Classesid;"
	if lg.Analysistype == 0 { //年级出勤统计分析
		ordersql = " group by ces.Enrollmentyear asc;"
		sql = sql + "avg(cccc.Toclassrate) as Analysisvalue,ces.Enrollmentyear as Analysisname"
	} else if lg.Analysistype == 1 { //学院出勤统计分析
		ordersql = " group by cg.Id;"
		sql = sql + "avg(cccc.Toclassrate) as Analysisvalue,cg.Collegename as Analysisname"
	} else if lg.Analysistype == 2 { //专业出勤统计分析
		ordersql = " group by mj.Id;"
		sql = sql + "avg(cccc.Toclassrate) as Analysisvalue,mj.Majorname as Analysisname"
	} else if lg.Analysistype == 3 { //班级出勤统计分析
		ordersql = " group by ccc.Classesid;"
		sql = sql + "avg(cccc.Toclassrate) as Analysisvalue,ces.Classesname as Analysisname"
	} else if lg.Analysistype == 4 { //课程出勤统计分析
		ordersql = " group by cc.Id;"
		sql = sql + "avg(cccc.Toclassrate) as Analysisvalue,cc.Curriculumname as Analysisname"
	} else {
		//core.CheckErr(selecterr, "curriculumDataAccess|QueryAttendanceAnalysisList|未传递数据类型")
		return list
	}
	sql = sql + " from curriculumsclasscentre as ccc inner join classes as ces on ccc.Classesid=ces.Id inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid "
	sql = sql + " inner join major as mj on ces.Majorid=mj.Id inner join College as cg on mj.Collegeid=cg.Id inner join teachingrecord as tr on cccc.Id=tr.Curriculumclassroomchaptercentreid"
	sql = sql + " inner join curriculums as cc on cc.Id=ccc.Curriculumsid "
	wheresql := " where tr.state=2"
	if lg.Collegeid > 0 {
		wheresql = wheresql + " and cg.Id=" + strconv.Itoa(lg.Collegeid)
	}
	if lg.Majorid > 0 {
		wheresql = wheresql + " and mj.Id=" + strconv.Itoa(lg.Majorid)
	}
	if lg.Curriculumsid > 0 {
		wheresql = wheresql + " and cc.Id=" + strconv.Itoa(lg.Curriculumsid)
	}
	if lg.Begindate != "" {
		wheresql = wheresql + " and cccc.Begindate>='" + lg.Begindate + "'"
	}
	if lg.Enddate != "" {
		wheresql = wheresql + " and cccc.Enddate<='" + lg.Enddate + "'"
	}
	if lg.Dateint > 0 {
		wheresql = wheresql + " and DATE_SUB(CURDATE(),INTERVAL " + strconv.Itoa(lg.Dateint) + " DAY)<=date(cccc.Begindate)"
	}
	if lg.Gradeint > 0 {
		wheresql = wheresql + " and ces.Enrollmentyear=" + strconv.Itoa(time.Now().Year()-lg.Gradeint)
	}

	fmt.Println(sql + wheresql + ordersql)
	_, selecterr := dbmap.Select(&list, sql+wheresql+ordersql)
	core.CheckErr(selecterr, "curriculumDataAccess|QueryAttendanceAnalysisList|管理者查看各种出勤统计分析")
	return list
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
