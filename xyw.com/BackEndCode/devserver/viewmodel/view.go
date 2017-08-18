package viewmodel

type Login struct {
	Loginuser    string
	Loginpwd     string
	Usersid      int
	Rolestype    int
	Truename     string
	Nickname     string
	Userheadimg  string
	Userphone    string
	Userstate    int
	Usersex      int
	Usermac      string
	Birthday     string
	Token        string
	ThirdPartyId string
	Os           string
}

type RespFilterCurriculum struct {
	Curriculumsid    int
	Curriculumname   string
	Curriculumicon   string
	Curriculumnature string
	Curriculumstype  string
}

type RespFilterClass struct {
	Classesid   int
	Classesname string
	Classesnum  int
	Classesicon string
	Classstate  int
}

type PostCollection struct {
	State       int //课程状态[0:收藏，1:取消]
	Classroomid int
}

type UpVideoFileCollect struct {
	Islive                             int
	Isondomian                         int
	Curriculumclassroomchaptercentreid int
}

type QueryAttentionRecordWhere struct {
	//	Usersid       int
	//	Rolestype     int
	//	Token         string
	Islive        int    //是否查询直播
	Isondemand    int    //是否查询录播
	Curriculumsid int    //课程ID
	Classesid     int    //班级ID//Classesid
	Subjectcode   string //学科分类
	State         int    //1:关注|0:取消关注
}

type QueryLivesWhere struct {
	Islive        int    //是否查询直播
	Isondemand    int    //是否查询录播
	Curriculumsid int    //课程ID
	Classesid     int    //班级ID//Classesid
	Subjectcode   string //学科分类
	Wherestring   string //搜索条件字符串
	Psstate       int    //点到的状态
	Trstate       int    //上课的状态
	PsUserid      int    //用户ID
	Begindate     string //开始时间
	Enddate       string //结束时间
	Isall         int    //1:代表查询二级学科下的所有课程
	Ccccid        int    //课程班级章节中间表Id
	Liveid        int    //视频的播放Id
	Collegeids    string //学院id数组
	Majorids      string //专业Id数组
	Teacherids    string //老师Id数组
	Classesids    string //班级Id数组
}

type QueryAttentionRecordList struct { //课程首页返回的数据结果集
	Curriculumsid   int
	Curriculumname  string
	Curriculumicon  string
	Curriculumstype string
	Chaptercount    int
	Newchapter      int
	Newlivechapter  int
	Subjectcode     string
	PlaySumnum      int
	DownloadSumnum  int
	FollowSumnum    int
	Classesid       int //Classesid
	PsState         int
	Classesname     string
	Majorname       string
	Collegename     string
	Nickname        string
	Lives           []QueryLiveInfos //子章节列表
}

type QueryLiveInfos struct { //获取课程章节详细
	Liveid         int    //视频播放的Id
	Ccccid         int    //课程班级章节中间表ID
	Chaptername    string //章节名称
	Chaptericon    string //章节图标
	Chapterdetails string //章节详情
	Liveinfo       string //直播/录播简介
	Livetitile     string //标题
	Livepath1      string //播放路径1
	Livepath2      string //播放路径2
	Trailerpath    string //预告播放的路径
	Coverimage     string //封面图片路径
	Livestate      int    //播放状态
	Livetype       int    //播放类型
	Whenlong       int    //时长
	Recommendread  string //推荐阅读
	Iscomment      int    //是否允许评论
	Ischeckcomment int    //是否审核评论
	Isdownload     int    //是否允许下载
	Downloadnum    int    //下载次数
	Playnum        int    //播放次数
	State          int    //是否缺课,0:缺课,1:未缺课
	Nickname       string //教师的名称
	Userheadimg    string //教师的头像
	Usersid        int    // 教师Id
	Begindate      string //开始时间
	Classroomsname string //所在教室
}

type GetStudentsinfo struct {
	//	Usersid          int
	//	Rolestype        int
	//	Token            string
	Studentsid       int
	SubjectclassCode string
}

