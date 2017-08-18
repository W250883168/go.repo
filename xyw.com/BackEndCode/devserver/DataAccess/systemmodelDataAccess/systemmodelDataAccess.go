package systemmodelDataAccess

import (
	"fmt"
	"strconv"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/systemmodel"
	"dev.project/BackEndCode/devserver/viewmodel"

	"gopkg.in/gorp.v1"
)

func GetQueryWhereStr(obj systemmodel.Systemmodule) (wherestr string) {
	//ModelName(like),Superiormoduleid(上级模块的Id),Id(具体模块的Id)]
	if obj.Id > 0 {
		wherestr = wherestr + " and sm.Id=" + strconv.Itoa(obj.Id)
	}
	if obj.Modulename != "" {
		wherestr = wherestr + " and(sm.Modulename='" + obj.Modulename + "' or sm.Modulename like '%" + obj.Modulename + "%' or sm.Moduledisplayname like '%" + obj.Modulename + "%')"
		wherestr = wherestr + " or(sm.Modulecode='" + obj.Modulename + "' or sm.Modulecode like '%" + obj.Modulename + "%' or sm.Moduledisplayname like '%" + obj.Modulename + "%')"
	}
	if obj.Modulecode != "" {
		wherestr = wherestr + " and(sm.Modulecode='" + obj.Modulecode + "' or sm.Modulecode like '%" + obj.Modulecode + "%')"
	}
	if obj.Superiormoduleid > 0 {
		wherestr = wherestr + " and sm.Superiormoduleid=" + strconv.Itoa(obj.Superiormoduleid)
	}
	return wherestr
}
func GetQueryWhereStrtwo(obj systemmodel.Systemmodulefunctions) (wherestr string) {
	//ModelName(like),Superiormoduleid(上级模块的Id),Id(具体模块的Id)]
	if obj.Id > 0 {
		wherestr = wherestr + " and smf.Id=" + strconv.Itoa(obj.Id)
	}
	if obj.Functionname != "" {
		wherestr = wherestr + " and(smf.Functionname='" + obj.Functionname + "' or smf.Functionname like '%" + obj.Functionname + "%')"
		wherestr = wherestr + " or(smf.Functioncode='" + obj.Functionname + "' or smf.Functioncode like '%" + obj.Functionname + "%')"
	}
	if obj.Functioncode != "" {
		wherestr = wherestr + " and(smf.Functioncode='" + obj.Functioncode + "' or smf.Functioncode like '%" + obj.Functioncode + "%')"
	}
	if obj.Systemmoduleid > 0 {
		wherestr = wherestr + " and smf.Systemmoduleid=" + strconv.Itoa(obj.Systemmoduleid)
	}
	return wherestr
}

//获取所有的模块数据
func QuerySystemModelList(pg core.PageData, obj systemmodel.Systemmodule, dbmap *gorp.DbMap) (rd core.Returndata) {
	sqlcount := "select count(*) from systemmodule as sm where 1=1" + GetQueryWhereStr(obj)
	fmt.Println(sqlcount)
	count, err := dbmap.SelectInt(sqlcount)
	if err == nil {
		pg.PageCount = int(count)
		sql := "select * from systemmodule as sm where 1=1"
		sql = sql + GetQueryWhereStr(obj) + " order by ModuleIndex " + core.GetLimitString(pg)
		var list []systemmodel.Systemmodule
		fmt.Println(sql)
		_, sserr1 := dbmap.Select(&list, sql)
		core.CheckErr(sserr1, "systemmodelDataAccess|QuerySystemModelList|获取所有的模块数据")
		if sserr1 == nil {
			pg.PageData = list
			rd.Result = pg
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1003"
			rd.Reason = "数据读取错误"
		}
	} else {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误"
	}
	return rd
}
func InsertSystemModel(obj systemmodel.Systemmodule, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Systemmodule{}, "systemmodule").SetKeys(true, "Id")
	count, _ := dbmap.SelectInt("select count(*) from systemmodule where (Modulecode='" + obj.Modulecode + "')")
	if count > 0 {
		rd.Rcode = "1004"
		rd.Reason = "数据提交错误，此模块代码已被占用"
	} else {
		inerr := dbmap.Insert(&obj)
		core.CheckErr(inerr, "systemmodelDataAccess|InsertSystemModel|添加新的模块")
		if inerr == nil {
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1003"
			rd.Reason = inerr.Error()
		}
	}
	return rd
}
func DeleteSystemModel(obj systemmodel.Systemmodule, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Systemmodule{}, "systemmodule").SetKeys(true, "Id")
	delsql1 := "delete from rolemodulefunctioncenter where Rolemodulecenterid in(select Id from rolemodulecenter where Systemmoduleid=?);"
	delsql2 := "delete from rolemodulecenter where Systemmoduleid=?;"
	_, inerr := dbmap.Exec(delsql1, obj.Id)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteSystemModel|删除角色模块功能中间表数据失败")
	_, inerr = dbmap.Exec(delsql2, obj.Id)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteSystemModel|删除角色模块中间表数据失败")
	_, inerr = dbmap.Delete(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteSystemModel|删除模块失败")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	return rd
}
func UpdateSystemModel(obj systemmodel.Systemmodule, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Systemmodule{}, "systemmodule").SetKeys(true, "Id")
	_, inerr := dbmap.Update(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|UpdateSystemModel|修改模块")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	return rd
}

