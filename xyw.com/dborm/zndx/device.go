package zndx

import (
	"time"

	"gopkg.in/gorp.v1"

	"xutils/xdebug"
	"xutils/xtext"
	"xutils/xtime"
)

/*************************************************************
CREATE TABLE `device` (
  `Id` varchar(50) NOT NULL,
  `Name` varchar(100) DEFAULT NULL,
  `Sn` varchar(50) DEFAULT NULL COMMENT '设备出厂时的序列号',
  `Code` varchar(50) DEFAULT NULL COMMENT '基于资产管理，学校会给每台设备分配一个唯一的编号',
  `Brand` varchar(50) DEFAULT NULL,
  `ModelId` varchar(50) DEFAULT NULL COMMENT 'DeviceModel.id',
  `ClassroomId` int(11) DEFAULT NULL,
  `PowerNodeId` varchar(50) DEFAULT NULL,
  `PowerSwitchId` varchar(50) DEFAULT NULL COMMENT '一个节点可能有多个开关，本字段记录设备电源所连接的节点开关',
  `JoinMethod` varchar(50) DEFAULT NULL,
  `JoinNodeId` varchar(50) DEFAULT NULL,
  `JoinSocketId` varchar(50) DEFAULT NULL COMMENT '一个节点可能有多个插口，每一个插口都可以连接一个设备，本字段记录设备所连接的插口id',
  `NodeSwitchStatus` varchar(50) DEFAULT NULL,
  `NodeSwitchStatusUpdateTime` varchar(50) DEFAULT NULL,
  `DeviceSelfStatus` varchar(50) DEFAULT NULL,
  `DeviceSelfStatusUpdateTime` varchar(50) DEFAULT NULL,
  `IsCanUse` varchar(1) DEFAULT NULL,
  `UseTimeBefore` bigint(20) DEFAULT NULL,
  `UseTimeAfter` bigint(20) DEFAULT NULL,
  `JoinNodeUpdateTime` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=gbk COMMENT='存储设备的基本信息，同时需指定设备属于的型号，以及设备所连接的盒子';

存储设备的基本信息，
1.设备连接方式有多种：
1）node连接(node)：即设备连接到新云网自己开发的节点（socketid）进行控制，此时需要设置nodeid and socketid
2）pjlink连接(pjlink)：通过网络(RJ45)连接，控制遵循pjlink的投影仪设备，此时不需要设置socketid,如果电源不插在节点上，则nodeid都不用设置
3）...(其它以后新发现的连接
*****************************************************************/

type Device struct {
	Id                         string
	Name                       string // 设备名称
	Sn                         string // 出厂序列号
	Code                       string // 设备编号
	Brand                      string // 设备品牌
	ModelId                    string // 设备型号ID
	ClassroomId                int    // 设备所在教室ID
	PowerNodeId                string // 设备电源节点ID
	PowerSwitchId              string // 设备连接的节点开关ID
	JoinMethod                 string // 设备接入方式
	JoinNodeId                 string // 设备接入-节点接入-节点ID
	JoinSocketId               string // 设备接入-节点接入-插口ID
	NodeSwitchStatus           string // 设备连接的节点开关状态
	NodeSwitchStatusUpdateTime string // 节点开关状态更新时间
	DeviceSelfStatus           string // 设备自身开关状态
	DeviceSelfStatusUpdateTime string // 设备自身开关状态更新时间
	IsCanUse                   string // 设备是否可用
	UseTimeBefore              int64  // 上系统前设备已使用时间（秒）
	UseTimeAfter               int64  // 上系统后设备已使用时间（秒）
	JoinNodeUpdateTime         string // 设备接入节点最后上报时间
}

// 根据ID获取设备
func Device_Get(id string, pDBMap *gorp.DbMap) (pDev *Device, err error) {
	pDBMap.AddTableWithName(Device{}, "device").SetKeys(false, "Id")
	pObject, err := pDBMap.Get(Device{}, id)
	xdebug.LogErrorText(err)

	pDev, _ = pObject.(*Device)
	return pDev, err
}