type PostQueryCurriculums struct {
	//	Usersid        int
	//	Rolestype      int
	//	Token          string
	Begindate                          string
	Enddate                            string
	Teacherids                         string //教师id数组
	Collegeids                         string //学院Id数组
	Majorids                           string //专业id数组
	Classesids                         string //班级id数组
	Searhtxt                           string //文字搜索接口
	Curriculumsids                     string //课程Id数组
	Campusid                           int    //校区Id
	Campusids                          string //校区Id数组
	Buildingid                         int    //楼栋Id
	Buildingids                        string //楼栋Id数组
	Floorsid                           int    //楼层Id
	Floorsids                          string //楼层Id数组
	Classroomid                        int    //教室Id
	State                              int    //课程状态[0：未开始，1：进行中，2：已完成,-1:不查询]
	Curriculumclassroomchaptercentreid int    //课程Id
}

type QueryClassroom struct {
	//	Usersid    int
	//	Rolestype  int
	//	Token      string
	Floorsid   int
	Buildingid int
	Campusid   int
	Pageindex  int
}

type GetAverageclassrate struct {
	//	Usersid       int
	//	Rolestype     int
	//	Token         string
	Curriculumsid int    //课程ID
	Classesid     int    //班级Id
	Teacherid     int    //教师Id
	Majorid       int    //专业ID
	Pattern       int    //模式[1:查询学生的平均到课率，2:查询班级的平均到课率，3:某班的近7次]
	Begindate     string //开始时间
	Enddate       string //结束时间
}

type ResponStudentsClassesAvg struct {
	Usersid        int
	Curriculumsid  int
	Curriculumname string
	Curriculumicon string
	Classesnum     int     //已上课数
	Sumcounts      int     //课程下章节数
	Absentnum      int     //旷课数
	AvgAttendance  float32 //出勤率
}

type ViewSubjectclass struct { //
	Subjectcode         string //
	Subjectname         string //
	Countnum            int    //
	Id                  int    // 学科ID
	Superiorsubjectcode string //上级学科代码
}

type CurriculumChapters struct { //课程章节详细信息
	Curriculumsclasscentreid int
	Curriculumsid            int
	Newchapter               int
	Averageclassrate         float32
	Curriculumname           string
	Infos                    []CurriculumChaptersInfo
}

type CurriculumChaptersInfo struct { //课程章节详细信息
	Curriculumclassroomchaptercentreid int
	Chaptersid                         int
	Chaptername                        string
	Plannumber                         int
	Actualnumber                       int
	Toclassrate                        float32
	Begindate                          string
	Enddate                            string
	State                              int
	ClassesId                          int
	Classesname                        string
}

type ResponAverageclassrate struct {
	Cccid            int
	Studentsid       int
	Classesid        int
	Classesname      string
	Classesnum       int
	Classesicon      string
	Classstate       int
	Truename         string
	Userheadimg      string
	Ddate            string  //日期
	Absenteeismnum   int     //旷课数
	Sumnum           int     //课程总数量
	Plannumber       int     //总人数
	Actualnumber     int     //到的人数
	Averageclassrate float32 //平均率
}

type QueryPeoples struct {
	BuildingId     int
	BuildingName   string
	FloorId        int
	FloorName      string
	Sumnumbers     int
	ClassroomId    int
	ClassroomName  string
	Seatsnumbers   int
	ClassroomState int
	Classroomstype string
	Classroomicon  string
	FloorsImage    string
}

type QueryResultClassroom struct {
	Classroomid       int
	Classroomsname    string
	Classroomicon     string
	Seatsnumbers      int
	Sumnumbers        int
	Classroomstype    string
	Classroomstate    int
	Collectionnumbers int
	Notes             string
	Floorname         string
	Buildingname      string
	Campusname        string
	State             int
}

type GetPointtos struct {
	//	Usersid                            int
	//	Rolestype                          int
	//	Token                              string
	Curriculumclassroomchaptercentreid  int
	Curriculumclassroomchaptercentreids string //中控端发生过来查询多班级上课
	Studentsid                          int
	State                               int
}

