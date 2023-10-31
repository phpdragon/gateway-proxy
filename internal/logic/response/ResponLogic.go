package response

import (
	"github.com/phpdragon/gateway-proxy/internal/config"
	"github.com/phpdragon/gateway-proxy/internal/consts/medietype"
	"github.com/phpdragon/gateway-proxy/internal/utils/json"
	"net/http"
)

func WriteJson(rw http.ResponseWriter, req *http.Request, response interface{}, isJson bool) {
	origin := req.Header.Get(medietype.ORIGIN)
	rw.Header().Set(medietype.CacheControl, "No-Cache")
	rw.Header().Set(medietype.ContentType, "application/json; charset=utf-8")
	rw.Header().Set(medietype.PRAGMA, "No-Cache")
	rw.Header().Set(medietype.EXPIRES, "Thu, 01 Jan 1970 00:00:00 GMT")
	if 0 < len(origin) {
		rw.Header().Set(medietype.AccessControlAllowOrigin, origin)
		rw.Header().Set(medietype.AccessControlAllowCredentials, "true")
	}

	var err error
	var dataBody []byte
	if isJson {
		dataBody, err = json.ToStringByte(response)
		if err != nil {
			config.Logger().Error(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		dataBody = []byte(response.(string))
	}

	_, err = rw.Write(dataBody)
	if err != nil {
		config.Logger().Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
