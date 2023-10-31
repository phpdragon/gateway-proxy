package main

import (
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/base"
	"github.com/phpdragon/gateway-proxy/internal/components/logger"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/errorcode"
	"github.com/phpdragon/gateway-proxy/internal/consts/medietype"
	"github.com/phpdragon/gateway-proxy/internal/logic"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	"github.com/phpdragon/gateway-proxy/internal/utils/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var gFaviconIco, _ = os.ReadFile("favicon.ico")

// 初始化方法
func init() {
	initSignalHandle()
}

/**
 *
 */
func initSignalHandle() {
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

func writeJsonResponse(rw http.ResponseWriter, req *http.Request, response interface{}, isJson bool) {
	origin := req.Header.Get(medietype.ORIGIN)
	rw.Header().Set(medietype.CacheControl, "No-Cache")
	rw.Header().Set(medietype.ContentType, "application/json; charset=utf-8")
	rw.Header().Set(medietype.PRAGMA, "No-Cache")
	rw.Header().Set(medietype.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(medietype.AccessControlAllowOrigin, origin)
		rw.Header().Set(medietype.AccessControlAllowCredentials, "true")
	}

	var err error
	var dataBody []byte
	if isJson {
		dataBody, err = json.ToJSONStringByte(response)
		if err != nil {
			logger.Info(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		dataBody = []byte(response.(string))
	}

	_, err = rw.Write(dataBody)
	if err != nil {
		logger.Info(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func shutdown() {
	eurekaClient := config.Eureka()
	if nil != eurekaClient {
		eurekaClient.Shutdown()
	}
}

func eurekaStateHandler() {
	eurekaClient := config.Eureka()
	http.HandleFunc("/actuator/info", func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, eurekaClient.ActuatorStatus(), true)
	})
	http.HandleFunc("/actuator/health", func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, eurekaClient.ActuatorHealth(), true)
	})
}

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	startTime := date.GetCurrentTimeMillis()
	logger.Info("")

	response, err := logic.HandleHttpRequest(req)
	if nil != err {
		logger.Error(err.Error())
		response = base.BuildFail(errorcode.SystemError, err.Error())
	}

	writeJsonResponse(rw, req, response, true)

	//打印方法执行耗时的信息
	endTime := date.GetCurrentTimeMillis()
	printExecTime(startTime, endTime)
}

// 打印方法执行耗时的信息
func printExecTime(startTime int64, endTime int64) {
	diffTime := endTime - startTime
	diffTimeStr := fmt.Sprintf("请求处理结束,耗时: %s ms\n", strconv.FormatInt(diffTime, 10))
	if diffTime > 1000 {
		logger.Warn(diffTimeStr)
	} else {
		logger.Info(diffTimeStr)
	}
}

// 程序入口
func main() {
	appConfig := config.GetAppConfig()

	// http server
	//处理站点图标
	http.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(gFaviconIco)
		if err != nil {
			logger.Info(err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	//请求入口
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		indexHandler(writer, request)
	})
	//监听日志级别设置
	http.HandleFunc("/handle/level", logger.GetAtomicLevel().ServeHTTP)

	//处理eureka的心跳等
	eurekaStateHandler()

	log.Printf("Listening on port %d", appConfig.Server.Port)
	log.Printf("Open http://localhost:%d in the browser", appConfig.Server.Port)

	// start http server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Server.Port), nil); err != nil {
		log.Fatal(err)
	}
}
