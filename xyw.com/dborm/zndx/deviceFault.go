package zndx

import (
	"strings"

	"gopkg.in/gorp.v1"

	"xutils/xdebug"
	"xutils/xhttp"
	"xutils/xtext"
)

/*
CREATE TABLE `devicefault` (
  `Id` varchar(50) NOT NULL,
  `DeviceId` varchar(50) DEFAULT NULL,
  `FaultSummary` varchar(500) DEFAULT NULL,
  `FaultDescription` varchar(1000) DEFAULT NULL,
  `HappenTime` varchar(50) DEFAULT NULL,
  `IsCanUse` varchar(1) DEFAULT NULL COMMENT '0-不可使用 1-可以使用',
  `InputUserId` int(11) DEFAULT NULL COMMENT '取当前登录用户，不可改',
  `InputTime` varchar(50) DEFAULT NULL COMMENT '取当前系统时间，不可改',
  `SubmitTime` varchar(50) DEFAULT NULL COMMENT '提交故障时间',
  `Status` varchar(1) DEFAULT NULL COMMENT '0-草稿 1-待受理 2-维修中 3-已维修',
  `AcceptanceRepairPerson` varchar(50) DEFAULT NULL COMMENT '手工录入',
  `AcceptanceRepairPersonTel` varchar(50) DEFAULT NULL COMMENT '手工录入，逗号隔开的多个手机，便于以后发送短信通知',
  `AcceptanceUserId` int(11) DEFAULT NULL COMMENT '取当前登录用户，不可改',
  `AcceptanceTime` varchar(50) DEFAULT NULL COMMENT '取当前系统时间，不可改',
  `RepairPerson` varchar(50) DEFAULT NULL COMMENT '默认为受理时指定的维修人(AcceptanceRepairPerson)，可修改',
  `RepairFinishTime` varchar(50) DEFAULT NULL COMMENT '默认取当前系统时间，可修改',
  `RepairDescription` varchar(1000) DEFAULT NULL,
  `RepairIsCanUse` varchar(1) DEFAULT NULL COMMENT '0-不可使用 1-可以使用',
  `RepairResult` varchar(1) DEFAULT NULL COMMENT '1-未修复 2-已修复',
  `RepairInputUserId` int(11) DEFAULT NULL COMMENT '取当前登录用户，不可改',
  `RepairInputTime` varchar(50) DEFAULT NULL COMMENT '取当前系统时间，不可改',
  `RepairSubmitTime` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=gbk;
*/
type DeviceFault struct {
	Id                        string `db: "size:50"`
	DeviceId                  string `db: "size:50"`   // 设备ID
	FaultSummary              string `db: "size:500"`  // 故障现象
	FaultDescription          string `db: "size:1000"` // 故障描述
	HappenTime                string `db: "size:50"`   // 故障发生时间
	IsCanUse                  string `db: "size:1"`    // 是否可用
	InputUserId               int    ``                // 申报人ID
	InputTime                 string `db: "size:50"`   // 申报时间
	SubmitTime                string `db: "size:50"`   // 提交时间
	Status                    string `db: "size:1"`    // 故障状态
	AcceptanceRepairPerson    string `db: "size:50"`   // 维修受理人
	AcceptanceRepairPersonTel string `db: "size:50"`   // 维修人电话
	AcceptanceUserId          int    ``                // 维修人ID
	AcceptanceTime            string `db: "size:50"`   // 维修受理时间
	RepairPerson              string `db: "size:50"`   // 维修人
	RepairFinishTime          string `db: "size:50"`   // 维修完成时间
	RepairDescription         string `db: "size:1000"` // 维修描述
	RepairIsCanUse            string `db: "size:1"`    // 维修后设备是否可用
	RepairResult              string `db: "size:1"`    // 维修结果
	RepairInputUserId         int    ``                // 维修登记人
	RepairInputTime           string `db: "size:50"`   // 维修登记时间
	RepairSubmitTime          string `db: "size:50"`   // 维修提交时间
}

// 插入
func (p *DeviceFault) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceFault{}, "devicefault").SetKeys(false, "Id")
	err = pDBMap.Insert(p)
	xdebug.LogErrorText(err)

	return err
}

// 删除
func (p *DeviceFault) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceFault{}, "devicefault").SetKeys(false, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 更新
func (p *DeviceFault) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceFault{}, "devicefault").SetKeys(false, "Id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

// 查询对应设备ID的故障存在?
func DeviceFault_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := "SELECT COUNT(1) FROM devicefault WHERE DeviceId = ?"
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

/***********************************
siteType: campus/building/floor/classroom
modelId: 设备型号ID
keyword: 查询关键词
pPage: 分页查询条件(输入输出参数，成功返回PageInfo.RowTotal; 用于计算总页数)
DeviceFault.InputUserId: 用户ID
DeviceFault.Status: 设备状态(-1/全部)
*********************************/
// 根据条件, 查询设备故障信息列表
func (p *DeviceFault) Query_DeviceFault(siteType, siteId, modelId, keyWord string, dbmap *gorp.DbMap, pPage *xhttp.PageInfo) (list []DeviceFaultView, err error) {
	var sql string
	var where string
	var sqlargs = map[string]interface{}{}

	//后面的记录统计和查询都要的SQL(虚拟成一个表)
	subtable := `
SELECT TFault.Id FaultId, TFault.DeviceId, TFault.HappenTime, TFault.FaultSummary, TFault.Status, TFault.InputUserId,
		CASE TFault.IsCanUse WHEN '0' THEN '不可用' WHEN '1' THEN '可用' END IsCanUse,
		CASE TFault.Status WHEN '0' THEN '待处理' WHEN '1' THEN '待处理' WHEN '2' THEN '处理中' WHEN '3' THEN '已处理' END StatusName,		
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
	if all := (p.Status == "-1"); all { // 全部状态
		where += `WHERE	(Status IN ('1', '2', '3') OR (Status = '0' AND InputUserId = :UserID)) `
		sqlargs["UserID"] = p.InputUserId // 草稿状态只能查询到创建人自己的
	} else if p.Status == "1" { // 待处理
		where += `WHERE	(Status = :Status OR (Status = '0' AND InputUserId = :UserID))	`
		sqlargs["UserID"] = p.InputUserId // 草稿状态只能查询到创建人自己的
		sqlargs["Status"] = p.Status
	} else {
		where += `WHERE	(Status = :Status)	`
		sqlargs["Status"] = p.Status
	}

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
		where += `AND (HappenTime LIKE :Keyword 
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

	// 计算分页信息
	sql = "select count(*) from (" + subtable + ") aa	" + where
	count, err := dbmap.SelectInt(sql, sqlargs)
	if err != nil {
		xdebug.LogErrorText(err)
		return list, err
	}

	// 获得具体数据
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
FROM ( @SubTable ) aa  
@SqlWhere  
ORDER BY HappenTime DESC
`
	sql += pPage.SQL_LimitString()
	sql = strings.Replace(sql, "@SubTable", subtable, 1)
	sql = strings.Replace(sql, "@SqlWhere", where, 1)
	_, err = dbmap.Select(&list, sql, sqlargs)
	xdebug.LogErrorText(err)

	pPage.RowTotal = int(count)
	return list, err
}

// 设备故障数据视图
type DeviceFaultView struct {
	FaultId       string
	DeviceId      string
	DeviceCode    string
	DeviceName    string
	DeviceModel   string
	DeviceSite    string
	HappenTime    string
	FaultSummary  string
	IsCanUse      string
	Status        string
	StatusName    string
	InputUserId   int
	InputUserName string
}
