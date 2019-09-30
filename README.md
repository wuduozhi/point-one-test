# Point One 小测验

简单实现微博好友推荐


## 设计

### 数据库设计

* Users 表

用户信息实体表

```
CREATE TABLE `users` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	`name` varchar(64) DEFAULT NULL,
	`create_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

* Weibos 表

微博信息实体表

```
CREATE TABLE `weibos` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	`user_id` bigint(20) DEFAULT NULL,
	`text`  text default null,
	`ats`   text default null,
	`create_at` datetime DEFAULT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

* at_user_weibo_refs 表

关系表。at_user_id 为一条微博中，被 at 的 user 的 ID，user_id 为发这条微博的用户 ID，weibo_id 为这条微博的 ID。

这个表在查找某一个用户的好友推荐的时候，起到类似倒排索引的作用，步骤如下：

1. 通过 at_user_id 找到曾经发微博提及此用户的 userIDs 与即 weiboIDs
2. 通过 userIDs ，找到相关的微博
3. 通过相关的微博，找到与 at_user_id 相关的 user


```sql
CREATE TABLE `at_user_weibo_refs` (
	`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
	`at_user_id` bigint(20) DEFAULT NULL,
	`user_id`  bigint(20) DEFAULT NULL,
	`weibo_id` bigint(20) DEFAULT NULL,
	PRIMARY KEY (`id`)
)ENGINE =InnoDB DEFAULT CHARSET=utf8mb4;
```

### 文件目录

```
server.go            -- 项目启动入口，以及一下路由信息
- config        
    - config.yml     -- 项目配置文件
- database           -- db 实体定义和相关 db 操作
    - database.go
    - user.go
    - weibo.go
- service            -- controller
    - handler.go     -- controller 实现
```

### 启动

- Golang version 1.12


### 接口信息

#### 添加微博

* POST http://127.0.0.1:8080/weibo

* request
```
{
	"text":"Hello,I am wuduozhi,How are you?I 澳门",
	"ats":[4],
	"user_id":1
}
```


#### 获取用户推荐列表

* GET http://127.0.0.1:8080/suggest/2

* response
```
{
    "msg": "ok",
    "userIDs": [
        3,
        4
    ]
}
```