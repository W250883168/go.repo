package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `eventsettable` (
  `EventSetTableId` int(11) NOT NULL AUTO_INCREMENT,
  `EventName` varchar(200) DEFAULT NULL COMMENT '事件名称',
  `EventContent` varchar(255) DEFAULT NULL COMMENT '事件执行的内容',
  `CampusId` int(11) DEFAULT '0',
  `BuildingId` int(11) DEFAULT '0',
  `FloorsId` int(11) DEFAULT '0',
  `ClassRoomId` int(11) DEFAULT '0',
  `NodeId` varchar(50) DEFAULT NULL COMMENT '节点Id',
  `DeviceId` varchar(50) DEFAULT NULL COMMENT '设备Id',
  `CmdId` varchar(50) DEFAULT NULL COMMENT '命令Id',
  PRIMARY KEY (`EventSetTableId`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='事件设置表';
*/
type EventSetTable struct {
	EventSetTableId int
	EventName       string
	EventContent    string
	CampusId        int
	BuildingId      int
	FloorsId        int
	ClassRoomId     int
	NodeId          string
	DeviceId        string
	CmdId           string
}

func EventSetTable_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM eventsettable WHERE DeviceId = ?"
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

func EventSetTable_Exists_ByNodeID(nodeID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM eventsettable WHERE NodeId = ?`
	count, err := pDBMap.SelectInt(query, nodeID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}
