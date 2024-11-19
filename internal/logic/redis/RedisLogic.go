package redis

import (
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	"strconv"
	"strings"
	"time"
)

const (
	CacheKeyRateLimit = "rateLimit"
	CacheKeyOverload  = "overload"
)

func GetAccessTotalCacheKey(routeId int, code string) string {
	return fmt.Sprintf("%s:%s:%d", config.GetAppConfig().AppName, code, routeId)
}

// GetAccessTotalAndTimeMillis 获取访问总数和计数时间
func GetAccessTotalAndTimeMillis(cacheKey string) (int, int64) {
	cache, err := config.Redis().Get(cacheKey).Result()
	if nil != err || 0 == len(cache) {
		return 0, date.GetCurrentTimeMillis()
	}

	val := strings.Split(cache, "|")
	timeMillis, _ := strconv.ParseInt(val[0], 10, 64)
	count, _ := strconv.Atoi(val[1])

	return count, timeMillis
}

// AccessTotalIncrBy 访问数量增加次数
func AccessTotalIncrBy(key string, total int, expiration int) {
	val := fmt.Sprintf("%d|%d", date.GetCurrentTimeMillis(), total+1)
	err := config.Redis().Set(key, val, time.Second*time.Duration(expiration)).Err()
	if nil != err {
		config.Logger().Errorf("访问数量增加次数异常：%v", err.Error())
	}
}
