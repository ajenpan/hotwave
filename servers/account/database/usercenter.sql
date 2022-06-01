create DATABASE if not exists account default character set = 'utf8mb4';

use account;

CREATE TABLE IF NOT EXISTS `users` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户唯一id',
  `uname` varchar(64) CHARACTER SET utf8mb4 NOT NULL COMMENT '用户名',
  `passwd` varchar(64) NOT NULL DEFAULT '' COMMENT '密码',
  `nickname` varchar(64) CHARACTER SET utf8mb4 NOT NULL DEFAULT '' COMMENT '昵称',
  `roleid` int(11) NOT NULL DEFAULT 1 COMMENT '角色码',
  `gender` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别',
  `avatar` varchar(1024) NOT NULL DEFAULT '' COMMENT '头像',
  `phone` varchar(32) NOT NULL DEFAULT '' COMMENT '电话号码',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '电子邮箱',
  `stat` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态码',
  `create_at` datetime NOT NULL DEFAULT now() COMMENT '创建时间',
  `last_login_at` datetime NOT NULL DEFAULT "1000-01-01 00:00:00" COMMENT '最后登录时间',
  PRIMARY KEY (`uid`),
  UNIQUE KEY `UQE_user_name` (`uname`)
) ENGINE = InnoDB AUTO_INCREMENT = 100000 DEFAULT CHARSET = utf8mb4;

-- CREATE TABLE IF NOT EXISTS `user_login_record` (
--   `id` bigint(20) NOT NULL AUTO_INCREMENT,
--   `uid` bigint(20) NOT NULL,
--   `login_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
--   `login_ip` varchar(32) NOT NULL DEFAULT '',
--   `login_device` varchar(32) NOT NULL DEFAULT '',
--   `login_device_code` varchar(256) NOT NULL DEFAULT '',
--   PRIMARY KEY (`id`),
--   INDEX KEY `login_at` (`login_at`)
-- ) ENGINE = InnoDB AUTO_INCREMENT = 1000 DEFAULT CHARSET = utf8mb4;