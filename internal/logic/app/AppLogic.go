package app

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts"
	"github.com/phpdragon/gateway-proxy/internal/consts/httpheader"
	"github.com/phpdragon/gateway-proxy/internal/mysql/models"
	"github.com/phpdragon/gateway-proxy/internal/utils/cipher"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	httpUtil "github.com/phpdragon/gateway-proxy/internal/utils/http"
	"net/http"
	"strconv"
)

// allActiveApps 所有应用配置
var allActiveApps *map[string]models.Application

func Refresh() {
	allActiveApps = nil
	GetAppConfMap()
}

// CheckAuth 鉴权
func CheckAuth(req *http.Request, route *models.RouteConf) bool {
	appConf := GetAppConf(route.AppId)
	if nil == appConf {
		return false
	}

	//不鉴权
	if consts.AuthModeNone == appConf.AuthMode {
		return true
	}

	var authToken = ""
	var authTimestamp = ""

	if consts.AuthModeHeader == appConf.AuthMode { //通过报头鉴权
		authToken = req.Header.Get(httpheader.XAuthToken)
		authTimestamp = req.Header.Get(httpheader.XAuthTimestamp)
	} else if consts.AuthModelUrl == appConf.AuthMode { //通过URL参数鉴权
		GetParams := httpUtil.ParseGetArgs(req.URL.RawQuery)
		authToken = GetParams[httpheader.XAuthToken]
		authTimestamp = GetParams[httpheader.XAuthTimestamp]
	}

	if len(authToken) <= 0 || len(authTimestamp) <= 0 {
		return false
	}

	//时间戳只能相对系统时间，延迟2两秒或超前2秒
	timestamp, _ := strconv.ParseInt(authTimestamp, 10, 64)
	currentTime := date.GetCurrentTimestamp() - 2
	laterTime := date.GetCurrentTimestamp() + 2
	if timestamp < currentTime || timestamp > laterTime {
		return false
	}

	token := cipher.Md5(authTimestamp + ":" + appConf.AuthCode)
	return authToken == token
}

// CheckOnline 校验应用是否在线
func CheckOnline(appId string) bool {
	return GetAppConf(appId) != nil
}

func GetAppConf(appId string) *models.Application {
	allAppMap := GetAppConfMap()
	app, ok := allAppMap[appId]
	if !ok {
		return nil
	}
	return &app
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
