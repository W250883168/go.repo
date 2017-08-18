package dal

import (
	"log"
	"runtime"
	"strings"

	"gopkg.in/gorp.v1"

	"dborm/zndx"
	"xutils/xdebug"
	"xutils/xerr"
	"xutils/xtime"
)

// 根据状态条件，更新设备预警值
func Refresh_DeviceAlert_ByStatusCode(joinNodeId string, joinSocketId string, statusCmd zndx.DeviceModelStatusCmd, statusValueCode string, dbmap *gorp.DbMap) {
	defer xerr.CatchPanic()

	// 先删除对应的预警记录
	devAlert := zndx.DeviceAlert{StatusCode: statusCmd.StatusCode}
	err := devAlert.Delete_ByStatusCodeAndDevice(&zndx.Device{JoinNodeId: joinNodeId, JoinSocketId: joinSocketId}, dbmap)
	xerr.ThrowPanic(err)

	//如果不需预警，如果预警条件为空，如果预警描述为空，则不作预警处理
	if statusCmd.IsAlert == "0" || statusCmd.AlertWhere == "" || statusCmd.AlertDescription == "" {
		return
	}

	//替换预警条件中的"{val}"为真实的值(statusValueCode),例如投影仪的灯泡预警条件可定义为：{val}>1500
	alertWhere := strings.Replace(statusCmd.AlertWhere, "{val}", statusValueCode, -1)

	//检查是否满足预警条件（通过构造一个SQL语句来判断条件是否满足）
	sql := "select 1 from ( select 1) a where ?"
	val, _ := dbmap.SelectNullInt(sql, alertWhere)
	if val.Valid {
		//程序能运行到这里，表示满足预警条件，接下来就是产生新的预警信息
		//创建新的预警消息(通过连接节点ID和插口ID获得设备ID，根据这两个ID查询设备，看起来好象会找到多个设备，但实际上只有一个设备，因为一个插口只能接一个设备)
		pAlert := &zndx.DeviceAlert{AlertDescription: statusCmd.AlertDescription, StatusCode: statusCmd.StatusCode}
		err = pAlert.Insert_ByStatusCodeAndDevice(&zndx.Device{JoinNodeId: joinNodeId, JoinSocketId: joinSocketId}, dbmap)
		xdebug.LogError(err)
	}
}

// 按具体状态值, 更新设备预警值
func Refresh_DeviceAlert_ByStatusValue(joinNodeId string, joinSocketId string, statusCmd zndx.DeviceModelStatusCmd, statusValueCode string, dbmap *gorp.DbMap) {
	defer xerr.CatchPanic()

	//先删除对应的预警记录
	devAlert := zndx.DeviceAlert{StatusCode: statusCmd.StatusCode, StatusValueCode: statusValueCode}
	err := devAlert.Delete_ByStatusValueAndDevice(&zndx.Device{JoinNodeId: joinNodeId, JoinSocketId: joinSocketId}, dbmap)
	xerr.ThrowPanic(err)

	if statusCmd.SelectValueFlag == "0" {
		return
	}

	//查询该状态代码是否需要预警
	sql := "select StatusValueName from DeviceModelStatusValueCode where ModelId=? and StatusCode=? and StatusValueCode=? and IsAlert='1'"
	valName, _ := dbmap.SelectNullStr(sql, statusCmd.ModelId, statusCmd.StatusCode, statusValueCode)
	if valName.Valid {
		//程序能运行到这里，表示满足预警条件，接下来就是产生新的预警信息
		//创建新的预警消息(通过连接节点ID和插口ID获得设备ID，根据这两个ID查询设备，看起来好象会找到多个设备，但实际上只有一个设备，因为一个插口只能接一个设备)
		pAlert := &zndx.DeviceAlert{AlertDescription: valName.String, LastAlertTime: xtime.NowString(), StatusCode: statusCmd.StatusCode, StatusValueCode: statusValueCode}
		err = pAlert.Insert_ByStatusValueAndDevice(&zndx.Device{JoinNodeId: joinNodeId, JoinSocketId: joinSocketId}, dbmap)
		xdebug.LogError(err)
	}
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

}
