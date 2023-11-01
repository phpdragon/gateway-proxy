package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/logic/router"
	"github.com/phpdragon/gateway-proxy/internal/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	//https://studygolang.com/articles/4490
	debugMode  = flag.Bool("d", false, "Debug mode: true or false")
	configPath = flag.String("c", "configs/app.yaml", "Config path: like this ../configs/app.yaml or absolute path")
)

// 初始化方法
func init() {
	flag.Parse()

	config.InitConf(*configPath, *debugMode)
	config.NewLogger()
	config.NewMySql()
	config.NewRedis()
	config.NewEureka()
	config.NewRabbit()
}

func listenSignal2(httpServer *http.Server, idleConnClose chan<- struct{}) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// We received an interrupt signal, shut down.
	if err := httpServer.Shutdown(ctx); err != nil {
		config.Logger().Error("Http server shutdown error", err)
	}

	config.Logger().Info("Http server shutdown！！！")

	close(idleConnClose)
}

func main() {
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
