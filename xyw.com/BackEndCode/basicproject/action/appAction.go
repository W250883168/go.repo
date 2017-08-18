package app

import (
	// "errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	PathUtil "path/filepath"
	"runtime"
	"strconv"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"

	"xutils/xerr"

	app "basicproject/Controllers"
	actiondata "basicproject/Controllers/actiondataControllers"
	basicset "basicproject/Controllers/basicsetControllers"
	//	curriculum "basicproject/Controllers/curriculumControllers"
	//	devctrl "basicproject/Controllers/deviceControllers"
	//	live "basicproject/Controllers/liveControllers"
	//	vodControl "basicproject/LiveControllers/vodControl"
	sysm "basicproject/Controllers/systemmodelControllers"
	users "basicproject/Controllers/usersControllers"
	core "xutils/xcore"
)

func LoadWebModel(r *gin.Engine) {
	r.Static("/web", "../../FrontEndCode/templates")
	r.Static("/web2", "../../FrontEndCode/templates2")
	r.StaticDows("/webStatic", "../../FrontEndCode/templates")
	r.StaticDows("/web2Static", "../../FrontEndCode/templates2")

	r.POST("/upfile", onUpfile_Handler)

}

func LoadAction(r *gin.Engine) {
	r.POST("/login", _AllowCrossDomain, app.GetLogin)              //登录
	r.POST("/getapp", _AllowCrossDomain, app.GetHomeApp)           //获取首页权限模块
	r.POST("/getServerTime", _AllowCrossDomain, app.GetServerTime) //获取服务器时间戳
	basicsetat := r.Group("/basicset", func(c *gin.Context) {})
	basicsetat.GET("/getall", _AllowCrossDomain, basicset.GetAll)                                           //获取校区、楼栋、楼层
	basicsetat.GET("/getcampus", _AllowCrossDomain, basicset.GetCampus)                                     //获取校区
	basicsetat.GET("/getbuilding", _AllowCrossDomain, basicset.GetBuilding)                                 //获取楼栋
	basicsetat.GET("/getfloorsandrooms", _AllowCrossDomain, basicset.GetFloorsandrooms)                     //获取楼层和教室
	basicsetat.POST("/queryclassroom", _AllowCrossDomain, basicset.QueryClassroom)                          //根据校区Id、楼栋Id、楼层id、用户Id、角色、token令牌查询出相关教室的情况
	basicsetat.POST("/setorcancelcollection", _AllowCrossDomain, actiondata.SetOrCancelClassroomCollection) //添加或取消教室收藏
	basicsetat.GET("/getclassroominfo", _AllowCrossDomain, basicset.GetClassRoomInfo)                       //查询教室当前状态
	basicsetat.GET("/queryteachers", _AllowCrossDomain, basicset.QueryTeachers)                             //查询教师数据

	basicsetat.POST("/getpeoples", _AllowCrossDomain, basicset.GetQueryPeoples)                        //根据楼层Id数组，楼栋Id数组，学院Id数组
	basicsetat.POST("/getclassroompeopleinfo", _AllowCrossDomain, basicset.GetClassRoomPeopleInfo)     //根据教室Id查询教室内所有人员列表信息
	basicsetat.POST("/getclassroompeoplecount", _AllowCrossDomain, basicset.GetClassRoomPeopleCount)   //根据教室Id查询教室内所有人员列表信息汇总
	basicsetat.POST("/getstreampeoplesanalysis", _AllowCrossDomain, basicset.GetStreamPeoplesAnalysis) //根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计
	userat := r.Group("/user", func(c *gin.Context) {})
	userat.POST("/mymessagelist", _AllowCrossDomain, users.GetMyMessageList) //获取推送用户消息
	userat.POST("/myisnew", _AllowCrossDomain, users.GetIsNewMessageList)    //获取推送用户消息
	userat.POST("/studentsinfo", _AllowCrossDomain, users.GetStudentsinfo)   //获取学生的详细信息

}

