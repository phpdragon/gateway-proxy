package logic

import (
	"../consts"
	"../core"
	log "../core/log"
	eureka "gitee.com/go-eurake-client"
	"../models"
	"../utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func HandleHttpRequest(req *http.Request, eurekaClient *eureka.EurekaClient) (interface{}, error) {
	body, _ := ioutil.ReadAll(req.Body)
	_ = req.Body.Close()

	if 0 == len(body) {
		return nil, errors.New("参数不能为空")
	}

	routeMap, err := models.QueryAllActiveRoutes()
	if nil != err {
		log.Info(err.Error())
		return nil, err
	}
	if nil == routeMap {
		return nil, errors.New("请开发人员配置转发设置")
	}

	route, ok := routeMap[req.URL.Path]
	if !ok {
		log.Info(err.Error())
		return nil, errors.New("请开发人员配置转发设置")
	}

	//请求频率检测
	if !CheckAccessRateLimit(route) {
		return nil, errors.New("请求过于频繁，请稍后再试")
	}

	//获取真实的链接
	httpUrl,err := eurekaClient.GetRealHttpUrl(route.ServiceUrl)
	if nil != err {
		return nil, err
	}

	//调用远程服务
	remoteData, err := callRemoteService(httpUrl, body)
	if nil != err {
		return nil, err
	}

	//访问数量增加一次
	AccessTotalIncr(route.Id)

	response := core.BuildOK()
	response.Data = remoteData
	return response, nil
}

func callRemoteService(httpUrl string, req []byte) (interface{}, error) {
	return utils.HttpPostByte(httpUrl, req)
}

func getPostParams(rw http.ResponseWriter, req *http.Request) (core.ApiRequest, error) {
	body, err := ioutil.ReadAll(req.Body)
	if nil != err {
		log.Info(err.Error())
		return core.ApiRequest{}, err
	}
	err = req.Body.Close()
	if nil != err {
		log.Info(err.Error())
		return core.ApiRequest{}, err
	}

	contentType := strings.ToLower(rw.Header().Get(consts.CONTENT_TYPE))

	var requestData = core.ApiRequest{}
	if strings.Contains(contentType, consts.APPLICATION_JSON) {
		//application/json协议
		err := json.Unmarshal([]byte(body), &requestData)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			log.Info(err.Error())
			return requestData, nil
		}
	} else if strings.Contains(contentType, consts.APPLICATION_X_WWW_FORM_URLENCODED) {
		//application/x-www-form-urlencoded协议
	} else if strings.Contains(contentType, consts.MULTIPART_FORM_DATA) {
		//multipart/form-data协议
	} else if strings.Contains(contentType, consts.APPLICATION_OCTET_STREAM) {
		//application/octet-stream协议
	} else {
		//其他文本协议，只当文本处理
	}
	return requestData, nil
}
