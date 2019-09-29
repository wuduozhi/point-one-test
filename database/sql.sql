CREATE TABLE `weibos` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	`user_id` bigint(20) DEFAULT NULL,
	`text`  text default null,
	`ats`   text default null,
	`create_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `users` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(64) DEFAULT NULL,
	`create_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `at_user_weibo_refs` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	`at_user_id` bigint(20) DEFAULT NULL,
	`weibo_id` bigint(20) DEFAULT NULL,
	PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;