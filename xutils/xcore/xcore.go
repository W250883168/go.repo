package xcore

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"

	"xutils/xconfig"
	"xutils/xerr"
	"xutils/xfile"
	"xutils/xtime"
)

var cfi xconfig.Config
var pLogFile *os.File

func GetLogFile() *os.File {
	return pLogFile
}
func InitDb() *gorp.DbMap {
	db, err := sql.Open("mysql", cfi.Read("data", "datastr"))
	CheckErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"MyISAM", "UTF8"}}
	CheckErr(err, "Create tables failed")

	//	if dbutil.ShowSQL() {
	//		dbmap.TraceOn("[gorp]", log.New(os.Stdout, "", log.LstdFlags))
	//	}

	return dbmap
}
func InitPointMysqlDb(configrow, configname string) *gorp.DbMap {
	db, err := sql.Open("mysql", cfi.Read(configrow, configname))
	CheckErr(err, "sql.Open failed")
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"MyISAM", "UTF8"}}
	CheckErr(err, "Create tables failed")

	return dbmap
}

// 检查数据库连接
func CheckDB() {
	dbmap := InitDb()
	defer dbmap.Db.Close()

	sql := "SELECT 1 "
	_, err := dbmap.Exec(sql)
	xerr.ThrowPanic(err)
}

func CheckErr(err error, msg string) bool {
	if err != nil {
		log.Println(msg, err)
		go WriteLog(msg + err.Error() + "\n")
		return false
	}
	return true
}

type BasicsToken struct {
	Usersid   int
	Rolestype int
	Token     string
	Os        string
}

type Returndata struct {
	Rcode  string      //返回状态
	Reason string      //返回消息内容
	Result interface{} //返回数据
}

type PageData struct {
	PageIndex int         //当前页
	PageCount int         //总数量
	PageSize  int         //每页大小
	PageData  interface{} //内容数据
}

//文件字符串串截取
func Substrs(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

//获取上级文件目录
func GetParentDirectory(dirctory string) string {
	return Substrs(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

//获取当前文件运行目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//判断文件或目录是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func GetLimitString(pg PageData) (limit string) {
	if pg.PageIndex > 1 {
		limit = " limit " + strconv.Itoa((pg.PageIndex-1)*pg.PageSize) + "," + strconv.Itoa(pg.PageSize) + ";"
	} else if pg.PageIndex == -1 {
		limit = ";"
	} else {
		if pg.PageSize == 0 {
			pg.PageSize = 10
		}
		limit = " limit 0," + strconv.Itoa(pg.PageSize) + ";"
	}
	return limit
}

func Timeaction(strtm string) string {
	tkbtarr1 := strings.Split(strtm, " ")
	tkbtarr2 := strings.Split(tkbtarr1[0], "/") //判断是否有/符号
	if len(tkbtarr2) == 1 {
		tkbtarr2 = strings.Split(tkbtarr1[0], "-")
	}
	if len(tkbtarr2) == 3 {
		if len(tkbtarr2[1]) < 2 {
			tkbtarr2[1] = "0" + tkbtarr2[1]
		}
		if len(tkbtarr2[2]) < 2 {
			tkbtarr2[2] = "0" + tkbtarr2[2]
		}
		strtm = tkbtarr2[0] + "-" + tkbtarr2[1] + "-" + tkbtarr2[2]
	} else {
		strtm = ""
	}
	if len(tkbtarr1) >= 2 { //判断是否精确到小时
		tkbtarr3 := strings.Split(tkbtarr1[1], ":") //判断是否有/符号
		if len(tkbtarr3) == 3 {
			if len(tkbtarr3[0]) < 2 {
				tkbtarr3[0] = "0" + tkbtarr3[0]
			}
			if len(tkbtarr3[1]) < 2 {
				tkbtarr3[1] = "0" + tkbtarr3[1]
			}
			if len(tkbtarr3[2]) < 2 {
				tkbtarr3[2] = "0" + tkbtarr3[2]
			}
			if strtm != "" {
				strtm = strtm + " " + tkbtarr3[0] + ":" + tkbtarr3[1] + ":" + tkbtarr3[2]
			}
		}
	}
	return strtm
}

//读取文本文件内容
func Readfile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

func WriteLog(errstr string) {
	timestr := time.Now().Format("20060102") + ".log"
	var filename = "./temporarylogfile/" + timestr
	var f *os.File
	defer f.Close()
	if CheckFileIsExist(filename) { //如果文件存在
		f, _ = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
	} else {
		f, _ = os.Create(filename) //创建文件
	}
	f.WriteString(errstr + "  " + time.Now().Format("2006-01-02 15:04:05") + " \n")
	f.Sync()
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func WriteFFmpegFile(filepath string, context string) {
	var f *os.File
	defer f.Close()
	var err1 error
	if CheckFileIsExist(filepath) { //如果文件存在
		f, err1 = os.OpenFile(filepath, os.O_APPEND, 0666) //打开文件
	} else {
		f, err1 = os.Create(filepath) //创建文件
	}
	CheckErr(err1, "关闭录制时错误")
	f.WriteString(context)
	f.Sync()
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

func init() {
	cfi.InitConfig("./config.ini")
	dir := "./log/"
	if !xfile.FileExist(dir) {
		os.Mkdir(dir, os.ModeDir)
	}

	if logfile, ok := xfile.CreateFile(dir + time.Now().Format(xtime.FORMAT_yyyyMMddHHmmss) + ".log"); ok {
		pLogFile = logfile
		os.Stdout = logfile
	}
}
