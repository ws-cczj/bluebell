--
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `user_id`     varchar(20)                            NOT NULL COMMENT 'id',
    `username`    varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码',
    `avatar`      varchar(32) COLLATE utf8mb4_general_ci NOT NULL COMMENT '头像',
    `email`       varchar(32) COLLATE utf8mb4_general_ci COMMENT '邮箱',
    `gender`      tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';


DROP TABLE IF EXISTS `community`;
CREATE TABLE `community`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT,
    `author_id`      varchar(20)                             NOT NULL COMMENT '创建者id',
    `author_name`    varchar(32) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '创建者用户名',
    `community_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
    `introduction`   varchar(256) COLLATE utf8mb4_general_ci NOT NULL COMMENT '简介',
    `status`         tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态',
    `create_time`    timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_community_name` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='社区表';


DROP TABLE IF EXISTS `post`;
CREATE TABLE `post`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `post_id`      varchar(20)                              NOT NULL COMMENT 'id',
    `title`        varchar(128) COLLATE utf8mb4_general_ci  NOT NULL COMMENT '标题',
    `content`      varchar(2048) COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
    `author_id`    bigint(20) NOT NULL COMMENT '作者id',
    `author_name`  varchar(32) COLLATE utf8mb4_general_ci   NOT NULL COMMENT '作者名称',
    `community_id` int(11) unsigned NOT NULL COMMENT '所属社区',
    `vote_num`     int(11) unsigned NOT NULL DEFAULT 0 COMMENT '最终票数',
    `status`       tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态',
    `create_time`  timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_post_id` (`post_id`),
    KEY            `idx_author_id` (`author_id`),
    KEY            `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='帖子表';

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT COMMENT '评论id',
    `father_id`      varchar(20)           DEFAULT NULL COMMENT '父评论id',
    `post_id`        varchar(20)  NOT NULL COMMENT '帖子id',
    `type`           tinyint(1) NOT NULL COMMENT '评论类型: 对人评论，对帖子评论',
    `author_id`      varchar(20)  NOT NULL COMMENT '评论作者id',
    `author_name`    varchar(32)  NOT NULL COLLATE utf8mb4_general_ci COMMENT '评论作者名称',
    `to_author_id`   varchar(20)           DEFAULT NULL COMMENT '回复评论作者id',
    `to_author_name` varchar(32)           DEFAULT '' COLLATE utf8mb4_general_ci COMMENT '被评论作者名称',
    `content`        varchar(256) NOT NULL COLLATE utf8mb4_general_ci COMMENT '内容',
    `create_time`    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`    timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY              `idx_post_id` (`post_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='评论表';

DROP TABLE IF EXISTS `user_follow`