func InsertSystemModelFunctions(obj systemmodel.Systemmodulefunctions, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Systemmodulefunctions{}, "systemmodulefunctions").SetKeys(true, "Id")
	//	count, _ := dbmap.SelectInt("select count(*) from systemmodulefunctions where (Functioncode='" + obj.Functioncode + "')")
	//	if count > 0 {
	//		rd.Rcode = "1004"
	//		rd.Reason = "数据提交错误，此模块功能代码已被占用"
	//	} else {
	inerr := dbmap.Insert(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|InsertSystemModelFunctions|添加模块功能")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	//	}
	return rd
}
func DeleteSystemModelFunctions(obj systemmodel.Systemmodulefunctions, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Systemmodulefunctions{}, "systemmodulefunctions").SetKeys(true, "Id")
	delsql1 := "delete from rolemodulefunctioncenter where Systemmodulefunctionsid=?;"
	_, inerr := dbmap.Exec(delsql1, obj.Id)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteSystemModelFunctions|删除模块功能中间表数据失败")
	_, inerr = dbmap.Delete(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteSystemModelFunctions|删除模块功能")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	return rd
}
func UpdateSystemModelFunctions(obj systemmodel.Systemmodulefunctions, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Systemmodulefunctions{}, "systemmodulefunctions").SetKeys(true, "Id")
	_, inerr := dbmap.Update(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|UpdateSystemModelFunctions|修改模块功能")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	return rd
}

//获取所有的模块下的功能数据
func QuerySystemModelFunctionsList(pg core.PageData, obj systemmodel.Systemmodulefunctions, dbmap *gorp.DbMap) (rd core.Returndata) {
	sqlcount := "select count(*) from systemmodulefunctions as smf where 1=1" + GetQueryWhereStrtwo(obj)
	count, err := dbmap.SelectInt(sqlcount)
	core.CheckErr(err, "systemmodelDataAccess|QuerySystemModelFunctionsList|获取所有的模块下的功能数据|获取所有模块数据总算")
	if err == nil {
		pg.PageCount = int(count)
		// sql := "select * from systemmodulefunctions as smf where 1=1"
		sql := `SELECT 	smf.Id, smf.Systemmoduleid, 
						IFNULL(smf.Functionname, '') AS Functionname,
						IFNULL(smf.Functionicon, '')  AS Functionicon,
						IFNULL(smf.Functioncode, '')  AS Functioncode,
						IFNULL(smf.Functionsurl, '')  AS Functionsurl,
						IFNULL(smf.Functionsattribute,  '')  AS Functionsattribute,
						IFNULL(smf.FunctionDescribe, '')  AS FunctionDescribe,
						IFNULL((SELECT sm.Modulename FROM systemmodule AS sm WHERE sm.Id = smf.Systemmoduleid),'') AS Systemmodulename
				FROM 	systemmodulefunctions AS smf where 1=1`
		sql = sql + GetQueryWhereStrtwo(obj) + " order by Systemmoduleid " + core.GetLimitString(pg)
		var list []viewmodel.SystemModuleFunctionView
		_, sserr1 := dbmap.Select(&list, sql)
		core.CheckErr(sserr1, "systemmodelDataAccess|QuerySystemModelFunctionsList|获取所有的模块下的功能数据|查询相关数据")
		if sserr1 == nil {
			pg.PageData = list
			rd.Result = pg
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1003"
			rd.Reason = "数据读取错误" + sserr1.Error()
		}
	} else {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误"
	}
	return rd
}
func InsertRoles(obj systemmodel.Roles, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Roles{}, "roles").SetKeys(true, "Id")
	count, _ := dbmap.SelectInt("select count(*) from roles where (Rolesname='" + obj.Rolesname + "')")
	if count > 0 {
		rd.Rcode = "1004"
		rd.Reason = "数据提交错误，此模块功能代码已被占用"
	} else {
		inerr := dbmap.Insert(&obj)
		core.CheckErr(inerr, "systemmodelDataAccess|InsertRoles|添加角色")
		if inerr == nil {
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1003"
			rd.Reason = inerr.Error()
		}
	}
	return rd
}
func DeleteRoles(obj systemmodel.Roles, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Roles{}, "roles").SetKeys(true, "Id")
	delsql1 := "delete from rolemodulefunctioncenter where Rolemodulecenterid in(select Id from rolemodulecenter where Rolesid=?);"
	delsql2 := "delete from rolemodulecenter where Rolesid=?;"
	_, inerr := dbmap.Exec(delsql1, obj.Id)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteRoles|删除角色模块功能中间表数据失败")
	_, inerr = dbmap.Exec(delsql2, obj.Id)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteRoles|删除角色模块中间表数据失败")
	_, inerr = dbmap.Delete(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|DeleteRoles|删除角色")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	return rd
}
func UpdateRoles(obj systemmodel.Roles, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Roles{}, "roles").SetKeys(true, "Id")
	_, inerr := dbmap.Update(&obj)
	core.CheckErr(inerr, "systemmodelDataAccess|UpdateRoles|修改角色")
	if inerr == nil {
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = inerr.Error()
	}
	return rd
}

