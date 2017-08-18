package systemmodelControllers

import (
	"encoding/json"
	"fmt"
	"time"
	//	"fmt"
	//	"time"
	//	"dev.project/BackEndCode/devserver/DataAccess/basicsetDataAccess"
	"dev.project/BackEndCode/devserver/DataAccess/curriculumDataAccess"
	"dev.project/BackEndCode/devserver/DataAccess/usersDataAccess"
	//	"dev.project/BackEndCode/devserver/model/basicset"
	"dev.project/BackEndCode/devserver/model/core"
	"dev.project/BackEndCode/devserver/model/curriculum"
	"dev.project/BackEndCode/devserver/viewmodel"

	"github.com/gin-gonic/gin"
)

// 查询学科类型数据
func GetSubjectclassListPG(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryCurriculumWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetSubjectclassList|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetSubjectclassList|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetSubjectclassListPG", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.GetSubjectclassListPG(obj, page, dbmap)
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

//添加学科类型数据
func AddSubjectclass(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Subjectclass
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddSubjectclass|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddSubjectclass", dbmap)
		if rd.Rcode == "1000" {
			isUnique := curriculumDataAccess.QueryUniqueSubjectcode(obj.Subjectcode, obj.Subjectname, dbmap) //检查唯一约束
			if isUnique > 0 {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:此编码已经存在"
			} else {
				adderr := curriculumDataAccess.AddSubjectclass(&obj, dbmap)
				core.CheckErr(adderr, "systemmodelControllers|AddSubjectclass|添加校区数据失败")
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

//修改学科类型数据
func ChangeSubjectclass(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Subjectclass
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeSubjectclass|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeSubjectclass", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.UpdateSubjectclass(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|ChangeSubjectclass|添加校区数据失败")
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

//删除学科类型数据
func DelSubjectclass(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Subjectclass
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelSubjectclass|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelSubjectclass", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.DeleteSubjectclass(&obj, dbmap)
			//			core.CheckErr(adderr, "systemmodelControllers|DelSubjectclass|删除校区数据失败")
			//			if adderr == nil {
			//				rd.Rcode = "1000"
			//				rd.Reason = "修改成功"
			//			} else {
			//				rd.Rcode = "1001"
			//				rd.Reason = "修改失败:" + adderr.Error()
			//			}
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

//查询详细课程
func GetCurriculumsInfo(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryCurriculumWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumsInfo|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumsInfo|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetCurriculumsInfo", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.GetCurriculumsInfo(obj, page, dbmap)
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

//查询课程
func GetCurriculumsListPG(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryCurriculumWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumsListPG|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumsListPG|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetCurriculumsListPG", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.GetCurriculumsListPG(obj, page, dbmap)
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

//添加课程
func AddCurriculums(c *gin.Context) {
	var rd core.Returndata
	type AddCurriculumsandChapters struct {
		Curriculumsinfo curriculum.Curriculums
		Chapterslist    []curriculum.Chapters
	}
	var ccinfo AddCurriculumsandChapters
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	fmt.Println(strdata.(string))
	errs1 := json.Unmarshal(data, &ccinfo)
	core.CheckErr(errs1, "systemmodelControllers|AddCurriculums|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddCurriculums", dbmap)
		if rd.Rcode == "1000" {
			rd.Rcode = ""
			ccinfo.Curriculumsinfo.Createdate = time.Now().Format("2006-01-02 15:04:05")
			adderr := curriculumDataAccess.AddCurriculums(&ccinfo.Curriculumsinfo, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|AddCurriculums|添加课程数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Result = ccinfo.Curriculumsinfo.Id
				for _, v := range ccinfo.Chapterslist {
					v.Curriculumsid = ccinfo.Curriculumsinfo.Id
					curriculumDataAccess.AddChapters(&v, dbmap)
				}
			} else {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:" + adderr.Error()
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

//修改课程
func ChangeCurriculums(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Curriculums
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeCurriculums|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeCurriculums", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.UpdateCurriculums(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|ChangeCurriculums|修改楼栋数据失败")
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

//删除课程
func DelCurriculums(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Curriculums
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelCurriculums|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelCurriculums", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.DeleteCurriculums(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelCurriculums|删除楼栋数据失败")
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

//获取所有的课程章节列表
func GetChaptersListPG(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryCurriculumWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetChaptersListPG|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetChaptersListPG|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetChaptersListPG", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.GetChaptersListPG(obj, page, dbmap)
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

//添加课程章节
func AddChapters(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Chapters
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddChapters|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddChapters", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.AddChapters(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|AddChapters|添加楼层数据失败")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Result = obj.Id
			} else {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改课程章节
func ChangeChapters(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Chapters
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeChapters|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeChapters", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.UpdateChapters(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|ChangeChapters|修改楼层数据失败")
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

//删除课程章节
func DelChapters(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Chapters
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelChapters|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelChapters", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.DeleteChapters(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelChapters|删除楼层数据失败")
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

//获取所有的课程班级中间数据
func GetCurriculumsClassCentreListPG(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryCurriculumWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumsClassCentreListPG|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumsClassCentreListPG|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetCurriculumsClassCentreListPG", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.GetCurriculumsClassCentreListPG(obj, page, dbmap)
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

//添加课程班级中间数据
func AddCurriculumsClassCentre(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.PostAddCurriculumsclasscentre //basicset.Classrooms
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddCurriculumsClassCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddCurriculumsClassCentre", dbmap)
		if rd.Rcode == "1000" {
			if obj.Classesid <= 0 {
				rd.Rcode = "1002"
				rd.Reason = "请选择上课的班级"
			} else if obj.Curriculumsid <= 0 {
				rd.Rcode = "1002"
				rd.Reason = "请选择所教学的课程"
			} else if obj.TeacherID <= 0 {
				rd.Rcode = "1002"
				rd.Reason = "请指定主讲的老师"
			} else {
				_, errs1 = dbmap.Exec("call Addkebd(?,?,?,?,?,?,?,?);", obj.Classesid, obj.Curriculumsid, obj.TeacherID, 0, obj.Islive, obj.Isondemand, 0, 0)
				if errs1 == nil {
					rd.Rcode = "1000"
				} else {
					rd.Rcode = "1001"
					rd.Reason = "数据保存错误"
				}
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改课程班级中间数据
func ChangeCurriculumsClassCentre(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.PostAddCurriculumsclasscentre //curriculum.Curriculumsclasscentre //basicset.Classrooms
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeCurriculumsClassCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeCurriculumsClassCentre", dbmap)
		if rd.Rcode == "1000" {
			objccc := curriculum.Curriculumsclasscentre{Id: obj.Id, Usersid: obj.TeacherID, Islive: obj.Islive, Isondemand: obj.Isondemand}
			rd = curriculumDataAccess.UpdateCurriculumsClassCentre(&objccc, dbmap)
			//			core.CheckErr(adderr, "systemmodelControllers|ChangeCurriculumsClassCentre|修改课程班级中间数据失败")
			//			if adderr == nil {
			//				rd.Rcode = "1000"
			//				rd.Reason = "修改成功"
			//			} else {
			//				rd.Rcode = "1001"
			//				rd.Reason = "修改失败:" + adderr.Error()
			//			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除课程班级中间数据
func DelCurriculumsClassCentre(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Curriculumsclasscentre //basicset.Classrooms
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelCurriculumsClassCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelCurriculumsClassCentre", dbmap)
		if rd.Rcode == "1000" {
			rd.Rcode = ""
			rd = curriculumDataAccess.DeleteCurriculumsClassCentre(&obj, dbmap)
			//			core.CheckErr(adderr, "systemmodelControllers|DelCurriculumsClassCentre|删除课程班级中间数据失败")
			//			if adderr == nil {
			//				rd.Rcode = "1000"
			//				rd.Reason = "修改成功"
			//			} else {
			//				rd.Rcode = "1001"
			//				rd.Reason = "修改失败:" + adderr.Error()
			//			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取所有的课程班级中间数据
func GetCurriculumClassroomChapterCentreListPG(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.QueryCurriculumWhere
	var page core.PageData //获取分页设置数据
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumClassroomChapterCentreListPG|模块数据转换失败")
	errs1 = json.Unmarshal(data, &page)
	core.CheckErr(errs1, "systemmodelControllers|GetCurriculumClassroomChapterCentreListPG|模块分页数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetCurriculumClassroomChapterCentreListPG", dbmap)
		if rd.Rcode == "1000" {
			rd = curriculumDataAccess.GetCurriculumclassroomchaptercentreListPG(obj, page, dbmap)
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

//添加课程班级章节中间数据
func GetCheckUpData(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.GetCurriculumclassroomchaptercentreList
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddCurriculumClassroomChapterCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		//		lgs, _ := c.Get("users")
		//		lg := lgs.(core.BasicsToken)
		list := curriculumDataAccess.QueryCheckUpData(obj, dbmap)
		if len(list) > 0 {
			rd.Rcode = "1001"
			reason := ""
			for _, v := range list {
				reason = reason + fmt.Sprintf("此课程时间安排上和 %s 老师课程冲突 \n", v.Truename)
			}
		} else {
			rd.Rcode = "1000"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//添加课程班级章节中间数据
func AddCurriculumClassroomChapterCentre(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.PostAddCurriculumclassroomchaptercentre
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|AddCurriculumClassroomChapterCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "AddCurriculumClassroomChapterCentre", dbmap)
		if rd.Rcode == "1000" {
			timestr := time.Now().Format("2006-01-02 15:04:05")
			ObjCCCC := curriculum.Curriculumclassroomchaptercentre{Curriculumsclasscentreid: obj.Curriculumsclasscentreid, Chaptersid: obj.Chaptersid, Usersid: obj.TeacherID, Begindate: obj.Begindate, Enddate: obj.Enddate, Islive: obj.Islive, Isondomian: obj.Isondomian, Createdate: timestr}
			adderr := curriculumDataAccess.AddCurriculumclassroomchaptercentre(&ObjCCCC, obj.Classroomid, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|AddCurriculumClassroomChapterCentre|添加课程班级中间数据")
			if adderr == nil {
				rd.Rcode = "1000"
				rd.Result = obj.Id
			} else {
				rd.Rcode = "1001"
				rd.Reason = "添加失败:" + adderr.Error()
			}
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//修改课程班级章节中间数据
func ChangeCurriculumClassroomChapterCentre(c *gin.Context) {
	var rd core.Returndata
	var obj viewmodel.PostAddCurriculumclassroomchaptercentre //curriculum.Curriculumclassroomchaptercentre
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|ChangeCurriculumClassroomChapterCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "ChangeCurriculumClassroomChapterCentre", dbmap)
		if rd.Rcode == "1000" {
			ObjCCCC := curriculum.Curriculumclassroomchaptercentre{Id: obj.Id, Curriculumsclasscentreid: obj.Curriculumsclasscentreid, Chaptersid: obj.Chaptersid, Usersid: obj.TeacherID, Begindate: obj.Begindate, Enddate: obj.Enddate, Islive: obj.Islive, Isondomian: obj.Isondomian}
			rd = curriculumDataAccess.UpdateCurriculumclassroomchaptercentre(&ObjCCCC, obj.Classroomid, dbmap)
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//删除课程班级章节中间数据
func DelCurriculumClassroomChapterCentre(c *gin.Context) {
	var rd core.Returndata
	var obj curriculum.Curriculumclassroomchaptercentre
	strdata, _ := c.Get("data")
	data := []byte(strdata.(string))
	errs1 := json.Unmarshal(data, &obj)
	core.CheckErr(errs1, "systemmodelControllers|DelCurriculumClassroomChapterCentre|模块数据转换失败")
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lgs, _ := c.Get("users")
		lg := lgs.(core.BasicsToken)
		rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "DelCurriculumClassroomChapterCentre", dbmap)
		if rd.Rcode == "1000" {
			adderr := curriculumDataAccess.DeleteCurriculumclassroomchaptercentre(&obj, dbmap)
			core.CheckErr(adderr, "systemmodelControllers|DelCurriculumClassroomChapterCentre|删除课程班级中间数据失败")
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
