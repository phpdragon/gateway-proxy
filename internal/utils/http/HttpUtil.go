package http

import (
	"bytes"
	"crypto/tls"
	"github.com/phpdragon/gateway-proxy/internal/utils/json"
	"io"
	"net/http"
	"time"
)

func PostByte(url string, postData []byte, timeout int64) (interface{}, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	reqBytes := bytes.NewBuffer(postData)
	request, _ := http.NewRequest(http.MethodPost, url, reqBytes)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-type", "application/json;charset=UTF-8")
	request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36")
	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	//将[]byte转JSON对象
	return json.ByteToJsonIfe(body)
}

func Post(url string, postData string, timeout int64) (interface{}, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	reqBytes := bytes.NewBuffer([]byte(postData))
	request, _ := http.NewRequest(http.MethodPost, url, reqBytes)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-type", "application/json;charset=UTF-8")
	request.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36")
	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return "", nil
	}

	body, err := io.ReadAll(response.Body)
	//将[]byte转JSON对象
	return json.ByteToJsonIfe(body)
}

func Get(url string, timeout int64) (interface{}, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 提交get请求
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Connection", "keep-alive")
	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	//将[]byte转JSON对象
	return json.ByteToJsonIfe(body)
}