type ChangeClassState struct {
	Classroomid int
	State       int
	Os          string //请求来源的平台
	Loginuser   string //请求的用户账号
	Ccccid      int    //正在上课的章节Id
	Ccccids     string //正在上课的章节Id
}

type PostUpdatePointtosData struct {
	CcccIds     []int
	StudentsIds []int
	States      []int
}

type Teacher struct { //获取所有的老师
	Usersid   int
	Nickname  string
	Loginuser string
	Collegeid int
	Majorid   int
}

type PointtosUsers struct {
	Usersid     int
	Truename    string
	Userheadimg string
	State       int
}

type GetCurriculumslist struct {
	Curriculumclassroomchaptercentreid int
	Curriculumname                     string
	Curriculumsid                      int
	Chaptersid                         int
	Chaptername                        string
	Begindate                          string
	Enddate                            string
	Classroomsname                     string
	Floorname                          string
	Buildingname                       string
	Campusname                         string
	Classesname                        string
	Majorname                          string
	Collegename                        string
	State                              int
	Nickname                           string
	Plannumber                         int
	Actualnumber                       int
	Toclassrate                        float32
	TeacherId                          int
}

type Getcampus struct { //
	Campusid   int
	Campusname string
	Campusicon string
	Campuscode string
	Campusnums int
}

type Getbuilding struct {
	Buildingid       int
	Campusid         int
	Buildingname     string
	Buildingicon     string
	Buildingcode     string
	Floorsnumber     int
	Classroomsnumber int
}

type Getfloors struct { //
	Floorsid        int
	Buildingid      int
	Floorname       string
	Floorscode      string
	Classroomnumber int
	FloorsImage     string
	Sumnumber       int
	Rooms           []Getclassrooms
}

type QueryClassroomInfo struct { //对接设备查询返回的数据结果
	Classroomid                        int
	Classroomsname                     string                         //教室名称
	Campusname                         string                         //校区名称
	Buildingname                       string                         //楼栋名称
	Floorname                          string                         //楼层名称
	Classroomicon                      string                         //教室图标
	Classroomstate                     int                            //教室当前的使用状态
	Curriculumname                     string                         //课程名称
	Nickname                           string                         //授课老师名称
	Chaptername                        string                         //章节名称
	Classesid                          int                            //班级Id
	Classesname                        string                         //班级名称
	Seatsnumbers                       int                            //教室内座位数
	Sumnumbers                         int                            //教室内人数
	Curriculumclassroomchaptercentreid int                            //当前上课的唯一Id
	Qccci                              []QueryClassroomCurriculumInfo //多个班级上级时返回
}

type QueryClassroomCurriculumInfo struct {
	Curriculumname                     string
	Nickname                           string
	Chaptername                        string
	Classesid                          int
	Classesname                        string
	Curriculumclassroomchaptercentreid int
}

type Getclassrooms struct {
	Classroomid       int
	Floorsid          int
	Classroomsname    string
	Classroomicon     string
	Seatsnumbers      int
	Sumnumbers        int
	Classroomstype    string
	Classroomstate    int
	Collectionnumbers int
	Notes             string
	Classroomscode    string
}

type Umessage struct { //
	Messageid       int
	Title           string
	Details         string
	State           int
	Usersid         int
	Createdate      string
	Megtype         string
	Readdate        string
	GoUrl           string
	GoParameter     string
	MessageImg      string
	MessageProfiles string
}

type Studentsinfo struct {
	Studentsid       int
	Enrollmentyear   int
	Homeaddress      string
	Nowaddress       string
	Classesid        int //班级Id
	Infostate        int //信息状态
	Currentstate     int //当前状态
	Truename         string
	Nickname         string
	Userheadimg      string
	Userphone        string
	Usersex          int
	Birthday         string
	Rolesid          int
	Userstate        int
	Attendance       float32 //出勤率
	Needcoursenum    int     //应上的课程数
	Alreadycoursenum int     //已经上过的课程数
	Absenteeism      int     //旷课数
}

