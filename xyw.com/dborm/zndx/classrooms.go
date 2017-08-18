package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `classrooms` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Floorsid` int(11) DEFAULT NULL,
  `Classroomsname` varchar(255) DEFAULT NULL,
  `Classroomicon` varchar(255) DEFAULT NULL,
  `Seatsnumbers` int(11) DEFAULT NULL,
  `Sumnumbers` int(11) DEFAULT NULL,
  `Classroomstype` varchar(255) DEFAULT NULL,
  `Classroomstate` int(11) DEFAULT NULL,
  `Collectionnumbers` int(11) DEFAULT NULL,
  `Maxy` double DEFAULT NULL,
  `Miny` double DEFAULT NULL,
  `Maxx` double DEFAULT NULL,
  `Minx` double DEFAULT NULL,
  `Notes` varchar(255) DEFAULT NULL,
  `Classroomscode` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=349 DEFAULT CHARSET=utf8;
*/
// 教室
type Classroom struct {
	Id                int
	Floorsid          int     // 楼层ID
	Classroomsname    string  // 教室名称
	Classroomicon     string  // 教室图标
	Seatsnumbers      int     // 座位数
	Sumnumbers        int     // 当前室内人数
	Classroomstype    string  // 教室属性
	Classroomstate    int     // 使用状态
	Collectionnumbers int     // 收藏次数
	Maxy              float64 //
	Miny              float64 //
	Maxx              float64 //
	Minx              float64 //
	Notes             string  // 备注
	Classroomscode    string  // 教室代码
}

//根据教室id得到教室详细名称（例如：新校区A栋1层101)
func (p *Classroom) GetClassroomDetailName(dbmap *gorp.DbMap) string {
	sql := `
SELECT CONCAT_WS('', c.Campusname, b.Buildingname, f.Floorname, r.Classroomsname) Name
FROM (SELECT Floorsid,Classroomsname FROM classrooms WHERE Id =?) r
	LEFT JOIN floors f ON f.Id = r.Floorsid
	LEFT JOIN building b ON b.Id = f.Buildingid
	LEFT JOIN campus c ON c.Id = b.Campusid
`
	args := []interface{}{p.Id}
	name, err := dbmap.SelectStr(sql, args...)
	xdebug.LogError(err)
	return name
}
