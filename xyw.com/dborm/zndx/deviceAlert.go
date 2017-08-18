package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
	"xutils/xtime"
)

/*
CREATE TABLE `devicealert` (
  `Id` bigint(50) NOT NULL AUTO_INCREMENT,
  `DeviceId` varchar(50) DEFAULT NULL,
  `AlertType` varchar(50) DEFAULT NULL COMMENT '1-超时使用预警 2-状态值条件预警 3-具体状态值预警',
  `StatusCode` varchar(50) DEFAULT NULL,
  `StatusValueCode` varchar(50) DEFAULT NULL,
  `AlertDescription` varchar(1000) DEFAULT NULL COMMENT '为预警信息生成一句话的描述',
  `LastAlertTime` varchar(50) DEFAULT NULL COMMENT '记录最后一次的警告时间',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=753 DEFAULT CHARSET=gbk COMMENT='对于状态值条件预警和具体状态值预警：一个设备的相同状态编码和状态值编码，只存在一条；\r\n对于超时使用警告：一个';

对于状态值条件预警和具体状态值预警：一个设备的相同状态编码和状态值编码，只存在一条；
对于超时使用警告：一个设备的超时使用警告只存在一条；
节点每10秒上传一次数据，就在上传数据时进行设备的预警处理
*/
type DeviceAlert struct {
	Id               int
	DeviceId         string // 设备ID
	AlertType        string // 预警分类
	StatusCode       string // 状态代码
	StatusValueCode  string // 状态值编码
	AlertDescription string // 预警描述
	LastAlertTime    string // 上次预警时间
}

// 插入
func (p *DeviceAlert) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceAlert{}, "devicealert").SetKeys(true, "Id")
	err = pDBMap.Insert(p)
	return err
}

// 删除
func (p *DeviceAlert) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceAlert{}, "devicealert").SetKeys(true, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *DeviceAlert) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceAlert{}, "devicealert").SetKeys(true, "Id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

// 删除设备对应的预警记录(按状态码, AlertType=2)
func (p *DeviceAlert) Delete_ByStatusCodeAndDevice(pDev *Device, pDBMap *gorp.DbMap) (err error) {
	sql := `
DELETE FROM DeviceAlert 
WHERE AlertType = '2'
		AND StatusCode =?
		AND DeviceId IN (SELECT Id FROM Device WHERE JoinNodeId =? AND JoinSocketId =?)
`
	_, err = pDBMap.Exec(sql, p.StatusCode, pDev.JoinNodeId, pDev.JoinSocketId)
	return err
}

// 删除设备对应的预警记录(按具体状态值, AlertType=3)
func (p *DeviceAlert) Delete_ByStatusValueAndDevice(pDev *Device, pDBMap *gorp.DbMap) (err error) {
	sql := `
DELETE FROM DeviceAlert 
WHERE AlertType = '3'
		AND StatusCode =?
		AND StatusValueCode =?
		AND DeviceId IN (SELECT Id FROM Device WHERE JoinNodeId =? AND JoinSocketId =?)		
`
	args := []interface{}{p.StatusCode, p.StatusValueCode, pDev.JoinNodeId, pDev.JoinSocketId}
	_, err = pDBMap.Exec(sql, args...)
	return err
}

// 创建预警消息(通过设备连接节点ID和插口ID获得设备ID，根据这两个ID查询设备，看起来好象会找到多个设备，但实际上只有一个设备，因为一个插口只能接一个设备)
// 状态值条件预警(AlertType=2)
func (p *DeviceAlert) Insert_ByStatusCodeAndDevice(pDev *Device, pDBMap *gorp.DbMap) (err error) {
	sql := `
INSERT INTO DeviceAlert (DeviceId, AlertType, AlertDescription, LastAlertTime, StatusCode ) 
	SELECT id,?,?,?,? FROM Device WHERE JoinNodeId =? AND JoinSocketId =?
`
	args := []interface{}{"2", p.AlertDescription, xtime.NowString(), p.StatusCode, pDev.JoinNodeId, pDev.JoinSocketId}
	_, err = pDBMap.Exec(sql, args...)
	return err
}

//创建新的预警消息(通过连接节点ID和插口ID获得设备ID，根据这两个ID查询设备，看起来好象会找到多个设备，但实际上只有一个设备，因为一个插口只能接一个设备)
// 具体状态值预警(AlertType=3)
func (p *DeviceAlert) Insert_ByStatusValueAndDevice(pDev *Device, pDBMap *gorp.DbMap) (err error) {
	sql := `
INSERT INTO DeviceAlert (DeviceId, AlertType, AlertDescription, LastAlertTime, StatusCode, StatusValueCode) 
	SELECT id,?,?,?,?,? FROM Device WHERE JoinNodeId =? AND JoinSocketId =?
`
	args := []interface{}{"3", p.AlertDescription, p.LastAlertTime, p.StatusCode, p.StatusValueCode, pDev.JoinNodeId, pDev.JoinSocketId}
	_, err = pDBMap.Exec(sql, args...)
	return err
}

// 删除设备超时预警(AlertType = '1')
func (p *DeviceAlert) Delete_DeviceOvertime(device_id string, pDBMap *gorp.DbMap) (err error) {
	sql := `DELETE FROM DeviceAlert WHERE DeviceId =? AND AlertType = '1'`
	_, err = pDBMap.Exec(sql, device_id)

	return err
}

func DeviceAlert_Get(id int, pDBMap *gorp.DbMap) (pAlert *DeviceAlert, err error) {
	pDBMap.AddTableWithName(DeviceAlert{}, "devicealert").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceAlert{}, id)

	pAlert, _ = pObj.(*DeviceAlert)
	return pAlert, err
}

// 告警存在?
func DeviceAlert_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicealert WHERE DeviceId = ?"
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}
