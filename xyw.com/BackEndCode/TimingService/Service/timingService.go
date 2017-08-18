package timingService

import (
	//	"TimingService/DataAccess"
	"TimingService/Model"
	"TimingService/Viewmodel"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	//	"math"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
	xconfig "xutils/xconfig"
	core "xutils/xcore"
	"xutils/xtime"

	"gopkg.in/gorp.v1"
)

var cf xconfig.Config
var cmdch = make(chan viewmodel.TTScanDevice) //发送命令消息队列
//后端第一次加入数据库的数据更新扫描间隔时间
var addUpdateTime int64 = 10
var execUpdateTime int64 = 60
var countdbmap *gorp.DbMap

func RunUpdateDate() { //更改下次服务运行时间
	//不管是新添加的还是已经运行过的，每隔一定时间，更新一个下次运行时间
	//	zndxdb := core.InitDb()
	//	defer zndxdb.Db.Close()
	Ticker := time.NewTicker(time.Second * time.Duration(addUpdateTime))      //定时10秒执行一下
	execTicker := time.NewTicker(time.Second * time.Duration(execUpdateTime)) //定时60秒执行一下
	for {
		select {
		case <-Ticker.C:
			log.Println("countdbmap.Db.Stats().OpenConnections:", countdbmap.Db.Stats().OpenConnections)
			var ttlist []taskModel.TimedTask //加载刚刚添加的任务数据
			_, listerr := countdbmap.Select(&ttlist, "select * from TimedTask where TaskState=0 and TaskIsOpen=1 and ExecBeginDate='' and ExecEndDate=''")
			log.Println("监测是否有刚刚添加的，并且开始使用的任务数据", len(ttlist))
			if listerr == nil {
				if len(ttlist) > 0 {
					//					dbmap := core.InitDb()
					log.Println("countdbmap.Db.Stats().OpenConnections1:", countdbmap.Db.Stats().OpenConnections)
					//					defer dbmap.Db.Close()
					//					dbmap.AddTableWithName(taskModel.TimedTask{}, "timedtask").SetKeys(true, "TaskId")
					nowtime := time.Now()
					for _, v := range ttlist {
						nowdatastr := time.Now().Format("2006-01-02") + " " + v.TimePoint + ":00" //设置执行时间
						nowDay := xtime.GetChainDay(nowtime)
						if v.RepeatType == "自定义" {
							if !strings.Contains(v.RepeatValue, nowDay) { //不等于今天的话
								continue //跳过此条数据 不予执行
							}
						}
						v.ExecBeginDate = nowdatastr
						_, err := countdbmap.Exec("update TimedTask set ExecBeginDate=? where TaskId=?", v.ExecBeginDate, v.TaskId)
						core.CheckErr(err, " 执行失败:")
						//						trans, err := dbmap.Begin()
						//						_, err = trans.Exec("update TimedTask set ExecBeginDate=? where TaskId=?", v.ExecBeginDate, v.TaskId)
						//						//						err = timingDataAccess.UpdateTimedTask(&v, trans) //更新数据
						//						if err != nil {
						//							trans.Rollback()
						//							core.CheckErr(err, fmt.Sprintf("%+v", v)+" \n 执行失败:["+err.Error()+"]")
						//							continue //跳过此条数据
						//						} else {
						//							err := trans.Commit()
						//							if err != nil {
						//								core.CheckErr(err, fmt.Sprintln(nowtime)+" 执行失败:["+err.Error()+"]")
						//							}
						//						}
					}

					//					dbmap.Db.Close()
				}
			} else {
				log.Println("定时任务服务模块，启动时错误：", listerr)
			}
		case <-execTicker.C:
			var ttlist []taskModel.TimedTask //加载已经运行过，需要重复运行的数据
			_, listerr := countdbmap.Select(&ttlist, "select * from TimedTask where TaskState=0 and TaskIsOpen=1 and ExecEndDate!='' and ExecEndDate<=now() and to_days(ExecBeginDate)!=to_days(now())")
			log.Println("检测到还有需要重复运行的数据", len(ttlist))
			if listerr == nil {
				if len(ttlist) > 0 {
					//					dbmap := core.InitDb()
					//					defer dbmap.Db.Close()
					log.Println("countdbmap.Db.Stats().OpenConnections2:", countdbmap.Db.Stats().OpenConnections)
					//					dbmap.AddTableWithName(taskModel.TimedTask{}, "timedtask").SetKeys(true, "TaskId")

					nowtime := time.Now()
					for _, v := range ttlist {
						nowdatastr := time.Now().Format("2006-01-02") + " " + v.TimePoint + ":00" //设置执行时间
						nowDay := xtime.GetChainDay(nowtime)
						if v.RepeatType == "自定义" {
							if !strings.Contains(v.RepeatValue, nowDay) { //不等于今天的话
								continue //跳过此条数据 不予执行
							}
						}
						v.ExecBeginDate = nowdatastr
						_, err := countdbmap.Exec("update TimedTask set ExecBeginDate=? where TaskId=?", v.ExecBeginDate, v.TaskId)
						core.CheckErr(err, " 执行失败:")
						//						trans, err := dbmap.Begin()
						//						_, err = trans.Exec("update TimedTask set ExecBeginDate=? where TaskId=?", v.ExecBeginDate, v.TaskId)
						//						//						err = timingDataAccess.UpdateTimedTask(&v, trans) //更新数据
						//						if err != nil {
						//							trans.Rollback()
						//							core.CheckErr(err, fmt.Sprintf("%+v", v)+" \n 执行失败:["+err.Error()+"]")
						//							continue //跳过此条数据
						//						} else {
						//							err := trans.Commit()
						//							if err != nil {
						//								core.CheckErr(err, fmt.Sprintln(nowtime)+" 执行失败:["+err.Error()+"]")
						//							}
						//						}
					}

					//					dbmap.Db.Close()
				}
			} else {
				log.Println("定时任务服务模块，启动时错误：", listerr)
			}
		}
	}
	//	zndxdb.Db.Close()
}

