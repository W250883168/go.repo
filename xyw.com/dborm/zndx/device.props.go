package zndx

import (
	"fmt"
	"log"
	"runtime"

	"gopkg.in/gorp.v1"
)

/*
CREATE TABLE `device_props` (
  `device_id` varchar(50) NOT NULL,
  `k` varchar(255) NOT NULL DEFAULT '' COMMENT '键名称',
  `v` varchar(1024) NOT NULL DEFAULT '' COMMENT '键值',
  `comments` varchar(255) NOT NULL,
  PRIMARY KEY (`device_id`,`k`),
  CONSTRAINT `device_props_fk_device_id` FOREIGN KEY (`device_id`) REFERENCES `device` (`Id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB DEFAULT CHARSET=utf8
	COMMENT='设备属性键值对表, 如:
	device.accumulate.time.seconds=300
	device.state.online=true/false
	device.used.time.before.seconds=500
	device.used.time.after.seconds=500';
*/
// 设备键值对属性
type DeviceProp struct {
	DeviceID string `db:"device_id"`
	K        string `db:"k"`
	V        string `db:"v"`
	Comments string `db:"comments"`
}

// 获取设备属性
func DeviceProp_Get(id, k string, pDBMap *gorp.DbMap) (pProp *DeviceProp, err error) {
	pDBMap.AddTableWithName(DeviceProp{}, "device_props").SetKeys(false, `device_id`, `k`)
	pObj, err := pDBMap.Get(DeviceProp{}, id, k)

	pProp, _ = pObj.(*DeviceProp)
	return pProp, err
}

// 设备属性存在？
func DeviceProp_Exists(id, k string, pDBMap *gorp.DbMap) (exist bool) {
	p, err := DeviceProp_Get(id, k, pDBMap)
	exist = (err == nil) && (p != nil)

	return exist
}

// 插入
func (p *DeviceProp) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceProp{}, "device_props").SetKeys(false, `device_id`, `k`)
	return pDBMap.Insert(p)
}

func (p *DeviceProp) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceProp{}, "device_props").SetKeys(false, `device_id`, `k`)
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

// 更新
func (p *DeviceProp) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceProp{}, "device_props").SetKeys(false, `device_id`, `k`)
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
