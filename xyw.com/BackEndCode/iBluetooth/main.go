package main

import (
	//	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	//	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zndx2/DataAccess/deviceDataAccess"
	"zndx2/model/core"
	//	"zndx2/model/deviceModel"

	"github.com/bugra/matrix"
	//	mx2 "github.com/tantanbei/matrix"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
)

var GlobalConfig Config
var ch = make(chan ScanDevice)
var bttmap = map[string]Bluetooth{}

var cmdch = make(chan TTScanDevice) //发送命令消息队列

func main() {

	GlobalConfig.InitConfig("./config.ini")
	selectdb := InitIbeaconDb()
	defer selectdb.Db.Close()
	zndxdb := InitZNDXDb()
	defer zndxdb.Db.Close()
	inertdb := InitIbeaconDb()
	defer inertdb.Db.Close()

	go func() {
		var bttlist []Bluetooth
		zndxdb.Select(&bttlist, "select * from `bluetooth`")
		for _, v := range bttlist {
			key := strconv.Itoa(v.Major) + "-" + strconv.Itoa(v.Minor)
			bttmap[key] = v
		}
	}()
	go AddOrbit() //一直处于连接状态的插入用户轨迹数据

	go RunExecCmd() //命令执行队列

	go RunBackgroundServer() //后台运行的定时检测服务
	//利用web方式来接收发给客户端信息
	sqlDevicemove := "select count(*) from devicemovetb where 1=1 "
	sqlScanDevice := "select count(*)as Ct from scandevice where DATE_FORMAT(NowDate,'%H')=DATE_FORMAT(now(),'%H') "
	sqlBluetooth := "select * from bluetooth where 1=1 "
	sqlqueryclassroom := "select * from classroomdetailsnums where Classroomid=? and Usersid=? and to_days(Createdate)=to_days(now())and DATE_FORMAT(Createdate,'%H')=DATE_FORMAT(now(),'%H');"
	go RunBackgroundBluetoothServer() //后台运行的蓝牙服务

	r := gin.Default()
	r.POST("/receiveibluetooth", func(c *gin.Context) { //接口,收集需要发出去的命令
		t := time.Now().UnixNano()
		fmt.Println("time:", t)
		var rd core.Returndata
		var lg ScanDeviceList
		data, _ := ioutil.ReadAll(c.Request.Body)
		errs1 := json.Unmarshal(data, &lg)
		core.CheckErr(errs1, "登录数据json格式转换错误")
		if errs1 == nil {
			bt := lg.UserOAuth
			tkbt := bt.Token
			tkbtarr := strings.Split(string(tkbt), "|")
			if len(tkbtarr) == 3 && (tkbtarr[2] == strconv.Itoa(bt.Usersid) && tkbtarr[1] == strconv.Itoa(bt.Rolestype)) { //判断账号是否正确
				//				inertdb := InitIbeaconDb()
				//				defer inertdb.Db.Close()
				go SaveReceiveIbluetooth(string(data), inertdb)
				var rpoint [][]float64
				//第一步验证有效的蓝牙基站，排除非有效的蓝牙基站后，上报的蓝牙基站集合是否大于等于3
				if len(lg.ScandeviceList) > 0 {
					//					newsdlist := make([]ScanDevice, 0)
					var newsdlist []ScanDevice
					var basescandevice ScanDevice //定位最近的蓝牙基站
					for newsdlistindex, v := range lg.ScandeviceList {
						key := strconv.Itoa(v.Major) + "-" + strconv.Itoa(v.Minor)
						if nv, isok := bttmap[key]; isok {
							v.LocationX = nv.LocationX
							v.LocationY = nv.LocationY
							v.LocationZ = nv.LocationZ
							v.FloorsId = nv.FloorsId
							if newsdlistindex <= 2 {
								newsdlist = append(newsdlist, v)
							}
							if basescandevice.Major == 0 && basescandevice.Minor == 0 && v.CurrentDistance > 0 {
								basescandevice = v
							} else if v.CurrentDistance < basescandevice.CurrentDistance && v.CurrentDistance > 0 {
								basescandevice = v
							}
						}
					}
					//第二步排除是否有位置距离是在0.5米以内的--如果有则人员x,y点就直接为蓝牙基站的位置,没有就执行下一步
					if basescandevice.CurrentDistance <= 1 {
						ts := time.Now().UnixNano()
						x := basescandevice.LocationX
						y := basescandevice.LocationY
						var rpointone []float64
						var rpointtwo []float64
						rpointone = append(rpointone, x)
						rpointtwo = append(rpointtwo, y)
						rpoint = append(rpoint, rpointone)
						rpoint = append(rpoint, rpointtwo)
						obs := Orbit{Usersid: lg.UserOAuth.Usersid, X: x, Y: y, TimeStamp: ts}
						go func() {
							ob <- obs
						}()
					} else if len(newsdlist) == 3 { //第三步执行RSSI算法[取距离最近的三条]
						rpoint = Calculate(newsdlist)
						ts := time.Now().UnixNano()
						x := rpoint[0][0]
						y := rpoint[1][0]
						obs := Orbit{Usersid: lg.UserOAuth.Usersid, X: x, Y: y, TimeStamp: ts}
						go func() {
							ob <- obs
						}()
					}

					//				}
					//				if len(lg.ScandeviceList) > 0 { //判断获取到了蓝牙信息为准
					//					go Calculate(lg)

					if basescandevice.CurrentDistance != 0 {
						basescandevice.Usersid = lg.UserOAuth.Usersid
						basescandevice.Uniqueid = lg.Devicemove.Uniqueid
						countint1, counterr := selectdb.SelectInt(sqlDevicemove + " and Uniqueid='" + lg.Devicemove.Uniqueid + "'")
						core.CheckErr(counterr, "未找到相关数据"+sqlDevicemove)
						if countint1 == 0 && counterr == nil {
							go inertdb.Insert(&lg.Devicemove)
						}
						type CtStruct struct {
							Ct int
						}
						var ctarr []CtStruct
						_, counterr = selectdb.Select(&ctarr, sqlScanDevice+" and to_days(NowDate)=to_days(now()) and Uniqueid='"+lg.Devicemove.Uniqueid+"' group by Usersid;")
						core.CheckErr(counterr, "未找到相关数据"+sqlScanDevice)
						//						if len(ctarr) < 2 && counterr == nil {

						var bluet Bluetooth
						execsqlBluetooth := sqlBluetooth
						if basescandevice.Mac != "" {
							execsqlBluetooth = execsqlBluetooth + " and Mac='" + basescandevice.Mac + "'"
						}
						if basescandevice.Major > 0 {
							execsqlBluetooth = execsqlBluetooth + " and Major=" + strconv.Itoa(basescandevice.Major) + " and Minor=" + strconv.Itoa(basescandevice.Minor)
						}
						zndxdb.SelectOne(&bluet, execsqlBluetooth)
						if bluet.Id > 0 {
							basescandevice.NowDate = time.Now().Format("2006-01-02 15:04:05")
							basescandevice.Classroomsid = bluet.Classroomsid
							//							go SaveClassroomdetailsnums(sqlqueryclassroom, basescandevice.Mac, basescandevice.Classroomsid, basescandevice.Usersid, zndxdb)
							go func() {
								var crdn Classroomdetailsnums
								zndxdb.SelectOne(&crdn, sqlqueryclassroom, basescandevice.Classroomsid, basescandevice.Usersid)
								rpXY := fmt.Sprintf("%v", rpoint)
								crdn.Xy = rpXY
								if len(rpoint) > 0 {
									crdn.X = rpoint[0][0]
									crdn.Y = rpoint[1][0]
								}
								if crdn.Id > 0 {
									crdn.Closedate = time.Now().Format("2006-01-02 15:04:05")
									crdn.Positionnum = crdn.Positionnum + 1
									zndxdb.Update(&crdn)
								} else {
									crdn.Classroomid = basescandevice.Classroomsid
									crdn.Createdate = time.Now().Format("2006-01-02 15:04:05")
									crdn.Closedate = time.Now().Format("2006-01-02 15:04:05")
									crdn.Usersid = basescandevice.Usersid
									crdn.Mac = basescandevice.Mac
									crdn.Positionnum = 1
									zndxdb.Insert(&crdn)
								}
							}()
							counterr = inertdb.Insert(&basescandevice)
							fmt.Println("counterr:", counterr)
							go func() {
								ch <- basescandevice //将值放到通道中
							}()
							rd.Rcode = "1000"
							rd.Result = bluet.Geographical
						} else {
							rd.Rcode = "1006"
							rd.Reason = "未找到此蓝牙数据"
						}
						//						} else {
						//							rd.Rcode = "1015"
						//							rd.Reason = "此终端已重复登录其他账号"
						//						}
					} else {
						rd.Rcode = "1025"
						rd.Reason = "数据提交不正确，采集数据出错"
					}
				} else {
					rd.Rcode = "1035"
					rd.Reason = "未获取到蓝牙定位信息"
				}
			} else {
				rd.Rcode = "1002"
				rd.Reason = "Token令牌不对"
			}
		} else {
			rd.Rcode = "1005"
			rd.Reason = "数据提交格式错误1"
		}
		et := time.Now().UnixNano()
		fmt.Println("endtime:", et)
		c.JSON(200, rd)
	})
	r.POST("/GetOrbit", func(c *gin.Context) { //查询单个的人某一时段轨迹记录

	})
	r.POST("/GetBLElist", func(c *gin.Context) { //查询所有的蓝牙基础信息反馈出去
		var rd core.Returndata
		var BLElist []Bluetooth
		zndxdb.Select(&BLElist, "select * from Bluetooth")
		rd.Rcode = "1000"
		rd.Result = BLElist
		c.JSON(200, rd)
	})
	r.Run(":8025")
}

