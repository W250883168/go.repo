package deviceDataAccess

import (
	dbsql "database/sql"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/gorp.v1"
	"gopkg.in/ini.v1"

	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xhttp"
	"xutils/xtext"

	"dborm/zndx"

	"dev.project/BackEndCode/devserver/commons"
	"dev.project/BackEndCode/devserver/model/core"
	model "dev.project/BackEndCode/devserver/model/deviceModel"
)

var (
	OfflineTime   = "30" //多长时间节点没有上传数据就算该节点为离线
	gResponseMsgs = commons.ResponseMsgSet_Instance()
)

//获得分页信息
func GetPageInfo(requestData model.RequestData, sql string, dbmap *gorp.DbMap) (pg model.PageData, rerr error) {
	//如果客户端未传入分页数据，或者传入分页数据中的当前页小于等于0，都当成不分页处理
	if requestData.Page.PageIndex <= 0 {
		return pg, nil
	}

	//校正分页尺寸
	pg.PageSize = requestData.Page.PageSize
	if pg.PageSize <= 0 {
		pg.PageSize = 10 //默认10条
	}

	//统计记录总数
	count, err := dbmap.SelectInt(sql)
	if err != nil {
		return pg, err
	}
	pg.RecordCount = int(count)

	//统计总页数
	if pg.RecordCount%pg.PageSize > 0 {
		pg.PageCount = pg.RecordCount/pg.PageSize + 1
	} else {
		pg.PageCount = pg.RecordCount / pg.PageSize
	}

	//校正当前页
	pg.PageIndex = requestData.Page.PageIndex
	if pg.PageIndex > pg.PageCount {
		pg.PageIndex = pg.PageCount
	}

	return pg, nil
}

//获得分页信息
func GetPageInfoEx(requestData model.RequestData, dbmap *gorp.DbMap, sql string, args ...interface{}) (pg model.PageData, rerr error) {
	//如果客户端未传入分页数据，或者传入分页数据中的当前页小于等于0，都当成不分页处理
	if requestData.Page.PageIndex <= 0 {
		return pg, nil
	}

	//校正分页尺寸
	pg.PageSize = requestData.Page.PageSize
	if pg.PageSize <= 0 {
		pg.PageSize = 10 //默认10条
	}

	//统计记录总数
	count, err := dbmap.SelectInt(sql, args...)
	if err != nil {
		return pg, err
	}
	pg.RecordCount = int(count)

	//统计总页数
	if pg.RecordCount%pg.PageSize > 0 {
		pg.PageCount = pg.RecordCount/pg.PageSize + 1
	} else {
		pg.PageCount = pg.RecordCount / pg.PageSize
	}

	//校正当前页
	pg.PageIndex = requestData.Page.PageIndex
	if pg.PageIndex > pg.PageCount {
		pg.PageIndex = pg.PageCount
	}

	return pg, nil
}

//获得sql limit字符串
func GetLimitString(pg model.PageData) string {
	if pg.PageIndex <= 0 {
		return ""
	}

	return " limit " + strconv.Itoa((pg.PageIndex-1)*pg.PageSize) + "," + strconv.Itoa(pg.PageSize)
}