//后台运行的定时检测服务
func RunBackgroundServer() {
	log.Println("定时服务后台运行服务：开启")
	//	zndxdb := core.InitDb()
	log.Println("countdbmap.Db.Stats().OpenConnections2:", countdbmap.Db.Stats().OpenConnections)
	//	defer zndxdb.Db.Close()
	Ticker := time.NewTicker(time.Second * (60)) //定时60秒执行一下
	for {
		select {
		case <-Ticker.C:
			//需要处理的事物
			//一、查询定时任务表，看看是否有需要处理的定时任务--例如[定时关灯/定时开灯]
			nowtimestr := time.Now().Add(-(time.Minute * 1)).Format("2006-01-02 15:04:00")
			var ttlist []taskModel.TimedTask //开始执行定时任务
			_, listerr := countdbmap.Select(&ttlist, "select * from TimedTask where TaskState=0 and TaskIsOpen=1 and to_days(ExecBeginDate)=to_days(now()) and ExecBeginDate>=? order by ExecBeginDate", nowtimestr)
			log.Println("开始执行定时任务，检查有多少条任务需要执行", len(ttlist))
			if listerr == nil {
				if len(ttlist) > 0 {
					//					dbmap := core.InitDb() //链接数据库
					//					defer dbmap.Db.Close()
					//					dbmap.AddTableWithName(taskModel.TimedTask{}, "timedtask").SetKeys(true, "TaskId")
					//

					for i, v := range ttlist {
						log.Printf("开始执行定时任务，执行%d [%+v]", i, v)
						NowTime := time.Now()               //获取准确的当前时间
						ebt := xtime.Parse(v.ExecBeginDate) //获取设置时间例如2017-02-23 11:55:00并转换成time.Time
						nt := xtime.Parse(NowTime.Format("2006-01-02 15:04:00"))
						absint := float64(nt.Unix() - ebt.Unix())
						//						fmt.Println("absint:", absint)
						if absint <= 60 && absint >= 0 { //只要是在一分钟内,满足条件就开始执行
							var ttsd viewmodel.TTScanDevice
							data := []byte(v.TaskContent)
							errs1 := json.Unmarshal(data, &ttsd)
							fmt.Println("errs1:", errs1)
							if errs1 != nil {
								core.CheckErr(errs1, xtime.NowString()+" 执行失败:["+errs1.Error()+"]")
								return
							}
							go func() {
								cmdch <- ttsd //将值放到通道中
							}()
							switch v.TaskType {
							case 0: //单次任务
								v.TaskState = 2
								v.TaskIsOpen = 0
								v.TaskExecNum = 0
							case 1: //多次任务
								v.TaskExecNum = v.TaskExecNum - 1
								if v.TaskExecNum == 0 {
									v.TaskState = 2
									v.TaskIsOpen = 0
								}
							case 2: //循环任务
								//暂时未想到需要更改的数据
							}
							v.ExecEndDate = xtime.NowString()
							_, err := countdbmap.Exec("update TimedTask set ExecEndDate=?,TaskState=?,TaskIsOpen=?,TaskExecNum=? where TaskId=?", v.ExecEndDate, v.TaskState, v.TaskIsOpen, v.TaskExecNum, v.TaskId)
							core.CheckErr(err, " 执行失败:")
							//							trans, err := dbmap.Begin() //开启事物
							//							_, err = trans.Exec("update TimedTask set ExecEndDate=?,TaskState=?,TaskIsOpen=?,TaskExecNum=? where TaskId=?", v.ExecEndDate, v.TaskState, v.TaskIsOpen, v.TaskExecNum, v.TaskId)
							//							//							err = timingDataAccess.UpdateTimedTask(&v, trans) //更新数据
							//							if err != nil {
							//								trans.Rollback()
							//								core.CheckErr(err, fmt.Sprintf("%+v", v)+" \n 执行失败:["+err.Error()+"]")
							//								continue //跳过此条数据
							//							} else {
							//								err := trans.Commit()
							//								if err != nil {
							//									core.CheckErr(err, " 执行失败:["+err.Error()+"]")
							//								}
							//							}
						} else {
							continue
						}
					}
					//					dbmap.Db.Close()
				}
			} else {
				log.Println("定时服务模块，启动时错误：", listerr)
			}
		}
	}
	//	zndxdb.Db.Close()
}

