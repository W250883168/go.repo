package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `nodesocketstatus` (
  `NodeId` varchar(50) NOT NULL,
  `SocketId` varchar(50) NOT NULL COMMENT '1/2/3',
  `SeqNo` int(11) NOT NULL COMMENT '顺序号（从1开始，依次2、3、4.。。）',
  `StatusValue` varchar(50) DEFAULT NULL COMMENT 'PWR=00 / ERR:',
  `UpdateTime` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`NodeId`,`SocketId`,`SeqNo`)
) ENGINE=MyISAM DEFAULT CHARSET=gbk COMMENT='\r\n';

*/
type NodeSocketStatus struct {
	NodeId      string //
	SocketId    string // 插口ID
	SeqNo       int    // 状态顺序编号
	StatusValue string // 状态取值
	UpdateTime  string // 更新时间
}

// 删除
func (p *NodeSocketStatus) DeleteByNodeID(dbmap *gorp.DbMap) (err error) {
	sql := "DELETE FROM NodeSocketStatus WHERE NodeId = ? "
	_, err = dbmap.Exec(sql, p.NodeId)
	return err
}

// 插入
func (p *NodeSocketStatus) Insert(dbmap *gorp.DbMap) (err error) {
	sql := "INSERT INTO NodeSocketStatus (NodeId,SocketId,SeqNo,StatusValue,UpdateTime) VALUES (?,?,?,?,?)"
	args := []interface{}{p.NodeId, p.SocketId, p.SeqNo, p.StatusValue, p.UpdateTime}
	_, err = dbmap.Exec(sql, args...)
	return err
}

// 查询节点的某个插口的状态
func (p *NodeSocketStatus) Query_ByNodeInfo(pDBMap *gorp.DbMap) (list []NodeSocketStatus, err error) {
	list = []NodeSocketStatus{}

	sql := "select NodeId,SocketId,SeqNo,StatusValue,UpdateTime from NodeSocketStatus where NodeId=? and SocketId=? order by SeqNo asc"
	_, err = pDBMap.Select(&list, sql, p.NodeId, p.SocketId)

	return list, err
}

// 查询存在NodeSocketStatus?
func NodeSocketStatus_Exists_ByNodeID(nodeID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM nodesocketstatus WHERE NodeId = ?`
	count, err := pDBMap.SelectInt(query, nodeID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}
