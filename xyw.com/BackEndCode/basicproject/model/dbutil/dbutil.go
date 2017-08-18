package dbutil

import (
	"fmt"
	"log"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"gopkg.in/ini.v1"

	"xutils/xerr"

	core "xutils/xcore"
)

var gDataSource string = "root:root@tcp(localhost:3306)/zndx"
var gShowSQL bool = false

func ShowSQL() bool {
	return gShowSQL
}

func InitDB() *xorm.Engine {
	pEngine, err := xorm.NewEngine("mysql", gDataSource)
	xerr.ThrowPanic(err)

	pEngine.ShowSQL(gShowSQL)
	return pEngine
}

func CheckDB() (ok bool) {
	pEngine := InitDB()
	defer pEngine.Close()

	_, err := pEngine.Exec("SELECT 1")
	xerr.ThrowPanic(err)

	ok = true
	return ok
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

	file, err := ini.Load(core.ConfigFile)
	xerr.ThrowPanic(err)
	gDataSource = file.Section("data").Key("datastr").Value()
	gShowSQL, _ = file.Section("data").Key("showsql").Bool()

	datasrc := fmt.Sprintf("	XORM/DATASOURCE: %s; SHOWSQL: %t", gDataSource, gShowSQL)
	log.Println(datasrc)
}
