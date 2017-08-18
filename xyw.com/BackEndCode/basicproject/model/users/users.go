package users

import (
	core "xutils/xcore"
	"fmt"
)

type Users struct { //
	Id           int
	Loginuser    string
	Loginpwd     string
	Rolesid      int
	Truename     string
	Nickname     string
	Userheadimg  string
	Userphone    string
	Userstate    int
	Usersex      int
	Usermac      string
	Birthday     string
	ThirdPartyId string
	Os           string
}

type Usermessage struct { //
	Id              int
	Title           string
	Details         string
	State           int
	Usersid         int
	Createdate      string
	Megtype         string
	Readdate        string
	GoUrl           string
	GoParameter     string
	MessageImg      string
	MessageProfiles string
}
type Students struct { //
	Id               int
	Enrollmentyear   int
	Homeaddress      string
	Nowaddress       string
	Classesid        int //Classesid
	Infostate        int
	Currentstate     int
	Attendance       float32 //出勤率
	Needcoursenum    int     //应上的课程数
	Alreadycoursenum int     //已经上过的课程数
	Absenteeism      int     //旷课数
}
type Teacher struct { //教师表
	Id        int
	Collegeid int
	Majorid   int
}

type LoginLog struct {
	Id         int
	LoginOS    string
	LoginUsers string
	LoginDate  string
	LoginState int
	IP         string
}

func init() {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(Users{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Usermessage{}, "usermessage").SetKeys(true, "Id")
	dbmap.AddTableWithName(Students{}, "students")
	dbmap.AddTableWithName(Teacher{}, "teacher")
	dbmap.AddTableWithName(LoginLog{}, "loginlog").SetKeys(true, "Id")
	err := dbmap.CreateTablesIfNotExists()
	core.CheckErr(err, "Create tables failed")
	fmt.Println("users模块开始初始化")
}
