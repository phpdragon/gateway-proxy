module github.com/phpdragon/gateway-proxy

go 1.21.3

require (
	github.com/astaxie/beego v1.12.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.7.1
	github.com/phpdragon/go-eureka-client v0.0.0-20231030062922-70801a3c4ab6
	github.com/streadway/amqp v1.1.0
	go.uber.org/zap v1.26.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)
