package route

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/mysql/models"
)

var routeConfMap *map[string]models.RouteConf

func QueryAllActiveRoutes() map[string]models.RouteConf {
	if nil != routeConfMap {
		return *routeConfMap
	}

	records, err := models.QueryAllActiveRouteConfMap()
	if nil != err {
		config.Logger().Errorf("获取路由配置异常, error: %v", err)
		return nil
	}

	var dataMap = make(map[string]models.RouteConf)
	for _, item := range records {
		dataMap[item.UrlPath] = item
	}
	routeConfMap = &dataMap
	return dataMap
}
