package logic

import (
	"encoding/json"
	"errors"
	"github.com/phpdragon/gateway-proxy/internal/base"
	"github.com/phpdragon/gateway-proxy/internal/components/logger"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/medietype"
	"github.com/phpdragon/gateway-proxy/internal/models"
	httpUtil "github.com/phpdragon/gateway-proxy/internal/utils/http"
	"io"
	"net/http"
	"strings"
)

func HandleHttpRequest(req *http.Request) (interface{}, error) {
	body, _ := io.ReadAll(req.Body)
	_ = req.Body.Close()

	if 0 == len(body) {
		return nil, errors.New("参数不能为空")
	}

	routeMap, err := models.QueryAllActiveRoutes()
	if nil != err {
		logger.Info(err.Error())
		return nil, err
	}
	if nil == routeMap {
		return nil, errors.New("请开发人员配置转发设置")
	}

	route, ok := routeMap[req.URL.Path]
	if !ok {
		logger.Info(err.Error())
		return nil, errors.New("请开发人员配置转发设置")
	}

	//请求频率检测
	if !CheckAccessRateLimit(route) {
		return nil, errors.New("请求过于频繁，请稍后再试")
	}

	//获取真实的链接
	eurekaClient := config.Eureka()
	httpUrl, err := eurekaClient.GetRealHttpUrl(route.ServiceUrl)
	if nil != err {
		return nil, err
	}

	//调用远程服务
	remoteData, err := callRemoteService(httpUrl, body, int64(route.Timeout))
	if nil != err {
		return nil, err
	}

	//访问数量增加一次
	AccessTotalIncr(route.Id)

	response := base.BuildOK()
	response.Data = remoteData
	return response, nil
}

func callRemoteService(httpUrl string, req []byte, timeout int64) (interface{}, error) {
	return httpUtil.PostByte(httpUrl, req, timeout)
}

func getPostParams(rw http.ResponseWriter, req *http.Request) (base.ApiRequest, error) {
	body, err := io.ReadAll(req.Body)
	if nil != err {
		logger.Info(err.Error())
		return base.ApiRequest{}, err
	}
	err = req.Body.Close()
	if nil != err {
		logger.Info(err.Error())
		return base.ApiRequest{}, err
	}

	contentType := strings.ToLower(rw.Header().Get(medietype.ContentType))

	var requestData = base.ApiRequest{}
	if strings.Contains(contentType, medietype.ApplicationJson) {
		//application/json协议
		err := json.Unmarshal(body, &requestData)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			logger.Info(err.Error())
			return requestData, nil
		}
	} else if strings.Contains(contentType, medietype.ApplicationXWwwFormUrlencoded) {
		//application/x-www-form-urlencoded协议
	} else if strings.Contains(contentType, medietype.MultipartFormData) {
		//multipart/form-data协议
	} else if strings.Contains(contentType, medietype.ApplicationOctetStream) {
		//application/octet-stream协议
	} else {
		//其他文本协议，只当文本处理
	}
	return requestData, nil
}
