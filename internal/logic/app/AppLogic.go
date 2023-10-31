package app

import (
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/base"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/errorcode"
	"github.com/phpdragon/gateway-proxy/internal/logic/request"
	"github.com/phpdragon/gateway-proxy/internal/logic/response"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func WatchSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		for sig := range ch {
			fmt.Println("Signal received:", sig)
			switch sig {
			case syscall.SIGHUP:
				println("Receive exit signal, clients instance going to de-register")
				fallthrough
			case syscall.SIGINT:
				println("Receive exit signal, clients instance going to de-register")
				fallthrough
			case syscall.SIGKILL:
				println("Receive exit signal, clients instance going to de-register")
				fallthrough
			case syscall.SIGTERM:
				log.Println("Receive exit signal, clients instance going to de-register")
				shutdown()
				os.Exit(0)
			}
		}
	}()
}

func HandleUrlRoute() {
	//处理eureka的心跳等
	eurekaStateHandler()
	//监听日志级别设置
	http.HandleFunc("/handle/level", config.GetAtomicLevel().ServeHTTP)
	//请求入口
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		indexHandler(writer, request)
	})

}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	startTime := date.GetCurrentTimeMillis()
	config.Logger().Infof("")

	rsp, err := request.HandleHttpRequest(req)
	if nil != err {
		config.Logger().Error(err.Error())
		rsp = base.BuildFail(errorcode.SystemError, err.Error())
	}

	response.WriteJson(rw, req, rsp, true)

	//打印方法执行耗时的信息
	endTime := date.GetCurrentTimeMillis()
	printExecTime(startTime, endTime)
}

func eurekaStateHandler() {
	eurekaClient := config.Eureka()
	http.HandleFunc("/actuator/info", func(writer http.ResponseWriter, request *http.Request) {
		response.WriteJson(writer, request, eurekaClient.ActuatorStatus(), true)
	})
	http.HandleFunc("/actuator/health", func(writer http.ResponseWriter, request *http.Request) {
		response.WriteJson(writer, request, eurekaClient.ActuatorHealth(), true)
	})
	http.HandleFunc("/actuator/**", func(writer http.ResponseWriter, request *http.Request) {
		response.WriteJson(writer, request, new(interface{}), true)
	})
}

// 打印方法执行耗时的信息
func printExecTime(startTime int64, endTime int64) {
	diffTime := endTime - startTime
	diffTimeStr := fmt.Sprintf("请求处理结束,耗时: %s ms\n", strconv.FormatInt(diffTime, 10))
	if diffTime > 1000 {
		config.Logger().Warn(diffTimeStr)
	} else {
		config.Logger().Infof(diffTimeStr)
	}
}

func shutdown() {
	eurekaClient := config.Eureka()
	if nil != eurekaClient {
		eurekaClient.Shutdown()
	}
}
