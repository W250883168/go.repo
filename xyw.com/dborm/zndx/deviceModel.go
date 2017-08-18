package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
	"xutils/xerr"
)

/*
CREATE TABLE `devicemodel` (
  `Id` 		varchar(50) NOT NULL,
  `PId` 	varchar(50) DEFAULT NULL,
  `Name` 	varchar(100) DEFAULT NULL,
  `Description` varchar(500) DEFAULT NULL,
  `Type` 		int(11) DEFAULT NULL COMMENT '1-分组 2-设备分类',
  `PageFileName` 	varchar(50) DEFAULT NULL,
  `ImgFileName` 	varchar(50) DEFAULT NULL,
  `ImgFileName2` 	varchar(50) DEFAULT NULL,
  `IsAlert` 		varchar(1) DEFAULT NULL,
  `MaxUseTime` 		bigint(20) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='只要是同一个型号，操作指令肯定是一样的';
*/
type DeviceModel struct {
	Id           string
	PId          string
	Name         string
	Description  string
	Type         int
	PageFileName string
	ImgFileName  string
	ImgFileName2 string
	IsAlert      string
	MaxUseTime   int
}

// 查询DeviceModel
func DeviceModel_Get(id string, pDBMap *gorp.DbMap) (pModel *DeviceModel, err error) {
	pDBMap.AddTableWithName(DeviceModel{}, "devicemodel").SetKeys(false, "Id")
	pObj, err := pDBMap.Get(DeviceModel{}, id)

	pModel, _ = pObj.(*DeviceModel)
	return pModel, err
}

// 查询父ID对应的DeviceModel存在?
func DeviceModel_Exists_ByParentID(id string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM devicemodel WHERE PId = ?`
	count, err := pDBMap.SelectInt(query, id)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 插入
func (p *DeviceModel) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceModel{}, "devicemodel").SetKeys(false, "Id")
	return pDBMap.Insert(p)
}

// 删除
func (p *DeviceModel) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModel{}, "devicemodel").SetKeys(false, "Id")
	rows, err := pDBMap.Delete(p)
	return int(rows), err
}

func (p *DeviceModel) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModel{}, "devicemodel").SetKeys(false, "Id")
	rows, err := pDBMap.Update(p)
	return int(rows), err
}

// 获取子集DeviceModel
func (p *DeviceModel) GetChildren(pDBMap *gorp.DbMap) (list []DeviceModel, err error) {
	list = []DeviceModel{}
	statement := "SELECT * FROM devicemodel WHERE PId = ?"
	_, err = pDBMap.Select(&list, statement, p.Id)
	return list, err
}

// 递归删除设备型号
func (p *DeviceModel) DeleteRecursively(pDBMap *gorp.DbMap) (ok bool) {
	chidren, err := p.GetChildren(pDBMap)
	xerr.ThrowPanic(err)
	for _, model := range chidren {
		model.DeleteRecursively(pDBMap)
	}

	_, err = p.Delete(pDBMap)
	xerr.ThrowPanic(err)
	return true
}
