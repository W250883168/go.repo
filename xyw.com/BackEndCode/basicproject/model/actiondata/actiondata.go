package actiondata

import (
	"dev.project/BackEndCode/devserver/model/core"
	"fmt"
)

type Attentionrecord struct {
	Id            int
	Usersid       int
	Curriculumsid int
	Createdate    string
	Classesid     int
	State         int
}

type Classroomcollection struct {
	Id          int
	Usersid     int
	Classroomid int
	Createdate  string
}
type Classroomdetailsnums struct {
	Id          int
	Classroomid int
	Loginuser   string
	Mac         string
	Createdate  string
	Positionnum int
	Xy          string
	X           float64
	Y           float64
	Ap          string
}
type Teachingrecord struct {
	Id                                 int
	Classroomid                        int
	Curriculumclassroomchaptercentreid int
	State                              int
}
type Pointtos struct {
	Id                                 int
	Curriculumclassroomchaptercentreid int
	Usersid                            int
	State                              int
	Ismodify                           int
}

func init() {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(Attentionrecord{}, "attentionrecord").SetKeys(true, "Id")
	dbmap.AddTableWithName(Classroomcollection{}, "classroomcollection").SetKeys(true, "Id")
	dbmap.AddTableWithName(Classroomdetailsnums{}, "classroomdetailsnums").SetKeys(true, "Id")
	dbmap.AddTableWithName(Teachingrecord{}, "teachingrecord").SetKeys(true, "Id")
	dbmap.AddTableWithName(Pointtos{}, "pointtos").SetKeys(true, "Id")
	err := dbmap.CreateTablesIfNotExists()
	core.CheckErr(err, "Create tables failed")
	fmt.Println("attentionrecord模块开始初始化")
}
