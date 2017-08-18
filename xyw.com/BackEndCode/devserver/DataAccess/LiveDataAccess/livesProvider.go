package liveDataAccess

import (
	"fmt"

	"gopkg.in/gorp.v1"

	"dev.project/BackEndCode/devserver/commons/xdebug"
	"dev.project/BackEndCode/devserver/model/live"
	"dev.project/BackEndCode/devserver/viewmodel/videoview"
)

// 查询课堂有关视频信息
func QueryLiveInfoByID(idCurriculumClassroomChapter int, dbmap *gorp.DbMap) (r live.Lives, ok bool) {
	sql := `SELECT TChapter.Chaptername AS Liveinfo, TChapter.Chaptername AS Livetitile, CCCC.Id AS Curriculumclassroomchaptercentreid
			FROM curriculumclassroomchaptercentre AS CCCC 
				JOIN curriculumsclasscentre AS CCC ON (CCCC.Curriculumsclasscentreid = CCC.Id)
				JOIN curriculums AS TCurr ON (CCC.Curriculumsid = TCurr.Id)
				JOIN chapters AS TChapter ON (CCCC.Chaptersid = TChapter.Id)
			WHERE CCCC.Id = ?;`

	err := dbmap.SelectOne(&r, sql, idCurriculumClassroomChapter)
	xdebug.DebugError(err)
	ok = (err == nil)

	return r, ok
}

// 查询视频信息列表
func QueryVideoInfos(condition *videoview.Request_VideoInfoView, dbmap *gorp.DbMap) (list []videoview.VideoInfoView, ok bool) {
	sqlcount := `SELECT COUNT(TLive.Id)
				 FROM lives AS TLive
					JOIN curriculumclassroomchaptercentre AS CCC ON (CCC.Id = TLive.Curriculumclassroomchaptercentreid)
					JOIN curriculumsclasscentre AS CC ON (CC.Id = CCC.Curriculumsclasscentreid)
					JOIN curriculums AS TCurr ON (TCurr.Id = CC.Curriculumsid)
					JOIN chapters AS TChap ON (TChap.Id = CCC.Chaptersid)
					JOIN users AS TUser ON (TUser.Id = CCC.Usersid) `
	sqlcount += condition.WhereCondition()
	if count, _ := dbmap.SelectInt(sqlcount); count > 0 {
		condition.PageInfo.RowTotal = int(count)
		sql := `
SELECT 	TLive.Id AS ID,
		TLive.Livepath1 AS VideoPath1,
		TLive.Livepath2 AS VideoPath2,
		TLive.Liveinfo AS VideoInfo,
		TLive.Livetitile AS VideoTitle,
		TLive.Whenlong AS VideoDuration,
		TLive.Begindate AS LiveBeginTime,
		TLive.Livestate AS LiveState,
		TLive.Livetype AS LiveType,
		TLive.Downloadnum AS DownloadNum,
		TLive.Playnum AS PlayNum,
		TLive.IsRelease AS IsRelease,
		TLive.Iscomment AS AllowComments,
		TLive.Isdownload AS AllowDownload,
		CCC.Islive AS AllowLive,			 
		TCurr.Id AS CurriculumID,
		TCurr.Curriculumname AS CurriculumName,
		TChap.Id AS ChapterID,
		TChap.Chaptername AS ChapterName,
		TUser.Id AS TeacherID,
		TUser.Nickname AS TeacherName,
		TLive.Curriculumclassroomchaptercentreid AS CurriculumClassroomChapterID
FROM lives AS TLive
		JOIN curriculumclassroomchaptercentre AS CCC ON (CCC.Id = TLive.Curriculumclassroomchaptercentreid)
		JOIN curriculumsclasscentre AS CC ON (CC.Id = CCC.Curriculumsclasscentreid)
		JOIN curriculums AS TCurr ON (TCurr.Id = CC.Curriculumsid)
		JOIN chapters AS TChap ON (TChap.Id = CCC.Chaptersid)
		JOIN users AS TUser ON (TUser.Id = CCC.Usersid) 
`
		sql += condition.WhereCondition() + condition.PageInfo.LimitString()
		_, err := dbmap.Select(&list, sql)
		xdebug.DebugError(err)
		ok = (err == nil)
	}

	return list, ok
}

