package auth

import "github.com/phpdragon/gateway-proxy/internal/mysql/models"

// CheckSession 鉴权
func CheckSession(route *models.RouteConf) bool {
	//TODO: 鉴权，可以url参数上鉴权或者报头上鉴权或者加解密鉴权
	return true
}
