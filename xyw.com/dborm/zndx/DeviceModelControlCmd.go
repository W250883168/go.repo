package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `devicemodelcontrolcmd` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ModelId` varchar(50) NOT NULL,
  `CmdCode` varchar(50) NOT NULL COMMENT '操作例子或设备的命令简码，如：打开盒子用on、关闭电源用off，获取状态用status',
  `CmdName` varchar(50) DEFAULT NULL,
  `RequestURI` varchar(50) DEFAULT NULL COMMENT '用于定义访问地址中<主机：端口>后面，<问号>前面的部分，如：config/smart_switch/projector/；可以用中括号来定义变量占位，如：[p1],命令调时由客户端传入对应值',
  `URIQuery` varchar(500) DEFAULT NULL COMMENT 'URI中问号后面的部分，可以用中括号来定义变量占位，如：[p1],命令调时由客户端传入对应值',
  `CmdDescription` varchar(500) DEFAULT NULL,
  `RequestType` varchar(50) DEFAULT NULL COMMENT 'get/post/put/delete',
  `Payload` varchar(500) DEFAULT NULL COMMENT '负载；可以用中括号来定义变量占位，如：[p1],命令调时由客户端传入对应值',
  `DelayMillisecond` int(20) DEFAULT NULL,
  `CloseCmdFlag` varchar(1) DEFAULT NULL,
  `OpenCmdFlag` varchar(1) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
*/
type DeviceModelControlCmd struct {
	Id               int
	ModelId          string
	CmdCode          string
	CmdName          string
	RequestURI       string
	URIQuery         string
	CmdDescription   string
	RequestType      string
	Payload          string
	DelayMillisecond int // 延迟(毫秒)
	CloseCmdFlag     string
	OpenCmdFlag      string
}

// 获取DeviceModelControlCmd
func DeviceModelControlCmd_Get(id int, pDBMap *gorp.DbMap) (pCmd *DeviceModelControlCmd, err error) {
	pDBMap.AddTableWithName(DeviceModelControlCmd{}, "devicemodelcontrolcmd").SetKeys(true, "id")
	pObj, err := pDBMap.Get(DeviceModelControlCmd{}, id)

	pCmd, _ = pObj.(*DeviceModelControlCmd)
	return pCmd, err
}

// 查询对应模型ID的DeviceModelControlCmd存在?
func DeviceModelControlCmd_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicemodelcontrolcmd WHERE ModelId = ?"
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 插入
func (p *DeviceModelControlCmd) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceModelControlCmd{}, "devicemodelcontrolcmd").SetKeys(true, "id")
	return pDBMap.Insert(p)
}

// 删除
func (p *DeviceModelControlCmd) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModelControlCmd{}, "devicemodelcontrolcmd").SetKeys(true, "id")
	rows, err := pDBMap.Delete(p)
	return int(rows), err
}

// 更新
func (p *DeviceModelControlCmd) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceModelControlCmd{}, "devicemodelcontrolcmd").SetKeys(true, "id")
	rows, err := pDBMap.Update(p)
	return int(rows), err
}

// 查询是否被引用?
func (p *DeviceModelControlCmd) Referenced_ByDevice(pDBMap *gorp.DbMap) (yes bool) {
	yes = true
	query := `
SELECT COUNT(1) FROM devicemodelcontrolcmd TCmd
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
