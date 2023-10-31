package utils

import (
	"encoding/json"
	"github.com/phpdragon/gateway-proxy/internal/core/log"
)

func ToJSONString(v interface{}) (string, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Info("解析失败:" + err.Error())
		return "", nil
	}

	return string(jsonStr), nil
}

func ToJSONStringByte(v interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Info("解析失败:" + err.Error())
		return nil, nil
	}
	return jsonStr, nil
}

func ToJSON(jsonStr string, v interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonStr), &v)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		log.Info("解析失败:" + err.Error())
		return nil, err
	}
	return v, nil
}
