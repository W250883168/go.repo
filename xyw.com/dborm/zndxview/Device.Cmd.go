package zndxview

// 设备命令
type DeviceCmd struct {
	ModelName          string
	DeviceName         string
	PowerNodeId        string
	PowerSwitchId      string
	JoinMethod         string
	JoinNodeId         string
	JoinSocketId       string
	CmdCode            string
	CmdName            string
	RequestURI         string
	URIQuery           string
	RequestType        string
	Payload            string
	IpType             string
	InRouteMappingPort string
	NodeCoapPort       string
	RouteIp            string
}
