package zndxview

import (
	gorp "gopkg.in/gorp.v1"
)

// 设备基本信息
type DeviceBasicView struct {
	DeviceID        string
	DeviceName      string
	DeviceRoomName  string
	DeviceModelName string
	DeviceModelDesc string
	PowerNodeID     string
	JoinNodeID      string
}

// 设备详细信息
type DeviceDetailView struct {
	DevID            string // 设备ID
	DevName          string // 设备名称
	PowerNodeID      string // 电源节点ID
	PowerSwitchIndex string // 电源开关序号
	JoinMethod       string // 接入方式
	JoinNodeID       string // 接入节点ID
	JoinSocketIndex  string // 接入插口序号
	DeviceSN         string // 设备序号
	DeviceCode       string // 设备代码
	DeviceBrand      string // 设备品牌
	RoomID           string // 位置
	IsCanUse         string // 能否使用
	UseTimeBefore    int64  // 使用时间(上线前)
	UseTimeAfter     int64  // 使用时间(上线后)

	DeviceModelID string //
	DevModelName  string // 设备型号名称
	ModelDesc     string //
	DevPageFile   string // 设备页面文件
	IsAlert       string //
	MaxUseTime    int64  // 最大使用时间
	DevImageFile  string // 设备图文件(PNG)
	DevImageFile2 string // 设备图文件(GIF)
}

// 设备详细信息（兼容获取设备接口）
type DeviceDetailView0 struct {
	Id                         string
	Name                       string
	Sn                         string
	Code                       string
	Brand                      string
	ModelId                    string
	ModelName                  string
	ClassroomId                int
	Campusname                 string
	Buildingname               string
	Classroomsname             string
	PowerNodeId                string
	PowerSwitchId              string
	Buildingid                 int
	Campusid                   int
	Floorsid                   int
	Floorname                  string
	JoinMethod                 string
	JoinNodeId                 string
	JoinSocketId               string
	NodeSwitchStatus           string
	NodeSwitchStatusUpdateTime string
	DeviceSelfStatus           string
	DeviceSelfStatusUpdateTime string
	IsCanUse                   string
	UseTimeBefore              int32
	UseTimeAfter               int32
	JoinNodeUpdateTime         string
}

// 设备扩展(节点)信息
type DeviceExNodeView struct {
	DeviceID                   string
	DeviceModelID              string
	PowerNodeID                string
	PowerSwitchIndex           string
	JoinNodeID                 string
	JoinSocketIndex            string
	JoinMethod                 string
	NodeSwitchStatus           string
	NodeSwitchStatusUpdateTime string
	DeviceSelfStatus           string
	DeviceSelfStatusUpdateTime string
	IsCanUse                   string
	UseTimeBefore              int64
	UseTimeAfter               int64
	IsAlert                    bool
	MaxUseTime                 int64
}

// 查询设备详细信息列表（兼容获取设备接口）
func QueryList_DeviceDetailView0_ByKeyword(keyword string, pDBMap *gorp.DbMap) (list []DeviceDetailView0, err error) {
	sql := `
SELECT TDev.Id, TDev.Name, TDev.Sn, TDev.Code, TDev.Brand, TDev.ModelId, TDev.ClassroomId, TDev.PowerNodeId,
	ifnull(TDev.JoinMethod, '') JoinMethod,
	ifnull(TDev.JoinNodeId, '') JoinNodeId,
	ifnull(TDev.JoinSocketId, '') JoinSocketId,
	ifnull(TDev.NodeSwitchStatus, '') NodeSwitchStatus,
	ifnull(TDev.NodeSwitchStatusUpdateTime, '') NodeSwitchStatusUpdateTime,
	ifnull(TDev.DeviceSelfStatus, '') DeviceSelfStatus,
	ifnull(TDev.DeviceSelfStatusUpdateTime, '') DeviceSelfStatusUpdateTime,
	ifnull(TDev.IsCanUse, '') IsCanUse,
	ifnull(TDev.UseTimeBefore, 0) UseTimeBefore,
	ifnull(TDev.UseTimeAfter, 0) UseTimeAfter,
	ifnull(TDev.JoinNodeUpdateTime, '') JoinNodeUpdateTime,
	ifnull(TDev.PowerSwitchId, '') PowerSwitchId,
	TModel.Name ModelName,	
	TRoom.Classroomsname,
	TFloor.Floorname,
	TBuilding.Buildingname,
	TCampus.Campusname,		
	TFloor.Buildingid,	
	TBuilding.Campusid,
	TRoom.Floorsid
FROM Device TDev
 JOIN devicemodel TModel ON TDev.ModelId = TModel.Id
 JOIN classrooms TRoom ON TDev.ClassroomId = TRoom.Id
 JOIN floors TFloor ON TFloor.Id = TRoom.Floorsid
 JOIN building TBuilding ON TBuilding.Id = TFloor.Buildingid
 JOIN campus TCampus ON TCampus.Id = TBuilding.Campusid
WHERE CONCAT(TDev.Name, TDev.Sn, TDev.Brand, TDev.Code, TModel.Name, TRoom.Classroomsname, TFloor.Floorname, TBuilding.Buildingname) LIKE ?;
`
	keyword = `%` + keyword + `%`
	_, err = pDBMap.Select(&list, sql, keyword)
	return list, err
}