type QueryBasicsetWhere struct {
	Usersid        int    //用户ID
	Campuscode     string //校区代码
	Campusid       int    //校区ID
	Campusids      string //校区ID数组
	Campusname     string //校区名称
	Buildingcode   string //楼栋代码
	Buildingid     int    //楼栋ID
	Buildingids    string //楼栋ID数组
	Buildingname   string //楼栋名称
	Floorscode     string //楼层代码
	Floorsids      string //楼层ID数组
	Floorsid       int    //楼层ID
	Classroomid    int    //教室ID
	Classroomscode string //教室代码
	Collegecode    string //学院代码
	Collegename    string
	Collegeid      int    //学院ID
	Majorid        int    //专业id
	Majorcode      string //专业代码
	Majorname      string
	Classesid      int    //班级id
	Classescode    string //班级代码
	Classesname    string
}

type Enclosure struct { //课程章节附件表
	Enclosureid          int //文件Id
	Ccccid               int
	Enclosurename        string //文件名称
	Enclosuretype        string //文件类型
	Enclosuresize        int    //文件大小以kb为单位
	EnclosureVirtualPath string //文件的虚拟路径
	Createdate           string //文件的创建时间
	Enclosureicon        string //文件的封面图片
}

type QueryClassRoomPeopleInfo struct { //查询教室内人数详情
	Begindate   string
	Enddate     string
	Classroomid int
	Nickname    string
	Usersid     int
	Userheadimg string
	Createdate  string
	Closedate   string
	DateLength  string
	Xy          string
	X           float64
	Y           float64
}
type PeopleCount struct { //查询教室内人数的统计总和
	Sumnumbers int
	Dateymd    string
	Dateh      string
}

type QueryStreamPeoplesWhere struct { //人流分析
	Begindate  string
	Enddate    string
	Campusid   int
	Buildingid int
	Floorsid   int
}

type ListStreamPeoplesAnalysis struct { //人流分析
	Valcount int
	Valname  string
}

type QueryAttendanceAnalysis struct { //管理者查看各种出勤统计分析
	Analysisvalue float32
	Analysisname  string
}

type AttendanceAnalysisWhere struct { //管理者查看各种出勤统计分析 查询条件
	Begindate     string
	Enddate       string
	Majorid       int
	Collegeid     int
	Curriculumsid int
	Dateint       int //时间差
	Gradeint      int //年级差
	Analysistype  int //查询数据类型
}

type QueryUsersWhere struct { //用户数据查询基本条件
	Studentsid       int    //学生Id
	SubjectclassCode string //班级编号
	Id               int    //用户人员ID
	Rolestype        int    //角色Id
	Userstate        int    //用户状态
	Usersex          int    //用户性别
	Usermac          string //用户终端的mac地址
	Birthday         string //用户生日
	Classesid        int    //学生所在班级
	Infostate        int    //学生信息状态
	Currentstate     int    //当前学生状态
	Searhtxt         string //文字搜索接口
	Classesids       string //查询的班级的数组
	Collegeid        int    //学院
	Majorid          int    //专业
	Collegeids       string //学院Id数组
	Majorids         string //专业Id数组
}

type UsersInfoAll struct {
	UsersId          int
	Loginuser        string
	Loginpwd         string
	Rolesid          int
	RoleName         string
	Truename         string
	Nickname         string
	Userheadimg      string
	Userphone        string
	Userstate        int
	Usersex          int
	Usermac          string
	Birthday         string
	StudentsId       int
	Enrollmentyear   int
	Homeaddress      string
	Nowaddress       string
	Classesid        int //Classesid
	Infostate        int
	Currentstate     int
	Attendance       float32 //出勤率
	Needcoursenum    int     //应上的课程数
	Alreadycoursenum int     //已经上过的课程数
	Absenteeism      int     //旷课数
	TeacherId        int
	Collegeid        int    // 学院ID
	CollegeName      string // 学院名称
	Majorid          int    // 专业ID
	MajorName        string // 专业名称
	Id               int
}

