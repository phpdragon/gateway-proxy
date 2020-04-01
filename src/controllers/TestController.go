package controllers

import "net/http"

type TestController struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

func (this *TestController) Get() {

}
