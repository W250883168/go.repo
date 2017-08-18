package deviceModel

import (
	"errors"

	"dev.project/BackEndCode/devserver/model/core"
)

//客户端请求数据
type RequestData struct {
	Auth core.BasicsToken  //身份认证
	Page RequestDataOfPage //分页信息
	Para interface{}       //查询参数
}

//客户端请求数据——分页
type RequestDataOfPage struct {
	PageIndex int //当前页(1-index)
	PageSize  int //每页大小(1-index)
}

//服务器响应客户端的数据
type ResultData struct {
	Page PageData    //分页数据
	Data interface{} //具体的数据
}

//服务器响应客户端的数据——分页
type PageData struct {
	PageIndex   int //当前页
	PageSize    int //每页大小
	PageCount   int //总页数
	RecordCount int //记录总数
}

//设备使用时间
type DeviceUseLog struct {
	OnTime  string //开始时间
	OffTime string //结束时间
	UseTime string //使用时长(秒)
}

//设备详细操作时间
type DeviceDetailLog struct {
	OperateTime string //操作时间
	UserCode    string //操作人（代码）
	UserName    string //操作人（姓名）
	CmdName     string //操作名称
	Para        string //操作参数
}

//设备预警信息
type DeviceAlertInfo struct {
	AlertTime        string //预警时间
	AlertDescription string //预警描述
}

//设备故障信息
type DeviceFaultInfo struct {
	Id               string
	HappenTime       string //发生时间
	FaultSummary     string //故障现象
	FaultDescription string //故障描述
	InputUserId      int    //录入用户(id)
	InputUserName    string //录入用户(姓名)
	InputTime        string //录入时间
	IsCanUse         string //故障发生时，该设备是否还能使用
	Status           string //故障记录状态
}

//教室状态数据
type ClassroomStatusData struct {
	BuildingId        int    //楼栋Id
	BuildingName      string //楼栋Name
	FloorId           int    //楼层id
	FloorName         string //楼层名称
	FloorImage        string //楼层图片
	ClassroomId       int    //教室id
	ClassroomName     string //教室名称
	ClassroomState    int    //教室状态
	CollectionNumbers int    //教室实时人数
	HaveStop          int    //停用：0-未停 1-停用
	HaveAlert         int    //预警：0-无 1-有
	HaveOffline       int    //离线：0-无 1-有
	HaveRun           int    //运行: 0-无 1-有
	ChangeTime        string //最后一次状态变更时间
}

//设备所有操作日志
type DeviceOperateLog struct {
	Id            int64
	OperateTime   string //操作时间
	UserCode      string //操作人（代码）
	UserName      string //操作人（姓名）
	OperateName   string //操作名称
	OperateObject string //操作对象
	ObjectName    string //对象名称

}

//设备型号树
type DeviceModelTree struct {
	Id   string
	PId  string
	Name string
	Type string
}

//设备所有预警信息
type DeviceAllAlertInfo struct {
	AlertId          int64  //预警id
	DeviceId         string //设备id
	DeviceCode       string //设备编码
	DeviceName       string //设备名称
	DeviceModel      string //设备型号
	DeviceSite       string //设备位置
	AlertTime        string //预警时间
	AlertDescription string //预警描述
}

//设备所有故障信息
type DeviceAllFaultInfo struct {
	FaultId       string
	DeviceId      string
	DeviceCode    string
	DeviceName    string
	DeviceModel   string
	DeviceSite    string
	HappenTime    string
	FaultSummary  string
	IsCanUse      string
	Status        string
	StatusName    string
	InputUserId   int
	InputUserName string
}

//设备数量
type DeviceQty struct {
	DeviceId    string //设备id
	ModelId     string //设备型号id
	StopFlag    int    //停用标识
	AlertFlag   int    //预警标识
	FaultFlag   int    //故障标识
	OfflineFlag int    //离线标识
}

//设备型号数量
type DeviceModelQty struct {
	ModelId     string
	ModelName   string
	SubModelIds string

	TotalQty   int
	StopQty    int
	AlertQty   int
	FaultQty   int
	OfflineQty int
}

//设备使用时间
type DeviceUseTime struct {
	DeviceId string
	ModelId  string
	UseTime  int64
}

//设备型号使用时间
type DeviceModelUseTime struct {
	ModelId     string
	ModelName   string
	SubModelIds string

	UseTime int64
}

//设备位置使用时间
type DeviceSiteUseTime struct {
	SiteId   string
	SiteName string
	UseTime  int64
}

