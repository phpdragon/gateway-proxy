package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"regexp"
)

func dumJsonStr(jsonStr string) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(jsonStr), "", "\t")

	if err != nil {
		log.Fatalln(err)
	}

	out.WriteTo(os.Stdout)
}

func jsonStr2Map(jsonStr string) map[string]interface{} {
	var mapResult map[string]interface{}
	//使用 json.Unmarshal(data []byte, v interface{})进行转换,返回 error 信息
	if err := json.Unmarshal([]byte(jsonStr), &mapResult); err != nil {
		log.Fatal(err)
		return nil
	}

	return mapResult
}

func isIpAddressPort(ipStr string) bool {
	if isOk, _ := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)(\:\d+)$`, ipStr); isOk {
		return isOk
	}
	return false
}
