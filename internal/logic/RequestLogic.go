package logic

import (
	"github.com/phpdragon/gateway-proxy/internal/core/log"
	"github.com/phpdragon/gateway-proxy/internal/models"
	"github.com/phpdragon/gateway-proxy/internal/utils"
	"go.uber.org/zap"
)

// CheckAccessRateLimit 检查访问频率
func CheckAccessRateLimit(route models.Route) bool {
	if 0 == route.RateLimit {
		return true
	}

	count, timeMillis := GetAccessTotal(route.Id)

	diffTime := utils.GetCurrentTimeMillis() - timeMillis
	//间隔时间大于等于一秒,重置
	if diffTime >= 1000 {
		AccessTotalIncrBy(route.Id, 1)
		return true
	}

	//一秒内超过最大次数
	if route.RateLimit <= count {
		log.Error("请求过于频繁,", zap.Any("router", route))
		return false
	}

	return true
}
