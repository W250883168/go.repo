package app

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"

	app "dev.project/BackEndCode/devserver/Controllers"
	vodControl "dev.project/BackEndCode/devserver/Controllers/LiveControllers/vodControl"
	actiondata "dev.project/BackEndCode/devserver/Controllers/actiondataControllers"
	curriculum "dev.project/BackEndCode/devserver/Controllers/curriculumControllers"
	devctrl "dev.project/BackEndCode/devserver/Controllers/deviceControllers"
	live "dev.project/BackEndCode/devserver/Controllers/liveControllers"
	sysm "dev.project/BackEndCode/devserver/Controllers/systemmodelControllers"
)

func LoadWebModel(r *gin.Engine) {
	//	r.Static("/web", "../../FrontEndCode/build/FrontEndCode/templates")
	//	r.Static("/web2", "../../FrontEndCode/build/FrontEndCode/templates2")
	//	r.StaticDows("/webStatic", "../../FrontEndCode/build/FrontEndCode/templates")
	//	r.StaticDows("/web2Static", "../../FrontEndCode/build/FrontEndCode/templates2")

	//	r.POST("/upfile", onUpfile_Handler)

}

func LoadAction(r *gin.Engine) {
	//	r.POST("/login", _AllowCrossDomain, app.GetLogin)              //登录
	//	r.POST("/getapp", _AllowCrossDomain, app.GetHomeApp)           //获取首页权限模块
	//	r.POST("/getServerTime", _AllowCrossDomain, app.GetServerTime) //获取服务器时间戳
	//	basicsetat := r.Group("/basicset", func(c *gin.Context) {})
	//	basicsetat.GET("/getall", _AllowCrossDomain, basicset.GetAll)                                           //获取校区、楼栋、楼层
	//	basicsetat.GET("/getcampus", _AllowCrossDomain, basicset.GetCampus)                                     //获取校区
	//	basicsetat.GET("/getbuilding", _AllowCrossDomain, basicset.GetBuilding)                                 //获取楼栋
	//	basicsetat.GET("/getfloorsandrooms", _AllowCrossDomain, basicset.GetFloorsandrooms)                     //获取楼层和教室
	//	basicsetat.POST("/queryclassroom", _AllowCrossDomain, basicset.QueryClassroom)                          //根据校区Id、楼栋Id、楼层id、用户Id、角色、token令牌查询出相关教室的情况
	//	basicsetat.POST("/setorcancelcollection", _AllowCrossDomain, actiondata.SetOrCancelClassroomCollection) //添加或取消教室收藏
	//	basicsetat.GET("/getclassroominfo", _AllowCrossDomain, basicset.GetClassRoomInfo)                       //查询教室当前状态
	//	basicsetat.GET("/queryteachers", _AllowCrossDomain, basicset.QueryTeachers)                             //查询教师数据

	//	basicsetat.POST("/getpeoples", _AllowCrossDomain, basicset.GetQueryPeoples)                      //根据楼层Id数组，楼栋Id数组，学院Id数组
	//	basicsetat.POST("/getclassroompeopleinfo", _AllowCrossDomain, basicset.GetClassRoomPeopleInfo)   //根据教室Id查询教室内所有人员列表信息
	//	basicsetat.POST("/getclassroompeoplecount", _AllowCrossDomain, basicset.GetClassRoomPeopleCount) //根据教室Id查询教室内所有人员列表信息汇总
	//	//根据时间段+校区+楼栋+楼层来综合查询,返回人流热力数据，人流统计数据，人群分析统计
	//	basicsetat.POST("/getstreampeoplesanalysis", _AllowCrossDomain, basicset.GetStreamPeoplesAnalysis)

	actionat := r.Group("/action", func(c *gin.Context) {})
	actionat.POST("/changeclassstate", _AllowCrossDomain, actiondata.ChangeClassState) //老师点击上课[老师在app上点击上课，老师在中控上点击上课,管理员在总控上点击上课]
	actionat.POST("/beginvideo", _AllowCrossDomain, actiondata.BeginVideo2)            //开始录制摄像
	actionat.POST("/endvideo", _AllowCrossDomain, actiondata.EndVideo2)
	actionat.POST("/closecomputer", _AllowCrossDomain, actiondata.CloseComputer) //关闭电脑
	actionat.POST("/udpclient", _AllowCrossDomain, actiondata.UdpClient)         //接受客户端发送过来的upd数据

	actionat.POST("/getpointtos", _AllowCrossDomain, actiondata.GetPointtos)                         //获取点到数据
	actionat.POST("/getccccid", _AllowCrossDomain, actiondata.GetCurriculumClassroomChapterCentreId) //获取教师马上要上课的课程章节等数据
	actionat.POST("/getclassclassroom", _AllowCrossDomain, actiondata.GetClassClassroomId)           //教师点击移动中控获取教室Id

	actionat.POST("/getcurriculums", _AllowCrossDomain, actiondata.GetCurriculumslist)           //根据时间段查询课表列表信息
	actionat.POST("/getwatchcurriculums", _AllowCrossDomain, actiondata.GetWatchCurriculumslist) //手表端根据时间段查询课表列表信息
	actionat.POST("/getcurriculumsinfo", _AllowCrossDomain, actiondata.GetCurriculumsinfo)       //根据课表Id查询详细信息
	actionat.POST("/getattendancelist", _AllowCrossDomain, actiondata.GetAttendanceQuerylist)    //根据时间段查询出勤记录
	actionat.POST("/gethistoryattendance", _AllowCrossDomain, actiondata.GetHistoryAttendance)   //根据课表Id查询详细信息

	actionat.POST("/updatepointtos", _AllowCrossDomain, actiondata.UpdatePointtos)              //更改学生的点到信息
	actionat.POST("/updatelistpointtos", _AllowCrossDomain, actiondata.UpdateListPointtos)      //更改学生的点到信息
	actionat.POST("/getmyfollow", _AllowCrossDomain, actiondata.GetMyAttentionRecord)           //获取我关注的课程信息
	actionat.POST("/getmyclassfollow", _AllowCrossDomain, actiondata.GetMyClassAttentionRecord) //获取我所在的班级的课程信息
	actionat.POST("/isfollowok", _AllowCrossDomain, actiondata.IsAttentionRecordOk)             //判断此课程是否关注了的课程
	actionat.POST("/setfollow", _AllowCrossDomain, actiondata.SetAttentionRecord)               //修改我关注的课程记录

	//	userat := r.Group("/user", func(c *gin.Context) {})
	//	userat.POST("/mymessagelist", _AllowCrossDomain, users.GetMyMessageList) //获取推送用户消息
	//	userat.POST("/myisnew", _AllowCrossDomain, users.GetIsNewMessageList)    //获取推送用户消息
	//	userat.POST("/studentsinfo", _AllowCrossDomain, users.GetStudentsinfo)   //获取学生的详细信息

	curriculumat := r.Group("/curriculum", func(c *gin.Context) {})
	curriculumat.POST("/attendanceanalysis", _AllowCrossDomain, curriculum.AttendanceAnalysis)               //管理者查看各种出勤统计分析统计年级人数比
	curriculumat.POST("/getageclassrate", _AllowCrossDomain, curriculum.GetEveryoneAverageclassrate)         //获取某课程下某个班所有的学生的平均到课率
	curriculumat.POST("/getcurriculumchaptersinfo", _AllowCrossDomain, curriculum.GetCurriculumChaptersInfo) //获取某教师下某班级所有课程包含章节所有到课数据
	curriculumat.POST("/getclassagerate", _AllowCrossDomain, curriculum.GetClassAveragerate)                 //获取班级的平均到课率
	curriculumat.POST("/getfilter", _AllowCrossDomain, curriculum.GetFilterData)                             //获取筛选条件页的数据
	curriculumat.POST("/getstudentscurriculm", _AllowCrossDomain, curriculum.GetStudentsClassesAvg)          //以学生Id查询到学生每个课程的平均出勤数据
	curriculumat.POST("/getsubjectclass", _AllowCrossDomain, curriculum.GetSubjectclassList)                 //查询学科分类信息

	liveat := r.Group("/live", func(c *gin.Context) {})
	liveat.POST("/querylivelist", _AllowCrossDomain, live.GetQueryLiveList)                 //视频搜索综合接口
	liveat.POST("/querylastweeklivelist", _AllowCrossDomain, live.GetQueryLastWeeklivelist) //视频搜索综合接口
	liveat.POST("/getabsenlist", _AllowCrossDomain, live.GetMyAbsentList)                   //获取我缺课的课程列表信息
	liveat.POST("/getliveinfo", _AllowCrossDomain, live.GetQueryLiveInfo)                   //获取我缺课的课程章节列表详细信息
	liveat.POST("/getliveenclosure", _AllowCrossDomain, live.GetQueryLiveEnclosure)         //查询课程下章节对应的附件数据
	liveat.POST("/updateplaynum", _AllowCrossDomain, live.UpdateLivePlayNum)                //更新播放次数
	liveat.POST("/updatedownloadnum", _AllowCrossDomain, live.UpdateLiveDownloadNum)        //更新播放次数

	vod := r.Group("/vod", func(c *gin.Context) {})
	vod.POST("/getvideolist", _AllowCrossDomain, vodControl.GetVideoList)
	vod.POST("/videodetails", _AllowCrossDomain, vodControl.GetVideoDetails)
	vod.POST("/updatevideodetails", _AllowCrossDomain, vodControl.UpdateVideoDetatils)
	vod.POST("/deletevideo", _AllowCrossDomain, vodControl.DeleteVideo)
	vod.POST("/deleleattachment", _AllowCrossDomain, vodControl.DeleteAttachment)

}

