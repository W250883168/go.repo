package curriculum

import (
	"log"
	"runtime"

	"xutils/xdebug"

	"dev.project/BackEndCode/devserver/model/core"
)

type Curriculums struct { //课程表
	Id                 int
	Curriculumname     string
	Curriculumicon     string
	Curriculumnature   string
	Curriculumstype    string
	Curriculumsdetails string
	Chaptercount       int
	Averageclassrate   float32
	Subjectcode        string
	Createdate         string
}

type Chapters struct { //课程章节表
	Id             int
	Curriculumsid  int
	Chaptername    string
	Chaptericon    string
	Chapterdetails string
	Createdate     string
	ChaptersIndex  int
}
type Subjectclass struct { //
	Id                  int
	Subjectcode         string
	Subjectname         string
	Superiorsubjectcode string
}
type Curriculumsclasscentre struct { //
	Id               int
	Curriculumsid    int
	Classesid        int
	Usersid          int //	TeacherID        int
	Averageclassrate float32
	Createdate       string
	Newchapter       int
	Newlivechapter   int
	Islive           int
	Isondemand       int
	Whenlongcount    int
	PlaySumnum       int //播放次数
	DownloadSumnum   int //下载次数
	FollowSumnum     int //关注总人数
}
type Curriculumclassroomchaptercentre struct { //
	Id                       int
	Curriculumsclasscentreid int
	Chaptersid               int
	Plannumber               int
	Actualnumber             int
	Toclassrate              float32
	Usersid                  int
	Createdate               string
	Begindate                string
	Enddate                  string
	Islive                   int
	Isondomian               int
	Whenlong                 int
}

type Enclosure struct { //课程章节附件表
	Id                                 int    //文件Id
	Curriculumclassroomchaptercentreid int    //
	Enclosurename                      string //文件名称
	Enclosuretype                      string //文件类型
	Enclosuresize                      int    //文件大小以kb为单位
	EnclosureVirtualPath               string //文件的虚拟路径
	Enclosurepath                      string //文件的实际路径
	Createdate                         string //文件的创建时间
	Enclosureicon                      string //文件的封面图片
	IsPublish                          int    //课程资源是否发布
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(Curriculums{}, "curriculums").SetKeys(true, "Id")
	dbmap.AddTableWithName(Chapters{}, "chapters").SetKeys(true, "Id")
	dbmap.AddTableWithName(Subjectclass{}, "subjectclass").SetKeys(true, "Id")
	dbmap.AddTableWithName(Curriculumsclasscentre{}, "curriculumsclasscentre").SetKeys(true, "Id")
	dbmap.AddTableWithName(Curriculumclassroomchaptercentre{}, "curriculumclassroomchaptercentre").SetKeys(true, "Id")
	dbmap.AddTableWithName(Enclosure{}, "enclosure").SetKeys(true, "Id")
	err := dbmap.CreateTablesIfNotExists()
	xdebug.LogError(err)
}
