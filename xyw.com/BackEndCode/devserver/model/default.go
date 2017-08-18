package model

import (
	"fmt"
	_ "dev.project/BackEndCode/devserver/model/actiondata"
	_ "dev.project/BackEndCode/devserver/model/basicset"
	_ "dev.project/BackEndCode/devserver/model/core"
	_ "dev.project/BackEndCode/devserver/model/curriculum"
	_ "dev.project/BackEndCode/devserver/model/live"
	_ "dev.project/BackEndCode/devserver/model/systemmodel"
	_ "dev.project/BackEndCode/devserver/model/users"
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