//后台运行的蓝牙服务
func RunBackgroundBluetoothServer() {
	zndxdb := InitZNDXDb()
	defer zndxdb.Db.Close()
	for {
		v := <-ch //从通道内将值取出

		if v.Major == 10810 && v.Minor == 2 {
			Idstr, _ := zndxdb.SelectStr("SELECT group_concat(Id) FROM building;")
			list := deviceDataAccess.GetClassroomStatusListPuls(Idstr, zndxdb)
			fmt.Println("Idstr:", Idstr)
			for _, lv := range list {
				if lv.ClassroomId == v.Classroomsid && lv.HaveRun <= 0 {
					ssh := TTScanDevice{}
					ssh.TTUrl = "http://192.168.0.201:8090/device/node/control/switch/on/room"
					ssh.SSW.UserID = strconv.Itoa(v.Usersid)
					ssh.SSW.RoomID = strconv.Itoa(v.Classroomsid)
					go func() {
						cmdch <- ssh //将值放到通道中
					}()
				}
			}
		}
		fmt.Println(v)
	}
}

//命令执行队列
func RunExecCmd() {
	//所有的命令执行都是在此方法中
	for {
		ssh := <-cmdch
		fmt.Printf("ssh:%+v \n", ssh)
		bte, _ := json.Marshal(&ssh.SSW)
		resp, err := http.Post(ssh.TTUrl, "application/json", strings.NewReader(string(bte)))
		if err == nil {
			respdata, resperr := ioutil.ReadAll(resp.Body)
			core.CheckErr(resperr, "BeginVideo|接受发送udp服务器响应回来的数据")
			fmt.Println(string(respdata))
		} else {
			core.CheckErr(err, "BeginVideo|commandsendlog|发送udp请求")
		}
	}
}

