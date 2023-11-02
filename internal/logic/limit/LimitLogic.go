package limit

import (
	"github.com/phpdragon/gateway-proxy/internal/logic/redis"
	"github.com/phpdragon/gateway-proxy/internal/models"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
)

// CheckAccessRateLimit 检查访问频率
func CheckAccessRateLimit(route *models.Route) (int, bool) {
	if 0 == route.RateLimit {
		return 0, true
	}

	count, timeMillis := redis.GetAccessTotal(route.Id)

	diffTime := date.GetCurrentTimeMillis() - timeMillis
	//间隔时间大于等于一秒,重置
	if diffTime >= 1000 {
		redis.AccessTotalIncrBy(route.Id, 1)
		return 0, true
	}

	//一秒内超过最大次数
	if route.RateLimit <= count {
		return count, false
	}

	return count, true
}

// AccessTotalIncr 访问数量增加一次
func AccessTotalIncr(route *models.Route, count int) {
	redis.AccessTotalIncrBy(route.Id, count)
}

// CHeckOverload 过载判断
func CHeckOverload(route *models.Route) bool {

	return true
}
