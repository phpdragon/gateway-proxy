package cross

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/route"
	appLogic "github.com/phpdragon/gateway-proxy/internal/logic/app"
	"github.com/phpdragon/gateway-proxy/internal/mysql/models"
)

var crossDomainConfMap *map[int]map[string]map[string]models.CrossDomain

func CheckDomain(routeConf *models.RouteConf, origin string) bool {
	//系统处理跨域，直接返回真
	if route.CrossModeAllow == routeConf.CrossMode {
		return true
	}

	//是系统配置模式，判断域名是否命中
	if route.CrossModeConfig == routeConf.CrossMode {
		if 0 >= len(origin) {
			return false
		}

		//是否命中域名
		configMap := getActiveDomains(routeConf.Id)
		_, ok := configMap[routeConf.AppId][origin]
		return ok
	}

	//判断全局跨域配置
	return checkGlobalDomain(routeConf.AppId, origin)
}

// checkGlobalDomain 判断全局跨域配置
func checkGlobalDomain(appId string, origin string) bool {
	//如果配置了全局跨域
	appMap := appLogic.GetAppConfMap()

	//全局跨域判断
	appConf, has := appMap[appId]
	if has && route.CrossModeAllow == appConf.CrossMode {
		return true
	}

	//是系统配置模式，判断域名是否命中
	if route.CrossModeConfig == appConf.CrossMode {
		if 0 >= len(origin) {
			return false
		}

		//是否命中域名
		configMap := getActiveDomains(0)
		_, ok := configMap[appId][origin]
		return ok
	}

	//否则默认为不处理跨域，由下游做处理
	return false
}

func getActiveDomains(routeId int) map[string]map[string]models.CrossDomain {
	confMap := getAllActiveDomains()
	data, ok := confMap[routeId]
	if ok {
		return data
	}
	return nil
}

func getAllActiveDomains() map[int]map[string]map[string]models.CrossDomain {
	if nil != crossDomainConfMap {
		return *crossDomainConfMap
	}

	records, err := models.QueryAllActiveDomains()
	if nil != err {
		config.Logger().Errorf("获取应用配置异常, error: %v", err)
		return nil
	}

	var allAppConfMap = make(map[int]map[string]map[string]models.CrossDomain)
	for _, item := range records {
		subAllAppConfMap, ok := allAppConfMap[item.RouteId]
		if !ok {
			subAllAppConfMap = make(map[string]map[string]models.CrossDomain)
			allAppConfMap[item.RouteId] = subAllAppConfMap
		}

		appConfMap, ok2 := subAllAppConfMap[item.AppId]
		if !ok2 {
			appConfMap = make(map[string]models.CrossDomain)
			subAllAppConfMap[item.AppId] = appConfMap
		}

		appConfMap[item.Origin] = item
	}

	crossDomainConfMap = &allAppConfMap
	return allAppConfMap
}
