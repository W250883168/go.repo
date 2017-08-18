package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `pjlink` (
  `DeviceId` varchar(50) NOT NULL,
  `Address` varchar(50) DEFAULT NULL,
  `Port` varchar(50) DEFAULT NULL,
  `Class` varchar(50) DEFAULT NULL,
  `Password` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`DeviceId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
type PJLink struct {
	DeviceId string
	Address  string
	Port     string
	Class    string
	Password string
}

// 查询存在PJLink?
func PJLink_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM pjlink WHERE DeviceId = ?"
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}
