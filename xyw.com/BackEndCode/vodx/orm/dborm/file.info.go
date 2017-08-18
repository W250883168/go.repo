package dborm

import (
	"fmt"
	"log"
	"runtime"

	"gopkg.in/gorp.v1"
	"vodx/ioutil/dbutil"
)

/*
CREATE TABLE `file_info` (
  `file_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `server_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '服务器ID',
  `file_name` varchar(255) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件路径(相对)',
  `file_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小',
  `check_sum` varchar(255) NOT NULL DEFAULT '' COMMENT '校验值',
  `check_mode` enum('none','crc32','md5','sha1') NOT NULL DEFAULT 'none' COMMENT '校验方式',
  PRIMARY KEY (`file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文件基本信息表';
*/

const (
	FileInfo_CheckMode_None  = "none"
	FileInfo_CheckMode_CRC32 = "crc32"
	FileInfo_CheckMode_MD5   = "md5"
	FileInfo_CheckMode_SHA1  = "sha1"
)

type FileInfo struct {
	FileID    int    `db:"file_id"`
	ServerID  int    `db:"server_id"`
	FileName  string `db:"file_name"`
	FilePath  string `db:"file_path"`
	FileSize  int    `db:"file_size"`
	CheckSum  string `db:"check_sum"`
	CheckMode string `db:"check_mode"`
}

func FileInfo_Query(id int, pDBMap *gorp.DbMap) (pInfo *FileInfo, err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "file_info").SetKeys(true, "file_id")
	pObj, err := pDBMap.Get(FileInfo{}, id)

	pInfo, _ = pObj.(*FileInfo)
	return pInfo, err
}

func (p *FileInfo) Insert(pDBMap *gorp.DbMap) (err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "file_info").SetKeys(true, "file_id")
	return pDBMap.Insert(p)
}

func (p *FileInfo) Delete(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "file_info").SetKeys(true, "file_id")
	rows, err := pDBMap.Delete(p)

	return int(rows), err
}

func (p *FileInfo) Update(pDBMap *gorp.DbMap) (affects int, err error) {
	pDBMap.AddTableWithName(MediaBasic{}, "file_info").SetKeys(true, "file_id")
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
	pDBMap.AddTableWithName(MediaBasic{}, "file_info").SetKeys(true, "file_id")
}
