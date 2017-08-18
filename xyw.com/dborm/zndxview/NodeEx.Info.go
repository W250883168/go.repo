package zndxview

import (
	gorp "gopkg.in/gorp.v1"

	"xutils/xhttp"
)

// 节点扩展信息
type NodeDetailView struct {
	Id                 string
	Name               string
	ModelId            string
	Campusid           int
	Buildingid         int
	Floorsid           int
	ClassRoomId        int
	IpType             string
	NodeCoapPort       string
	InRouteMappingPort string
	RouteIp            string
	UploadTime         string
	NodeModelName      string
	Classroomsname     string
	Buildingname       string
	Campusname         string
	FloorName          string
}

// 查询节点扩展信息
func QueryList_NodeDetailView_ByNodeID(node_id string, pDBMap *gorp.DbMap) (list []NodeDetailView, err error) {
	sql := `
SELECT TNode.Id,
	TNode.Name,
	TNode.ModelId,
	TNode.ClassRoomId,
	TNode.IpType,
	TNode.NodeCoapPort,
	TNode.InRouteMappingPort,
	TNode.RouteIp,
	TNode.UploadTime,
	IFNULL(TModel.Name, '') NodeModelName,
	IFNULL(TRoom.Classroomsname, '') Classroomsname,
	IFNULL(TBuilding.Buildingname, '') Buildingname,
	IFNULL(TCampus.Campusname, '') Campusname,
	IFNULL(TFloor.Floorname,'') FloorName
FROM node TNode
	JOIN nodemodel TModel ON (TNode.ModelId = TModel.Id)
	JOIN Classrooms TRoom ON (TNode.ClassRoomId = TRoom.Id)
	JOIN floors TFloor ON (TRoom.Floorsid = TFloor.Id)
	JOIN building TBuilding ON (TFloor.Buildingid = TBuilding.Id)
	JOIN campus TCampus ON (TBuilding.Campusid = TCampus.Id)
WHERE 	TNode.Id =?
`
	list = []NodeDetailView{}
	sqlargs := []interface{}{node_id}
	_, err = pDBMap.Select(&list, sql, sqlargs...)

	return list, err
}

// 查询节点详细信息(ByKeyword)
func Query_NodeDetailView_ByKeyword(keyword string, pPage *xhttp.PageInfo, pDBMap *gorp.DbMap) (list []NodeDetailView, err error) {
	sql := `
SELECT TNode.Id,
	TNode.Name,
	TNode.ModelId,
	TNode.ClassRoomId,
	TNode.IpType,
	TNode.NodeCoapPort,
	TNode.InRouteMappingPort,
	TNode.RouteIp,
	TNode.UploadTime,
	IFNULL(TModel.Name, '') NodeModelName,
	IFNULL(TRoom.Classroomsname, '') Classroomsname,
	IFNULL(TBuilding.Buildingname, '') Buildingname,
	IFNULL(TCampus.Campusname, '') Campusname,
	IFNULL(TFloor.Floorname,'') FloorName
FROM node TNode
	JOIN nodemodel TModel ON (TNode.ModelId = TModel.Id)
	JOIN Classrooms TRoom ON (TNode.ClassRoomId = TRoom.Id)
	JOIN floors TFloor ON (TRoom.Floorsid = TFloor.Id)
	JOIN building TBuilding ON (TFloor.Buildingid = TBuilding.Id)
	JOIN campus TCampus ON (TBuilding.Campusid = TCampus.Id)
WHERE CONCAT_WS('',TNode.Id, TNode.Name, TModel.Name, TRoom.Classroomsname, TFloor.Floorscode, TBuilding.Buildingname, TCampus.Campusname) LIKE ?
`
	sql += pPage.SQL_LimitString()
	keyword = "%" + keyword + "%"
	_, err = pDBMap.Select(&list, sql, keyword)
	return list, err
}

// 查询节点详细信息(ByID)
func Query_NodeDetailView_ByID(NodeID string, pPage *xhttp.PageInfo, pDBMap *gorp.DbMap) (list NodeDetailView, err error) {
	sql := `
SELECT TNode.Id,
	TNode.Name,
	TNode.ModelId,
	TNode.ClassRoomId,
	TNode.IpType,
	TNode.NodeCoapPort,
	TNode.InRouteMappingPort,
	TNode.RouteIp,
	TNode.UploadTime,
	IFNULL(TModel.Name, '') NodeModelName,
	IFNULL(TRoom.Classroomsname, '') Classroomsname,
	IFNULL(TBuilding.Buildingname, '') Buildingname,
	IFNULL(TCampus.Campusname, '') Campusname,
	IFNULL(TFloor.Floorname,'') FloorName
FROM node TNode
	JOIN nodemodel TModel ON (TNode.ModelId = TModel.Id)
	JOIN Classrooms TRoom ON (TNode.ClassRoomId = TRoom.Id)
	JOIN floors TFloor ON (TRoom.Floorsid = TFloor.Id)
	JOIN building TBuilding ON (TFloor.Buildingid = TBuilding.Id)
	JOIN campus TCampus ON (TBuilding.Campusid = TCampus.Id)
WHERE TNode.Id=?
`
	err = pDBMap.SelectOne(&list, sql, NodeID)
	return list, err
}
