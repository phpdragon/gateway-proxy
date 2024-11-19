package request

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	routeConst "github.com/phpdragon/gateway-proxy/internal/consts/route"
	"github.com/phpdragon/gateway-proxy/internal/logic/app"
	"github.com/phpdragon/gateway-proxy/internal/logic/cross"
	"github.com/phpdragon/gateway-proxy/internal/logic/route"
	httpUtil "github.com/phpdragon/gateway-proxy/internal/utils/http"
	"net/http"
	"strings"
)

func HandleSystemRequest(req *http.Request) {
	action := strings.Replace(req.URL.Path, routeConst.RouterSystem, "", 1)

	param := httpUtil.ParseGetArgs(req.URL.RawQuery)
	if "refresh" == action {
		refreshCache(param)
	}
}

// refreshCache 刷新缓存
func refreshCache(param map[string]string) {
	if routeConst.SysRefreshKey != param["key"] {
		config.Logger().Errorf("刷新系统配置秘钥非法, key: %s", param["key"])
		return
	}
	app.Refresh()
	route.Refresh()
	cross.Refresh()
}
