package app

import "github.com/phpdragon/gateway-proxy/internal/mysql/models/application"

// CheckAppIsOnline 校验应用是否在线
func CheckAppIsOnline(appId string) bool {
	return application.IsExist(appId)
}
