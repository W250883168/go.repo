package coap

import (
	"fmt"
	"log"
	"runtime"
	"sync"
)

const (
	OfflineTime      = "30s" //多长时间节点没有上传数据就算该节点为离线
	AckTimeout       = 2     //应答超时时间，超过此时间将重新发送命令(如果MaxSendCound>0的话)
	MaxSendCound     = 5     //最大发送次数(即使设置为0，系统也会至少发送一次）
	AckTimeoutPJLink = 2     //pjlink超时时间
)

var pOTAContext *OTAContext

// OTA进度Context, 封装了并发安全的map
type OTAContext struct {
	KValue map[string]string
	lock   sync.Mutex
}

func GetOTAContext() *OTAContext {
	return pOTAContext
}

func (p *OTAContext) GetValue(key string) string {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.KValue[key]
}

func (p *OTAContext) PutValue(key, value string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.KValue[key] = value
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	pOTAContext = &OTAContext{KValue: map[string]string{}}
}
