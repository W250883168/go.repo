package live

import (
	"dev.project/BackEndCode/devserver/model/core"
	"fmt"
)

type Lives struct { //
	Id                                 int
	Curriculumclassroomchaptercentreid int
	Liveinfo                           string
	Livetitile                         string
	Livepath1                          string //播放1路视频
	Livepath2                          string //播放2路视频
	Trailerpath                        string
	Coverimage                         string
	Livestate                          int
	Livetype                           int
	Begindate                          string
	Whenlong                           int
	Recommendread                      string
	Iscomment                          int
	Ischeckcomment                     int
	Isdownload                         int
	Downloadnum                        int //下载次数
	Playnum                            int //播放次数
	IsRelease                          int
}

type Playrecord struct { //
	Id         int
	Usersid    int
	Livesid    int
	Createdate string
}

type Comments struct { //
	Id             int
	Livesid        int
	Commentdetails string
	Ischeckok      int
	createdate     string
	Usersid        int
	Goodnum        int
}

type Goodrecord struct { //
	Id         int
	Usersid    int
	Createdate string
	Commentid  int
}

type Livedownloadlog struct { //下载记录表
	Id         int
	CreateDate string
	UsersId    int
	LivesId    int
}

func init() {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(Lives{}, "lives").SetKeys(true, "Id")
	dbmap.AddTableWithName(Playrecord{}, "playrecord").SetKeys(true, "Id")
	dbmap.AddTableWithName(Comments{}, "comments").SetKeys(true, "Id")
	dbmap.AddTableWithName(Goodrecord{}, "goodrecord").SetKeys(true, "Id")
	dbmap.AddTableWithName(Livedownloadlog{}, "livedownloadlog").SetKeys(true, "Id")
	err := dbmap.CreateTablesIfNotExists()
	core.CheckErr(err, "Create tables failed")
	fmt.Println("lives模块开始初始化")
}
