package basicset

import (
	"dev.project/BackEndCode/devserver/model/core"
	"fmt"
)

type Campus struct { //
	Id         int
	Campusname string
	Campusicon string
	Campuscode string
	Campusnums int // 校区内的总人数(反写值)
}

type College struct { //
	Id          int
	Campusid    int
	Collegename string
	Collegeicon string
	Collegecode string
	Collegenum  int // 学院内的总人数(反写值)
}
type Major struct { //
	Id        int
	Collegeid int
	Majorname string
	Majoricon string
	Majorcode string
	Majornum  int
}
type Classes struct { //
	Id             int
	Majorid        int
	Classesname    string
	Classescode    string
	Classesnum     int // 班级内的总人数(反写值)
	Classesicon    string
	Classstate     int
	Enrollmentyear int
}
type Building struct { //
	Id               int
	Campusid         int
	Buildingname     string
	Buildingicon     string
	Buildingcode     string
	Floorsnumber     int
	Classroomsnumber int
}
type Floors struct { //
	Id              int
	Buildingid      int
	Floorname       string
	Floorscode      string
	FloorsImage     string
	Classroomnumber int
	Maxy            float32
	Miny            float32
	Maxx            float32
	Minx            float32
	Sumnumber       int
}
type Classrooms struct {
	Id                int
	Floorsid          int
	Classroomsname    string
	Classroomicon     string
	Seatsnumbers      int
	Sumnumbers        int
	Classroomstype    string
	Classroomstate    int
	Collectionnumbers int
	Maxy              float32
	Miny              float32
	Maxx              float32
	Minx              float32
	Notes             string
	Classroomscode    string
}

func init() {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(Campus{}, "campus").SetKeys(true, "Id")
	dbmap.AddTableWithName(College{}, "college").SetKeys(true, "Id")
	dbmap.AddTableWithName(Major{}, "major").SetKeys(true, "Id")
	dbmap.AddTableWithName(Classes{}, "classes").SetKeys(true, "Id")
	dbmap.AddTableWithName(Building{}, "building").SetKeys(true, "Id")
	dbmap.AddTableWithName(Floors{}, "floors").SetKeys(true, "Id")
	dbmap.AddTableWithName(Classrooms{}, "classrooms").SetKeys(true, "Id")
	err := dbmap.CreateTablesIfNotExists()
	core.CheckErr(err, "Create tables failed")
	fmt.Println("campus模块开始初始化")
}
