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
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return nil, err
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
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return "", nil
	}

	body, err := ioutil.ReadAll(response.Body)
	return string(body), err
}

func HttpGet(url string) (string, error) {
	// 提交get请求
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Connection", "keep-alive")
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		return "", nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil || response.StatusCode != 200 {
		return "", err
	}
	return string(body), err
}
