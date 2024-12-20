package base

import jsonUtil "github.com/phpdragon/gateway-proxy/internal/utils/json"

type ApiResponse struct {
	Code interface{} `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

type ApiRequest struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

func BuildOK() ApiResponse {
	response := ApiResponse{}
	response.Code = "0"
	response.Msg = "ok"
	return response
}

func BuildFail(code string, msg string) ApiResponse {
	response := ApiResponse{}
	response.Code = code
	response.Msg = msg
	return response
}

func BuildFailByte(code string, msg string) ([]byte, error) {
	return jsonUtil.Ife2Byte(BuildFail(code, msg))
}
