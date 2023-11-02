package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/phpdragon/gateway-proxy/internal/base"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/httpheader"
	"github.com/phpdragon/gateway-proxy/internal/consts/medietype"
	"github.com/phpdragon/gateway-proxy/internal/consts/rspmode"
	"github.com/phpdragon/gateway-proxy/internal/logic/auth"
	"github.com/phpdragon/gateway-proxy/internal/logic/limit"
	"github.com/phpdragon/gateway-proxy/internal/mysql/dao"
	"github.com/phpdragon/gateway-proxy/internal/utils/net"
	"io"
	"net/http"
	"strings"
	"time"
)

func HandleHttpRequest(req *http.Request) ([]byte, http.Header, error) {
	body, _ := io.ReadAll(req.Body)
	_ = req.Body.Close()

	if 0 == len(body) {
		return nil, nil, errors.New("参数不能为空")
	}

	routeConfMap, err := dao.QueryAllActiveRoutes()
	if nil != err {
		config.Logger().Errorf("系统无法路由当前请求,请联系开发人员进行配置, urlPath: %s, error: %v", req.URL.Path, err.Error())
		return nil, nil, errors.New("系统无法路由当前请求,请联系开发人员进行配置")
	}
	if nil == routeConfMap {
		config.Logger().Errorf("系统无法路由当前请求,请联系开发人员进行配置, urlPath: %s", req.URL.Path)
		return nil, nil, errors.New("系统无法路由当前请求,请开发人员配置转发设置")
	}

	routeConf, ok := routeConfMap[req.URL.Path]
	if !ok {
		config.Logger().Errorf("请开发人员配置转发设置, urlPath: %s", req.URL.Path)
		return nil, nil, errors.New("请开发人员配置转发设置")
	}

	//校验App是否已经下线
	if !checkAppIsOnline(routeConf.AppId) {
		config.Logger().Warnf("当前服务已下线, app id: %s", routeConf.AppId)
		return nil, nil, errors.New("当前服务已下线")
	}

	//鉴权
	if !auth.CheckSession(&routeConf) {
		config.Logger().Warnf("当前会话鉴权无效, routeConf id: %d", routeConf.Id)
		return nil, nil, errors.New("当前会话鉴权无效")
	}

	//请求频率检测
	accessTotal, checked := limit.CheckAccessRateLimit(&routeConf)
	if !checked {
		config.Logger().Warnf("请求过于频繁，请稍后再试, routeConf id: %d", routeConf.Id)
		return nil, nil, errors.New("请求过于频繁，请稍后再试")
	}

	//过载保护
	overload, chk := limit.CheckOverloadLimit(&routeConf)
	if !chk {
		config.Logger().Warnf("服务器请求繁忙，请稍后重试, routeConf id: %d", routeConf.Id)
		return nil, nil, errors.New("服务器请求繁忙，请稍后重试")
	}

	//获取真实的链接
	eurekaClient := config.Eureka()
	httpUrl, err := eurekaClient.GetRealHttpUrl(routeConf.ServiceUrl)
	if nil != err {
		config.Logger().Errorf("获取下游真实地址异常，请稍后重试, routeConf id: %d, error: %v", routeConf.Id, err)
		return nil, nil, errors.New("服务器请求异常，请稍后重试")
	}

	//调用远程服务
	remoteRsp, rspHeader, err := callRemoteService(httpUrl, body, int64(routeConf.Timeout), req)
	if nil != err {
		config.Logger().Errorf("转发请求至下游异常, routeConf id: %d, error: %v", routeConf.Id, err)
		return nil, nil, errors.New("服务器转发异常，请稍后重试")
	}

	//访问数量增加一次
	limit.TotalIncr(&routeConf, accessTotal, overload)

	if rspmode.ENCRYPT != routeConf.RspMode {
		remoteRsp = encryptRsp(remoteRsp)
	}

	return remoteRsp, rspHeader, nil
}

func callRemoteService(url string, postData []byte, timeout int64, req *http.Request) ([]byte, http.Header, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	reqBytes := bytes.NewBuffer(postData)
	request, _ := http.NewRequest(http.MethodPost, url, reqBytes)

	if nil != req.Header {
		for key := range req.Header {
			request.Header.Set(key, req.Header.Get(key))
		}
	} else {
		request.Header.Set(httpheader.Connection, "keep-alive")
		request.Header.Set(httpheader.ContentType, "application/json;charset=UTF-8")
		request.Header.Set(httpheader.UserAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36")
	}

	//设置真实IP地址
	request.Header.Set(httpheader.RemoteAddr, req.RemoteAddr) //兼容PHP
	request.Header.Set(httpheader.XRealIp, req.RemoteAddr)
	request.Header.Set(httpheader.XForwardedFor, buildXForwardedForHeader(req.Header.Get(httpheader.XForwardedFor)))

	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return []byte(""), nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, response.Header, nil
}

// buildXForwardedForHeader 构造 X-Forwarded-For 报头
func buildXForwardedForHeader(xForwardedFor string) string {
	localIp := net.GetLocalIp()
	if 0 == len(xForwardedFor) {
		return localIp
	}

	return xForwardedFor + "," + localIp
}

// checkAppIsOnline 校验应用是否在线
func checkAppIsOnline(appId string) bool {
	//TODO: 校验应用是否在线
	return true
}

// encryptRsp 加密应答
func encryptRsp(rsp []byte) []byte {
	//TODO 加密应答
	return rsp
}

func getPostParams(rw http.ResponseWriter, req *http.Request) (base.ApiRequest, error) {
	body, err := io.ReadAll(req.Body)
	if nil != err {
		config.Logger().Error(err.Error())
		return base.ApiRequest{}, err
	}
	err = req.Body.Close()
	if nil != err {
		config.Logger().Error(err.Error())
		return base.ApiRequest{}, err
	}

	contentType := strings.ToLower(rw.Header().Get(httpheader.ContentType))

	var requestData = base.ApiRequest{}
	if strings.Contains(contentType, medietype.ApplicationJson) {
		//application/json协议
		err := json.Unmarshal(body, &requestData)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			config.Logger().Error(err.Error())
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