func LoadSystemAction(r *gin.Engine) {
	basesys := r.Group("/system", _AllowCrossDomain, app.SystemGroupFunc)
	//	sm := basesys.Group("/sm", func(c *gin.Context) { //配置用户和权限等相关模块
	//	})
	//	sm.POST("/getsystemmodel", _AllowCrossDomain, sysm.GetSystemModelList)   //获取所有的系统模块
	//	sm.POST("/addsystemmodel", _AllowCrossDomain, sysm.AddSystemModel)       //添加系统模块
	//	sm.POST("/delsystemmodel", _AllowCrossDomain, sysm.DelSystemModel)       //删除系统模块
	//	sm.POST("/updatesystemmodel", _AllowCrossDomain, sysm.ChangeSystemModel) //修改系统模块

	//	sm.POST("/getsystemmodelfunc", _AllowCrossDomain, sysm.GetSystemModelFunctionsList)   //获取所有的系统模块下的功能
	//	sm.POST("/addsystemmodelfunc", _AllowCrossDomain, sysm.AddSystemModelFunctions)       //添加系统模块功能
	//	sm.POST("/delsystemmodelfunc", _AllowCrossDomain, sysm.DelSystemModelFunctions)       //删除系统模块功能
	//	sm.POST("/updatesystemmodelfunc", _AllowCrossDomain, sysm.ChangeSystemModelFunctions) //修改系统模块功能

	//	sm.POST("/getroles", _AllowCrossDomain, sysm.GetRolesList)                    //获取所有的系统角色
	//	sm.POST("/addroles", _AllowCrossDomain, sysm.AddRoles)                        //添加系统角色
	//	sm.POST("/delroles", _AllowCrossDomain, sysm.DelRoles)                        //删除系统角色
	//	sm.POST("/updateroles", _AllowCrossDomain, sysm.ChangeRoles)                  //修改系统角色
	//	sm.POST("/getsetsystemconfigall", _AllowCrossDomain, sysm.GetSetSystemConfig) //获取权限设置中间表的数据集合
	//	sm.POST("/setsystemconfig", _AllowCrossDomain, sysm.SaveSystemConfig)         //保存权限设置中间表的数据集合

	//	bs := basesys.Group("/bs", func(c *gin.Context) {})

	//	bs.POST("/campuslist", _AllowCrossDomain, sysm.GetCampusList)
	//	bs.POST("/campusadd", _AllowCrossDomain, sysm.AddCampus)
	//	bs.POST("/campuschange", _AllowCrossDomain, sysm.CampusUpdate)
	//	bs.POST("/campusdel", _AllowCrossDomain, sysm.DelCampus)

	//	bs.POST("/buildinglist", _AllowCrossDomain, sysm.GetBuildingList)
	//	bs.POST("/buildingadd", _AllowCrossDomain, sysm.AddBuilding)
	//	bs.POST("/buildingchange", _AllowCrossDomain, sysm.ChangeBuilding)
	//	bs.POST("/buildingdel", _AllowCrossDomain, sysm.DelBuilding)

	//	bs.POST("/floorslist", _AllowCrossDomain, sysm.GetFloorsList)
	//	bs.POST("/floorsadd", _AllowCrossDomain, sysm.AddFloors)
	//	bs.POST("/floorschange", _AllowCrossDomain, sysm.ChangeFloors)
	//	bs.POST("/floorsdel", _AllowCrossDomain, sysm.DelFloors)

	//	bs.POST("/classroomslist", _AllowCrossDomain, sysm.GetClassroomsList)
	//	bs.POST("/classroomsadd", _AllowCrossDomain, sysm.AddClassrooms)
	//	bs.POST("/classroomschange", _AllowCrossDomain, sysm.ChangeClassrooms)
	//	bs.POST("/classroomsdel", _AllowCrossDomain, sysm.DelClassrooms)

	//	bs.POST("/collegelist", _AllowCrossDomain, sysm.GetCollegeList)
	//	bs.POST("/collegeadd", _AllowCrossDomain, sysm.AddCollege)
	//	bs.POST("/collegechange", _AllowCrossDomain, sysm.ChangeCollege)
	//	bs.POST("/collegedel", _AllowCrossDomain, sysm.DelCollege)

	//	bs.POST("/majorlist", _AllowCrossDomain, sysm.GetMajorList)
	//	bs.POST("/majoradd", _AllowCrossDomain, sysm.AddMajor)
	//	bs.POST("/majorchange", _AllowCrossDomain, sysm.ChangeMajor)
	//	bs.POST("/majordel", _AllowCrossDomain, sysm.DelMajor)

	//	bs.POST("/classeslist", _AllowCrossDomain, sysm.GetClassesList)
	//	bs.POST("/classesadd", _AllowCrossDomain, sysm.AddClasses)
	//	bs.POST("/classeschange", _AllowCrossDomain, sysm.ChangeClasses)
	//	bs.POST("/classesdel", _AllowCrossDomain, sysm.DelClasses)

	us := basesys.Group("/us", func(c *gin.Context) {})
	//	us.POST("/userslist", _AllowCrossDomain, sysm.GetUsersList2)
	//	us.POST("/usersadd", _AllowCrossDomain, sysm.AddUsers)
	//	us.POST("/userschange", _AllowCrossDomain, sysm.ChangeUsers)
	//	us.POST("/usersdel", _AllowCrossDomain, sysm.DelUsers)

	//	us.POST("/studentslist", _AllowCrossDomain, sysm.GetStudentsList)
	//	us.POST("/studentsadd", _AllowCrossDomain, sysm.AddStudents)
	//	us.POST("/studentschange", _AllowCrossDomain, sysm.ChangeStudents)
	//	us.POST("/studentsdel", _AllowCrossDomain, sysm.DelStudents)

	//	us.POST("/teacherlist", _AllowCrossDomain, sysm.GetTeacherList)
	//	us.POST("/teacheradd", _AllowCrossDomain, sysm.AddTeacher)
	//	us.POST("/teacherchange", _AllowCrossDomain, sysm.ChangeTeacher)
	//	us.POST("/teacherdel", _AllowCrossDomain, sysm.DelTeacher)

	/*
		校区管理
			校区查询
			校区增加、修改、删除
			校区楼栋查询
			校区楼栋增加、修改、删除
			校区楼层查询
			校区楼层增加、修改、删除
			校区教室查询
			校区教室增加、修改、删除
		学院管理
			学院查询
			学院增加、修改、删除
			专业查询
			专业增加、修改、删除
			班级查询
			班级增加、修改、删除
		用户管理
			用户查询
			用户增加、修改、删除
			学生查询
			学生增加、修改、删除
			教师查询
			教师增加、修改、删除
	*/
	us.POST("/subjectclasslist", _AllowCrossDomain, sysm.GetSubjectclassListPG)
	us.POST("/subjectclassadd", _AllowCrossDomain, sysm.AddSubjectclass)
	us.POST("/subjectclasschange", _AllowCrossDomain, sysm.ChangeSubjectclass)
	us.POST("/subjectclassdel", _AllowCrossDomain, sysm.DelSubjectclass)

	us.POST("/curriculumsinfo", _AllowCrossDomain, sysm.GetCurriculumsInfo)
	us.POST("/curriculumslist", _AllowCrossDomain, sysm.GetCurriculumsListPG)
	us.POST("/curriculumsadd", _AllowCrossDomain, sysm.AddCurriculums)
	us.POST("/curriculumschange", _AllowCrossDomain, sysm.ChangeCurriculums)
	us.POST("/curriculumsdel", _AllowCrossDomain, sysm.DelCurriculums)

	us.POST("/chapterslist", _AllowCrossDomain, sysm.GetChaptersListPG)
	us.POST("/chaptersadd", _AllowCrossDomain, sysm.AddChapters)
	us.POST("/chapterschange", _AllowCrossDomain, sysm.ChangeChapters)
	us.POST("/chaptersdel", _AllowCrossDomain, sysm.DelChapters)

	us.POST("/curriculumsclasscentrelist", _AllowCrossDomain, sysm.GetCurriculumsClassCentreListPG)
	us.POST("/curriculumsclasscentreadd", _AllowCrossDomain, sysm.AddCurriculumsClassCentre)
	us.POST("/curriculumsclasscentrechange", _AllowCrossDomain, sysm.ChangeCurriculumsClassCentre)
	us.POST("/curriculumsclasscentredel", _AllowCrossDomain, sysm.DelCurriculumsClassCentre)

	us.POST("/curriculumclassroomchaptercentrelist", _AllowCrossDomain, sysm.GetCurriculumClassroomChapterCentreListPG)
	us.POST("/getcheckup", _AllowCrossDomain, sysm.GetCheckUpData) //验证排课时间、地点、教师是否会起冲突
	us.POST("/curriculumclassroomchaptercentreadd", _AllowCrossDomain, sysm.AddCurriculumClassroomChapterCentre)
	us.POST("/curriculumclassroomchaptercentrechange", _AllowCrossDomain, sysm.ChangeCurriculumClassroomChapterCentre) //修改教室、修改上课时间
	us.POST("/curriculumclassroomchaptercentredel", _AllowCrossDomain, sysm.DelCurriculumClassroomChapterCentre)

	/*
		学科管理
			学科查询
			学科增加、修改、删除
		课程查询
		课程添加、修改、删除
		课程章节查询
		课程章节添加、修改、删除
		课程班级中间数据查询，课程班级章节中间数据查询
		课程班级中间数据添加、修改、删除，课程班级章节中间数据添加、修改、删除
	*/

}

