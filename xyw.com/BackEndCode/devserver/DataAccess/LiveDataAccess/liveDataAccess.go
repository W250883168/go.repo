package liveDataAccess

import (
	"fmt"
	"strconv"

	"gopkg.in/gorp.v1"

	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/live"
	"dev.project/BackEndCode/devserver/viewmodel"
)

//添加课程班级章节附件关联表数据
func AddLives(ct *live.Lives, dbmap *gorp.DbMap) (inerr error) {
	dbmap.AddTableWithName(live.Lives{}, "lives").SetKeys(true, "Id")
	inerr = dbmap.Insert(ct)
	core.CheckErr(inerr, "liveDataAccess|AddLives|添加课程班级章节附件关联表数据")
	return inerr
}

/*
查询直播/录播数据列表
*/
func QueryLiveWhereAll(lg viewmodel.QueryLivesWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.QueryAttentionRecordList) {
	sql1 := "select ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,ccc.Classesid,cs.Classesname,mr.Majorname,cg.Collegename,us.Nickname"
	sql1 = sql1 + " from curriculums as cc inner join curriculumsclasscentre as ccc on cc.Id=ccc.Curriculumsid inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid"
	sql1 = sql1 + " inner join chapters as cts on (cts.Curriculumsid=cc.Id and cts.Id=cccc.Chaptersid) inner join lives as ls on ls.Curriculumclassroomchaptercentreid=cccc.Id inner join classes as cs on cs.Id=ccc.Classesid inner join major as mr on mr.Id=cs.Majorid inner join college as cg on cg.Id=mr.Collegeid inner join users as us on us.Id=ccc.Usersid where 1=1 "
	sql1 = sql1 + GetQueryLiveWhere(lg) + " group by ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,ccc.Classesid,cs.Classesname" + core.GetLimitString(pg)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "liveDataAccess|QueryLiveWhereAll|查询直播/录播数据列表")
	return list
}

/*
查询我缺课的课程视频一
*/
func QueryMyAbsentList(lg viewmodel.QueryLivesWhere, pg core.PageData, dbmap *gorp.DbMap) (list []viewmodel.QueryAttentionRecordList) {
	sql1 := "select ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,ccc.Classesid as Classesid,ps.State as PsState"
	sql1 = sql1 + " from curriculums as cc inner join curriculumsclasscentre as ccc on cc.Id=ccc.Curriculumsid "
	sql1 = sql1 + " inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid inner join lives as ls on ls.Curriculumclassroomchaptercentreid=cccc.Id"
	sql1 = sql1 + " inner join pointtos as ps on ps.Curriculumclassroomchaptercentreid=cccc.Id inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id where 1=1 "
	sql1 = sql1 + GetQueryLiveWhere(lg) + " group by ccc.Curriculumsid,cc.Curriculumname,cc.Curriculumicon,cc.Curriculumstype,cc.Chaptercount,ccc.Newchapter,ccc.Newlivechapter,cc.Subjectcode,ccc.PlaySumnum,ccc.DownloadSumnum,ccc.FollowSumnum,ccc.Classesid" + core.GetLimitString(pg)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "liveDataAccess|QueryMyAbsentList|查询我缺课的课程视频一")
	return list
}

/*
查询我缺课的课程视频二
*/
func QueryMyAbsentInfo(lg viewmodel.QueryLivesWhere, dbmap *gorp.DbMap) (list []viewmodel.QueryLiveInfos) {
	sql1 := "select ls.Id as Liveid,cccc.Id as Ccccid,cts.Chaptername,cts.Chaptericon,cts.Chapterdetails,ls.Liveinfo,ls.Livetitile,ls.Livepath1,ls.Livepath2,ls.Trailerpath,ls.Coverimage,"
	sql1 = sql1 + " ls.Livestate,ls.Livetype,ls.Whenlong,ls.Recommendread,ls.Iscomment,ls.Ischeckcomment,ls.Isdownload,ls.Downloadnum,ls.Playnum,ps.State,us.Nickname,us.Userheadimg,cccc.Usersid,cccc.Begindate,cr.Classroomsname"
	sql1 = sql1 + " from curriculumsclasscentre as ccc inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid"
	sql1 = sql1 + " inner join chapters as cts on cts.Id=cccc.Chaptersid inner join lives as ls on ls.Curriculumclassroomchaptercentreid=cccc.Id "
	sql1 = sql1 + " inner join pointtos as ps on ps.Curriculumclassroomchaptercentreid=cccc.Id inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id inner join users as us on cccc.Usersid=us.Id "
	sql1 = sql1 + " inner join classrooms as cr on tr.Classroomid=cr.Id where 1=1"
	sql1 = sql1 + GetQueryLiveWhere(lg) + ";"
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "liveDataAccess|QueryMyAbsentInfo|查询我缺课的课程视频二")
	return list
}

