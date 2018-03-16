/*
Navicat MySQL Data Transfer

Source Server         : docker
Source Server Version : 50721
Source Host           : 192.168.182.151:3306
Source Database       : simulation

Target Server Type    : MYSQL
Target Server Version : 50721
File Encoding         : 65001

Date: 2018-03-05 20:31:15
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for registry
-- ----------------------------
DROP TABLE IF EXISTS `registry`;
CREATE TABLE `registry` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `ip` varchar(32) NOT NULL,
  `port` int(6) NOT NULL,
  `version` varchar(10) NOT NULL,
  `major` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of registry
-- ----------------------------
INSERT INTO `registry` VALUES ('5', 'major', 'localdockerhub.com', '5000', 'v2', '0');
INSERT INTO `registry` VALUES ('7', 'registry', 'localdockerhub.com', '5001', 'v2', '1');

-- ----------------------------
-- Table structure for right
-- ----------------------------
DROP TABLE IF EXISTS `right`;
CREATE TABLE `right` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rightname` varchar(30) NOT NULL,
  `righturl` varchar(50) NOT NULL,
  `icon` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of right
-- ----------------------------
INSERT INTO `right` VALUES ('1', '资源封装', '/encapsulationview', 'glyphicon glyphicon-plus-sign');
INSERT INTO `right` VALUES ('2', '本地镜像', '/localimage', 'glyphicon glyphicon-modal-window');
INSERT INTO `right` VALUES ('3', '仓库管理', '/registrymanage', 'glyphicon glyphicon-inbox');
INSERT INTO `right` VALUES ('4', '仓库镜像', '/registry', 'glyphicon glyphicon-home');
INSERT INTO `right` VALUES ('5', '仿真开发', '/simulation/toolview', 'glyphicon glyphicon-wrench');
INSERT INTO `right` VALUES ('6', '文件管理', '/filemanage', 'glyphicon glyphicon-folder-close');
INSERT INTO `right` VALUES ('7', '仿真任务', '/simulation/tasksview', 'glyphicon glyphicon-play');
INSERT INTO `right` VALUES ('8', '用户管理', '/usermanage', 'glyphicon glyphicon-user');
INSERT INTO `right` VALUES ('9', '系统设置', '/setting', 'glyphicon glyphicon-cog');

-- ----------------------------
-- Table structure for rightusermap
-- ----------------------------
DROP TABLE IF EXISTS `rightusermap`;
CREATE TABLE `rightusermap` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rid` int(11) NOT NULL,
  `uid` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of rightusermap
-- ----------------------------
INSERT INTO `rightusermap` VALUES ('33', '1', '1');
INSERT INTO `rightusermap` VALUES ('34', '2', '1');
INSERT INTO `rightusermap` VALUES ('35', '3', '1');
INSERT INTO `rightusermap` VALUES ('36', '4', '1');
INSERT INTO `rightusermap` VALUES ('37', '5', '1');
INSERT INTO `rightusermap` VALUES ('38', '6', '1');
INSERT INTO `rightusermap` VALUES ('39', '7', '1');
INSERT INTO `rightusermap` VALUES ('40', '8', '1');
INSERT INTO `rightusermap` VALUES ('41', '9', '1');
INSERT INTO `rightusermap` VALUES ('42', '1', '2');
INSERT INTO `rightusermap` VALUES ('43', '2', '2');
INSERT INTO `rightusermap` VALUES ('44', '5', '2');
INSERT INTO `rightusermap` VALUES ('45', '6', '2');
INSERT INTO `rightusermap` VALUES ('46', '3', '3');
INSERT INTO `rightusermap` VALUES ('47', '4', '3');
INSERT INTO `rightusermap` VALUES ('48', '7', '3');
INSERT INTO `rightusermap` VALUES ('49', '3', '4');
INSERT INTO `rightusermap` VALUES ('50', '7', '4');
INSERT INTO `rightusermap` VALUES ('51', '8', '4');
INSERT INTO `rightusermap` VALUES ('52', '9', '4');

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(30) NOT NULL,
  `password` varchar(30) NOT NULL,
  `status` int(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user` VALUES ('1', 'default', '123456', '1');
INSERT INTO `user` VALUES ('2', 'developer', '123456', '1');
INSERT INTO `user` VALUES ('3', 'Operator', '123456', '1');
INSERT INTO `user` VALUES ('4', 'administrator', '123456', '1');