//获取所有的模块下的功能数据
func QueryRolesList(pg core.PageData, obj systemmodel.Roles, dbmap *gorp.DbMap) (rd core.Returndata) {
	sqlcount := "select count(*) from Roles as smf where 1=1"
	sqlwhere := ""
	if obj.Rolesname != "" {
		sqlwhere = " and Rolesname like '%" + obj.Rolesname + "%'"
	}
	if obj.Id > 0 {
		sqlwhere = " and Id=" + strconv.Itoa(obj.Id)
	}
	count, err := dbmap.SelectInt(sqlcount + sqlwhere)
	core.CheckErr(err, "systemmodelDataAccess|QueryRolesList|获取所有的角色数据")
	if err == nil {
		pg.PageCount = int(count)
		sql := "select * from Roles as smf where 1=1"
		sql = sql + sqlwhere + core.GetLimitString(pg)
		var list []systemmodel.Roles
		_, err = dbmap.Select(&list, sql)
		core.CheckErr(err, "systemmodelDataAccess|QueryRolesList|获取所有的角色数据")
		if err == nil {
			pg.PageData = list
			rd.Result = pg
			rd.Rcode = "1000"
		} else {
			rd.Rcode = "1003"
			rd.Reason = "数据读取错误"
		}
	} else {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误"
	}
	return rd
}