// 查询【有效】设备详细信息列表（已绑定到节点的有效设备）
func QueryList_ValidDeviceDetailView0_ByKeyword(keyword string, pDBMap *gorp.DbMap) (list []DeviceDetailView0, err error) {
	sql := `
SELECT TDev.Id, TDev.Name, TDev.Sn, TDev.Code, TDev.Brand, TDev.ModelId, TDev.ClassroomId, TDev.PowerNodeId,
	ifnull(TDev.JoinMethod, '') JoinMethod,
	ifnull(TDev.JoinNodeId, '') JoinNodeId,
	ifnull(TDev.JoinSocketId, '') JoinSocketId,
	ifnull(TDev.NodeSwitchStatus, '') NodeSwitchStatus,
	ifnull(TDev.NodeSwitchStatusUpdateTime, '') NodeSwitchStatusUpdateTime,
	ifnull(TDev.DeviceSelfStatus, '') DeviceSelfStatus,
	ifnull(TDev.DeviceSelfStatusUpdateTime, '') DeviceSelfStatusUpdateTime,
	ifnull(TDev.IsCanUse, '') IsCanUse,
	ifnull(TDev.UseTimeBefore, 0) UseTimeBefore,
	ifnull(TDev.UseTimeAfter, 0) UseTimeAfter,
	ifnull(TDev.JoinNodeUpdateTime, '') JoinNodeUpdateTime,
	ifnull(TDev.PowerSwitchId, '') PowerSwitchId,
	TModel.Name ModelName,	
	TRoom.Classroomsname,
	TFloor.Floorname,
	TBuilding.Buildingname,
	TCampus.Campusname,		
	TFloor.Buildingid,	
	TBuilding.Campusid,
	TRoom.Floorsid
FROM Device TDev
 JOIN devicemodel TModel ON TDev.ModelId = TModel.Id
 JOIN classrooms TRoom ON TDev.ClassroomId = TRoom.Id
 JOIN floors TFloor ON TFloor.Id = TRoom.Floorsid
 JOIN building TBuilding ON TBuilding.Id = TFloor.Buildingid
 JOIN campus TCampus ON TCampus.Id = TBuilding.Campusid
WHERE (TDev.PowerNodeId IN(SELECT Id FROM node) OR TDev.JoinNodeId IN(SELECT Id FROM node))
		AND CONCAT(TDev.Name, TDev.Sn, TDev.Brand, TDev.Code, TModel.Name, TRoom.Classroomsname, TFloor.Floorname, TBuilding.Buildingname) LIKE ?
`
	list = []DeviceDetailView0{}
	keyword = `%` + keyword + `%`
	_, err = pDBMap.Select(&list, sql, keyword)
	return list, err
}

// 查询设备详细信息列表（兼容获取设备接口）
func QueryList_DeviceDetailView0_ById(DeviceId string, pDBMap *gorp.DbMap) (view DeviceDetailView0, err error) {
	sql := `
SELECT TDev.Id, TDev.Name, TDev.Sn, TDev.Code, TDev.Brand, TDev.ModelId, TDev.ClassroomId, TDev.PowerNodeId,
	ifnull(TDev.JoinMethod, '') JoinMethod,
	ifnull(TDev.JoinNodeId, '') JoinNodeId,
	ifnull(TDev.JoinSocketId, '') JoinSocketId,
	ifnull(TDev.NodeSwitchStatus, '') NodeSwitchStatus,
	ifnull(TDev.NodeSwitchStatusUpdateTime, '') NodeSwitchStatusUpdateTime,
	ifnull(TDev.DeviceSelfStatus, '') DeviceSelfStatus,
	ifnull(TDev.DeviceSelfStatusUpdateTime, '') DeviceSelfStatusUpdateTime,
	ifnull(TDev.IsCanUse, '') IsCanUse,
	ifnull(TDev.UseTimeBefore, 0) UseTimeBefore,
	ifnull(TDev.UseTimeAfter, 0) UseTimeAfter,
	ifnull(TDev.JoinNodeUpdateTime, '') JoinNodeUpdateTime,
	ifnull(TDev.PowerSwitchId, '') PowerSwitchId,
	TModel.Name ModelName,	
	TRoom.Classroomsname,
	TFloor.Floorname,
	TBuilding.Buildingname,
	TCampus.Campusname,		
	TFloor.Buildingid,	
	TBuilding.Campusid,
	TRoom.Floorsid
FROM Device TDev
 JOIN devicemodel TModel ON TDev.ModelId = TModel.Id
 JOIN classrooms TRoom ON TDev.ClassroomId = TRoom.Id
 JOIN floors TFloor ON TFloor.Id = TRoom.Floorsid
 JOIN building TBuilding ON TBuilding.Id = TFloor.Buildingid
 JOIN campus TCampus ON TCampus.Id = TBuilding.Campusid
WHERE TDev.Id=?;
`
	err = pDBMap.SelectOne(&view, sql, DeviceId)
	return view, err
}

