package limit

import (
	"github.com/phpdragon/gateway-proxy/internal/logic/redis"
	"github.com/phpdragon/gateway-proxy/internal/mysql/entity"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
)

// CheckAccessRateLimit 检查访问频率
func CheckAccessRateLimit(routeConf *entity.RouteConf) (int, bool) {
	cacheKey := redis.GetAccessTotalCacheKey(routeConf.Id, redis.CacheKeyRateLimit)
	return checkLimit(cacheKey, routeConf.RateLimit, 1)
}

// CheckOverloadLimit 过载保护
func CheckOverloadLimit(routeConf *entity.RouteConf) (int, bool) {
	cacheKey := redis.GetAccessTotalCacheKey(routeConf.Id, redis.CacheKeyOverload)
	return checkLimit(cacheKey, routeConf.Limit, routeConf.Interval)
}

// TotalIncr 访问、过载计数增加一次
func TotalIncr(routeConf *entity.RouteConf, accessTotal int, overload int) {
	accessTotalKey := redis.GetAccessTotalCacheKey(routeConf.Id, redis.CacheKeyRateLimit)
	overloadKey := redis.GetAccessTotalCacheKey(routeConf.Id, redis.CacheKeyOverload)
	//访问计数增1
	redis.AccessTotalIncrBy(accessTotalKey, accessTotal, 5)
	//过载计数增1
	redis.AccessTotalIncrBy(overloadKey, overload, routeConf.Interval+5)
}

// checkLimit 检查访问频率
func checkLimit(cacheKey string, limit int, intervalTime int) (int, bool) {
	if 0 >= limit {
		return 0, true
	}

	count, timeMillis := redis.GetAccessTotalAndTimeMillis(cacheKey)

	diffTime := date.GetCurrentTimeMillis() - timeMillis
	intervalMillis := int64(intervalTime * 1000)
	//当前区间大于等于间隔时间,重置计数
	if diffTime >= intervalMillis {
		redis.AccessTotalIncrBy(cacheKey, 1, intervalTime)
		return 0, true
	}

	//间隔时间内访问计数超过最大限制次数
	if limit <= count {
		return count, false
	}

	return count, true
}
