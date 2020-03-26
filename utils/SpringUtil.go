package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

func ToJSONString(v interface{})string {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Fatal( "Resp Json Err:", fmt.Sprint(err))
	}
	return string(jsonStr)
}

func ToJSONStringByte(v interface{})[]byte {
	jsonStr, err := json.Marshal(v)
	if err != nil {
		log.Fatal( "Resp Json Err:", fmt.Sprint(err))
	}
	return jsonStr
}

func ToJSON(jsonStr string,v interface{}) interface{} {
	err := json.Unmarshal([]byte(jsonStr), &v)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		fmt.Println(err)
	}
	return v
}