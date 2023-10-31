package config

import (
	"github.com/phpdragon/gateway-proxy/internal/components/logger"
	"github.com/phpdragon/go-eureka-client"
)

var eurekaClient *eureka.Client

func init() {
	// create eureka clients
	eurekaClient = eureka.NewClientWithLog("configs/app.yaml", logger.GetLogger())
	eurekaClient.Run()
}

func Eureka() *eureka.Client {
	return eurekaClient
}
