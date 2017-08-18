package dborm

import (
	"fmt"
	"log"
	"runtime"

	"gopkg.in/gorp.v1"
	"vodx/ioutil/dbutil"
)

/*
CREATE TABLE `computer_info` (
  `computer_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `location_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '位置ID',
  `computer_name` varchar(255) NOT NULL DEFAULT '' COMMENT '电脑名称',
  `computer_ip` varchar(255) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `computer_port` int(11) NOT NULL DEFAULT '0' COMMENT '端口号',
  `login_account` varchar(255) NOT NULL DEFAULT '' COMMENT '登录账号',
  `login_password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码 ',
  PRIMARY KEY (`computer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='教学电脑信息表';

*/

type ComputerInfo struct {
	ComputerID    int    `db:"computer_id"`
	LocationID    int    `db:"location_id"`
	ComputerName  string `db:"computer_name"`
	ComputerIP    string `db:"computer_ip"`
	ComputerPort  int    `db:"computer_port"`
	LoginAccount  string `db:"login_account"`
	LoginPassword string `db:"login_password"`
}

func ComputerInfo_Query(id int, pDBMap *gorp.DbMap) (pInfo *ComputerInfo, err error) {
	pDBMap.AddTableWithName(ComputerInfo{}, "computer_info").SetKeys(true, "computer_id")
	pObj, err := pDBMap.Get(ComputerInfo{}, id)
	pInfo, _ = pObj.(*ComputerInfo)

	return pInfo, err
}

func ComputerInfo_QueryByLocation(locationID int, pDBMap *gorp.DbMap) (list []ComputerInfo, err error) {
	query := `
SELECT
		TComputer.computer_id,
		TComputer.location_id,
		TComputer.computer_name,
		TComputer.computer_ip,
		TComputer.computer_port,
		TComputer.login_account,
		TComputer.login_password
FROM computer_info AS TComputer
WHERE TComputer.location_id = ?`
	list = []ComputerInfo{}
	_, err = pDBMap.Select(&list, query, locationID)
	return list, err
}

func (p *ComputerInfo) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(ComputerInfo{}, "computer_info").SetKeys(true, "computer_id")
	return pDBMap.Insert(p)
}

func (p *ComputerInfo) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(ComputerInfo{}, "computer_info").SetKeys(true, "computer_id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *ComputerInfo) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(ComputerInfo{}, "computer_info").SetKeys(true, "computer_id")
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
	pDBMap.AddTableWithName(ComputerInfo{}, "computer_info").SetKeys(true, "computer_id")
}
