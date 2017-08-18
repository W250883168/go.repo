package coapview

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"

	"xutils/xtext"
)

type DeviceDataView struct {
	R0, R1, R2, R3, R4, R5, R6, R7, R8, R9 string
}

func (p *DeviceDataView) StringArray() [10]string {
	return [10]string{
		p.R0, p.R1, p.R2, p.R3,
		p.R4, p.R5, p.R6, p.R7,
		p.R8, p.R9,
	}
}

// 设备状态信息
type DeviceStatView struct {
	Id   string         // 设备编号
	Type string         // 设备控制类型
	Data DeviceDataView // 设备数据
}

// 节点状态信息
type NodeStatView struct {
	Eui  string // 设备编号(必选项)
	Rssi string // 信号强度值(必选项)
	Vdd  string // 电压值(可选项)
	Type int    // 节点类型(enum=1:广播/2:计量/3:非计量/4:红外/5:单火)

	S_addr string // 上报服务器地址
	B_cast string // 广播标志
}

// 是否为广播
func (p *NodeStatView) IsBroadcast() (yes bool) {
	return xtext.IsNotBlank(p.S_addr) && (p.B_cast == "1")
}

// 验证节点EUI
func (p *NodeStatView) Validate() {
	if xtext.IsBlank(p.Eui) {
		panic(errors.New("收到节点上报数据,eui=空 （系统不处理）"))
	}
}

// 开关状态信息
type SwitchStatView struct {
	Id   string // 编号(必选项)
	Stat string // 状态(必选项)
	I    string // 电流值(可选项)
}

// 计量信息
type PowerMeterView struct {
	I int     // 电流（单位：mA）
	V int     // 电压（单位：V）
	F int     // 频率（单位：Hz）
	P float64 // 功率（单位：W）
	W float64 // 电能（单位：KWh）
}

// 负载信息(节点上传数据）
type DevPayloadView struct {
	Node NodeStatView     // 设备节点(必选项)
	Sw   []SwitchStatView // 开关状态列表(可选项)
	Dev  []DeviceStatView // 设备状态列表(可选项)
}

// 节点信息请求(数据报)
type NodeUpinfo_Request struct {
	Node NodeStatView     // 设备节点(必选项)
	Sw   []SwitchStatView // 开关状态列表(可选项)
	Dev  []DeviceStatView // 设备状态列表(可选项)
	Pe   PowerMeterView   // 计量信息(可选项)
}

func (p *NodeUpinfo_Request) NodeID() string {
	return p.Node.Eui
}

func (p *NodeUpinfo_Request) ToJson() string {
	data, _ := json.Marshal(p)
	return string(data)
}

// json format
func foo() {
	bytes, _ := json.Marshal(&NodeUpinfo_Request{})
	log.Println(string(bytes))

}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

}