//获取设备使用日志（开关机）
func GetDeviceUseLogList(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//拼接查询条件(where)
	sqlWhere := " where DeviceId = '" + deviceId + "'"

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from DeviceUseLog" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		xdebug.LogError(err)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select OnTime,IFNULL(OffTime,'') OffTime,IFNULL(SEC_TO_TIME(UseTime),'') UseTime from DeviceUseLog " + sqlWhere + " order by OnTime Desc " + GetLimitString(pg)
	var list []model.DeviceUseLog
	_, err = dbmap.Select(&list, sql)
	xdebug.LogError(err)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

// 获取设备操作日志
func GetDeviceOperateLogList(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//拼接查询条件(where)
	sqlWhere := " where DeviceId = '" + deviceId + "'"

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from DeviceDetailLog" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select d.OperateTime,IFNULL(u.Loginuser,'') UserCode,IFNULL(u.TrueName,'') UserName,IFNULL(d.CmdName,'') CmdName,IFNULL(d.Para,'') Para from DeviceDetailLog d left join Users u on d.OperateUserId=u.Id " + sqlWhere + " order by d.OperateTime Desc " + GetLimitString(pg)
	var list []model.DeviceDetailLog
	_, err = dbmap.Select(&list, sql)
	xdebug.LogError(err)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

// 获取设备操作日志
func GetDeviceOperateLogList2(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	responses := commons.ResponseMsgSet_Instance()
	rd.Rcode = responses.FAIL.CodeText()
	rd.Reason = responses.FAIL.Text

	//
	devlog := zndx.DeviceDetailLog{DeviceId: deviceId}
	pInfo := &xhttp.PageInfo{PageIndex: requestData.Page.PageIndex, PageSize: requestData.Page.PageSize}
	list, err := devlog.GetByDeviceID(pInfo, dbmap)
	xdebug.LogError(err)
	if err == nil {
		rd.Rcode = responses.SUCCESS.CodeText()
		rd.Reason = responses.SUCCESS.Text
		rd.Result = &model.ResultData{model.PageData{
			PageIndex:   pInfo.PageIndex,
			PageSize:    pInfo.PageSize,
			RecordCount: pInfo.RowTotal,
			PageCount:   pInfo.PageTotal(),
		}, list}
	}

	return rd
}

//获取设备预警信息
func GetDeviceAlertInfoList(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//拼接查询条件(where)
	sqlWhere := " where DeviceId = '" + deviceId + "'"

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from DeviceAlert" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select LastAlertTime AlertTime,AlertDescription from DeviceAlert  " + sqlWhere + " order by LastAlertTime Desc " + GetLimitString(pg)
	var list []model.DeviceAlertInfo
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//获取设备故障信息
func GetDeviceFaultInfoList(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//拼接查询条件(where)
	sqlWhere := " where DeviceId = '" + deviceId + "'"

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from DeviceFault" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select d.Id, d.HappenTime,IFNULL(d.FaultSummary,'') FaultSummary,IFNULL(d.FaultDescription,'') FaultDescription,IFNULL(d.InputUserId,0) InputUserId,IFNULL(u.TrueName,'') InputUserName,IFNULL(d.InputTime,'') InputTime,IFNULL(d.IsCanUse,'') IsCanUse,Status from DeviceFault d left join Users u on d.InputUserId=u.Id " + sqlWhere + " order by d.HappenTime Desc " + GetLimitString(pg)
	sql = `
SELECT d.Id, d.HappenTime, Status,
	IFNULL(d.FaultSummary, '') FaultSummary,
	IFNULL(d.FaultDescription, '') FaultDescription,
	IFNULL(d.InputUserId, 0) InputUserId,
	IFNULL(u.TrueName, '') InputUserName,
	IFNULL(d.InputTime, '') InputTime,
	IFNULL(d.IsCanUse, '') IsCanUse
FROM DeviceFault d
	LEFT JOIN Users u on d.InputUserId=u.Id
WHERE DeviceId = ?
ORDER BY d.HappenTime DESC
`

	var list []model.DeviceFaultInfo
	_, err = dbmap.Select(&list, sql, deviceId)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}
	return rd
}

//获取教室状态信息
func GetClassroomStatusList(requestData model.RequestData, floorIds string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//获得具体数据
	//说明：在查询设备状态时，用0表示关，1表示开，9999表示离线的方式表示节点或设备的状态，是为了方便在外层进行计算
	sql = "select"
	sql = sql + " b.id BuildingId,b.Buildingname BuildingName,"
	sql = sql + " f.Id FloorId,f.Floorname FloorName,IFNULL(f.FloorsImage,'') FloorImage, "
	sql = sql + " c.id ClassroomId,c.ClassroomsName ClassroomName,c.Classroomstate ClassroomState,c.Collectionnumbers CollectionNumbers,"
	sql = sql + " IFNULL(aaa.HaveStop,-1) HaveStop,IFNULL(aaa.HaveAlert,-1) HaveAlert,IFNULL(aaa.HaveOffline,-1) HaveOffline,IFNULL(aaa.HaveRun,-1) HaveRun"
	sql = sql + " from classrooms c "
	sql = sql + " left join floors f on f.id = c.floorsId"
	sql = sql + "      left join building b on b.id = f.Buildingid"
	sql = sql + " 	   left join ("
	sql = sql + " 				select  ClassroomId,"
	sql = sql + " 					    case when sum(StopFlag) >0 then 1 else 0 end HaveStop,"   //只要有一个设备被停止，则表示教室有设备停止
	sql = sql + " 						case when sum(AlertFlag) >0 then 1 else 0 end HaveAlert,"    //只要有一个设备有预警，则表示教室有设备预警
	sql = sql + " 						case when sum(OfflineFlag)>0 then 1 else 0 end HaveOffline," //只要有一个设备离线，则表示教室有设备离线
	sql = sql + " 						case when sum(RunFlag)>0 then 1 else 0 end HaveRun"          //只要有一个设备运行，则表示教室有设备运行
	sql = sql + " 				from (   "
	sql = sql + " 						select d.ClassroomId,"
	sql = sql + " 							   d.Id DeviceId,"
	sql = sql + " 							   case when IFNULL(IsCanUse,'1')='1' then 0 else 1 end StopFlag,"
	sql = sql + " 							   case when EXISTS (select 1 from DeviceAlert  where DeviceId=d.Id) then '1' else '0' end AlertFlag,"
	//	sql = sql + " 							   case when IFNULL(PowerNodeId,'')='' or IFNULL(PowerSwitchId,'')='' then 0" //0表示在线
	//	sql = sql + " 									when IFNULL(NodeSwitchStatusUpdateTime,'')=''  then 9999"                   //离线
	//	sql = sql + " 									when TIMESTAMPDIFF(SECOND,NodeSwitchStatusUpdateTime,now())<" + OfflineTime + " then  0"
	//	sql = sql + " 									else 9999 "                                       //离线
	//	sql = sql + " 								    end NodeSwitch1,"                              //提供电源给设备使用的节点的状态
	//	sql = sql + " 							   case when IFNULL(JoinNodeId,'')='' then 0"       //0表示在线
	//	sql = sql + " 									when IFNULL(JoinNodeUpdateTime,'')=''  then 9999" //离线
	//	sql = sql + " 									when TIMESTAMPDIFF(SECOND,JoinNodeUpdateTime,now())<" + OfflineTime + " then  0"
	//	sql = sql + " 									else 9999 "          //离线
	//	sql = sql + " 								    end NodeSwitch2," //提供插槽给设备使用的节点的状态
	sql = sql + " 		                       case when IFNULL(PowerNodeId,'')!='' and IFNULL(NodeSwitchStatusUpdateTime,'')=''  then 1"                                                                               //1-离线
	sql = sql + " 		                            when IFNULL(JoinNodeId,'')!='' and IFNULL(JoinNodeUpdateTime,'')=''  then 1"                                                                                        //1-离线
	sql = sql + "  			                        when IFNULL(PowerNodeId,'')!='' and IFNULL(NodeSwitchStatusUpdateTime,'')!='' and  TIMESTAMPDIFF(SECOND,NodeSwitchStatusUpdateTime,now())>" + OfflineTime + " then 1" //1-离线
	sql = sql + "  			                        when IFNULL(JoinNodeId,'')!='' and IFNULL(JoinNodeUpdateTime,'')!='' and  TIMESTAMPDIFF(SECOND,JoinNodeUpdateTime,now())>" + OfflineTime + " then 1"                  //1-离线
	sql = sql + "                                   else 0"
	sql = sql + "                                   end OfflineFlag,"
	sql = sql + " 							   case when DeviceSelfStatus='on' THEN 1"
	sql = sql + " 									when DeviceSelfStatus='off' THEN 0"
	sql = sql + " 									else 0"
	sql = sql + " 								    end RunFlag  " //form-input-select
	sql = sql + " 						from 	Device d  "
	sql = sql + " 						where d.ClassroomId in (select id from classrooms where Floorsid in (select id from floors where Buildingid in (select id from building where id in (" + floorIds + "))))"
	sql = sql + " 					    ) aa"
	sql = sql + " 				group by ClassroomId"
	sql = sql + "               ) aaa on aaa.ClassroomId = c.Id"

	sql = sql + " where c.Floorsid in (select id from floors where Buildingid in (select id from building where id in (" + floorIds + ")))"
	sql = sql + " order by b.id,f.id,c.id"

	sql = `
SELECT  TRoom.id ClassroomId, TRoom.ClassroomsName ClassroomName, TRoom.Classroomstate ClassroomState, TRoom.Collectionnumbers CollectionNumbers,
  IFNULL(TBuildiing.id, 0) BuildingId,
	IFNULL(TBuildiing.Buildingname, '') BuildingName,
	IFNULL(TFloor.Id, 0) FloorId,
	IFNULL(TFloor.Floorname, '') FloorName,
	IFNULL(TFloor.FloorsImage, '') FloorImage,	
	IFNULL(aaa.HaveStop, -1) HaveStop,
	IFNULL(aaa.HaveAlert, -1) HaveAlert,
	IFNULL(aaa.HaveOffline, -1) HaveOffline,
	IFNULL(aaa.HaveRun, -1) HaveRun
FROM classrooms TRoom
	LEFT JOIN floors TFloor ON TFloor.id = TRoom.floorsId
	LEFT JOIN building TBuildiing ON TBuildiing.id = TFloor.Buildingid
	LEFT JOIN (SELECT ClassroomId, 
										CASE WHEN SUM(StopFlag) > 0 THEN 1 ELSE 0 END HaveStop,
										CASE WHEN SUM(AlertFlag) > 0 THEN 1 ELSE 0 END HaveAlert,
										CASE WHEN SUM(OfflineFlag) > 0 THEN 1 ELSE 0 END HaveOffline,
										CASE WHEN SUM(RunFlag) > 0 THEN 1 ELSE 0 END HaveRun
						 FROM (SELECT TDev.ClassroomId, TDev.Id DeviceId,
													CASE WHEN IFNULL(IsCanUse, '1') = '1' THEN 0 WHEN LENGTH(IsCanUse) = 0 THEN 0 ELSE 1 END StopFlag,
													CASE WHEN EXISTS (SELECT 1 FROM DeviceAlert WHERE DeviceId = TDev.Id) THEN '1' ELSE '0' END AlertFlag,
													CASE 	WHEN IFNULL(PowerNodeId, '') != '' AND IFNULL(NodeSwitchStatusUpdateTime, '') = '' THEN 1 
																WHEN IFNULL(JoinNodeId, '') != '' AND IFNULL(JoinNodeUpdateTime, '') = '' THEN 1 
																WHEN IFNULL(PowerNodeId, '') != '' AND IFNULL(NodeSwitchStatusUpdateTime, '') != '' AND TIMESTAMPDIFF(SECOND, NodeSwitchStatusUpdateTime, now()) > :OfflineTime THEN 1
																WHEN IFNULL(JoinNodeId, '') != '' AND IFNULL(JoinNodeUpdateTime, '') != '' AND TIMESTAMPDIFF(SECOND, JoinNodeUpdateTime, now()) > :OfflineTime THEN 1
																ELSE 0
													END OfflineFlag,
													CASE WHEN DeviceSelfStatus = 'on' THEN 1 WHEN DeviceSelfStatus = 'off' THEN 0 ELSE 0 END RunFlag
										FROM Device TDev
										WHERE TDev.ClassroomId IN (SELECT id FROM classrooms WHERE Floorsid IN (SELECT id FROM floors WHERE Buildingid IN (SELECT id FROM building WHERE id IN (@FloorIDs) )))) aa
										GROUP BY ClassroomId ) aaa ON aaa.ClassroomId = TRoom.Id
WHERE TRoom.Floorsid IN (SELECT id FROM floors WHERE Buildingid IN (SELECT id FROM building WHERE id IN (@FloorIDs)))
ORDER BY TBuildiing.Buildingname, TFloor.Floorname, TRoom.Classroomsname
`
	sql = strings.Replace(sql, "@FloorIDs", floorIds, 2)
	var list []model.ClassroomStatusData
	maparg := map[string]interface{}{}
	maparg["OfflineTime"] = OfflineTime
	_, err = dbmap.Select(&list, sql, maparg)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//获取教室状态信息[外部插件所用]
func GetClassroomStatusListPuls(floorIds string, dbmap *gorp.DbMap) (list []model.ClassroomStatusData) {
	var sql string
	//获得具体数据
	//说明：在查询设备状态时，用0表示关，1表示开，9999表示离线的方式表示节点或设备的状态，是为了方便在外层进行计算
	//	sql = "select"
	//	sql = sql + " b.id BuildingId,b.Buildingname BuildingName,"
	//	sql = sql + " f.Id FloorId,f.Floorname FloorName,IFNULL(f.FloorsImage,'') FloorImage, "
	//	sql = sql + " c.id ClassroomId,c.ClassroomsName ClassroomName,c.Classroomstate ClassroomState,c.Collectionnumbers CollectionNumbers,"
	//	sql = sql + " IFNULL(aaa.HaveStop,-1) HaveStop,IFNULL(aaa.HaveAlert,-1) HaveAlert,IFNULL(aaa.HaveOffline,-1) HaveOffline,IFNULL(aaa.HaveRun,-1) HaveRun"
	//	sql = sql + " from classrooms c "
	//	sql = sql + " left join floors f on f.id = c.floorsId"
	//	sql = sql + "      left join building b on b.id = f.Buildingid"
	//	sql = sql + " 	   left join ("
	//	sql = sql + " 				select  ClassroomId,"
	//	sql = sql + " 					    case when sum(StopFlag) >0 then 1 else 0 end HaveStop,"   //只要有一个设备被停止，则表示教室有设备停止
	//	sql = sql + " 						case when sum(AlertFlag) >0 then 1 else 0 end HaveAlert,"    //只要有一个设备有预警，则表示教室有设备预警
	//	sql = sql + " 						case when sum(OfflineFlag)>0 then 1 else 0 end HaveOffline," //只要有一个设备离线，则表示教室有设备离线
	//	sql = sql + " 						case when sum(RunFlag)>0 then 1 else 0 end HaveRun"          //只要有一个设备运行，则表示教室有设备运行
	//	sql = sql + " 				from (   "
	//	sql = sql + " 						select d.ClassroomId,"
	//	sql = sql + " 							   d.Id DeviceId,"
	//	sql = sql + " 							   case when IFNULL(IsCanUse,'1')='1' then 0 else 1 end StopFlag,"
	//	sql = sql + " 							   case when EXISTS (select 1 from DeviceAlert  where DeviceId=d.Id) then '1' else '0' end AlertFlag,"
	//	sql = sql + " 		                       case when IFNULL(PowerNodeId,'')!='' and IFNULL(NodeSwitchStatusUpdateTime,'')=''  then 1"                                                                               //1-离线
	//	sql = sql + " 		                            when IFNULL(JoinNodeId,'')!='' and IFNULL(JoinNodeUpdateTime,'')=''  then 1"                                                                                        //1-离线
	//	sql = sql + "  			                        when IFNULL(PowerNodeId,'')!='' and IFNULL(NodeSwitchStatusUpdateTime,'')!='' and  TIMESTAMPDIFF(SECOND,NodeSwitchStatusUpdateTime,now())>" + OfflineTime + " then 1" //1-离线
	//	sql = sql + "  			                        when IFNULL(JoinNodeId,'')!='' and IFNULL(JoinNodeUpdateTime,'')!='' and  TIMESTAMPDIFF(SECOND,JoinNodeUpdateTime,now())>" + OfflineTime + " then 1"                  //1-离线
	//	sql = sql + "                                   else 0"
	//	sql = sql + "                                   end OfflineFlag,"
	//	sql = sql + " 							   case when DeviceSelfStatus='on' THEN 1"
	//	sql = sql + " 									when DeviceSelfStatus='off' THEN 0"
	//	sql = sql + " 									else 0"
	//	sql = sql + " 								    end RunFlag  " //form-input-select
	//	sql = sql + " 						from 	Device d  "
	//	sql = sql + " 						where d.ClassroomId in (select id from classrooms where Floorsid in (select id from floors where Buildingid in (select id from building where id in (" + floorIds + "))))"
	//	sql = sql + " 					    ) aa"
	//	sql = sql + " 				group by ClassroomId"
	//	sql = sql + "               ) aaa on aaa.ClassroomId = c.Id"

	//	sql = sql + " where c.Floorsid in (select id from floors where Buildingid in (select id from building where id in (" + floorIds + ")))"
	//	sql = sql + " order by b.id,f.id,c.id"

	sql = `
SELECT TBuildiing.id BuildingId,
	TBuildiing.Buildingname BuildingName,
	TFloor.Id FloorId,
	TFloor.Floorname FloorName,
	IFNULL(TFloor.FloorsImage, '') FloorImage,
	TRoom.id ClassroomId,
	TRoom.ClassroomsName ClassroomName,
	TRoom.Classroomstate ClassroomState,
	TRoom.Collectionnumbers CollectionNumbers,
	IFNULL(aaa.HaveStop, -1) HaveStop,
	IFNULL(aaa.HaveAlert, -1) HaveAlert,
	IFNULL(aaa.HaveOffline, -1) HaveOffline,
	IFNULL(aaa.HaveRun, -1) HaveRun,ChangeTime
FROM classrooms TRoom
	LEFT JOIN floors TFloor ON TFloor.id = TRoom.floorsId
	LEFT JOIN building TBuildiing ON TBuildiing.id = TFloor.Buildingid
	LEFT JOIN (SELECT ClassroomId,ChangeTime, 
										CASE WHEN SUM(StopFlag) > 0 THEN 1 ELSE 0 END HaveStop,
										CASE WHEN SUM(AlertFlag) > 0 THEN 1 ELSE 0 END HaveAlert,
										CASE WHEN SUM(OfflineFlag) > 0 THEN 1 ELSE 0 END HaveOffline,
										CASE WHEN SUM(RunFlag) > 0 THEN 1 ELSE 0 END HaveRun
						 FROM (SELECT TDev.ClassroomId, TDev.Id DeviceId,ifnull(max(DeviceSelfStatusUpdateTime),'') ChangeTime,
													CASE WHEN IFNULL(IsCanUse, '1') = '1' THEN 0 WHEN LENGTH(IsCanUse) = 0 THEN 0 ELSE 1 END StopFlag,
													CASE WHEN EXISTS (SELECT 1 FROM DeviceAlert WHERE DeviceId = TDev.Id) THEN '1' ELSE '0' END AlertFlag,
													CASE 	WHEN IFNULL(PowerNodeId, '') != '' AND IFNULL(NodeSwitchStatusUpdateTime, '') = '' THEN 1 
																WHEN IFNULL(JoinNodeId, '') != '' AND IFNULL(JoinNodeUpdateTime, '') = '' THEN 1 
																WHEN IFNULL(PowerNodeId, '') != '' AND IFNULL(NodeSwitchStatusUpdateTime, '') != '' AND TIMESTAMPDIFF(SECOND, NodeSwitchStatusUpdateTime, now()) > :OfflineTime THEN 1
																WHEN IFNULL(JoinNodeId, '') != '' AND IFNULL(JoinNodeUpdateTime, '') != '' AND TIMESTAMPDIFF(SECOND, JoinNodeUpdateTime, now()) > :OfflineTime THEN 1
																ELSE 0
													END OfflineFlag,
													CASE WHEN DeviceSelfStatus = 'on' THEN 1 WHEN DeviceSelfStatus = 'off' THEN 0 ELSE 0 END RunFlag
										FROM Device TDev
										WHERE TDev.ClassroomId IN (SELECT id FROM classrooms WHERE Floorsid IN (SELECT id FROM floors WHERE Buildingid IN (SELECT id FROM building WHERE id IN (@FloorIDs) )))) aa
										GROUP BY ClassroomId,ChangeTime ) aaa ON aaa.ClassroomId = TRoom.Id
WHERE TRoom.Floorsid IN (SELECT id FROM floors WHERE Buildingid IN (SELECT id FROM building WHERE id IN (@FloorIDs)))
ORDER BY TBuildiing.Buildingname, TFloor.Floorname, TRoom.Classroomsname
`
	sql = strings.Replace(sql, "@FloorIDs", floorIds, 2)
	//	var list []model.ClassroomStatusData
	maparg := map[string]interface{}{}
	maparg["OfflineTime"] = OfflineTime
	_, err := dbmap.Select(&list, sql, maparg)
	if err != nil {
		log.Println("err:", err)
	}
	return list
}

//获取所有设备操作日志
func GetAllOperateLogList(requestData model.RequestData, fromTime string, toTime string, keyWord string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	// 后面的记录统计和查询都要的SQL(虚拟成一个表)
	// sqlTable := "select d.Id,d.OperateTime,IFNULL(u.Loginuser,'') UserCode,IFNULL(u.TrueName,'') UserName,IFNULL(d.CmdName,'') OperateName,case d.OperateObject when 'device' then '设备' when 'classroom' then '教室' when 'floor' then '楼层' end OperateObject,IFNULL(d.ObjectName,'') ObjectName from DeviceOperateLog d left join Users u on d.OperateUserId=u.Id"
	sqlTable := `
SELECT TLog.Id, TLog.OperateTime,
	IFNULL(TUser.Loginuser, '') UserCode,
	IFNULL(TUser.TrueName, '') UserName,
	IFNULL(TLog.CmdName, '') OperateName,
	IFNULL(TLog.ObjectName, '') ObjectName,
	CASE TLog.OperateObject 	WHEN 'device' THEN '设备'
							WHEN 'classroom' THEN '教室'
							WHEN 'floor' THEN '楼层'
	END OperateObject,
	IFNULL(TModel.ImgFileName, '') ImgFileName
FROM DeviceOperateLog TLog
	LEFT JOIN device TDev ON(TDev.Id = TLog.ObjectId)
	LEFT JOIN devicemodel TModel ON(TModel.Id = TDev.ModelId)	
	LEFT JOIN Users TUser ON (TLog.OperateUserId = TUser.Id)
`
	sqlTable = `
SELECT TLog.Id, TLog.OperateTime,
	IFNULL(TUser.Loginuser, '') UserCode,
	IFNULL(TUser.TrueName, '') UserName,
	IFNULL(TLog.CmdName, '') OperateName,
	IFNULL(TLog.ObjectName, '') ObjectName,
	CASE TLog.OperateObject 	WHEN 'device' THEN '设备'
							WHEN 'classroom' THEN '教室'
							WHEN 'floor' THEN '楼层'
	END OperateObject 
FROM DeviceOperateLog TLog
	LEFT JOIN Users TUser ON (TLog.OperateUserId = TUser.Id)
`
	// 拼接查询条件(where)
	maparg := map[string]interface{}{}
	sqlWhere := " WHERE 1=1"
	if xtext.IsNotBlank(fromTime) {
		sqlWhere += ` AND OperateTime >= :FromTime `
		maparg["FromTime"] = fromTime
	}
	if xtext.IsNotBlank(toTime) {
		sqlWhere += ` and OperateTime <= :ToTime `
		maparg["ToTime"] = toTime
	}
	if xtext.IsNotBlank(keyWord) {
		sqlWhere += ` AND CONCAT_WS('', OperateTime, UserCode, Username, OperateName, OperateObject, ObjectName) LIKE :Keyword `
		maparg["Keyword"] = "%" + keyWord + "%"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfoEx(requestData, dbmap, sql, maparg)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + " order by OperateTime Desc " + GetLimitString(pg)
	var list []model.DeviceOperateLog
	_, err = dbmap.Select(&list, sql, maparg)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}
	return rd
}

//获取所有设备警告信息
func GetAllAlertInfoList(requestData model.RequestData, siteType string, siteId string, modelId string, keyWord string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL(虚拟成一个表)
	sqlTable := ""
	sqlTable = sqlTable + " select a.Id AlertId,a.DeviceId,a.LastAlertTime AlertTime,a.AlertDescription,IFNULL(d.`Code`,'') DeviceCode,d.ModelId,d.`Name` DeviceName,"
	sqlTable = sqlTable + "        CONCAT(c.Campusname,b.Buildingname,f.Floorname,r.Classroomsname) DeviceSite,m.`Name` DeviceModel"
	sqlTable = sqlTable + " from DeviceAlert a "
	sqlTable = sqlTable + "      left join device d on d.id = a.DeviceId"
	sqlTable = sqlTable + "      left join devicemodel m on m.id = d.ModelId"
	sqlTable = sqlTable + "	     left join classrooms r on r.Id = d.ClassroomId"
	sqlTable = sqlTable + "	     left join floors f on f.Id = r.Floorsid"
	sqlTable = sqlTable + "	     left join building b on b.Id = f.Buildingid"
	sqlTable = sqlTable + "	     left join campus c on c.Id=b.Campusid"

	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	//1）拼接安装位置
	if siteType != "" && siteId != "" {
		switch siteType {
		case "campus":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + "))))"
		case "building":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + ")))"
		case "floor":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid=" + siteId + "))"
		case "classroom":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId=" + siteId + " )"
		}

	}
	//2)拼接设备型号
	if modelId != "" {
		sqlWhere = sqlWhere + " and FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "'))>0" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
	}
	//3)拼接关键字
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "AlertTime" + s + "or DeviceModel" + s + "or DeviceSite" + s + "or DeviceCode" + s + "or DeviceName" + s + "or AlertDescription" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select AlertId,DeviceId,DeviceCode,DeviceName,DeviceModel,DeviceSite,AlertTime,AlertDescription from (" + sqlTable + ") aa " + sqlWhere + " order by AlertTime Desc " + GetLimitString(pg)
	var list []model.DeviceAllAlertInfo
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//获取所有设备警告信息
func GetAllAlertInfoList2(requestData model.RequestData, siteType string, siteId string, modelId string, keyWord string, dbmap *gorp.DbMap) (rd core.Returndata) {
	defer xerr.CatchPanic()
	rd.Rcode = gResponseMsgs.EXE_CMD_FAIL.CodeText()
	rd.Reason = gResponseMsgs.EXE_CMD_FAIL.Text

	//后面的记录统计和查询都要的SQL(虚拟成一个表)
	sqlTable := `SELECT a.Id AlertId, a.DeviceId, a.LastAlertTime AlertTime, a.AlertDescription, 
						 IFNULL(d.Code, '') DeviceCode, d.ModelId, d.Name DeviceName,
						 CONCAT(c.Campusname,b.Buildingname,f.Floorname,r.Classroomsname) DeviceSite,
						 m.Name DeviceModel
				FROM DeviceAlert a
					JOIN device d ON d.id = a.DeviceId
					JOIN devicemodel m ON m.id = d.ModelId
					JOIN classrooms r ON r.Id = d.ClassroomId
					JOIN floors f ON f.Id = r.Floorsid
					JOIN building b ON b.Id = f.Buildingid
					JOIN campus c ON c.Id = b.Campusid	
					`
	sqlWhere := ` where 1=1	`

	//拼接查询条件(where)
	//1）拼接安装位置
	if siteType != "" && siteId != "" {
		switch siteType {
		case "campus":
			sqlWhere += " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + "))))"
		case "building":
			sqlWhere += " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + ")))"
		case "floor":
			sqlWhere += " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid=" + siteId + "))"
		case "classroom":
			sqlWhere += " and DeviceId in (select id from device where ClassroomId=" + siteId + " )"
		}
	}

	//2)拼接设备型号
	if modelId != "" {
		sqlWhere += " and FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "'))>0" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
	}

	//3)拼接关键字
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "AlertTime" + s + "or DeviceModel" + s + "or DeviceSite" + s + "or DeviceCode" + s + "or DeviceName" + s + "or AlertDescription" + s
		sqlWhere += " and (" + s + ")"
	}

	//计算分页信息
	sql := "select count(*) from (" + sqlTable + ") aa" + sqlWhere
	pg, _ := GetPageInfo(requestData, sql, dbmap)
	count, err := dbmap.SelectInt(sql)
	xerr.ThrowPanic(err)

	//获得具体数据
	sql = "select AlertId,DeviceId,DeviceCode,DeviceName,DeviceModel,DeviceSite,AlertTime,AlertDescription from (" + sqlTable + ") aa " + sqlWhere + " order by AlertTime Desc " + GetLimitString(pg)
	var list []model.DeviceAllAlertInfo
	_, err = dbmap.Select(&list, sql)
	xerr.ThrowPanic(err)

	pg.RecordCount = int(count)
	rd.Rcode = gResponseMsgs.SUCCESS.CodeText()
	rd.Reason = gResponseMsgs.SUCCESS.Text
	rd.Result = &model.ResultData{pg, list}
	return rd
}

