/*
Navicat MySQL Data Transfer

Source Server         : MySQL.Server
Source Server Version : 50634
Source Host           : 192.168.0.205:3306
Source Database       : zndx2

Target Server Type    : MYSQL
Target Server Version : 50634
File Encoding         : 65001

Date: 2017-04-21 11:04:13
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for device_props
-- ----------------------------
DROP TABLE IF EXISTS `device_props`;
CREATE TABLE `device_props` (
  `device_id` varchar(50) NOT NULL,
  `k` varchar(255) NOT NULL DEFAULT '' COMMENT '键名称',
  `v` varchar(1024) NOT NULL DEFAULT '' COMMENT '键值',
  `comments` varchar(255) NOT NULL,
  PRIMARY KEY (`device_id`,`k`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备属性键值对表, 如: \r\ndevice.accumulate.time.seconds=300\r\ndevice.state.online=true/false\r\ndevice.used.time.before.seconds=500\r\ndevice.used.time.after.seconds=500\r\n\r\n\r\n\r\n';

-- ----------------------------
-- Records of device_props
-- ----------------------------
INSERT INTO `device_props` VALUES ('f28c5c7f-7dc4-4a8a-8e83-4b1a5e469b25', 'device.accumulate.time.seconds', '1220', '设备累计使用时间');