func SaveSystemUsersmmodule(Rolesid int, systemmodelId int, smlist []systemmodel.Systemmodule, dbmap *gorp.DbMap) {
	//判断模块是否还有上级模块,并且递归调用
	for _, v := range smlist {
		if v.Id == systemmodelId {
			if v.Superiormoduleid > 0 {
				SaveSystemUsersmmodule(Rolesid, v.Superiormoduleid, smlist, dbmap)
			}
			var vv systemmodel.Rolemodulecenter
			dbmap.SelectOne(&vv, "select * from Rolemodulecenter where Rolesid=? and Systemmoduleid=?;", Rolesid, v.Id)
			if vv.Id == 0 { //未找到数据
				vv.Rolesid = Rolesid
				vv.Systemmoduleid = v.Id
				vv.State = 0
				dbmap.Insert(&vv)
			}
		}
	}
	//判断模块是否添加数据到模块角色中间表中
}
func SaveSystemConfig(obj1 []systemmodel.Rolemodulefunctioncenter, obj2 []systemmodel.Rolemodulecenter, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Rolemodulecenter{}, "rolemodulecenter").SetKeys(true, "Id")
	dbmap.AddTableWithName(systemmodel.Rolemodulefunctioncenter{}, "rolemodulefunctioncenter").SetKeys(true, "Id")
	if len(obj2) == 0 {
		rd.Rcode = "1004"
		rd.Reason = "数据提交错误，此模块功能代码已被占用"
	} else {
		var funclist []systemmodel.Systemmodulefunctions
		var smlist []systemmodel.Systemmodule
		sql := "select * from systemmodulefunctions"
		_, seleerr := dbmap.Select(&funclist, sql)
		core.CheckErr(seleerr, "systemmodelDataAccess|SaveSystemConfig|查找功能列表错误")
		sql = "select * from systemmodule"
		_, seleerr = dbmap.Select(&smlist, sql)
		core.CheckErr(seleerr, "systemmodelDataAccess|SaveSystemConfig|查找模块列表错误")
		for _, v := range obj2 {
			if v.Id == 0 {
				SaveSystemUsersmmodule(v.Rolesid, v.Systemmoduleid, smlist, dbmap)
				var vv systemmodel.Rolemodulecenter
				//判断此角色和模块是否有数据
				seleerr = dbmap.SelectOne(&vv, "select * from Rolemodulecenter where Rolesid=? and Systemmoduleid=?;", v.Rolesid, v.Systemmoduleid)
				core.CheckErr(seleerr, "systemmodelDataAccess|SaveSystemConfig|判断此角色和模块是否有数据")
				if vv.Id > 0 { //找到数据
					for _, iv := range obj1 { //找到角色模块功能中间数据
						for _, fcv := range funclist { //循环判断功能是否存在
							if iv.Systemmodulefunctionsid == fcv.Id && vv.Systemmoduleid == fcv.Systemmoduleid {
								var objiv systemmodel.Rolemodulefunctioncenter
								seleerr = dbmap.SelectOne(&objiv, "select * from Rolemodulefunctioncenter where Rolemodulecenterid=? and Systemmodulefunctionsid=?;", vv.Id, iv.Systemmodulefunctionsid)
								core.CheckErr(seleerr, "systemmodelDataAccess|SaveSystemConfig|找到角色模块功能中间数据")
								if objiv.Id == 0 { //未找到角色模块功能中间数据没找到
									iv.Rolemodulecenterid = vv.Id
									seleerr = dbmap.Insert(&iv)
									core.CheckErr(seleerr, "systemmodelDataAccess|SaveSystemConfig|将角色模块功能中间数据插入数据库中去")
								}
							}
						}
					}
				} else {
					seleerr = dbmap.Insert(&v)
					for _, iv := range obj1 {
						for _, fcv := range funclist {
							if iv.Systemmodulefunctionsid == fcv.Id && v.Systemmoduleid == fcv.Systemmoduleid {
								if core.CheckErr(seleerr, "插入角色模块中间表数据错误") {
									iv.Rolemodulecenterid = v.Id
									seleerr = dbmap.Insert(&iv)
								}
							}
						}
					}
				}
			}
		}
		rd.Rcode = "1000"
	}
	return rd
}
func DelSystemConfig(obj1 []systemmodel.Rolemodulefunctioncenter, obj2 []systemmodel.Rolemodulecenter, dbmap *gorp.DbMap) (rd core.Returndata) {
	dbmap.AddTableWithName(systemmodel.Rolemodulecenter{}, "rolemodulecenter").SetKeys(true, "Id")
	dbmap.AddTableWithName(systemmodel.Rolemodulefunctioncenter{}, "rolemodulefunctioncenter").SetKeys(true, "Id")
	var seleerr error
	if len(obj1) == 0 && len(obj2) == 0 {
		rd.Rcode = "1004"
		rd.Reason = "数据提交错误，此模块功能代码已被占用"
	} else {
		for _, v := range obj1 {
			_, seleerr = dbmap.Delete(&v)
			core.CheckErr(seleerr, "删除角色模块中间数据")
		}
		for _, iv := range obj2 {
			_, seleerr = dbmap.Delete(&iv)
			core.CheckErr(seleerr, "删除角色模块功能中间数据")
		}
		rd.Rcode = "1000"
	}
	return rd
}

//获取所有的模块下的功能数据
func QuerySetSystemConfig(Id int, dbmap *gorp.DbMap) (rd core.Returndata) {
	sql1 := "select rmc.Id,rmc.Rolesid,rmc.Systemmoduleid,rmc.State from roles as rs inner join rolemodulecenter as rmc on rs.Id=rmc.Rolesid inner join systemmodule as sm on rmc.Systemmoduleid=sm.Id where rs.Id=? order by ModuleIndex;"
	sql2 := "select rmfc.Id,rmfc.Systemmodulefunctionsid,rmfc.Rolemodulecenterid from roles as rs inner join rolemodulecenter as rmc on rs.Id=rmc.Rolesid"
	sql2 = sql2 + " inner join systemmodule as sm on rmc.Systemmoduleid=sm.Id inner join rolemodulefunctioncenter as rmfc on rmfc.Rolemodulecenterid=rmc.Id"
	sql2 = sql2 + " inner join systemmodulefunctions as smf on rmfc.Systemmodulefunctionsid=smf.Id where rs.Id=?;"
	var list1 []systemmodel.Rolemodulecenter
	var list2 []systemmodel.Rolemodulefunctioncenter
	_, sserr1 := dbmap.Select(&list1, sql1, Id)
	_, sserr2 := dbmap.Select(&list2, sql2, Id)
	core.CheckErr(sserr1, "systemmodelDataAccess|QuerySetSystemConfig")
	core.CheckErr(sserr2, "systemmodelDataAccess|QuerySetSystemConfig")
	if sserr1 == nil && sserr2 == nil {
		arrlist := make([]interface{}, 2)
		arrlist[0] = list1
		arrlist[1] = list2
		rd.Result = arrlist
		rd.Rcode = "1000"
	} else {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误"
		if sserr1 != nil {
			rd.Reason = rd.Reason + sserr1.Error()
		}
		if sserr2 != nil {
			rd.Reason = rd.Reason + sserr2.Error()
		}
	}
	return rd
}
