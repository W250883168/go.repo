package zndx

import (
	"gopkg.in/gorp.v1"

	"xutils/xdebug"
)

/*
CREATE TABLE `floors` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Buildingid` int(11) DEFAULT NULL,
  `Floorname` varchar(255) DEFAULT NULL,
  `Floorscode` varchar(255) DEFAULT NULL,
  `FloorsImage` varchar(255) DEFAULT NULL COMMENT '楼层平面图',
  `Classroomnumber` int(11) DEFAULT NULL,
  `Maxy` double DEFAULT NULL,
  `Miny` double DEFAULT NULL,
  `Maxx` double DEFAULT NULL,
  `Minx` double DEFAULT NULL,
  `Sumnumber` int(11) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=114 DEFAULT CHARSET=utf8;
*/
type Floor struct {
	Id              int
	Buildingid      int     // 楼栋ID
	Floorname       string  // 楼层名称
	Floorscode      string  // 楼层代码
	FloorsImage     string  // 楼层平面图
	Classroomnumber int     // 教室数
	Maxy            float64 // 最大Y
	Miny            float64 // 最小Y
	Maxx            float64 // 最大X
	Minx            float64 // 最小X
	Sumnumber       int     // 楼层人数
}

//根据楼层id得到楼层详细名称（例如：新校区A栋1层)
func (p *Floor) GetFloorDetailName(dbmap *gorp.DbMap) string {
	sql := `
SELECT CONCAT_WS('', c.Campusname,b.Buildingname,f.Floorname) Name
FROM (SELECT Buildingid, Floorname FROM floors WHERE Id =?) f
	LEFT JOIN building b ON b.Id = f.Buildingid
	LEFT JOIN campus c ON c.Id = b.Campusid
`
	args := []interface{}{p.Id}
	name, err := dbmap.SelectStr(sql, args...)
	xdebug.LogError(err)
	return name
}
