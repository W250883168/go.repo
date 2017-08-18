package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `node` (
  `Id` varchar(50) NOT NULL,
  `Name` varchar(100) DEFAULT NULL COMMENT '给同一个教室的节点取个名字，便于区分。',
  `ModelId` varchar(50) DEFAULT NULL,
  `ClassRoomId` int(11) DEFAULT NULL,
  `IpType` varchar(50) DEFAULT NULL COMMENT 'ip4/ip6（使用ip4时，读取本表Ip字段值，使用ip6时，根据Eui转化',
  `NodeCoapPort` varchar(50) DEFAULT NULL COMMENT 'coap端口，默认5683，IPV6时使用',
  `InRouteMappingPort` varchar(50) DEFAULT NULL COMMENT '盒子每10秒上传时自动更新：从协议中读取，如20005，IPV4时使用',
  `RouteIp` varchar(50) DEFAULT NULL COMMENT '盒子每10秒上传时自动更新：从协议中读取',
  `UploadTime` varchar(50) DEFAULT NULL COMMENT '盒子每10秒上传时自动更新：取当前时间',
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=gbk COMMENT='';

这里的节点指由新云网研发的智能节点，它一端通过无线连接路由器与服务器通讯，一端连接设备对其进行控制（红外线控制或其它控制如RS232)，
节点每10秒钟自动向服务器报告自己的状态，该表始终存储最新的状态
*/
type Node struct {
	Id                 string
	Name               string // 节点名称
	ModelId            string // 节点型号ID
	ClassRoomId        int    // 节点所在教室ID
	IpType             string // 使用IP类型(ipv4/ipv6)
	NodeCoapPort       string // CoAP端口号
	InRouteMappingPort string // 节点在路由器上的映射端口
	RouteIp            string // 节点连接的路由器IP
	UploadTime         string // 上报时间
}

// 查询节点是否存在（通过ModelID）
func Node_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM node WHERE ModelId = ?`
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 获取Node
func Node_Get(id string, pDBMap *gorp.DbMap) (pNode *Node, err error) {
	pDBMap.AddTableWithName(Node{}, "node").SetKeys(false, "Id")
	pObj, err := pDBMap.Get(Node{}, id)

	pNode, _ = pObj.(*Node)
	return pNode, err
}

// 存在?
func Node_Exists(id string, dbmap *gorp.DbMap) (exist bool) {
	sql := "SELECT COUNT(1) FROM node WHERE Id = ? "
	count, err := dbmap.SelectInt(sql, id)
	if err == nil {
		exist = (count > 0)
	}

	return exist
}

// 保存
func (p *Node) Save(dbmap *gorp.DbMap) (err error) {
	if Node_Exists(p.Id, dbmap) {
		_, err = p.Update(dbmap)
		return
	}

	err = p.Insert(dbmap)
	return err
}

// 保存上报信息
func (p *Node) SavePingInfo(dbmap *gorp.DbMap) (err error) {
	if !Node_Exists(p.Id, dbmap) {
		return p.Insert(dbmap)
	}

	sql := `
UPDATE node 
SET IpType = ?,
	NodeCoapPort = ?,
	InRouteMappingPort = ?,
	RouteIp = ?,
	UploadTime = ?
WHERE (Id = ?);`
	args := []interface{}{p.IpType, p.NodeCoapPort, p.InRouteMappingPort, p.RouteIp, p.UploadTime, p.Id}
	_, err = dbmap.Exec(sql, args...)
	return err
}

// 更新
func (p *Node) Update(dbmap *gorp.DbMap) (affects int, err error) {
	dbmap.AddTableWithName(Node{}, "node").SetKeys(false, "Id")
	rows, err := dbmap.Update(p)

	return int(rows), err
}

// 插入
func (p *Node) Insert(dbmap *gorp.DbMap) (err error) {
	dbmap.AddTableWithName(Node{}, "node").SetKeys(false, "Id")
	err = dbmap.Insert(p)

	return err
}

// 删除
func (p *Node) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(Node{}, "node").SetKeys(false, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}
