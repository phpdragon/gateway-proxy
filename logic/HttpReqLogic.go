package logic

import (
	"../consts"
	"../core"
	"../utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HandleHttpRequest(rw http.ResponseWriter, req *http.Request)(core.ApiResponse, error) {
	body, _ := ioutil.ReadAll(req.Body)
	_ = req.Body.Close()

	requestData := core.ApiRequest{}
	err := json.Unmarshal([]byte(body), &requestData)
	//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
	if err != nil {
		fmt.Println(err)
		return core.BuildFail(core.SYSTEM_ERROR,err.Error()), nil
	}

	//调用远程服务
	remoteData := callRemoteService()

	response := core.BuildOK()
	response.Data = remoteData
	return response, nil
}

func callRemoteService()interface{}{
	remoteDataStr := utils.HttpPost("http://172.16.1.124:12247/qcloud/getCosAuthConfig","{\"bucket\":\"adfa\"}")

	var jsonData = utils.ToJSON(remoteDataStr, core.ApiResponse{})
	return jsonData
}

func getPostParams(rw http.ResponseWriter, req *http.Request) core.ApiRequest{
	body, _ := ioutil.ReadAll(req.Body)
	_ = req.Body.Close()
	contentType := strings.ToLower(rw.Header().Get(consts.CONTENT_TYPE))

	var requestData = core.ApiRequest{}
	if(strings.Contains(contentType,consts.APPLICATION_JSON)){
		//application/json协议
		err := json.Unmarshal([]byte(body), &requestData)
		//解析失败会报错，如json字符串格式不对，缺"号，缺}等。
		if err != nil {
			fmt.Println(err)
			return requestData
		}
	}else if(strings.Contains(contentType,consts.APPLICATION_X_WWW_FORM_URLENCODED)){
		//application/x-www-form-urlencoded协议
	}else if(strings.Contains(contentType,consts.MULTIPART_FORM_DATA)){
		//multipart/form-data协议
	}else if(strings.Contains(contentType,consts.APPLICATION_OCTET_STREAM)){
		//application/octet-stream协议
	}else{
		//其他文本协议，只当文本处理
	}
	return requestData
}
