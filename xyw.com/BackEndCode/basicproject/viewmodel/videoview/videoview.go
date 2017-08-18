package videoview

import (
	"bytes"
	"fmt"
	"runtime"

	"basicproject/commons/xtext"
)

type PageInfo struct {
	PageIndex int // 当前页
	PageSize  int // 每页大小
	RowTotal  int // 总数量
}

// 视频
type VideoInfoView struct {
	ID             int
	CurriculumID   int    // 课程ID
	CurriculumName string // 课程名称
	ChapterID      int    // 章节ID
	ChapterName    string // 章节名称
	TeacherID      int    // 教师ID
	TeacherName    string // 教师名称
	AllowLive      int    // 允许直播否
	AllowVOD       int    // 允许点播否
	VideoPath1     string // 点播视频1
	VideoPath2     string // 点播视频2
	VideoTitle     string // 视频标题
	VideoInfo      string // 视频信息
	LiveState      int    // 直播状态
	LiveType       int    // 播放类型
	VideoDuration  int    // 视频时长(秒)
	LiveBeginTime  string // 直播开始时间
	DownloadNum    int    // 下载次数
	PlayNum        int    // 播放次数
	IsRelease      int    // 是否发布
	AllowComments  int    // 允许评论
	AllowDownload  int    // 允许下载

	CurriculumClassroomChapterID int
}

// 视频详细
type VideoDetailView struct {
	ID              int
	VideoTitle      string           // 视频标题
	VideoInfo       string           // 视频信息
	CoverImage      string           // 封面
	LiveState       int              // 直播状态
	BeginTime       string           // 开始时间
	EndTime         string           // 结束时间
	VideoDuration   int              // 时长(秒)
	RecommendReads  string           // 推荐读物
	VodPath1        string           // 视频路径1
	VodPath2        string           // 视频路径2
	AllowVOD        int              // 允许点播否
	IsRelease       int              // 是否发布点播
	AllowComments   int              // 允许评论
	IsCheckComments int              // 是否审核评论
	AllowDownload   int              // 允许下载
	PlayNum         int              // 播放次数
	DownloadNum     int              // 下载次数
	AttachmentList  []AttachmentView // 附件

	LessonInfo LessonView // 课堂信息
}

type LessonView struct {
	ID             int
	CurriculumID   int    // 课程ID
	CurriculumName string // 课程名称
	ChapterID      int    // 章节ID
	ChapterName    string // 章节名称
	TeacherID      int    // 教师ID
	TeacherName    string // 教师名称
	ClassID        int    // 班级ID
	ClassName      string // 班级名称
	ClassroomID    int    // 教室ID
	ClassroomName  string // 教室名称
	BuildingID     int    // 教学楼ID
	BuildingName   string // 教学楼名称
	CampusID       int    // 校区ID
	CampusName     string // 校区名称
	BeginTime      string // 开始时间
	EndTime        string // 结束时间
}

// 课程附件
type AttachmentView struct {
	ID            int
	EnclosureName string // 附件名称
	EnclosureType string // 附件类型
	EnclosureSize int    // 附件大小
	VirtualPath   string // HTTP路径
	EnclosurePath string // 相对路径
	CreateDate    string // 创建日期
	EnclosuerIcon string // 附件图标
	IsPublish     int    // 是否发布

	CurriculumClassroomChapterID int // 课程班级章节ID
}

type Request_VideoInfoView struct {
	KeyWords string
	PageInfo
}

func (p *Request_VideoInfoView) WhereCondition() string {
	buff := bytes.NewBufferString(" WHERE (TLive.Whenlong > 0) ")
	if !xtext.IsBlank(p.KeyWords) {
		str := percentstr(p.KeyWords)
		txt := fmt.Sprintf(` AND ((TLive.Liveinfo LIKE '%s') OR (TLive.Livetitile LIKE '%s') OR (TLive.Livepath1 LIKE '%s') OR (TLive.Livepath2 LIKE '%s')) `, str, str, str, str)
		buff.WriteString(txt)
	}

	return buff.String()
}

func (p *PageInfo) LimitString() string {
	var limit string
	if p.PageSize > 0 {
		offset := (p.PageIndex - 1) * p.PageSize
		count := p.PageSize
		limit = fmt.Sprintf(`LIMIT %d, %d `, offset, count)
	}

	return limit
}

func percentstr(txt string) (ret string) {
	if len(txt) > 0 {
		ret = "%" + txt + "%"
	}

	return ret
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		fmt.Printf("初始化： %s /%d\n", fun.Name(), line)
	}
}