func LoadSystemAction(r *gin.Engine) {
	basesys := r.Group("/system", _AllowCrossDomain, app.SystemGroupFunc)
	sm := basesys.Group("/sm", func(c *gin.Context) { //配置用户和权限等相关模块
	})
	sm.POST("/getsystemmodel", _AllowCrossDomain, sysm.GetSystemModelList)   //获取所有的系统模块
	sm.POST("/addsystemmodel", _AllowCrossDomain, sysm.AddSystemModel)       //添加系统模块
	sm.POST("/delsystemmodel", _AllowCrossDomain, sysm.DelSystemModel)       //删除系统模块
	sm.POST("/updatesystemmodel", _AllowCrossDomain, sysm.ChangeSystemModel) //修改系统模块

	sm.POST("/getsystemmodelfunc", _AllowCrossDomain, sysm.GetSystemModelFunctionsList)   //获取所有的系统模块下的功能
	sm.POST("/addsystemmodelfunc", _AllowCrossDomain, sysm.AddSystemModelFunctions)       //添加系统模块功能
	sm.POST("/delsystemmodelfunc", _AllowCrossDomain, sysm.DelSystemModelFunctions)       //删除系统模块功能
	sm.POST("/updatesystemmodelfunc", _AllowCrossDomain, sysm.ChangeSystemModelFunctions) //修改系统模块功能

	sm.POST("/getroles", _AllowCrossDomain, sysm.GetRolesList)                    //获取所有的系统角色
	sm.POST("/addroles", _AllowCrossDomain, sysm.AddRoles)                        //添加系统角色
	sm.POST("/delroles", _AllowCrossDomain, sysm.DelRoles)                        //删除系统角色
	sm.POST("/updateroles", _AllowCrossDomain, sysm.ChangeRoles)                  //修改系统角色
	sm.POST("/getsetsystemconfigall", _AllowCrossDomain, sysm.GetSetSystemConfig) //获取权限设置中间表的数据集合
	sm.POST("/setsystemconfig", _AllowCrossDomain, sysm.SaveSystemConfig)         //保存权限设置中间表的数据集合

	bs := basesys.Group("/bs", func(c *gin.Context) {})

	bs.POST("/campuslist", _AllowCrossDomain, sysm.GetCampusList)
	bs.POST("/campusadd", _AllowCrossDomain, sysm.AddCampus)
	bs.POST("/campuschange", _AllowCrossDomain, sysm.CampusUpdate)
	bs.POST("/campusdel", _AllowCrossDomain, sysm.DelCampus)

	bs.POST("/buildinglist", _AllowCrossDomain, sysm.GetBuildingList)
	bs.POST("/buildingadd", _AllowCrossDomain, sysm.AddBuilding)
	bs.POST("/buildingchange", _AllowCrossDomain, sysm.ChangeBuilding)
	bs.POST("/buildingdel", _AllowCrossDomain, sysm.DelBuilding)

	bs.POST("/floorslist", _AllowCrossDomain, sysm.GetFloorsList)
	bs.POST("/floorsadd", _AllowCrossDomain, sysm.AddFloors)
	bs.POST("/floorschange", _AllowCrossDomain, sysm.ChangeFloors)
	bs.POST("/floorsdel", _AllowCrossDomain, sysm.DelFloors)

	bs.POST("/classroomslist", _AllowCrossDomain, sysm.GetClassroomsList)
	bs.POST("/classroomsadd", _AllowCrossDomain, sysm.AddClassrooms)
	bs.POST("/classroomschange", _AllowCrossDomain, sysm.ChangeClassrooms)
	bs.POST("/classroomsdel", _AllowCrossDomain, sysm.DelClassrooms)

	bs.POST("/collegelist", _AllowCrossDomain, sysm.GetCollegeList)
	bs.POST("/collegeadd", _AllowCrossDomain, sysm.AddCollege)
	bs.POST("/collegechange", _AllowCrossDomain, sysm.ChangeCollege)
	bs.POST("/collegedel", _AllowCrossDomain, sysm.DelCollege)

	bs.POST("/majorlist", _AllowCrossDomain, sysm.GetMajorList)
	bs.POST("/majoradd", _AllowCrossDomain, sysm.AddMajor)
	bs.POST("/majorchange", _AllowCrossDomain, sysm.ChangeMajor)
	bs.POST("/majordel", _AllowCrossDomain, sysm.DelMajor)

	bs.POST("/classeslist", _AllowCrossDomain, sysm.GetClassesList)
	bs.POST("/classesadd", _AllowCrossDomain, sysm.AddClasses)
	bs.POST("/classeschange", _AllowCrossDomain, sysm.ChangeClasses)
	bs.POST("/classesdel", _AllowCrossDomain, sysm.DelClasses)

	us := basesys.Group("/us", func(c *gin.Context) {})
	us.POST("/userslist", _AllowCrossDomain, sysm.GetUsersList2)
	us.POST("/usersadd", _AllowCrossDomain, sysm.AddUsers)
	us.POST("/userschange", _AllowCrossDomain, sysm.ChangeUsers)
	us.POST("/usersdel", _AllowCrossDomain, sysm.DelUsers)

	us.POST("/studentslist", _AllowCrossDomain, sysm.GetStudentsList)
	us.POST("/studentsadd", _AllowCrossDomain, sysm.AddStudents)
	us.POST("/studentschange", _AllowCrossDomain, sysm.ChangeStudents)
	us.POST("/studentsdel", _AllowCrossDomain, sysm.DelStudents)

	us.POST("/teacherlist", _AllowCrossDomain, sysm.GetTeacherList)
	us.POST("/teacheradd", _AllowCrossDomain, sysm.AddTeacher)
	us.POST("/teacherchange", _AllowCrossDomain, sysm.ChangeTeacher)
	us.POST("/teacherdel", _AllowCrossDomain, sysm.DelTeacher)

}

