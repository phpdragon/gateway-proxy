package config

import (
	"github.com/phpdragon/gateway-proxy/internal/components/logger"
	"github.com/phpdragon/gateway-proxy/internal/components/rabbit"
)

var rabbitClient *rabbit.MqClient

func init() {
	config := GetRabbitConfig()
	if config != nil {
		logger.Warn("Rabbit config missing.")
		return
	}

	rabbitClient = rabbit.NewClient(&rabbit.Options{
		Host:     config.Host,
		Password: config.Password,
		User:     config.User,
	})
}

func Rabbit() *rabbit.MqClient {
	return rabbitClient
}