//获取所有设备故障信息
func GetAllFaultInfoList(requestData model.RequestData, siteType string, siteId string, modelId string, keyWord string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//得到当前用户（后面过滤使用）
	uid := strconv.Itoa(requestData.Auth.Usersid)

	//后面的记录统计和查询都要的SQL(虚拟成一个表)
	sqlTable := ""
	sqlTable = sqlTable + " select a.Id FaultId,a.DeviceId,a.HappenTime,a.FaultSummary,case a.IsCanUse when '0' then '不可使用' when '1' then '可以使用' end IsCanUse,a.Status,case a.Status when '0' then '草稿' when '1' then '待受理' when '2' then '维修中' when '3' then '已维修' end StatusName,a.InputUserId,u.TrueName InputUserName,IFNULL(d.`Code`,'') DeviceCode,d.ModelId,d.`Name` DeviceName,"
	sqlTable = sqlTable + "        CONCAT(c.Campusname,b.Buildingname,f.Floorname,r.Classroomsname) DeviceSite,m.`Name` DeviceModel"
	sqlTable = sqlTable + " from DeviceFault a "
	sqlTable = sqlTable + "      left join device      d on d.id = a.DeviceId"
	sqlTable = sqlTable + "      left join devicemodel m on m.id = d.ModelId"
	sqlTable = sqlTable + "	     left join classrooms  r on r.Id = d.ClassroomId"
	sqlTable = sqlTable + "	     left join floors      f on f.Id = r.Floorsid"
	sqlTable = sqlTable + "	     left join building    b on b.Id = f.Buildingid"
	sqlTable = sqlTable + "	     left join campus      c on c.Id=b.Campusid"
	sqlTable = sqlTable + "      left join Users       u on u.id = a.InputUserId"

	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	//0) 过滤用户(非草稿状态的全部记录，草稿状态的只查当前用记录）
	sqlWhere = sqlWhere + " and ((Status='0' and InputUserId='" + uid + "') or (Status!='0')) "

	//1）拼接安装位置
	if siteType != "" && siteId != "" {
		switch siteType {
		case "campus":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + "))))"
		case "building":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + ")))"
		case "floor":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid=" + siteId + "))"
		case "classroom":
			sqlWhere = sqlWhere + " and DeviceId in (select id from device where ClassroomId=" + siteId + " )"
		}

	}
	//2)拼接设备型号
	if modelId != "" {
		sqlWhere = sqlWhere + " and FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "'))>0" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
	}
	//3)拼接关键字
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "HappenTime" + s + "or DeviceModel" + s + "or DeviceSite" + s + "or DeviceCode" + s + "or DeviceName" + s + "or FaultSummary" + s + "or StatusName" + s + "or IsCanUse" + s + "or InputUserName" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		xdebug.LogError(err)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	// sql = "select FaultId,DeviceId,DeviceCode,DeviceName,DeviceModel,DeviceSite,HappenTime,FaultSummary,IsCanUse,Status,StatusName,InputUserId,InputUserName from (" + sqlTable + ") aa " + sqlWhere + " order by HappenTime Desc " + GetLimitString(pg)
	sql = `SELECT
				FaultId,
				DeviceId,
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
			FROM (` + sqlTable + ") aa " + sqlWhere + " order by HappenTime Desc " + GetLimitString(pg)
	var list []model.DeviceAllFaultInfo
	_, err = dbmap.Select(&list, sql)
	xdebug.LogError(err)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误11:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//获取设备数量
