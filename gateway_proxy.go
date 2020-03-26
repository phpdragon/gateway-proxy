package main

import (
	"./consts"
	ctl "./controllers"
	eureka "./eureka-client"
	"./logic"
	"./utils"
	"fmt"
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

			signal.Notify(ch, syscall.SIGHUP,syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
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

func returnResponse(rw http.ResponseWriter, req *http.Request, data interface{}) {
	origin := req.Header.Get(consts.ORIGIN)
	rw.Header().Set(consts.CACHE_CONTROL, "No-Cache")
	rw.Header().Set(consts.CONTENT_TYPE, "application/json; charset=utf-8")
	rw.Header().Set(consts.PRAGMA, "No-Cache")
	rw.Header().Set(consts.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(consts.ACCESS_CONTROL_ALLOW_ORIGIN, origin)
		rw.Header().Set(consts.ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	}

	_, err := rw.Write(utils.ToJSONStringByte(data))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
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
		returnResponse(writer, request, ctl.ActuatorStatus(port, appName))
	})
	http.HandleFunc(healthCheckUrl, func(writer http.ResponseWriter, request *http.Request) {
		returnResponse(writer, request, ctl.ActuatorHealth())
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// full applications from eureka server
		//apps := client.Applications
		//b, _ := json.Marshal(apps)
		//_, _ = writer.Write(b)

		indexHandler(writer, request, client)
	})

	log.Printf("Listening on port %s", port)
	log.Printf("Open http://localhost:%s in the browser", port)

	// start http server
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func indexHandler(rw http.ResponseWriter, req *http.Request, client *eureka.EurekaClient) {
	if req.URL.Path != "/" {
		http.NotFound(rw, req)
		return
	}

	instance := client.GetNextServerFromEureka("FUNDS-ACCOUNT")
	println("instance: " + utils.ToJSONString(instance))

	response, _ := logic.HandleHttpRequest(rw, req)

	origin := req.Header.Get(consts.ORIGIN)
	rw.Header().Set(consts.CACHE_CONTROL, "No-Cache")
	rw.Header().Set(consts.CONTENT_TYPE, "application/json; charset=utf-8")
	rw.Header().Set(consts.PRAGMA, "No-Cache")
	rw.Header().Set(consts.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(consts.ACCESS_CONTROL_ALLOW_ORIGIN, origin)
		rw.Header().Set(consts.ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
	}

	_, err := rw.Write([]byte(utils.ToJSONString(response)))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
