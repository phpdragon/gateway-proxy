package models

const (
	STATUS_DISABLE = 0
	STATUS_ENABLE  = 1
)

//定义模型结构体，属性大写开头表示public权限，小写为私有
type Route struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`
	UrlPath    string `orm:"column(url_path)" json:"url_path"`
	ServiceUrl string `orm:"column(service_url)" json:"service"`
	RateLimit int    `orm:"column(rate_limit)" json:"rate_limit"`
	Timeout    int    `orm:"column(timeout)" json:"timeout"`
	Status     int    `orm:"column(status)" json:"status"`
	Timestamp  string `orm:"column(timestamp)" json:"timestamp"`
}
