package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `devicemodelfaultword` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `Name` varchar(50) DEFAULT NULL,
  `ModelId` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8 COMMENT='为了方便录入故障现象而定义常用词条';
*/
type DeviceModelFaultWord struct {
	Id      int
	Name    string
	ModelId string
}

// 获取DeviceModelFaultWord
func DeviceModelFaultWord_Get(id int, pDBMap *gorp.DbMap) (pWord *DeviceModelFaultWord, err error) {
	pDBMap.AddTableWithName(DeviceModelFaultWord{}, "devicemodelfaultword").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceModelFaultWord{}, id)

	pWord, _ = pObj.(*DeviceModelFaultWord)
	return pWord, err
}

// 查询对应模型ID的DeviceModelFaultWord存在?
func DeviceModelFaultWord_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicemodelfaultword WHERE ModelId = ?"
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 插入
func (p *DeviceModelFaultWord) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceModelFaultWord{}, "devicemodelfaultword").SetKeys(true, "Id")
	return pDBMap.Insert(p)
}

// 删除
func (p *DeviceModelFaultWord) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModelFaultWord{}, "devicemodelfaultword").SetKeys(true, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 更新
func (p *DeviceModelFaultWord) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModelFaultWord{}, "devicemodelfaultword").SetKeys(true, "Id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}