type QueryFilterWhere struct { //获取筛选条件页所需的查询条件
	TeacherIds string
}

type QueryCurriculumWhere struct { //系统后台课程模块所属查询添加
	CurriculumsclasscentreId int //课程班级中间主Id
	Subjectcode              string
	Subjectname              string
	Seacrchtxt               string //搜索字段
	Curriculumname           string //课程名称
	Curriculumnature         string //课程性质
	Curriculumstype          string //课程类型
	Curriculumsid            int    //课程Id
	Chaptername              string //章节名称
	TeacherId                int    //教师Id
	Classesid                int    //班级Id
	Begindate                string //查询开始时间
	Enddate                  string //查询结束时间

}

type ViewCurriculums struct { //课程表
	CurriculumsId      int
	Curriculumname     string
	Curriculumicon     string
	Curriculumnature   string
	Curriculumstype    string
	Curriculumsdetails string
	Chaptercount       int
	Averageclassrate   float32
	Subjectcode        string
	Subjectname        string
	Createdate         string
}

type StudentInfoDetail struct {
	UserID      int    // 用户ID
	RoleID      int    // 角色ID
	RoleName    string // 角色名称
	LoginUser   string // 登录用户
	LoginPwd    string // 登录密码
	TrueName    string // 实名
	NickName    string // 昵称
	UserHeadImg string // 用户头像
	UserPhone   string // 手机号码
	UserState   int    // 用户状态
	UserSex     int    // 性别
	Usermac     string // 用户mac
	Birthday    string // 出生日期

	StudentID        int     // 学生ID
	ClasseID         int     // 班级ID
	ClassName        string  // 班级ID
	EnrollmentYear   int     // 入学年份
	HomeAddress      string  // 家庭住址
	NowAddress       string  // 现住地址
	InfoState        int     // 信息状态
	CurrentState     int     // 当前状态
	Attendance       float32 // 出勤率
	NeedCourseNum    int     // 应上课程数
	AlreadyCourseNum int     // 已上课程数
	Absenteeism      int     // 旷课数
}
type GetCurriculumsClassCentreList struct { //后台课程主计划
	CurriculumsclasscentreId int    //课程班级中间主Id
	Curriculumsid            int    //课程Id
	Curriculumname           string //课程名称
	Curriculumicon           string //课程图标
	Curriculumnature         string //课程性质
	Curriculumstype          string //课程类型
	Truename                 string //真实姓名
	TeacherId                int    //Usersid                  int    //教师ID
	Chaptercount             int    //章节总数
	Subjectcode              string //学科代码
	Subjectname              string //学科名称
	Classesname              string //班级名称
	Classesid                int
	Isondemand               int
	Islive                   int
}
type GetCurriculumclassroomchaptercentreList struct { //后台课程子计划
	CurriculumclassroomchaptercentreId int
	Chaptername                        string
	Classesname                        string
	Classesid                          int
	Begindate                          string
	Enddate                            string
	Islive                             int
	Isondomian                         int
	Classroomsname                     string
	Buildingname                       string
	Campusname                         string
	State                              int
	Truename                           string
	Classroomid                        int
	TeacherId                          int
	Chaptersid                         int
}

type SystemModuleFunctionView struct { //系统模块功能表
	Id                 int
	Systemmoduleid     int
	Systemmodulename   string
	Functionname       string
	Functionicon       string
	Functioncode       string
	Functionsurl       string
	Functionsattribute string
	FunctionDescribe   string
}
type PostAddCurriculumsclasscentre struct { //
	Id            int
	Curriculumsid int
	Classesid     int
	TeacherID     int
	Islive        int
	Isondemand    int
}
type PostAddCurriculumclassroomchaptercentre struct { //
	Id                       int
	Curriculumsclasscentreid int
	Chaptersid               int
	TeacherID                int
	Classroomid              int
	Begindate                string
	Enddate                  string
	Islive                   int
	Isondomian               int
}
