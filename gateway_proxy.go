package main

import (
	"./consts"
	ctl "./controllers"
	"./core"
	eureka "./eureka-client"
	"./logic"
	"./utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//var (
//	g_mqaddr = flag.String("mqaddr", "amqp://jyb:root@172.16.1.8:5672/", "mq server addr")
//
//	g_mysqlconnect = flag.String("mysqipaddr", "172.16.1.13:3306", "myssql host")
//	g_redisaddr    = flag.String("redisaddr", "172.16.1.13:6379", "redis mq server addr")
//
//	g_srvport = flag.String("srvport", "19959", "server port")
//
//	g_jmfpathbase = flag.String("jmfpathbase", "com.jyblife.com.bg", "jmf path base")
//	g_jmfapp      = flag.String("jmfapp", "userService", "jmf app name")
//	g_verion      = flag.String("version", "1.0.0", "version")
//	g_owner       = flag.String("owner", "zhonghua.hzh", "server owner")
//	g_group       = flag.String("group", "*", "server group")
//)

//初始化方法
func init() {
	initSignalHandle()
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

	var err interface{}
	var dataBody []byte
	if isJson {
		dataBody, err = utils.ToJSONStringByte(response)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		dataBody = []byte(response.(string))
	}

	_, err = rw.Write(dataBody)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

//程序入口
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	var statusPageURL = "/actuator/info"
	var healthCheckUrl = "/actuator/health"
	var appName = "go-example"

	// create eureka client
	var client = eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://172.16.1.155:8761/eureka/",
		App:                   appName,
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
	}) // start client, register、heartbeat、refresh
	client.Start()

	// http server
	http.HandleFunc(statusPageURL, func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, ctl.ActuatorStatus(port, appName), true)
	})
	http.HandleFunc(healthCheckUrl, func(writer http.ResponseWriter, request *http.Request) {
		writeJsonResponse(writer, request, ctl.ActuatorHealth(), true)
	})
	http.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		//TODO: 缓存起来
		ff, _ := ioutil.ReadFile("favicon.ico")
		_, err := writer.Write(ff)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		indexHandler(writer, request, client)
	})

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)

	// start http server
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(rw http.ResponseWriter, req *http.Request, client *eureka.EurekaClient) {
	response, err := logic.HandleHttpRequest(req, client)
	if nil != err {
		log.Println(err)
		response = core.BuildFail(core.SYSTEM_ERROR, err.Error())
	}

	writeJsonResponse(rw, req, response, true)
}