//后台运行的定时检测服务
func RunBackgroundServer() {

	selectdbs := InitIbeaconDb()
	defer selectdbs.Db.Close()
	zndxdb := InitZNDXDb()
	defer zndxdb.Db.Close()
	Ticker := time.NewTicker(time.Second * (5)) //定时60秒执行一下
	for {
		select {
		case <-Ticker.C:
			Idstr, _ := zndxdb.SelectStr("SELECT group_concat(Id) FROM building;")
			list := deviceDataAccess.GetClassroomStatusListPuls(Idstr, zndxdb)
			fmt.Println("Idstr:", Idstr)
			//			//服务器响应客户端的数据
			//			type ResultData struct {
			//				Page deviceModel.PageData              //分页数据
			//				Data []deviceModel.ClassroomStatusData //具体的数据
			//			}
			if list != nil && len(list) > 0 {
				for _, v := range list {
					if v.HaveRun > 0 { //这间教室有设备在运行
						fmt.Println("HaveRun:", v)
						//						go func() {
						//							zndxdb := InitZNDXDb()
						//							defer zndxdb.Db.Close()
						nowtime := time.Now().Add(-(time.Second * 5)).Format("2006-01-02 15:04:05") //查询30分钟以内有没有人进入
						fmt.Println("nowtime:", nowtime, "v.ClassroomId:", v.ClassroomId)
						countcmd, counterr := selectdbs.SelectInt("select count(*) from scandevice where NowDate>=? and Classroomsid=? and CurrentDistance<=3", nowtime, v.ClassroomId)
						//						countcmd, counterr := zndxdb.SelectInt("select count(*) from Classroomdetailsnums where Closedate>=? and Classroomid=?", nowtime, v.ClassroomId)
						fmt.Println("countcmd:", countcmd)
						fmt.Println("counterr:", counterr)
						if counterr == nil {
							if countcmd > 0 { //判断是否有人在教室内 有 则不操作
								continue
							} else { //无 则操作关闭设备
								ttsd := TTScanDevice{}
								ttsd.SSW.RoomID = strconv.Itoa(v.ClassroomId)
								ttsd.SSW.UserID = "9"
								ttsd.SSW.Params = ""
								ttsd.TTUrl = "http://192.168.0.201:8090/device/node/control/switch/off/room"
								go func() {
									fmt.Println(ttsd)
									cmdch <- ttsd //将数值塞进通道中去
								}()
							}
						} else {
							continue
						}
						//						}()
					}
				}
			}
			//二、查询教室内根据一段时间内蓝牙定位是否还存在人
			/*
				1、以教室内的设备运行为单位来寻找教室内是否有在运行的设备{http://192.168.0.201:8080/basicset/getclassroominfo?id=343}
				2、教室内如果有运行的设备则查看是否有未执行完的定时任务，如果没有则判断蓝牙上报的信息中是否还有人在教室内
				3、如果没有则添加一条任务命令，等待任务机制去处理发送
			*/
		}
	}
}

