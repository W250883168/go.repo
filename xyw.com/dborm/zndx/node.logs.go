package zndx

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"gopkg.in/gorp.v1"
)

/*
CREATE TABLE `node_logs` (
  `log_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `log_time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `log_topic` varchar(255) NOT NULL DEFAULT '',
  `log_content` varchar(1024) NOT NULL DEFAULT '',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;*/
type NodeLog struct {
	LogID      int       `db:"log_id"`
	LogTime    time.Time `db:"log_time"`
	LogTopic   string    `db:"log_topic"`
	LogContent string    `db:"log_content"`
}

// 查询
func NodeLog_Get(id int, pDBMap *gorp.DbMap) (pLog *NodeLog, err error) {
	pDBMap.AddTableWithName(NodeLog{}, "node_logs").SetKeys(true, "log_id")
	pObj, err := pDBMap.Get(NodeLog{}, id)

	pLog, _ = pObj.(*NodeLog)
	return pLog, err
}

// 插入
func (p *NodeLog) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(NodeLog{}, "node_logs").SetKeys(true, "log_id")
	return pDBMap.Insert(p)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
