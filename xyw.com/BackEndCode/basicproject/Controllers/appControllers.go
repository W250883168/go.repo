package appControllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"xutils/xtext"

	//	"basicproject/DataAccess/actiondataDataAccess"
	//	"basicproject/DataAccess/basicsetDataAccess"
	//	"basicproject/DataAccess/curriculumDataAccess"
	//	"basicproject/DataAccess/liveDataAccess"
	"basicproject/DataAccess/usersDataAccess"
	//	"basicproject/model/actiondata"
	//	"basicproject/model/basicset"
	//	"basicproject/model/curriculum"
	"basicproject/model/systemmodel"
	"basicproject/model/users"
	"basicproject/viewmodel"
	core "xutils/xcore"

	"github.com/gin-gonic/gin"
	//	"github.com/tealeg/xlsx"

	"net"
)

func GetLogin(c *gin.Context) { //登录方法
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Content-Type", "application/json")
	var rd core.Returndata
	var lg viewmodel.Login
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	fmt.Printf("%+v \n", lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	if errs1 == nil {
		//		if lg.Os != "" && (lg.Os == "IOS" || lg.Os == "Android" || lg.Os == "Controlpanel" || lg.Os == "PcWEB") {
		ThirdPartyId := lg.ThirdPartyId
		Os := lg.Os
		sql := "select Id as Usersid,Loginuser,Rolesid as Rolestype,Truename,Nickname,Userheadimg,Userphone,Userstate,Usersex,Usermac,Birthday from users where Loginuser=? and Loginpwd=?;"
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		lg.Loginpwd = xtext.SpaceTrim(lg.Loginpwd)
		lg.Loginuser = xtext.SpaceTrim(lg.Loginuser)
		errs2 := dbmap.SelectOne(&lg, sql, lg.Loginuser, lg.Loginpwd)
		core.CheckErr(errs2, "appControllers|GetLogin|登录方法中根据账号密码查询错误")
		if errs2 == nil {
			fmt.Println(lg.Loginuser + "|" + strconv.Itoa(lg.Rolestype) + "|" + strconv.Itoa(lg.Usersid))
			//tkbt, _ := core.RsaEncrypt([]byte(lg.Loginuser + "|" + strconv.Itoa(lg.Rolestype) + "|" + strconv.Itoa(lg.Usersid)))
			lg.Token = lg.Loginuser + "|" + strconv.Itoa(lg.Rolestype) + "|" + strconv.Itoa(lg.Usersid) //string(tkbt)
			//			fmt.Println(ThirdPartyId)
			if ThirdPartyId != "" {
				fmt.Println(ThirdPartyId)
				_, errs2 = dbmap.Exec("update users set ThirdPartyId='' where ThirdPartyId='" + ThirdPartyId + "'")
				fmt.Println(errs2)
				_, errs2 = dbmap.Exec("update users set ThirdPartyId='" + ThirdPartyId + "' where Id=" + strconv.Itoa(lg.Usersid))
				fmt.Println(errs2)
			}
			if Os != "" {
				_, errs2 = dbmap.Exec("update users set Os='" + Os + "' where Id=" + strconv.Itoa(lg.Usersid))
				fmt.Println(errs2)
			}
			rd.Result = lg
			rd.Rcode = "1000"
			rd.Reason = "成功"

			if Os == "Controlpanel" || Os == "PcWEB" { //来源控制面板
				if lg.Rolestype == 3 {
					rd.Result = nil
					rd.Rcode = "1005"
					rd.Reason = "您权限不够，不能登录"
				}
			}

			// 记录登录成功日志
			defer func() {
				var log users.LoginLog
				log.LoginOS = Os
				log.LoginDate = time.Now().String()
				log.LoginUsers = lg.Loginuser
				log.LoginState = 1
				log.IP = GetLocalIPAddr()
				// fmt.Println(log)
				dbmap.AddTableWithName(users.LoginLog{}, "loginlog").SetKeys(true, "Id")
				usersDataAccess.AddLoginlog(&log, dbmap)
			}()
		} else {
			rd.Rcode = "1001"
			rd.Reason = "账户或者密码错误"
		}
		//		} else {
		//			rd.Rcode = "1002"
		//			rd.Reason = "数据提交错误"
		//		}
	} else {
		rd.Rcode = "1003"
		rd.Reason = "提交数据格式不正确"
	}
	c.JSON(200, rd)
}

func SystemGroupFunc(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "appControllers|SystemGroupFunc|登录数据json格式转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			rd.Rcode = "1000"
			c.Set("users", lg)
			c.Set("data", string(data))
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	if rd.Rcode != "1000" {
		c.JSON(200, rd)
	}
}

/*
SELECT
	sm.Modulename,
	sm.Moduledisplayname,
	sm.Moduledisplayterminal,
	sm.Modulecode,
	sm.Moduleicon,
	sm.Superiormoduleid,
	sm.Id
FROM systemmodule AS sm
		INNER JOIN rolemodulecenter rmc ON sm.Id = rmc.Systemmoduleid
		INNER JOIN roles AS rl ON rmc.Rolesid = rl.Id
WHERE rl.Id = @RoleTypeID AND find_in_set(@LoginOS, sm.Moduledisplayterminal) > 0
GROUP BY sm.Modulename, sm.Modulecode, sm.Moduleicon
ORDER BY sm.ModuleIndex;
*/
func GetHomeApp(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken //viewmodel.Getapp
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "appControllers|GetHomeApp|登录数据json格式转换错误")
	if errs1 == nil {
		tkbt := lg.Token //core.RsaDecrypt([]byte(lg.Token))
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			sql := "select sm.Modulename,sm.Moduledisplayname,sm.Moduledisplayterminal,sm.Modulecode,sm.Moduleicon,sm.Superiormoduleid,sm.Id from systemmodule as sm inner join rolemodulecenter rmc on sm.Id=rmc.Systemmoduleid inner join roles as rl on rmc.Rolesid=rl.Id"
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			var ssarr []systemmodel.Systemmodule
			sql = sql + " where rl.Id=? and find_in_set('" + lg.Os + "',sm.Moduledisplayterminal)>0 group by sm.Modulename,sm.Modulecode,sm.Moduleicon order by sm.ModuleIndex;"
			_, errs2 := dbmap.Select(&ssarr, sql, lg.Rolestype) //查询权限模块
			core.CheckErr(errs2, "appControllers|GetHomeApp|登录方法中根据账号密码查询错误")
			if errs2 == nil {
				rd.Result = ssarr
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "没有查询到此角色的权限"
			}
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}

//获取服务器的时间戳
func GetServerTime(c *gin.Context) {
	var rd core.Returndata
	rd.Rcode = "1000"
	t := time.Now().Unix()
	//	fmt.Println(t)
	rd.Result = t
	//	//时间戳到具体显示的转化
	//	fmt.Println(time.Unix(t, 0).String())
	//	//带纳秒的时间戳
	//	t = time.Now().UnixNano()
	//	fmt.Println(t)
	//	fmt.Println("------------------")
	//	//基本格式化的时间表示
	//	fmt.Println(time.Now().String())
	//	fmt.Println(time.Now().Format("2006year 01month 02day"))
	c.JSON(200, rd)
}

/*
		初始化begin
		主要业务：导入学科分类，导入课程信息，导入章节信息，导入课程班级信息，导入课程班级章节信息，导入课程班级章节点到数据
		初始化end

func Insertdata(c *gin.Context) {

	dbmap2 := core.InitDb() //InitZndxDb() //wifi定位数据分析后
	defer dbmap2.Db.Close()
	excelFileName := "F:/textdata8-4.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic(err)
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Println(sheet.Name)
		switch sheet.Name {
		case "模块表":
			break
		case "功能表":
			break
		case "角色表":
			break
		case "角色模块中间表":
			break
		case "角色模块功能中间表":
			break
		case "校区表":
			if len(sheet.Rows) > 2 {
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 3 {
						cas := basicset.Campus{Campusname: row.Cells[0].Value, Campusicon: row.Cells[1].Value, Campuscode: row.Cells[2].Value, Campusnums: 0}
						prterr := basicsetDataAccess.AddCampus(&cas, dbmap2)
						if prterr != nil {
							fmt.Println("添加校区数据", prterr)
						}
					}
				}
			}
			break
		case "学院表":
			if len(sheet.Rows) > 2 {
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 4 {
						Campusidstr := row.Cells[0].Value
						ws := viewmodel.QueryBasicsetWhere{Campuscode: Campusidstr}
						Campuslist := basicsetDataAccess.QueryCampus(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap2)
						var Campusidint int
						if len(Campuslist) == 1 {
							Campusidint = Campuslist[0].Campusid
						}
						cas := basicset.College{Campusid: Campusidint, Collegecode: row.Cells[1].Value, Collegename: row.Cells[2].Value, Collegeicon: row.Cells[3].Value, Collegenum: 0}
						basicsetDataAccess.AddCollege(&cas, dbmap2)
					}
				}
			}
			break
		case "楼栋表":
			if len(sheet.Rows) > 2 {
				fmt.Println("进入到添加楼栋表数据中。。。")
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 3 {
						fmt.Println("开始添加楼栋表数据。。。")
						Campusidstr := row.Cells[0].Value
						ws := viewmodel.QueryBasicsetWhere{Campuscode: Campusidstr}
						Campuslist := basicsetDataAccess.QueryCampus(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap2)
						var Campusidint int
						if len(Campuslist) == 1 {
							Campusidint = Campuslist[0].Campusid
						}
						fmt.Println("添加楼栋数据时，获取的校区数据Campusidint", Campusidint)
						Buildingiconstr := strings.Trim(row.Cells[5].Value, " ")
						cas := basicset.Building{Campusid: Campusidint, Buildingname: row.Cells[1].Value, Buildingcode: row.Cells[2].Value, Buildingicon: Buildingiconstr, Floorsnumber: 0, Classroomsnumber: 0}
						inerr := basicsetDataAccess.AddBuilding(&cas, dbmap2)
						if inerr != nil {
							fmt.Println("添加楼栋数据", inerr)
						}
					}
				}
			}
			break
		case "楼层表":
			if len(sheet.Rows) > 2 {
				fmt.Println("进入到添加楼层表数据中。。。")
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 3 {
						fmt.Println("开始添加楼层表数据。。。")
						Buildingstr := row.Cells[0].Value
						ws := viewmodel.QueryBasicsetWhere{Buildingcode: Buildingstr}
						Buildinglist := basicsetDataAccess.QueryBuilding(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap2)
						var Buildingid int
						if len(Buildinglist) == 1 {
							Buildingid = Buildinglist[0].Buildingid
						}
						cas := basicset.Floors{Buildingid: Buildingid, Floorname: row.Cells[1].Value, Floorscode: row.Cells[2].Value}
						inerr := basicsetDataAccess.AddFloors(&cas, dbmap2)
						if inerr != nil {
							fmt.Println("添加楼层数据", inerr)
						}
					}
				}
			}
			break
		case "教室表":
			if len(sheet.Rows) > 2 {
				fmt.Println("进入到添加教室表数据中。。。")
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 5 {
						fmt.Println("开始添加教室表数据。。。")
						Floorsstr := row.Cells[0].Value
						ws := viewmodel.QueryBasicsetWhere{Floorscode: Floorsstr}
						Floorslist := basicsetDataAccess.QueryFloors(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap2)
						var Floorsidint int
						if len(Floorslist) == 1 {
							Floorsidint = Floorslist[0].Floorsid
						}
						Seatsnumbersint, _ := strconv.Atoi(row.Cells[3].Value)
						ssarr := strings.Split(row.Cells[1].Value, "-")
						var Classroomsnamestr string
						if len(ssarr) == 3 {
							Classroomsnamestr = strings.Replace(ssarr[1], "栋", "", 0)
							Classroomsnamestr = Classroomsnamestr + ssarr[2]
						}
						cas := basicset.Classrooms{Floorsid: Floorsidint, Classroomsname: Classroomsnamestr, Classroomscode: row.Cells[1].Value, Classroomicon: row.Cells[2].Value, Seatsnumbers: int(Seatsnumbersint), Classroomstype: row.Cells[4].Value, Notes: row.Cells[6].Value}
						inerr := basicsetDataAccess.AddClassrooms(&cas, dbmap2)
						if inerr != nil {
							fmt.Println("添加教室表数据", inerr)
						}
					}
				}
			}
			break
		case "科系专业表":
			if len(sheet.Rows) > 2 {
				fmt.Println("进入到添加科系专业表数据中。。。")
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 3 {
						fmt.Println("开始添加科系专业表数据。。。")
						Collegestr := row.Cells[0].Value
						ws := viewmodel.QueryBasicsetWhere{Collegecode: Collegestr}
						Collegelist := basicsetDataAccess.QueryCollege(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap2)
						var Collegeidint int
						if len(Collegelist) == 1 {
							Collegeidint = Collegelist[0].Id
						}
						Majoricon := ""
						if len(row.Cells) >= 4 {
							Majoricon = row.Cells[3].Value
						}
						cas := basicset.Major{Collegeid: Collegeidint, Majorcode: row.Cells[1].Value, Majorname: row.Cells[2].Value, Majoricon: Majoricon, Majornum: 0}
						inerr := basicsetDataAccess.AddMajor(&cas, dbmap2)
						if inerr != nil {
							fmt.Println("添加科系专业表数据", inerr)
						}
					}
				}
			}
			break
		case "班级表":
			if len(sheet.Rows) > 2 {
				fmt.Println("进入到添加班级数据中。。。")
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 4 {
						fmt.Println("开始添加班级数据。。。")
						Majorstr := row.Cells[2].Value
						ws := viewmodel.QueryBasicsetWhere{Majorcode: Majorstr}
						Majorlist := basicsetDataAccess.QueryMajor(ws, core.PageData{PageIndex: 0, PageSize: 10000}, dbmap2)
						var Majorid int
						if len(Majorlist) == 1 {
							Majorid = Majorlist[0].Id
						}
						classericon := ""
						if len(row.Cells) >= 5 {
							classericon = row.Cells[4].Value
						}
						cas := basicset.Classes{Majorid: Majorid, Classesname: row.Cells[0].Value, Classesnum: 0, Classesicon: classericon, Classstate: 0}
						inerr := basicsetDataAccess.AddClasses(&cas, dbmap2)
						if inerr != nil {
							fmt.Println("添加班级表数据", inerr)
						}
					}
				}
			}
			break
		case "人员信息表":
			if len(sheet.Rows) > 2 {
				for rls, row := range sheet.Rows {
					if rls == 0 {
						continue
					}
					if len(row.Cells) >= 11 {
						Rolesstr := row.Cells[4].Value
						var Usersexint int
						if row.Cells[2].Value == "男" {
							Usersexint = 0
						} else if row.Cells[2].Value == "女" {
							Usersexint = 1
						} else {
							Usersexint = 2
						}
						userheadimg := ""
						if len(row.Cells) >= 19 {
							if row.Cells[18].Value != "" {
								userheadimg = row.Cells[18].Value
							}
						}
						rid, _ := dbmap2.SelectInt("select Id from roles where Rolesname=?", Rolesstr)
						us := users.Users{Loginuser: row.Cells[1].Value, Loginpwd: row.Cells[2].Value, Rolesid: int(rid), Truename: row.Cells[5].Value, Nickname: row.Cells[0].Value, Userheadimg: userheadimg, Userphone: "", Userstate: 0, Usersex: Usersexint, Birthday: row.Cells[8].Value}
						var st users.Students
						if rid == 3 {
							Enrollmentyearint, _ := strconv.Atoi(row.Cells[9].Value)
							Classstr := row.Cells[12].Value
							classid, _ := dbmap2.SelectInt("select Id from classes where Classesname=?;", Classstr)
							st = users.Students{Enrollmentyear: int(Enrollmentyearint), Homeaddress: row.Cells[10].Value, Nowaddress: row.Cells[11].Value, Classesid: int(classid)}
						}
						usersDataAccess.AddUsers(&us, &st, dbmap2)
					}
				}
			}
			break
		case "学科表":
			for rls, row := range sheet.Rows {
				if rls == 0 {
					continue
				}
				cell1 := row.Cells[0].Value
				cell2 := row.Cells[1].Value
				var sc curriculum.Subjectclass
				if len(cell1) > 2 { //有上一级
					sc.Subjectcode = cell1
					sc.Subjectname = cell2
					sc.Superiorsubjectcode = core.Substr(cell1, 0, len(cell1)-2)
				} else { //无上一级
					sc.Subjectcode = cell1
					sc.Subjectname = cell2
				}
				//				curriculumDataAccess.AddSubjectclass(&sc, dbmap2)
			}
			break
		case "班级课程详细表":
			var cc curriculum.Curriculums              //当前课程
			var cccc curriculum.Curriculumsclasscentre //课程班级中间表
			for rls, row := range sheet.Rows {
				if rls == 0 {
					continue
				}
				if len(row.Cells) < 18 {
					continue
				}
				cell1 := row.Cells[0].Value
				cell12 := row.Cells[11].Value
				usersid, _ := dbmap2.SelectInt("select Id from users where Loginuser=?;", cell12)                       //获取老师ID
				classid, _ := dbmap2.SelectInt("select Id from classes where Classesname=?;", row.Cells[12].Value)      //班级ID
				roomid, _ := dbmap2.SelectInt("select Id from Classrooms where Classroomscode=?;", row.Cells[13].Value) //教室ID
				Cell15 := row.Cells[14].Value                                                                           //开始时间
				Cell16 := row.Cells[15].Value                                                                           //结束时间
				Cell17 := row.Cells[16].Value                                                                           //是否允许直播
				Islive, _ := strconv.Atoi(Cell17)
				Cell18 := row.Cells[17].Value //是否允许录播
				Isondemand, _ := strconv.Atoi(Cell18)
				if cell1 != "" { //添加课程记录
					cell2 := row.Cells[1].Value
					cell3 := row.Cells[2].Value
					cell4 := row.Cells[3].Value
					cell5 := row.Cells[4].Value
					cell6 := row.Cells[5].Value
					Chaptercount, _ := strconv.Atoi(cell6)
					//cell7 := row.Cells[6].Value
					cell8 := row.Cells[7].Value
					//cell9 := row.Cells[8].Value
					cc.Curriculumname = cell1                                //课程名称
					cc.Curriculumicon = cell2                                //课程图标
					cc.Curriculumnature = cell3                              //课程性质
					cc.Curriculumstype = cell4                               //课程类型
					cc.Curriculumsdetails = cell5                            //课程详情
					cc.Chaptercount = Chaptercount                           //总章节数
					cc.Subjectcode = cell8                                   //学科分类
					cc.Createdate = time.Now().Format("2006-01-02 15:04:05") //记录时间

					//					curriculumDataAccess.AddCurriculums(&cc, dbmap2)
					cccc.Curriculumsid = cc.Id
					cccc.Classesid = int(classid)
					cccc.Usersid = int(usersid)
					cccc.Createdate = time.Now().Format("2006-01-02 15:04:05")
					cccc.Newchapter = 0
					cccc.Islive = Islive
					cccc.Isondemand = Isondemand
					cccc.Whenlongcount = 0

					//					curriculumDataAccess.AddCurriculumsClassCentre(&cccc, dbmap2)
					//dbmap2.Insert(&cccc) //添加班级课程中间表
				}
				var ct curriculum.Chapters //添加章节
				//				var lv live.Lives          //添加直播记录表
				cell9 := row.Cells[8].Value
				cell10 := row.Cells[9].Value
				cell11 := row.Cells[10].Value
				ct.Chaptername = cell9
				ct.Chaptericon = cell10
				ct.Curriculumsid = cc.Id
				ct.Chapterdetails = cell11
				ct.Createdate = time.Now().Format("2006-01-02 15:04:05") //记录时间
				//				curriculumDataAccess.AddChapters(&ct, dbmap2)

				var sdarr []users.Students
				dbmap2.Select(&sdarr, "select * from students where Classid=?;", classid)
				var ccccct curriculum.Curriculumclassroomchaptercentre //课程班级章节中间表
				ccccct.Curriculumsclasscentreid = cccc.Id
				ccccct.Chaptersid = ct.Id
				ccccct.Usersid = int(usersid)
				ccccct.Createdate = time.Now().Format("2006-01-02 15:04:05")
				ccccct.Begindate = core.Timeaction(Cell15)
				ccccct.Enddate = core.Timeaction(Cell16)
				ccccct.Isondomian = Isondemand
				ccccct.Islive = Islive
				ccccct.Whenlong = 0
				ccccct.Plannumber = len(sdarr)
				//				curriculumDataAccess.AddCurriculumclassroomchaptercentre(&ccccct, 0, dbmap2)
				//dbmap2.Insert(&ccccct)            //课程班级章节中间表
				for p := 0; p < len(sdarr); p++ { //循环插入预置点到信息记录
					pt := actiondata.Pointtos{Curriculumclassroomchaptercentreid: ccccct.Id, Usersid: sdarr[p].Id, State: 0}
					actiondataDataAccess.AddPointtos(&pt, dbmap2)
					//dbmap2.Insert(&pt)
				}
				//插入教室授课记录表
				tr := actiondata.Teachingrecord{Classroomid: int(roomid), Curriculumclassroomchaptercentreid: ccccct.Id, State: 0}
				actiondataDataAccess.AddTeachingrecord(&tr, dbmap2)
				//				if ccccct.Islive == 1 || ccccct.Isondomian == 1 { //判断是直播或者是录播
				//					lv.Begindate = ccccct.Begindate
				//					lv.Curriculumclassroomchaptercentreid = ccccct.Id
				//					lv.Coverimage = ""
				//					if len(row.Cells) >= 28 {
				//						lv.Coverimage = row.Cells[27].Value
				//					}
				//					if len(row.Cells) >= 25 {
				//						lv.Livepath1 = row.Cells[24].Value
				//					}
				//					if len(row.Cells) >= 29 {
				//						Whenlong, wlerr := strconv.Atoi(row.Cells[28].Value)
				//						fmt.Println(wlerr)
				//						lv.Whenlong = Whenlong
				//					}
				//					lv.Liveinfo = row.Cells[19].Value      //播放内容简介
				//					lv.Recommendread = row.Cells[20].Value //推荐阅读
				//					lv.Livetitile = row.Cells[21].Value    //播放标题
				//					lv.Iscomment = 0                       //是否允许评论
				//					lv.Ischeckcomment = 0                  //是否审核评论
				//					lv.Isdownload = 0                      //是否允许下载
				//					if len(row.Cells) >= 23 {
				//						Iscomment, _ := strconv.Atoi(row.Cells[22].Value)
				//						lv.Iscomment = Iscomment
				//					}
				//					if len(row.Cells) >= 24 {
				//						Ischeckcomment, _ := strconv.Atoi(row.Cells[23].Value)
				//						lv.Ischeckcomment = Ischeckcomment //是否审核评论
				//					}
				//					if len(row.Cells) >= 19 {
				//						Isdownload, _ := strconv.Atoi(row.Cells[18].Value)
				//						lv.Isdownload = Isdownload //是否审核评论
				//					}
				//Isdownload, _ := strconv.Atoi(row.Cells[18].Value)
				//lv.Isdownload = Isdownload //是否允许下载
				//					liveDataAccess.AddLives(&lv, dbmap2)
				//				}
			}
			break

		}
	}
	c.String(200, "Ok")
}
*/
/*
func InitData(c *gin.Context) {

	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	smarr := make([]systemmodel.Systemmodule, 8)
	smarr[0] = systemmodel.Systemmodule{Modulename: "课程管理", Modulecode: "kcgl", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[1] = systemmodel.Systemmodule{Modulename: "远程教学", Modulecode: "ycjx", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[2] = systemmodel.Systemmodule{Modulename: "教室导流", Modulecode: "jsdl", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[3] = systemmodel.Systemmodule{Modulename: "查课表", Modulecode: "ckb", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[4] = systemmodel.Systemmodule{Modulename: "视频录播", Modulecode: "splb", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[5] = systemmodel.Systemmodule{Modulename: "设备管理", Modulecode: "sbgl", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[6] = systemmodel.Systemmodule{Modulename: "上课点到", Modulecode: "skdd", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	smarr[7] = systemmodel.Systemmodule{Modulename: "出勤统计", Modulecode: "cqtj", Moduleicon: "", Superiormoduleid: 0, Moduleurl: "", Moduleattribute: ""}
	dbmap.AddTableWithName(systemmodel.Systemmodule{}, "systemmodule").SetKeys(true, "Id")
	err1 := dbmap.Insert(&smarr[0], &smarr[1], &smarr[2], &smarr[3], &smarr[4], &smarr[5], &smarr[6], &smarr[7])
	core.CheckErr(err1, "静态模块数据设置失败")
	rlarr := make([]systemmodel.Roles, 4)
	rlarr[0] = systemmodel.Roles{Rolesname: "管理人员"}
	rlarr[1] = systemmodel.Roles{Rolesname: "教师人员"}
	rlarr[2] = systemmodel.Roles{Rolesname: "学生"}
	rlarr[3] = systemmodel.Roles{Rolesname: "测试人员"}
	dbmap.AddTableWithName(systemmodel.Roles{}, "Roles").SetKeys(true, "Id")
	err3 := dbmap.Insert(&rlarr[0], &rlarr[1], &rlarr[2], &rlarr[3])
	core.CheckErr(err3, "静态角色数据设置失败")
	rmcarr := make([]systemmodel.Rolemodulecenter, 27)
	rmcarr[0] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[0].Id}
	rmcarr[1] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[1].Id}
	rmcarr[2] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[2].Id}
	rmcarr[3] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[3].Id}
	rmcarr[4] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[4].Id}
	rmcarr[5] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[5].Id}
	rmcarr[6] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[6].Id}
	rmcarr[7] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[7].Id}
	rmcarr[8] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[0].Id}
	rmcarr[9] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[1].Id}
	rmcarr[10] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[3].Id}
	rmcarr[11] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[4].Id}
	rmcarr[12] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[5].Id}
	rmcarr[13] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[6].Id}
	rmcarr[14] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[7].Id}
	rmcarr[15] = systemmodel.Rolemodulecenter{Rolesid: rlarr[2].Id, Systemmoduleid: smarr[1].Id}
	rmcarr[16] = systemmodel.Rolemodulecenter{Rolesid: rlarr[2].Id, Systemmoduleid: smarr[2].Id}
	rmcarr[17] = systemmodel.Rolemodulecenter{Rolesid: rlarr[2].Id, Systemmoduleid: smarr[3].Id}
	rmcarr[18] = systemmodel.Rolemodulecenter{Rolesid: rlarr[2].Id, Systemmoduleid: smarr[4].Id}
	rmcarr[19] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[0].Id}
	rmcarr[20] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[1].Id}
	rmcarr[21] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[2].Id}
	rmcarr[22] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[3].Id}
	rmcarr[23] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[4].Id}
	rmcarr[24] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[5].Id}
	rmcarr[25] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[6].Id}
	rmcarr[26] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[7].Id}
	dbmap.AddTableWithName(systemmodel.Rolemodulecenter{}, "rolemodulecenter").SetKeys(true, "Id")
	for p := 0; p < len(rmcarr); p++ {
		err3 = dbmap.Insert(&rmcarr[p])
		core.CheckErr(err3, "静态角色模块数据设置失败")
	}
	var smfarr []systemmodel.Systemmodulefunctions
	smfarr0 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[0].Id, Functionname: "查看课程", Functionicon: "", Functioncode: "queryall"}
	smfarr1 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[1].Id, Functionname: "远程教学列表", Functionicon: "", Functioncode: "ycjxlist"}
	smfarr2 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[2].Id, Functionname: "教室导流查看", Functionicon: "", Functioncode: "jsdlquery"}
	smfarr3 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[3].Id, Functionname: "查课表查看", Functionicon: "", Functioncode: "ckbquery"}
	smfarr4 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[4].Id, Functionname: "视频录播查看", Functionicon: "", Functioncode: "splbquery"}
	smfarr5 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[5].Id, Functionname: "设备管理列表", Functionicon: "", Functioncode: "sbgllist"}
	smfarr6 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[6].Id, Functionname: "上课点到查看", Functionicon: "", Functioncode: "skddquery"}
	smfarr7 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[7].Id, Functionname: "出勤统计查看", Functionicon: "", Functioncode: "cqtjquery"}
	smfarr = append(smfarr, smfarr0, smfarr1, smfarr2, smfarr3, smfarr4, smfarr5, smfarr6, smfarr7)
	dbmap.AddTableWithName(systemmodel.Systemmodulefunctions{}, "systemmodulefunctions").SetKeys(true, "Id")
	err6 := dbmap.Insert(&smfarr[0], &smfarr[1], &smfarr[2], &smfarr[3], &smfarr[4], &smfarr[5], &smfarr[6], &smfarr[7])
	core.CheckErr(err6, "静态功能数据设置失败")
	var rmfcarr []systemmodel.Rolemodulefunctioncenter
	var rmfc systemmodel.Rolemodulefunctioncenter
	dbmap.AddTableWithName(systemmodel.Rolemodulefunctioncenter{}, "rolemodulefunctioncenter").SetKeys(true, "Id")
	for i := 0; i < len(rmcarr); i++ {
		for j := 0; j < len(smfarr); j++ {
			if rmcarr[i].Systemmoduleid == smfarr[j].Systemmoduleid {
				rmfc = systemmodel.Rolemodulefunctioncenter{Rolemodulecenterid: rmcarr[i].Id, Systemmodulefunctionsid: smfarr[j].Id}
				err3 = dbmap.Insert(&rmfc)
				core.CheckErr(err3, "静态角色模块功能数据设置失败")
				rmfcarr = append(rmfcarr, rmfc)
			}
		}
	}

}
*/
func GetLocalIPAddr() string {
	var ip string = ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		// 检查ip地址判断是否回环地址
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// fmt.Println(ipnet.IP.String())
					ip = ipnet.IP.String()
				}
			}
		}
	}

	return ip
}

