package app

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"xutils/xerr"

	"gopkg.in/ini.v1"
)

const ConfigFile = "vod.config.ini"

var pConfig = &_AppConfig{}

func GetConfig() _AppConfig {
	return *pConfig
}

func init() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d\n", fun.Name(), line)
		log.Printf(str)
	}

	file, err := ini.Load(ConfigFile)
	xerr.ThrowPanic(err)
	name := file.Section("vodx").Key("app.name").String()
	publish_state := file.Section("vodx").Key("app.publish.state").String()
	version := file.Section("vodx").Key("app.version").String()
	datasrc := file.Section("vodx").Key("database.mysql.datasrc").String()
	showsql, _ := file.Section("vodx").Key("database.mysql.showsql").Bool()
	http_port, _ := file.Section("vodx").Key("app.http.port").Int()
	debug_port, _ := file.Section("vodx").Key("app.debug.http.port").Int()
	mqConnString := file.Section("vodx").Key("rabbitmq.connection.string").String()

	pConfig.AppName = name
	pConfig.DataSource = datasrc
	pConfig.ShowSQL = showsql
	pConfig.PublishState = publish_state
	pConfig.Version = version
	pConfig.HttpPort = http_port
	pConfig.ProfHttpPort = debug_port
	pConfig.MQConnString = mqConnString

	buff, _ := json.Marshal(pConfig)
	fmt.Printf("\t    AppConfig: %+v\n", string(buff))
}
