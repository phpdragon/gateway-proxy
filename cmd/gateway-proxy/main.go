package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/router"
	"github.com/phpdragon/gateway-proxy/internal/server"
	"net/http"
	"time"
)

var (
	//https://studygolang.com/articles/4490
	debugMode  = flag.Bool("d", false, "Debug mode: true or false")
	configPath = flag.String("c", "configs/app.yaml", "Config path: like this ../configs/app.yaml or absolute path")
)

// 初始化方法
func init() {
	config.InitConf(*configPath, *debugMode)
	config.NewLogger()
	config.NewMySql()
	config.NewRedis()
	config.NewEureka()
	config.NewRabbit()
}

func main() {
	flag.Parse()

	appConfig := config.GetAppConfig()

	httpServer := &server.HttpServer{
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", appConfig.Server.Port),
			Handler: router.Handler(),
		},
	}

	//监听退出信号，实现安全退出
	go httpServer.WaitForExitingSignal(5 * time.Second)

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		config.Logger().Fatal("Http server listen error, ", err.Error())
	}
}