/*
查询课程下的详细视频记录
*/
func QueryLiveInfo(lg viewmodel.QueryLivesWhere, dbmap *gorp.DbMap) (list []viewmodel.QueryLiveInfos) {
	sql1 := "select ls.Id as Liveid,cccc.Id as Ccccid,cts.Chaptername,cts.Chaptericon,cts.Chapterdetails,ls.Liveinfo,ls.Livetitile,ls.Livepath1,ls.Livepath2,ls.Trailerpath,ls.Coverimage,"
	sql1 = sql1 + " ls.Livestate,ls.Livetype,ls.Whenlong,ls.Recommendread,ls.Iscomment,ls.Ischeckcomment,ls.Isdownload,ls.Downloadnum,ls.Playnum,us.Nickname,us.Userheadimg,cccc.Usersid,cccc.Begindate,cr.Classroomsname"
	sql1 = sql1 + " from curriculumsclasscentre as ccc inner join curriculumclassroomchaptercentre as cccc on ccc.Id=cccc.Curriculumsclasscentreid"
	sql1 = sql1 + " inner join chapters as cts on cts.Id=cccc.Chaptersid inner join lives as ls on ls.Curriculumclassroomchaptercentreid=cccc.Id inner join users as us on cccc.Usersid=us.Id"
	sql1 = sql1 + " inner join teachingrecord as tr on tr.Curriculumclassroomchaptercentreid=cccc.Id inner join classrooms as cr on tr.Classroomid=cr.Id where 1=1 "
	sql1 = sql1 + GetQueryLiveWhere(lg) + ";"
	fmt.Println(sql1)
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "liveDataAccess|QueryLiveInfo|查询课程下的详细视频记录")
	return list
}

/*
查询课程上传的附件
*/
func QueryLiveEnclosure(lg viewmodel.QueryLivesWhere, dbmap *gorp.DbMap) (list []viewmodel.Enclosure) {
	sql1 := "select ccc.Id as Enclosureid,ccc.Enclosurename,ccc.Enclosuretype,ccc.Enclosuresize,ccc.EnclosureVirtualPath,ccc.Createdate,ccc.Enclosureicon from enclosure as ccc where 1=1 "
	sql1 = sql1 + GetQueryLiveWhere(lg) + ";"
	_, sserr1 := dbmap.Select(&list, sql1)
	core.CheckErr(sserr1, "liveDataAccess|QueryLiveEnclosure|查询课程上传的附件")
	return list
}

/*
更新课程章节的播放次数
*/
func UpdateLivePlayNum(lg viewmodel.QueryLivesWhere, dbmap *gorp.DbMap) (rd core.Returndata) {
	if lg.Ccccid > 0 && lg.PsUserid > 0 && lg.Liveid > 0 {
		ct, inerr := dbmap.SelectInt("select count(*) as Ct from playrecord where Livesid=? and Usersid=?;", lg.Liveid, lg.PsUserid)
		core.CheckErr(inerr, "liveDataAccess|UpdateLivePlayNum|更新课程章节的播放次数")
		if ct < 1 {
			dbmap.AddTableWithName(live.Playrecord{}, "playrecord").SetKeys(true, "Id")
			pr := live.Playrecord{Livesid: lg.Liveid, Usersid: lg.PsUserid}
			inerr = dbmap.Insert(&pr)
			core.CheckErr(inerr, "liveDataAccess|UpdateLivePlayNum|更新课程章节的播放次数|插入播放记录")
			if inerr == nil {
				_, inerr = dbmap.Exec("update lives set Playnum=Playnum+1 where Id=?;", lg.Liveid)
				if inerr == nil {
					cccid, inerr := dbmap.SelectInt("select Curriculumsclasscentreid from curriculumclassroomchaptercentre where Id=?;", lg.Ccccid)
					if cccid > 0 {
						_, inerr = dbmap.Exec("update curriculumsclasscentre set PlaySumnum=PlaySumnum+1 where Id=?;", cccid)
						core.CheckErr(inerr, "liveDataAccess|UpdateLivePlayNum|更新课程章节的播放次数|更新播放次数")
					}
				} else {
					rd.Rcode = "2003"
					rd.Reason = "数据记录失败"
				}
			} else {
				rd.Rcode = "2003"
				rd.Reason = "数据记录失败"
			}
		} else {
			rd.Rcode = "1000"
			rd.Reason = "已经记录过"
		}
	} else {
		rd.Rcode = "2002"
		rd.Reason = "数据提交错误"
	}
	return rd
}

