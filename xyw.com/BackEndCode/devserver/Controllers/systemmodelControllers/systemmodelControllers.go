package systemmodelControllers

import (
	"encoding/json"
	//	"fmt"
	"dev.project/BackEndCode/devserver/DataAccess/systemmodelDataAccess"
	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/systemmodel"
	//	"dev.project/BackEndCode/devserver/viewmodel"
	//	"strconv"
	//	"strings"

	"github.com/gin-gonic/gin"
)

//获取所有的模块数据
//参数[PageIndex,PageSize,ModelName(like),Superiormoduleid(上级模块的Id),systemmoduleId(具体模块的Id)]
func GetSystemModelList(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodule //viewmodel.GetStudentsinfo
	var page core.PageData           //获取分页设置数据
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetSystemModelList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetSystemModelList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetSystemModelList", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.QuerySystemModelList(page, obj, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加系统模块
func AddSystemModel(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodule
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddSystemModel|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddSystemModel", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.InsertSystemModel(obj, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改系统模块
func ChangeSystemModel(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodule
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeSystemModel|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "UpdateSystemModel", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.UpdateSystemModel(obj, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除系统模块
func DelSystemModel(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodule
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelSystemModel|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DeleteSystemModel", dbmap)
		if rd.Rcode == "1000" {
			//第一步删除模块和功能中间相关数据
			rd = systemmodelDataAccess.DeleteSystemModel(obj, dbmap) //第二步删除模块数据
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取所有的功能列表
//参数[PageIndex,PageSize,ModelName(like),Superiormoduleid(上级模块的Id),systemmoduleId(具体模块的Id)]
func GetSystemModelFunctionsList(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodulefunctions //viewmodel.GetStudentsinfo
	var page core.PageData                    //获取分页设置数据

	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetSystemModelFunctionsList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetSystemModelFunctionsList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetSystemModelFunctionsList", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.QuerySystemModelFunctionsList(page, obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加系统模块功能
func AddSystemModelFunctions(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodulefunctions
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddSystemModelFunctions|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddSystemModelFunctions", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.InsertSystemModelFunctions(obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改系统模块功能
func ChangeSystemModelFunctions(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodulefunctions
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")

	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeSystemModelFunctions|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "UpdateSystemModelFunctions", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.UpdateSystemModelFunctions(obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除系统模块功能
func DelSystemModelFunctions(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Systemmodulefunctions
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelSystemModelFunctions|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DeleteSystemModelFunctions", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.DeleteSystemModelFunctions(obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取所有的功能列表
func GetRolesList(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Roles //viewmodel.GetStudentsinfo
	var page core.PageData    //获取分页设置数据
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetRolesList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetRolesList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetRolesList", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.QueryRolesList(page, obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加系统角色
func AddRoles(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Roles
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddRoles|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddRoles", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.InsertRoles(obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改系统角色
func ChangeRoles(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Roles
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeRoles|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "UpdateRoles", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.UpdateRoles(obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除系统角色
func DelRoles(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Roles
	//data, _ := ioutil.ReadAll(c.Request.Body)
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelRoles|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DeleteRoles", dbmap)
		if rd.Rcode == "1000" {
			rd = systemmodelDataAccess.DeleteRoles(obj, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取所有的功能列表
func GetSetSystemConfig(c *gin.Context) {
	var rd core.Returndata
	var obj systemmodel.Roles
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetSetSystemConfig|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetSetSystemConfig", dbmap)
		if rd.Rcode == "1000" {
			if obj.Id > 0 {
				rd = systemmodelDataAccess.QuerySetSystemConfig(obj.Id, dbmap)
			} else {
				rd.Rcode = "1002"
				rd.Reason = "数据提交格式错误"
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加系统角色
func SaveSystemConfig(c *gin.Context) {
	var rd core.Returndata
	//	var obj systemmodel.Roles
	//	var obj1 []systemmodel.Rolemodulefunctioncenter
	//	var obj2 []systemmodel.Rolemodulecenter
	type PostSmlist struct {
		Rmfclist    []systemmodel.Rolemodulefunctioncenter
		Rmclist     []systemmodel.Rolemodulecenter
		DelRmfclist []systemmodel.Rolemodulefunctioncenter
		DelRmclist  []systemmodel.Rolemodulecenter
	}
	var podata PostSmlist
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &podata)
	core.CheckErr(errs1, "systemmodelControllers|SaveSystemConfig|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "SaveSystemConfig", dbmap)
		if rd.Rcode == "1000" {
			rd.Rcode = "1003"
			if len(podata.DelRmclist) > 0 || len(podata.DelRmfclist) > 0 {
				rd = systemmodelDataAccess.DelSystemConfig(podata.DelRmfclist, podata.DelRmclist, dbmap)
			}
			if len(podata.Rmclist) > 0 || len(podata.Rmfclist) > 0 {
				rd = systemmodelDataAccess.SaveSystemConfig(podata.Rmfclist, podata.Rmclist, dbmap)
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}
