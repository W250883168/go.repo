package zndx

import (
	"fmt"
	"log"
	"runtime"

	"gopkg.in/gorp.v1"
)

/*
CREATE TABLE `device_status_templates` (
  `template_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `device_model_id` varchar(50) NOT NULL DEFAULT '' COMMENT '设备型号ID',
  `status_index` int(11) NOT NULL DEFAULT '0' COMMENT '设备状态序',
  `status_name` varchar(255) NOT NULL DEFAULT '' COMMENT '设备状态名称',
  `status_code` varchar(255) NOT NULL DEFAULT '' COMMENT '设备状态键名称',
  `status_value` varchar(1024) NOT NULL DEFAULT '' COMMENT '设备状态值',
  `normal_value` varchar(1024) NOT NULL DEFAULT '' COMMENT '正常状态值(或范围)',
  `match_mode` enum('none','=','!=','>','<','>=','<=','between','in') NOT NULL DEFAULT 'none' COMMENT '匹配模式(设备状态值比较)',
  PRIMARY KEY (`template_id`),
  KEY `device_alert_templates_fk_device_id` (`device_model_id`) USING BTREE,
  CONSTRAINT `device_alert_templates_fk_device_id` FOREIGN KEY (`device_model_id`) REFERENCES `device` (`Id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='设备状态值告警表， 通过设备当前状态值与正常标准值作比较，得出设备告警状态';
*/

const (
	DeviceStatusTemplate_EMatchMode_none = iota
	DeviceStatusTemplate_EMatchMode_eq
	DeviceStatusTemplate_EMatchMode_neq
	DeviceStatusTemplate_EMatchMode_gt
	DeviceStatusTemplate_EMatchMode_lt
	DeviceStatusTemplate_EMatchMode_gteq
	DeviceStatusTemplate_EMatchMode_lteq
	DeviceStatusTemplate_EMatchMode_between
	DeviceStatusTemplate_EMatchMode_in
)

type DeviceStatusTemplate_EMatchMode int

func (p DeviceStatusTemplate_EMatchMode) String() (str string) {
	switch p {
	case DeviceStatusTemplate_EMatchMode_none:
		str = "none"
	case DeviceStatusTemplate_EMatchMode_eq:
		str = "="
	case DeviceStatusTemplate_EMatchMode_neq:
		str = "!="
	case DeviceStatusTemplate_EMatchMode_lt:
		str = "<"
	case DeviceStatusTemplate_EMatchMode_gt:
		str = ">"
	case DeviceStatusTemplate_EMatchMode_gteq:
		str = ">="
	case DeviceStatusTemplate_EMatchMode_lteq:
		str = "<="
	case DeviceStatusTemplate_EMatchMode_between:
		str = "between"
	case DeviceStatusTemplate_EMatchMode_in:
		str = "in"
	}

	return str
}

type DeviceStatusTemplate struct {
	TemplateID    int    `db:"template_id"`
	DeviceModelID string `db:"device_model_id"`
	StatusIndex   int    `db:"status_index"`
	StatusName    string `db:"status_name"`
	StatusCode    string `db:"status_code"`
	StatusValue   string `db:"status_value"`
	NormalValue   string `db:"normal_value"`
	MatchMode     string `db:"match_mode"`
}

func DeviceStatusTemplate_Get(tempID int, pDBMap *gorp.DbMap) (pTemplate *DeviceStatusTemplate, err error) {
	pDBMap.AddTableWithName(DeviceStatusTemplate{}, `device_status_templates`).SetKeys(true, `template_id`)
	pObj, err := pDBMap.Get(DeviceStatusTemplate{}, tempID)

	pTemplate, _ = pObj.(*DeviceStatusTemplate)
	return pTemplate, err
}

func DeviceStatusTemplate_QueryByModelID(modelID string, pDBMap *gorp.DbMap) (list []DeviceStatusTemplate, err error) {
	list = []DeviceStatusTemplate{}
	query := `
SELECT * FROM device_status_templates 
WHERE device_model_id = ?
ORDER BY status_index
`
	_, err = pDBMap.Select(&list, query, modelID)
	return list, err
}

func (p *DeviceStatusTemplate) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(DeviceStatusTemplate{}, `device_status_templates`).SetKeys(true, `template_id`)
	return pDBMap.Insert(p)
}

func (p *DeviceStatusTemplate) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceStatusTemplate{}, `device_status_templates`).SetKeys(true, `template_id`)
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *DeviceStatusTemplate) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(DeviceStatusTemplate{}, `device_status_templates`).SetKeys(true, `template_id`)
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

func (p *DeviceStatusTemplate) IsAlert() (yes bool) {
	mode := DeviceStatusTemplate_EMatchMode(DeviceStatusTemplate_EMatchMode_none)
	yes = (p.MatchMode == mode.String())

	return yes
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
