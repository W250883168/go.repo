package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `devicemodelstatusvaluecode` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ModelId` varchar(100) DEFAULT NULL,
  `StatusCode` varchar(50) DEFAULT NULL,
  `StatusValueCode` varchar(50) DEFAULT NULL,
  `StatusValueName` varchar(500) DEFAULT NULL,
  `IsAlert` varchar(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=14 DEFAULT CHARSET=gbk;
*/
type DeviceModelStatusValueCode struct {
	Id              int64  // `"BIGINT(20)"`
	ModelId         string // `"VARCHAR(100)"`
	StatusCode      string // `"VARCHAR(50)"`
	StatusValueCode string // `"VARCHAR(50)"`
	StatusValueName string // `"VARCHAR(500)"`
	IsAlert         string // `"VARCHAR(1)"`
}

// 获取DeviceModelStatusValueCode
func DeviceModelStatusValueCode_Get(id string, pDBMap *gorp.DbMap) (pValueCode *DeviceModelStatusValueCode, err error) {
	pDBMap.AddTableWithName(DeviceModelStatusValueCode{}, "devicemodelstatusvaluecode").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceModelStatusValueCode{}, id)

	pValueCode, _ = pObj.(*DeviceModelStatusValueCode)
	return pValueCode, err
}

// 查询是否存在DeviceModelStatusValueCode(通过模型ID)
func DeviceModelStatusValueCode_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicemodelstatusvaluecode WHERE ModelId = ?"
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 查询是否存在DeviceModelStatusValueCode(ByStatusCmd)
func DeviceModelStatusValueCode_Exists_ByStatusCmd(pCmd *DeviceModelStatusCmd, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM devicemodelstatusvaluecode TCode WHERE TCode.ModelId = ? AND TCode.StatusCode = ?`
	count, err := pDBMap.SelectInt(query, pCmd.ModelId, pCmd.StatusCode)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

//根据状态值代码从状态表中查询参数值名称
func (p *DeviceModelStatusValueCode) Query_StatusValueName(dbmap *gorp.DbMap) (name string, err error) {
	data := struct{ StatusValueName string }{}
	sql := "select  StatusValueName from DeviceModelStatusValueCode where ModelId=? and StatusCode=? and StatusValueCode=?"
	err = dbmap.SelectOne(&data, sql, p.ModelId, p.StatusCode, p.StatusValueCode)

	return data.StatusValueName, err
}
