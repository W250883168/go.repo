package deviceDataAccess

import (
	"strconv"
	"strings"

	"gopkg.in/gorp.v1"

	"dborm/zndx"
	"dborm/zndxview"
	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xtext"

	"dev.project/BackEndCode/devserver/model/core"
	devmodel "dev.project/BackEndCode/devserver/model/deviceModel"
)

type r struct {
	NodeCmds   interface{}
	DeviceCmds interface{}
}

func Query_DeviceCmd_ByRoom(roomID int, cmd string, dbmap *gorp.DbMap) (ret interface{}, err error) {
	defer xerr.CatchPanic()
	ret = &r{DeviceCmds: []zndxview.DeviceCmd{}, NodeCmds: []zndxview.NodeCmd{}}

	sql := `
SELECT TModel.Name ModelName, TNode.Name NodeName, TNode.Id NodeID, TNode.IpType, TNode.InRouteMappingPort, TNode.NodeCoapPort, TNode.RouteIp,
			TCmd.CmdCode, TCmd.CmdName, TCmd.RequestURI, TCmd.URIQuery, TCmd.RequestType, TCmd.Payload
FROM nodemodelcmd TCmd
	JOIN nodemodel TModel ON(TModel.Id = TCmd.ModelId)
	JOIN node TNode ON(TNode.ModelId = TModel.Id)	
WHERE TNode.ClassRoomId = :RoomID AND TCmd.CmdCode = :CmdCode;
`
	var node_cmds = []zndxview.NodeCmd{}

	maparg := map[string]interface{}{}
	maparg["RoomID"] = roomID
	maparg["CmdCode"] = cmd
	_, err = dbmap.Select(&node_cmds, sql, maparg)
	xerr.ThrowPanic(err)

	sql2 := `
SELECT TModel.Name ModelName, TDev.Name DeviceName,  TDev.PowerNodeId, TDev.PowerSwitchId, TDev.JoinMethod, TDev.JoinNodeId, TDev.JoinSocketId, 
			 TCmd.CmdCode, TCmd.CmdName, TCmd.RequestURI, TCmd.URIQuery, TCmd.RequestType, TCmd.Payload,
			 TNode.IpType, TNode.InRouteMappingPort, TNode.NodeCoapPort, TNode.RouteIp
FROM devicemodelcontrolcmd TCmd
	JOIN devicemodel TModel ON(TModel.Id = TCmd.ModelId)
	JOIN device TDev ON(TDev.ModelId = TModel.Id)
	JOIN node TNode ON(TNode.Id = TDev.JoinNodeId)
WHERE TNode.ClassRoomId = :RoomID AND TCmd.CmdCode = :CmdCode;
`
	var dev_cmds = []zndxview.DeviceCmd{}
	maparg2 := map[string]interface{}{}
	maparg2["RoomID"] = roomID
	maparg2["CmdCode"] = cmd
	_, err = dbmap.Select(&dev_cmds, sql2, maparg2)
	xerr.ThrowPanic(err)

	ret = &r{DeviceCmds: dev_cmds, NodeCmds: node_cmds}
	return ret, err
}

