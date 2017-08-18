package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `nodemodelcmd` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `ModelId` varchar(50) NOT NULL,
  `CmdCode` varchar(50) NOT NULL COMMENT '操作例子或设备的命令简码，如：打开盒子用on、关闭电源用off，获取状态用status',
  `CmdName` varchar(50) DEFAULT NULL,
  `RequestURI` varchar(50) DEFAULT NULL COMMENT '用于定义访问地址中<主机：端口>后面，<问号>前面的部分，如：config/smart_switch/projector/；可以用中括号来定义变量占位，如：[p1],命令调时由客户端传入对应值',
  `URIQuery` varchar(500) DEFAULT NULL COMMENT 'URI中问号后面的部分，可以用中括号来定义变量占位，如：[p1],命令调时由客户端传入对应值',
  `CmdDescription` varchar(500) DEFAULT NULL,
  `RequestType` varchar(50) DEFAULT NULL COMMENT 'get/post/put/delete',
  `Payload` varchar(500) DEFAULT NULL COMMENT '负载；可以用中括号来定义变量占位，如：[p1],命令调时由客户端传入对应值',
  `CloseCmdFlag` varchar(1) DEFAULT NULL,
  `OpenCmdFlag` varchar(1) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;
*/
type NodeModelCmd struct {
	Id             int
	ModelId        string
	CmdCode        string
	CmdName        string
	RequestURI     string
	URIQuery       string
	CmdDescription string
	RequestType    string
	Payload        string
	CloseCmdFlag   string
	OpenCmdFlag    string
}

// 获取NodeModelCmd
func NodeModelCmd_Get(id int, pDBMap *gorp.DbMap) (pCmd *NodeModelCmd, err error) {
	pDBMap.AddTableWithName(NodeModelCmd{}, "nodemodelcmd").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(NodeModelCmd{}, id)

	pCmd, _ = pObj.(*NodeModelCmd)
	return pCmd, err
}

// 查询存在NodeModelCmd?(ByModelID)
func NodeModelCmd_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM nodemodelcmd WHERE ModelId = ?`
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 插入
func (p *NodeModelCmd) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(NodeModelCmd{}, "nodemodelcmd").SetKeys(true, "Id")
	return pDBMap.Insert(p)
}

// 删除
func (p *NodeModelCmd) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(NodeModelCmd{}, "nodemodelcmd").SetKeys(true, "Id")
	rows, err := pDBMap.Delete(p)
	return int(rows), err
}

// 更新
func (p *NodeModelCmd) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(NodeModelCmd{}, "nodemodelcmd").SetKeys(true, "Id")
	rows, err := pDBMap.Update(p)
	return int(rows), err
}

// 保存
func (p *NodeModelCmd) Save(pDBMap *gorp.DbMap) (err error) {
	if p.Exists(pDBMap) {
		sql := `UPDATE nodemodelcmd SET ModelId = ?, CmdCode = ?, CmdName = ?, RequestURI = ?, URIQuery = ?, CmdDescription = ?, RequestType = ? WHERE (Id = ?);`
		args := []interface{}{p.ModelId, p.CmdCode, p.CmdName, p.RequestURI, p.URIQuery, p.CmdDescription, p.RequestType, p.Id}
		_, err = pDBMap.Exec(sql, args...)
		return err
	}

	err = p.Insert(pDBMap)
	return err
}

// 存在?
func (p *NodeModelCmd) Exists(dbmap *gorp.DbMap) (exist bool) {
	sql := "SELECT 1 FROM nodemodelcmd WHERE Id = ? "
	nullInt, err := dbmap.SelectNullInt(sql, p.Id)
	exist = (err == nil) && (nullInt.Valid)
	return exist
}

// 查询命令被设备引用否
func (p *NodeModelCmd) Referenced_ByDevice(pDBMap *gorp.DbMap) (yes bool) {
	yes = true
	query := `
SELECT COUNT(1) FROM nodemodelcmd TCmd
	JOIN nodemodel TModel ON(TModel.Id = TCmd.ModelId)
	JOIN node TNode ON(TNode.Id IN(SELECT DISTINCT PowerNodeId FROM device) OR TNode.Id IN(SELECT DISTINCT JoinNodeId FROM device))
WHERE TCmd.Id = ?
`
	count, err := pDBMap.SelectInt(query, p.Id)
	xdebug.LogError(err)
	if err == nil {
		yes = count > 0
	}

	return yes
}
