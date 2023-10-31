package main

import (
	"flag"
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/logic/app"
	"net/http"
	"os"
)

var (
	gFaviconIco, _ = os.ReadFile("favicon.ico")
	//https://studygolang.com/articles/4490
	debugMode  = flag.Bool("d", false, "Debug mode: true or false")
	configPath = flag.String("c", "configs/app.yaml", "Config path: like this ../configs/app.yaml or absolute path")
)

// 初始化方法
func init() {
	flag.Parse()

	app.WatchSignal()

	config.InitConf(*configPath, *debugMode)
	config.NewLogger()
	config.NewMySql()
	config.NewRedis()
	config.NewEureka()
	config.NewRabbit()
}

// 程序入口
func main() {
	appConfig := config.GetAppConfig()

	// http server
	//处理站点图标
	http.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(gFaviconIco)
		if err != nil {
			config.Logger().Infof(err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	//请求入口
	app.HandleUrlRoute()

	config.Logger().Infof("Listening on port %d", appConfig.Server.Port)
	config.Logger().Infof("Open http://localhost:%d in the browser", appConfig.Server.Port)

	// start http server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Server.Port), nil); err != nil {
		config.Logger().Fatal(err.Error())
	}
}