func Device_Exists_ByNodeID(nodeID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM device WHERE PowerNodeId = ? OR JoinNodeId = ?`
	count, err := pDBMap.SelectInt(query, nodeID, nodeID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

func Device_Exists_ByModelID(modelID string, pDBMap *gorp.DbMap) (exist bool) {
	query := `SELECT COUNT(1) FROM device WHERE ModelId = ?`
	count, err := pDBMap.SelectInt(query, modelID)
	xdebug.LogError(err)

	return count > 0
}

// 按节点查找设备
func Device_QueryByNodeID(nodeID string, pDBMap *gorp.DbMap) (list []Device, err error) {
	list = []Device{}
	query := `SELECT TDev.* FROM device TDev WHERE ? IN(TDev.PowerNodeId, TDev.JoinNodeId)`
	_, err = pDBMap.Select(&list, query, nodeID)

	return list, err
}

// 查找房间中有电源节点的设备
func Device_QueryByRoomID(room_id string, pDBMap *gorp.DbMap) (list []Device, err error) {
	sql := `
SELECT TDev.* FROM device TDev JOIN node TNode ON (TNode.Id = TDev.PowerNodeId)
WHERE TDev.ClassroomId = :ClassroomID
`
	params := map[string]interface{}{}
	params["ClassroomID"] = room_id
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}

// 按楼层查找设备
func Device_QueryByFloorID(floor_id string, pDBMap *gorp.DbMap) (list []Device, err error) {
	sql := `
SELECT TDev.* FROM device TDev 
	JOIN node TNode ON (TNode.Id = TDev.PowerNodeId)
	JOIN classrooms TRoom ON (TRoom.Id=TDev.ClassroomId)	
WHERE TRoom.Floorsid = :FloorID
`
	params := map[string]interface{}{}
	params["FloorID"] = floor_id
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}

func (p *Device) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(Device{}, "device").SetKeys(false, "Id")
	return pDBMap.Insert(p)
}

func (p *Device) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(Device{}, "device").SetKeys(false, "Id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

func (p *Device) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(Device{}, "device").SetKeys(false, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 根据设备ID得到设备详细名称（例如：新校区A栋1层101投影仪)
func (p *Device) GetDeviceDetailName(dbmap *gorp.DbMap) string {
	sql := `SELECT CONCAT_WS('', c.Campusname, b.Buildingname, f.Floorname, r.Classroomsname, d.DeviceName) Name
			FROM (SELECT ClassroomId,Name DeviceName FROM device WHERE Id = ?) d
				LEFT JOIN classrooms r ON r.Id = d.ClassroomId
				LEFT JOIN floors f ON f.Id = r.Floorsid
				LEFT JOIN building b ON b.Id = f.Buildingid
				LEFT JOIN campus c ON c.Id = b.Campusid
			`
	args := []interface{}{p.Id}
	name, err := dbmap.SelectStr(sql, args...)
	xdebug.LogError(err)
	return name
}

// 检查给设备供电的节点是否在线
func (p *Device) CheckSwitchOnlined(offTimeout time.Duration) (onlined bool) {
	if xtext.IsBlank(p.PowerNodeId) {
		onlined = true
		return onlined
	}

	if xtext.IsNotBlank(p.NodeSwitchStatusUpdateTime) {
		tFormat := xtime.FormatString()
		tUpdate, err := time.Parse(tFormat, p.NodeSwitchStatusUpdateTime)
		if err == nil {
			onlined = time.Now().After(tUpdate.Add(offTimeout))
		}
	}

	return onlined
}

// 检查设备接入节点（RS232资源、红外资源等）是否在线
func (p *Device) CheckNodeOnlined(offTimeout time.Duration) (onlined bool) {
	if xtext.IsBlank(p.JoinNodeId) {
		onlined = true
		return onlined
	}

	if xtext.IsNotBlank(p.JoinNodeUpdateTime) { //不为空，说明有更新过
		tFormat := xtime.FormatString()
		tUpdate, err := time.Parse(tFormat, p.JoinNodeUpdateTime)
		if err == nil {
			onlined = time.Now().After(tUpdate.Add(offTimeout))
		}
	}

	return onlined
}

// 根据离线时长，判断设备是否在线
func (p *Device) Onlined(offTimeout time.Duration) bool {
	return p.CheckSwitchOnlined(offTimeout) && p.CheckNodeOnlined(offTimeout)
}

// 更新设备自身状态
func (p *Device) Update_DeviceSelfStatus(pDBMap *gorp.DbMap) (err error) {
	sql := "UPDATE device SET DeviceSelfStatus = ?, DeviceSelfStatusUpdateTime = ? WHERE Id = ?"
	sqlargs := []interface{}{p.DeviceSelfStatus, xtime.NowString(), p.Id}

	_, err = pDBMap.Exec(sql, sqlargs...)
	return err
}

// 更新给设备供电的节点开关状态
func (p *Device) Update_NodeSwitchStatus(pDBMap *gorp.DbMap) (err error) {
	sql := "UPDATE device SET NodeSwitchStatus =?, NodeSwitchStatusUpdateTime =? WHERE PowerNodeId =? AND PowerSwitchId =?"
	args := []interface{}{p.NodeSwitchStatus, p.NodeSwitchStatusUpdateTime, p.PowerNodeId, p.PowerSwitchId}
	_, err = pDBMap.Exec(sql, args...)

	return err
}

// 更新接入设备的最后上报时间
func (p *Device) Update_JoinNodeUpdateTime(pDBMap *gorp.DbMap) (err error) {
	sql := "UPDATE device SET JoinNodeUpdateTime=? WHERE IFNULL(JoinMethod,'')!='' and IFNULL(JoinNodeId,'')=?"
	_, err = pDBMap.Exec(sql, p.JoinNodeUpdateTime, p.JoinNodeId)

	return err
}

// 更新设备累计使用时间
func (p *Device) Update_DeviceUseTime(pDBMap *gorp.DbMap) (err error) {
	sql := `UPDATE Device SET UseTimeAfter = IFNULL(UseTimeAfter, 0) +? WHERE Id =?`
	_, err = pDBMap.Exec(sql, p.UseTimeAfter, p.Id)
	return err
}
