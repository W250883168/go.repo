package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"

	"xutils/xerr"

	"dev.project/BackEndCode/devcontrol/app"
)

var (
	pDBMap *gorp.DbMap // 全局数据库连接池
)

// 获取数据库连接
func GetDBMap() *gorp.DbMap {
	if pDBMap == nil {
		db, err := sql.Open("mysql", app.GetConfig().DataSource)
		xerr.ThrowPanic(err)

		pDBMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
		if app.GetConfig().LogShowSQL {
			if pLogFile := app.GetLogFile(); pLogFile != nil {
				pDBMap.TraceOn("", log.New(pLogFile, "[gorp]", log.LstdFlags|log.Lshortfile))
			}
		}
	}

	return pDBMap
}

// 检查数据库连接
func CheckDB() {
	pDBMap := GetDBMap()

	_, err := pDBMap.Exec("Select 1")
	xerr.ThrowPanic(err)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

}
