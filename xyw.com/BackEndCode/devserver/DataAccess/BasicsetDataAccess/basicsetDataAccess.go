package basicsetDataAccess

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"dev.project/BackEndCode/devserver/model/basicset"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/viewmodel"

	"gopkg.in/gorp.v1"
)

func GetWhereString(ws viewmodel.QueryBasicsetWhere) (where string) {
	if ws.Campuscode != "" {
		where = where + " and Campuscode='" + ws.Campuscode + "'"
	}
	if ws.Campusid > 0 {
		where = where + " and Campusid=" + strconv.Itoa(ws.Campusid)
	}
	if ws.Campusname != "" {
		where = where + " and(Campusname like '%" + ws.Campusname + "%' or Campuscode like '%" + ws.Campusname + "%')"
	}
	if ws.Buildingcode != "" {
		where = where + " and Buildingcode='" + ws.Buildingcode + "'"
	}
	if ws.Buildingname != "" {
		where = where + " and(Buildingname like '%" + ws.Buildingname + "%' or Buildingcode like '%" + ws.Buildingname + "%')"
	}
	if ws.Buildingid > 0 {
		where = where + " and Buildingid=" + strconv.Itoa(ws.Buildingid)
	}
	if ws.Floorscode != "" {
		where = where + " and Floorscode='" + ws.Floorscode + "'"
	}
	if ws.Floorsid > 0 {
		where = where + " and Floorsid=" + strconv.Itoa(ws.Floorsid)
	}
	if ws.Classroomid > 0 {
		where = where + " and Classroomid=" + strconv.Itoa(ws.Classroomid)
	}
	if ws.Classroomscode != "" {
		where = where + " and Classroomscode='" + ws.Classroomscode + "'"
	}
	if ws.Collegecode != "" {
		where = where + " and Collegecode='" + ws.Collegecode + "'"
	}
	if ws.Collegename != "" {
		where = where + " and Collegename like '%" + ws.Collegename + "%'"
	}
	if ws.Collegeid > 0 {
		where = where + " and Collegeid=" + strconv.Itoa(ws.Collegeid)
	}
	if ws.Majorid > 0 {
		where = where + " and Majorid=" + strconv.Itoa(ws.Majorid)
	}
	if ws.Majorcode != "" {
		where = where + " and Majorcode='" + ws.Majorcode + "'"
	}
	if ws.Majorname != "" {
		where = where + " and Majorname like '%" + ws.Majorname + "%'"
	}
	if ws.Classesid > 0 {
		where = where + " and Classesid=" + strconv.Itoa(ws.Classesid)
	}
	if ws.Classescode != "" {
		where = where + " and Classescode='" + ws.Classescode + "'"
	}
	if ws.Classesname != "" {
		where = where + " and Classesname like '%" + ws.Classesname + "%'"
	}

	return where
}

//添加校区
func AddCampus(cps *basicset.Campus, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Campus{}, "campus").SetKeys(true, "Id")
	inerr = dbmap.Insert(cps)
	core.CheckErr(inerr, "basicsetDataAccess|AddCampus|添加校区:")
	return inerr
}

//查询校区代码是否唯一
func QueryUniqueCampuscode(Campuscode string, Campusname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Campuscode) from campus WHERE (Campuscode=? or Campusname=?);", Campuscode, Campusname)
	return num
}

//修改校区
func UpdateCampus(cps *basicset.Campus, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Campus{}, "campus").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateCampus|修改校区:")
	return inerr
}

