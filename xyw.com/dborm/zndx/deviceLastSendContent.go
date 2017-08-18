package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `devicelastsendcontent` (
  `DeviceId` varchar(50) NOT NULL,
  `CmdCode` varchar(50) NOT NULL,
  `LastSendContent` varchar(500) DEFAULT NULL,
  `SendTime` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`DeviceId`,`CmdCode`)
) ENGINE=MyISAM DEFAULT CHARSET=gbk;
*/
type DeviceLastSendContent struct {
	DeviceId        string // 设备ID
	CmdCode         string // 命令代码
	LastSendContent string // 最后发送内容
	SendTime        string // 发送时间
}

// 获取DeviceLastSendContent
func DeviceLastSendContent_Get(id, cmd string, pDBMap *gorp.DbMap) (pContent *DeviceLastSendContent, err error) {
	pDBMap.AddTableWithName(DeviceLastSendContent{}, "devicelastsendcontent").SetKeys(false, "DeviceId", "CmdCode")
	pObj, err := pDBMap.Get(DeviceLastSendContent{}, id, cmd)

	pContent, _ = pObj.(*DeviceLastSendContent)
	return pContent, err
}

// 插入
func (p *DeviceLastSendContent) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceLastSendContent{}, "devicelastsendcontent").SetKeys(false, "DeviceId", "CmdCode")
	return pDBMap.Insert(p)
}

// 删除
func (p *DeviceLastSendContent) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceLastSendContent{}, "devicelastsendcontent").SetKeys(false, "DeviceId", "CmdCode")
	rows, err := pDBMap.Delete(p)
	return int(rows), err
}

// 更新
func (p *DeviceLastSendContent) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceLastSendContent{}, "devicelastsendcontent").SetKeys(false, "DeviceId", "CmdCode")
	rows, err := pDBMap.Update(p)
	return int(rows), err
}

// 保存
func (p *DeviceLastSendContent) Save(dbmap *gorp.DbMap) bool {
	//先删除
	sql := " delete from DeviceLastSendContent where DeviceId=? and CmdCode=? "
	_, err1 := dbmap.Exec(sql, p.DeviceId, p.CmdCode)
	xdebug.LogError(err1)

	//再保存
	sql = " insert into DeviceLastSendContent (DeviceId,CmdCode,LastSendContent,SendTime) values (?,?,?,?) "
	_, err2 := dbmap.Exec(sql, p.DeviceId, p.CmdCode, p.LastSendContent, p.SendTime)
	xdebug.LogError(err2)

	return (err1 == nil) && (err2 == nil)
}

// 查询设备ID对应的DeviceLastSendContent存在?
func DeviceLastSendContent_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM devicelastsendcontent WHERE ObjectId = ?`
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 查询设备ID对应的DeviceLastSendContent
func DeviceLastSendContent_Query_ByDevice(device_id string, dbmap *gorp.DbMap) (list []DeviceLastSendContent, err error) {
	sql := "select * from DeviceLastSendContent where DeviceId=? "
	_, err = dbmap.Select(&list, sql, device_id)

	return list, err
}