//后台运行的定时检测数据的完整性和有效性
func RunBackgroundCheckValid() {

	Ticker := time.NewTicker(time.Second * (60)) //定时60分钟执行一下
	for {
		select {
		case <-Ticker.C:
			//需要处理的事物
			//			zndxdb := core.InitDb()
			//			defer zndxdb.Db.Close()
			//后台运行的定时检测数据的完整性和有效性
			_, listerr := countdbmap.Exec("update TimedTask set ExecBeginDate='',ExecEndDate='' where TaskState=0 and TaskIsOpen=1 and to_days(ExecBeginDate)!=to_days(now()) and ExecBeginDate!='' and ExecEndDate=''")
			if listerr != nil {
				core.CheckErr(listerr, xtime.NowString()+" 执行失败:["+listerr.Error()+"]")
			}
		}
	}
}

//命令执行队列
func RunExecCmd() {
	//所有的命令执行都是在此方法中
	for {
		ssh := <-cmdch
		sp, sperr := json.Marshal(&ssh.Para)
		fmt.Println("sperr", sperr)
		resp, err := http.Post(ssh.TTUrl, "application/json", strings.NewReader(string(sp)))
		if err == nil {
			respdata, resperr := ioutil.ReadAll(resp.Body)
			core.CheckErr(resperr, "timingService|接受发送udp服务器响应回来的数据"+string(respdata))
		} else {
			core.CheckErr(err, "timingService|RunExecCmd|发送udp请求")
		}
	}
}

func init() {

	cf.InitConfig("./config.ini")
	addUpdateTimestr := cf.Read("timingset", "addUpdateTime")
	if addUpdateTimestr != "" { //判断是否存在设置
		addUpdateTimeint, err := strconv.ParseInt(addUpdateTimestr, 10, 0)
		if err != nil || addUpdateTimeint <= 0 {
			log.Printf("新数据扫描间隔时间设置无效，系统自动默认间隔时间10秒")
		} else {
			addUpdateTime = addUpdateTimeint
		}
	}
	execUpdateTimestr := cf.Read("timingset", "execUpdateTime")
	if execUpdateTimestr != "" { //判断是否存在设置
		execUpdateTimeint, err := strconv.ParseInt(execUpdateTimestr, 10, 0)
		if err != nil || execUpdateTimeint <= 0 {
			log.Printf("新数据扫描间隔时间设置无效，系统自动默认间隔时间60秒")
		} else {
			execUpdateTime = execUpdateTimeint
		}
	}
	countdbmap = core.InitDb()
	//	defer countdbmap.Db.Close()
	//	RunUpdateDate()
	//	RunBackgroundServer()
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
