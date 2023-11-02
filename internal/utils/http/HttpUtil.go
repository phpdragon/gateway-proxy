package http

import (
	"bytes"
	"crypto/tls"
	"github.com/phpdragon/gateway-proxy/internal/consts/httpheader"
	"io"
	"net/http"
	"time"
)

func Post(url string, postData string, timeout int64) ([]byte, http.Header, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	reqBytes := bytes.NewBuffer([]byte(postData))
	request, _ := http.NewRequest(http.MethodPost, url, reqBytes)
	request.Header.Set(httpheader.Connection, "keep-alive")
	request.Header.Set(httpheader.ContentType, "application/json;charset=UTF-8")
	request.Header.Set(httpheader.UserAgent, "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36")
	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return []byte(""), nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, response.Header, nil
}

func Get(url string, timeout int64) ([]byte, http.Header, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 提交get请求
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set(httpheader.Connection, "keep-alive")
	response, err := httpClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		return []byte(""), nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, response.Header, nil
}
