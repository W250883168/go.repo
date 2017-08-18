package model

import (
	"log"
	"runtime"
)

// 设备基本信息视图
type DeviceBasciInfoView struct {
	DeviceId                   string
	DeviceName                 string
	DeviceImg                  string //关闭时显示的图片
	DeviceImg2                 string //开启时显示的图片
	DevicePage                 string
	PowerNodeId                string
	PowerSwitchId              string
	JoinMethod                 string //node/pjlink
	JoinNodeId                 string
	JoinSocketId               string
	NodeSwitchStatus           string //on/off/offline(无应答时为offline,离线的意思)
	NodeSwitchStatusUpdateTime string
	DeviceSelfStatus           string // on/off(无应答时为off)
	DeviceSelfStatusUpdateTime string
	UseTimeBefore              int64  //上系统前已使用时间 单位：秒
	UseTimeAfter               int64  //上系统后累计使用时间 单位：秒
	IsCanUse                   string //0-不可用(停用) 1-可用
	IsHaveAlert                string //0-没有预警消息 1-有预警消息
	JoinNodeUpdateTime         string //接入设备的节点最后上报时间
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		log.Printf("初始化： %s /%d\n", fun.Name(), line)
	}

}
