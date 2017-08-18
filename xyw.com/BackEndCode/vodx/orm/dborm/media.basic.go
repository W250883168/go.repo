package dborm

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"gopkg.in/gorp.v1"
	"vodx/ioutil/dbutil"
)

/*
CREATE TABLE `media_basic` (
  `media_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `file_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件ID',
  `media_type` enum('unknown','picture','audio','video') NOT NULL DEFAULT 'unknown' COMMENT '媒体类型',
  `media_name` varchar(255) NOT NULL DEFAULT '' COMMENT '媒体内容名称',
  `content_desc` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `creator_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `create_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建日期',
  `play_sum` int(11) NOT NULL DEFAULT '0' COMMENT '播放次数',
  `download_sum` int(11) NOT NULL DEFAULT '0' COMMENT '下载次数',
  `keep_days` int(11) NOT NULL DEFAULT '0' COMMENT '保留天数',
  PRIMARY KEY (`media_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='媒体文件基本信息表';
*/

const (
	MediaBasic_MediaType_Unknown = "unknown"
	MediaBasic_MediaType_Picture = "picture"
	MediaBasic_MediaType_Audio   = "audio"
	MediaBasic_MediaType_Video   = "video"
)

type MediaBasic struct {
	MediaID     int       `db:"media_id"`
	FileID      int       `db:"file_id"`
	CreatorID   int       `db:"creator_id"`
	MediaType   string    `db:"media_type"`
	MediaName   string    `db:"media_name"`
	ContentDesc string    `db:"content_desc"`
	CreateDate  time.Time `db:"create_date"`
	PlaySum     int       `db:"play_sum"`
	DownloadSum int       `db:"download_sum"`
	KeepDays    int       `db:"keep_days"`
}

func MediaBasic_Query(id int, pDBMap *gorp.DbMap) (pMedia *MediaBasic, err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "media_basic").SetKeys(true, "media_id")
	pObj, err := pDBMap.Get(MediaBasic{}, id)

	pMedia, _ = pObj.(*MediaBasic)
	return pMedia, err
}

func (p *MediaBasic) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "media_basic").SetKeys(true, "media_id")
	return pDBMap.Insert(p)
}

func (p *MediaBasic) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "media_basic").SetKeys(true, "media_id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *MediaBasic) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "media_basic").SetKeys(true, "media_id")
	rows, err := pDBMap.Update(p)

	return int(rows), err
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}

	pDBMap := dbutil.GetDBMap()
	pDBMap.AddTableWithName(MediaBasic{}, "media_basic").SetKeys(true, "media_id")
}
