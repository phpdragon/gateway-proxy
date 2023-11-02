package limit

import (
	"github.com/phpdragon/gateway-proxy/internal/logic/redis"
	"github.com/phpdragon/gateway-proxy/internal/mysql/entity"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
)

// CheckAccessRateLimit 检查访问频率
func CheckAccessRateLimit(route *entity.RouteConf, limit int, intervalTime int64) (int, bool) {
	if 0 >= route.RateLimit {
		return 0, true
	}

	count, timeMillis := redis.GetAccessTotal(route.Id)

	diffTime := date.GetCurrentTimeMillis() - timeMillis
	//间隔时间大于等于一秒,重置
	if diffTime >= intervalTime {
		redis.AccessTotalIncrBy(route.Id, 1)
		return 0, true
	}

	//一秒内超过最大次数
	if limit <= count {
		return count, false
	}

	return count, true
}

// TotalIncr 访问数量增加一次
func TotalIncr(route *entity.RouteConf, access int, overload int) {
	//访问计数增1
	redis.AccessTotalIncrBy(route.Id, access)
}
