package app

import (
	"time"
)

const (
	_INFO = 1 << iota
	_DEBUG
	_WARN
	_ERROR
	_FATAL
)

// 调试级别
type EDebugLevel int

func (level EDebugLevel) String() (str string) {
	switch int(level) {
	case _INFO:
		str = "info"
	case _DEBUG:
		str = "debug"
	case _WARN:
		str = "warn"
	case _ERROR:
		str = "error"
	case _FATAL:
		str = "fatal"
	}

	return str
}

// CoAP配置
type CoapConfig struct {
	OffTimeout        int // 离线超时(秒)
	AckTimeout        int // 确认超时(秒)
	MaxSendCount      int // 最大发送次数
	HeartbeatInterval int // 心跳间隔(秒)
}

// 离线超时字符串
func (p *CoapConfig) OffTimeout_String() string {
	dur := time.Duration(p.AckTimeout) * time.Second
	return dur.String()
}

// App配置结构
type AppConfig struct {
	ConfigFile string // 配置文件
	DataSource string // 数据源连接字符串

	GinHTTPPort   int // gin.http端口
	PprofHTTPPort int // pprof.http端口

	LogFile    string // 日志文件
	LogShowSQL bool   // 日志显示SQL
	DebugLevel EDebugLevel

	BroadcastingSupport bool   // 广播功能支持
	ThisHostAddr        string // 本机地址(广播功能)

	PJLinkTimeout int
	Version       string
	PublistState  string
	CoapConfig
}
