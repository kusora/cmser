CREATE DATABASE IF NOT EXISTS test default charset utf8mb4 COLLATE utf8mb4_general_ci;
create user 'test'@localhost identified by 'test';
grant all on test.* to 'test'@'localhost';
flush privileges;


CREATE TABLE `feedback` (
  `feedback_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL COMMENT '反馈发送者的id',
  `feedback` varchar(2048) NOT NULL DEFAULT '',
  `latitude` decimal(13,10) NOT NULL DEFAULT '999.0000000000' COMMENT '纬度',
  `longitude` decimal(13,10) NOT NULL DEFAULT '999.0000000000' COMMENT '经度',
  `service_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '内部客服id，用户发起的反馈默认值0，其余情况由管理后台设置，大于0',
  `related_feedback_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '客服回复内容关联用户上行反馈的id',
  `feedback_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0是用户上行的内容，1是内部客服下行的内容',
  `device_guid` char(32) NOT NULL DEFAULT '' COMMENT '用户设备唯一Id',
  `created_at` datetime NOT NULL,
  `platform` varchar(10) NOT NULL DEFAULT '',
  `app_version` varchar(20) NOT NULL DEFAULT '',
  `device_model` varchar(200) NOT NULL DEFAULT '',
  `user_agent` varchar(200) NOT NULL DEFAULT '',
  `status` int(11) NOT NULL DEFAULT '0',
  `service_name` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`feedback_id`),
  KEY `user_id_idx` (`user_id`),
  KEY `device_guid` (`device_guid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;