package zndx

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"xutils/xtime"

	"gopkg.in/gorp.v1"
)

/*
CREATE TABLE `node_upinfo` (
  `node_id` varchar(50) NOT NULL COMMENT 'Eui',
  `prev_uptime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '前一次上报时间',
  `last_uptime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '最后上报时间',
  `last_upinfo` varchar(500) NOT NULL DEFAULT '' COMMENT '最后一次上报信息',
  `upinfo_list` varchar(5000) NOT NULL DEFAULT '' COMMENT '节点上报信息列表(Json数组形式，用来保存一段时间内的节点上报信息)',
  PRIMARY KEY (`node_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='节点上报信息表';
*/
type NodeUpinfo struct {
	NodeID     string `db:"node_id"`
	PrevTime   string `db:"prev_uptime"`
	LastTime   string `db:"last_uptime"`
	LastUpinfo string `db:"last_upinfo"`
	UpinfoList string `db:"upinfo_list"` // json
}

// 获取NodeUpinfo
func NodeUpinfo_Get(nodeID string, pDBMap *gorp.DbMap) (pUpinfo *NodeUpinfo, err error) {
	pDBMap.AddTableWithName(NodeUpinfo{}, "node_upinfo").SetKeys(false, "node_id")
	pObj, err := pDBMap.Get(NodeUpinfo{}, nodeID)
	pUpinfo, _ = pObj.(*NodeUpinfo)

	return pUpinfo, err
}

// 查询存在?
func NodeUpinfo_Exists(nodeID string, pDBMap *gorp.DbMap) (yes bool) {
	query := `SELECT COUNT(1) FROM node_upinfo TInfo WHERE TInfo.node_id = ?`
	count, err := pDBMap.SelectInt(query, nodeID)
	if err == nil {
		yes = count > 0
	}

	return yes
}

// 插入
func (p *NodeUpinfo) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(NodeUpinfo{}, "node_upinfo").SetKeys(false, "node_id")
	return pDBMap.Insert(p)
}

// 删除
func (p *NodeUpinfo) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(NodeUpinfo{}, "node_upinfo").SetKeys(false, "node_id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 更新
func (p *NodeUpinfo) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(NodeUpinfo{}, "node_upinfo").SetKeys(false, "node_id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

// 保存
func (p *NodeUpinfo) Save(pDBMap *gorp.DbMap) (err error) {
	if NodeUpinfo_Exists(p.NodeID, pDBMap) {
		_, err = p.Update(pDBMap)
		return err
	}

	return p.Insert(pDBMap)
}

// 在线?
func (p *NodeUpinfo) Onlined(d time.Duration) (onlined bool) {
	if last_time, err := time.Parse(xtime.FormatString(), p.LastTime); err == nil {
		prev_time, err := time.Parse(xtime.FormatString(), p.PrevTime)
		if (err == nil) && (last_time.Sub(prev_time) < d) {
			onlined = true
		}
	}

	return onlined
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
