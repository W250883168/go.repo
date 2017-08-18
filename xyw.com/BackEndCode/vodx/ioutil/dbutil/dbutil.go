package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"

	"gopkg.in/gorp.v1"

	_ "github.com/go-sql-driver/mysql"

	"xutils/xerr"

	"vodx/app"
)

var (
	gConnString      = ""
	gShowSQL    bool = false

	pDBMap *gorp.DbMap = nil
)

func GetDBMap() *gorp.DbMap {
	if pDBMap == nil {
		db, err := sql.Open("mysql", app.GetConfig().DataSource)
		xerr.ThrowPanic(err)

		pDBMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
		if app.GetConfig().ShowSQL {
			pDBMap.TraceOn("", log.New(os.Stdout, "[gorp]", log.LstdFlags|log.Lshortfile))
		}
	}

	return pDBMap
}

func CheckDB() {
	dbmap := GetDBMap()

	_, err := dbmap.Exec("Select 1")
	xerr.ThrowPanic(err)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
