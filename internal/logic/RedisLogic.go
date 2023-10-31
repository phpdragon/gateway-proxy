package logic

import (
	"github.com/phpdragon/gateway-proxy/internal/client"
	"github.com/phpdragon/gateway-proxy/internal/core/log"
	"github.com/phpdragon/gateway-proxy/internal/utils"
	"strconv"
	"strings"
	"time"
)

const AppName = "gateway-proxy:"

func getAccessTotalCacheKey(key string) string {
	return AppName + "access_total:" + key
}

// GetAccessTotal 访问数量增加一次
func GetAccessTotal(routeId int) (int, int64) {
	key := getAccessTotalCacheKey(strconv.Itoa(routeId))
	cache, err := client.Redis().Get(key).Result()
	if nil != err || 0 == len(cache) {
		return 0, utils.GetCurrentTimeMillis()
	}

	val := strings.Split(cache, "|")
	timeMillis, _ := strconv.ParseInt(val[0], 10, 64)
	count, _ := strconv.Atoi(val[1])

	return count, timeMillis
}

// AccessTotalIncr 访问数量增加一次
func AccessTotalIncr(routeId int) {
	key := getAccessTotalCacheKey(strconv.Itoa(routeId))
	cache, err := client.Redis().Get(key).Result()
	if nil != err || 0 == len(cache) {
		return
	}

	val := strings.Split(cache, "|")
	count, _ := strconv.Atoi(val[1])
	//
	count++
	//
	AccessTotalIncrBy(routeId, count)
}

// AccessTotalIncrBy 访问数量增加次数
func AccessTotalIncrBy(routeId int, total int) {
	key := getAccessTotalCacheKey(strconv.Itoa(routeId))
	val := strconv.FormatInt(utils.GetCurrentTimeMillis(), 10) + "|" + strconv.Itoa(total)
	err := client.Redis().Set(key, val, 180*time.Second).Err()
	if nil != err {
		log.Error(err.Error())
	}
}
