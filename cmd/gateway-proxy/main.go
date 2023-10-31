package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/consts"
	"github.com/phpdragon/gateway-proxy/internal/core"
	logger "github.com/phpdragon/gateway-proxy/internal/core/log"
	"github.com/phpdragon/gateway-proxy/internal/logic"
	"github.com/phpdragon/gateway-proxy/internal/utils"
	"github.com/phpdragon/go-eureka-client"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	gFaviconIco, _ = os.ReadFile("favicon.ico")
	gEurekaClient  *eureka.Client
)

// 初始化方法
func init() {
	initSignalHandle()
	initDB()
	iniLog()
}

func initDB() {
	dbConfig := core.GetDatabaseConfig()
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.DbName, dbConfig.Charset)

	// set default database
	if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
		log.Println("Init db failed. err: ", fmt.Sprint(err))
		os.Exit(1)
	}

	log.Println("Init db success. host: ", dataSource)
}

func iniLog() {
	logConfig := core.GetLogConfig()
	logger.InitLog(logConfig.GetLogFilePath())
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
				println("Receive exit signal, client instance going to de-register")
				fallthrough
			case syscall.SIGINT:
				println("Receive exit signal, client instance going to de-register")
				fallthrough
			case syscall.SIGKILL:
				println("Receive exit signal, client instance going to de-register")
				fallthrough
			case syscall.SIGTERM:
				log.Println("Receive exit signal, client instance going to de-register")
				shutdown()
				os.Exit(0)
			}
		}
	}()
}

func writeJsonResponse(rw http.ResponseWriter, req *http.Request, response interface{}, isJson bool) {
	origin := req.Header.Get(consts.ORIGIN)
	rw.Header().Set(consts.CacheControl, "No-Cache")
	rw.Header().Set(consts.ContentType, "application/json; charset=utf-8")
	rw.Header().Set(consts.PRAGMA, "No-Cache")
	rw.Header().Set(consts.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(consts.AccessControlAllowOrigin, origin)
		rw.Header().Set(consts.AccessControlAllowCredentials, "true")
	}

	var err error
	var dataBody []byte
	if isJson {
		dataBody, err = utils.ToJSONStringByte(response)
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
	if nil != gEurekaClient {
		gEurekaClient.Shutdown()
	}
}

func eurekaStateHandler() {
	http.HandleFunc("/actuator/info", func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, gEurekaClient.ActuatorStatus(), true)
	})
	http.HandleFunc("/actuator/health", func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, gEurekaClient.ActuatorHealth(), true)
	})
}

func indexHandler(rw http.ResponseWriter, req *http.Request, client *eureka.Client) {
	startTime := utils.GetCurrentTimeMillis()
	logger.Info("")

	response, err := logic.HandleHttpRequest(req, client)
	if nil != err {
		logger.Error(err.Error())
		response = core.BuildFail(core.SystemError, err.Error())
	}

	writeJsonResponse(rw, req, response, true)

	//打印方法执行耗时的信息
	endTime := utils.GetCurrentTimeMillis()
	printExecTime(startTime, endTime)
}

// 打印方法执行耗时的信息
func printExecTime(startTime int64, endTime int64) {
	diffTime := endTime - startTime
	diffTimeStr := strings.Replace("请求处理结束,耗时: time ms \n\n=========================================>>>", "time", strconv.FormatInt(diffTime, 10), -1)
	if diffTime > 1000 {
		logger.Warn(diffTimeStr)
	} else {
		logger.Info(diffTimeStr)
	}
}

// 程序入口
func main() {
	appConfig := core.GetAppConfig()

	// create eureka client
	gEurekaClient = eureka.NewClientWithLog("configs/app.yaml", logger.GetLogger())
	gEurekaClient.Run()

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
		indexHandler(writer, request, gEurekaClient)
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
