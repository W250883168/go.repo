package systemmodel

import (
	"dev.project/BackEndCode/devserver/model/core"
	"fmt"
)

type Systemmodule struct { //系统模块表
	Id                    int //
	Modulename            string
	Modulecode            string
	Moduleicon            string
	Moduleurl             string
	Moduleattribute       string
	Superiormoduleid      int
	ModuleIndex           int
	Moduledisplayname     string //模块显示名称
	Moduledisplayterminal string //模块适用终端
}

type Systemmodulefunctions struct { //系统模块功能表
	Id                 int
	Systemmoduleid     int
	Functionname       string
	Functionicon       string
	Functioncode       string
	Functionsurl       string
	Functionsattribute string
	FunctionDescribe   string
}
type Rolemodulecenter struct { //角色模块中间表
	Id             int
	Rolesid        int
	Systemmoduleid int
	State          int
}
type Rolemodulefunctioncenter struct { //角色模块功能中间表
	Id                      int
	Systemmodulefunctionsid int
	Rolemodulecenterid      int
}
type Roles struct { //角色表
	Id        int
	Rolesname string
}
type Roleuserscenter struct { //角色模块中间表
	Roleid  int
	Usersid int
	State   int
}

func init() {
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	dbmap.AddTableWithName(Systemmodule{}, "systemmodule").SetKeys(true, "Id")
	dbmap.AddTableWithName(Systemmodulefunctions{}, "Systemmodulefunctions").SetKeys(true, "Id")
	dbmap.AddTableWithName(Rolemodulecenter{}, "rolemodulecenter").SetKeys(true, "Id")
	dbmap.AddTableWithName(Rolemodulefunctioncenter{}, "rolemodulefunctioncenter").SetKeys(true, "Id")
	dbmap.AddTableWithName(Roles{}, "roles").SetKeys(true, "Id")
	dbmap.AddTableWithName(Roleuserscenter{}, "roleuserscenter")
	err := dbmap.CreateTablesIfNotExists()
	core.CheckErr(err, "Create tables failed")
	fmt.Println("systemmodule模块开始初始化")
}
