package dborm

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"gopkg.in/gorp.v1"
	"vodx/ioutil/dbutil"
)

/*
CREATE TABLE `operate_log` (
  `log_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `log_type` int(11) NOT NULL DEFAULT '0' COMMENT '0/query; 1/add; 2/delete; 3/update;',
  `log_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP COMMENT '日志日期',
  `log_content` varchar(1024) NOT NULL DEFAULT '' COMMENT '日志内容',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

type OperateLog struct {
	LogID      int       `db:"log_id"`
	LogType    int       `db:"log_type"`
	LogDate    time.Time `db:"log_date"`
	LogContent string    `db:"log_content"`
}

func OperateLog_Query(id int, pDBMap *gorp.DbMap) (pLog *OperateLog, err error) {
	pDBMap.AddTableWithName(OperateLog{}, "operate_log").SetKeys(true, "log_id")
	pObj, err := pDBMap.Get(OperateLog{}, id)

	pLog, _ = pObj.(*OperateLog)
	return pLog, err
}

func (p *OperateLog) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(OperateLog{}, "operate_log").SetKeys(true, "log_id")
	return pDBMap.Insert(p)
}

func (p *OperateLog) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(OperateLog{}, "operate_log").SetKeys(true, "log_id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *OperateLog) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(OperateLog{}, "operate_log").SetKeys(true, "log_id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	pDBMap := dbutil.GetDBMap()
	pDBMap.AddTableWithName(OperateLog{}, "operate_log").SetKeys(true, "log_id")
}
