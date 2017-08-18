package dal

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"

	"dev.project/BackEndCode/devcontrol/model"
)

//查询设备基本信息
func Query_DeviceBaseInfo(pid, ptype string, dbmap *gorp.DbMap) (data []model.DeviceBasciInfoView) {
	sql := `
SELECT a.Id DeviceId,a. NAME DeviceName,
	IFNULL(b.imgFileName, '') DeviceImg,
	IFNULL(b.imgFileName2, '') DeviceImg2,
	IFNULL(b.PageFileName, '') DevicePage,
	IFNULL(PowerNodeId, '') PowerNodeId,
	IFNULL(PowerSwitchId, '') PowerSwitchId,
	IFNULL(JoinMethod, '') JoinMethod,
	IFNULL(JoinNodeId, '') JoinNodeId,
	IFNULL(JoinSocketId, '') JoinSocketId,
	IFNULL(NodeSwitchStatus, '') NodeSwitchStatus,
	IFNULL(NodeSwitchStatusUpdateTime, '') NodeSwitchStatusUpdateTime,
	IFNULL(DeviceSelfStatus, 'off') DeviceSelfStatus,
	IFNULL(DeviceSelfStatusUpdateTime,'') DeviceSelfStatusUpdateTime,
	IFNULL(IsCanUse, '1') IsCanUse,
	IFNULL(UseTimeBefore, 0) UseTimeBefore,
	IFNULL(UseTimeAfter, 0) UseTimeAfter,
	IFNULL(JoinNodeUpdateTime, '') JoinNodeUpdateTime,
	IFNULL(alert.IsHaveAlert, '0') IsHaveAlert
FROM device a
	LEFT JOIN devicemodel b ON a.ModelId = b.id
	LEFT JOIN (SELECT DISTINCT DeviceId,'1' IsHaveAlert FROM DeviceAlert) alert ON alert.DeviceId = a.id
`
	switch ptype {
	case "device":
		sql += " where a.Id='" + pid + "'"
	case "classroom":
		sql += " where a.ClassroomId=" + pid
	}

	_, err := dbmap.Select(&data, sql)
	xdebug.LogError(err)
	return data
}
