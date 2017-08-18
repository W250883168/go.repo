package coapview

import (
	"encoding/json"
	"time"

	"gopkg.in/gorp.v1"

	"dborm/zndx"
	"xutils/xtime"
)

type CoapResponseView struct {
	Msg string
}

func (p *CoapResponseView) ToJson() string {
	data, _ := json.Marshal(p)
	return string(data)
}

// 节点上报信息视图
type NodeUpinfoView struct {
	NodeUpinfo_Request

	Uptime   string
	NodeAddr string
}

func NodeUpinfoView_GetListFrom(pUpinfo *zndx.NodeUpinfo) (list []NodeUpinfoView, err error) {
	list = []NodeUpinfoView{}
	err = json.Unmarshal([]byte(pUpinfo.UpinfoList), &list)
	return list, err
}

func NodeUpinfoView_GetLastFrom(pUpinfo *zndx.NodeUpinfo) (pView *NodeUpinfoView, err error) {
	pView = &NodeUpinfoView{}
	err = json.Unmarshal([]byte(pUpinfo.LastUpinfo), pView)
	return pView, err
}

func (p *NodeUpinfoView) Onlined(offTimeout time.Duration) (onlined bool) {
	uptime, err := time.Parse(xtime.FormatString(), p.Uptime)
	if err == nil {
		onlined = uptime.Add(offTimeout).After(time.Now())
	}

	return onlined
}

func (p *NodeUpinfoView) Save(pDBMap *gorp.DbMap) (err error) {
	var pInfo *zndx.NodeUpinfo
	if !zndx.NodeUpinfo_Exists(p.NodeID(), pDBMap) {
		pInfo = &zndx.NodeUpinfo{
			NodeID:   p.NodeID(),
			PrevTime: xtime.TimeString(&time.Time{}),
			LastTime: xtime.NowString()}
		bytes, _ := json.Marshal(&p.NodeUpinfo_Request)
		pInfo.LastUpinfo = string(bytes)
		list := []NodeUpinfoView{*p}
		bytes, _ = json.Marshal(&list)
		pInfo.UpinfoList = string(bytes)
		return pInfo.Save(pDBMap)
	}

	pInfo, err = zndx.NodeUpinfo_Get(p.NodeID(), pDBMap)
	if (err == nil) && (pInfo != nil) {
		pInfo.PrevTime = pInfo.LastTime
		pInfo.LastTime = xtime.NowString()
		bytes, _ := json.Marshal(&p.NodeUpinfo_Request)
		pInfo.LastUpinfo = string(bytes)

		list := []NodeUpinfoView{}
		if err = json.Unmarshal([]byte(pInfo.UpinfoList), &list); err != nil {
			return err
		}

		const Max_List_Length = 10 // 最大列表大小
		list = append(list, *p)
		if length := len(list); length > Max_List_Length {
			list = list[length-Max_List_Length:]
		}

		data, _ := json.Marshal(&list)
		pInfo.UpinfoList = string(data)
		err = pInfo.Save(pDBMap)
	}

	return err
}
