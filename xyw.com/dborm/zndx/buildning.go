package zndx

/*
CREATE TABLE `building` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Campusid` int(11) DEFAULT NULL,
  `Buildingname` varchar(255) DEFAULT NULL,
  `Buildingicon` varchar(255) DEFAULT NULL,
  `Buildingcode` varchar(255) DEFAULT NULL,
  `Floorsnumber` int(11) DEFAULT NULL,
  `Classroomsnumber` int(11) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM AUTO_INCREMENT=54 DEFAULT CHARSET=utf8;
*/
// 楼栋
type Building struct {
	Id               int    //
	Campusid         int    // 校区ID
	Buildingname     string // 楼栋名称
	Buildingicon     string // 楼栋图标
	Buildingcode     string // 楼栋代码
	Floorsnumber     int    // 楼层数
	Classroomsnumber int    // 教室数
}