func GetDeviceQty(requestData model.RequestData, siteType string, siteId string, modelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//获得设备数量---------------------------------------------------------------------------
	//拼接查询条件(where)
	sqlWhere := " where 1=1"

	//1)拼接安装位置
	if siteType != "" && siteId != "" {
		switch siteType {
		case "campus":
			sqlWhere = sqlWhere + " and ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + ")))"
		case "building":
			sqlWhere = sqlWhere + " and ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + "))"
		case "floor":
			sqlWhere = sqlWhere + " and ClassroomId in (select id from classrooms where Floorsid=" + siteId + ")"
		case "classroom":
			sqlWhere = sqlWhere + " and ClassroomId=" + siteId
		}

	}
	//2)拼接设备型号
	if modelId != "" {
		sqlWhere = sqlWhere + " and FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "'))>0" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
	}

	//查询
	sql = ""
	sql = sql + " select Id DeviceId,IFNULL(ModelId,'') ModelId,"
	sql = sql + "  		case when IFNULL(IsCanUse,'1')='1' then 0 else 1 end StopFlag,"
	sql = sql + "  		case when EXISTS (select 1 from DeviceAlert  where DeviceId=d.Id) then 1 else 0 end AlertFlag,"
	sql = sql + "       case when EXISTS (select 1 from DeviceFault  where  status!='0' and !(status='3' and RepairResult='1') and DeviceId=d.Id) then 1 else 0 end FaultFlag,"
	sql = sql + " 		case when IFNULL(PowerNodeId,'')!='' and IFNULL(NodeSwitchStatusUpdateTime,'')=''  then 1"                                                                               //1-离线
	sql = sql + " 		     when IFNULL(JoinNodeId,'')!='' and IFNULL(JoinNodeUpdateTime,'')=''  then 1"                                                                                        //1-离线
	sql = sql + "  			 when IFNULL(PowerNodeId,'')!='' and IFNULL(NodeSwitchStatusUpdateTime,'')!='' and  TIMESTAMPDIFF(SECOND,NodeSwitchStatusUpdateTime,now())>" + OfflineTime + " then 1" //1-离线
	sql = sql + "  			 when IFNULL(JoinNodeId,'')!='' and IFNULL(JoinNodeUpdateTime,'')!='' and  TIMESTAMPDIFF(SECOND,JoinNodeUpdateTime,now())>" + OfflineTime + " then 1"                  //1-离线
	sql = sql + "            else 0"
	sql = sql + "  			 end OfflineFlag "
	sql = sql + " from 	Device d "

	var deviceQty []model.DeviceQty
	_, err = dbmap.Select(&deviceQty, sql+sqlWhere)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "获取设备数量错误:" + err.Error()
		return rd
	}

	//获得设备型号---------------------------------------------------------------------------
	switch modelId {
	case "":
		sqlWhere = " where IFNULL(PId,'')=''"
	default:
		sqlWhere = " where PId='" + modelId + "'"
	}

	sql = "select Id ModelId,Name ModelName,getDeviceModelChildNodes(Id) SubModelIds,0 TotalQty,0 StopQty,0 FaultQty,0 OfflineQty,0 AlertQty from DeviceModel "

	var deviceModelQty []model.DeviceModelQty
	_, err = dbmap.Select(&deviceModelQty, sql+sqlWhere)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "获取设备型号数据错误:" + err.Error()
		return rd
	}

	//使用以上两个数据集进行计算
	for i, m := range deviceModelQty {
		t, a, f, o, s := 0, 0, 0, 0, 0
		for _, d := range deviceQty {
			if d.ModelId != "" && strings.Index(m.SubModelIds, d.ModelId) >= 0 {
				t++
				a += d.AlertFlag
				f += d.FaultFlag
				o += d.OfflineFlag
				s += d.StopFlag
			}
		}
		deviceModelQty[i].TotalQty = t
		deviceModelQty[i].AlertQty = a
		deviceModelQty[i].FaultQty = f
		deviceModelQty[i].StopQty = s
		deviceModelQty[i].OfflineQty = o
		deviceModelQty[i].SubModelIds = ""
		m.SubModelIds = ""
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, deviceModelQty}

	return rd
}

//设备分析-按设备型号统计使用时间
func GetUseTimeByModel(requestData model.RequestData, fromTime string, toTime string, siteType string, siteId string, modelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//获得设备数量---------------------------------------------------------------------------
	//拼接查询条件(where)
	sqlWhere := " where 1=1"

	//1)拼接安装位置
	if siteType != "" && siteId != "" {
		switch siteType {
		case "campus":
			sqlWhere = sqlWhere + " and ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + ")))"
		case "building":
			sqlWhere = sqlWhere + " and ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + "))"
		case "floor":
			sqlWhere = sqlWhere + " and ClassroomId in (select id from classrooms where Floorsid=" + siteId + ")"
		case "classroom":
			sqlWhere = sqlWhere + " and ClassroomId=" + siteId
		}

	}
	//2)拼接设备型号
	if modelId != "" {
		sqlWhere = sqlWhere + " and FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "'))>0" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
	}

	//查询
	sql = ""
	sql = sql + " select DeviceId,ModelId,TIMESTAMPDIFF(SECOND,OnTime,OffTime) UseTime"
	sql = sql + " from ("
	sql = sql + " select DeviceId,"
	sql = sql + " case when OnTime>='" + fromTime + "' then OnTime"
	sql = sql + "      else '" + fromTime + "'"
	sql = sql + "      end OnTime,"
	sql = sql + " case when IFNULL(OffTime,'')='' then '" + toTime + "'"
	sql = sql + "      when OffTime<='" + toTime + "' then OffTime"
	sql = sql + "      else '" + toTime + "'"
	sql = sql + "      end OffTime"
	sql = sql + " from deviceuselog"
	sql = sql + " where OffTime>='" + fromTime + "' and OnTime<='" + toTime + "'"
	sql = sql + " ) a"
	sql = sql + " left join device d on d.Id = a.DeviceId"

	var deviceUseTime []model.DeviceUseTime
	_, err = dbmap.Select(&deviceUseTime, sql+sqlWhere)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "获取设备使用时间错误:" + err.Error()
		return rd
	}

	//获得设备型号---------------------------------------------------------------------------
	switch modelId {
	case "":
		sqlWhere = " where IFNULL(PId,'')=''"
	default:
		sqlWhere = " where PId='" + modelId + "'"
	}

	sql = "select Id ModelId,Name ModelName,getDeviceModelChildNodes(Id) SubModelIds,0 UseTime from DeviceModel "

	var deviceModelUseTime []model.DeviceModelUseTime
	_, err = dbmap.Select(&deviceModelUseTime, sql+sqlWhere)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "获取设备型号数据错误:" + err.Error()
		return rd
	}

	//使用以上两个数据集进行计算
	for i, m := range deviceModelUseTime {
		var t int64
		t = 0
		for _, d := range deviceUseTime {
			if d.ModelId != "" && strings.Index(m.SubModelIds, d.ModelId) >= 0 {
				t += d.UseTime
			}
		}
		deviceModelUseTime[i].UseTime = t
		m.SubModelIds = ""
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, deviceModelUseTime}
	return rd
}

//设备分析-按设备型号统计使用时间
func GetUseTimeByModel2(requestData model.RequestData, fromTime string, toTime string, siteType string, siteId string, modelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	sql := `SELECT DeviceId, ModelId, TIMESTAMPDIFF(SECOND, OnTime, OffTime) UseTime
			FROM (SELECT DeviceId, 
						CASE WHEN OnTime >= '@FromTime' THEN OnTime ELSE '@FromTime' END OnTime,
						CASE WHEN IFNULL(OffTime, '') = '' THEN '@ToTime' WHEN OffTime <= '@ToTime' THEN OffTime ELSE '@ToTime' END OffTime
					FROM deviceuselog WHERE OffTime >= '@FromTime' AND OnTime <= '@ToTime') a  JOIN device d ON (d.Id = a.DeviceId)
			WHERE (1=1)
			`
	sql = strings.Replace(sql, "@FromTime", fromTime, -1)
	sql = strings.Replace(sql, "@ToTime", toTime, -1)
	sqlWhere := ""

	//获得设备数量---------------------------------------------------------------------------
	//拼接查询条件(where)
	//1)拼接安装位置
	if siteType != "" && siteId != "" {
		switch siteType {
		case "campus":
			sqlWhere += " and ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + ")))"
		case "building":
			sqlWhere += " and ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + "))"
		case "floor":
			sqlWhere += " and ClassroomId in (select id from classrooms where Floorsid=" + siteId + ")"
		case "classroom":
			sqlWhere += " and ClassroomId=" + siteId
		}
	}

	//2)拼接设备型号
	if modelId != "" {
		sqlWhere += " AND FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "'))>0" //编写一个MYSQL函数，传入一个节点ID,返回该ID及所有子节点ID
	}

	var deviceUseTime []model.DeviceUseTime
	// fmt.Println("<<<<<<<<<<<\n:		", sql+sqlWhere)
	_, err := dbmap.Select(&deviceUseTime, sql+sqlWhere)
	xdebug.LogError(err)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "获取设备使用时间错误:" + err.Error()
		return rd
	}

	//获得设备型号---------------------------------------------------------------------------
	sql = "select Id ModelId,Name ModelName,getDeviceModelChildNodes(Id) SubModelIds,0 UseTime from DeviceModel "
	sqlWhere = " where IFNULL(PId,'')=''"
	if xtext.IsNotBlank(modelId) {
		sqlWhere = " where PId='" + modelId + "'"
	}
	var deviceModelUseTime []model.DeviceModelUseTime
	_, err = dbmap.Select(&deviceModelUseTime, sql+sqlWhere)
	xdebug.LogError(err)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "获取设备型号数据错误:" + err.Error()
		return rd
	}

	//使用以上两个数据集进行计算
	for i, m := range deviceModelUseTime {
		var t int64
		t = 0
		for _, d := range deviceUseTime {
			if d.ModelId != "" && strings.Index(m.SubModelIds, d.ModelId) >= 0 {
				t += d.UseTime
			}
		}
		deviceModelUseTime[i].UseTime = t
		m.SubModelIds = ""
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{Data: deviceModelUseTime}
	return rd
}

