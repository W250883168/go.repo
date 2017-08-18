package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `devicemodelstatuscmd` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ModelId` varchar(100) DEFAULT NULL,
  `Payload` varchar(100) DEFAULT NULL COMMENT '对于投影仪：PWR?\r    ERR?\r   SOURCE?\r  LAMP?\r',
  `StatusName` varchar(50) DEFAULT NULL COMMENT '灯泡/错误/信号/开关',
  `StatusCode` varchar(50) DEFAULT NULL COMMENT 'LAMP/ERR/SOURCE/SWITCH',
  `StatusValueMatchString` varchar(200) DEFAULT NULL COMMENT '使用用该字符串从返回的设备数据中取出状态值，匹配串格式如：PWR={val}',
  `SwitchStatusFlag` varchar(1) DEFAULT NULL COMMENT '是否是开关命令：0-否 1-是 ，一个型号所有状态命令中最多只有一条命令设为是(1)',
  `OnValue` varchar(100) DEFAULT NULL COMMENT '设备开启状态的值是多少，例如投影仪，返回的PWR=01，01表示开',
  `OffValue` varchar(100) DEFAULT NULL COMMENT '设备关闭状态的值是多少，例如投影仪，返回的PWR=00时，00表示关闭',
  `SeqNo` int(11) DEFAULT NULL COMMENT '每一个设备型号从1开始升序编号，依次对应设备返回数据的cmd0_res  cmd1_res cmd2_res ...',
  `SelectValueFlag` varchar(1) DEFAULT NULL COMMENT '0-不从编码表取值 1-从编码表取值  (按匹配串取出值后，是否从编码表取对应的值)',
  `IsAlert` varchar(1) DEFAULT NULL,
  `AlertWhere` varchar(200) DEFAULT NULL,
  `AlertDescription` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='这里定义的命令，会驻留在节点上（由管理员配置好后发送到节点保存），由节点自动获取，并每十秒自动上报';
*/
type DeviceModelStatusCmd struct {
	Id                     int64  // `"BIGINT(20)"`
	ModelId                string // `"VARCHAR(100)"`
	Payload                string // `"VARCHAR(100)"`
	StatusName             string // `"VARCHAR(50)"`
	StatusCode             string // `"VARCHAR(50)"`
	StatusValueMatchString string // `"VARCHAR(200)"`
	SwitchStatusFlag       string // `"VARCHAR(1)"`
	OnValue                string // `"VARCHAR(100)"`
	OffValue               string // `"VARCHAR(100)"`
	SeqNo                  int    // `"INT(11)"`
	SelectValueFlag        string // `"VARCHAR(1)"`
	IsAlert                string // `"VARCHAR(1)"`
	AlertWhere             string // `"VARCHAR(200)"`
	AlertDescription       string // `"VARCHAR(500)"`
}

// 查询是否被设备引用?
func (p *DeviceModelStatusCmd) Referenced_ByDevice(pDBMap *gorp.DbMap) (yes bool) {
	yes = true
	query := `
SELECT COUNT(1) FROM devicemodelstatuscmd TCmd
	JOIN devicemodel TModel ON(TModel.Id = TCmd.ModelId)
	JOIN device TDev ON(TDev.ModelId = TModel.Id)
WHERE TCmd.id = ?
`
	count, err := pDBMap.SelectInt(query, p.Id)
	xdebug.LogError(err)
	if err == nil {
		yes = count > 0
	}

	return yes
}

// 获取DeviceModelStatusCmd
func DeviceModelStatusCmd_Get(id string, pDBMap *gorp.DbMap) (pCmd *DeviceModelStatusCmd, err error) {
	pDBMap.AddTableWithName(DeviceModelStatusCmd{}, "devicemodelstatuscmd").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceModelStatusCmd{}, id)

	pCmd, _ = pObj.(*DeviceModelStatusCmd)
	return pCmd, err
}

// 通过模型ID, 查询存在DeviceModelStatusCmd
func DeviceModelStatusCmd_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicemodelstatuscmd WHERE ModelId = ?"
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 查询状态命令配置（报警使用）(根据设备接入节点id和插口id)
func DeviceModelStatusCmd_Query_ByDevice(pDev *Device, dbmap *gorp.DbMap) (list []DeviceModelStatusCmd, err error) {
	sql := `
SELECT Id, ModelId, SeqNo, SwitchStatusFlag,
	IFNULL(Payload, '') Payload,
	IFNULL(StatusName, '') StatusName,
	IFNULL(StatusCode, '') StatusCode,
	IFNULL(StatusValueMatchString, '') StatusValueMatchString,	
	IFNULL(OnValue, '') OnValue,
	IFNULL(OffValue, '') OffValue,	
	IFNULL(SelectValueFlag, '0') SelectValueFlag,
	IFNULL(IsAlert, '0') IsAlert,
	IFNULL(AlertWhere, '') AlertWhere,
	IFNULL(AlertDescription, '') AlertDescription
FROM DeviceModelStatusCmd
WHERE	ModelId IN (SELECT ModelId FROM Device WHERE JoinNodeId =? AND JoinSocketId =?)
ORDER BY SeqNo
`
	_, err = dbmap.Select(&list, sql, pDev.JoinNodeId, pDev.JoinSocketId)
	return list, err
}
