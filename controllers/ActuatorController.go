package controllers

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

func ActuatorStatus(port string, appName string) interface{} {
	appStatus := status{}
	appStatus.Name = appName
	appStatus.Server.Port = port
	return appStatus
}

func ActuatorHealth() interface{} {
	appHealth := health{}
	appHealth.Status = "UP"
	appHealth.Details = Details{}
	return appHealth
}