//用户轨迹通道
var ob = make(chan Orbit)

//后台运行用户轨迹添加服务
func AddOrbit() {
	inertdb := InitIbeaconDb()
	defer inertdb.Db.Close()
	inertdb.AddTableWithName(Orbit{}, "orbit")
	for {
		o := <-ob
		ersr := inertdb.Insert(&o)
		fmt.Println("ersr:", ersr)
	}
}

//RSSI算法
func Calculate(sdl []ScanDevice) (rdata [][]float64) {
	/*距离数组*/
	baseNum := len(sdl)
	fmt.Println("--------------------------------begin")
	fmt.Println("******************************")
	fmt.Println(sdl)
	fmt.Println("******************************")
	/*对应的权值*/
	//	weight := 0.0
	//	for i := 0; i < 3; i++ {
	//		weight += (1.0 / sdl[i].CurrentDistance)
	//	}
	baseNums := int(baseNum - 1)
	a := make([][]float64, 0)
	b := make([][]float64, 0)
	LastScandevice := sdl[baseNums]
	for i, v := range sdl {
		liba := []float64{(v.LocationX - LastScandevice.LocationX), (v.LocationY - LastScandevice.LocationY)}
		a = append(a, liba)
		if i == 1 {
			break
		}
	}
	for ib, vb := range sdl {
		libb := []float64{(math.Pow(vb.LocationX, 2) - math.Pow((LastScandevice.LocationX), 2) + math.Pow(vb.LocationY, 2) - math.Pow((LastScandevice.LocationY), 2) + math.Pow(LastScandevice.CurrentDistance, 2) - math.Pow(vb.CurrentDistance, 2))}
		b = append(b, libb)
		if ib == 1 {
			break
		}
	}
	b1 := b
	a1 := a
	a2 := matrix.Transpose(a1)
	tmpMatrix1 := Times(a1, a2)

	//	tmpMatrix1s := mx2.NewMatrix(tmpMatrix1)
	//	reTmpMatrix1s := tmpMatrix1s.Negative()
	//	reTmpMatrix1 := reTmpMatrix1s.GetRealData()

	reTmpMatrix1 := Inverse(tmpMatrix1)
	tmpMatrix2 := Times(a2, reTmpMatrix1)
	resultMatrix := Times(b1, tmpMatrix2)
	var resultMatrixs [][]float64
	for p := 0; p < len(resultMatrix); p++ {
		var iresultMatrix []float64
		for i := 0; i < len(resultMatrix[p]); i++ {
			s := float64(resultMatrix[p][i] / 2) //math.Sqrt(resultMatrix[p][i])
			//			s := float64(resultMatrix[p][i] * weight)
			s = math.Abs(s)
			iresultMatrix = append(iresultMatrix, s)
		}
		resultMatrixs = append(resultMatrixs, iresultMatrix)
	}
	//	fmt.Println("resultMatrixs:", resultMatrixs, "weight:", weight)
	fmt.Println("resultMatrixs:", resultMatrixs, "weight:", 0.5)
	fmt.Println("--------------------------------end")
	return resultMatrixs
}
func SaveReceiveIbluetooth(str string, inertdb *gorp.DbMap) {
	var rib ReceiveIbluetooth
	rib.ReceiveDate = time.Now().Format("2006-01-02 15:04:05")
	rib.ReceiveStr = str
	inertdb.Insert(&rib)
}
func SaveClassroomdetailsnums(sqlqueryclassroom, Mac string, Classroomsid, Usersid int, zndxdb *gorp.DbMap) {
	var crdn Classroomdetailsnums
	zndxdb.SelectOne(&crdn, sqlqueryclassroom, Classroomsid, Usersid)
	if crdn.Id > 0 {
		crdn.Closedate = time.Now().Format("2006-01-02 15:04:05")
		crdn.Positionnum = crdn.Positionnum + 1
		zndxdb.Update(&crdn)
	} else {
		crdn.Classroomid = Classroomsid
		crdn.Createdate = time.Now().Format("2006-01-02 15:04:05")
		crdn.Closedate = time.Now().Format("2006-01-02 15:04:05")
		crdn.Usersid = Usersid
		crdn.Mac = Mac
		crdn.Positionnum = 1
		zndxdb.Insert(&crdn)
	}
}
func Times(b [][]float64, a [][]float64) (result [][]float64) {
	if len(b) != len(a[0]) {
		fmt.Println("len(b)--len(a[0])", len(b), len(a[0]))
		//		return
	}
	fmt.Println("len(b)--len(b[0])--len(a)--len(a[0])", len(b), len(b[0]), len(a), len(a[0]))
	var libresult [][]float64
	Bcolj := make([]float64, 0)
	for j := 0; j < len(b[0]); j++ {
		for k := 0; k < len(a[0]); k++ {
			if j > 0 {
				Bcolj[k] = b[k][j]
			} else {
				Bcolj = append(Bcolj, b[k][j])
			}
		}
		//		fmt.Println("Bcolj:", Bcolj)
		for i := 0; i < len(a); i++ {
			var libresult1 []float64
			if j > 0 {
				libresult1 = libresult[i]
			}
			Arowi := a[i]
			var s float64
			s = 0
			for k := 0; k < len(a[0]); k++ {
				s += float64(float64(Arowi[k]) * float64(Bcolj[k]))
			}
			libresult1 = append(libresult1, s)
			if j > 0 {
				libresult[i] = libresult1
			} else {
				libresult = append(libresult, libresult1)
			}
			//			fmt.Println("libresult:", libresult)
		}
	}
	return libresult
}
func Inverse(a [][]float64) (result [][]float64) {
	ai := Identity(len(a), len(a))
	return Solve(a, ai)
}
func Identity(m int, n int) (result [][]float64) {
	for i := 0; i < m; i++ {
		//var rt []float64
		rt := make([]float64, 0)
		for j := 0; j < n; j++ {
			if i == j {
				rt = append(rt, 1)
			} else {
				rt = append(rt, 0)
			}
		}
		result = append(result, rt)
	}
	return result
}
func Solve(olda [][]float64, newb [][]float64) (result [][]float64) {
	if len(olda) == len(olda[0]) {
		l := LUD{}
		l.LUDecomposition(olda)
		return l.LUDSolve(newb)
		//		return LUDSolve(newb)
	} else {
		q := QRD{}
		q.QRDecomposition(olda)
		return q.QRDsolve(newb)
		//		return LUDSolve(newb)
	}
	return newb
}

