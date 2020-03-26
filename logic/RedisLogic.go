package logic

import (
	"../client"
	"../core/log"
	"../utils"
	"strconv"
	"strings"
	"time"
)

const APP_NAME = "gateway_proxy:"

func getAccessTotalCacheKey(key string) string {
	return APP_NAME + "access_total:" + key
}

//访问数量增加一次
func GetAccessTotal(routeId int)(int,int64){
	key  := getAccessTotalCacheKey(strconv.Itoa(routeId))
	cache,err  := client.Redis().Get(key).Result()
	if nil != err || 0 == len(cache){
		return 0,utils.GetCurrentTimeMillis()
	}

	val := strings.Split(cache,"|")
	timeMillis, _ := strconv.ParseInt(val[0], 10, 64)
	count, _ := strconv.Atoi(val[1])

	return count,timeMillis
}

//访问数量增加一次
func AccessTotalIncr(routeId int){
	key  := getAccessTotalCacheKey(strconv.Itoa(routeId))
	cache,err  := client.Redis().Get(key).Result()
	if nil != err || 0 == len(cache){
		return
	}

	val := strings.Split(cache,"|")
	count, _ := strconv.Atoi(val[1])
	//
	count ++
	//
	AccessTotalIncrBy(routeId,count)
}

//访问数量增加次数
func AccessTotalIncrBy(routeId int, total int){
	key  := getAccessTotalCacheKey(strconv.Itoa(routeId))
	val := strconv.FormatInt(utils.GetCurrentTimeMillis(),10) + "|" + strconv.Itoa(total)
	err := client.Redis().Set(key, val,time.Minute).Err()
	if nil != err{
		log.Error(err.Error())
	}
}