//故障管理---------------------
//故障
type DeviceFault struct {
	Id                        string //位置id：字符
	DeviceId                  string //设备id：字符
	DeviceName                string //设备名称：字符
	DeviceSite                string //设备位置：字符
	FaultSummary              string //故障现象：字符
	FaultDescription          string //故障描述：字符
	HappenTime                string //发生时间：字符
	IsCanUse                  string //是否可用：字符(0-不可使用 1-可以使用)
	InputUserId               int    //申报人:整型
	InputUserName             string //申报人名称：字符
	InputTime                 string //申报时间：字符
	SubmitTime                string //提交时间：字符
	Status                    string //状态：字符（0-草稿 1-待受理 2-维修中 3-已维修）
	AcceptanceRepairPerson    string //指定维修人：字符
	AcceptanceRepairPersonTel string //维修人电话：字符
	AcceptanceUserId          int    //受理人id：整型
	AcceptanceUserName        string //受理人姓名：字符
	AcceptanceTime            string //受理时间：字符
	RepairPerson              string //维修人：字符
	RepairFinishTime          string //维修完成时间：字符
	RepairDescription         string //维修描述：字符
	RepairIsCanUse            string //设备是否可用：字符（0-不可使用 1-可以使用）
	RepairResult              string //维修结果：字符（1-未修复 2-已修复）
	RepairInputUserId         int    //维修登记人id：整型
	RepairInputUserName       string //维修登记人姓名：字符
	RepairInputTime           string //维修登记时间：字符
	RepairSubmitTime          string //维修提交时间
	RepairFaultType           []RepairFaultType
}

//维修故障类型
type RepairFaultType struct {
	FaultTypeId   string //故障类型id
	FaultTypeName string //故障类型name
}

//故障注册数据
type RequestRegisterFaultData struct {
	Auth core.BasicsToken  //身份认证
	Para RegisterFaultData //故障注册数据
}

//故障注册数据
type RegisterFaultData struct {
	Id               string //故障id：字符，不能为空
	DeviceId         string //设备id：字符，不能为空
	FaultSummary     string //故障现象：字符，不能为空
	FaultDescription string //故障描述：字符，可以为空
	HappenTime       string //发生时间：字符，不能为空
	IsCanUse         string //是否可用：字符，不能为空（"0"/"1"）
	OT               string //操作类型(暂存-save，提交-submit)
}

func (p *RegisterFaultData) DataValidate() {
	if p.IsCanUse != "0" && p.IsCanUse != "1" {
		panic(errors.New("设备是否可以使用的值只能是0或1"))
	}

	if p.OT != "save" && p.OT != "submit" {
		panic(errors.New("操作类型的值只能是save(暂存)和submit(提交)"))
	}
}

//故障表(与数据库表一致）
type DeviceFaultTable struct {
	Id                        string //故障id：字符
	DeviceId                  string //设备id：字符
	FaultSummary              string //故障现象：字符
	FaultDescription          string //故障描述：字符
	HappenTime                string //发生时间：字符
	IsCanUse                  string //是否可用：字符(0-不可使用 1-可以使用)
	InputUserId               int    //申报人:整型
	InputTime                 string //申报时间：字符
	SubmitTime                string //提交时间：字符
	Status                    string //状态：字符（0-草稿 1-待受理 2-维修中 3-已维修）
	AcceptanceRepairPerson    string //指定维修人：字符
	AcceptanceRepairPersonTel string //维修人电话：字符
	AcceptanceUserId          int    //受理人id：整型
	AcceptanceTime            string //受理时间：字符
	RepairPerson              string //维修人：字符
	RepairFinishTime          string //维修完成时间：字符
	RepairDescription         string //维修描述：字符
	RepairIsCanUse            string //设备是否可用：字符（0-不可使用 1-可以使用）
	RepairResult              string //维修结果：字符（1-未修复 2-已修复）
	RepairInputUserId         int    //维修登记人id：整型
	RepairInputTime           string //维修登记时间：字符
	RepairSubmitTime          string //维修提交时间
}

//故障id数据
type RequestFaultIdData struct {
	Auth core.BasicsToken //身份认证
	Para FaultIdData      //故障Id数据
}

//故障id数据
type FaultIdData struct {
	Id string
}

//教室内设备
type ClassroomDevice struct {
	DeviceId   string
	DeviceName string
}

//设备对应型号分类记录
type DeviceFaultType struct {
	FaultTypeId   string
	FaultTypeName string
}

//设备对应型号词条
type DeviceFaultWord struct {
	Name string
}

//故障受理数据
type RequestAcceptanceFaultData struct {
	Auth core.BasicsToken    //身份认证
	Para AcceptanceFaultData //故障受理数据
}

//故障受理数据
type AcceptanceFaultData struct {
	Id              string //故障id
	RepairPerson    string //维修人
	RepairPersonTel string //维修电话
}

//维修注册数据
type RequestRegisterRepairData struct {
	Auth core.BasicsToken   //身份认证
	Para RegisterRepairData //维修注册数据
}

//维修注册数据
type RegisterRepairData struct {
	Id                string              //故障id
	RepairPerson      string              //维修人
	RepairFinishTime  string              //维修完成时间
	RepairDescription string              //维修描述
	RepairIsCanUse    string              //维修后设备是否可用
	RepairResult      string              //维修结果
	FaultType         []RepairFaultTypeId //故障类型
	OT                string              //操作类型 (暂存-save 提交-submit)
}

