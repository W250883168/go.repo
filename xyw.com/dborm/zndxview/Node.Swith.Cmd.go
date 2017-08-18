package zndxview

// 节点开关命令信息
type NodeSwitchCmdView struct {
	NodeID              string // zndx.Node.Id
	NodeName            string // 节点名称
	NodeModelID         string // 节点型号ID
	ClassRoomID         int    // 节点所在教室ID
	IPType              string // 使用IP类型(ipv4/ipv6)
	RouterIP            string // 节点连接的路由器IP
	NodeCoapPort        string // CoAP端口号
	InRouterMappingPort string // 节点在路由器上的映射端口
	UploadTime          string // 上报时间

	CmdID          int64 // zndx.NodeModelCmd.Id
	CmdCode        string
	CmdName        string
	RequestURI     string
	URIQuery       string
	CmdDescription string
	RequestType    string
	Payload        string
	CloseCmdFlag   string
	OpenCmdFlag    string
}