// 查询设备扩展相关信息(节点)
func Query_DeviceExNodeView(device_id string, pDBMap *gorp.DbMap) (view DeviceExNodeView, err error) {
	var sql = `
SELECT TDev.Id DeviceID, 
	TDev.ModelID AS DeviceModelID,
	IFNULL(TDev.PowerNodeId, '') PowerNodeID,
	IFNULL(TDev.PowerSwitchId, '') PowerSwitchIndex,
	IFNULL(TDev.JoinNodeId, '') JoinNodeID,
	IFNULL(TDev.JoinSocketId, '') JoinSocketIndex,
	IFNULL(TDev.JoinMethod, '') JoinMethod,
	IFNULL(TDev.NodeSwitchStatus, '') NodeSwitchStatus,
	IFNULL(TDev.NodeSwitchStatusUpdateTime, '') NodeSwitchStatusUpdateTime,
	IFNULL(TDev.DeviceSelfStatus, '') DeviceSelfStatus,
	IFNULL(TDev.DeviceSelfStatusUpdateTime, '') DeviceSelfStatusUpdateTime,
	IFNULL(TDev.IsCanUse, '1') IsCanUse,
	IFNULL(TDev.UseTimeBefore, 0) UseTimeBefore,
	IFNULL(TDev.UseTimeAfter, 0) UseTimeAfter,
	IFNULL(TModel.IsAlert, FALSE) IsAlert,
	IFNULL(TModel.MaxUseTime, 0) MaxUseTime
FROM Device TDev JOIN DeviceModel TModel ON TDev.ModelId = TModel.Id
WHERE TDev.Id =? 
`
	err = pDBMap.SelectOne(&view, sql, device_id)
	return view, err
}

// 查询设备基本信息
func QueryList_DeviceBasicView(keyword string, pDBMap *gorp.DbMap) (list []DeviceBasicView, err error) {
	sql := `
SELECT 	TDevModel.Name AS DeviceModelName, 
				TDevModel.Description AS DeviceModelDesc,
				TDevRoom.Classroomsname AS DeviceRoomName,
				TDev.Id AS DeviceID, 
				TDev.Name AS DeviceName,
				TDev.PowerNodeId AS PowerNodeID,
				TDev.JoinNodeId AS JoinNodeID
FROM device AS TDev 
	JOIN node TNode ON(TNode.Id IN(TDev.PowerNodeId,TDev.JoinNodeId))
	JOIN devicemodel AS TDevModel ON (TDevModel.Id=TDev.ModelId)
	JOIN classrooms AS TDevRoom ON (TDevRoom.Id=TDev.ClassroomId)	
WHERE TDev.Id LIKE :Keyword
		OR TDev.Name LIKE :Keyword
		OR TDev.Brand LIKE :Keyword
		OR TDevModel.Name LIKE :Keyword
		OR TDevModel.Description LIKE :Keyword
		OR TDevRoom.Classroomsname LIKE :Keyword
`
	params := map[string]interface{}{}
	params["Keyword"] = "%" + keyword + "%"
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}

// 查询设备详细信息
func QueryList_DeviceDetailView(nodeID string, pDBMap *gorp.DbMap) (list []DeviceDetailView, err error) {
	sql := `SELECT 	TDev.ID AS DevID, 
				TDev.Name AS DevName,
				TDev.PowerNodeId AS PowerNodeID, 
				TDev.PowerSwitchId AS PowerSwitchIndex,
				TDev.JoinMethod AS JoinMethod, 
				TDev.JoinNodeId AS JoinNodeID, 
				TDev.JoinSocketId AS JoinSocketIndex,
				TDev.Sn AS DeviceSN,
				TDev.Code AS DeviceCode,
				TDev.Brand AS DeviceBrand,
				TDev.ClassroomId AS RoomID,
				TDev.IsCanUse,
				TDev.UseTimeBefore,
				TDev.UseTimeAfter,				
				TModel.Id AS DeviceModelID,
				TModel.Name AS DevModelName,
				TModel.Description AS ModelDesc,
				TModel.IsAlert,
				TModel.MaxUseTime,
				TModel.PageFileName AS DevPageFile,
				TModel.ImgFileName AS DevImageFile,
				TModel.ImgFileName2 AS DevImageFile2	
FROM device TDev
		JOIN devicemodel TModel ON(TModel.Id=TDev.ModelId)
WHERE  :NodeID IN(TDev.PowerNodeId, TDev.JoinNodeId)
`
	params := map[string]interface{}{}
	params["NodeID"] = nodeID
	_, err = pDBMap.Select(&list, sql, params)
	return list, err
}
