package json

import (
	"bytes"
	"encoding/json"
	"github.com/phpdragon/gateway-proxy/internal/config"
	"os"
)

// Ife2JsonStr interface转JSON字符串, 类似于对象转JSON字符串
func Ife2JsonStr(object interface{}) (string, error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		config.Logger().Error("解析失败:", err.Error())
		return "", nil
	}

	return string(jsonStr), nil
}

// Ife2Byte interface 转 []byte, 类似于对象转比特
func Ife2Byte(object interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		config.Logger().Error("解析失败:", err.Error())
		return nil, nil
	}
	return jsonStr, nil
}

// ByteToJsonIfe []byte 转 interface, 类似于比特转对象
func ByteToJsonIfe(jsonByte []byte) (interface{}, error) {
	object := new(interface{})
	err := json.Unmarshal(jsonByte, &object)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		config.Logger().Error("解析失败:" + err.Error())
		return nil, err
	}
	return object, nil
}

// ByteToJsonIfe2 []byte 转 interface, 类似于比特转对象
func ByteToJsonIfe2(jsonByte []byte, object interface{}) (interface{}, error) {
	err := json.Unmarshal(jsonByte, &object)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		config.Logger().Error("解析失败:" + err.Error())
		return nil, err
	}
	return object, nil
}

// Str2JsonIfe 字符串转 interface, 类似于字符串转对象
func Str2JsonIfe(jsonStr string, obj interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(jsonStr), &obj)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		config.Logger().Error("解析失败:" + err.Error())
		return nil, err
	}
	return obj, nil
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
