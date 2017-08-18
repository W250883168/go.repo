package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xhttp"
)

/*
CREATE TABLE `deviceuselog` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `DeviceId` varchar(50) DEFAULT NULL,
  `OnTime` varchar(50) DEFAULT NULL,
  `OffTime` varchar(50) DEFAULT NULL,
  `UseTime` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=1217 DEFAULT CHARSET=gbk COMMENT='1.通过每一次的开机关机命令获取设备使用的开始时间和结束时间；\r\n2.通过解析节点每10秒上报一次数据获取设备';
*/
type DeviceUseLog struct {
	Id       int64
	DeviceId string // 设备ID
	OnTime   string // 开始时间
	OffTime  string // 结束时间
	UseTime  int64  // 使用时间(秒)
}

// 获取DeviceUseLog
func DeviceUseLog_Get(id int, pDBMap *gorp.DbMap) (pLog *DeviceUseLog, err error) {
	pDBMap.AddTableWithName(DeviceUseLog{}, "deviceuselog").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceUseLog{}, id)

	pLog, _ = pObj.(*DeviceUseLog)
	return pLog, err
}

// 查询是否存在日志，该日志设备ID与指定的参数相同
func DeviceUseLog_Exists_ByDeviceID(deviceID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM deviceuselog WHERE DeviceId = ?`
	count, err := pDBMap.SelectInt(query, deviceID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 查询存在否
func (p *DeviceUseLog) Exists(pDBMap *gorp.DbMap) (existed bool) {
	sql := "SELECT COUNT(*) FROM deviceuselog WHERE Id = ? "
	nullInt, err := pDBMap.SelectNullInt(sql, p.Id)
	existed = (err == nil) && (nullInt.Valid)
	return existed
}

//插入
func (p *DeviceUseLog) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceUseLog{}, "deviceuselog").SetKeys(true, "Id")
	err = pDBMap.Insert(p)

	return err
}

// 删除
func (p *DeviceUseLog) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceUseLog{}, "deviceuselog").SetKeys(true, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 更新
func (p *DeviceUseLog) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceUseLog{}, "deviceuselog").SetKeys(true, "Id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

// 更新日志设备ID为默认值
func DeviceUseLog_UpdateDefault_DeviceID(deviceID string, pDBMap *gorp.DbMap) (affects int, err error) {
	statement := `UPDATE deviceuselog SET DeviceId = '' WHERE DeviceId = ?`
	ret, err := pDBMap.Exec(statement, deviceID)

	rows, _ := ret.RowsAffected()
	return int(rows), err
}

// 获取设备日志（通过DeviceID）
func (p *DeviceUseLog) GetByDeviceID(page *xhttp.PageInfo, dbmap *gorp.DbMap) (list []DeviceUseLogView, err error) {
	defer xerr.CatchPanic()
	sqlcount := `SELECT COUNT(*) FROM DeviceUseLog WHERE DeviceId = ?`
	count, err := dbmap.SelectInt(sqlcount, p.DeviceId)
	xerr.ThrowPanic(err)
	page.RowTotal = int(count)

	sql := `
SELECT OnTime, IFNULL(OffTime, '') OffTime, IFNULL(SEC_TO_TIME(UseTime), '') UseTime
FROM 	DeviceUseLog
WHERE DeviceId = ?
ORDER BY OnTime DESC
`
	sql += page.SQL_LimitString()
	_, err = dbmap.Select(&list, sql, p.DeviceId)
	return list, err
}

// 查询设备是否有开机记录
func DeviceUseLog_Exists_DeviceOpenLog(device_id string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	sql := `SELECT COUNT(1) FROM DeviceUseLog WHERE DeviceId =? AND IFNULL(OffTime, '') = ''`
	count, err := pDBMap.SelectInt(sql, device_id)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 查询设备开机记录
func DeviceUseLog_Query_DeviceOpenLog(device_id string, pDBMap *gorp.DbMap) (log DeviceUseLog, err error) {
	sql := `SELECT * FROM DeviceUseLog WHERE DeviceId =? AND IFNULL(OffTime, '') = ''`
	err = pDBMap.SelectOne(&log, sql, device_id)
	return log, err
}

// 插入设备开机记录
func DeviceUseLog_Insert_DeviceOpenLog(device_id, onTime string, pDBMap *gorp.DbMap) (err error) {
	log := DeviceUseLog{DeviceId: device_id, OnTime: onTime}
	err = log.Insert(pDBMap)

	return err
}

// 设备使用时间日志
type DeviceUseLogView struct {
	Id       int64
	DeviceId string // 设备ID
	OnTime   string // 开始时间
	OffTime  string // 结束时间
	UseTime  string // 使用时间('HH:MM:SS')
}
