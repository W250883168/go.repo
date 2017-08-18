package dal

import (
	gorp "gopkg.in/gorp.v1"
)

type DevicePowerCmdView struct {
	// zndx.Device
	DeviceID         string `zndx.Device.Id`
	DeviceName       string `zndx.Device.Name`    // 设备名称
	DeviceModelID    string `zndx.Device.ModelId` // 设备型号ID
	PowerNodeID      string // 设备电源节点ID
	PowerSwitchIndex string `zndx.Device.PowerSwitchId` // 设备连接的节点开关ID
	JoinMethod       string // 设备接入方式
	JoinNodeID       string // 设备接入-节点接入-节点ID
	JoinSocketID     string // 设备接入-节点接入-插口ID

	// zndx.Node(接入节点)
	NodeID             string `zndx.Node.Id`
	NodeName           string `zndx.Node.Name`        // 节点名称
	NodeModelID        string `zndx.Node.NodeModelId` // 节点型号ID
	NodeClassroomID    int    `zndx.Node.ClassroomId` // 节点所在教室ID
	IPType             string `zndx.Node.IpType`      // 使用IP类型(ipv4/ipv6)
	NodeCoapPort       string // CoAP端口号
	InRouteMappingPort string // 节点在路由器上的映射端口
	RouterIP           string `zndx.Node.RouteIp` // 节点连接的路由器IP

	// zndx.DeviceModelControlCmd
	CmdID            int64
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

// 查询设备命令
func Query_DevicePowerCmdView_ByDeviceID(device_id, cmdcode string, pDBMap *gorp.DbMap) (list []DevicePowerCmdView, err error) {
	sql := `
SELECT TDev.Id AS DeviceID,
	TDev.Name AS DeviceName,
	TDev.ModelId AS DeviceModelID,
	TDev.PowerNodeId AS PowerNodeID,
	TDev.PowerSwitchId AS PowerSwitchIndex,
	TDev.JoinMethod,
	TDev.JoinNodeId AS JoinNodeID,
	TDev.JoinSocketId AS JoinSocketID,
	TCmd.id AS CmdID,
	TCmd.CmdCode,
	TCmd.CmdName,
	TCmd.RequestType,
	TCmd.RequestURI,
	TCmd.URIQuery,
	TCmd.Payload,
	TCmd.CmdDescription,
	TCmd.DelayMillisecond,
	TCmd.CloseCmdFlag,
	TCmd.OpenCmdFlag,
	TNode.Id AS NodeID,
	TNode.Name AS NodeName,
	TNode.ModelId AS NodeModelID,
	TNode.ClassRoomId AS NodeClassroomID,
	TNode.IpType AS IPType,
	TNode.RouteIp AS RouterIP,
	TNode.NodeCoapPort,
	TNode.InRouteMappingPort	
FROM device TDev
	JOIN devicemodel TDevModel ON(TDevModel.Id=TDev.ModelId)
	JOIN devicemodelcontrolcmd TCmd ON(TCmd.ModelId=TDev.ModelId)
	JOIN node TNode ON(TNode.Id = TDev.JoinNodeId AND TDev.JoinMethod = 'node')
WHERE TDev.Id = ? AND TCmd.CmdCode = ?
`
	list = []DevicePowerCmdView{}
	_, err = pDBMap.Select(&list, sql, device_id, cmdcode)
	return list, err
}
