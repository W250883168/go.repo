package systemmodelControllers

import (
	"encoding/json"

	"basicproject/DataAccess/basicsetDataAccess"
	"basicproject/DataAccess/usersDataAccess"
	"basicproject/model/basicset"
	"basicproject/viewmodel"
	core "xutils/xcore"

	"github.com/gin-gonic/gin"
)

//查询校区数据
func GetCampusList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetCampusList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetCampusList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetCampusList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryCampusPG(obj, page, dbmap)
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
func AddCampus(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Campus
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddCampus|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddCampus", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueCampuscode(obj.Campuscode, obj.Campusname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddCampus(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddCampus|添加校区数据失败：")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败:" + adderr.Error()
				}
			}
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

//修改校区模块
func CampusUpdate(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Campus
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|CampusUpdate|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "CampusUpdate", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueCampuscode(obj.Campuscode, obj.Campusname, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:不能修改校区代码"
			} else {
				adderr := basicsetDataAccess.UpdateCampus(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|CampusUpdate|添加校区数据失败：")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
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

//删除校区
func DelCampus(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Campus
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelCampus|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelCampus", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteCampus(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DeleteCampus|删除校区数据失败：")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
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

//查询楼栋
func GetBuildingList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetBuildingList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetBuildingList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetBuildingList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryBuildingPG(obj, page, dbmap)
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

//添加楼栋
func AddBuilding(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Building
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddBuilding|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddBuilding", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueBuildingcode(obj.Buildingcode, obj.Buildingname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddBuilding(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddBuilding|添加楼栋数据失败：")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败"
				}
			}
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

//修改楼栋
func ChangeBuilding(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Building
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeBuilding|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeBuilding", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueBuildingcode(obj.Buildingcode, obj.Buildingname, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:不能修改楼栋代码"
			} else {
				adderr := basicsetDataAccess.UpdateBuilding(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|ChangeBuilding|修改楼栋数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
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

//删除楼栋
func DelBuilding(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Building
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelBuilding|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelBuilding", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteBuilding(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelBuilding|删除楼栋数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
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

//获取所有的楼层列表
func GetFloorsList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetFloorsList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetFloorsList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetFloorsList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryFloorsPG(obj, page, dbmap)
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

//添加楼层
func AddFloors(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Floors
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddFloors|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddFloors", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueFloorscode(obj.Floorscode, obj.Floorname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddFloors(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddFloors|添加楼层数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改楼层
func ChangeFloors(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Floors
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeFloors|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeFloors", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueFloorscode(obj.Floorscode, obj.Floorname, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:此编码不能再修改"
			} else {
				adderr := basicsetDataAccess.UpdateFloors(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|ChangeFloors|修改楼层数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除楼层
func DelFloors(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Floors
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelFloors|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelFloors", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteFloors(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelFloors|删除楼层数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取所有的教室列表
func GetClassroomsList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetClassroomsList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetClassroomsList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetClassroomsList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryClassroomsPG(obj, page, dbmap)
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

//添加教室
func AddClassrooms(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Classrooms
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddClassrooms|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddClassrooms", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueClassroomscode(obj.Classroomscode, obj.Classroomsname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddClassrooms(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddClassrooms|添加教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改教室
func ChangeClassrooms(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Classrooms
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeClassrooms|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeClassrooms", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueClassroomscode(obj.Classroomscode, obj.Classroomsname, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:此编码不允许修改"
			} else {
				adderr := basicsetDataAccess.UpdateClassrooms(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|ChangeClassrooms|修改教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除教室
func DelClassrooms(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Classrooms
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelClassrooms|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelClassrooms", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteClassrooms(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelClassrooms|删除教室数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取学院数据列表
func GetCollegeList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetCollegeList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetCollegeList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetCollegeList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryCollegePG(obj, page, dbmap)
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

//添加学院
func AddCollege(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.College
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddCollege|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddCollege", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueCollegecode(obj.Collegecode, obj.Collegename, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddCollege(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddCollege|添加教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改学院
func ChangeCollege(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.College
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeCollege|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeCollege", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueCollegecode(obj.Collegecode, obj.Collegename, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:此编码不允许修改"
			} else {
				adderr := basicsetDataAccess.UpdateCollege(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|ChangeCollege|修改教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除学院
func DelCollege(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.College
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelCollege|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelCollege", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteCollege(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelCollege|删除教室数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取专业数据列表
func GetMajorList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetMajorList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetMajorList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetMajorList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryMajorPG(obj, page, dbmap)
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

//添加专业
func AddMajor(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Major
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddMajor|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddMajor", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueMajorcode(obj.Majorcode, obj.Majorname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1002"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddMajor(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddMajor|添加教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改专业
func ChangeMajor(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Major
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeMajor|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeMajor", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueMajorcode(obj.Majorcode, obj.Majorname, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:此编码不允许修改"
			} else {
				adderr := basicsetDataAccess.UpdateMajor(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|ChangeMajor|修改教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}

	c.JSON(200, rd)
}

//删除专业
func DelMajor(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Major
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelMajor|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelMajor", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteMajor(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelMajor|删除教室数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取班级数据列表
func GetClassesList(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryBasicsetWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetClassesList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetClassesList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetClassesList", dbmap)
		if rd.Rcode == "1000" {
			rd = basicsetDataAccess.QueryClassesPG(obj, page, dbmap)
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

//添加班级
func AddClasses(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Classes
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddClasses|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddClasses", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueClassescode(obj.Classescode, obj.Classesname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1002"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := basicsetDataAccess.AddClasses(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddClasses|添加教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Result = obj.Id
				} else {
					rd.Rcode = "1001"
					rd.Reason = "添加失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改班级
func ChangeClasses(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Classes
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeClasses|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeClasses", dbmap)
		if rd.Rcode == "1000" {
			isUnique := basicsetDataAccess.QueryUniqueClassescode(obj.Classescode, obj.Classesname, dbmap) //检查唯一约束
			if isUnique != 1 {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:此编码不允许修改"
			} else {
				adderr := basicsetDataAccess.UpdateClasses(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|ChangeClasses|修改教室数据失败")
				if adderr == nil {
					rd.Rcode = "1000"
					rd.Reason = "修改成功"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "修改失败:" + adderr.Error()
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除班级
func DelClasses(c *gin.Context) {
	var rd core.Returndata
	var obj basicset.Classes
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelClasses|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelClasses", dbmap)
		if rd.Rcode == "1000" {
			adderr := basicsetDataAccess.DeleteClasses(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelClasses|删除教室数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Reason = "修改成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "修改失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}