type LUD struct {
	LU      [][]float64
	M       int
	N       int
	Pivsign int
	Piv     []int
}
type QRD struct {
	QR    [][]float64
	M     int
	N     int
	Rdiag []float64
}

func (o *LUD) LUDecomposition(a [][]float64) {
	o.LU = a
	o.M = len(a)
	o.N = len(a[0])
	piv := make([]int, 0)
	o.Piv = piv
	for i := 0; i < o.M; i++ {
		o.Piv = append(o.Piv, i)
		//o.Piv[i] = i
	}

	o.Pivsign = 1
	var LUrowi []float64
	LUcolj := make([]float64, o.M)

	for j := 0; j < o.N; j++ {
		for i := 0; i < o.M; i++ {
			LUcolj[i] = o.LU[i][j]
		}
		for i := 0; i < o.M; i++ {
			LUrowi = o.LU[i]
			kmax := math.Min(float64(i), float64(j))
			var s float64
			s = 0
			for k := 0; k < int(kmax); k++ {
				s += LUrowi[k] * LUcolj[k]
			}
			LUcolj[i] = LUcolj[i] - s

			LUrowi[j] = LUcolj[i]
		}
		p := j
		for i := j + 1; i < o.M; i++ {
			luji := math.Abs(LUcolj[i])
			lujp := math.Abs(LUcolj[p])
			if luji > lujp {
				p = i
			}
		}
		if p != j {
			for k := 0; k < o.N; k++ {
				t := o.LU[p][k]
				o.LU[p][k] = o.LU[j][k]
				o.LU[j][k] = t
			}
			k := o.Piv[p]
			o.Piv[p] = o.Piv[j]
			o.Piv[j] = k
			o.Pivsign = -o.Pivsign
		}
		if j < o.M && o.LU[j][j] != 0.0 {
			for i := j + 1; i < o.M; i++ {
				o.LU[i][j] /= o.LU[j][j]
			}
		}
	}
}

