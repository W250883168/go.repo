package dal

import (
	"gopkg.in/gorp.v1"

	"dborm/zndx"
)

// 节点开关命令
type NodeSwitchCmdView struct {
	// zndx.Node
	NodeID              string // zndx.Node.Id
	NodeName            string // 节点名称
	NodeModelID         string // 节点型号ID
	ClassRoomID         int    // 节点所在教室ID
	IPType              string // 使用IP类型(ipv4/ipv6)
	RouterIP            string // 节点连接的路由器IP
	NodeCoapPort        string // CoAP端口号
	InRouterMappingPort string // 节点在路由器上的映射端口
	UploadTime          string // 上报时间

	// zndx.NodeModelCmd
	CmdID          int64 // zndx.NodeModelCmd.Id
	CmdCode        string
	CmdName        string
	RequestURI     string
	URIQuery       string
	RequestType    string
	Payload        string
	CloseCmdFlag   string
	OpenCmdFlag    string
	CmdDescription string
}

// 查询节点开关命令(ByDevice)
func Query_NodeSwitchCmd_ByDevice(pDevice *zndx.Device, cmdCode string, pDBMap *gorp.DbMap) (view NodeSwitchCmdView, err error) {
	sql := `
SELECT 	TNode.Id AS NodeID, 
		TNode.Name AS NodeName,
		TNode.ModelId AS NodeModelID,
		TNode.ClassRoomId AS ClassRoomID,
		TNode.IpType AS IPType,
		TNode.NodeCoapPort,
		TNode.InRouteMappingPort AS InRouterMappingPort,
		TNode.RouteIp AS RouterIP,
		TNode.UploadTime,
		TCmd.Id AS CmdID,
		TCmd.CmdCode,
		TCmd.CmdName,
		TCmd.RequestURI,
		TCmd.RequestType,
		TCmd.URIQuery,
		IFNULL(TCmd.CmdDescription, '') AS CmdDescription,
		TCmd.Payload,
		TCmd.CloseCmdFlag,
		TCmd.OpenCmdFlag
FROM node TNode 
	JOIN nodemodel TModel ON (TModel.Id = TNode.ModelId)
	JOIN nodemodelcmd TCmd ON (TCmd.ModelId = TModel.Id)
WHERE TNode.Id IN(SELECT PowerNodeId FROM device WHERE Id = :DeviceID)
	AND TCmd.CmdCode = :CmdCode		
`
	maparg := map[string]interface{}{}
	maparg["DeviceID"] = pDevice.Id
	maparg["CmdCode"] = cmdCode
	err = pDBMap.SelectOne(&view, sql, maparg)
	return view, err
}

// 查询节点开关命令(ByNode)
func QueryList_NodeSwitchCmd_ByNode(pNode *zndx.Node, cmd string, pDBMap *gorp.DbMap) (list []NodeSwitchCmdView, err error) {
	sql := `
SELECT 	TNode.Id AS NodeID, 
				TNode.Name AS NodeName,
				TNode.ModelId AS NodeModelID,
				TNode.ClassRoomId AS ClassRoomID,
				TNode.IpType AS IPType,
				TNode.NodeCoapPort,
				TNode.InRouteMappingPort AS InRouterMappingPort,
				TNode.RouteIp AS RouterIP,
				TNode.UploadTime,
				TCmd.Id AS CmdID,
				TCmd.CmdCode,
				TCmd.CmdName,
				TCmd.RequestURI,
				TCmd.RequestType,
				TCmd.URIQuery,
				IFNULL(TCmd.CmdDescription, '') AS CmdDescription,
				TCmd.Payload,
				TCmd.CloseCmdFlag,
				TCmd.OpenCmdFlag
FROM node AS TNode 
		JOIN nodemodel AS TModel ON (TModel.Id = TNode.ModelId)
		JOIN nodemodelcmd AS TCmd ON (TCmd.ModelId = TModel.Id)
WHERE TNode.Id = :NodeID AND TCmd.CmdCode = :CmdCode;
`
	list = []NodeSwitchCmdView{}
	params := map[string]interface{}{}
	params["NodeID"] = pNode.Id
	params["CmdCode"] = cmd
	_, err = pDBMap.Select(&list, sql, params)

	return list, err
}

// 查询节点开关命令(ByRoom)
func Query_NodeSwitchCmd_ByRoom(cmd, room_id string, pDBMap *gorp.DbMap) (list []NodeSwitchCmdView, err error) {
	sql := `
SELECT 	TNode.Id AS NodeID, 
				TNode.Name AS NodeName,
				TNode.ModelId AS NodeModelID,
				TNode.ClassRoomId AS ClassRoomID,
				TNode.IpType AS IPType,
				TNode.NodeCoapPort,
				TNode.InRouteMappingPort AS InRouterMappingPort,
				TNode.RouteIp AS RouterIP,
				TNode.UploadTime,
				TCmd.Id AS CmdID,
				TCmd.CmdCode,
				TCmd.CmdName,
				TCmd.RequestURI,
				TCmd.RequestType,
				TCmd.URIQuery,
				IFNULL(TCmd.CmdDescription, '') AS CmdDescription,
				TCmd.Payload,
				TCmd.CloseCmdFlag,
				TCmd.OpenCmdFlag
FROM node AS TNode 
		JOIN nodemodel AS TModel ON (TModel.Id = TNode.ModelId)
		JOIN nodemodelcmd AS TCmd ON (TCmd.ModelId = TModel.Id)
WHERE TCmd.CmdCode = :CmdCode 
		AND TNode.Id IN(SELECT device.PowerNodeId FROM device WHERE ClassRoomId= :ClassroomID)
`
	params := map[string]interface{}{}
	params["CmdCode"] = cmd
	params["ClassroomID"] = room_id
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}

// 查询节点开关命令(ByRoom)
func Query_NodeSwitchCmd_ByFloor(cmd, floor_id, room_id string, pDBMap *gorp.DbMap) (list []NodeSwitchCmdView, err error) {
	sql := `
SELECT 	TNode.Id AS NodeID, 
				TNode.Name AS NodeName,
				TNode.ModelId AS NodeModelID,
				TNode.ClassRoomId AS ClassRoomID,
				TNode.IpType AS IPType,
				TNode.NodeCoapPort,
				TNode.InRouteMappingPort AS InRouterMappingPort,
				TNode.RouteIp AS RouterIP,
				TNode.UploadTime,
				TCmd.Id AS CmdID,
				TCmd.CmdCode,
				TCmd.CmdName,
				TCmd.RequestURI,
				TCmd.RequestType,
				TCmd.URIQuery,
				IFNULL(TCmd.CmdDescription, '') AS CmdDescription,
				TCmd.Payload,
				TCmd.CloseCmdFlag,
				TCmd.OpenCmdFlag
FROM node AS TNode 
		JOIN nodemodel AS TModel ON (TModel.Id = TNode.ModelId)
		JOIN nodemodelcmd AS TCmd ON (TCmd.ModelId = TModel.Id)
WHERE TCmd.CmdCode = :CmdCode 
		AND TNode.Id IN(SELECT TDev.PowerNodeId 
										FROM device TDev 
											JOIN classrooms TRoom ON(TRoom.Id = TDev.ClassroomId)
										WHERE ClassRoomId=:ClassroomID AND TRoom.Floorsid=:FloorID)
`
	params := map[string]interface{}{}
	params["CmdCode"] = cmd
	params["ClassroomID"] = room_id
	params["FloorID"] = floor_id
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}