//维修故障分类
type RepairFaultTypeId struct {
	FaultTypeId string
}

//节点型号
type NodeModel struct {
	Id          string //节点型号Id
	Name        string //节点型号名称
	Description string //节点型号描述
}

//节点型号命令
type NodeModelCMD struct {
	Id             int
	ModelId        string
	CmdCode        string
	CmdName        string
	RequestURI     string
	URIQuery       string
	CmdDescription string
	RequestType    string
	Payload        string
	CloseCmdFlag   string
	OpenCmdFlag    string
	NodeModelName  string
}

//节点型号(请求保存）
type RequestNodeModelData struct {
	Auth core.BasicsToken //身份认证
	Para NodeModel        //节点型号
}

//节点型号命令(请求保存）
type RequestNodeModelCMDData struct {
	Auth core.BasicsToken //身份认证
	Para NodeModelCMD     //节点型号命令
}

//节点
type Node struct {
	Id                 string
	Name               string
	ModelId            string
	Campusid           int
	Buildingid         int
	Floorsid           int
	ClassRoomId        int
	IpType             string
	NodeCoapPort       string
	InRouteMappingPort string
	RouteIp            string
	UploadTime         string
	NodeModelName      string
	Classroomsname     string
	Buildingname       string
	Campusname         string
}

//节点(请求保存）
type RequestNodeData struct {
	Auth core.BasicsToken //身份认证
	Para Node             //节点型号命令
}

//设备类型
type DeviceModel struct {
	Id           string
	PId          string
	Name         string
	Description  string
	Type         float64
	PageFileName string
	ImgFileName  string
	ImgFileName2 string
	IsAlert      string
	MaxUseTime   int
	TypeName     string
}

//设备类型(请求保存）
type RequestDeviceModelData struct {
	Auth core.BasicsToken //身份认证
	Para DeviceModel      //节点型号命令
}

//设备状态命令
type DeviceModelStatusCMD struct {
	Id                     int
	ModelId                string
	ModelName              string
	Payload                string
	StatusName             string
	StatusCode             string
	StatusValueMatchString string
	SwitchStatusFlag       string
	OnValue                string
	OffValue               string
	SeqNo                  int
	SelectValueFlag        string
	IsAlert                string
	AlertWhere             string
	AlertDescription       string
}

//设备状态命令(请求保存）
type RequestDeviceModelStatusCMDData struct {
	Auth core.BasicsToken     //身份认证
	Para DeviceModelStatusCMD //设备状态命令
}

//设备状态编码
type DeviceModelStatusValueCode struct {
	Id              int
	ModelId         string
	StatusCode      string
	StatusName      string
	ModelName       string
	StatusValueCode string
	StatusValueName string
	IsAlert         string
}

//设备状态编码(请求保存）
type RequestDeviceModelStatusValueCodeData struct {
	Auth core.BasicsToken           //身份认证
	Para DeviceModelStatusValueCode //设备状态编码
}

//设备状态编码
type DeviceModelControlCMD struct {
	Id               int
	ModelId          string
	CmdCode          string
	CmdName          string
	RequestURI       string
	URIQuery         string
	CmdDescription   string
	RequestType      string
	Payload          string
	DelayMillisecond int32
	CloseCmdFlag     string
	OpenCmdFlag      string
	ModelIdName      string
}

//设备状态编码(请求保存）
type RequestDeviceModelControlCMDData struct {
	Auth core.BasicsToken      //身份认证
	Para DeviceModelControlCMD //设备状态编码
}

//设备
type Device struct {
	Id                         string
	Name                       string
	Sn                         string
	Code                       string
	Brand                      string
	ModelId                    string
	ModelName                  string
	ClassroomId                int
	Campusname                 string
	Buildingname               string
	Classroomsname             string
	PowerNodeId                string
	PowerSwitchId              string
	Buildingid                 int
	Campusid                   int
	Floorsid                   int
	JoinMethod                 string
	JoinNodeId                 string
	JoinSocketId               string
	NodeSwitchStatus           string
	NodeSwitchStatusUpdateTime string
	DeviceSelfStatus           string
	DeviceSelfStatusUpdateTime string
	IsCanUse                   string
	UseTimeBefore              int32
	UseTimeAfter               int32
	JoinNodeUpdateTime         string
}

//设备
type RequestDeviceData struct {
	Auth core.BasicsToken //身份认证
	Para Device           //设备
}

//设备型号故障分类
type DeviceModelFaultType struct {
	Id        string
	Name      string
	ModelId   string
	ModelName string
}

//设备型号故障分类
type RequestDeviceModelFaultTypeData struct {
	Auth core.BasicsToken     //身份认证
	Para DeviceModelFaultType //设备
}

//设备型号故障现象常用词条
type DeviceModelFaultWord struct {
	Id        int
	Name      string
	ModelId   string
	ModelName string
}

//设备型号故障现象常用词条
type RequestDeviceModelFaultWordData struct {
	Auth core.BasicsToken     //身份认证
	Para DeviceModelFaultWord //设备
}