func (o *LUD) LUDSolve(b [][]float64) (result [][]float64) {
	if len(b) != o.M {
		fmt.Println("len(b) != o.M:", len(b), o.M)
	}
	if !o.IsNonsingular() {
		fmt.Println("o.IsNonsingular():", o.IsNonsingular())
	}

	nx := len(b)
	Xmat := make([][]float64, 0)
	for i := 0; i < len(b); i++ {
		Xmatb := make([]float64, 0)
		for j := 0; j < nx; j++ {
			Xmatb = append(Xmatb, b[o.Piv[i]][j])
		}
		Xmat = append(Xmat, Xmatb)
	}

	for k := 0; k < o.N; k++ {
		for i := k + 1; i < o.N; i++ {
			for j := 0; j < nx; j++ {
				Xmat[i][j] -= Xmat[k][j] * o.LU[i][k]
			}
		}
	}
	for k := o.N - 1; k >= 0; k-- {
		for j := 0; j < nx; j++ {
			Xmat[k][j] /= o.LU[k][k]
		}
		for i := 0; i < k; i++ {
			for j := 0; j < nx; j++ {
				Xmat[i][j] -= Xmat[k][j] * o.LU[i][k]
			}
		}
	}
	return Xmat
}

func (o *LUD) IsNonsingular() (result bool) {
	for j := 0; j < o.N; j++ {
		if o.LU[j][j] == 0 {
			return false
		}
	}
	return true
}

