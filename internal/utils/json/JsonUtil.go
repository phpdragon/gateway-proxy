package json

import (
	"bytes"
	"encoding/json"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"os"
)

func ToJsonString(v interface{}) (string, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		config.Logger().Error("解析失败:", err.Error())
		return "", nil
	}

	return string(jsonStr), nil
}

func ToStringByte(v interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		config.Logger().Error("解析失败:", err.Error())
		return nil, nil
	}
	return jsonStr, nil
}

func ToJSON(jsonStr string, v interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonStr), &v)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		config.Logger().Error("解析失败:" + err.Error())
		return nil, err
	}
	return v, nil
}

func DumJsonStr(jsonStr string) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonStr), "", "\t")

	if err != nil {
		config.Logger().Error(err.Error())
	}

	_, _ = out.WriteTo(os.Stdout)
}

func Str2Map(jsonStr string) (map[string]interface{}, error) {
	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 errorcode 信息
	if err := json.Unmarshal([]byte(jsonStr), &mapResult); err != nil {
		config.Logger().Error(err.Error())
		return nil, err
	}

	return mapResult, nil
}
