# gateway_proxy

#### 介绍
转发代理

------------------------------------

#### 软件架构
软件架构说明


------------------------------------

#### 安装教程

1.  DB

```sql
CREATE TABLE `t_route` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `app_id` int(11) DEFAULT '0' COMMENT 'appId',
  `url_path` varchar(75) DEFAULT NULL COMMENT 'URI路径',
  `service_url` varchar(300) DEFAULT NULL COMMENT '服务名',
  `rate_limit` int(11) DEFAULT '10' COMMENT '频率限制每秒次数',
  `timeout` int(11) DEFAULT '10' COMMENT '微服务调用超时时间，秒',
  `status` int(1) DEFAULT NULL COMMENT '1:启用,0:禁用',
  `timestamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `url_path` (`url_path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='cmd 服务名映射表';

INSERT INTO `db_gateway_proxy`.`t_route`(`id`, `app_id`, `url_path`, `service_url`, `rate_limit`, `timeout`, `status`, `timestamp`)
VALUES (1, 100001, '/cfg/testCtl/demoAction', 'http://FILE-PROXY/testCtl/demoAction', 2, 60, 1, '2020-03-26 19:50:59');

```


2.  xxxx
3.  xxxx


------------------------------------


#### 使用说明

1.编译打包
```shell script
[root@localhost gateway_proxy]# yum install -y make
[root@localhost gateway_proxy]# make help

 Choose a command run in gateway_proxy:

 ############################################
  Go项目编译脚本
  参考：https  //studygolang.com/articles/14919?fr=sidebar
 ############################################
  deps          Install missing dependencies.
  build         Compile the binary.
  clean         Clean build files. Runs `go clean` internally.
  package       Package the app
  deploy        Deploy package to server site
```

2.启动服务
> 部署目录
```shell script
gateway_proxy/
├── bin
│   ├── gateway_proxy
│   └── server.sh
├── etc
│   └── app.yaml
├── favicon.ico
├── LICENSE
├── log
├── README.en.md
└── README.md
```

>执行脚本
```shell script
[root@localhost gateway_proxy]# dos2unix bin/*.sh

[root@localhost gateway_proxy]# bin/server.sh
USAGE:bin/server.sh {start|stop|restart|status}

[root@localhost gateway_proxy]# bin/server.sh status
INFO: the app gateway_proxy is running , pid:727 !
```

------------------------------------


#### 压力测试
1. echo '{"a":"b"}' > data.json
2. ab -c 10 -t 60 -T 'application/json' -p data.json http://192.16.1.2:8080/cfg/testCtl/demoAction

------------------------------------

#### TODO
1. 需要添加日志上报功能
2. 抽离功能部分成为代码库、代码优化
3. 测试用例编写


#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 码云特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  码云官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解码云上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是码云最有价值开源项目，是码云综合评定出的优秀开源项目
5.  码云官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  码云封面人物是一档用来展示码云会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
