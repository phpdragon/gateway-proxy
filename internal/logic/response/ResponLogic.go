package response

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/medietype"
	"github.com/phpdragon/gateway-proxy/internal/utils/json"
	"net/http"
	"strings"
)

func WriteRsp(rw http.ResponseWriter, req *http.Request, response interface{}) {
	origin := req.Header.Get(medietype.ORIGIN)
	rw.Header().Set(medietype.CacheControl, "No-Cache")
	rw.Header().Set(medietype.ContentType, "application/json; charset=utf-8")
	rw.Header().Set(medietype.PRAGMA, "No-Cache")
	rw.Header().Set(medietype.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(medietype.AccessControlAllowOrigin, origin)
		rw.Header().Set(medietype.AccessControlAllowCredentials, "true")
	}

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

func WriteByteRsp(rw http.ResponseWriter, req *http.Request, response []byte, rspHeader http.Header) {
	if rspHeader != nil {
		for key := range rspHeader {
			keyLower := strings.ToLower(key)
			if keyLower == medietype.ContentLength || keyLower == medietype.TransferEncoding {
				continue
			}
			rw.Header().Set(key, rspHeader.Get(key))
		}
	} else {
		rw.Header().Set(medietype.CacheControl, "No-Cache")
		rw.Header().Set(medietype.ContentType, "application/json; charset=utf-8")
		rw.Header().Set(medietype.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
		rw.Header().Set(medietype.PRAGMA, "No-Cache")
	}

	//设置跨域报头
	origin := req.Header.Get(medietype.ORIGIN)
	if 0 < len(origin) {
		rw.Header().Set(medietype.AccessControlAllowOrigin, origin)
		rw.Header().Set(medietype.AccessControlAllowCredentials, "true")
	}

	_, err := rw.Write(response)
	if err != nil {
		config.Logger().Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
