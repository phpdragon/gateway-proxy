# gateway-proxy

#### 介绍
转发代理

------------------------------------

#### 软件架构
##### 1.软件架构说明

##### 2.代码结构目录
```text
gateway-proxy/
├── cmd  #项目主要的应用程序
│   └── gateway-proxy
│       └── main.go  #程序入口
├── configs  #项目配置文件目录
│   └── app.yaml
├── favicon.ico  #站点图标
├── go.mod
├── go.sum
├── internal  #私有的应用程序代码库
├── LICENSE
├── logs
│   └── nohup.log
├── Makefile  #编译文件
├── README.en.md
├── README.md
└── scripts  #项目脚本
    └── server.sh
```
参考： [Go 项目标准布局](https://zhuanlan.zhihu.com/p/662397116)、[Go 项目目录结构](https://blog.csdn.net/wohu1104/article/details/123209272)、[go项目标准化工程结构解析](https://blog.csdn.net/keenw/article/details/126352773)

------------------------------------

#### 安装教程

1.  DB

```mysql
CREATE DATABASE db_gateway_proxy;

CREATE TABLE `db_gateway_proxy`.`t_route` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `app_id` int(11) DEFAULT '0' COMMENT '应用ID',
    `url_path` varchar(75) DEFAULT '' COMMENT '请求路径，建议:/appName/module/action',
    `service_url` varchar(300) DEFAULT '' COMMENT '下游Url，支持eureka模式和域名、ip端口模式',
    `rate_limit` int(11) DEFAULT '10' COMMENT '频率限制，每秒次数',
    `timeout` int(11) DEFAULT '10' COMMENT '超时时间，单位秒',
    `rsp_mode` int(1) DEFAULT '0' COMMENT '应答模式：0-明文,1-加密',
    `remark` varchar(255) DEFAULT '' COMMENT '请求路径描述',
    `state` int(1) DEFAULT NULL COMMENT '1:启用,0:禁用',
    `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `url_path` (`url_path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='服务路由映射表';

CREATE TABLE  `db_gateway_proxy`.`t_overload` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `app_id` int(11) DEFAULT '0' COMMENT '应用ID',
  `url_path` varchar(75) DEFAULT '' COMMENT '请求路径',
  `limit` int(11) DEFAULT '10' COMMENT '限制次数',
  `interval` int(11) DEFAULT '10' COMMENT '间隔时间，单位秒',
  `remark` varchar(255) DEFAULT '' COMMENT '请求路径描述',
  `state` int(1) DEFAULT NULL COMMENT '1:启用,0:禁用',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `url_path` (`url_path`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='过载配置表';

INSERT INTO `db_gateway_proxy`.`t_route`(`id`, `app_id`, `url_path`, `service_url`, `rate_limit`, `timeout`, `rsp_mode`, `state`, `update_time`, `create_time`)
VALUES (1, 100001, '/fileProxy/testCtl/demoAction', 'http://FILE-PROXY/testCtl/demoAction', 2, 60, 0, 1, '2020-03-26 19:50:59', '2020-03-26 19:50:59');

INSERT INTO `db_gateway_proxy`.`t_overload`(`id`, `app_id`, `url_path`, `limit`, `interval`, `remark`, `state`, `update_time`, `create_time`) 
VALUES (1, 100001, '/fileProxy/testCtl/demoAction', 5, 10, '', 1, '2023-11-02 18:16:53', '2023-11-02 18:16:53');


```


2.  项目初始化

Linux：
```shell script 
yum install make
make deps
```

Windows：
```shell script
# 设置七牛云代理
go env -w GOPROXY=https://goproxy.cn,direct
# 开启module功能
go env -w GO111MODULE=on
# 依赖安装
go mod download
```

3.  xxxx


------------------------------------


#### 使用说明

1.编译打包
```shell script
[root@localhost gateway-proxy]# yum install -y make
[root@localhost gateway-proxy]# make help

 Choose a command run in gateway-proxy:

 ########################################################
 # Go项目编译脚本
 # 参考：https://studygolang.com/articles/14919?fr=sidebar
 ########################################################
  deps          Install missing dependencies.
  build         Compile the binary.
  clean         Clean build files. Runs `go clean` internally.
  package       Package the app
  deploy        Deploy package to server site
```

2.启动服务
> 部署目录
```text
gateway-proxy/
├── configs
│   └── app.yaml
├── favicon.ico
├── gateway-proxy
├── logs
│   └── nohup.log
└── scripts
    └── server.sh
```

>执行脚本
```shell script
[root@localhost gateway-proxy]# dos2unix scripts/*.sh

[root@localhost gateway-proxy]# ./scripts/server.sh
USAGE:scripts/server.sh {start|stop|restart|status}

[root@localhost gateway-proxy]# ./scripts/server.sh status
INFO: the app gateway-proxy is running , pid:727 !
```

------------------------------------


#### 压力测试
1. echo '{"a":"b"}' > data.json
2. ab -c 10 -t 60 -T 'application/json' -p data.json http://127.0.0.1:8080/cfg/testCtl/demoAction

------------------------------------

#### TODO
1. 支持apollo获取配置
2. 需要添加日志上报功能
3. 代码优化
4. 测试用例编写


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
