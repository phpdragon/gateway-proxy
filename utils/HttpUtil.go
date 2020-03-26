package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HttpPostByte(url string, postData []byte) ([]byte, error) {
	client := &http.Client{}
	reqBytes := bytes.NewBuffer(postData)
	request, _ := http.NewRequest("POST", url, reqBytes)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-type", "application/json;charset=UTF-8")
	request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36")
	response, error := client.Do(request)
	if error != nil || response.StatusCode != 200 {
		return nil, error
	}

	return ioutil.ReadAll(response.Body)
}

func HttpPost(url string, postData string) (string, error) {
	client := &http.Client{}
	reqBytes := bytes.NewBuffer([]byte(postData))
	request, _ := http.NewRequest("POST", url, reqBytes)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-type", "application/json;charset=UTF-8")
	request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36")
	response, error := client.Do(request)
	if error != nil || response.StatusCode != 200 {
		return "", nil
	}

	body, error := ioutil.ReadAll(response.Body)
	return string(body), error
}

func HttpGet(url string) (string, error) {
	// 提交get请求
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Connection", "keep-alive")
	response, error := client.Do(request)
	if error != nil || response.StatusCode != 200 {
		return "", nil
	}

	body, error := ioutil.ReadAll(response.Body)
	if error != nil || response.StatusCode != 200 {
		return "", error
	}
	return string(body), error
}
