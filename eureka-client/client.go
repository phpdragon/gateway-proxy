package eureka_client

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// EurekaClient客户端
type EurekaClient struct {
	// for monitor system signal
	signalChan chan os.Signal
	mutex      sync.RWMutex
	Running    bool
	Config     *Config
	// eureka服务中注册的应用
	Applications *Applications
	//TODO:增器
	autoInc *AutoInc
	//
	cache map[string]interface{}
}

// Start 启动时注册客户端，并后台刷新服务列表，以及心跳
func (client *EurekaClient) Start() {
	client.mutex.Lock()
	client.Running = true
	client.mutex.Unlock()

	//TODO: 设值自增器
	client.autoInc = NewAutoInc(0,1)

	// 注册
	if err := client.doRegister(); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("Register application instance successful")
	// 刷新服务列表
	go client.refresh()
	// 心跳
	go client.heartbeat()
	// 监听退出信号，自动删除注册信息
	go client.handleSignal()
}

// refresh 刷新服务列表
func (client *EurekaClient) refresh() {
	for {
		if client.Running {
			if err := client.doRefresh(); err != nil {
				log.Println(err)
			} else {
				log.Println("Refresh application instance successful")
			}
		} else {
			break
		}
		sleep := time.Duration(client.Config.RegistryFetchIntervalSeconds)
		time.Sleep(sleep * time.Second)
	}
}

// heartbeat 心跳
func (client *EurekaClient) heartbeat() {
	for {
		if client.Running {
			if err := client.doHeartbeat(); err != nil {
				log.Println(err)
			} else {
				log.Println("Heartbeat application instance successful")
			}
		} else {
			break
		}
		sleep := time.Duration(client.Config.RenewalIntervalInSecs)
		time.Sleep(sleep * time.Second)
	}
}

func (client *EurekaClient) doRegister() error {
	instance := client.Config.instance
	return Register(client.Config.DefaultZone, client.Config.App, instance)
}

func (client *EurekaClient) doUnRegister() error {
	instance := client.Config.instance
	return UnRegister(client.Config.DefaultZone, instance.App, instance.InstanceID)
}

func (client *EurekaClient) doHeartbeat() error {
	instance := client.Config.instance
	return Heartbeat(client.Config.DefaultZone, instance.App, instance.InstanceID)
}

func (client *EurekaClient) doRefresh() error {
	// TODO: If the delta is disabled or if it is the first time, get all applications

	// get all applications
	applications, err := Refresh(client.Config.DefaultZone)
	if err != nil {
		return err
	}

	// set applications
	client.mutex.Lock()
	client.Applications = applications
	client.cache = nil
	client.mutex.Unlock()
	return nil
}

// handleSignal 监听退出信号，删除注册的实例
func (client *EurekaClient) handleSignal() {
	if client.signalChan == nil {
		client.signalChan = make(chan os.Signal)
	}
	signal.Notify(client.signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	for {
		switch <-client.signalChan {
		case syscall.SIGINT:
			println("Receive exit signal, client instance going to de-egisterdd")
			fallthrough
		case syscall.SIGKILL:
			println("Receive exit signal, client instance going to de-egister123123")
			fallthrough
		case syscall.SIGTERM:
			log.Println("Receive exit signal, client instance going to de-egister")
			err := client.doUnRegister()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("UnRegister application instance successful")
			}
			os.Exit(0)
		}
	}
}

// NewClient 创建客户端
func NewClient(config *Config) *EurekaClient {
	defaultConfig(config)
	config.instance = NewInstance(getLocalIP(), config)
	return &EurekaClient{Config: config}
}

func defaultConfig(config *Config) {
	if config.DefaultZone == "" {
		config.DefaultZone = "http://localhost:8761/eureka/"
	}
	if config.RenewalIntervalInSecs == 0 {
		config.RenewalIntervalInSecs = 30
	}
	if config.RegistryFetchIntervalSeconds == 0 {
		config.RegistryFetchIntervalSeconds = 15
	}
	if config.DurationInSecs == 0 {
		config.DurationInSecs = 90
	}
	if config.App == "" {
		config.App = "server"
	} else {
		config.App = strings.ToLower(config.App)
	}
	if config.Port == 0 {
		config.Port = 80
	}
}

func (client *EurekaClient) GetNextServerFromEureka(appName string) Instance {
	apps := client.Applications.Applications

	if nil == client.cache {
		appMap := make(map[string]interface{})
		for _, app := range apps {
			appName := app.Name
			appInstances := app.Instances

			instanceMap := make(map[int]interface{})
			for key, instance := range appInstances {
				instanceMap[key] = instance
			}
			appMap[appName] = instanceMap
		}
		client.cache = appMap
	}

	appList := client.cache[appName].(map[int]interface {})
	var incrementAndGet = client.autoInc.IncrementAndGet()
	var index = incrementAndGet % len(appList)
	return appList[index].(Instance)
}

//TODO: 优化 比如https
func (client *EurekaClient) GetRealHttpUrl(httpUrl string) string {
	httpUrlTmp := strings.Replace(httpUrl, "http://", "", -1)
	httpUrlTmp = strings.Replace(httpUrlTmp, "https://", "", -1)
	urls := strings.Split(httpUrlTmp, "/")
	appName := urls[0]

	instance := client.GetNextServerFromEureka(appName)
	realIpPort := instance.InstanceID

	return strings.Replace(httpUrl, appName, realIpPort, -1)
}