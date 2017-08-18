package model

import (
	_ "basicproject/model/actiondata"
	_ "basicproject/model/basicset"
	_ "basicproject/model/curriculum"
	_ "basicproject/model/live"
	_ "basicproject/model/systemmodel"
	_ "basicproject/model/users"
	"fmt"
	_ "xutils/xcore"
)

func init() {
	fmt.Println("开始初始化相关数据表")
	/*
		dbmap := core.InitDb()
		dbmap.AddTableWithName(Systemmodule{}, "systemmodule").SetKeys(true, "Id")
		dbmap.AddTableWithName(Systemmodulefunctions{}, "Systemmodulefunctions").SetKeys(true, "Id")
		dbmap.AddTableWithName(Rolemodulecenter{}, "rolemodulecenter").SetKeys(true, "Id")
		dbmap.AddTableWithName(Rolemodulefunctioncenter{}, "rolemodulefunctioncenter").SetKeys(true, "Id")
		dbmap.AddTableWithName(Roles{}, "roles").SetKeys(true, "Id")
		err := dbmap.CreateTablesIfNotExists()
		core.CheckErr(err, "Create tables failed")
	*/
}