/*
更新课程章节的下载次数
*/
func UpdateLiveDownloadNum(lg viewmodel.QueryLivesWhere, dbmap *gorp.DbMap) (rd core.Returndata) {
	if lg.Ccccid > 0 && lg.PsUserid > 0 && lg.Liveid > 0 {
		ct, inerr := dbmap.SelectInt("select count(*) as Ct from livedownloadlog where LivesId=? and UsersId=?;", lg.Liveid, lg.PsUserid)
		core.CheckErr(inerr, "liveDataAccess|UpdateLiveDownloadNum|更新课程章节的下载次数|获取下载次数总和")
		if ct < 1 {
			dbmap.AddTableWithName(live.Livedownloadlog{}, "livedownloadlog").SetKeys(true, "Id")
			pr := live.Livedownloadlog{LivesId: lg.Liveid, UsersId: lg.PsUserid}
			inerr = dbmap.Insert(&pr)
			core.CheckErr(inerr, "liveDataAccess|UpdateLiveDownloadNum|更新课程章节的下载次数|添加下载记录")
			if inerr == nil {
				_, inerr = dbmap.Exec("update lives set Downloadnum=Downloadnum+1 where Id=?;", lg.Liveid)
				if inerr == nil {
					cccid, inerr := dbmap.SelectInt("select Curriculumsclasscentreid from curriculumclassroomchaptercentre where Id=?;", lg.Ccccid)
					if cccid > 0 {
						_, inerr = dbmap.Exec("update curriculumsclasscentre set DownloadSumnum=DownloadSumnum+1 where Id=?;", cccid)
						core.CheckErr(inerr, "liveDataAccess|UpdateLiveDownloadNum|更新课程章节的下载次数|更新下载次数")
					}
				} else {
					rd.Rcode = "2003"
					rd.Reason = "数据记录失败"
				}
			} else {
				rd.Rcode = "2003"
				rd.Reason = "数据记录失败"
			}
		} else {
			rd.Rcode = "1000"
			rd.Reason = "已经记录过"
		}
	} else {
		rd.Rcode = "2002"
		rd.Reason = "数据提交错误"
	}
	return rd
}

//判断条件获取
func GetQueryLiveWhere(lg viewmodel.QueryLivesWhere) (where string) {
	if lg.Islive > 0 {
		where = where + " and ccc.Islive=" + strconv.Itoa(lg.Islive)
	}
	if lg.Isondemand > 0 {
		where = where + " and ccc.Isondemand=" + strconv.Itoa(lg.Isondemand)
	}
	if lg.Classesid > 0 {
		where = where + " and ccc.Classesid=" + strconv.Itoa(lg.Classesid)
	}
	if lg.Ccccid > 0 {
		where = where + " and ccc.Curriculumclassroomchaptercentreid=" + strconv.Itoa(lg.Ccccid)
	}
	if lg.Subjectcode != "" {
		if lg.Isall == 1 {
			where = where + " and(cc.Subjectcode='" + lg.Subjectcode + "' or cc.Subjectcode like '%" + lg.Subjectcode + "%')"
		} else {
			where = where + " and(cc.Subjectcode='" + lg.Subjectcode + "')"
		}
	}
	if lg.Wherestring != "" {
		where = where + " and (cc.Curriculumname like '%" + lg.Wherestring + "%' or cts.Chaptername like '%" + lg.Wherestring + "%' or ls.Livetitile like '%" + lg.Wherestring + "%' or ls.Liveinfo like '%" + lg.Wherestring + "%')"
	}
	if lg.Curriculumsid > 0 {
		where = where + " and ccc.Curriculumsid=" + strconv.Itoa(lg.Curriculumsid)
	}
	if lg.Begindate != "" {
		lg.Begindate = core.Timeaction(lg.Begindate)
		where = where + " and cccc.Begindate>='" + lg.Begindate + "'"
	}
	if lg.Enddate != "" {
		lg.Enddate = core.Timeaction(lg.Enddate)
		where = where + " and cccc.Begindate<='" + lg.Enddate + "'"
	}
	if lg.Psstate == 0 {
		where = where + " and ps.State=0"
	}
	if lg.Psstate == 1 {
		where = where + " and ps.State=1"
	}
	if lg.PsUserid > 0 {
		where = where + " and ps.Usersid=" + strconv.Itoa(lg.PsUserid)
	}
	if lg.Trstate == 0 { //查询未播放
		where = where + " and tr.State=0"
	}
	if lg.Trstate == 1 { //查询正在上课的
		where = where + " and tr.State=1"
	}
	if lg.Trstate == 2 { //查询未播放
		where = where + " and tr.State=2"
	}
	if lg.Collegeids != "" { //学院
		where = where + " and Collegeid in(" + lg.Collegeids + ")"
	}
	if lg.Majorids != "" { //专业
		where = where + " and Majorid in(" + lg.Majorids + ")"
	}
	if lg.Teacherids != "" { //老师
		where = where + " and ccc.Usersid in(" + lg.Teacherids + ")"
	}
	if lg.Classesids != "" { //班级
		where = where + " and ccc.Classesid in(" + lg.Classesids + ")"
	}
	return where
}