func LoadDeviceAction(r *gin.Engine) {
	d := r.Group("/device", func(c *gin.Context) {})
	d.POST("/getUseLogList", _AllowCrossDomain, devctrl.GetDeviceUseLogList)         //获取设备使用日志
	d.POST("/getOperateLogList", _AllowCrossDomain, devctrl.GetDeviceOperateLogList) //获取设备操作日志
	d.POST("/getAlertInfoList", _AllowCrossDomain, devctrl.GetDeviceAlertInfoList)   //获取设备预警信息
	d.POST("/getFaultInfoList", _AllowCrossDomain, devctrl.GetDeviceFaultInfoList)   //获取设备故障信息

	d.POST("/getClassroomStatusList", _AllowCrossDomain, devctrl.GetClassroomStatusList) //获取教室状态信息

	d.POST("/getAllOperateLogList", _AllowCrossDomain, devctrl.GetAllOperateLogList)       //获取所有设备操作日志
	d.POST("/getAllAlertInfoList", _AllowCrossDomain, devctrl.GetAllAlertInfoList)         //获取所有设备预警消息
	d.POST("/getAllFaultInfoList", _AllowCrossDomain, devctrl.GetAllFaultInfoList)         //获取所有设备故障信息
	d.POST("/getAllFaultInfoList4App", _AllowCrossDomain, devctrl.GetAllFaultInfoList4App) // 获取所有设备故障信息

	d.POST("/getDeviceQty", _AllowCrossDomain, devctrl.GetDeviceQty)           //设备分析-获取设备数量
	d.POST("/getUseTimeByModel", _AllowCrossDomain, devctrl.GetUseTimeByModel) //设备分析-按设备型号统计使用时间
	d.POST("/getUseTimeBySite", _AllowCrossDomain, devctrl.GetUseTimeBySite)   //设备分析-按设备位置统计使用时间

	d.POST("/getDeviceModelTree", _AllowCrossDomain, devctrl.GetDeviceModelTree) //获取设备型树

	d.POST("/getFault", _AllowCrossDomain, devctrl.GetFault)               //故障管理-获取故障记录
	d.POST("/registerFault", _AllowCrossDomain, devctrl.RegisterFault)     //故障管理-故障登记
	d.POST("/submitFault", _AllowCrossDomain, devctrl.SubmitFault)         //故障管理-故障提交
	d.POST("/deleteFault", _AllowCrossDomain, devctrl.DeleteFault)         //故障管理-故障删除
	d.POST("/acceptanceFault", _AllowCrossDomain, devctrl.AcceptanceFault) //故障管理-故障受理
	d.POST("/registerRepair", _AllowCrossDomain, devctrl.RegisterRepair)   //故障管理-维修登记

	d.POST("/getClassroomDevice", _AllowCrossDomain, devctrl.GetClassroomDevice)       //获取教室设备
	d.POST("/getDeviceAllFaultType", _AllowCrossDomain, devctrl.GetDeviceAllFaultType) //获取设备对应型号的所有故障故障分类
	d.POST("/getDevicFaultWord", _AllowCrossDomain, devctrl.GetDevicFaultWord)         //获取设备对应型号的故障现象词条

	d.POST("/getNodeModelList", _AllowCrossDomain, devctrl.GetNodeModelList) //节点配置-获取节点型号列表
	d.POST("/getNodeModel", _AllowCrossDomain, devctrl.GetNodeModel)         //节点配置-获取节点型号
	d.POST("/saveNodeModel", _AllowCrossDomain, devctrl.SaveNodeModel)       //节点配置-保存节点型号
	d.POST("/deleteNodeModel", _AllowCrossDomain, devctrl.DeleteNodeModel)   //节点配置-删除节点型号
	d.POST("/onDeletingNodeModel", _AllowCrossDomain, devctrl.OnDeletingNodeModel)

	d.POST("/getNodeModelCMDList", _AllowCrossDomain, devctrl.GetNodeModelCMDList) //节点配置-获取节点型号命令列表
	d.POST("/getNodeModelCMD", _AllowCrossDomain, devctrl.GetNodeModelCMD)         //节点配置-获取节点型号命令
	d.POST("/saveNodeModelCMD", _AllowCrossDomain, devctrl.SaveNodeModelCMD)       //节点配置-保存节点型号命令
	d.POST("/deleteNodeModelCMD", _AllowCrossDomain, devctrl.DeleteNodeModelCMD)   //节点配置-删除节点型号命令
	d.POST("/onDeletingNodeModelCmd", _AllowCrossDomain, devctrl.OnDeletingNodeModelCmd)

	d.POST("/getNodeList", _AllowCrossDomain, devctrl.GetNodeList) //节点配置-获取节点列表
	d.POST("/getNode", _AllowCrossDomain, devctrl.GetNode)         //节点配置-获取节点
	d.POST("/saveNode", _AllowCrossDomain, devctrl.SaveNode)       //节点配置-保存节点
	d.POST("/deleteNode", _AllowCrossDomain, devctrl.DeleteNode)   //节点配置-删除节点
	d.POST("/onDeletingNode", _AllowCrossDomain, devctrl.OnDeletingNode)

	d.POST("/getDeviceModelList", _AllowCrossDomain, devctrl.GetDeviceModelList) //设备配置-获取设备类型列表
	d.POST("/getDeviceModel", _AllowCrossDomain, devctrl.GetDeviceModel)         //设备配置-获取设备类型
	d.POST("/saveDeviceModel", _AllowCrossDomain, devctrl.SaveDeviceModel)       //设备配置-保存设备类型
	d.POST("/deleteDeviceModel", _AllowCrossDomain, devctrl.DeleteDeviceModel)   //设备配置-删除设备类型
	d.POST("/onDeletingDeviceModel", _AllowCrossDomain, devctrl.OnDeletingDeviceModel)

	/* 设备配置 设备状态管理*/
	d.POST("/getDeviceModelStatusCMDList", _AllowCrossDomain, devctrl.GetDeviceModelStatusCMDList) //设备配置-获取设备类型状态命令列表
	d.POST("/getDeviceModelStatusCMD", _AllowCrossDomain, devctrl.GetDeviceModelStatusCMD)         //设备配置-获取设备类型状态命令
	d.POST("/saveDeviceModelStatusCMD", _AllowCrossDomain, devctrl.SaveDeviceModelStatusCMD)       //设备配置-保存设备类型状态命令
	d.POST("/deleteDeviceModelStatusCMD", _AllowCrossDomain, devctrl.DeleteDeviceModelStatusCMD)   //设备配置-删除设备类型状态命令
	d.POST("/onDeletingDeviceModelStatusCmd", _AllowCrossDomain, devctrl.OnDeletingDeviceModelStatusCmd)

	d.POST("/getDeviceModelStatusValueCodeList", _AllowCrossDomain, devctrl.GetDeviceModelStatusValueCodeList) //设备配置-获取设备类型状态值编码列表
	d.POST("/getDeviceModelStatusValueCode", _AllowCrossDomain, devctrl.GetDeviceModelStatusValueCode)         //设备配置-获取设备类型状态值编码
	d.POST("/saveDeviceModelStatusValueCode", _AllowCrossDomain, devctrl.SaveDeviceModelStatusValueCode)       //设备配置-保存设备类型状态值编码
	d.POST("/deleteDeviceModelStatusValueCode", _AllowCrossDomain, devctrl.DeleteDeviceModelStatusValueCode)   //设备配置-删除设备类型状态值编码
	/* 设备配置 设备状态管理*/

	/* 设备配置 控制命令管理*/
	d.POST("/getDeviceModelControlCMDList", _AllowCrossDomain, devctrl.GetDeviceModelControlCMDList) //设备配置-获取设备型号控制命令列表
	d.POST("/getDeviceModelControlCMD", _AllowCrossDomain, devctrl.GetDeviceModelControlCMD)         //设备配置-获取设备型号控制命令
	d.POST("/saveDeviceModelControlCMD", _AllowCrossDomain, devctrl.SaveDeviceModelControlCMD)       //设备配置-保存设备型号控制命令
	d.POST("/deleteDeviceModelControlCMD", _AllowCrossDomain, devctrl.DeleteDeviceModelControlCMD)   //设备配置-删除设备型号控制命令
	d.POST("/onDeletingDeviceModelControlCmd", _AllowCrossDomain, devctrl.OnDeletingDeviceModelControlCmd)
	/* 设备配置 控制命令管理*/

	/* 设备配置 设备管理*/
	d.POST("/getDeviceList", _AllowCrossDomain, devctrl.GetDeviceList) //设备配置-获取设备列表
	d.POST("/getDevice", _AllowCrossDomain, devctrl.GetDevice)         //设备配置-获取设备
	d.POST("/saveDevice", _AllowCrossDomain, devctrl.SaveDevice)       //设备配置-保存设备
	d.POST("/deleteDevice", _AllowCrossDomain, devctrl.DeleteDevice)   //设备配置-删除设备
	d.POST("/onDeletingDevice", _AllowCrossDomain, devctrl.OnDeletingDevice)
	/* 设备配置 设备管理*/

	/* 设备配置 设备型号故障分类*/
	d.POST("/getDeviceModelFaultTypeList", _AllowCrossDomain, devctrl.GetDeviceModelFaultTypeList) //设备配置-故障分类列表
	d.POST("/getDeviceModelFaultType", _AllowCrossDomain, devctrl.GetDeviceModelFaultType)         //设备配置-获取故障分类
	d.POST("/saveDeviceModelFaultType", _AllowCrossDomain, devctrl.SaveDeviceModelFaultType)       //设备配置-保存故障分类
	d.POST("/deleteDeviceModelFaultType", _AllowCrossDomain, devctrl.DeleteDeviceModelFaultType)   //设备配置-删除故障分类
	d.POST("/onDeletingDeviceModelFaultType", _AllowCrossDomain, devctrl.OnDeletingDeviceModelFaultType)
	/* 设备配置 设备型号故障分类*/

	/* 设备配置 设备型号故障现象常用词条*/
	d.POST("/getDeviceModelFaultWordList", _AllowCrossDomain, devctrl.GetDeviceModelFaultWordList) //设备配置-获取故障现象常用词条列表
	d.POST("/getDeviceModelFaultWord", _AllowCrossDomain, devctrl.GetDeviceModelFaultWord)         //设备配置-获取故障现象常用词条
	d.POST("/saveDeviceModelFaultWord", _AllowCrossDomain, devctrl.SaveDeviceModelFaultWord)       //设备配置-保存故障现象常用词条
	d.POST("/deleteDeviceModelFaultWord", _AllowCrossDomain, devctrl.DeleteDeviceModelFaultWord)   //设备配置-删除故障现象常用词条
	/* 设备配置 设备型号故障现象常用词条*/

	r.POST("/device/node/deviceinfo/detail/query", _AllowCrossDomain, devctrl.GetNode_DeviceDetailInfo_Handler) // 查询节点下设备详细信息
	r.POST("/device/node_device/basicinfo/query", _AllowCrossDomain, devctrl.Get_NodeDevBasicInfo_Handler)      // 查询节点&设备基本信息
	r.POST("/device/node/basicinfo/query", _AllowCrossDomain, devctrl.Get_NodeBasicInfo_Handler)                // 节点基本信息查询
	r.POST("/device/device/basicinfo/query", _AllowCrossDomain, devctrl.Get_DeviceBasicInfo_Handler)            // 设备基本信息查询
	r.POST("/device/management/device/unbind", _AllowCrossDomain, devctrl.Do_DeviceUnbind_Handler)              // 设备解绑定
	r.POST("/device/node/detailinfo/query", _AllowCrossDomain, devctrl.Query_NodeDetail_Handler)                // 查询节点详细信息列表
	r.POST("/device/device/detailinfo/query", _AllowCrossDomain, devctrl.Query_DeviceDetail_Handler)            // 查询设备详细信息列表
	r.POST("/device/device/detailinfo/query_valid", _AllowCrossDomain, devctrl.QueryValid_DeviceDetail_Handler) // 查询【有效】设备详细信息列表
	r.POST("/device/node/detailinfo/queryinfo", _AllowCrossDomain, devctrl.Query_NodeDetail_Info_Handler)       // 查询节点具体详细信息
	r.POST("/device/device/detailinfo/queryinfo", _AllowCrossDomain, devctrl.Query_DeviceDetail_Info_Handler)   // 查询设备具体详细信息
	r.POST("/device/cmd/query_byroom", _AllowCrossDomain, devctrl.GetDeviceCmd_ByRoom)                          // 获取房间内的所有设备命令

	r.POST("/upload", func(c *gin.Context) {
		//		defer xerr.CatchPanic()
		//		r := c.Request

		//		log.Printf("%v\n", r)
		//		log.Println(r.Header)
		//		log.Println(r.Form)
		//		log.Println(c.Params)
		//		log.Println(c.Query("Filename"))
		//		file_length := r.Header.Get("Content-Length")
		//		log.Println(file_length)

		//		obj := struct{ File, FileSize string }{"file.up", file_length}
		//		c.JSON(http.StatusOK, obj)
	})
}

/*
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

	log.Printf("%+v", c.Request)
	log.Println(c.Request.FormValue("filename"))
	xerr.ThrowPanic(err)
	if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
		for k, pFileArr := range c.Request.MultipartForm.File {
			log.Println(k)
			if len(pFileArr) > 0 {
				var inFile multipart.File
				inFile, err = pFileArr[0].Open()
				xerr.ThrowPanic(err)
				defer inFile.Close()

				curdir := core.GetCurrentDirectory() //获取当前运行的目录
				curdir += "/../../FrontEndCode/build/FrontEndCode/templates/upfile"
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

				rd.Rcode = "1002"
				rd.Rcode = "操作失败!"
				_, err = io.Copy(outFile, inFile)
				xerr.ThrowPanic(err)

				rd.Rcode = "1000"
				rd.Reason = "操作成功!"
				rd.Result = "/web/upfile/" + timestr + "/" + filename
			}
		}
	}
}
*/
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
