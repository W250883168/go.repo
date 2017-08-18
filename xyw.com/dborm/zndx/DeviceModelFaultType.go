package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `devicemodelfaulttype` (
  `Id` varchar(50) NOT NULL,
  `Name` varchar(50) DEFAULT NULL,
  `ModelId` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
type DeviceModelFaultType struct {
	Id      string
	Name    string
	ModelId string
}

// 获取DeviceModelFaultType
func DeviceModelFaultType_Get(id string, pDBMap *gorp.DbMap) (pType *DeviceModelFaultType, err error) {
	pDBMap.AddTableWithName(DeviceModelFaultType{}, "devicemodelfaulttype").SetKeys(false, "Id")
	pObj, err := pDBMap.Get(DeviceModelFaultType{}, id)

	pType, _ = pObj.(*DeviceModelFaultType)
	return pType, err
}

// 查询对应模型ID的DeviceModelFaultType存在?
func DeviceModelFaultType_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicemodelfaulttype WHERE ModelId = ?"
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 插入
func (p *DeviceModelFaultType) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceModelFaultType{}, "devicemodelfaulttype").SetKeys(false, "Id")
	return pDBMap.Insert(p)
}

// 删除
func (p *DeviceModelFaultType) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModelFaultType{}, "devicemodelfaulttype").SetKeys(false, "Id")
	rows, err := pDBMap.Delete(p)
	return int(rows), err
}

// 更新
func (p *DeviceModelFaultType) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModelFaultType{}, "devicemodelfaulttype").SetKeys(false, "Id")
	rows, err := pDBMap.Update(p)
	return int(rows), err
}

// 查询是否被设备故障引用?
func (p *DeviceModelFaultType) Referenced_ByDeviceFault(pDBMap *gorp.DbMap) (yes bool) {
	yes = true
	query := `SELECT COUNT(1) FROM devicemodelfaulttype TModelFaultType WHERE  ? IN (SELECT DISTINCT FaultTypeId FROM devicefaulttype)`
	count, err := pDBMap.SelectInt(query, p.Id)
	xdebug.LogError(err)
	if err == nil {
		yes = count > 0
	}

	return yes
}