//删除校区
func DeleteCampus(cps *basicset.Campus, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//第一步删除校区下所有的教室
		_, inerr = tran.Exec("delete from classrooms where Floorsid in(select Id from floors where Buildingid in (select Id from building where Campusid in (select Id from campus where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第二步删除校区下所有的楼层
		_, inerr = tran.Exec("delete from floors where Buildingid in(select Id from building where Campusid in (select Id from campus where Id=?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第三步删除校区下所有的楼栋
		_, inerr = tran.Exec("delete from building where Campusid in(select Id from campus where Id=?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第四步删除校区表中的记录
		_, inerr = tran.Exec("delete from campus where Id =?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第五步修改绑定校区的学院记录
		_, inerr = tran.Exec("update college set Campusid=0 where Campusid=?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

//添加楼栋
func AddBuilding(bd *basicset.Building, dbmap *gorp.DbMap) (inerr error) {
	if bd.Campusid > 0 {
		dbmap.AddTableWithName(basicset.Building{}, "building").SetKeys(true, "Id")
		inerr = dbmap.Insert(bd)
		fmt.Println(inerr)
		core.CheckErr(inerr, "basicsetDataAccess|AddBuilding|添加楼栋:")
	}
	return inerr
}

//查询楼栋代码是否唯一
func QueryUniqueBuildingcode(Buildingcode, Buildingname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Buildingcode) from building WHERE (Buildingcode=? or Buildingname=?);", Buildingcode, Buildingname)
	return num
}

//修改楼栋
func UpdateBuilding(cps *basicset.Building, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Building{}, "building").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateBuilding|修改校区:")
	return inerr
}

//删除楼栋
func DeleteBuilding(cps *basicset.Building, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//第一步删除校区下所有的教室
		_, inerr = tran.Exec("delete from classrooms where Floorsid in(select Id from floors where Buildingid in (select Id from building where Id =?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第二步删除校区下所有的楼层
		_, inerr = tran.Exec("delete from floors where Buildingid in(select Id from building where Id =?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第三步删除校区下所有的楼栋
		_, inerr = tran.Exec("delete from building where Id =?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

//添加楼层
func AddFloors(fl *basicset.Floors, dbmap *gorp.DbMap) (inerr error) {
	fmt.Println("AddFloors添加楼层方法进入中......")
	if fl.Buildingid > 0 {
		dbmap.AddTableWithName(basicset.Floors{}, "floors").SetKeys(true, "Id")
		inerr = dbmap.Insert(fl)
		if inerr == nil {
			dbmap.Exec("update building set Floorsnumber=Floorsnumber+1 where Id=?;", fl.Buildingid)
		}
		fmt.Println(inerr)
		core.CheckErr(inerr, "basicsetDataAccess|AddFloors|添加楼层:")
	}
	return inerr
}

//查询楼层代码是否唯一
func QueryUniqueFloorscode(Floorscode, Floorname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Floorscode) from floors WHERE (Floorscode=? or Floorname=?);", Floorscode, Floorname)
	return num
}

//修改楼层
func UpdateFloors(cps *basicset.Floors, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Floors{}, "floors").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateFloors|修改楼层:")
	return inerr
}

//删除楼层
func DeleteFloors(cps *basicset.Floors, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//第一步删除校区下所有的教室
		_, inerr = tran.Exec("delete from classrooms where Floorsid in(select Id from floors where Id =?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第二步删除校区下所有的楼层
		_, inerr = tran.Exec("delete from floors where Id =?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

//添加教室
func AddClassrooms(cr *basicset.Classrooms, dbmap *gorp.DbMap) (inerr error) {
	fmt.Println("AddClassrooms添加教室方法进入中......")
	if cr.Floorsid > 0 {
		dbmap.AddTableWithName(basicset.Classrooms{}, "classrooms").SetKeys(true, "Id")
		inerr = dbmap.Insert(cr)
		if inerr == nil {
			dbmap.Exec("update floors set Classroomnumber=Classroomnumber+1 where Id=?;", cr.Floorsid)    //更新楼层的教室数
			Buildingid, _ := dbmap.SelectInt("select Buildingid from floors where Id=?;", cr.Floorsid)    //获取楼栋的ID
			dbmap.Exec("update building set Classroomsnumber=Classroomsnumber+1 where Id=?;", Buildingid) //更新楼栋中的教室数
		}
		core.CheckErr(inerr, "basicsetDataAccess|AddClassrooms|添加教室:")
		//fmt.Println(inerr)
	}
	return inerr
}

//查询教室代码是否唯一
func QueryUniqueClassroomscode(Classroomscode, Classroomsname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Classroomscode) from Classrooms WHERE (Classroomscode=? or Classroomsname=?);", Classroomscode, Classroomsname)
	return num
}

//修改教室
func UpdateClassrooms(cps *basicset.Classrooms, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Classrooms{}, "classrooms").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateClassrooms|修改教室:")
	return inerr
}

//删除教室
func DeleteClassrooms(cps *basicset.Classrooms, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Classrooms{}, "classrooms").SetKeys(true, "Id")
	_, inerr = dbmap.Delete(cps)
	core.CheckErr(inerr, "basicsetDataAccess|DeleteClassrooms|删除教室:")
	return inerr
}

//添加学院
func AddCollege(bd *basicset.College, dbmap *gorp.DbMap) (inerr error) {
	if bd.Campusid > 0 {
		dbmap.AddTableWithName(basicset.College{}, "college").SetKeys(true, "Id")
		inerr = dbmap.Insert(bd)
		core.CheckErr(inerr, "basicsetDataAccess|AddCollege|添加学院:")
	}
	return inerr
}

//查询学院代码是否唯一
func QueryUniqueCollegecode(Collegecode, Collegename string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Collegecode) from College WHERE (Collegecode=? or Collegename=?);", Collegecode, Collegename)
	return num
}

//修改学院
func UpdateCollege(cps *basicset.College, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.College{}, "college").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateCollege|修改学院:")
	return inerr
}

//删除学院
func DeleteCollege(cps *basicset.College, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//第一步删除学院下所有的学生用户
		_, inerr = tran.Exec("delete from users where Id in(select Id from students where Classesid in(select Id from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?))))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第二步删除学院下所有的学生
		_, inerr = tran.Exec("delete from students where Classesid in(select Id from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第三步删除学院下所有班级的课程点到数据
		_, inerr = tran.Exec("delete from pointtos where curriculumclassroomchaptercentreId in (select Id from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?)))))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第四步删除学院下所有班级的教室上课记录
		_, inerr = tran.Exec("delete from teachingrecord where curriculumclassroomchaptercentreId in (select Id from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?)))))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第五步删除学院下所有班级的课程计划详情数据
		_, inerr = tran.Exec("delete from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?))))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第六步删除学院下所有班级的课程计划主数据
		_, inerr = tran.Exec("delete from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第七步删除学院下所有的班级
		_, inerr = tran.Exec("delete from Classes where Majorid in(select Id from Major where Collegeid in(select Id from College where Id=?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第八步删除学院下所有的专业
		_, inerr = tran.Exec("delete from Major where Collegeid in(select Id from College where Id=?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第九步删除学院下所有的学院
		_, inerr = tran.Exec("delete from College where Id=?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

//查询专业代码是否唯一
func QueryUniqueMajorcode(Majorcode, Majorname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Majorcode) from Major WHERE (Majorcode=? or Majorname=?);", Majorcode, Majorname)
	return num
}

//添加专业
func AddMajor(fl *basicset.Major, dbmap *gorp.DbMap) (inerr error) {
	if fl.Collegeid > 0 {
		dbmap.AddTableWithName(basicset.Major{}, "major").SetKeys(true, "Id")
		inerr = dbmap.Insert(fl)
		core.CheckErr(inerr, "basicsetDataAccess|AddMajor|添加专业:")
	}
	return inerr
}

//修改专业
func UpdateMajor(cps *basicset.Major, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Major{}, "major").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateMajor|修改专业:")
	return inerr
}

//删除专业
func DeleteMajor(cps *basicset.Major, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//第一步删除学院下所有的学生用户
		_, inerr = tran.Exec("delete from users where Id in(select Id from students where Classesid in(select Id from Classes where Majorid in(select Id from Major where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第二步删除学院下所有的学生
		_, inerr = tran.Exec("delete from students where Classesid in(select Id from Classes where Majorid in(select Id from Major where Id=?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第三步删除学院下所有班级的课程点到数据
		_, inerr = tran.Exec("delete from pointtos where curriculumclassroomchaptercentreId in (select Id from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Id=?))))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第四步删除学院下所有班级的教室上课记录
		_, inerr = tran.Exec("delete from teachingrecord where curriculumclassroomchaptercentreId in (select Id from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Id=?))))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第五步删除学院下所有班级的课程计划详情数据
		_, inerr = tran.Exec("delete from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第六步删除学院下所有班级的课程计划主数据
		_, inerr = tran.Exec("delete from curriculumsclasscentre where ClassesId in(select Id from Classes where Majorid in(select Id from Major where Id=?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第七步删除学院下所有的班级
		_, inerr = tran.Exec("delete from Classes where Majorid in(select Id from Major where Id=?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第八步删除学院下所有的专业
		_, inerr = tran.Exec("delete from Major where Id=?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

//查询专业代码是否唯一
func QueryUniqueClassescode(Classescode, Classesname string, dbmap *gorp.DbMap) (num int64) {
	num, _ = dbmap.SelectInt("select count(Classescode) from Classes WHERE (Classescode=? or Classesname=?);", Classescode, Classesname)
	return num
}

//添加班级
func AddClasses(cr *basicset.Classes, dbmap *gorp.DbMap) (inerr error) {
	if cr.Majorid > 0 {
		dbmap.AddTableWithName(basicset.Classes{}, "classes").SetKeys(true, "Id")
		inerr = dbmap.Insert(cr)
		core.CheckErr(inerr, "basicsetDataAccess|AddClasses|添加班级:")
	}
	return inerr
}

//修改班级
func UpdateClasses(cps *basicset.Classes, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(basicset.Classes{}, "classes").SetKeys(true, "Id")
	_, inerr = dbmap.Update(cps)
	core.CheckErr(inerr, "basicsetDataAccess|UpdateClasses|修改班级:")
	return inerr
}

//删除班级
func DeleteClasses(cps *basicset.Classes, dbmap *gorp.DbMap) (inerr error) {
	tran, inerr := dbmap.Begin()
	if inerr == nil {
		//第一步删除学院下所有的学生用户
		_, inerr = tran.Exec("delete from users where Id in(select Id from students where Classesid in(select Id from Classes where Id=?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第二步删除学院下所有的学生
		_, inerr = tran.Exec("delete from students where Classesid in(select Id from Classes where Id=?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第三步删除学院下所有班级的课程点到数据
		_, inerr = tran.Exec("delete from pointtos where curriculumclassroomchaptercentreId in (select Id from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第四步删除学院下所有班级的教室上课记录
		_, inerr = tran.Exec("delete from teachingrecord where curriculumclassroomchaptercentreId in (select Id from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Id=?)))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第五步删除学院下所有班级的课程计划详情数据
		_, inerr = tran.Exec("delete from curriculumclassroomchaptercentre where curriculumsclasscentreid in(select Id from curriculumsclasscentre where ClassesId in(select Id from Classes where Id=?))", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第六步删除学院下所有班级的课程计划主数据
		_, inerr = tran.Exec("delete from curriculumsclasscentre where ClassesId in(select Id from Classes where Id=?)", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		//第七步删除学院下所有的班级
		_, inerr = tran.Exec("delete from Classes where Id=?", cps.Id)
		if inerr != nil {
			tran.Rollback()
			return inerr
		}
		inerr = tran.Commit()
	}
	return inerr
}

/*
获取所有校区
根据校区代码查询某校区
查询所有校区
*/
func QueryCampus(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.Getcampus) {
	sql := "select Id as Campusid,Campusname,Campusicon,Campuscode,Campusnums from campus where 1=1"
	sql = sql + GetWhereString(ws) + core.GetLimitString(pg) + ";"
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryCampus|获取所有校区:")
	return list
}

/*
获取所有校区
根据校区代码查询某校区
查询所有校区
*/
func QueryCampusPG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from campus where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryCampusPG|系统后台获取所有校区:")
	if sqlerr == nil {
		if countint > 0 {
			var list []viewmodel.Getcampus
			sql := "select Id as Campusid,Campusname,Campusicon,Campuscode,Campusnums from campus where 1=1"
			sql = sql + wheresql + core.GetLimitString(pg) + ";"
			_, sqlerr := dbmap.Select(&list, sql)
			core.CheckErr(sqlerr, "basicsetDataAccess|QueryCampusPG|系统后台获取所有校区:")
			rd.Rcode = "1000"
			pg.PageCount = int(countint)
			pg.PageData = list
			rd.Result = pg
		} else {
			rd.Rcode = "1009"
			rd.Reason = "未找到数据"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未找到数据"
	}

	return rd
}

/*
获取所有楼栋
根据楼栋代码或者校区ID查询楼栋信息
*/
func QueryBuilding(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.Getbuilding) {
	sql := "select Id as Buildingid,Campusid,Buildingname,Buildingicon,Floorsnumber,Classroomsnumber,Buildingcode from building where 1=1 "
	sql = sql + GetWhereString(ws) + " order by Campusid" + core.GetLimitString(pg) + ";"
	fmt.Println(sql)
	_, Campuserr := dbmap.Select(&list, sql)
	if Campuserr != nil {
		fmt.Println("QueryBuilding:", Campuserr)
		core.CheckErr(Campuserr, "basicsetDataAccess|QueryBuilding|获取所有楼栋:")
		return nil
	}
	return list
}

/*
获取所有楼栋
根据楼栋代码或者校区ID查询楼栋信息
*/
func QueryBuildingPG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from building where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryBuildingPG|系统后台获取所有校区:")
	if sqlerr == nil {
		if countint > 0 {
			var list []viewmodel.Getbuilding
			sql := "select Id as Buildingid,Campusid,Buildingname,Buildingicon,Floorsnumber,Classroomsnumber,Buildingcode from building where 1=1"
			sql = sql + wheresql + " order by Campusid" + core.GetLimitString(pg) + ";"
			_, sqlerr := dbmap.Select(&list, sql)
			core.CheckErr(sqlerr, "basicsetDataAccess|QueryBuildingPG|系统后台获取所有校区:")
			rd.Rcode = "1000"
			pg.PageCount = int(countint)
			pg.PageData = list
			rd.Result = pg
		} else {
			rd.Rcode = "1000"
			rd.Reason = "未找到数据"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "未找到数据"
	}

	return rd
}

//获取所有楼层
func QueryFloors(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.Getfloors) {
	sql := "select Id as Floorsid,Floorscode,Buildingid,Floorname,Classroomnumber,Sumnumber from floors where 1=1 "
	sql = sql + GetWhereString(ws) + " order by Buildingid" + core.GetLimitString(pg) + ";"
	fmt.Println(sql)
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryFloors|获取所有楼层:")
	return list
}

//获取所有楼层
func QueryFloorsPG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from floors where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryFloorsPG|系统后台获取所有楼层:")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.Getfloors
		sql := "select Id as Floorsid,Floorscode,Buildingid,Floorname,Classroomnumber,FloorsImage, Sumnumber from floors where 1=1 "
		sql = sql + wheresql + " order by Buildingid" + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "basicsetDataAccess|QueryFloorsPG|系统后台获取所有楼层:")
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

//获取所有教室
func QueryClassrooms(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.Getclassrooms) {
	sql := "select Id as Classroomid,Classroomsname,Floorsid,Classroomicon,Seatsnumbers,Sumnumbers,Classroomstype,Classroomstate,Collectionnumbers from classrooms where 1=1 "
	sql = sql + GetWhereString(ws) + " order by Floorsid" + core.GetLimitString(pg) + ";"
	fmt.Println(sql)
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryClassrooms|获取所有教室:")
	return list
}

//获取所有教室
func QueryClassroomsPG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from classrooms where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryClassroomsPG|系统后台获取所有教室:")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.Getclassrooms
		sql := "select Id as Classroomid,Classroomsname,Floorsid,Classroomicon,Seatsnumbers,Sumnumbers,Classroomstype,Classroomstate,Collectionnumbers,Classroomscode from classrooms where 1=1 "
		sql = sql + wheresql + " order by Floorsid" + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "basicsetDataAccess|QueryClassroomsPG|系统后台获取所有教室:")
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

//查询某一教室的相关信息
func QueryClassroomsInfo(ws viewmodel.QueryBasicsetWhere, dbmap *gorp.DbMap) (list viewmodel.QueryClassroomInfo) {
	sql := "select cs.Id as Classroomid,cs.Classroomsname,c.Campusname,bd.Buildingname,fl.Floorname,cs.Classroomicon,cs.Classroomstate,cs.Seatsnumbers,cs.Sumnumbers from classrooms as cs inner join floors as fl on cs.Floorsid=fl.Id inner join building as bd on bd.Id=fl.Buildingid inner join Campus as c on c.Id=bd.Campusid where cs.Id=?;"
	Campuserr := dbmap.SelectOne(&list, sql, ws.Classroomid)
	fmt.Println(sql, ws.Classroomid)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryClassroomsInfo|查询某一教室的相关信息:")
	return list
}

/*
获取所有学院
*/
func QueryCollege(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []basicset.College) {
	sql := "select * from college where 1=1 "
	sql = sql + GetWhereString(ws) + " order by Campusid" + core.GetLimitString(pg) + ";"
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryCollege|获取所有学院:")
	return list
}

/*
获取所有学院
*/
func QueryCollegePG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from college where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryCollegePG|系统后台获取所有学院:")
	if sqlerr == nil && countint > 0 {
		var list []basicset.College
		sql := "select * from college where 1=1 "
		sql = sql + wheresql + " order by Campusid" + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "basicsetDataAccess|QueryCollegePG|系统后台获取所有学院:")
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

//获取所有专业
func QueryMajor(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []basicset.Major) {
	sql := "select * from major where 1=1 "
	sql = sql + GetWhereString(ws) + " order by Collegeid" + core.GetLimitString(pg) + ";"
	fmt.Println(sql)
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryMajor|获取所有专业:")
	return list
}

//获取所有专业
func QueryMajorPG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from major where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryMajorPG|系统后台获取所有专业:")
	if sqlerr == nil && countint > 0 {
		var list []basicset.Major
		sql := "select * from major where 1=1 "
		sql = sql + wheresql + " order by Collegeid" + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "basicsetDataAccess|QueryMajorPG|系统后台获取所有专业:")
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

//获取所有班级
func QueryClasses(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []basicset.Classes) {
	sql := "select * from classes where 1=1 "
	sql = sql + GetWhereString(ws) + " order by Majorid" + core.GetLimitString(pg) + ";"
	fmt.Println(sql)
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryClasses|获取所有班级:")
	return list
}

//获取所有班级
func QueryClassesPG(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from classes where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryClassesPG|系统后台获取所有班级:")
	if sqlerr == nil && countint > 0 {
		var list []basicset.Classes
		sql := "select * from classes where 1=1 "
		sql = sql + wheresql + " order by Majorid" + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql)
		core.CheckErr(sqlerr, "basicsetDataAccess|QueryClassesPG|系统后台获取所有班级:")
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

//获取所有授课老师
func QueryTeachers(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.Teacher) {
	sql := "select t.Id as Usersid,us.Loginuser,us.Nickname,t.Collegeid,ifnull(t.Majorid,0)as Majorid from teacher as t inner join users as us on t.Id=us.Id "
	sql = sql + " where 1=1 "
	if ws.Collegeid > 0 {
		sql = sql + " and t.Collegeid=" + strconv.Itoa(ws.Collegeid)
	}
	if ws.Majorid > 0 {
		sql = sql + " and t.Majorid=" + strconv.Itoa(ws.Majorid)
	}
	sql = sql + " group by t.Id,us.Loginuser,us.Nickname" + core.GetLimitString(pg) + ";"
	fmt.Println(sql)
	_, Campuserr := dbmap.Select(&list, sql)
	core.CheckErr(Campuserr, "basicsetDataAccess|QueryTeachers|获取所有授课老师:")
	return list
}

/*
根据校区Id、楼栋id、楼层id、用户Id等相关信息查询
*/
func QueryClassroominfo(ws viewmodel.QueryBasicsetWhere, pg core.PageData, dbmap *gorp.DbMap) (rd core.Returndata) {
	countsql := "select count(*) from classrooms as cr inner join floors as fs on cr.Floorsid=fs.Id inner join building as bd on fs.Buildingid=bd.Id inner join campus as cs on bd.Campusid=cs.Id left join classroomcollection as crc on (crc.Classroomid=cr.Id and crc.Usersid=" + strconv.Itoa(ws.Usersid) + ") where 1=1"
	wheresql := GetWhereString(ws)
	var sqlerr error
	countint, sqlerr := dbmap.SelectInt(countsql + wheresql)
	core.CheckErr(sqlerr, "basicsetDataAccess|QueryClassroominfo|根据校区Id、楼栋id、楼层id、用户Id等相关信息查询:")
	if sqlerr == nil && countint > 0 {
		var list []viewmodel.QueryResultClassroom
		sql1 := "select cr.Id as Classroomid,cr.Classroomsname,cr.Classroomicon,cr.Seatsnumbers,cr.Sumnumbers,cr.Classroomstype,cr.Classroomstate,cr.Collectionnumbers,cr.Notes,fs.Floorname,bd.Buildingname,cs.Campusname,CASE when crc.Id>0 then 1 else 0  end as State"
		sql1 = sql1 + " from classrooms as cr inner join floors as fs on cr.Floorsid=fs.Id inner join building as bd on fs.Buildingid=bd.Id inner join campus as cs on bd.Campusid=cs.Id left join classroomcollection as crc on (crc.Classroomid=cr.Id and crc.Usersid=" + strconv.Itoa(ws.Usersid) + ") where 1=1 "
		sql1 = sql1 + wheresql + core.GetLimitString(pg) + ";"
		_, sqlerr := dbmap.Select(&list, sql1)
		core.CheckErr(sqlerr, "basicsetDataAccess|QueryClassroominfo|根据校区Id、楼栋id、楼层id、用户Id等相关信息查询:")
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
根据校区Id、楼栋id、楼层id、查询教室相关数据状态
*/
func GetSelectPeoples(ws viewmodel.QueryBasicsetWhere, dbmap *gorp.DbMap) (list []viewmodel.QueryPeoples) {
	sql1 := "select cr.Sumnumbers,cr.Id as Classroomid,bd.Id as BuildingId,bd.Buildingname as BuildingName,fs.Id as FloorId,fs.Floorname as FloorName,fs.FloorsImage,cr.Classroomsname as ClassroomName,cr.Seatsnumbers,cr.Classroomstate as ClassroomState,cr.Classroomstype,cr.Classroomicon from "
	sql1 = sql1 + " classrooms as cr inner join floors as fs on cr.Floorsid=fs.Id inner join building as bd on bd.Id=fs.Buildingid"
	sql1 = sql1 + " inner join campus as cp on cp.Id=bd.Campusid where 1=1"
	sql1 = `
SELECT TRoom.Sumnumbers,
	TRoom.Id AS Classroomid,
	TBuilding.Id AS BuildingId,
	TBuilding.Buildingname AS BuildingName,
	TFloor.Id AS FloorId,
	TFloor.Floorname AS FloorName,
	TFloor.FloorsImage,
	TRoom.Classroomsname AS ClassroomName,
	TRoom.Seatsnumbers,
	TRoom.Classroomstate AS ClassroomState,
	TRoom.Classroomstype,
	TRoom.Classroomicon
FROM classrooms AS TRoom
	INNER JOIN floors AS TFloor ON TRoom.Floorsid = TFloor.Id
	INNER JOIN building AS TBuilding ON TBuilding.Id = TFloor.Buildingid
	INNER JOIN campus AS TCampus ON TCampus.Id = TBuilding.Campusid
WHERE 1 = 1
`
	if ws.Campusids != "" {
		// sql1 = sql1 + " and cp.Id in(" + ws.Campusids + ")"
		sql1 += fmt.Sprintf(" AND TCampus.Id IN (%s) \n", ws.Campusids)
	}
	if ws.Buildingids != "" {
		// sql1 = sql1 + " and bd.Id in(" + ws.Buildingids + ")"
		sql1 += fmt.Sprintf(" AND TBuilding.Id IN (%s) \n", ws.Buildingids)
	}
	if ws.Floorsids != "" {
		// sql1 = sql1 + " and fs.Id in(" + ws.Floorsids + ")"
		sql1 += fmt.Sprintf(" AND TFloor.Id IN (%s) \n", ws.Floorsids)
	}
	// sql1 = sql1 + " group by bd.Id,fs.Id,cr.Id,cr.Classroomsname,cr.Seatsnumbers,cr.Classroomstate,cr.Classroomstype,cr.Classroomicon;"
	sql1 += `
GROUP BY
	TBuilding.Id,
	TFloor.Id,
	TRoom.Id,
	TRoom.Classroomsname,
	TRoom.Seatsnumbers,
	TRoom.Classroomstate,
	TRoom.Classroomstype,
	TRoom.Classroomicon
ORDER BY TBuilding.Buildingname, TFloor.Floorname, TRoom.Classroomsname;
`
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "basicsetDataAccess|GetSelectPeoples|根据校区Id、楼栋id、楼层id、查询教室相关数据状态")
	return list
}

/*
查询教室内所有人的记录
*/
func QueryClassRoomPeopleInfo(ws viewmodel.QueryClassRoomPeopleInfo, dbmap *gorp.DbMap) (list []viewmodel.QueryClassRoomPeopleInfo) {
	sql1 := "select us.Userheadimg,us.Nickname,crdn.Usersid,min(crdn.Createdate)as Createdate,max(crdn.Closedate)as Closedate,((UNIX_TIMESTAMP(max(crdn.Closedate))-UNIX_TIMESTAMP(min(crdn.Createdate))))as DateLength,crdn.Xy,crdn.X,crdn.Y "
	sql1 = sql1 + " from classroomdetailsnums as crdn inner join users as us on us.Id=crdn.Usersid where 1=1"
	if ws.Begindate != "" {
		sql1 = sql1 + " and crdn.Createdate>='" + ws.Begindate + "'"
	} else {
		sql1 = sql1 + " and crdn.Createdate>='" + time.Now().Format("2006-01-02") + " 00:00:00'"
	}
	if ws.Enddate != "" {
		sql1 = sql1 + " and crdn.Closedate<='" + ws.Enddate + "'"
	} else {
		sql1 = sql1 + " and crdn.Closedate<='" + time.Now().Format("2006-01-02") + " 23:59:59'"
	}
	sql1 = sql1 + " and crdn.Classroomid=" + strconv.Itoa(ws.Classroomid)
	sql1 = sql1 + " group by crdn.Usersid;"
	fmt.Println(sql1)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryClassRoomPeopleInfo|查询教室内所有人的记录")
	return list
}

/*
查询教室内所有人的记录的汇总统计
*/
func QueryClassRoomPeopleCount(ws viewmodel.QueryClassRoomPeopleInfo, dbmap *gorp.DbMap) (list []viewmodel.PeopleCount) {
	sql1 := "select count(crdn.Usersid) as Sumnumbers,DATE_FORMAT(Closedate,'%Y%m%d')as Dateymd,DATE_FORMAT(Closedate,'%H') as Dateh"
	sql1 = sql1 + " from classroomdetailsnums as crdn inner join classrooms as cr on crdn.Classroomid=cr.Id "
	if ws.Begindate != "" {
		sql1 = sql1 + " and crdn.Createdate>='" + ws.Begindate + "'"
	} else {
		//		time.Now().AddDate(0,0,-3)
		sql1 = sql1 + " and crdn.Createdate>='" + time.Now().AddDate(0, 0, -7).Format("2006-01-02") + " 00:00:00'"
	}
	if ws.Enddate != "" {
		sql1 = sql1 + " and crdn.Closedate<='" + ws.Enddate + "'"
	} else {
		sql1 = sql1 + " and crdn.Closedate<='" + time.Now().Format("2006-01-02") + " 23:59:59'"
	}
	sql1 = sql1 + " and cr.Id=" + strconv.Itoa(ws.Classroomid)
	sql1 = sql1 + " group by DATE_FORMAT(Closedate,'%Y%m%d'),DATE_FORMAT(Closedate,'%H');"
	fmt.Println(sql1)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryClassRoomPeopleCount|查询教室内所有人的记录的汇总统计")
	return list
}

/*
根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计
*/
func QueryStreamPeoplesAnalysis(ws viewmodel.QueryStreamPeoplesWhere, dbmap *gorp.DbMap) (rd core.Returndata) {
	sql1 := "select Valname,count(*)as Valcount from (select crdn.Usersid as Valcount,{0} as Valname from classroomdetailsnums as crdn inner join users as us on crdn.Usersid=us.Id right join students as sds on us.Id=sds.Id "
	sql1 = sql1 + " inner join classes as cs on sds.Classesid=cs.Id inner join major as mj on cs.Majorid=mj.Id inner join College as cg on cg.Id=mj.Collegeid inner join classrooms as cr on crdn.Classroomid=cr.Id right join "
	sql1 = sql1 + " floors as fs on cr.Floorsid=fs.Id right join building as bd on fs.Buildingid=bd.Id right join campus as cp on bd.campusid=cp.Id where 1=1 "
	Replacestr := "cp.Campusname"
	if ws.Campusid > 0 {
		sql1 = sql1 + " and cp.Id=" + strconv.Itoa(ws.Campusid)
		Replacestr = "bd.Buildingname"
	}
	if ws.Buildingid > 0 {
		sql1 = sql1 + " and bd.Id=" + strconv.Itoa(ws.Buildingid)
		Replacestr = "fs.Floorname"
	}
	if ws.Floorsid > 0 {
		sql1 = sql1 + " and fs.Id=" + strconv.Itoa(ws.Floorsid)
	}
	if ws.Begindate != "" {
		sql1 = sql1 + " and crdn.Createdate>='" + ws.Begindate + "'"
	}
	if ws.Enddate != "" {
		sql1 = sql1 + " and crdn.Createdate<='" + ws.Enddate + "'"
	}
	sql1 = sql1 + " group by crdn.Usersid,{0}) as ttl group by Valname;"
	arrlist := make([]interface{}, 5)
	execsql := strings.Replace(sql1, "{0}", Replacestr, -1)
	var list []viewmodel.ListStreamPeoplesAnalysis
	_, sserr1 := dbmap.Select(&list, execsql)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryStreamPeoplesAnalysis|根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计")
	arrlist[0] = list
	Replacestr = "cg.collegename" //按学院
	execsql = strings.Replace(sql1, "{0}", Replacestr, -1)
	list = nil
	_, sserr1 = dbmap.Select(&list, execsql)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryStreamPeoplesAnalysis|根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计")
	arrlist[1] = list
	Replacestr = "mj.Majorname" //按专业
	execsql = strings.Replace(sql1, "{0}", Replacestr, -1)
	list = nil

	_, sserr1 = dbmap.Select(&list, execsql)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryStreamPeoplesAnalysis|根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计")
	arrlist[2] = list
	Replacestr = "cs.Classesname" //按班级
	execsql = strings.Replace(sql1, "{0}", Replacestr, -1)
	list = nil

	_, sserr1 = dbmap.Select(&list, execsql)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryStreamPeoplesAnalysis|根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计")
	arrlist[3] = list
	Replacestr = "us.Usersex" //按性别
	execsql = strings.Replace(sql1, "{0}", Replacestr, -1)
	list = nil

	_, sserr1 = dbmap.Select(&list, execsql)
	core.CheckErr(sserr1, "basicsetDataAccess|QueryStreamPeoplesAnalysis|根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计")
	arrlist[4] = list
	rd.Result = arrlist
	rd.Rcode = "1000"
	return rd
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
