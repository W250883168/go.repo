package zndx

import (
	"strconv"
	"strings"

	"gopkg.in/gorp.v1"

	"xutils/xdebug"
	"xutils/xhttp"
)

/*
CREATE TABLE `devicedetaillog` (
  `Id` bigint(20) NOT NULL AUTO_INCREMENT,
  `OperateTime` varchar(50) DEFAULT NULL COMMENT 'yyyy-MM-dd HH:mm:ss',
  `OperateUserId` int(11) DEFAULT NULL,
  `DeviceId` varchar(50) DEFAULT NULL,
  `CmdCode` varchar(50) DEFAULT NULL,
  `CmdName` varchar(50) DEFAULT NULL,
  `Para` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=4102 DEFAULT CHARSET=gbk COMMENT='记录设备的每一次操作';
*/
type DeviceDetailLog struct {
	Id            int
	OperateUserId int    // 操作人ID
	OperateTime   string // 操作时间
	DeviceId      string // 设备编号
	CmdCode       string // 操作代码
	CmdName       string // 操作名称
	Para          string // 参数
}

// 设备详细操作时间
type DeviceDetailLogView struct {
	OperateTime string //操作时间
	UserCode    string //操作人（代码）
	UserName    string //操作人（姓名）
	CmdName     string //操作名称
	Para        string //操作参数
}

// 插入
func (o *DeviceDetailLog) Insert(dbmap *gorp.DbMap) error {
	sql := `INSERT INTO DeviceDetailLog (OperateTime, OperateUserId, DeviceId, CmdCode, CmdName, Para) VALUES (?,?,?,?,?,?)`
	args := []interface{}{o.OperateTime, o.OperateUserId, o.DeviceId, o.CmdCode, o.CmdName, o.Para}
	_, err := dbmap.Exec(sql, args...)
	return err
}

// 删除
func (p *DeviceDetailLog) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceDetailLog{}, "devicedetaillog").SetKeys(true, "Id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 更新
func (p *DeviceDetailLog) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceDetailLog{}, "devicedetaillog").SetKeys(true, "Id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

// 获取设备详细日志
func DeviceDetailLog_Get(id int, pDBMap *gorp.DbMap) (pLog *DeviceDetailLog, err error) {
	pDBMap.AddTableWithName(DeviceDetailLog{}, "devicedetaillog").SetKeys(true, "Id")
	pObj, err := pDBMap.Get(DeviceDetailLog{}, id)

	pLog, _ = pObj.(*DeviceDetailLog)
	return pLog, err
}

// 日志存在?
func DeviceDetailLog_Exists_ByDeviceID(devID string, pDBMap *gorp.DbMap) (exist bool) {
	exist = true
	query := `SELECT COUNT(1) FROM devicedetaillog WHERE DeviceId = ?`
	count, err := pDBMap.SelectInt(query, devID)
	xdebug.LogError(err)
	if err == nil {
		exist = count > 0
	}

	return exist
}

// 获取设备操作日志
func (p *DeviceDetailLog) GetByDeviceID(pPage *xhttp.PageInfo, dbmap *gorp.DbMap) (list []DeviceDetailLogView, err error) {
	// 分页数据
	sqlcount := `SELECT COUNT(*) FROM DeviceDetailLog WHERE DeviceId = ? `
	count, err := dbmap.SelectInt(sqlcount, p.DeviceId)
	if err == nil {
		pPage.RowTotal = int(count)

		//获得具体数据
		sql := `SELECT d.OperateTime, IFNULL(u.Loginuser, '') UserCode, IFNULL(u.TrueName, '') UserName,  IFNULL(d.CmdName, '') CmdName, IFNULL(d.Para, '') Para
				FROM DeviceDetailLog d JOIN Users u ON (d.OperateUserId = u.Id)
				WHERE DeviceId = ?
				ORDER BY d.OperateTime DESC
				` + pPage.SQL_LimitString()
		_, err = dbmap.Select(&list, sql, p.DeviceId)
	}

	return list, err
}

//对教室或设备进行开关操作时，使用以下结构和方法生成设备的详细日志
type DeviceOnOffLogView struct {
	OperateTime     string
	OperateUserId   string
	OperateObject   string // classroom/device
	ObjectId        string // classroomid/deviceid
	OperateType     string // on/off
	OperateTypeName string // 开启/关闭
	Para            string
}

// 创建日志
func (o *DeviceOnOffLogView) CreateLog(dbmap *gorp.DbMap) {
	fFieldName := ""
	if o.OperateType == "on" {
		fFieldName = "OpenCmdFlag"
	} else {
		fFieldName = "CloseCmdFlag"
	}

	switch o.OperateObject {
	case "classroom":
		sql := `
INSERT INTO DeviceDetailLog (OperateTime, OperateUserId, DeviceId, CmdCode, CmdName, Para ) 
	SELECT ?, ?, d.id, ?, ?, ? FROM device d 
	WHERE d.classroomid =? 
			AND (EXISTS (SELECT 1 FROM DeviceModelControlCmd WHERE ModelId = d.ModelId AND @fFieldName = '1') 
					OR EXISTS (SELECT 1 FROM NodeModelCmd WHERE ModelId IN (SELECT ModelId FROM Node WHERE Id = d.PowerNodeId) AND @fFieldName = '1'))
`
		sql = strings.Replace(sql, "@fFieldName", fFieldName, -1)
		_, err := dbmap.Exec(sql, o.OperateTime, o.OperateUserId, o.OperateType, o.OperateTypeName, o.Para, o.ObjectId)
		xdebug.LogError(err)
	case "device":
		userid, _ := strconv.Atoi(o.OperateUserId)
		detailLog := DeviceDetailLog{Id: 0,
			OperateUserId: userid,
			OperateTime:   o.OperateTime,
			DeviceId:      o.ObjectId,
			CmdCode:       o.OperateType,
			CmdName:       o.OperateTypeName,
			Para:          o.Para}
		err := detailLog.Insert(dbmap)
		xdebug.LogError(err)
	case "floor":
		sql := `
INSERT INTO DeviceDetailLog (OperateTime, OperateUserId, DeviceId, CmdCode, CmdName, Para) 
		SELECT ?, ?, d.id, ?, ?, ? 
		FROM device d 
		WHERE d.classroomid IN (SELECT id FROM classrooms WHERE Floorsid =?) 
				AND (EXISTS (SELECT 1 FROM DeviceModelControlCmd WHERE ModelId = d.ModelId AND @fFieldName = '1')
						OR EXISTS (SELECT 1 FROM NodeModelCmd WHERE ModelId IN (SELECT ModelId FROM Node WHERE Id = d.PowerNodeId) AND @fFieldName = '1'))
`
		sql = strings.Replace(sql, "@fFieldName", fFieldName, -1)
		_, err := dbmap.Exec(sql, o.OperateTime, o.OperateUserId, o.OperateType, o.OperateTypeName, o.Para, o.ObjectId)
		xdebug.LogError(err)
	}

}
