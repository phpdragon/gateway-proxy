package main

import (
	"./consts"
	ctl "./controllers"
	"./core"
	logger "./core/log"
	eureka "./eureka-client"
	"./logic"
	"./utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	//根目录为执行可执行文件命令的所在目录： favicon.ico => $(pwd)/favicon.ico
	gFaviconIco, _ = ioutil.ReadFile("favicon.ico")
)

//初始化方法
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
	go func() {
		for {
			ch := make(chan os.Signal)

			signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
			sig := <-ch
			fmt.Println("Signal received:", sig, " \n")
			switch sig {
			case syscall.SIGHUP:
				println("Receive exit signal, client instance going to de-egister")
				fallthrough
			case syscall.SIGINT:
				println("Receive exit signal, client instance going to de-egister")
				fallthrough
			case syscall.SIGKILL:
				println("Receive exit signal, client instance going to de-egister")
				fallthrough
			case syscall.SIGTERM:
				log.Println("Receive exit signal, client instance going to de-egister")
				os.Exit(0)
			}
		}
	}()
}

func writeJsonResponse(rw http.ResponseWriter, req *http.Request, response interface{}, isJson bool) {
	origin := req.Header.Get(consts.ORIGIN)
	rw.Header().Set(consts.CACHE_CONTROL, "No-Cache")
	rw.Header().Set(consts.CONTENT_TYPE, "application/json; charset=utf-8")
	rw.Header().Set(consts.PRAGMA, "No-Cache")
	rw.Header().Set(consts.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(consts.ACCESS_CONTROL_ALLOW_ORIGIN, origin)
		rw.Header().Set(consts.ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
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

//程序入口
func main() {
	appConfig := core.GetAppConfig()

	var statusPageURL = "/actuator/info"
	var healthCheckUrl = "/actuator/health"

	// create eureka client
	var eurekaClient = eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://172.16.1.155:8761/eureka/",
		App:                   appConfig.AppName,
		Port:                  10000,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
		StatusPageURL:  statusPageURL,
		HealthCheckUrl: healthCheckUrl,
	}) // start eurekaClient, register、heartbeat、refresh
	eurekaClient.Start()

	//监听日志级别设置
	http.HandleFunc("/handle/level", logger.GetAtomicLevel().ServeHTTP)

	// http server
	http.HandleFunc(statusPageURL, func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, ctl.ActuatorStatus(appConfig.Server.Port, appConfig.AppName), true)
	})
	http.HandleFunc(healthCheckUrl, func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, ctl.ActuatorHealth(), true)
	})
	http.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write(gFaviconIco)
		if err != nil {
			logger.Info(err.Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		indexHandler(writer, request, eurekaClient)
	})

	log.Printf("Listening on port %d", appConfig.Server.Port)
	log.Printf("Open http://localhost:%d in the browser", appConfig.Server.Port)

	// start http server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", appConfig.Server.Port), nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(rw http.ResponseWriter, req *http.Request, client *eureka.EurekaClient) {
	startTime := utils.GetCurrentTimeMillis()
	logger.Info("")

	response, err := logic.HandleHttpRequest(req, client)
	if nil != err {
		logger.Error(err.Error())
		response = core.BuildFail(core.SYSTEM_ERROR, err.Error())
	}

	writeJsonResponse(rw, req, response, true)

	//打印方法执行耗时的信息
	endTime := utils.GetCurrentTimeMillis()
	printExecTime(startTime, endTime)
}

//打印方法执行耗时的信息
func printExecTime(startTime int64, endTime int64) {
	diffTime := endTime - startTime
	diffTimeStr := strings.Replace("请求处理结束,耗时: time ms \n\n=========================================>>>", "time", strconv.FormatInt(diffTime, 10), -1)
	if diffTime > 1000 {
		logger.Warn(diffTimeStr)
	} else {
		logger.Info(diffTimeStr)
	}
}