//共用的接受上传图片
func onUpfile_Handler(c *gin.Context) {
	defer xerr.CatchPanic() // 捕获异常
	// 允许跨域访问
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// HTTP响应数据
	var err error
	var rd = core.Returndata{Rcode: "1003", Reason: "未提交数据文件"}
	defer func() {
		if err != nil {
			rd.Reason = err.Error()
		}
		c.JSON(http.StatusOK, rd)
	}()

	f, fh, err := c.Request.FormFile("file")
	log.Printf("f1: %+v\n fh: %+v \n %s", f, fh, err)
	xerr.ThrowPanic(err)

	log.Printf("%+v", c.Request)
	log.Println(c.Request.FormValue("filename"))
	content_len := c.Request.Header.Get("Content-Length")
	log.Println("文件大小(字节)： ", content_len)
	filesize, _ := strconv.Atoi(content_len)
	file_oversize := filesize > 1024*1024*2
	if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
		for k, pFileArr := range c.Request.MultipartForm.File {
			log.Println(k)
			if len(pFileArr) > 0 {
				var inFile multipart.File
				inFile, err = pFileArr[0].Open()
				xerr.ThrowPanic(err)
				defer inFile.Close()

				curdir := core.GetCurrentDirectory() //获取当前运行的目录
				curdir += "/../../FrontEndCode/templates/upfile"
				timestr := time.Now().Format("20060102")
				curdir += "/" + timestr
				log.Println(curdir)
				if !core.Exist(curdir) { //判断文件目录是否存在
					err = os.Mkdir(curdir, os.ModePerm) //创建文件目录
					xerr.ThrowPanic(err)
				}

				filename := strconv.FormatInt(time.Now().UnixNano(), 10) + PathUtil.Ext(fh.Filename)
				outFile, err := os.Create(curdir + "/" + filename)
				xerr.ThrowPanic(err)
				defer outFile.Close()
				log.Println(outFile.Name())

				rd.Rcode = "1002"
				rd.Rcode = "操作失败!"
				_, err = io.Copy(outFile, inFile)
				xerr.ThrowPanic(err)

				rd.Rcode = "1000"
				rd.Reason = "操作成功!"
				rd.Result = "/web/upfile/" + timestr + "/" + filename
				if file_oversize {
					outFile.Close()
					os.Remove(outFile.Name())

					rd.Rcode = "1000"
					rd.Reason = "文件大小超限！"
					rd.Result = ""
					break
				}
			}
		}
	}

}

func _AllowCrossDomain(c *gin.Context) {
	//设置参数，允许跨域调用
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	c.Writer.Header().Set("Content-Type", "application/x-www-form-urlencoded")
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
