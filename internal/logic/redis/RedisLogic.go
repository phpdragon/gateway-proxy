package redis

import (
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	"strconv"
	"strings"
	"time"
)

func getAccessTotalCacheKey(key int) string {
	return fmt.Sprintf("%s:access_total:%d", config.GetAppConfig().AppName, key)
}

// GetAccessTotal 访问数量增加一次
func GetAccessTotal(routeId int) (int, int64) {
	key := getAccessTotalCacheKey(routeId)
	cache, err := config.Redis().Get(key).Result()
	if nil != err || 0 == len(cache) {
		return 0, date.GetCurrentTimeMillis()
	}

	val := strings.Split(cache, "|")
	timeMillis, _ := strconv.ParseInt(val[0], 10, 64)
	count, _ := strconv.Atoi(val[1])

	return count, timeMillis
}

// AccessTotalIncrBy 访问数量增加次数
func AccessTotalIncrBy(routeId int, total int) {
	key := getAccessTotalCacheKey(routeId)
	val := fmt.Sprintf("%d|%d", date.GetCurrentTimeMillis(), total+1)
	err := config.Redis().Set(key, val, 180*time.Second).Err()
	if nil != err {
		config.Logger().Errorf("访问数量增加次数异常：%v", err.Error())
	}
}
