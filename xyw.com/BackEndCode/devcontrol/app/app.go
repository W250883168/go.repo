package app

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"gopkg.in/ini.v1"

	"xutils/xerr"
	"xutils/xfile"
	"xutils/xtime"
)

const (
	_CONFIG_FILE = "config.ini"

	// 默认值
	_GIN_HTTP_PORT   = 8090
	_PPROF_HTTP_PORT = 8091
	_DATASOURCE      = `root:root@tcp(localhost:3306)/zndx2`
	_DEBUG_LEVEL     = _WARN
	_LOG_SHOWSQL     = false

	_COAP_ACK_TIMEOUT        = 2
	_COAP_MAX_SEND_COUNT     = 5
	_COAP_OFFLINE_TIMEOUT    = 30
	_COAP_HEARTBEAT_INTERVAL = 10
	_PJLINK_ACK_TIMEOUT      = 2
)

var (
	pConfig  *AppConfig
	pLogFile *os.File
)

func GetConfig() AppConfig {
	return *pConfig
}

func GetLogFile() *os.File {
	return pLogFile
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	// 默认配置
	pLogFile = os.Stdout
	pConfig = &AppConfig{
		ConfigFile:    _CONFIG_FILE,
		DataSource:    _DATASOURCE,
		GinHTTPPort:   _GIN_HTTP_PORT,
		PprofHTTPPort: _PPROF_HTTP_PORT,
		LogShowSQL:    _LOG_SHOWSQL,
		DebugLevel:    _DEBUG_LEVEL,
		PJLinkTimeout: _PJLINK_ACK_TIMEOUT,
		CoapConfig: CoapConfig{
			OffTimeout:        _COAP_OFFLINE_TIMEOUT,
			AckTimeout:        _COAP_ACK_TIMEOUT,
			MaxSendCount:      _COAP_MAX_SEND_COUNT,
			HeartbeatInterval: _COAP_HEARTBEAT_INTERVAL},
	}

	// 日志文件
	dir := "./log/"
	if !xfile.FileExist(dir) {
		os.MkdirAll(dir, os.ModeDir)
	}
	filepath := dir + time.Now().Format(xtime.FORMAT_yyyyMMddHHmmss) + ".log"
	if logfile, ok := xfile.CreateFile(filepath); ok {
		pLogFile = logfile
		pConfig.LogFile = logfile.Name()
	}

	// 加载配置
	file, err := ini.Load(_CONFIG_FILE)
	xerr.ThrowPanic(err)
	section := file.Section("app.config")
	pConfig.DataSource = section.Key(`database.mysql.datasrc`).String()
	pConfig.ThisHostAddr = section.Key(`app.this.host.addr`).String()
	pConfig.PublistState = section.Key(`app.publish.state`).String()
	pConfig.Version = section.Key(`app.version`).String()
	if port, err := section.Key(`gin.http.port`).Int(); err == nil {
		pConfig.GinHTTPPort = port
	}
	if port, err := section.Key(`pprof.http.port`).Int(); err == nil {
		pConfig.PprofHTTPPort = port
	}
	if timeout, err := section.Key(`coap.offline.timeout.seconds`).Int(); err == nil {
		pConfig.CoapConfig.OffTimeout = timeout
	}
	if timeout, err := section.Key(`coap.ack.timeout.seconds`).Int(); err == nil {
		pConfig.CoapConfig.AckTimeout = timeout
	}
	if max, err := section.Key(`coap.max.send.count`).Int(); err == nil {
		pConfig.CoapConfig.MaxSendCount = max
	}
	if showsql, err := section.Key(`app.log.showsql`).Bool(); err == nil {
		pConfig.LogShowSQL = showsql
	}
	if level, err := section.Key(`app.debug.level`).Int(); err == nil {
		pConfig.DebugLevel = EDebugLevel(level)
	}
	if support, err := section.Key(`app.device.broadcasting.support`).Bool(); err == nil {
		pConfig.BroadcastingSupport = support
	}
	if interval, err := section.Key(`coap.node.heartbeat.interval.seconds`).Int(); err == nil {
		pConfig.CoapConfig.HeartbeatInterval = interval
	}

}
