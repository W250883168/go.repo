package dborm

import (
	"fmt"
	"log"
	"runtime"

	"gopkg.in/gorp.v1"
	"vodx/ioutil/dbutil"
)

/*
CREATE TABLE `camera_info` (
  `camera_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `location_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '位置ID',
  `camera_name` varchar(255) NOT NULL DEFAULT '' COMMENT '摄像头名称',
  `camera_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '摄像头IP',
  `camera_port` int(11) NOT NULL DEFAULT '0' COMMENT '摄像头端口',
  `login_account` varchar(255) NOT NULL DEFAULT '' COMMENT '登录账号',
  `login_password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码',
  PRIMARY KEY (`camera_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='摄像机信息表';
*/

type CameraInfo struct {
	CameraID      int    `db:"camera_id"`
	LocationID    int    `db:"location_id"`
	CameraName    string `db:"camera_name"`
	CameraIP      string `db:"camera_ip"`
	CameraPort    int    `db:"camera_port"`
	LoginAccount  string `db:"login_account"`
	LoginPassword string `db:"login_password"`
}

func CameraInfo_Query(id int, pDBMap *gorp.DbMap) (pCamera *CameraInfo, err error) {
	pDBMap.AddTableWithName(CameraInfo{}, "camera_info").SetKeys(true, "camera_id")
	pObj, err := pDBMap.Get(CameraInfo{}, id)
	pCamera, _ = pObj.(*CameraInfo)

	return pCamera, err
}

func CameraInfo_QueryByLocation(locationID int, pDBMap *gorp.DbMap) (list []CameraInfo, err error) {
	query := `
SELECT
	TCamera.camera_id,
	TCamera.location_id,
	TCamera.camera_name,
	TCamera.camera_ip,
	TCamera.camera_port,
	TCamera.login_account,
	TCamera.login_password
FROM camera_info TCamera
WHERE location_id = ?
`
	list = []CameraInfo{}
	_, err = pDBMap.Select(&list, query, locationID)
	return list, err
}

func (p *CameraInfo) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(CameraInfo{}, "camera_info").SetKeys(true, "camera_id")
	return pDBMap.Insert(p)
}

func (p *CameraInfo) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(CameraInfo{}, "camera_info").SetKeys(true, "camera_id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *CameraInfo) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(CameraInfo{}, "camera_info").SetKeys(true, "camera_id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	pDBMap := dbutil.GetDBMap()
	pDBMap.AddTableWithName(CameraInfo{}, "camera_info").SetKeys(true, "camera_id")
}