//设备分析-按设备位置统计使用时间
func GetUseTimeBySite(requestData model.RequestData, fromTime string, toTime string, siteType string, siteId string, modelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//位置过滤语句
	sqlSite := ""
	switch siteType {
	case "":
		sqlSite = " 1=1"
	case "campus":
		sqlSite = " DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid in ( select Id  from building where Campusid=" + siteId + "))))"
	case "building":
		sqlSite = " DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid in (select Id from floors where Buildingid =" + siteId + ")))"
	case "floor":
		sqlSite = " DeviceId in (select id from device where ClassroomId in (select id from classrooms where Floorsid=" + siteId + "))"
	}

	//型号过滤条件
	sqlModel := " 1=1"
	if modelId != "" {
		sqlModel = " deviceId in (select Id from device where FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + modelId + "')) )"
	}

	//主体查询语句（作为子表）
	sqlFrom := ""
	sqlFrom = sqlFrom + " from ("
	sqlFrom = sqlFrom + "      select DeviceId,"
	sqlFrom = sqlFrom + "             case when OnTime>='" + fromTime + "' then OnTime"
	sqlFrom = sqlFrom + "                  else '" + fromTime + "'"
	sqlFrom = sqlFrom + "             end OnTime,"
	sqlFrom = sqlFrom + "             case when IFNULL(OffTime,'')='' then '" + toTime + "'"
	sqlFrom = sqlFrom + "                  when OffTime<='" + toTime + "' then OffTime"
	sqlFrom = sqlFrom + "                  else '" + toTime + "'"
	sqlFrom = sqlFrom + "             end OffTime"
	sqlFrom = sqlFrom + "      from   deviceuselog"
	sqlFrom = sqlFrom + "      where  OffTime>='" + fromTime + "' and OnTime<='" + toTime + "'"
	sqlFrom = sqlFrom + " 	          and " + sqlModel
	sqlFrom = sqlFrom + " 	          and " + sqlSite
	sqlFrom = sqlFrom + "      ) a"
	sqlFrom = sqlFrom + "     left join device     d on d.Id = a.DeviceId"
	sqlFrom = sqlFrom + "     left join classrooms r on r.Id = d.ClassroomId"

	//拼接语句
	switch siteType {
	case "":
		sql = sql + "select c.Id SiteId,TIMESTAMPDIFF(SECOND,OnTime,OffTime) UseTime"
		sql = sql + sqlFrom
		sql = sql + " left join floors   f   on f.Id = r.Floorsid"
		sql = sql + " left join building b   on b.Id = f.Buildingid"
		sql = sql + " left join campus   c   on c.Id = b.Campusid"
	case "campus":
		sql = sql + "select b.Id SiteId,TIMESTAMPDIFF(SECOND,OnTime,OffTime) UseTime"
		sql = sql + sqlFrom
		sql = sql + " left join floors   f   on f.Id = r.Floorsid"
		sql = sql + " left join building b   on b.Id = f.Buildingid"
	case "building":
		sql = sql + "select f.Id SiteId,TIMESTAMPDIFF(SECOND,OnTime,OffTime) UseTime"
		sql = sql + sqlFrom
		sql = sql + " left join floors   f   on f.Id = r.Floorsid"
	case "floor":
		sql = sql + "select r.Id SiteId,TIMESTAMPDIFF(SECOND,OnTime,OffTime) UseTime"
		sql = sql + sqlFrom
	}

	//套一层，按SiteId分组求和
	sqlSiteSum := " select SiteId,sum(UseTime) UseTime"
	sqlSiteSum = sqlSiteSum + " from (" + sql + ") aa"
	sqlSiteSum = sqlSiteSum + " group by SiteId"

	//再以位置为依据，左连接上一步按SiteId分组求得的各位置的总使用时间
	sqlAll := ""
	switch siteType {
	case "":
		sqlAll = sqlAll + " select cast(aaa.Id as char) SiteId,aaa.Campusname SiteName,IFNULL(bbb.UseTime,0) UseTime"
		sqlAll = sqlAll + " from campus aaa "
		sqlAll = sqlAll + "      left join (" + sqlSiteSum + ") bbb on bbb.SiteId = aaa.Id "
	case "campus":
		sqlAll = sqlAll + " select cast(aaa.Id as char) SiteId,aaa.Buildingname SiteName,IFNULL(bbb.UseTime,0) UseTime"
		sqlAll = sqlAll + " from building aaa "
		sqlAll = sqlAll + "      left join (" + sqlSiteSum + ") bbb on bbb.SiteId = aaa.Id "
		sqlAll = sqlAll + " where aaa.CampusId=" + siteId
	case "building":
		sqlAll = sqlAll + " select cast(aaa.Id as char) SiteId,aaa.Floorname SiteName,IFNULL(bbb.UseTime,0) UseTime"
		sqlAll = sqlAll + " from floors aaa "
		sqlAll = sqlAll + "      left join (" + sqlSiteSum + ") bbb on bbb.SiteId = aaa.Id "
		sqlAll = sqlAll + " where aaa.BuildingId=" + siteId
	case "floor":
		sqlAll = sqlAll + " select cast(aaa.Id as char) SiteId,aaa.Classroomsname SiteName,IFNULL(bbb.UseTime,0) UseTime"
		sqlAll = sqlAll + " from classrooms aaa "
		sqlAll = sqlAll + "      left join (" + sqlSiteSum + ") bbb on bbb.SiteId = aaa.Id "
		sqlAll = sqlAll + " where aaa.FloorsId=" + siteId
	}

	var deviceSiteUseTime []model.DeviceSiteUseTime
	_, err = dbmap.Select(&deviceSiteUseTime, sqlAll)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "按位置统计设备使用时间错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, deviceSiteUseTime}

	return rd
}

