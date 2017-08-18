package actiondataControllers

import (
	"encoding/json"
	"fmt"
	//	"io"
	"io/ioutil"
	//	"log"
	//	"net/http"
	//	"os"
	"os/exec"
	//	"reflect"
	"strconv"
	"strings"
	//	"time"
)

import (
	"github.com/gin-gonic/gin"
	//	"github.com/go-fsnotify/fsnotify"
)

import (
	"basicproject/DataAccess/actiondataDataAccess"
	//	"basicproject/DataAccess/curriculumDataAccess"
	//	"basicproject/DataAccess/liveDataAccess"
	"basicproject/DataAccess/usersDataAccess"

	//	"basicproject/commons"
	//	"basicproject/commons/xdebug"

	"basicproject/model/actiondata"
	xcon "xutils/xconfig"
	core "xutils/xcore"
	//	"xutils/xdebug"
	//	"basicproject/model/equipment"
	//	"basicproject/model/live"
	//	"basicproject/videosrv"
	"basicproject/viewmodel"
)

var rclist []RecordCmd
var cfc xcon.Config
var IsMonitorFile int

//var onlycmdsend []equipment.CommandSendlog

func init() {
	cfc.InitConfig("./config.ini")
}

type RecordCmd struct {
	Classroomid int
	Cmd         *exec.Cmd
}

