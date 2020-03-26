package logic

import (
	"../consts"
	"../core"
	eureka "../eureka-client"
	"../utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HandleHttpRequest(req *http.Request, client *eureka.EurekaClient) (interface{}, error) {
	body, _ := ioutil.ReadAll(req.Body)
	_ = req.Body.Close()

	httpUrl := client.GetRealHttpUrl("http://CLOUD-FILE-PROXY/qcloud/getCosAuthConfig")

	requestData := core.ApiRequest{}
	err := json.Unmarshal([]byte(body), &requestData)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//调用远程服务
	//"http://172.16.1.124:12247/qcloud/getCosAuthConfig"
	remoteData,err := callRemoteService(httpUrl, body)
	if nil != err{
		return nil, err
	}

	response := core.BuildOK()
	response.Data = remoteData
	return response, nil
}

func callRemoteService(httpUrl string, req []byte) (interface{},error) {
	return utils.HttpPostByte(httpUrl, req)
}

func getPostParams(rw http.ResponseWriter, req *http.Request) (core.ApiRequest, error) {
	body, err := ioutil.ReadAll(req.Body)
	if nil != err {
		return core.ApiRequest{}, err
	}
	err = req.Body.Close()
	if nil != err {
		return core.ApiRequest{}, err
	}

	contentType := strings.ToLower(rw.Header().Get(consts.CONTENT_TYPE))

	var requestData = core.ApiRequest{}
	if strings.Contains(contentType, consts.APPLICATION_JSON) {
		//application/json协议
		err := json.Unmarshal([]byte(body), &requestData)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			fmt.Println(err)
			return requestData, nil
		}
	} else if strings.Contains(contentType, consts.APPLICATION_X_WWW_FORM_URLENCODED) {
		//application/x-www-form-urlencoded协议
	} else if strings.Contains(contentType, consts.MULTIPART_FORM_DATA) {
		//multipart/form-data协议
	} else if strings.Contains(contentType, consts.APPLICATION_OCTET_STREAM) {
		//application/octet-stream协议
	} else {
		//其他文本协议，只当文本处理
	}
	return requestData, nil
}