//func ImageFile(filetype int, filepath string) (newfilepath string) {
//	switch filetype {
//	case 0:
//		//		files, fileerr1 := os.Lstat(filepath)
//		//		if fileerr1 == nil {
//		//			files.IsDir()
//		//		}
//		break //校区图标
//	case 1:
//		break //楼栋图标
//	default: //默认
//		break
//	}
//}

//func InitData(c *gin.Context) {
//	//第一步添加设置数据
//	/*
//		模块表设置
//		功能表设置
//		角色表设置
//		角色模块中间表
//		角色模块功能中间表

//		1.1添加模块
//		1.2添加角色功能
//		1.3添加模块角色中间设置功能
//		1.4添加模块下功能数据
//		1.5添加模块角色中间表与功能关联数据

//		校区表
//		学院表
//		科系/专业表
//		班级表
//		楼栋表
//		楼层表
//		教室表
//		用户表
//		学生表
//		2.1添加校区数据
//		2.2添加学院表[需更变学院人数统计]
//		2.3添加科系/专业表[需更变科系人数统计]
//		2.4添加班级表[人数不是人工输入进入的]
//		2.5添加楼栋表[需要反写楼层数以及教室的数量]
//		2.6添加楼层表[需要反写教室的数量]
//		2.7添加教室表数据[需要反写教室内的人数，使用状态，收藏次数]
//		2.8添加用户表数据
//		2.9添加学生表数据
//		学科分类表
//		课程表
//		课程章节表
//		课程班级中间表
//		课程班级章节中间表
//		课程章节附件表
//	*/
//	//第二步测试设置数据的业务关联性
//	//第三步测试设置数据的调整
//	//第四步添加业务数据
//	//第五步测试业务数据的业务关联性
//	//第六步整体数据调整
//	dbmap := core.InitDb()
//	defer dbmap.Db.Close()
//	//	var smarr []systemmodel.Systemmodule
//	smarr := make([]systemmodel.Systemmodule, 8)
//	smarr[0] = systemmodel.Systemmodule{Modulename: "课程管理", Modulecode: "kcgl", Moduleicon: "", Superiormoduleid: 0}
//	smarr[1] = systemmodel.Systemmodule{Modulename: "远程教学", Modulecode: "ycjx", Moduleicon: "", Superiormoduleid: 0}
//	smarr[2] = systemmodel.Systemmodule{Modulename: "教室导流", Modulecode: "jsdl", Moduleicon: "", Superiormoduleid: 0}
//	smarr[3] = systemmodel.Systemmodule{Modulename: "查课表", Modulecode: "ckb", Moduleicon: "", Superiormoduleid: 0}
//	smarr[4] = systemmodel.Systemmodule{Modulename: "视频录播", Modulecode: "splb", Moduleicon: "", Superiormoduleid: 0}
//	smarr[5] = systemmodel.Systemmodule{Modulename: "设备管理", Modulecode: "sbgl", Moduleicon: "", Superiormoduleid: 0}
//	smarr[6] = systemmodel.Systemmodule{Modulename: "上课点到", Modulecode: "skdd", Moduleicon: "", Superiormoduleid: 0}
//	smarr[7] = systemmodel.Systemmodule{Modulename: "出勤统计", Modulecode: "cqtj", Moduleicon: "", Superiormoduleid: 0}
//	dbmap.AddTableWithName(systemmodel.Systemmodule{}, "systemmodule").SetKeys(true, "Id")
//	err1 := dbmap.Insert(&smarr[0], &smarr[1], &smarr[2], &smarr[3], &smarr[4], &smarr[5], &smarr[6], &smarr[7])
//	core.CheckErr(err1, "静态模块数据设置失败")
//	rlarr := make([]systemmodel.Roles, 4)
//	rlarr[0] = systemmodel.Roles{Rolesname: "管理人员"}
//	rlarr[1] = systemmodel.Roles{Rolesname: "教师人员"}
//	rlarr[2] = systemmodel.Roles{Rolesname: "学生"}
//	rlarr[3] = systemmodel.Roles{Rolesname: "测试人员"}
//	dbmap.AddTableWithName(systemmodel.Roles{}, "Roles").SetKeys(true, "Id")
//	err3 := dbmap.Insert(&rlarr[0], &rlarr[1], &rlarr[2], &rlarr[3])
//	core.CheckErr(err3, "静态角色数据设置失败")
//	rmcarr := make([]systemmodel.Rolemodulecenter, 25)
//	rmcarr[0] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[0].Id}
//	rmcarr[1] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[1].Id}
//	rmcarr[2] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[2].Id}
//	rmcarr[3] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[3].Id}
//	rmcarr[4] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[4].Id}
//	rmcarr[5] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[5].Id}
//	rmcarr[6] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[6].Id}
//	rmcarr[7] = systemmodel.Rolemodulecenter{Rolesid: rlarr[0].Id, Systemmoduleid: smarr[7].Id}
//	rmcarr[8] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[0].Id}
//	rmcarr[9] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[1].Id}
//	rmcarr[10] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[3].Id}
//	rmcarr[11] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[4].Id}
//	rmcarr[12] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[5].Id}
//	rmcarr[13] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[6].Id}
//	rmcarr[14] = systemmodel.Rolemodulecenter{Rolesid: rlarr[1].Id, Systemmoduleid: smarr[7].Id}
//	rmcarr[15] = systemmodel.Rolemodulecenter{Rolesid: rlarr[2].Id, Systemmoduleid: smarr[2].Id}
//	rmcarr[16] = systemmodel.Rolemodulecenter{Rolesid: rlarr[2].Id, Systemmoduleid: smarr[3].Id}
//	rmcarr[17] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[0].Id}
//	rmcarr[18] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[1].Id}
//	rmcarr[19] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[2].Id}
//	rmcarr[20] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[3].Id}
//	rmcarr[21] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[4].Id}
//	rmcarr[22] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[5].Id}
//	rmcarr[23] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[6].Id}
//	rmcarr[24] = systemmodel.Rolemodulecenter{Rolesid: rlarr[3].Id, Systemmoduleid: smarr[7].Id}
//	dbmap.AddTableWithName(systemmodel.Rolemodulecenter{}, "rolemodulecenter").SetKeys(true, "Id")
//	for p := 0; p < len(rmcarr); p++ {
//		err3 = dbmap.Insert(&rmcarr[p])
//		core.CheckErr(err3, "静态角色模块数据设置失败")
//	}
//	var smfarr []systemmodel.Systemmodulefunctions
//	smfarr0 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[0].Id, Functionname: "查看课程", Functionicon: "", Functioncode: "queryall"}
//	smfarr1 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[1].Id, Functionname: "远程教学列表", Functionicon: "", Functioncode: "ycjxlist"}
//	smfarr2 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[2].Id, Functionname: "教室导流查看", Functionicon: "", Functioncode: "jsdlquery"}
//	smfarr3 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[3].Id, Functionname: "查课表查看", Functionicon: "", Functioncode: "ckbquery"}
//	smfarr4 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[4].Id, Functionname: "视频录播查看", Functionicon: "", Functioncode: "splbquery"}
//	smfarr5 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[5].Id, Functionname: "设备管理列表", Functionicon: "", Functioncode: "sbgllist"}
//	smfarr6 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[6].Id, Functionname: "上课点到查看", Functionicon: "", Functioncode: "skddquery"}
//	smfarr7 := systemmodel.Systemmodulefunctions{Systemmoduleid: smarr[7].Id, Functionname: "出勤统计查看", Functionicon: "", Functioncode: "cqtjquery"}
//	smfarr = append(smfarr, smfarr0, smfarr1, smfarr2, smfarr3, smfarr4, smfarr5, smfarr6, smfarr7)
//	dbmap.AddTableWithName(systemmodel.Systemmodulefunctions{}, "systemmodulefunctions").SetKeys(true, "Id")
//	err6 := dbmap.Insert(&smfarr[0], &smfarr[1], &smfarr[2], &smfarr[3], &smfarr[4], &smfarr[5], &smfarr[6], &smfarr[7])
//	core.CheckErr(err6, "静态功能数据设置失败")
//	var rmfcarr []systemmodel.Rolemodulefunctioncenter
//	var rmfc systemmodel.Rolemodulefunctioncenter
//	dbmap.AddTableWithName(systemmodel.Rolemodulefunctioncenter{}, "rolemodulefunctioncenter").SetKeys(true, "Id")
//	for i := 0; i < len(rmcarr); i++ {
//		for j := 0; j < len(smfarr); j++ {
//			if rmcarr[i].Systemmoduleid == smfarr[j].Systemmoduleid {
//				rmfc = systemmodel.Rolemodulefunctioncenter{Rolemodulecenterid: rmcarr[i].Id, Systemmodulefunctionsid: smfarr[j].Id}
//				err3 = dbmap.Insert(&rmfc)
//				core.CheckErr(err3, "静态角色模块功能数据设置失败")
//				rmfcarr = append(rmfcarr, rmfc)
//			}
//		}
//	}
//	/*添加校区数据*/
//	cps := basicset.Campus{Campusname: "新校区", Campusicon: "", Campuscode: "xxq", Campusnums: 0}
//	dbmap.AddTableWithName(basicset.Campus{}, "campus").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&cps)
//	core.CheckErr(err3, "校区数据设置失败")
//	cg1 := basicset.College{Campusid: cps.Id, Collegename: "外国语学院", Collegenum: 0}
//	cg2 := basicset.College{Campusid: cps.Id, Collegename: "化工学院", Collegenum: 0}
//	dbmap.AddTableWithName(basicset.College{}, "college").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&cg1, &cg2)
//	core.CheckErr(err3, "校区学院数据设置失败")
//	mj1 := basicset.Major{Collegeid: cg1.Id, Majorname: "英语系", Majoricon: "", Majornum: 0}
//	mj2 := basicset.Major{Collegeid: cg1.Id, Majorname: "法语系", Majoricon: "", Majornum: 0}
//	mj3 := basicset.Major{Collegeid: cg1.Id, Majorname: "日语系", Majoricon: "", Majornum: 0}
//	mj4 := basicset.Major{Collegeid: cg1.Id, Majorname: "西班牙系", Majoricon: "", Majornum: 0}
//	mj5 := basicset.Major{Collegeid: cg2.Id, Majorname: "应用化学系", Majoricon: "", Majornum: 0}
//	mj6 := basicset.Major{Collegeid: cg2.Id, Majorname: "无机化学系", Majoricon: "", Majornum: 0}
//	mj7 := basicset.Major{Collegeid: cg2.Id, Majorname: "分析科学系", Majoricon: "", Majornum: 0}
//	mj8 := basicset.Major{Collegeid: cg2.Id, Majorname: "制药工程系", Majoricon: "", Majornum: 0}
//	dbmap.AddTableWithName(basicset.Major{}, "major").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&mj1, &mj2, &mj3, &mj4, &mj5, &mj6, &mj7, &mj8)
//	core.CheckErr(err3, "校区学院科系/专业数据设置失败")
//	cs1 := basicset.Classes{Majorid: mj1.Id, Classesname: "英语1班", Classesnum: 30, Classesicon: "", Classstate: 1}
//	cs2 := basicset.Classes{Majorid: mj1.Id, Classesname: "英语2班", Classesnum: 30, Classesicon: "", Classstate: 0}
//	cs3 := basicset.Classes{Majorid: mj2.Id, Classesname: "法语1班", Classesnum: 50, Classesicon: "", Classstate: 1}
//	cs4 := basicset.Classes{Majorid: mj7.Id, Classesname: "分析科学1班", Classesnum: 45, Classesicon: "", Classstate: 1}
//	cs5 := basicset.Classes{Majorid: mj8.Id, Classesname: "制药工程1班", Classesnum: 55, Classesicon: "", Classstate: 1}
//	cs6 := basicset.Classes{Majorid: mj7.Id, Classesname: "分析科学2班", Classesnum: 38, Classesicon: "", Classstate: 1}
//	cs7 := basicset.Classes{Majorid: mj8.Id, Classesname: "制药工程2班", Classesnum: 40, Classesicon: "", Classstate: 1}
//	cs8 := basicset.Classes{Majorid: mj7.Id, Classesname: "分析科学3班", Classesnum: 50, Classesicon: "", Classstate: 0}
//	cs9 := basicset.Classes{Majorid: mj8.Id, Classesname: "制药工程3班", Classesnum: 60, Classesicon: "", Classstate: 0}
//	cs10 := basicset.Classes{Majorid: mj7.Id, Classesname: "分析科学4班", Classesnum: 45, Classesicon: "", Classstate: 0}
//	cs11 := basicset.Classes{Majorid: mj8.Id, Classesname: "制药工程4班", Classesnum: 35, Classesicon: "", Classstate: 0}
//	dbmap.AddTableWithName(basicset.Classes{}, "classes").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&cs1, &cs2, &cs3, &cs4, &cs5, &cs6, &cs7, &cs8, &cs9, &cs10, &cs11)
//	core.CheckErr(err3, "校区学院科系/专业班级数据设置失败")
//	/*添加校区楼栋数据*/
//	bd1 := basicset.Building{Campusid: cps.Id, Buildingname: "A栋", Floorsnumber: 8, Classroomsnumber: 200}
//	bd2 := basicset.Building{Campusid: cps.Id, Buildingname: "B栋", Floorsnumber: 9, Classroomsnumber: 300}
//	bd3 := basicset.Building{Campusid: cps.Id, Buildingname: "C栋", Floorsnumber: 5, Classroomsnumber: 180}
//	dbmap.AddTableWithName(basicset.Building{}, "building").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&bd1, &bd2, &bd3)
//	core.CheckErr(err3, "校区楼栋数据设置失败")
//	fs1 := basicset.Floors{Buildingid: bd1.Id, Floorname: "1楼", Classroomnumber: 20}
//	fs2 := basicset.Floors{Buildingid: bd1.Id, Floorname: "2楼", Classroomnumber: 25}
//	fs3 := basicset.Floors{Buildingid: bd1.Id, Floorname: "3楼", Classroomnumber: 25}
//	fs4 := basicset.Floors{Buildingid: bd1.Id, Floorname: "4楼", Classroomnumber: 25}
//	fs5 := basicset.Floors{Buildingid: bd1.Id, Floorname: "5楼", Classroomnumber: 25}
//	dbmap.AddTableWithName(basicset.Floors{}, "floors").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&fs1, &fs2, &fs3, &fs4, &fs5)
//	core.CheckErr(err3, "校区楼栋楼层数据设置失败")
//	cr1 := basicset.Classrooms{Floorsid: fs1.Id, Classroomsname: "A104教室", Classroomstype: "智慧教室", Seatsnumbers: 65, Classroomstate: 0, Maxy: 0.773777, Miny: 0.459113, Maxx: 1, Minx: 0.5}
//	cr2 := basicset.Classrooms{Floorsid: fs1.Id, Classroomsname: "A102教室", Classroomstype: "智慧教室", Seatsnumbers: 55, Classroomstate: 1, Maxy: 1, Miny: 0.668516, Maxx: 0.267500, Minx: 0}
//	cr3 := basicset.Classrooms{Floorsid: fs2.Id, Classroomsname: "A103教室", Classroomstype: "普通教室", Seatsnumbers: 55, Classroomstate: 0, Maxy: 0.668516, Miny: 0.353855, Maxx: 0.267500, Minx: 0}
//	cr4 := basicset.Classrooms{Floorsid: fs2.Id, Classroomsname: "A124教室", Classroomstype: "智慧教室", Seatsnumbers: 65, Classroomstate: 1, Maxy: 0.354973, Miny: 0, Maxx: 0.704547, Minx: 0.410834}
//	cr5 := basicset.Classrooms{Floorsid: fs2.Id, Classroomsname: "A123教室", Classroomstype: "普通教室", Seatsnumbers: 65, Classroomstate: 0, Maxy: 0.354262, Miny: 0, Maxx: 1, Minx: 0.704547}
//	cr6 := basicset.Classrooms{Floorsid: fs3.Id, Classroomsname: "A301教室", Classroomstype: "普通教室", Seatsnumbers: 65, Classroomstate: 0, Maxy: 0.773777, Miny: 0.459113, Maxx: 1, Minx: 0.5}
//	cr7 := basicset.Classrooms{Floorsid: fs3.Id, Classroomsname: "A302教室", Classroomstype: "智慧教室", Seatsnumbers: 85, Classroomstate: 1, Maxy: 0.773777, Miny: 0.459113, Maxx: 1, Minx: 0.5}
//	cr8 := basicset.Classrooms{Floorsid: fs3.Id, Classroomsname: "A303教室", Classroomstype: "智慧教室", Seatsnumbers: 150, Classroomstate: 0, Maxy: 0.773777, Miny: 0.459113, Maxx: 1, Minx: 0.5}
//	cr9 := basicset.Classrooms{Floorsid: fs3.Id, Classroomsname: "A304教室", Classroomstype: "普通教室", Seatsnumbers: 80, Classroomstate: 2, Maxy: 0.773777, Miny: 0.459113, Maxx: 1, Minx: 0.5}
//	dbmap.AddTableWithName(basicset.Classrooms{}, "classrooms").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&cr1, &cr2, &cr3, &cr4, &cr5, &cr6, &cr7, &cr8, &cr9)
//	core.CheckErr(err3, "校区楼栋楼层教室数据设置失败")
//	//管理员
//	uu1 := users.Users{Loginuser: "glz001", Loginpwd: "123456", Rolesid: rlarr[0].Id, Truename: "管理员1", Nickname: "管理员1", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 1, Usermac: "", Birthday: ""}
//	uu2 := users.Users{Loginuser: "glz002", Loginpwd: "123456", Rolesid: rlarr[0].Id, Truename: "管理员2", Nickname: "管理员2", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 2, Usermac: "", Birthday: ""}
//	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&uu1, &uu2)
//	core.CheckErr(err3, "管理员数据设置失败")
//	//老师
//	uu3 := users.Users{Loginuser: "js001", Loginpwd: "123456", Rolesid: rlarr[1].Id, Truename: "教师人员1", Nickname: "教师人员1", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 1, Usermac: "", Birthday: ""}
//	uu4 := users.Users{Loginuser: "js002", Loginpwd: "123456", Rolesid: rlarr[1].Id, Truename: "教师人员2", Nickname: "教师人员2", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 2, Usermac: "", Birthday: ""}
//	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&uu3, &uu4)
//	core.CheckErr(err3, "老师数据设置失败")
//	//学生
//	uu5 := users.Users{Loginuser: "2015wgy01001", Loginpwd: "123456", Rolesid: rlarr[2].Id, Truename: "15届外国语英语1班A", Nickname: "15届外国语英语1班A", Userheadimg: "", Userphone: "", Userstate: 2, Usersex: 1, Usermac: "", Birthday: "1989.5.5"}
//	uu6 := users.Users{Loginuser: "2015wgy01002", Loginpwd: "123456", Rolesid: rlarr[2].Id, Truename: "15届外国语英语1班B", Nickname: "15届外国语英语1班B", Userheadimg: "", Userphone: "", Userstate: 2, Usersex: 2, Usermac: "", Birthday: "1989.6.5"}
//	uu7 := users.Users{Loginuser: "2016wgy02001", Loginpwd: "123456", Rolesid: rlarr[2].Id, Truename: "16届外国语英语1班C", Nickname: "16届外国语英语1班C", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 1, Usermac: "", Birthday: "1989.6.5"}
//	uu8 := users.Users{Loginuser: "2016wgy02002", Loginpwd: "123456", Rolesid: rlarr[2].Id, Truename: "16届外国语英语1班D", Nickname: "16届外国语英语1班D", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 2, Usermac: "", Birthday: "1989.6.5"}
//	uu9 := users.Users{Loginuser: "2016wgy02003", Loginpwd: "123456", Rolesid: rlarr[2].Id, Truename: "16届外国语英语1班E", Nickname: "16届外国语英语1班E", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 1, Usermac: "", Birthday: "1989.6.5"}
//	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&uu5, &uu6, &uu7, &uu8, &uu9)
//	core.CheckErr(err3, "学生数据设置一失败")
//	us1 := users.Students{Id: uu5.Id, Enrollmentyear: 2015, Homeaddress: "湖南省长沙市天心区", Nowaddress: "湖南省长沙市岳麓区", Classid: cs1.Id, Infostate: 1, Currentstate: 1}
//	us2 := users.Students{Id: uu6.Id, Enrollmentyear: 2015, Homeaddress: "湖南省长沙市芙蓉区", Nowaddress: "湖南省长沙市岳麓区", Classid: cs1.Id, Infostate: 1, Currentstate: 1}
//	us3 := users.Students{Id: uu7.Id, Enrollmentyear: 2016, Homeaddress: "湖南省长沙市天心区", Nowaddress: "湖南省长沙市芙蓉区", Classid: cs2.Id, Infostate: 1, Currentstate: 1}
//	us4 := users.Students{Id: uu8.Id, Enrollmentyear: 2016, Homeaddress: "湖南省长沙市天心区", Nowaddress: "湖南省长沙市岳麓区", Classid: cs2.Id, Infostate: 1, Currentstate: 1}
//	us5 := users.Students{Id: uu9.Id, Enrollmentyear: 2016, Homeaddress: "湖南省长沙市岳麓区", Nowaddress: "湖南省长沙市天心区", Classid: cs2.Id, Infostate: 1, Currentstate: 1}
//	dbmap.AddTableWithName(users.Students{}, "students")
//	err3 = dbmap.Insert(&us1, &us2, &us3, &us4, &us5)
//	core.CheckErr(err3, "学生数据设置二失败")
//	//测试用户
//	uu10 := users.Users{Loginuser: "text001", Loginpwd: "123456", Rolesid: rlarr[3].Id, Truename: "测试人员1", Nickname: "测试人员1", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 1, Usermac: "", Birthday: ""}
//	uu11 := users.Users{Loginuser: "text002", Loginpwd: "123456", Rolesid: rlarr[3].Id, Truename: "测试人员2", Nickname: "测试人员2", Userheadimg: "", Userphone: "", Userstate: 1, Usersex: 1, Usermac: "", Birthday: ""}
//	dbmap.AddTableWithName(users.Users{}, "users").SetKeys(true, "Id")
//	err3 = dbmap.Insert(&uu10, &uu11)
//	core.CheckErr(err3, "测试用户数据设置一失败")

//	/*
//		学科分类表
//		课程表
//		课程章节表
//		课程班级中间表
//		课程班级章节中间表
//		课程章节附件表
//	*/
//}