// 查询视频详细
func QueryVideoDetials(vid int, dbmap *gorp.DbMap) (video videoview.VideoDetailView, ok bool) {
	defer xdebug.DoRecover()
	sql := `SELECT 	TLive.Id AS ID,
				TLive.Livetitile AS VideoTitle,
				TLive.Liveinfo AS VideoInfo,
				TLive.Coverimage AS CoverImage,
				TLive.Livestate AS LiveState,					
				TLive.Whenlong AS VideoDuration,
				TLive.Iscomment AS AllowComments,
				TLive.Isdownload AS AllowDownload,
				TLive.Livepath1 AS VodPath1,
				TLive.Livepath2 AS VodPath2,
				TLive.Recommendread AS RecommendReads,
				TLive.Ischeckcomment AS IsCheckComments,
				CCC.Begindate AS BeginTime,
				CCC.Enddate AS EndTime,
				CCC.Isondomian AS AllowVOD,
				TLive.Downloadnum AS DownloadNum,
				TLive.Playnum AS PlayNum,
				TLive.IsRelease AS IsRelease
			FROM lives AS TLive
					JOIN curriculumclassroomchaptercentre AS CCC ON (CCC.Id = TLive.Curriculumclassroomchaptercentreid)
					JOIN curriculumsclasscentre AS CC ON (CC.Id = CCC.Curriculumsclasscentreid)
					JOIN curriculums AS TCurr ON (TCurr.Id = CC.Curriculumsid)
					JOIN chapters AS TChap ON (TChap.Id = CCC.Chaptersid)
					JOIN users AS TUser ON (TUser.Id = CCC.Usersid)
			WHERE TLive.Id = ? `

	// fmt.Println(sql, vid)
	err := dbmap.SelectOne(&video, sql, vid)
	xdebug.HandleError(err)

	sql = `SELECT  TEn.Id AS ID,
					TEn.Enclosurename AS EnclosureName,
					TEn.Enclosuretype AS EnclosureType, 
					TEn.Enclosuresize AS EnclosureSize,
					TEn.EnclosureVirtualPath AS VirtualPath,
					TEn.Enclosurepath AS EnclosurePath,
					TEn.Createdate AS CreateDate,
					TEn.Enclosureicon AS EnclosuerIcon, 
					TEn.IsPublish AS IsPublish,
					TLive.Curriculumclassroomchaptercentreid AS CurriculumClassroomChapterID
			FROM lives AS TLive
				JOIN curriculumclassroomchaptercentre AS CCC ON (CCC.Id = TLive.Curriculumclassroomchaptercentreid)
				JOIN enclosure AS TEn ON(TEn.Curriculumclassroomchaptercentreid = TLive.Curriculumclassroomchaptercentreid) 
			WHERE TLive.Id = ?`
	// fmt.Println(sql, vid)
	var arr = []videoview.AttachmentView{}
	_, err = dbmap.Select(&arr, sql, vid)
	xdebug.HandleError(err)

	sql = `	SELECT 	TLive.Id AS ID,
					CCC.Begindate AS BeginTime,
					CCC.Enddate AS EndTime,
					TCurriculum.Id AS CurriculumID,
					TCurriculum.Curriculumname AS CurriculumName,
					TChapter.Id AS ChapterID,
					TChapter.Chaptername AS ChapterName,
					TUser.Id AS TeacherID,
					TUser.Nickname AS TeacherName,
					TClass.Id AS ClassID,
					TClass.Classesname AS ClassName,
					TRoom.Id AS ClassroomID,
					TRoom.Classroomsname AS ClassroomName,
					TBuild.Id AS BuildingID,
					TBuild.Buildingname AS BuildingName,
					TCampus.Id AS CampusID,
					TCampus.Campusname AS CampusName
			FROM lives AS TLive
					JOIN curriculumclassroomchaptercentre AS CCC ON (CCC.Id = TLive.Curriculumclassroomchaptercentreid)		
					JOIN curriculumsclasscentre AS CC ON (CC.Id = CCC.Curriculumsclasscentreid)
					JOIN curriculums AS TCurriculum ON (TCurriculum.Id = CC.Curriculumsid)
					JOIN chapters AS TChapter ON (TChapter.Id = CCC.Chaptersid)
					JOIN users AS TUser ON (TUser.Id = CCC.Usersid)
					JOIN classes AS TClass ON (TClass.Id = CC.Classesid)
					JOIN teachingrecord AS TRecord ON (TRecord.Curriculumclassroomchaptercentreid = CCC.Id)
					JOIN classrooms AS TRoom ON (TRoom.Id = TRecord.Classroomid)
					JOIN floors AS TFloor ON (TFloor.Id = TRoom.Floorsid)
					JOIN building AS TBuild ON (TBuild.Id = TFloor.Buildingid)
					JOIN campus AS TCampus ON (TCampus.Id = TBuild.Campusid)
			WHERE TLive.Id = ?`
	// fmt.Println(sql, vid)
	var lesson videoview.LessonView
	err = dbmap.SelectOne(&lesson, sql, vid)
	xdebug.HandleError(err)

	ok = true
	video.AttachmentList = arr
	video.LessonInfo = lesson
	return video, ok
}

// 更新视频详细
func UpdateVideoDetails(video videoview.VideoDetailView, dbmap *gorp.DbMap) (ok bool) {
	sql := `UPDATE lives SET Livestate = ?,Iscomment = ?,Ischeckcomment = ?,Isdownload = ?,Downloadnum = ?,Playnum = ?,IsRelease = ?,
				Livetitile = ?, Liveinfo = ?, Coverimage = ?, Recommendread = ? 
			WHERE (Id = ?);`
	fmt.Println(sql)
	fmt.Println(video)
	defer xdebug.DoRecover()
	_, err := dbmap.Exec(sql, video.LiveState, video.AllowComments, video.IsCheckComments,
		video.AllowDownload, video.DownloadNum, video.PlayNum, video.IsRelease,
		video.VideoTitle, video.VideoInfo, video.CoverImage, video.RecommendReads, video.ID)
	xdebug.HandleError(err)

	ok = true
	return ok
}

// 删除视频
func DeleteVideo(vid int, dbmap *gorp.DbMap) (ok bool) {
	sql := "DELETE FROM lives WHERE Id = ?"
	if _, err := dbmap.Exec(sql, vid); err == nil {
		ok = true
	}

	return ok
}

// 删除附件
func DeleteAttachment(rowid int, dbmap *gorp.DbMap) (ok bool) {
	sql := "DELETE FROM enclosure WHERE Id = ?"
	if _, err := dbmap.Exec(sql, rowid); err == nil {
		ok = true
	}

	return ok
}

func foo() {
	fmt.Print("")
}
