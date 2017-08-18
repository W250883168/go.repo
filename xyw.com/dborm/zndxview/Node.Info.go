package zndxview

import (
	gorp "gopkg.in/gorp.v1"
)

// 节点基本信息
type NodeBasicView struct {
	NodeID        string
	NodeName      string
	NodeRoomName  string
	NodeModelName string
	NodeModelDesc string
}

// 节点或设备基本信息
type NodeDevBasicView struct {
	ID        string
	Name      string
	Room      string
	Model     string // 型号
	ModelDesc string // 型号描述
	What      string // node/device
}

// 查询节点设备基本信息
func QueryList_NodeDevBasicView(keyword string, pDBMap *gorp.DbMap) (list []NodeDevBasicView, err error) {
	sql := `
SELECT A.DeviceID ID, A.DeviceName Name, A.DeviceModelName Model, A.DeviceRoomName Room, A.DeviceModelDesc ModelDesc, A.What
FROM (SELECT TDevModel.Name AS DeviceModelName, TDevModel.Description AS DeviceModelDesc,
					TDevRoom.Classroomsname AS DeviceRoomName,
					TDev.Id AS DeviceID, TDev.Name AS DeviceName,
					'device' AS What 
			FROM device AS TDev 
				JOIN devicemodel AS TDevModel ON (TDevModel.Id=TDev.ModelId)
				JOIN classrooms AS TDevRoom ON (TDevRoom.Id=TDev.ClassroomId)	
			WHERE TDev.Id LIKE :Keyword
					OR TDev.Name LIKE :Keyword
					OR TDev.Brand LIKE :Keyword
					OR TDevModel.Name LIKE :Keyword
					OR TDevModel.Description LIKE :Keyword
					OR TDevRoom.Classroomsname LIKE :Keyword) AS A
UNION
SELECT B.NodeID ID, B.NodeName Name, B.NodeModelName Model, B.NodeRoomName Room, B.NodeModelDesc ModelDesc, B.What
FROM (SELECT TNodeModel.Name AS NodeModelName, TNodeModel.Description AS NodeModelDesc,
					TDevRoom.Classroomsname AS NodeRoomName,
					TNode.Id AS NodeID, TNode.Name AS NodeName,
					'node' AS What
			FROM node AS TNode
				JOIN nodemodel AS TNodeModel ON(TNodeModel.Id=TNode.ModelId)	
				JOIN classrooms AS TDevRoom ON (TDevRoom.Id=TNode.ClassroomId)	
			WHERE TNode.Id LIKE :Keyword
					OR TNode.Name LIKE :Keyword
					OR TNodeModel.Name LIKE :Keyword
					OR TNodeModel.Description LIKE :Keyword
					OR TDevRoom.Classroomsname LIKE :Keyword) B
`
	params := map[string]interface{}{}
	params["Keyword"] = "%" + keyword + "%"
	_, err = pDBMap.Select(&list, sql, params)

	return list, err
}

// 查询节点基本信息
func QueryList_NodeBasicView(keyword string, pDBMap *gorp.DbMap) (list []NodeBasicView, err error) {
	sql := `
SELECT TNodeModel.Name AS NodeModelName, TNodeModel.Description AS NodeModelDesc,
		TDevRoom.Classroomsname AS NodeRoomName,
		TNode.Id AS NodeID, TNode.Name AS NodeName
FROM node AS TNode
	JOIN nodemodel AS TNodeModel ON(TNodeModel.Id=TNode.ModelId)	
	JOIN classrooms AS TDevRoom ON (TDevRoom.Id=TNode.ClassroomId)	
WHERE TNode.Id LIKE :Keyword
		OR TNode.Name LIKE :Keyword	
		OR TNodeModel.Name LIKE :Keyword
		OR TNodeModel.Description LIKE :Keyword
		OR TDevRoom.Classroomsname LIKE :Keyword
`
	params := map[string]interface{}{}
	params["Keyword"] = "%" + keyword + "%"
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}
