package zndx

import (
	"strconv"

	"xutils/xdebug"

	"gopkg.in/gorp.v1"
)

/*
CREATE TABLE `deviceoperatelog` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `OperateTime` varchar(50) DEFAULT NULL,
  `OperateUserId` varchar(50) DEFAULT NULL,
  `OperateType` varchar(50) DEFAULT NULL,
  `OperateObject` varchar(50) DEFAULT NULL,
  `ObjectId` varchar(50) DEFAULT NULL,
  `ObjectName` varchar(200) DEFAULT NULL,
  `UseWhoseCmd` varchar(50) DEFAULT NULL,
  `CmdCode` varchar(50) DEFAULT NULL,
  `CmdName` varchar(50) DEFAULT NULL,
  `Para` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=3149 DEFAULT CHARSET=gbk COMMENT='存储向节点发送的每一个命令';
记录对教室、楼栋、设备的每一次控制操作
*/
type DeviceOperateLog struct {
	Id            int
	OperateTime   string // 操作时间
	OperateUserId string // 操作人ID
	OperateType   string // 操作类型
	OperateObject string // 操作对象
	ObjectId      string // 对象ID
	ObjectName    string // 对象名称
	UseWhoseCmd   string // 使用谁的命令
	CmdCode       string // 命令代码
	CmdName       string // 命令代码名称
	Para          string // 参数
}

// 获取DeviceOperateLog
func DeviceOperateLog_Get(id int, pDBMap *gorp.DbMap) (pLog *DeviceOperateLog, err error) {
	pDBMap.AddTableWithName(DeviceOperateLog{}, "deviceoperatelog").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceOperateLog{}, id)

	pLog, _ = pObj.(*DeviceOperateLog)
	return pLog, err
}

// 删除
func (p *DeviceOperateLog) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceOperateLog{}, "deviceoperatelog").SetKeys(true, "Id")
	rows, err := pDBMap.Delete(p)
	return int(rows), err
}

// 更新
func (p *DeviceOperateLog) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceOperateLog{}, "deviceoperatelog").SetKeys(true, "Id")
	rows, err := pDBMap.Update(p)
	return int(rows), err
}

// 查询是否存在DeviceOperateLog(通过DeviceID)
func DeviceOperateLog_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM deviceoperatelog WHERE ObjectId = ?"
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 插入
func (p *DeviceOperateLog) Insert(dbmap *gorp.DbMap) (err error) {
	switch p.OperateObject {
	case "device":
		device := Device{Id: p.ObjectId}
		p.ObjectName = device.GetDeviceDetailName(dbmap)
	case "classroom":
		id, _ := strconv.Atoi(p.ObjectId)
		room := Classroom{Id: id}
		p.ObjectName = room.GetClassroomDetailName(dbmap)
	case "floor":
		id, _ := strconv.Atoi(p.ObjectId)
		floor := Floor{Id: id}
		p.ObjectName = floor.GetFloorDetailName(dbmap)
	}

	dbmap.AddTableWithName(DeviceOperateLog{}, "deviceoperatelog").SetKeys(true, "Id")
	err = dbmap.Insert(p)
	return err
}