/*
获取现场考勤数据
根据课程班级章节中间表ID查询
[班级名称，班级人数，班级图标,班级内人数列表[学生Id，学生姓名，学生头像，是否到堂]]
*/
func GetPointtos(c *gin.Context) {
	var rd core.Returndata
	/*
		提交课程班级章节中间表ID:Curriculumclassroomchaptercentreid
	*/
	var lg viewmodel.GetPointtos
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "GetPointtos|登录数据json格式转换错误:"+string(data)+":")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "GetPointtos|令牌验证数据转换失败:"+string(data)+":")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "getpointtos", dbmap)
			if rd.Rcode == "1000" {
				rd.Result = actiondataDataAccess.QueryClassPointtos(lg, dbmap)
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
1：更改课程章节的上课状态
2：更改教室的使用状态
*/
func ChangeClassState(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.ChangeClassState
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	fmt.Println(string(data))
	core.CheckErr(errs1, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")
	IsExec := false
	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			IsExec = true
		} else {
			rd.Rcode = "1002"
			rd.Reason = "Token令牌不对"
		}
	} else if lg.Os == "PC" {
		us := usersDataAccess.QueryUsersInfo(lg.Loginuser, dbmap)
		if us.Id > 0 {
			bt.Usersid = us.Id
			bt.Rolestype = us.Rolesid
			IsExec = true
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	if IsExec {
		rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "ChangeClassState", dbmap)
		if rd.Rcode == "1000" {
			rd = actiondataDataAccess.ChangeClassState(lg.Classroomid, bt.Usersid, lg.State, lg.Ccccid, lg.Ccccids, dbmap)
		} else {
			rd.Rcode = "1105"
			rd.Reason = "此功能未授权,不可用"
		}
	}
	c.JSON(200, rd)
}

/*
开始录制视频
需要处理的业务如下
1:服务器开始拉去视频流
2:发送录制命令给客户端进行视频录制
需要参数[教室Id，教室内的监控ip地址，教室电脑的ip地址，发送给监控的命令，调用服务器端本地ffmege应用]
3:文件保存，并更新到对应的课程章节数据中
*/
/*
func BeginVideo(c *gin.Context) {
	var rd core.Returndata
	var pc viewmodel.PostCollection
	data, _ := ioutil.ReadAll(c.Request.Body)
	var bt core.BasicsToken
	errs1 := json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "BeginVideo|令牌验证数据转换失败")
	errs1 = json.Unmarshal(data, &pc)
	core.CheckErr(errs1, "BeginVideo|获取教室Id和录制控制状态")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "BeginVideo", dbmap)
			if rd.Rcode == "1000" {
				rd.Rcode = ""
				//第一步：查询教室内的相关设备的Ip地址
				var vc equipment.VideoConfig
				var cc equipment.ClassroomComputerConfig
				selecterr := dbmap.SelectOne(&vc, "select * from videoconfig where Classroomid=?;", pc.Classroomid)
				core.CheckErr(selecterr, "BeginVideo|查询教室内录像设备ip地址等数据")
				if vc.Id == 0 {
					rd.Rcode = "1003"
					rd.Reason = "提交的数据不正确"
				}
				selecterr = dbmap.SelectOne(&cc, "select * from classroomcomputerconfig where Classroomid=?;", pc.Classroomid)
				core.CheckErr(selecterr, "BeginVideo|获取教室内的电脑配置等相关数据")
				if cc.Id == 0 {
					rd.Rcode = "1003"
					rd.Reason = "提交的数据不正确"
				}
				if rd.Rcode == "" {
					if pc.State == 1 { //开启录制
						for _, v := range rclist {
							if v.Classroomid == pc.Classroomid { //如果找到相等的值，则抛出正在录制的提升
								rd.Rcode = "1005"
								rd.Reason = "视频正在录制"
							}
						}
						if rd.Rcode == "" {
							jsonstr := "{\"CmdStr\":\"" + "runprogram|beginvdieo" + "\",\"CmdType\":\"runprogram\",\"CmdUsersId\":" + strconv.Itoa(bt.Usersid) + ",\"Classroomid\":" + strconv.Itoa(pc.Classroomid) + "}"
							resp, err := http.Post(cfc.Read("udpserver", "httpurl")+"/collectcmd", "application/json", strings.NewReader(jsonstr))
							if err == nil {
								respdata, resperr := ioutil.ReadAll(resp.Body)
								respdatastr := string(respdata)
								json.Unmarshal(respdata, &rd)
								if resperr == nil && rd.Rcode == "ok" {
									//fmt.Println(respdatastr)
									rd.Rcode = "1000"
									filename := time.Now().Format("2006t01m02d15t04t05") + ".mp4"
									path := cfc.Read("file", "upfilevideo")                                                         //"E:\\workgo\\src\\zndx2\\templates\\upfile"                  //需要写到配置文件中去
									path = path + "\\" + time.Now().Format("20060102") + "\\" + filename                            //实际的物理路径
									virtualpath := cfc.Read("file", "virtualfile") + time.Now().Format("20060102") + "/" + filename //"http://172.17.70.222:8081/web/upfile/" + time.Now().Format("20060102") + "/" + filename
									//go ExecCommand(path, pc.Classroomid) //发送指令给摄像机开始录像
									ExecCommand(path, pc.Classroomid) //发送指令给摄像机开始录像
									go MonitorFile()
									rd = curriculumDataAccess.UpdateLivePath(bt.Usersid, pc.Classroomid, path, virtualpath, 1, dbmap)
									vc.CameraState = 1
									_, selecterr = dbmap.Exec("update videoconfig set CameraState=? where Classroomid=?;", vc.CameraState, vc.Classroomid)
									core.CheckErr(selecterr, "BeginVideo|更新教室内摄像机状态")
									cc.ComputerState = 1
									_, selecterr = dbmap.Exec("update classroomcomputerconfig set ComputerState=? where Classroomid=?;", cc.ComputerState, cc.Classroomid)
									core.CheckErr(selecterr, "BeginVideo|更新教室内电脑状态")
								} else {
									rd.Rcode = "1004"
									rd.Reason = "发送指令失败:" + rd.Reason
									core.CheckErr(resperr, "发送开启录屏指令失败：")
								}
							} else {
								rd.Rcode = "1004"
								rd.Reason = "发送指令失败:" + rd.Reason
								core.CheckErr(err, "发送开启录屏指令失败：")
							}
						}
					} else { //关闭录制
						for i, v := range rclist {
							if v.Classroomid == pc.Classroomid { //如果找到相等的值
								rd.Rcode = "1"
								var filename = "./temporarylogfile/" + strconv.Itoa(pc.Classroomid) + ".txt"
								var f *os.File
								defer f.Close()
								var err1 error
								if core.CheckFileIsExist(filename) { //如果文件存在
									f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
								} else {
									f, err1 = os.Create(filename) //创建文件
								}
								core.CheckErr(err1, "BeginVideo|关闭录制时错误")
								f.WriteString("q")
								f.Sync()
								rclist = append(rclist[:i], rclist[i+1:]...) //移除自身
								_, selecterr = dbmap.Exec("update videoconfig set CameraState=? where Classroomid=?;", 0, v.Classroomid)
								core.CheckErr(selecterr, "BeginVideo|更新教室内摄像机状态")
								_, selecterr = dbmap.Exec("update classroomcomputerconfig set ComputerState=? where Classroomid=?;", 0, v.Classroomid)
								core.CheckErr(selecterr, "BeginVideo|更新教室内电脑状态")
								break
							}
						}
						if rd.Rcode == "1" {
							jsonstr := "{\"CmdStr\":\"" + "runprogram|endvideo" + "\",\"CmdType\":\"runprogram\",\"CmdUsersId\":" + strconv.Itoa(bt.Usersid) + ",\"Classroomid\":" + strconv.Itoa(pc.Classroomid) + "}"
							resp, err := http.Post(cfc.Read("udpserver", "httpurl")+"/collectcmd", "application/json", strings.NewReader(jsonstr))
							if err == nil {
								respdata, resperr := ioutil.ReadAll(resp.Body)
								respdatastr := string(respdata)
								json.Unmarshal(respdata, &rd)
								if resperr == nil && rd.Rcode == "ok" {
									fmt.Println(respdatastr)
									rd.Rcode = "1000"
								} else {
									rd.Rcode = "1004"
									rd.Reason = "发送指令失败:" + rd.Reason
								}
							} else {
								rd.Rcode = "1004"
								rd.Reason = "指令发送失败"
								core.CheckErr(err, "指令发送失败:")
							}
						} else {
							rd.Rcode = "1003"
							rd.Reason = "此教室没有在录制视频"
						}
					}
				}
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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
*/
/*
func BeginVideo(c *gin.Context) {
	var rd core.Returndata
	var pc viewmodel.PostCollection
	data, _ := ioutil.ReadAll(c.Request.Body)
	var bt core.BasicsToken
	errs1 := json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "BeginVideo|令牌验证数据转换失败")
	errs1 = json.Unmarshal(data, &pc)
	core.CheckErr(errs1, "BeginVideo|获取教室Id和录制控制状态")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "BeginVideo", dbmap)
			if rd.Rcode == "1000" {
				rd.Rcode = ""
				var bte []byte                       //定义发送的内容
				var arrcs []equipment.CommandSendlog //真正发送过去的数据
				var vc equipment.VideoConfig
				var cc equipment.ClassroomComputerConfig
				errs1 = dbmap.SelectOne(&vc, "select * from videoconfig where Classroomid=?;", pc.Classroomid)
				core.CheckErr(errs1, "BeginVideo|查询教室内录像设备ip地址等数据")
				if vc.Id == 0 {
					rd.Rcode = "1013"
					rd.Reason = "此教室未配置录像机"
				}
				errs1 = dbmap.SelectOne(&cc, "select * from classroomcomputerconfig where Classroomid=?;", pc.Classroomid)
				core.CheckErr(errs1, "BeginVideo|获取教室内的电脑配置等相关数据")
				if cc.Id == 0 {
					rd.Rcode = "1013"
					rd.Reason = "此教室未配置电脑"
				}
				path1 := cfc.Read("file", "upfilevideo") //需要写到配置文件中去
				path2 := cfc.Read("file", "clientvideo") //需要写到配置文件中去
				virtualpath1 := cfc.Read("file", "virtualfile")
				virtualpath2 := cfc.Read("file", "virtualfile")
				if pc.State == 1 { //开始录制
					//第一步：查询教室内的相关设备的Ip地址
					if vc.CameraState == -1 {
						rd.Rcode = "1011"
						rd.Reason = "此教室摄像机连接异常无法发送指令"
					}
					if vc.CameraState == 1 {
						rd.Rcode = "1012"
						rd.Reason = "此教室摄像机正在录制视频"
					}
					if cc.ComputerState == -1 {
						rd.Rcode = "1014"
						rd.Reason = "此教室电脑连接异常无法发送指令"
					}
					if cc.ComputerState == 1 {
						rd.Rcode = "1015"
						rd.Reason = "此教室电脑正在录制视频"
					}

					if rd.Rcode == "" {
						var arrcr []equipment.CommandRecord //需要发送的命令集合
						_, errs1 = dbmap.Select(&arrcr, "select * from CommandRecord where CmdCode in('v001','c002');")
						core.CheckErr(errs1, "BeginVideo|获取命令配置数据")
						var cs equipment.CommandSendlog
						path1 = path1 + "\\" + time.Now().Format("20060102") + "\\" //实际的物理路径
						virtualpath1 := virtualpath1 + time.Now().Format("20060102")
						virtualpath2 := virtualpath2 + time.Now().Format("20060102")
						for _, v := range arrcr {
							cs.Classroomid = pc.Classroomid
							cs.CmdDate = time.Now().Format("2006-01-02 15:04:05")
							filename := time.Now().Format("2006t01m02d15t04t05") + ".mp4"
							if v.CmdCode == "v001" { //摄像机
								path1 = path1 + filename //实际的物理路径
								virtualpath1 = virtualpath1 + filename
								//path1 = "C:\\" + filename
								cs.CmdIp = "172.16.100.102" //vc.CameraIp
								cs.CmdPort = 1514           //vc.CameraPort
								cs.CmdStr = v.CmdStr
								cs.CmdStr = strings.Replace(cs.CmdStr, "{0}", vc.CameraLoginUser, -1)
								cs.CmdStr = strings.Replace(cs.CmdStr, "{1}", vc.CameraLoginPass, -1)
								cs.CmdStr = strings.Replace(cs.CmdStr, "{2}", vc.CameraIp, -1)
								cs.CmdStr = strings.Replace(cs.CmdStr, "{3}", strconv.Itoa(vc.CameraPort), -1)
								cs.CmdStr = strings.Replace(cs.CmdStr, "{4}", "300", -1)
								cs.CmdStr = strings.Replace(cs.CmdStr, "{5}", path1, -1)
								cs.CmdType = "0"
								cs.CmdState = 0
							} else if v.CmdCode == "c002" { //电脑
								path2 = path2 + filename //实际的物理路径
								virtualpath2 = virtualpath2 + filename
								cs.CmdIp = cc.Computerip
								cs.CmdPort = cc.Computerupdaccept
								cs.CmdStr = v.CmdStr
								cs.CmdStr = strings.Replace(cs.CmdStr, "{0}", "300", -1)
								cs.CmdStr = strings.Replace(cs.CmdStr, "{1}", path2, -1)
								cs.CmdType = "0"
								cs.CmdState = 0
							}
							cs.CmdUsersId = bt.Usersid
							dbmap.AddTableWithName(equipment.CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
							//							if v.CmdCode == "v001" {
							errs1 = dbmap.Insert(&cs)
							core.CheckErr(errs1, "BeginVideo|commandsendlog|将命令插入到数据库中去")
							arrcs = append(arrcs, cs)
							//							}
							cs = equipment.CommandSendlog{}
						}
					}
				} else { //结束录制
					for _, kv := range onlycmdsend {
						if kv.Classroomid == pc.Classroomid {
							kv.CmdStr = "q"
							arrcs = append(arrcs, kv)
						}
					}
				}
				if rd.Rcode == "" {
					bte, errs1 = json.Marshal(&arrcs)
					fmt.Println("发送udp数据：", string(bte))
					core.CheckErr(errs1, "BeginVideo|commandsendlog|数据转换失败")
					fmt.Println(string(bte))

					resp, err := http.Post(cfc.Read("udpserver", "httpurl")+"/collectcmd", "application/json", strings.NewReader(string(bte)))
					if err == nil {
						respdata, resperr := ioutil.ReadAll(resp.Body)
						core.CheckErr(resperr, "BeginVideo|接受发送udp服务器响应回来的数据")
						fmt.Println(string(respdata))
						var resparrcs []equipment.CommandSendlog //接受响应回来的数据
						resperrjosn := json.Unmarshal(respdata, &resparrcs)
						core.CheckErr(resperrjosn, "BeginVideo|转换发送udp服务器响应回来的数据")
						if resperrjosn == nil && resperr == nil && len(resparrcs) > 0 {
							rd.Rcode = "1000"
							dbmap.AddTableWithName(equipment.CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
							for _, v := range resparrcs {
								_, errs1 = dbmap.Update(&v)
								core.CheckErr(errs1, "BeginVideo|commandsendlog|将数据更新到数据库中去")
							}
							//							rd = curriculumDataAccess.UpdateLivePath(bt.Usersid, pc.Classroomid, path1, virtualpath1, 1, dbmap)
							vc.CameraState = pc.State
							_, errs1 = dbmap.Exec("update videoconfig set CameraState=? where Classroomid=?;", vc.CameraState, vc.Classroomid)
							core.CheckErr(errs1, "BeginVideo|更新教室内摄像机状态")
							cc.ComputerState = pc.State
							_, errs1 = dbmap.Exec("update classroomcomputerconfig set ComputerState=? where Classroomid=?;", cc.ComputerState, cc.Classroomid)
							core.CheckErr(errs1, "BeginVideo|更新教室内电脑状态")
							for _, kv := range resparrcs {
								if kv.CmdState > 4 { //判断是否完结[5,6是代表已完成或者已终止]
									for kis, kvs := range onlycmdsend {
										if kv.Id == kvs.Id {
											onlycmdsend = append(onlycmdsend[:kis], onlycmdsend[kis+1:]...)
										}
									}
								} else {
									onlycmdsend = append(onlycmdsend, kv)
								}
								dbmap.Update(&kv) //将信息更新到服务器中
							}
						} else {
							rd.Rcode = "1004"
							rd.Reason = "发送指令失败:" + rd.Reason
							core.CheckErr(resperr, "发送开启录屏指令失败：")
						}
					} else {
						rd.Rcode = "1023"
						rd.Reason = "发送udp请求失败"
						core.CheckErr(err, "BeginVideo|commandsendlog|发送udp请求")
					}
				}
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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
*/
/*
func BeginVideo2(c *gin.Context) {
	var responses = commons.ResponseMsgSet_Instance()
	var rd = core.Returndata{Result: ""}
	var token core.BasicsToken
	data, err := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &token); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		c.JSON(http.StatusOK, rd)
		xdebug.DebugError(err)
		return
	}

	var request videosrv.VideoCaptureRequest
	if err = json.Unmarshal(data, &request); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		c.JSON(http.StatusOK, rd)
		xdebug.DebugError(err)
		return
	}

	rd.Rcode = strconv.Itoa(responses.TOKEN_INCORRECT.Code)
	rd.Reason = responses.TOKEN_INCORRECT.Text
	tokenText := token.Token
	tokenArray := strings.Split(string(tokenText), "|")
	if len(tokenArray) == 3 && (tokenArray[2] == strconv.Itoa(token.Usersid) && tokenArray[1] == strconv.Itoa(token.Rolestype)) { //判断账号是否正确
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		rd = usersDataAccess.CheckVaild(token.Rolestype, token.Usersid, "BeginVideo", dbmap)
		if rd.Rcode == strconv.Itoa(responses.SUCCESS.Code) {
			rd.Reason = responses.SUCCESS.Text
			var vConfig equipment.VideoConfig
			err = dbmap.SelectOne(&vConfig, "select * from videoconfig where Classroomid=?;", request.ClassroomID)
			xdebug.DebugError(err)
			cameraFile := time.Now().Format("20060102_150405.999") + ".Camera.mp4"
			if err == nil { // 查找到教室摄像机信息 发送命令
				var cmd = videosrv.VideoCaptureCommand{
					ClassroomID: vConfig.Classroomid,
					CmdType:     videosrv.CmdAction_BeginVideo,
					TargetArgs: videosrv.TargetArgs{
						TargetIP:   vConfig.CameraIp,
						TargetPort: vConfig.CameraPort,
						TargetUser: vConfig.CameraLoginUser,
						TargetPass: vConfig.CameraLoginPass},
					CurriculumArgs: request.CurriculumArgs,
					FFmpegVideoArgs: videosrv.FFmpegVideoArgs{
						VideoFile:     cameraFile,
						VideoDuration: request.CurriculumArgs.CurriculumDuration}}
				cmddata, _ := json.Marshal(&cmd)
				var msg = videosrv.CmdMessage{
					CmdID:    "Camera",
					CmdType:  reflect.TypeOf(videosrv.VideoCaptureCommand{}).Name(),
					JsonText: string(cmddata)}
				addr := videosrv.FFmpegServer_UDPAddr()
				go videosrv.SendCmdMessage(msg, addr)

				dbmap.AddTableWithName(equipment.CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
				cmdlog := equipment.CommandSendlog{
					CmdIp:       cmd.TargetArgs.TargetIP,
					CmdPort:     cmd.TargetArgs.TargetPort,
					Classroomid: cmd.ClassroomID,
					CmdStr:      "BeginVideo"}
				err = dbmap.Insert(&cmdlog)
				go xdebug.DebugError(err)
			} else {
				rd.Rcode = "1013"
				rd.Reason = "此教室未配置录像机"
			}

			scrennFile := time.Now().Format("20060102_150405.999") + ".Screen.mp4"
			var computer equipment.ClassroomComputerConfig
			err = dbmap.SelectOne(&computer, "select * from classroomcomputerconfig where Classroomid=?;", request.ClassroomID)
			xdebug.DebugError(err)
			if err == nil {
				var cmd = videosrv.VideoCaptureCommand{
					ClassroomID: request.ClassroomID,
					CmdType:     videosrv.CmdAction_BeginVideo,
					TargetArgs: videosrv.TargetArgs{
						TargetIP:   computer.Computerip, // 教学电脑IP地址
						TargetPort: computer.Computerupdaccept },//UDP端口
					CurriculumArgs: request.CurriculumArgs,
					FFmpegVideoArgs: videosrv.FFmpegVideoArgs{
						VideoFile:     scrennFile,
						VideoDuration: request.CurriculumArgs.CurriculumDuration}}
				cmddata, _ := json.Marshal(&cmd)
				var msg = videosrv.CmdMessage{
					CmdID:    "Screen",
					CmdType:  reflect.TypeOf(videosrv.VideoCaptureCommand{}).Name(),
					JsonText: string(cmddata)}
				addr := cmd.TargetIP + ":" + strconv.Itoa(cmd.TargetPort)
				go videosrv.SendCmdMessage(msg, addr)

				dbmap.AddTableWithName(equipment.CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
				cmdlog := equipment.CommandSendlog{
					CmdIp:       cmd.TargetArgs.TargetIP,
					CmdPort:     cmd.TargetArgs.TargetPort,
					Classroomid: cmd.ClassroomID,
					CmdStr:      "BeginVideo"}
				err = dbmap.Insert(&cmdlog)
				go xdebug.DebugError(err)
			}

			tid := request.CurriculumClassroomChapterID
			if liveObj, ok := liveDataAccess.QueryLiveInfoByID(tid, dbmap); ok {
				liveObj.Whenlong = request.CurriculumArgs.CurriculumDuration
				liveObj.Livepath1 = videosrv.VOD_HttpPath() + cameraFile
				liveObj.Livepath2 = videosrv.VOD_HttpPath() + scrennFile
				liveObj.Begindate = time.Now().Format("2006-01-02 15:04:05")
				dbmap.AddTableWithName(live.Lives{}, "lives").SetKeys(true, "Id")
				err := dbmap.Insert(&liveObj)
				xdebug.DebugError(err)
			}
		} else {
			rd.Rcode = strconv.Itoa(responses.AUTH_LIMITED.Code)
			rd.Reason = responses.AUTH_LIMITED.Text
		}
	}

	c.JSON(http.StatusOK, rd)
}
*/
/*
func EndVideo2(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")

	var rd core.Returndata
	data, err := ioutil.ReadAll(c.Request.Body)
	xdebug.HandleError(err)

	var responses = commons.ResponseMsgSet_Instance()
	var pc videosrv.VideoCaptureRequest
	if err = json.Unmarshal(data, &pc); err != nil {
		rd.Rcode = strconv.Itoa(responses.DATA_MALFORMED.Code)
		rd.Reason = responses.DATA_MALFORMED.Text
		c.JSON(http.StatusOK, rd)
		xdebug.DebugError(err)
		return
	}

	dbmap := core.InitDb()
	defer dbmap.Db.Close()
	var vConfig equipment.VideoConfig
	err = dbmap.SelectOne(&vConfig, "select * from videoconfig where Classroomid=?;", pc.ClassroomID)
	if err == nil { // 查找到教室摄像机信息 发送拉流命令
		var cmd = videosrv.VideoCaptureCommand{
			ClassroomID: vConfig.Classroomid,
			CmdType:     2, // StopVideo
			TargetArgs: videosrv.TargetArgs{
				TargetIP:   vConfig.CameraIp,
				TargetPort: vConfig.CameraPort,
				TargetUser: vConfig.CameraLoginUser,
				TargetPass: vConfig.CameraLoginPass},
			CurriculumArgs: pc.CurriculumArgs}
		data, _ := json.Marshal(&cmd)
		var msg = videosrv.CmdMessage{
			CmdType:  reflect.TypeOf(videosrv.VideoCaptureCommand{}).Name(),
			JsonText: string(data)}
		addr := videosrv.FFmpegServer_UDPAddr()
		go videosrv.SendCmdMessage(msg, addr)

		dbmap.AddTableWithName(equipment.CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
		commandlog := equipment.CommandSendlog{
			CmdIp:       cmd.TargetIP,
			CmdPort:     cmd.TargetPort,
			Classroomid: vConfig.Classroomid,
			CmdStr:      "StopVideo"}
		err = dbmap.Insert(&commandlog)
		xdebug.DebugError(err)
	}

	var computer equipment.ClassroomComputerConfig
	err = dbmap.SelectOne(&computer, "select * from classroomcomputerconfig where Classroomid=?;", pc.ClassroomID)
	xdebug.DebugError(err)
	if err == nil {
		var cmd = videosrv.VideoCaptureCommand{
			ClassroomID: pc.ClassroomID,
			CmdType:     videosrv.CmdAction_StopVideo,
			TargetArgs: videosrv.TargetArgs{
				TargetIP:   computer.Computerip, // 教学电脑IP地址
				TargetPort: computer.Computerupdaccept },//UDP端口
			CurriculumArgs: pc.CurriculumArgs,
			FFmpegVideoArgs: videosrv.FFmpegVideoArgs{
				VideoDuration: pc.CurriculumArgs.CurriculumDuration}}
		cmddata, _ := json.Marshal(&cmd)
		var msg = videosrv.CmdMessage{
			CmdID:    "Screen",
			CmdType:  reflect.TypeOf(videosrv.VideoCaptureCommand{}).Name(),
			JsonText: string(cmddata)}
		addr := cmd.TargetIP + ":" + strconv.Itoa(cmd.TargetPort)
		go videosrv.SendCmdMessage(msg, addr)

		dbmap.AddTableWithName(equipment.CommandSendlog{}, "commandsendlog").SetKeys(true, "Id")
		cmdlog := equipment.CommandSendlog{
			CmdIp:       cmd.TargetArgs.TargetIP,
			CmdPort:     cmd.TargetArgs.TargetPort,
			Classroomid: cmd.ClassroomID,
			CmdStr:      "EndVideo"}
		err = dbmap.Insert(&cmdlog)
		go xdebug.DebugError(err)
	}

	rd.Rcode = strconv.Itoa(responses.SUCCESS.Code)
	rd.Reason = responses.SUCCESS.Text
	c.JSON(http.StatusOK, rd)
}
*/
/*
//上传视频文件
func UpVideoFile(c *gin.Context) {
	var rd core.Returndata
	Classroomstrid := c.Query("classroomid")
	Loginip := c.Query("loginip")
	LoginMac := c.Query("loginmac")
	fmt.Println(Loginip, LoginMac, Classroomstrid)
	Classroomid, interr := strconv.Atoi(Classroomstrid)
	if Classroomid > 0 && interr == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		path := cfc.Read("file", "upfilevideo")                                        //"E:\\workgo\\src\\zndx2\\templates\\upfile" //需要写到配置文件中去
		path = path + "\\" + time.Now().Format("20060102")                             //实际的物理路径
		os.MkdirAll(path, 0777)                                                        //创建文件夹
		virtualpath := cfc.Read("file", "virtualfile") + time.Now().Format("20060102") //"http://172.17.70.222:8081/web/upfile/" + time.Now().Format("20060102")
		fmt.Printf("%+v\n", c.Request.Body)
		fl, hdfl, flerr := c.Request.FormFile("FileData")
		fmt.Printf("%+v \n %+v \n %s", fl, hdfl, flerr)
		if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
		F:
			for i, v := range c.Request.MultipartForm.File {
				fmt.Println(i)
				if len(v) > 0 {
					f, errf := v[0].Open()
					defer f.Close()
					fmt.Println(errf)
					fname := strings.TrimLeft(v[0].Filename, "c:")
					path = path + "\\" + fname
					virtualpath = virtualpath + "/" + fname
					fv, errfv := os.Create(path)
					if errfv != nil {
						fmt.Println(errfv.Error())
					}
					defer fv.Close()
					_, errcp := io.Copy(fv, f)
					if errcp != nil {
						fmt.Println(errcp.Error())
					}
					break F
				}
			}
			rd = curriculumDataAccess.UpdateLivePath(0, Classroomid, path, virtualpath, 0, dbmap)
		} else {
			rd.Rcode = "1003"
			rd.Reason = "未提交数据文件"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}
*/
/*
//接收客户端传入过来的消息
func UdpClient(c *gin.Context) {
	var rd core.Returndata
	var csl equipment.CommandSendlog
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &csl)
	core.CheckErr(errs1, "接受udp客户端发送过来的信息")
	fmt.Printf("%+v\n", csl)
	if errs1 == nil {
		dbmap := core.InitDb()
		defer dbmap.Db.Close()
		//第一步：查询教室内的相关设备的Ip地址
		var cc equipment.ClassroomComputerConfig
		selecterr := dbmap.SelectOne(&cc, "select * from ClassroomComputerConfig where Computermac='"+csl.CmdMac+"';")
		core.CheckErr(selecterr, "查询教室内电脑设备ip地址等数据")
		fmt.Println("-------------------------")
		fmt.Printf("%+v\n", cc)
		if cc.Id == 0 {
			rd.Rcode = "1003"
			rd.Reason = "提交的数据不正确"
		}
		if rd.Rcode == "" {
			dbmap.AddTableWithName(equipment.ClassroomComputerConfig{}, "classroomcomputerconfig").SetKeys(true, "Id")
			if csl.CmdStr == "on" { //打开
				if cc.ComputerState < 0 {
					cc.ComputerState = 0
					cc.OpenDate = time.Now().Format("2006-01-02 15:04:05")
				} else {
					cc.UpdateDate = time.Now().Format("2006-01-02 15:04:05")
				}
			} else if csl.CmdStr == "off" { //关闭
				cc.ComputerState = -1
			}
			dbmap.Update(&cc)
			rd.Rcode = "1000"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)

}
*/
/*
//关闭电脑
func CloseComputer(c *gin.Context) {
	var rd core.Returndata
	var pc viewmodel.PostCollection
	data, _ := ioutil.ReadAll(c.Request.Body)
	var bt core.BasicsToken
	errs1 := json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")
	errs1 = json.Unmarshal(data, &pc)
	core.CheckErr(errs1, "获取教室Id和录制控制状态")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "CloseComputer", dbmap)
			if rd.Rcode == "1000" {
				rd.Rcode = ""
				//第一步：查询教室内的相关设备的Ip地址
				var cc equipment.ClassroomComputerConfig
				selecterr := dbmap.SelectOne(&cc, "select * from classroomcomputerconfig where Classroomid=?;", pc.Classroomid)
				core.CheckErr(selecterr, "查询教室内电脑设备ip地址等数据")
				fmt.Println(cc.Id)
				if cc.Id == 0 {
					rd.Rcode = "1003"
					rd.Reason = "提交的数据不正确"
				}
				if rd.Rcode == "" {
					jsonstr := "{\"CmdStr\":\"" + "cmd|shutdown ~s ~t 00" + "\",\"CmdType\":\"run\",\"CmdUsersId\":" + strconv.Itoa(bt.Usersid) + ",\"Classroomid\":" + strconv.Itoa(pc.Classroomid) + "}"
					resp, err := http.Post(cfc.Read("udpserver", "httpurl")+"/closecmd", "application/json", strings.NewReader(jsonstr))
					if err == nil {
						respdata, resperr := ioutil.ReadAll(resp.Body)
						respdatastr := string(respdata)
						json.Unmarshal(respdata, &rd)
						if resperr == nil && rd.Rcode == "ok" {
							fmt.Println(respdatastr)
							rd.Rcode = "1000"
						} else {
							rd.Rcode = "1004"
							rd.Reason = "发送指令失败:" + rd.Reason
						}
					}
				}
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
*/

/*
//监控文件，查看文件发生了什么变化
//若是文件内容变成exit说明对应的视频录制管道需要移除
func MonitorFile() {
	if len(rclist) > 0 && IsMonitorFile == 0 {
		IsMonitorFile = 1
		watcher, err := fsnotify.NewWatcher()
		core.CheckErr(err, "MonitorFile|创建文件夹监控服务：")
		defer watcher.Close()
		done := make(chan bool)
		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op&fsnotify.Write == fsnotify.Write {
						log.Println("modified file:", event.Name)
						filestr := core.Readfile(event.Name)
						if filestr == "exit" {
							filenamearr := strings.Split(event.Name, "\\")
							idstrarr := filenamearr[len(filenamearr)-1]
							idstr := strings.Split(idstrarr, ".")
							for i, v := range rclist {
								vstr := strconv.Itoa(v.Classroomid)
								if vstr == idstr[0] { //如果找到相等的值
									rclist = append(rclist[:i], rclist[i+1:]...) //移除自身
									dbmap := core.InitDb()
									defer dbmap.Db.Close()
									_, selecterr := dbmap.Exec("update videoconfig set CameraState=? where Classroomid=?;", 0, v.Classroomid)
									core.CheckErr(selecterr, "MonitorFile|更新教室内摄像机状态")
									_, selecterr = dbmap.Exec("update classroomcomputerconfig set ComputerState=? where Classroomid=?;", 0, v.Classroomid)
									core.CheckErr(selecterr, "MonitorFile|更新教室内电脑状态")
									break
								}
							}
						}
					}
				case err := <-watcher.Errors:
					core.CheckErr(err, "MonitorFile|文件监控运行时错误：")
				}
			}
		}()
		err = watcher.Add("temporarylogfile/ffmpeglog")
		core.CheckErr(err, "MonitorFile|设置文件夹目录：")
		<-done
	}
}
*/
/*
func ExecCommand(filename string, Classroomid int) {
	//execcmd := exec.Command("C:\\ffmpegwin64\\bin\\ffmpeg.exe", "-i", "rtsp://admin:xywadmin@172.16.100.200:554/cam/realmonitor?channel=1&subtype=1", "-t", "180", "-vcodec", "copy", filename)
	//execcmd := exec.Command("E:\\ffmpegwin64\\bin\\ffmpeg.exe", "-i", "rtsp://admin:lq011189@172.17.92.10:554/cam/realmonitor?channel=1&subtype=1", "-t", "180", "-vcodec", "copy", filename)
	//execcmd := exec.Command("E:\\ffmpegwin64\\bin\\ffmpeg.exe", "-i", "rtsp://admin:xywadmin@172.17.92.67:554/cam/realmonitor?channel=1&subtype=1", "-t", "180", "-vcodec", "copy", filename)
	//execcmd := exec.Command("C:\\ffmpegwin64\\bin\\ffmpeg.exe", "-f", "gdigrab", "-framerate", "5", "-offset_x", "0", "-offset_y", "0", "-video_size", "1364x768", "-i", "desktop", "-t", "2700", "-vcodec", "libx264", "-pix_fmt", "yuv420p", filename)
	//	rc := RecordCmd{Classroomid: Classroomid, Cmd: execcmd}
	//	rclist = append(rclist, rc)
	//	var out bytes.Buffer
	//	var stderr bytes.Buffer
	//	execcmd.Stdout = &out
	//	execcmd.Stderr = &stderr
	//	errrun := execcmd.Run()
	//	if errrun != nil {
	//		fmt.Println(fmt.Sprint(errrun) + ":" + stderr.String())
	//	}
	//	outstr := out.String()
	//	fmt.Println("Result:" + outstr)
	//	for i, v := range rclist {
	//		if v.Classroomid == Classroomid { //如果找到相等的值
	//			rclist = append(rclist[:i], rclist[i+1:]...) //然后在毁尸灭迹
	//		}
	//	}

	cmdstr := "ffmpeg|C:\\ffmpegwin64\\bin\\ffmpeg.exe -i rtsp://admin:xywadmin@172.16.100.200:554/cam/realmonitor?channel=1&subtype=1 -t 180 -vcodec copy " + filename
	var createfilename = "./temporarylogfile/ffmpeglog/" + strconv.Itoa(Classroomid) + ".txt"
	var f *os.File
	defer f.Close()
	var err1 error
	fmt.Println(core.CheckFileIsExist(createfilename))
	if core.CheckFileIsExist(createfilename) { //如果文件存在
		f, err1 = os.OpenFile(createfilename, os.O_APPEND, 0666) //打开文件
	} else {
		f, err1 = os.Create(createfilename) //创建文件
	}
	core.CheckErr(err1, "ExecCommand|获取文件失败：")
	var arrbt []byte
	_, ierr := f.Read(arrbt)
	core.CheckErr(ierr, "ExecCommand|读取录制文件时错误：")
	strbt := string(arrbt)
	if strbt == "exit" || strbt == "q" || strbt == "" {
		//core.CheckErr(err1, "ExecCommand|操作文件：")
		_, err1 := io.WriteString(f, cmdstr) //写入文件(字符串)
		core.CheckErr(err1, "ExecCommand|操作文件：")
		rc := RecordCmd{Classroomid: Classroomid, Cmd: nil}
		rclist = append(rclist, rc)
	}
}
*/

/*
教师点击移动中控获取教室Id
*/
func GetClassClassroomId(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	fmt.Println(string(data))
	core.CheckErr(errs1, "登录数据json格式转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			if lg.Rolestype == 2 {
				rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "GetClassClassroomId", dbmap)
				if rd.Rcode == "1000" {
					rd = actiondataDataAccess.QueryClassClassroomId(lg, dbmap)
				} else {
					rd.Rcode = "1105"
					rd.Reason = "此功能未授权,不可用"
				}
			} else {
				rd.Rcode = "1005"
				rd.Reason = "此角色无法调用此功能，请确认角色"
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

/*
在无Id传入的情况下获取当前老师的状况下的课程章节点到

func GetCurriculumClassroomChapterCentreId(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	fmt.Println(string(data))
	core.CheckErr(errs1, "登录数据json格式转换错误")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			if lg.Rolestype == 2 {
				rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "getcurriculumclassroomchaptercentreid", dbmap)
				if rd.Rcode == "1000" {
					rd = actiondataDataAccess.GetQueryCurriculumClassroomChapterCentreId(lg, dbmap)
				} else {
					rd.Rcode = "1105"
					rd.Reason = "此功能未授权,不可用"
				}
			} else {
				rd.Rcode = "1005"
				rd.Reason = "此角色无法调用此功能，请确认角色"
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
*/
/*
更改学生点到的状态

提交课程班级章节中间表ID:Curriculumclassroomchaptercentreid,学生ID,状态[0:未到,1:已到]

func UpdatePointtos(c *gin.Context) {
	var rd core.Returndata

	var lg viewmodel.GetPointtos
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "updatepointtos", dbmap)
			if rd.Rcode == "1000" {
				rd = actiondataDataAccess.UpdateStudentsPointtos(lg, lg.Studentsid, lg.State, dbmap)
			} else {
				rd.Rcode = "1005"
				rd.Reason = "此功能未授权,不可用"
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
*/
/*
批量更改学生点到的状态
*/
func UpdateListPointtos(c *gin.Context) {
	var rd core.Returndata
	/*
		提交课程班级章节中间表ID:Curriculumclassroomchaptercentreid,学生ID,状态[0:未到,1:已到]
	*/
	var lg core.BasicsToken
	var updata viewmodel.PostUpdatePointtosData
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &updata)
	core.CheckErr(errs1, "修改数据获取失败")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(lg.Rolestype, lg.Usersid, "updatelistpointtos", dbmap)
			if rd.Rcode == "1000" {
				rd = actiondataDataAccess.UpdateListStudentsPointtos(updata, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
根据楼栋查看每间教室的上课状态
请求数据
[Token令牌、用户ID、角色ID，楼栋Id数组]
响应数据
[课程名称，课程ID，章节名称，章节ID，开始时间，结束时间，所在校区，所在楼栋，所在楼层，所在教室,课程状态]
*/
func GetAttendanceQuerylist(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.PostQueryCurriculums
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs2 := json.Unmarshal(data, &pg)
	core.CheckErr(errs2, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "GetAttendanceQuerylist", dbmap)
			if rd.Rcode == "1000" {
				lg.Floorsids = strings.Trim(lg.Floorsids, "|")
				lg.Floorsids = strings.Trim(lg.Floorsids, " ")
				lg.Buildingids = strings.Trim(lg.Buildingids, "|")
				lg.Buildingids = strings.Trim(lg.Buildingids, " ")
				lg.Campusids = strings.Trim(lg.Campusids, "|")
				lg.Campusids = strings.Trim(lg.Campusids, " ")
				rd = actiondataDataAccess.QueryAttendancelist(lg, pg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
获取课表数据
根据角色不同查看的数据也不同
学生：只看班级课程数据
教师：只看安排自己将要上的课程数据
管理者：看所有的数据
可以根据教室ID查询此教室的所有课程
请求数据
[开始时间、结束时间、Token令牌、用户ID、角色ID，教室Id]
响应数据
[课程名称，课程ID，章节名称，章节ID，开始时间，结束时间，所在校区，所在楼栋，所在楼层，所在教室,课程状态]
*/
func GetCurriculumslist(c *gin.Context) {
	_AllowCrossDomain(c)

	var rd core.Returndata
	var lg viewmodel.PostQueryCurriculums
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs2 := json.Unmarshal(data, &pg)
	core.CheckErr(errs2, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "getcurriculumslist", dbmap)
			if rd.Rcode == "1000" {
				rd = actiondataDataAccess.QueryCurriculumsTable(lg, bt, pg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
获取课表数据
根据角色不同查看的数据也不同
学生：只看班级课程数据
教师：只看安排自己将要上的课程数据
管理者：看所有的数据
可以根据教室ID查询此教室的所有课程
请求数据
[开始时间、结束时间、Token令牌、用户ID、角色ID，教室Id]
响应数据
[课程名称，课程ID，章节名称，章节ID，开始时间，结束时间，所在校区，所在楼栋，所在楼层，所在教室,课程状态]
*/
func GetWatchCurriculumslist(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.PostQueryCurriculums
	var pg core.PageData
	var ul viewmodel.Login
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs2 := json.Unmarshal(data, &pg)
	core.CheckErr(errs2, "登录数据json格式转换错误")
	errs2 = json.Unmarshal(data, &ul)
	core.CheckErr(errs2, "登录数据json格式转换错误")

	if errs1 == nil {
		if ul.Loginuser != "" {
			//		tkbt := bt.Token
			//		tkbtarr := strings.Split(string(tkbt), "|")
			//		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			fmt.Println(ul.Loginuser)
			us := usersDataAccess.QueryUsersInfo(ul.Loginuser, dbmap)
			if us.Id > 0 {
				//		rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "getcurriculumslist", dbmap)
				//			if rd.Rcode == "1000" {
				bt := core.BasicsToken{Rolestype: us.Rolesid, Usersid: us.Id}
				rd = actiondataDataAccess.QueryCurriculumsTable(lg, bt, pg, dbmap)
				//			} else {
				//				rd.Rcode = "1105"
				//				rd.Reason = "此功能未授权,不可用"
				//			}
				//		} else {
				//			rd.Rcode = "1002"
				//			rd.Reason = "Token令牌不对"
				//		}
			} else {
				rd.Rcode = "1003"
			}
		} else {
			rd.Rcode = "1002"
			rd.Reason = "数据提交格式错误"
		}
	} else {
		rd.Rcode = "1002"
		rd.Reason = "数据提交格式错误"
	}
	c.JSON(200, rd)
}
func GetCurriculumsinfo(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.PostQueryCurriculums
	//	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	//	errs2 := json.Unmarshal(data, &pg)
	//	core.CheckErr(errs2, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")
	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "GetCurriculumsinfo", dbmap)
			if rd.Rcode == "1000" {
				rd = actiondataDataAccess.QueryCurriculumsinfo(lg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
获取课表数据
根据角色不同查看的数据也不同
学生：只看班级课程数据
教师：只看安排自己将要上的课程数据
管理者：看所有的数据
可以根据教室ID查询此教室的所有课程
请求数据
[开始时间、结束时间、Token令牌、用户ID、角色ID，教室Id]
响应数据
[课程名称，课程ID，章节名称，章节ID，开始时间，结束时间，所在校区，所在楼栋，所在楼层，所在教室,课程状态]
*/
func GetHistoryAttendance(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.PostQueryCurriculums
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs2 := json.Unmarshal(data, &pg)
	core.CheckErr(errs2, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			rd = usersDataAccess.CheckVaild(bt.Rolestype, bt.Usersid, "GetHistoryAttendance", dbmap)
			if rd.Rcode == "1000" {
				rd = actiondataDataAccess.QueryHistoryAttendance(lg, bt, pg, dbmap)
			} else {
				rd.Rcode = "1105"
				rd.Reason = "此功能未授权,不可用"
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

/*
设置或者取消教室收藏记录
*/
func SetOrCancelClassroomCollection(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.PostCollection
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			if actiondataDataAccess.SetOrCancelClassroomCollection(lg, bt, dbmap) {
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "取消失败"
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

/*
获取我关注的课程视频接口
*/
func GetMyAttentionRecord(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryAttentionRecordWhere
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	//data2, _ := ioutil.ReadAll(c.Request.Body)
	errs2 := json.Unmarshal(data, &pg)
	core.CheckErr(errs2, "分页数据转换失败")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			list := actiondataDataAccess.QueryMyAttentionRecord(lg, bt, pg, dbmap)
			if list != nil {
				rd.Result = list
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "查询失败"
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

/*
获取我所在的班级的课程视频接口
*/
func GetMyClassAttentionRecord(c *gin.Context) {
	var rd core.Returndata
	var lg viewmodel.QueryAttentionRecordWhere
	var pg core.PageData
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	//data2, _ := ioutil.ReadAll(c.Request.Body)
	errs2 := json.Unmarshal(data, &pg)
	core.CheckErr(errs2, "分页数据转换失败")
	var bt core.BasicsToken
	errs1 = json.Unmarshal(data, &bt)
	core.CheckErr(errs1, "令牌验证数据转换失败")

	if errs1 == nil {
		tkbt := bt.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			list := actiondataDataAccess.QueryMyClassAttentionRecord(lg, bt, pg, dbmap)
			if list != nil {
				rd.Result = list
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "查询失败"
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

/*
更改我关注的课程信息[添加或取消我的关注]
*/
func SetAttentionRecord(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken //viewmodel.BasicsToken
	var at actiondata.Attentionrecord
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &at)
	core.CheckErr(errs1, "关注数据转换失败")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			errs1 := actiondataDataAccess.SetAttentionRecord(&at, dbmap)
			if errs1 == nil {
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "查询失败"
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

/*
更改我关注的课程信息[添加或取消我的关注]
*/
func IsAttentionRecordOk(c *gin.Context) {
	var rd core.Returndata
	var lg core.BasicsToken //viewmodel.BasicsToken
	var at actiondata.Attentionrecord
	data, _ := ioutil.ReadAll(c.Request.Body)
	errs1 := json.Unmarshal(data, &lg)
	core.CheckErr(errs1, "登录数据json格式转换错误")
	errs1 = json.Unmarshal(data, &at)
	core.CheckErr(errs1, "关注数据转换失败")
	if errs1 == nil {
		tkbt := lg.Token
		tkbtarr := strings.Split(string(tkbt), "|")
		if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(lg.Usersid) && tkbtarr[1] == strconv.Itoa(lg.Rolestype)) { //判断账号是否正确
			dbmap := core.InitDb()
			defer dbmap.Db.Close()
			b := actiondataDataAccess.IsAttentionRecordOk(&at, dbmap)
			if b {
				rd.Rcode = "1000"
				rd.Reason = "成功"
			} else {
				rd.Rcode = "1001"
				rd.Reason = "查询失败"
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
