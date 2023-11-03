package response

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/httpheader"
	"github.com/phpdragon/gateway-proxy/internal/utils/json"
	"net/http"
	"strings"
)

func WriteRsp(rw http.ResponseWriter, req *http.Request, response interface{}) {
	rw.Header().Set(httpheader.CacheControl, "No-Cache")
	rw.Header().Set(httpheader.ContentType, "application/json; charset=utf-8")
	rw.Header().Set(httpheader.PRAGMA, "No-Cache")
	rw.Header().Set(httpheader.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")

	//设置跨域报头
	setCrossDomainHeaders(rw, req)

	dataBody, err := json.Ife2Byte(response)
	if err != nil {
		config.Logger().Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(dataBody)
	if err != nil {
		config.Logger().Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteByteRsp(rw http.ResponseWriter, req *http.Request, response []byte, rspHeader http.Header, crossDomain bool) {
	if rspHeader != nil {
		for key := range rspHeader {
			keyLower := strings.ToLower(key)
			if keyLower == strings.ToLower(httpheader.ContentLength) ||
				keyLower == strings.ToLower(httpheader.TransferEncoding) {
				continue
			}
			rw.Header().Set(key, rspHeader.Get(key))
		}
	} else {
		rw.Header().Set(httpheader.CacheControl, "No-Cache")
		rw.Header().Set(httpheader.ContentType, "application/json; charset=utf-8")
		rw.Header().Set(httpheader.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
		rw.Header().Set(httpheader.PRAGMA, "No-Cache")
	}

	//设置跨域报头
	if crossDomain {
		setCrossDomainHeaders(rw, req)
	}

	_, err := rw.Write(response)
	if err != nil {
		config.Logger().Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// setCrossDomainHeaders 设置跨域报头
func setCrossDomainHeaders(rw http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get(httpheader.ORIGIN)
	if 0 < len(origin) {
		rw.Header().Set(httpheader.AccessControlAllowOrigin, origin)
		rw.Header().Set(httpheader.AccessControlAllowMethods, req.Method)
		rw.Header().Set(httpheader.AccessControlAllowHeaders, "*")
		rw.Header().Set(httpheader.AccessControlAllowCredentials, "true")
	}
}

func WriteStatusCode(rw http.ResponseWriter, req *http.Request, statusCode int) {
	//设置跨域报头
	setCrossDomainHeaders(rw, req)
	rw.Header().Set(httpheader.ContentLength, "0")
	rw.WriteHeader(statusCode)
}
