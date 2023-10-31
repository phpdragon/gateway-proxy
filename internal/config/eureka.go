package config

import (
	"github.com/phpdragon/go-eureka-client"
)

var eurekaClient *eureka.Client

func NewEureka() {
	// create eureka clients
	eurekaClient = eureka.NewClientWithLog("configs/app.yaml", Logger())
	eurekaClient.Run()
}

func Eureka() *eureka.Client {
	return eurekaClient
}