//获取所有设备故障信息
func GetAllFaultInfoList2(requestData devmodel.RequestData, siteType string, siteId string, modelId string, keyWord string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var sqlargs = map[string]interface{}{}
	var where = ""
	var err error
	page := xhttp.PageInfo{PageIndex: requestData.Page.PageIndex, PageSize: requestData.Page.PageSize}

	//后面的记录统计和查询都要的SQL(虚拟成一个表)
	subtable := `
SELECT TFault.Id FaultId, TFault.DeviceId, TFault.HappenTime, TFault.FaultSummary, TFault.Status, TFault.InputUserId,
		CASE TFault.IsCanUse WHEN '0' THEN '不可使用' WHEN '1' THEN '可以使用' END IsCanUse,
		CASE TFault.Status WHEN '0' THEN '草稿' WHEN '1' THEN '待受理' WHEN '2' THEN '维修中' WHEN '3' THEN '已维修' END StatusName,		
		IFNULL(TUser.TrueName, '') InputUserName,
		IFNULL(TDev.Code, '') DeviceCode,
		IFNULL(TDev.ModelId, '') ModelId,
		IFNULL(TDev.Name, '') DeviceName,
		CONCAT_WS('', TCampus.Campusname, TBuild.Buildingname, TFloor.Floorname, TRoom.Classroomsname) DeviceSite,
		IFNULL(TModel.Name, '') DeviceModel
FROM DeviceFault TFault
		LEFT JOIN device TDev ON TDev.id = TFault.DeviceId
		LEFT JOIN devicemodel TModel ON TModel.id = TDev.ModelId
		LEFT JOIN classrooms TRoom ON TRoom.Id = TDev.ClassroomId
		LEFT JOIN floors TFloor ON TFloor.Id = TRoom.Floorsid
		LEFT JOIN building TBuild ON TBuild.Id = TFloor.Buildingid
		LEFT JOIN campus TCampus ON TCampus.Id = TBuild.Campusid
		LEFT JOIN Users TUser ON TUser.id = TFault.InputUserId
`
	where += `	WHERE (1 = 1) AND ((Status = '0' AND InputUserId = :UserID) OR (Status != '0'))		`
	uid := strconv.Itoa(requestData.Auth.Usersid) //得到当前用户
	sqlargs["UserID"] = uid

	//拼接查询条件(where)
	//1）拼接安装位置
	if xtext.IsNotBlank(siteType) && xtext.IsNotBlank(siteId) {
		switch siteType {
		case "campus":
			where += " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=:SiteID))))	"
		case "building":
			where += " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =:SiteID)))	"
		case "floor":
			where += " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid=:SiteID))	"
		case "classroom":
			where += " and DeviceId in (select id from device where ClassroomId=:SiteID)	"
		}

		sqlargs["SiteID"] = siteId
	}
	//2)拼接设备型号
	if xtext.IsNotBlank(modelId) {
		where += " AND FIND_IN_SET(ModelId,getDeviceModelChildNodes( :ModelID ))>0	" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
		sqlargs["ModelID"] = modelId
	}

	//3)拼接关键字
	if xtext.IsNotBlank(keyWord) {
		where += `	AND (HappenTime LIKE :Keyword 
						OR DeviceModel LIKE :Keyword 
						OR DeviceSite LIKE :Keyword
						OR DeviceCode LIKE :Keyword 
						OR DeviceName LIKE :Keyword 
						OR FaultSummary LIKE :Keyword 
						OR StatusName LIKE :Keyword 
						OR IsCanUse LIKE :Keyword 
						OR InputUserName LIKE :Keyword)
			 `
		sqlargs["Keyword"] = "%" + keyWord + "%"
	}

	rd.Rcode = gResponseMsgs.FAIL.CodeText()
	rd.Reason = "数据读取错误"

	//计算分页信息

	sql = "select count(*) from (" + subtable + ") aa" + where
	count, err := dbmap.SelectInt(sql, sqlargs)
	page.RowTotal = int(count)
	if err != nil {
		xdebug.LogError(err)
		return rd
	}

	//获得具体数据
	sql = `
SELECT FaultId, DeviceId,
		IFNULL(DeviceCode, '') DeviceCode,
		IFNULL(DeviceName, '') DeviceName,
		IFNULL(DeviceModel, '') DeviceModel,
		IFNULL(DeviceSite, '') DeviceSite,
		HappenTime,
		FaultSummary,
		IsCanUse,
		Status,
		StatusName,
		InputUserId,
		InputUserName 
FROM ( @SubTable ) aa  @SqlWhere  ORDER BY HappenTime DESC ` + page.SQL_LimitString()
	sql = strings.Replace(sql, "@SubTable", subtable, 1)
	sql = strings.Replace(sql, "@SqlWhere", where, 1)

	var list []zndx.DeviceFaultView
	if _, err = dbmap.Select(&list, sql, sqlargs); err != nil {
		xdebug.LogError(err)
		return rd
	}

	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &devmodel.ResultData{devmodel.PageData{
		PageIndex:   page.PageIndex,
		PageSize:    page.PageSize,
		RecordCount: page.RowTotal,
		PageCount:   page.PageTotal()}, list}
	return rd
}