func (o *QRD) QRDecomposition(a [][]float64) {
	o.QR = a
	o.M = len(a)
	o.N = len(a[0])

	for k := 0; k < o.N; k++ {

		var nrm float64
		for i := k; i < o.M; i++ {
			nrm = math.Hypot(nrm, o.QR[i][k])
		}

		if nrm != 0.0 {
			if o.QR[k][k] < 0 {
				nrm = -nrm
			}
			for i := k; i < o.M; i++ {
				o.QR[i][k] /= nrm
			}
			o.QR[k][k] += 1.0

			for j := k + 1; j < o.N; j++ {
				var s float64
				for i := k; i < o.M; i++ {
					s += o.QR[i][k] * o.QR[i][j]
				}
				s = -s / o.QR[k][k]
				for i := k; i < o.M; i++ {
					o.QR[i][j] += s * o.QR[i][k]
				}
			}
		}
		valk := o.Rdiag[k]
		valk = -nrm
		o.Rdiag[k] = valk
	}
}
func (o *QRD) IsFullRank() bool {
	for j := 0; j < o.N; j++ {
		if o.Rdiag[j] == 0 {
			return false
		}
	}
	return true
}
func (o *QRD) QRDsolve(b [][]float64) (result [][]float64) {
	if len(b) != o.M {
		fmt.Println("(len(b) != o.M):", len(b), o.M)
		//         throw new IllegalArgumentException("Matrix row dimensions must agree.");
	}
	if !o.IsFullRank() {
		fmt.Println("o.IsFullRank():", o.IsFullRank())
		//         throw new RuntimeException("Matrix is rank deficient.");
	}

	// Copy right hand side
	nx := len(b[0])
	//      double[][] X = B.getArrayCopy();

	for k := 0; k < o.N; k++ {
		for j := 0; j < nx; j++ {
			//            double s = 0.0;
			var s float64
			for i := k; i < o.M; i++ {
				s += o.QR[i][k] * b[i][j]
			}
			s = -s / o.QR[k][k]
			for i := k; i < o.M; i++ {
				b[i][j] += s * o.QR[i][k]
			}
		}
	}
	for k := o.N - 1; k >= 0; k-- {
		for j := 0; j < nx; j++ {
			b[k][j] /= o.Rdiag[k]
		}
		for i := 0; i < k; i++ {
			for j := 0; j < nx; j++ {
				b[i][j] -= b[k][j] * o.QR[i][k]
			}
		}
	}
	resultss := make([][]float64, 0)
	for i := 0; i <= (o.N - 1); i++ {
		results := make([]float64, 0)
		for j := 0; j <= (nx - 1); j++ {
			results = append(results, b[i][j])
			//               B[i-i0][j-j0] = A[i][j];
		}
		resultss = append(resultss, results)
	}
	return b //(new Matrix(X,n,nx).getMatrix(0,n-1,0,nx-1));
}

type SendSwitch struct {
	UId     string
	Id      string
	Type    string
	CmdCode string
	Para    string
	UserID  string
	RoomID  string
	Params  string
	FloorID string
}
type ReceiveIbluetooth struct {
	Id          int
	ReceiveStr  string
	ReceiveDate string
}
type ScanDeviceList struct {
	ScandeviceList []ScanDevice
	UserOAuth      core.BasicsToken
	Devicemove     DeviceMoveTB
}

type DeviceMoveTB struct { //app端采集到的详细数据
	Id       int
	Brand    string //Xiaomi//品牌
	Uniqueid string //ffffffff-895c-dd2a-ffff-ffff95cbb2e5//唯一识别Id
	Model    string //Redmi Note 3//手机系列
	Os       string //Android //操作系统
	Product  string //kenzo //产品型号
}
type Bluetooth struct { //蓝牙表Mac地址教室对照表
	Id                int
	Classroomsid      int     //教室Id
	Mac               string  //蓝牙的Mac地址
	BeaconUuid        string  //FDA50693-A4E2-4FB1-AFCF-C6EB07647825,//蓝牙基站逻辑生成的UUid
	Major             int     //微信摇一摇 基站所用
	Minor             int     //微信摇一摇 基站所用
	Name              string  //MTBeaconAC设备的名称
	Battery           int     //100,//电量
	Setlev            string  //A //功率模式
	Txpower           int     //恒定值，自动校正
	Geographical      string  //蓝牙基站所在的具体位置
	LocationX         float64 //基站所在的X轴
	LocationY         float64 //基站所在的Y轴
	LocationZ         float64 //基站所在的Z轴
	FloorsId          int     //楼层Id
	AttenuationFactor float32 //环境衰减因子
}
type ScanDevice struct { //定位数据详细记录表
	Id              int     //主键序列
	Usersid         int     //当前登录人
	Classroomsid    int     //教室Id
	Uniqueid        string  //登录设备的唯一识别码
	Battery         int     //100,//电量
	BeaconUuid      string  //FDA50693-A4E2-4FB1-AFCF-C6EB07647825,//蓝牙基站逻辑生成的UUid
	Beaconid        int     //0,//?
	CurrentDistance float64 //0.056313514709472656,//目前距离[手持的终端离基站的距离]
	CurrentRssi     int     //-48,//信号强度
	Devicetype      string  //5,//设备类型
	Mac             string  //20:C3:8F:E0:61:AC//MAC地址,理想的条件下是唯一的
	Major           int     //10000,//微信摇一摇 基站所用
	Minor           int     //1,////微信摇一摇 基站所用
	Name            string  //MTBeaconAC设备的名称
	Noscancount     int     //0,没有扫描计数
	Setlev          string  //A //功率模式
	Txpower         int     //-64,//功率设定
	SetlevName      string  //可连接
	NowDate         string  //添加时间
	LocationX       float64 //基站所在的X轴
	LocationY       float64 //基站所在的Y轴
	LocationZ       float64 //基站所在的Z轴
	FloorsId        int     //楼层Id
}
type Orbit struct { //用户轨迹表
	Usersid   int     //当前登录人
	TimeStamp int64   //时间戳
	X         float64 //X轴
	Y         float64 //Y轴
	Z         float64 //Z轴
}
type Classroomdetailsnums struct { //教室详情记录表
	Id          int
	Classroomid int
	Usersid     int
	Loginuser   string
	Mac         string
	Createdate  string
	Closedate   string
	Positionnum int
	Xy          string
	X           float64
	Y           float64
	Ap          string
}

