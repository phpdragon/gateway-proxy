package app

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/mysql/models"
)

// allActiveApps 所有应用配置
var allActiveApps *map[string]models.Application

func Refresh() {
	allActiveApps = nil
	GetAppConfMap()
}

// CheckAppIsOnline 校验应用是否在线
func CheckAppIsOnline(appId string) bool {
	configMap := GetAppConfMap()
	_, ok := configMap[appId]
	return ok
}

func GetAppConfMap() map[string]models.Application {
	if nil != allActiveApps {
		return *allActiveApps
	}

	records, err := models.QueryAllActiveApp()
	if nil != err {
		config.Logger().Errorf("获取应用配置异常, error: %v", err)
		return nil
	}

	var dataMap = make(map[string]models.Application)
	for _, item := range records {
		dataMap[item.AppId] = item
	}

	allActiveApps = &dataMap
	return dataMap
}
