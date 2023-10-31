package controllers

import "github.com/phpdragon/gateway-proxy/internal/utils"

type status struct {
	Name   string `json:"name"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

type health struct {
	Status  string  `json:"status"`
	Details Details `json:"details"`
}

type Details struct {
}

func ActuatorStatus(port int, appName string) interface{} {
	appStatus := status{}
	appStatus.Name = appName
	appStatus.Server.Port = utils.Int2Str(port)
	return appStatus
}

func ActuatorHealth() interface{} {
	appHealth := health{}
	appHealth.Status = "UP"
	appHealth.Details = Details{}
	return appHealth
}
