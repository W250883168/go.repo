/*
Navicat MySQL Data Transfer

Source Server         : LOCALHOST
Source Server Version : 50634
Source Host           : localhost:3306
Source Database       : media_info

Target Server Type    : MYSQL
Target Server Version : 50634
File Encoding         : 65001

Date: 2017-03-28 08:40:17
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for camera_info
-- ----------------------------
DROP TABLE IF EXISTS `camera_info`;
CREATE TABLE `camera_info` (
  `camera_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `location_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '位置ID',
  `camera_name` varchar(255) NOT NULL DEFAULT '' COMMENT '摄像头名称',
  `camera_ip` varchar(255) NOT NULL DEFAULT '' COMMENT '摄像头IP',
  `camera_port` int(11) NOT NULL DEFAULT '0' COMMENT '摄像头端口',
  `login_account` varchar(255) NOT NULL DEFAULT '' COMMENT '登录账号',
  `login_password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码',
  PRIMARY KEY (`camera_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='摄像头信息表';

-- ----------------------------
-- Records of camera_info
-- ----------------------------
INSERT INTO `camera_info` VALUES ('1', '0', '', '', '0', '', '');

-- ----------------------------
-- Table structure for computer_info
-- ----------------------------
DROP TABLE IF EXISTS `computer_info`;
CREATE TABLE `computer_info` (
  `computer_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `location_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '位置ID',
  `computer_name` varchar(255) NOT NULL DEFAULT '' COMMENT '电脑名称',
  `computer_ip` varchar(255) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `computer_port` int(11) NOT NULL DEFAULT '0' COMMENT '端口号',
  `login_account` varchar(255) NOT NULL DEFAULT '' COMMENT '登录账号',
  `login_password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码 ',
  PRIMARY KEY (`computer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='教学电脑信息表';

-- ----------------------------
-- Records of computer_info
-- ----------------------------

-- ----------------------------
-- Table structure for file_info
-- ----------------------------
DROP TABLE IF EXISTS `file_info`;
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

-- ----------------------------
-- Records of file_info
-- ----------------------------

-- ----------------------------
-- Table structure for file_server
-- ----------------------------
DROP TABLE IF EXISTS `file_server`;
CREATE TABLE `file_server` (
  `server_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `server_name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务器名称',
  `server_address` varchar(255) NOT NULL DEFAULT '' COMMENT '服务器地址',
  `file_root_path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件根目录',
  `web_path` varchar(255) NOT NULL DEFAULT '' COMMENT 'Web路径(如: http://vod.server:8000/vod)',
  PRIMARY KEY (`server_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文件服务器表';

-- ----------------------------
-- Records of file_server
-- ----------------------------

-- ----------------------------
-- Table structure for location_info
-- ----------------------------
DROP TABLE IF EXISTS `location_info`;
CREATE TABLE `location_info` (
  `location_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `city_name` varchar(255) NOT NULL DEFAULT '' COMMENT '城市名',
  `district_name` varchar(255) NOT NULL DEFAULT '' COMMENT '行政区名称',
  `building_name` varchar(255) NOT NULL DEFAULT '' COMMENT '楼栋名称',
  `floor_name` varchar(255) NOT NULL DEFAULT '' COMMENT '楼层名称',
  `room_name` varchar(255) NOT NULL DEFAULT '' COMMENT '房间名称',
  `display_text` varchar(255) NOT NULL DEFAULT '' COMMENT '显示名称',
  PRIMARY KEY (`location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of location_info
-- ----------------------------

-- ----------------------------
-- Table structure for media_attr
-- ----------------------------
DROP TABLE IF EXISTS `media_attr`;
CREATE TABLE `media_attr` (
  `media_id` bigint(20) NOT NULL,
  `audio_format` enum('unknown','wav','amr','mp3','aac','g711') NOT NULL DEFAULT 'unknown',
  `audio_duration` int(11) NOT NULL DEFAULT '0' COMMENT '音频时长(ms)',
  `auido_channels` int(11) NOT NULL DEFAULT '0' COMMENT '声道数',
  `audio_bitrate` int(11) NOT NULL DEFAULT '0' COMMENT '音频码流率(bps)',
  `audio_sampling_rate` int(11) NOT NULL DEFAULT '0' COMMENT '音频采样频率(Hz)',
  `audio_bitdepth` int(11) NOT NULL DEFAULT '0' COMMENT '位深度(bit)',
  `video_format` enum('unknown','h264','h263','mpeg2') NOT NULL DEFAULT 'unknown',
  `video_width` int(11) NOT NULL DEFAULT '0' COMMENT '宽度(pixel)',
  `video_height` int(11) NOT NULL DEFAULT '0' COMMENT '高度(pixel)',
  `video_resolution` varchar(255) NOT NULL DEFAULT '' COMMENT '分辨率(如: 1920x1080p)',
  `video_duration` int(11) NOT NULL DEFAULT '0' COMMENT '视频时长(ms)',
  `video_frame_rate` int(11) NOT NULL DEFAULT '0' COMMENT '视频帧率(FPS)',
  `video_aspect_ratio` varchar(255) NOT NULL DEFAULT '' COMMENT '宽高比',
  `video_bitrate` int(11) NOT NULL DEFAULT '0' COMMENT '码流率(bps)',
  `media_desc` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`media_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of media_attr
-- ----------------------------

-- ----------------------------
-- Table structure for media_basic
-- ----------------------------
DROP TABLE IF EXISTS `media_basic`;
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

-- ----------------------------
-- Records of media_basic
-- ----------------------------

-- ----------------------------
-- Table structure for media_cat
-- ----------------------------
DROP TABLE IF EXISTS `media_cat`;
CREATE TABLE `media_cat` (
  `category_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父ID',
  `category_name` varchar(255) NOT NULL DEFAULT '' COMMENT '类别名称',
  `category_desc` varchar(255) NOT NULL DEFAULT '' COMMENT '类别名称',
  PRIMARY KEY (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='媒体文件类别表';

-- ----------------------------
-- Records of media_cat
-- ----------------------------

-- ----------------------------
-- Table structure for media_live
-- ----------------------------
DROP TABLE IF EXISTS `media_live`;
CREATE TABLE `media_live` (
  `media_id` bigint(20) NOT NULL,
  `live_path` varchar(1024) NOT NULL,
  PRIMARY KEY (`media_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of media_live
-- ----------------------------

-- ----------------------------
-- Table structure for media_vod
-- ----------------------------
DROP TABLE IF EXISTS `media_vod`;
CREATE TABLE `media_vod` (
  `media_id` bigint(20) NOT NULL,
  `vod_path` varchar(1024) NOT NULL,
  PRIMARY KEY (`media_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of media_vod
-- ----------------------------

-- ----------------------------
-- Table structure for operate_log
-- ----------------------------
DROP TABLE IF EXISTS `operate_log`;
CREATE TABLE `operate_log` (
  `log_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `log_type` int(11) NOT NULL DEFAULT '0' COMMENT '0/query; 1/add; 2/delete; 3/update;',
  `log_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP COMMENT '日志日期',
  `log_content` varchar(1024) NOT NULL DEFAULT '' COMMENT '日志内容',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of operate_log
-- ----------------------------