//获取设备型号树
func GetDeviceModelTree(requestData model.RequestData, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//拼接查询条件(where)

	//计算分页信息

	//获得具体数据
	sql = "select Id,IFNULL(PId,'') PId,Name,Type from DeviceModel order by Id "
	var list []model.DeviceModelTree
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//获取故障记录
func GetFault(requestData model.RequestData, id string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var pg model.PageData
	var err error

	//查故障记录
	sql := `
SELECT g.Id, DeviceId,  ifnull(d.Name, '') DeviceName, ifnull(r.Id,0) AS ClassroomID,
	ifnull(CONCAT(c.Campusname, b.Buildingname, f.Floorname, r.Classroomsname), '') DeviceSite,			
	ifnull(FaultSummary, '') FaultSummary,
	ifnull(FaultDescription, '') FaultDescription,
	ifnull(HappenTime, '') HappenTime,
	ifnull(g.IsCanUse, '') IsCanUse,
	ifnull(InputUserId, 0) InputUserId,
	ifnull(u1.Truename, '') InputUserName,
	ifnull(InputTime, '') InputTime,
	ifnull(SubmitTime, '') SubmitTime,
	ifnull(Status, '') Status,
	ifnull(AcceptanceRepairPerson, '') AcceptanceRepairPerson,
	ifnull(AcceptanceRepairPersonTel, '') AcceptanceRepairPersonTel,
	ifnull(AcceptanceUserId, 0) AcceptanceUserId,
	ifnull(u2.Truename, '') AcceptanceUserName,
	ifnull(AcceptanceTime, '') AcceptanceTime,
	ifnull(RepairPerson, '') RepairPerson,
	ifnull(RepairFinishTime, '') RepairFinishTime,
	ifnull(RepairDescription, '') RepairDescription,
	ifnull(RepairIsCanUse, '') RepairIsCanUse,
	ifnull(RepairResult, '') RepairResult,
	ifnull(RepairInputUserId, 0) RepairInputUserId,
	ifnull(u3.Truename, '') RepairInputUserName,
	ifnull(RepairInputTime, '') RepairInputTime,
	ifnull(RepairSubmitTime, '') RepairSubmitTime
FROM DeviceFault g
	LEFT JOIN Device d ON d.Id = g.DeviceId
	LEFT JOIN classrooms r ON r.Id = d.ClassroomId
	LEFT JOIN floors f ON f.Id = r.Floorsid
	LEFT JOIN building b ON b.Id = f.Buildingid
	LEFT JOIN campus c ON c.Id = b.Campusid
	LEFT JOIN Users u1 ON u1.Id = g.InputUserId
	LEFT JOIN Users u2 ON u2.Id = g.AcceptanceUserId
	LEFT JOIN Users u3 ON u3.Id = g.RepairInputUserId
WHERE g.Id = ?
`
	var list = []struct {
		model.DeviceFault
		ClassroomID int // 教室ID
	}{}
	_, err = dbmap.Select(&list, sql, id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	} else if len(list) <= 0 {
		rd.Rcode = "1003"
		rd.Reason = "未查到数据"
		return rd
	}

	//查故障对应的故障类型
	sql = `
SELECT ft.FaultTypeId, ifnull(mft. NAME, '') FaultTypeName
FROM DeviceFaultType ft 
	LEFT JOIN DeviceModelFaultType mft ON mft.Id = ft.FaultTypeId
WHERE ft.FaultId = ?
`
	var list2 []model.RepairFaultType
	_, err = dbmap.Select(&list2, sql, id)
	if err != nil {
		rd.Rcode = "1004"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	list[0].RepairFaultType = list2
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list[0]}
	return rd
}

//获取教室设备
func GetClassroomDevice(requestData model.RequestData, classroomId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = " select Id DeviceId,Name DeviceName from Device where ClassroomId=" + classroomId
	var data []model.ClassroomDevice
	_, err = dbmap.Select(&data, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}
	return rd
}

//获取设备型号对应的所有故障分类
func GetDeviceAllFaultType(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = " select Id FaultTypeId,Name FaultTypeName from DeviceModelFaultType where ModelId in (select ModelId from Device where Id=?)"
	var data []model.DeviceFaultType
	_, err = dbmap.Select(&data, sql, deviceId)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd

}

//获取设备型号对应的所有故障现象词条
func DeviceSiteName(requestData model.RequestData, deviceId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = " select Name  from DeviceModelFaultWord where ModelId in (select ModelId from Device where Id=?)"
	var data []model.DeviceFaultWord
	_, err = dbmap.Select(&data, sql, deviceId)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd

}

//查询故障表信息
func QueryFaultTableInfo(id string, dbmap *gorp.DbMap) (data model.DeviceFaultTable, err error) {
	sql := ""
	sql = sql + " select Id,DeviceId,"
	sql = sql + "        ifnull(FaultSummary,'') FaultSummary,ifnull(FaultDescription,'') FaultDescription,ifnull(HappenTime,'') HappenTime,ifnull(IsCanUse,'') IsCanUse,ifnull(InputUserId,0) InputUserId,ifnull(InputTime,'') InputTime,ifnull(SubmitTime,'') SubmitTime,ifnull(Status,'') Status,"
	sql = sql + "        ifnull(AcceptanceRepairPerson,'') AcceptanceRepairPerson,ifnull(AcceptanceRepairPersonTel,'') AcceptanceRepairPersonTel,ifnull(AcceptanceUserId,0) AcceptanceUserId,ifnull(AcceptanceTime,'') AcceptanceTime,"
	sql = sql + "        ifnull(RepairPerson,'') RepairPerson,ifnull(RepairFinishTime,'') RepairFinishTime,ifnull(RepairDescription,'') RepairDescription,ifnull(RepairIsCanUse,'') RepairIsCanUse,ifnull(RepairResult,'') RepairResult,ifnull(RepairInputUserId,0) RepairInputUserId,ifnull(RepairInputTime,'') RepairInputTime,ifnull(RepairSubmitTime,'') RepairSubmitTime"
	sql = sql + " from   DeviceFault"
	sql = sql + " where  Id = ?"

	err = dbmap.SelectOne(&data, sql, id)
	if err == dbsql.ErrNoRows {
		return data, nil
	}

	return data, err
}

//故障登记——新增
func RegisterFault_Add(t string, r model.RequestRegisterFaultData, d model.DeviceFaultTable, dbmap *gorp.Transaction) error {
	sql := " insert into DeviceFault (Id,DeviceId,FaultSummary,FaultDescription,HappenTime,IsCanUse,InputUserId,InputTime,Status) values(?,?,?,?,?,?,?,?,?) "
	_, err := dbmap.Exec(sql, r.Para.Id, r.Para.DeviceId, r.Para.FaultSummary, r.Para.FaultDescription, r.Para.HappenTime, r.Para.IsCanUse, r.Auth.Usersid, t, "0")
	return err
}

//故障登记——编辑
func RegisterFault_Edit(r model.RequestRegisterFaultData, d model.DeviceFaultTable, dbmap *gorp.Transaction) error {
	//编辑记录时，只能自己编辑自己的，所以where条件中需判断用户
	sql := " update DeviceFault set DeviceId=?,FaultSummary=?,FaultDescription=?,HappenTime=?,IsCanUse=? where Id=? and InputUserId=? and Status='0' "
	sql = `
UPDATE DeviceFault SET DeviceId =?, FaultSummary =?, FaultDescription =?, HappenTime =?, IsCanUse =? 
WHERE Id =? AND Status = '0'
`
	_, err := dbmap.Exec(sql, r.Para.DeviceId, r.Para.FaultSummary, r.Para.FaultDescription, r.Para.HappenTime, r.Para.IsCanUse, r.Para.Id)
	return err
}

//故障登记——提交
func RegisterFault_Submit(t string, deviceId string, faultId string, isCanUse string, userId int, dbmap *gorp.Transaction) error {
	//更改故障状态，填写故障提交时间
	sql := " update DeviceFault set SubmitTime=?,Status='1' where Id=? and Status='0' "
	_, err := dbmap.Exec(sql, t, faultId)
	if err != nil {
		return err
	}

	//更改设备表的设备是否可用字段(Device.IsCanUse)
	sql = " update Device set IsCanUse = ? where Id=? "
	_, err = dbmap.Exec(sql, isCanUse, deviceId)
	return err
}

//故障管理——故障删除
func DeleteFault(faultId string, userId int, dbmap *gorp.Transaction) error {
	sql := " delete from DeviceFault where Id=? and status='0' " //只能删除status='0'(即草稿状态）的故障记录,而且只能是故障申报人自己删除自己的
	_, err := dbmap.Exec(sql, faultId)
	return err
}

//故障受理
func AcceptanceFault(t string, r model.RequestAcceptanceFaultData, dbmap *gorp.Transaction) error {
	//更改故障受理信息
	sql := " update DeviceFault set AcceptanceRepairPerson=?,AcceptanceRepairPersonTel=?,AcceptanceUserId=?,AcceptanceTime=?,Status='2' where Id=? and Status='1' "
	_, err := dbmap.Exec(sql, r.Para.RepairPerson, r.Para.RepairPersonTel, r.Auth.Usersid, t, r.Para.Id)
	return err
}

//维修登记——编辑
func RegisterRepair_Edit(t string, r model.RequestRegisterRepairData, d model.DeviceFaultTable, dbmap *gorp.Transaction) error {
	var sql string

	//暂存维修记录
	if d.RepairInputTime == "" { //第一次暂存，需处理RepairInputTime和RepairInputUserId这两个字段
		sql = " update DeviceFault set RepairPerson=?,RepairFinishTime=?,RepairDescription=?,RepairIsCanUse=?,RepairResult=?,RepairInputTime=?,RepairInputUserId=? where Id=? and Status='2' "
		_, err := dbmap.Exec(sql, r.Para.RepairPerson, r.Para.RepairFinishTime, r.Para.RepairDescription, r.Para.RepairIsCanUse, r.Para.RepairResult, t, r.Auth.Usersid, r.Para.Id)
		if err != nil {
			return err
		}
	} else {
		sql = " update DeviceFault set RepairPerson=?,RepairFinishTime=?,RepairDescription=?,RepairIsCanUse=?,RepairResult=? where Id=? and Status='2' "
		_, err := dbmap.Exec(sql, r.Para.RepairPerson, r.Para.RepairFinishTime, r.Para.RepairDescription, r.Para.RepairIsCanUse, r.Para.RepairResult, r.Para.Id)
		if err != nil {
			return err
		}
	}

	//暂存维修分类记录(先删除)
	sql = " delete from DeviceFaultType where FaultId=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}

	//暂存维修分类记录(再插入)
	if len(r.Para.FaultType) > 0 {
		sql = " insert into DeviceFaultType(FaultId,FaultTypeId)"
		v := ""
		for _, s := range r.Para.FaultType {
			if v != "" {
				v = v + " union "
			}
			v = v + "select '" + r.Para.Id + "','" + s.FaultTypeId + "'"
		}
		_, err := dbmap.Exec(sql + v)
		if err != nil {
			return err
		}
	}

	return nil
}

//维修登记——提交
func RegisterRepair_Submit(t string, r model.RequestRegisterRepairData, d model.DeviceFaultTable, dbmap *gorp.Transaction) error {
	//更改故障状态，填写维修提交时间
	sql := " update DeviceFault set RepairSubmitTime=?,Status='3' where Id=? and Status='2' and (InputUserId=? or AcceptanceUserId=?)"
	_, err := dbmap.Exec(sql, t, r.Para.Id, r.Auth.Usersid, r.Auth.Usersid)
	if err != nil {
		return err
	}

	//更改设备表的设备是否可用字段(Device.IsCanUse)
	sql = " update Device set IsCanUse=? where Id=? "
	_, err = dbmap.Exec(sql, r.Para.RepairIsCanUse, d.DeviceId)
	if err != nil {
		return err
	}

	return nil
}

//节点配置---------------------------------------------

//节点配置-获取节点型号列表
func GetNodeModelList(requestData model.RequestData, keyWord string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := "select Id,Name,ifnull(Description,'') Description  from NodeModel"

	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "Name" + s + "or Description" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.NodeModel
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//节点配置-获取节点型号
func GetNodeModel(requestData model.RequestData, id string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = "select Id,Name,ifnull(Description,'') Description  from NodeModel where Id=?"

	var data model.NodeModel
	err = dbmap.SelectOne(&data, sql, id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//节点配置——保存节点节点型号
func SaveNodeModel(r model.RequestNodeModelData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from NodeModel where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}

	//再插入
	sql = " insert into NodeModel(Id,Name,Description) values(?,?,?)"
	_, err = dbmap.Exec(sql, r.Para.Id, r.Para.Name, r.Para.Description)
	if err != nil {
		return err
	}

	return nil
}

//节点配置——删除节点节点型号
func DeleteNodeModel(id string, dbmap *gorp.Transaction) error {
	var sql string

	//先删除节点型号命令
	sql = " delete from NodeModelCmd where ModelId=?"
	_, err := dbmap.Exec(sql, id)
	if err != nil {
		return err
	}

	//清空节点表中的节点型号字段
	sql = " update Node set ModelId=null where ModelId=?"
	_, err = dbmap.Exec(sql, id)
	if err != nil {
		return err
	}

	//再删除节点型号
	sql = " delete from NodeModel where Id=?"
	_, err = dbmap.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

//节点配置-获取节点型号列表
func GetNodeModelCMDList(requestData model.RequestData, keyWord string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := ` select nmcmd.Id,ModelId,CmdCode,CmdName,RequestURI,URIQuery,ifnull(CmdDescription,'')CmdDescription,
					ifnull(RequestType,'')RequestType,ifnull(Payload,'')Payload,ifnull(CloseCmdFlag,'')CloseCmdFlag,ifnull(OpenCmdFlag,'')OpenCmdFlag,
					nm.Name NodeModelName from nodemodelcmd nmcmd inner join nodemodel nm on nm.Id=nmcmd.ModelId`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "CmdCode" + s + " or CmdName" + s + " or CmdDescription" + s + " or NodeModelName" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		m := " ModelId='" + ModelId + "' "
		sqlWhere = sqlWhere + " and (" + m + ")"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.NodeModelCMD
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//节点配置-获取节点型号
func GetNodeModelCMD(requestData model.RequestData, id int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = ` select nmcmd.Id,ModelId,CmdCode,CmdName,RequestURI,URIQuery,ifnull(CmdDescription,'')CmdDescription,
					ifnull(RequestType,'')RequestType,ifnull(Payload,'')Payload,ifnull(CloseCmdFlag,'')CloseCmdFlag,ifnull(OpenCmdFlag,'')OpenCmdFlag,
					nm.Name NodeModelName from nodemodelcmd nmcmd inner join nodemodel nm on nm.Id=nmcmd.ModelId where nmcmd.Id=?`
	var data model.NodeModelCMD
	err = dbmap.SelectOne(&data, sql, id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//节点配置——保存节点节点型号
func SaveNodeModelCMD(r model.RequestNodeModelCMDData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from NodeModelCMD where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}

	//再插入
	sql = ` insert into NodeModelCMD(ModelId,CmdCode,CmdName,RequestURI,URIQuery,CmdDescription,RequestType,Payload,CloseCmdFlag,OpenCmdFlag)
	values(?,?,?,?,?,?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.ModelId, r.Para.CmdCode, r.Para.CmdName, r.Para.RequestURI, r.Para.URIQuery, r.Para.CmdDescription, r.Para.RequestType, r.Para.Payload, r.Para.CloseCmdFlag, r.Para.OpenCmdFlag)
	if err != nil {
		return err
	}

	return nil
}

//节点配置——删除节点节点型号
func DeleteNodeModelCMD(id int, dbmap *gorp.Transaction) error {
	var sql string

	//先删除节点型号命令
	sql = " delete from NodeModelCmd where Id=?"
	_, err := dbmap.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

//节点配置-获取节点列表
func GetNodeList(requestData model.RequestData, keyWord, NodeId, Campusids, Buildingids, Floorsids, ClassRoomIds, IsNoSave string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := ` select nd.Id,ifnull(nd.Name,'')Name,ifnull(nm.Name,'') NodeModelName,ifnull(crs.Classroomsname,'')Classroomsname,ifnull(bd.Buildingname,'')Buildingname,ifnull(cps.Campusname,'')Campusname,
				ifnull(nd.ModelId,'')ModelId,ifnull(nd.ClassRoomId,0)ClassRoomId,ifnull(bd.Campusid,0)Campusid,ifnull(fs.Buildingid,0)Buildingid,ifnull(crs.Floorsid,0)Floorsid
				,nd.IpType,nd.NodeCoapPort,nd.InRouteMappingPort,nd.RouteIp,nd.UploadTime 
				from node nd left join nodeModel nm on nd.ModelId=nm.Id 
				left join Classrooms crs on nd.ClassRoomId=crs.Id
				left join floors fs on crs.Floorsid=fs.Id
				left join building bd on fs.Buildingid=bd.Id
				left join campus cps on bd.Campusid=cps.Id`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "Name" + s + "or NodeModelName" + s + "or Classroomsname" + s + "or Buildingname" + s + "or Campusname" + s + "or RouteIp='" + keyWord + "' or Id " + s + " or Id='" + keyWord + "'"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if NodeId != "" {
		if NodeId == "-1" {
			m := " (ModelId='') "
			sqlWhere = sqlWhere + " and (" + m + ")"
		} else {
			m := " ModelId='" + NodeId + "' "
			sqlWhere = sqlWhere + " and (" + m + ")"
		}
	}
	var m string
	if IsNoSave != "1" {
		if Campusids != "" {
			m = " Campusid in(" + Campusids + ")"
			sqlWhere = sqlWhere + " and (" + m + ")"
		}
		if Buildingids != "" {
			m = " Buildingid in(" + Buildingids + ") "
			sqlWhere = sqlWhere + " and (" + m + ")"
		}
		if Floorsids != "" {
			m = " Floorsid in(" + Floorsids + ") "
			sqlWhere = sqlWhere + " and (" + m + ")"
		}
		if ClassRoomIds != "" {
			m = " ClassRoomId in(" + ClassRoomIds + ") "
			sqlWhere = sqlWhere + " and (" + m + ")"
		}
	} else {
		m = " ClassRoomId=0"
		sqlWhere = sqlWhere + " and (" + m + ")"
	}
	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.Node
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//节点配置-获取节点
func GetNode(requestData model.RequestData, NodeId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = `select nd.Id,nd.Name,ifnull(nm.Name,'') NodeModelName,ifnull(crs.Classroomsname,'')Classroomsname,ifnull(bd.Buildingname,'')Buildingname,ifnull(cps.Campusname,'')Campusname,
		nd.ModelId,nd.ClassRoomId,nd.IpType,nd.NodeCoapPort,nd.InRouteMappingPort,nd.RouteIp,nd.UploadTime 
		from node nd left join nodeModel nm on nd.ModelId=nm.Id 
		left join Classrooms crs on nd.ClassRoomId=crs.Id
		left join floors fs on crs.Floorsid=fs.Id
		left join building bd on fs.Buildingid=bd.Id
		left join campus cps on bd.Campusid=cps.Id where nd.Id=?`
	var data model.Node
	err = dbmap.SelectOne(&data, sql, NodeId)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//节点配置——保存节点
func SaveNode(r model.RequestNodeData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = "DELETE FROM Node WHERE Id =? "
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}

	//再插入
	sql = `INSERT INTO Node (Id, Name, ModelId, ClassRoomId, IpType, NodeCoapPort, InRouteMappingPort, RouteIp, UploadTime) VALUES (?,?,?,?,?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.Id, r.Para.Name, r.Para.ModelId, r.Para.ClassRoomId, r.Para.IpType, r.Para.NodeCoapPort, r.Para.InRouteMappingPort, r.Para.RouteIp, r.Para.UploadTime)
	if err != nil {
		return err
	}

	return nil
}

//设备配置-获取设备型号列表
func GetDeviceModelList(requestData model.RequestData, keyWord string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := `select Id,ifnull(PId,'')PId,Name,ifnull(Description,'')Description,case Type when 1 then '分类' when 2 then '型号' end TypeName,Type,ifnull(PageFileName,'')PageFileName,
ifnull(ImgFileName,'')ImgFileName,ifnull(ImgFileName2,'')ImgFileName2,ifnull(IsAlert,'')IsAlert,ifnull(MaxUseTime,0)MaxUseTime from devicemodel`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "Name" + s + "or Id='" + keyWord + "'"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		m := "  FIND_IN_SET(Id,getDeviceModelChildNodes('" + ModelId + "'))>0 "
		sqlWhere = sqlWhere + " and (" + m + ")"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.DeviceModel
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取设备型号
func GetDeviceModel(requestData model.RequestData, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	sql = `select Id,PId,Name,Description,case Type when 1 then '分类' when 2 then '型号' end TypeName,Type,PageFileName,ImgFileName,ImgFileName2,IsAlert,MaxUseTime from devicemodel where Id=?`
	var data model.DeviceModel
	err = dbmap.SelectOne(&data, sql, ModelId)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存设备型号
func SaveDeviceModel(r model.RequestDeviceModelData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from DeviceModel where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}

	//再插入
	sql = ` insert into DeviceModel(Id,PId,Name,Description,Type,PageFileName,ImgFileName,ImgFileName2,IsAlert,MaxUseTime)
	values(?,?,?,?,?,?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.Id, r.Para.PId, r.Para.Name, r.Para.Description, r.Para.Type, r.Para.PageFileName, r.Para.ImgFileName, r.Para.ImgFileName2, r.Para.IsAlertValue(), r.Para.MaxUseTime)
	if err != nil {
		return err
	}

	return nil
}

//设备配置——删除设备型号
func DeleteDeviceModel(ModelId string, dbmap *gorp.Transaction) error {
	var sql string

	//先删除设备型号控制命令
	sql = " delete from devicemodelcontrolcmd where ModelId=?"
	_, err := dbmap.Exec(sql, ModelId)
	if err != nil {
		return err
	}
	//先删除设备型号状态值编码
	sql = " delete from devicemodelstatusvaluecode where ModelId=?"
	_, err = dbmap.Exec(sql, ModelId)
	if err != nil {
		return err
	}
	//先删除设备型号状态命令
	sql = " delete from devicemodelstatuscmd where ModelId=?"
	_, err = dbmap.Exec(sql, ModelId)
	if err != nil {
		return err
	}
	//清空设备表中的设备类型字段
	sql = " update Device set ModelId='' where ModelId=?"
	_, err = dbmap.Exec(sql, ModelId)
	if err != nil {
		return err
	}

	//再删除设备类型
	sql = " delete from DeviceModel where Id=?"
	_, err = dbmap.Exec(sql, ModelId)
	if err != nil {
		return err
	}

	return nil
}

//设备配置-获取设备状态列表
func GetDeviceModelStatusCMDList(requestData model.RequestData, keyWord string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := `select dmsc.Id,dmsc.ModelId,dm.Name ModelName,dmsc.Payload,dmsc.StatusName,dmsc.StatusCode,dmsc.StatusValueMatchString,dmsc.SwitchStatusFlag,ifnull(dmsc.OnValue,'')OnValue,
ifnull(dmsc.OffValue,'')OffValue,dmsc.SeqNo,dmsc.SelectValueFlag,ifnull(dmsc.IsAlert,'')IsAlert,ifnull(dmsc.AlertWhere,'')AlertWhere,ifnull(dmsc.AlertDescription,'')AlertDescription 
from devicemodelstatuscmd dmsc inner join devicemodel dm on dmsc.ModelId=dm.Id`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "StatusName" + s + "or ModelName='" + keyWord + "'" + "or StatusCode='" + keyWord + "'"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		m := "  FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + ModelId + "'))>0 "
		sqlWhere = sqlWhere + " and (" + m + ")"
	}

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.DeviceModelStatusCMD
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取设备状态
func GetDeviceModelStatusCMD(requestData model.RequestData, Id int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error
	sql = `select dmsc.Id,dmsc.ModelId,dm.Name ModelName,dmsc.Payload,dmsc.StatusName,dmsc.StatusCode,dmsc.StatusValueMatchString,dmsc.SwitchStatusFlag,ifnull(dmsc.OnValue,'')OnValue,
ifnull(dmsc.OffValue,'')OffValue,dmsc.SeqNo,dmsc.SelectValueFlag,ifnull(dmsc.IsAlert,'')IsAlert,ifnull(dmsc.AlertWhere,'')AlertWhere,ifnull(dmsc.AlertDescription,'')AlertDescription 
from devicemodelstatuscmd dmsc inner join devicemodel dm on dmsc.ModelId=dm.Id where dmsc.Id=?`
	var data model.DeviceModelStatusCMD
	err = dbmap.SelectOne(&data, sql, Id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存设备状态
func SaveDeviceModelStatusCMD(r model.RequestDeviceModelStatusCMDData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from DeviceModelStatusCMD where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}
	//再插入
	sql = ` insert into DeviceModelStatusCMD(ModelId,Payload,StatusName,StatusCode,StatusValueMatchString,SwitchStatusFlag,OnValue,OffValue,SeqNo,SelectValueFlag,IsAlert,AlertWhere,AlertDescription)
	values(?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.ModelId, r.Para.Payload, r.Para.StatusName, r.Para.StatusCode, r.Para.StatusValueMatchString, r.Para.SwitchStatusFlag,
		r.Para.OnValue, r.Para.OffValue, r.Para.SeqNo, r.Para.SelectValueFlag, r.Para.IsAlert, r.Para.AlertWhere, r.Para.AlertDescription)
	if err != nil {
		return err
	}

	return nil
}

//设备配置——删除设备状态
func DeleteDeviceModelStatusCMD(ModelId string, StatusCode string, dbmap *gorp.Transaction) error {
	var sql string

	//先删除设备型号状态值编码
	sql = " delete from devicemodelstatusvaluecode where ModelId=? and StatusCode=?"
	_, err := dbmap.Exec(sql, ModelId, StatusCode)
	if err != nil {
		return err
	}
	//先删除设备型号状态命令
	sql = " delete from devicemodelstatuscmd where ModelId=? and StatusCode=?"
	_, err = dbmap.Exec(sql, ModelId, StatusCode)
	if err != nil {
		return err
	}

	return nil
}

//设备配置-获取设备型号状态编码列表
func GetDeviceModelStatusValueCodeList(requestData model.RequestData, keyWord string, StatusCode string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := `select dmsvc.Id,dmsvc.ModelId,dmsvc.StatusCode,dmsc.StatusName,dm.Name ModelName,dmsvc.StatusValueCode,dmsvc.StatusValueName,ifnull(dmsvc.IsAlert,'')IsAlert  
				from devicemodelstatusvaluecode dmsvc inner join devicemodelstatuscmd dmsc 
				on (dmsvc.ModelId=dmsc.ModelId and dmsvc.StatusCode=dmsc.StatusCode)
				inner join DeviceModel dm on dmsvc.ModelId=dm.Id`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "StatusName" + s + "or ModelName='" + keyWord + "'" + "or StatusCode='" + keyWord + "'" + "or StatusValueCode='" + keyWord + "'" + "or StatusValueName='" + keyWord + "'"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if StatusCode != "" {
		m := " StatusCode='" + StatusCode + "' "
		sqlWhere = sqlWhere + " and (" + m + ")"
	}
	if ModelId != "" {
		s := " ModelId='" + ModelId + "' "
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	//ModelId

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.DeviceModelStatusValueCode
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取设备型号状态编码
func GetDeviceModelStatusValueCode(requestData model.RequestData, Id int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error
	sql = `select dmsvc.Id,dmsvc.ModelId,dmsvc.StatusCode,dmsc.StatusName,dm.Name ModelName,dmsvc.StatusValueCode,dmsvc.StatusValueName,dmsvc.IsAlert 
		from devicemodelstatusvaluecode dmsvc inner join devicemodelstatuscmd dmsc 
		on (dmsvc.ModelId=dmsc.ModelId and dmsvc.StatusCode=dmsc.StatusCode)
		inner join DeviceModel dm on dmsvc.ModelId=dm.Id where dmsvc.Id=?`
	var data model.DeviceModelStatusValueCode
	err = dbmap.SelectOne(&data, sql, Id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存设备型号状态编码
func SaveDeviceModelStatusValueCode(r model.RequestDeviceModelStatusValueCodeData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from devicemodelstatusvaluecode where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}
	//再插入
	sql = ` insert into devicemodelstatusvaluecode(ModelId,StatusCode,StatusValueCode,StatusValueName,IsAlert)
	values(?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.ModelId, r.Para.StatusCode, r.Para.StatusValueCode, r.Para.StatusValueName, r.Para.IsAlert)
	if err != nil {
		return err
	}

	return nil
}

//设备配置——删除设备型号状态编码
func DeleteDeviceModelStatusValueCode(Id int, dbmap *gorp.Transaction) error {
	var sql string

	//先删除设备型号状态值编码
	sql = " delete from devicemodelstatusvaluecode where Id=?"
	_, err := dbmap.Exec(sql, Id)
	if err != nil {
		return err
	}

	return nil
}

//设备配置-获取设备型号状态编码列表
func GetDeviceModelControlCMDList(requestData model.RequestData, keyWord string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	//	sqlTable := `select dmcc.Id,dmcc.ModelId,dmcc.CmdCode,dmcc.CmdName,dmcc.RequestURI,dmcc.URIQuery
	//,ifnull(dmcc.CmdDescription,'')CmdDescription,dmcc.RequestType,dmcc.Payload,ifnull(dmcc.DelayMillisecond,0)DelayMillisecond,
	//ifnull(dmcc.CloseCmdFlag,'')CloseCmdFlag,ifnull(dmcc.OpenCmdFlag,'')OpenCmdFlag,dm.Name ModelIdName
	//from DeviceModelControlCMD dmcc inner join DeviceModel dm on dmcc.ModelId=dm.Id`
	sqlTable := `
SELECT dmcc.Id,
	dmcc.ModelId,
	dmcc.CmdCode,
	dmcc.CmdName,
	dmcc.RequestURI,
	dmcc.URIQuery,
	ifnull(dmcc.CmdDescription, '') CmdDescription,
	dmcc.RequestType,
	dmcc.Payload,
	ifnull(dmcc.DelayMillisecond, 0) DelayMillisecond,
	ifnull(dmcc.CloseCmdFlag, '') CloseCmdFlag,
	ifnull(dmcc.OpenCmdFlag, '') OpenCmdFlag,
	dm.Name ModelIdName
FROM DeviceModelControlCMD dmcc
		INNER JOIN DeviceModel dm ON dmcc.ModelId = dm.Id
`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "CmdCode" + s + "or CmdName='" + keyWord + "'" + "or CmdDescription='" + keyWord + "'" + "or ModelIdName='" + keyWord + "'"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		s := " ModelId='" + ModelId + "' "
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	//ModelId

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.DeviceModelControlCMD
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取设备型号状态编码
func GetDeviceModelControlCMD(requestData model.RequestData, Id int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error
	sql = `select dmcc.Id,dmcc.ModelId,dmcc.CmdCode,dmcc.CmdName,dmcc.RequestURI,dmcc.URIQuery
		,ifnull(dmcc.CmdDescription,'')CmdDescription,dmcc.RequestType,dmcc.Payload,ifnull(dmcc.DelayMillisecond,0)DelayMillisecond,
		ifnull(dmcc.CloseCmdFlag,'')CloseCmdFlag,ifnull(dmcc.OpenCmdFlag,'')OpenCmdFlag,dm.Name ModelIdName 
		from DeviceModelControlCMD dmcc inner join DeviceModel dm on dmcc.ModelId=dm.Id where dmcc.Id=?`
	var data model.DeviceModelControlCMD
	err = dbmap.SelectOne(&data, sql, Id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存设备型号状态编码
func SaveDeviceModelControlCMD(r model.RequestDeviceModelControlCMDData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from DeviceModelControlCMD where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}
	//再插入
	sql = ` insert into DeviceModelControlCMD(ModelId,CmdCode,CmdName,RequestURI,URIQuery,CmdDescription,RequestType,Payload,DelayMillisecond,CloseCmdFlag,OpenCmdFlag)
	values(?,?,?,?,?,?,?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.ModelId, r.Para.CmdCode, r.Para.CmdName, r.Para.RequestURI, r.Para.URIQuery, r.Para.CmdDescription, r.Para.RequestType, r.Para.Payload, r.Para.DelayMillisecond, r.Para.CloseCmdFlag, r.Para.OpenCmdFlag)
	if err != nil {
		return err
	}

	return nil
}

//设备配置——删除设备型号状态编码
func DeleteDeviceModelControlCMD(Id int, dbmap *gorp.Transaction) error {
	var sql string

	//先删除设备型号状态值编码
	sql = " delete from DeviceModelControlCMD where Id=?"
	_, err := dbmap.Exec(sql, Id)
	if err != nil {
		return err
	}

	return nil
}

//设备配置-获取设备列表
func GetDeviceList(requestData model.RequestData, keyWord string, ModelId string, Buildingid int, Floorsid int, Campusid int, ClassroomId int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := `
SELECT d.Id, d.Name, d.Sn, d.Code, d.Brand, d.ModelId,
			dm.Name ModelName,
			IFNULL(d.ClassroomId, 0) ClassroomId,
			IFNULL(cps.Campusname, '') Campusname,
			IFNULL(bd.Buildingname, '') Buildingname,
			IFNULL(crs.Classroomsname, '') Classroomsname,
			d.PowerNodeId,
			ifnull(d.PowerSwitchId, '') PowerSwitchId,
			IFNULL(fl.Buildingid, 0) Buildingid,
			IFNULL(bd.Campusid, 0)  Campusid,
			IFNULL(crs.Floorsid, 0) Floorsid,
			ifnull(d.JoinMethod, '') JoinMethod,
			ifnull(d.JoinNodeId, '') JoinNodeId,
			ifnull(d.JoinSocketId, '') JoinSocketId,
			ifnull(d.NodeSwitchStatus, '') NodeSwitchStatus,
			ifnull(d.NodeSwitchStatusUpdateTime, '') NodeSwitchStatusUpdateTime,
			ifnull(d.DeviceSelfStatus, '') DeviceSelfStatus,
			ifnull(d.DeviceSelfStatusUpdateTime, '') DeviceSelfStatusUpdateTime,
			ifnull(d.IsCanUse, '') IsCanUse,
			ifnull(d.UseTimeBefore, 0) UseTimeBefore,
			ifnull(d.UseTimeAfter, 0) UseTimeAfter,
			ifnull(d.JoinNodeUpdateTime, '') JoinNodeUpdateTime
		FROM Device d
			INNER JOIN DeviceModel dm ON d.ModelId = dm.Id
			LEFT JOIN classrooms crs ON d.ClassroomId = crs.Id
			LEFT JOIN floors fl ON fl.Id = crs.Floorsid
			LEFT JOIN building bd ON bd.Id = fl.Buildingid
			LEFT JOIN campus cps ON cps.Id = bd.Campusid
`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "ModelName" + s + "or Name" + s + "or Brand" + s + "or PowerNodeId" + s + "or JoinNodeId" + s + "or Code" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		s := " FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + ModelId + "'))>0"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if Buildingid > 0 {
		b := " Buildingid=" + strconv.Itoa(Buildingid)
		sqlWhere = sqlWhere + " and (" + b + ")"
	}
	if Floorsid > 0 {
		f := " Floorsid=" + strconv.Itoa(Floorsid)
		sqlWhere = sqlWhere + " and (" + f + ")"
	}
	if Campusid > 0 {
		p := " Campusid=" + strconv.Itoa(Campusid)
		sqlWhere = sqlWhere + " and (" + p + ")"
	}
	if ClassroomId > 0 {
		c := " ClassroomId=" + strconv.Itoa(ClassroomId)
		sqlWhere = sqlWhere + " and (" + c + ")"
	}
	//ModelId

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.Device
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取设备
func GetDevice(requestData model.RequestData, Id string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error
	sql = `select d.Id,d.Name,d.Sn,d.Code,d.Brand,d.ModelId,dm.Name ModelName,d.ClassroomId,cps.Campusname,bd.Buildingname,crs.Classroomsname,d.PowerNodeId,
				ifnull(d.PowerSwitchId,'')PowerSwitchId,fl.Buildingid,bd.Campusid,crs.Floorsid,
				ifnull(d.JoinMethod,'')JoinMethod,ifnull(d.JoinNodeId,'')JoinNodeId,ifnull(d.JoinSocketId,'')JoinSocketId,ifnull(d.NodeSwitchStatus,'')NodeSwitchStatus,
				ifnull(d.NodeSwitchStatusUpdateTime,'')NodeSwitchStatusUpdateTime,ifnull(d.DeviceSelfStatus,'')DeviceSelfStatus,
				ifnull(d.DeviceSelfStatusUpdateTime,'')DeviceSelfStatusUpdateTime,ifnull(d.IsCanUse,'')IsCanUse,ifnull(d.UseTimeBefore,0)UseTimeBefore,
				ifnull(d.UseTimeAfter,0)UseTimeAfter,ifnull(d.JoinNodeUpdateTime,'')JoinNodeUpdateTime
				from Device d inner join DeviceModel dm on d.ModelId=dm.Id
				left join classrooms crs on d.ClassroomId=crs.Id
				left join floors fl on fl.Id=crs.Floorsid
				left join building bd on bd.Id=fl.Buildingid
				left join campus cps on cps.Id=bd.Campusid where d.Id=?`
	var data model.Device
	err = dbmap.SelectOne(&data, sql, Id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存设备
func SaveDevice(r model.RequestDeviceData, dbmap *gorp.Transaction) (err error) {
	//先删除
	sql := " delete from Device where Id=?"
	_, err = dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}

	// 查询绑定节点位置
	query := `SELECT ClassRoomId FROM node WHERE Id = ?`
	room_id, _ := dbmap.SelectInt(query, r.Para.PowerNodeId)
	r.Para.ClassroomId = int(room_id)

	//再插入
	sql = `
INSERT INTO Device (Id, Name, Sn, Code, Brand, ModelId, ClassroomId, PowerNodeId,PowerSwitchId, JoinMethod, 
		JoinNodeId, JoinSocketId, NodeSwitchStatus, NodeSwitchStatusUpdateTime, DeviceSelfStatus, 
		DeviceSelfStatusUpdateTime, IsCanUse, UseTimeBefore, UseTimeAfter, JoinNodeUpdateTime)
VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.Id, r.Para.Name, r.Para.Sn, r.Para.Code, r.Para.Brand,
		r.Para.ModelId, r.Para.ClassroomId, r.Para.PowerNodeId, r.Para.PowerSwitchId,
		r.Para.JoinMethod, r.Para.JoinNodeId, r.Para.JoinSocketId, r.Para.NodeSwitchStatus,
		r.Para.NodeSwitchStatusUpdateTime, r.Para.DeviceSelfStatus, r.Para.DeviceSelfStatusUpdateTime,
		r.Para.IsCanUse, r.Para.UseTimeBefore, r.Para.UseTimeAfter, r.Para.JoinNodeUpdateTime)

	return err
}

//设备配置-获取故障分类列表
func GetDeviceModelFaultTypeList(requestData model.RequestData, keyWord string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := `select dmft.Id,dmft.Name,dm.Name ModelName,dmft.ModelId from devicemodelfaulttype dmft inner join devicemodel dm on dmft.ModelId=dm.Id`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "Name" + s + "or ModelName" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		s := " FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + ModelId + "'))>0"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	//ModelId

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.DeviceModelFaultType
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取故障分类
func GetDeviceModelFaultType(requestData model.RequestData, Id string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error
	sql = `select dmft.Id,dmft.Name,dm.Name ModelName,dmft.ModelId from devicemodelfaulttype dmft inner join devicemodel dm on dmft.ModelId=dm.Id where dmft.Id=?`
	var data model.DeviceModelFaultType
	err = dbmap.SelectOne(&data, sql, Id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存故障分类
func SaveDeviceModelFaultType(r model.RequestDeviceModelFaultTypeData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from DeviceModelFaultType where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}
	//再插入
	sql = ` insert into DeviceModelFaultType(Id,Name,ModelId)
	values(?,?,?)`
	_, err = dbmap.Exec(sql, r.Para.Id, r.Para.Name, r.Para.ModelId)
	if err != nil {
		return err
	}

	return nil
}

//设备配置——删除故障分类
func DeleteDeviceModelFaultType(Id string, dbmap *gorp.Transaction) error {
	var sql string

	//先删除设备型号状态值编码
	sql = " delete from DeviceModelFaultType where Id=?"
	_, err := dbmap.Exec(sql, Id)
	if err != nil {
		return err
	}

	return nil
}

//设备配置-获取故障现象常用词条列表
func GetDeviceModelFaultWordList(requestData model.RequestData, keyWord string, ModelId string, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error

	//后面的记录统计和查询都要的SQL基础语句
	sqlTable := `select dmfw.Id,dmfw.Name,dm.Name ModelName,dmfw.ModelId from devicemodelfaultword dmfw inner join devicemodel dm on dmfw.ModelId=dm.Id`
	//拼接查询条件(where)
	sqlWhere := " where 1=1"
	if keyWord != "" {
		s := " like '%" + keyWord + "%' "
		s = "Name" + s + "or ModelName" + s
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	if ModelId != "" {
		s := " FIND_IN_SET(ModelId,getDeviceModelChildNodes('" + ModelId + "'))>0"
		sqlWhere = sqlWhere + " and (" + s + ")"
	}
	//ModelId

	//计算分页信息
	if requestData.Page.PageIndex > 0 {
		sql = "select count(*) from (" + sqlTable + ") aa" + sqlWhere
		pg, err = GetPageInfo(requestData, sql, dbmap)
		if err != nil {
			rd.Rcode = "1003"
			rd.Reason = "获得分页数据错误:" + err.Error()
			return rd
		}
	}

	//获得具体数据
	sql = "select * from (" + sqlTable + ") aa " + sqlWhere + GetLimitString(pg)
	var list []model.DeviceModelFaultWord
	_, err = dbmap.Select(&list, sql)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}
	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, list}

	return rd
}

//设备配置-获取故障现象常用词条
func GetDeviceModelFaultWord(requestData model.RequestData, Id int, dbmap *gorp.DbMap) (rd core.Returndata) {
	var sql string
	var pg model.PageData
	var err error
	sql = `select dmfw.Id,dmfw.Name,dm.Name ModelName,dmfw.ModelId from devicemodelfaultword dmfw inner join devicemodel dm on dmfw.ModelId=dm.Id where dmfw.Id=?`
	var data model.DeviceModelFaultWord
	err = dbmap.SelectOne(&data, sql, Id)
	if err != nil {
		rd.Rcode = "1003"
		rd.Reason = "数据读取错误:" + err.Error()
		return rd
	}

	rd.Rcode = "1000"
	rd.Reason = "操作成功"
	rd.Result = &model.ResultData{pg, data}

	return rd
}

//设备配置——保存故障现象常用词条
func SaveDeviceModelFaultWord(r model.RequestDeviceModelFaultWordData, dbmap *gorp.Transaction) error {
	var sql string

	//先删除
	sql = " delete from DeviceModelFaultWord where Id=?"
	_, err := dbmap.Exec(sql, r.Para.Id)
	if err != nil {
		return err
	}
	//再插入
	sql = ` insert into DeviceModelFaultWord(Name,ModelId)
	values(?,?)`
	_, err = dbmap.Exec(sql, r.Para.Name, r.Para.ModelId)
	if err != nil {
		return err
	}

	return nil
}

//设备配置——删除故障现象常用词条
func DeleteDeviceModelFaultWord(Id int, dbmap *gorp.Transaction) error {
	var sql string

	//
	sql = " delete from DeviceModelFaultWord where Id=?"
	_, err := dbmap.Exec(sql, Id)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	file, err := ini.Load(core.ConfigFile)
	xerr.ThrowPanic(err)
	if seconds, err := file.Section("coap").Key("offline.timeout.seconds").Int(); err == nil {
		OfflineTime = fmt.Sprintf("%d", seconds)
	}
	log.Printf("<<<<<<<<<<\t	CoAP.Offline.Timeout.Seconds=%v", OfflineTime)
}
