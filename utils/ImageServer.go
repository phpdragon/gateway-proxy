package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type HttpDealImg struct{}

//实现File和FileInfo接口的类
type ReadImg struct {
	buf      *bytes.Reader
	fileUrl  string
	fileData []byte
}

func (r ReadImg) Close() error {
	panic("implement me")
}

func (r ReadImg) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (r ReadImg) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (r ReadImg) Readdir(count int) ([]os.FileInfo, error) {
	panic("implement me")
}

func (r ReadImg) Stat() (os.FileInfo, error) {
	panic("implement me")
}

//获取C的图片数据
func ReadImgData(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return pix
}

func (self HttpDealImg) Open(name string) (http.File, error) {
	img_name := name[1:]
	fmt.Println(img_name)
	//C(文件服务器地址)
	img_url := "https://www.baidu.com/favicon.ico"
	img_data := ReadImgData(img_url)  //向服务器气球图片数据
	if len(img_data) == 0 {
		fmt.Println("file access forbidden:", name)
		return nil, os.ErrNotExist
	}
	fmt.Println("get img file:", img_url)
	//标红的可以查看标准库bytes的Reader类型，NewReader(p []byte)可返回*Reader，然后调用和http.File相同的Seek()和Read()方法
	var f http.File = &ReadImg{buf: bytes.NewReader(img_data), fileUrl: img_name, fileData: img_data}

	return f, nil
}