type TTScanDevice struct {
	TTUrl string
	SSW   SendSwitch
	TTK   TimedTask
}
type TimedTask struct { //定时任务
	TaskId          int
	TaskState       int    //:定时任务的状态[0:未启动,1:已启动,2:已结束]
	TaskType        int    //:任务类型[0:单次任务,1:循环任务,2:多次任务]
	TaskExecNum     int    //:执行次数[多次任务时:填写执行的次数，执行一次-1，循环任务时：默认-1，单次任务时：1]
	MakeUsersId     int    //:制定任务人Id
	MakeDate        string //:任务制定时间
	ExecBeginDate   string //:任务开始的时间[精确到时分秒]
	ExecEndDate     string //:任务结束的时间[精确到时分秒]
	TaskName        string //:任务的名称
	TaskContent     string //:任务的内容[执行的代码]
	TimeLong        int    //:任务定的时长[示例：定10分钟后执行]
	TimePoint       string //:任务定的时间触发的[示例：每天的下午2点 ,存储特定的时间格式数据]
	RepeatType      string //:重复类型[执行一次、每天、工作日、周末、自定义]
	RepeatValue     string //:重复类型的值
	EventSetTableId int    //事件定义Id
}
type EventSetTable struct { //事件设置表
	EventSetTableId int    //事件定义Id
	EventName       string //事件名称
	EventContent    string //事件执行的内容
}

func InitIbeaconDb() *gorp.DbMap {
	db, _ := sql.Open("mysql", GlobalConfig.Read("Data", "datastr")) //GlobalConfig.Read("Data", "datastr")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(ReceiveIbluetooth{}, "receiveibluetooth").SetKeys(true, "Id")
	dbmap.AddTableWithName(ScanDevice{}, "scandevice").SetKeys(true, "Id")
	dbmap.AddTableWithName(DeviceMoveTB{}, "devicemovetb").SetKeys(true, "Id")
	dbmap.AddTableWithName(Orbit{}, "orbit")
	dbmap.CreateTablesIfNotExists()
	return dbmap
}
func InitZNDXDb() *gorp.DbMap {
	db, _ := sql.Open("mysql", GlobalConfig.Read("Data", "datastr2")) //GlobalConfig.Read("Data", "datastr")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dbmap.AddTableWithName(Bluetooth{}, "bluetooth").SetKeys(true, "Id")
	dbmap.AddTableWithName(Classroomdetailsnums{}, "classroomdetailsnums").SetKeys(true, "Id")
	dbmap.AddTableWithName(TimedTask{}, "timedtask").SetKeys(true, "TaskId")
	dbmap.CreateTablesIfNotExists()
	return dbmap
}
