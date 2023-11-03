package router

import (
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/base"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/errorcode"
	"github.com/phpdragon/gateway-proxy/internal/consts/httpcode"
	"github.com/phpdragon/gateway-proxy/internal/consts/route"
	"github.com/phpdragon/gateway-proxy/internal/logic/request"
	"github.com/phpdragon/gateway-proxy/internal/logic/response"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	jsonUtil "github.com/phpdragon/gateway-proxy/internal/utils/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type routeInfo struct {
	path    string
	handler http.HandlerFunc
}

var (
	gFaviconIco, _ = os.ReadFile("favicon.ico")
	routePath      []routeInfo
)

func buildRoutes() []routeInfo {
	if nil != routePath {
		return routePath
	}

	routePath := []routeInfo{
		{path: "^/favicon.ico$", handler: favicon},
		//处理eureka的心跳等
		{path: "^/actuator\\w*", handler: config.Eureka().ServeHTTP},
		//监听日志级别设置
		{path: fmt.Sprintf("^%slogLevel$", route.RouterSystem), handler: config.GetAtomicLevel().ServeHTTP},
		{path: fmt.Sprintf("^%s\\w*", route.RouterSystem), handler: systemHandle},
		//请求入口
		{path: "^/\\w+$", handler: indexHandle}, // \w：匹配字母、数字、下划线
	}
	return routePath
}

func WebHandler() http.HandlerFunc {
	routePath := buildRoutes()
	return func(rspWriter http.ResponseWriter, request *http.Request) {
		for _, route := range routePath {
			ok, err := regexp.Match(route.path, []byte(request.URL.Path))
			if err != nil {
				config.Logger().Error(err.Error())
			}
			if ok {
				route.handler(rspWriter, request)
				return
			}
		}
		_, _ = rspWriter.Write([]byte("404 not found"))
	}
}

func favicon(writer http.ResponseWriter, _ *http.Request) {
	_, err := writer.Write(gFaviconIco)
	if err != nil {
		config.Logger().Infof(err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func indexHandle(rw http.ResponseWriter, req *http.Request) {
	startTime := date.GetCurrentTimeMillis()
	config.Logger().Infof("")

	//处理请求
	rsp, rspHeader, crossDomain, err := request.HandleHttpRequest(req)
	if nil != err {
		//当前系统的报错支持跨域处理
		crossDomain = true
		rsp, _ = base.BuildFailByte(errorcode.SystemError, err.Error())
	}

	//处理跨域返回
	if http.MethodOptions == req.Method && crossDomain {
		response.WriteStatusCode(rw, req, httpcode.NoContent)
		return
	}

	response.WriteByteRsp(rw, req, rsp, rspHeader, crossDomain)

	//打印方法执行耗时的信息
	endTime := date.GetCurrentTimeMillis()
	printExecTime(startTime, endTime)
}

func systemHandle(rw http.ResponseWriter, req *http.Request) {
	request.HandleSystemRequest(req)

	rsp, _ := jsonUtil.Ife2Byte(base.BuildOK())
	_, err := rw.Write(rsp)
	if err != nil {
		config.Logger().Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// 打印方法执行耗时的信息
func printExecTime(startTime int64, endTime int64) {
	diffTime := endTime - startTime
	diffTimeStr := fmt.Sprintf("请求处理结束,耗时: %s ms\n", strconv.FormatInt(diffTime, 10))
	if diffTime > 1000 {
		config.Logger().Warn(diffTimeStr)
	} else {
		config.Logger().Infof(diffTimeStr)
	}
}
