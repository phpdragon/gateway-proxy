package entity

// RouteConf 服务路由配置
type RouteConf struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`         //应用ID
	UrlPath    string `orm:"column(url_path)" json:"url_path"`     //请求路径，建议:/appName/module/action
	ServiceUrl string `orm:"column(service_url)" json:"service"`   //下游Url，支持eureka模式和域名、ip端口模式
	RateLimit  int    `orm:"column(rate_limit)" json:"rate_limit"` //频率限制，每秒次数
	Timeout    int    `orm:"column(timeout)" json:"timeout"`       //超时时间，单位秒
	RspMode    int    `orm:"column(rsp_mode)" json:"rsp_mode"`     //应答模式：0-明文,1-加密
	Limit      int    `orm:"column(limit)" json:"limit"`           //限制次数
	Interval   int    `orm:"column(interval)" json:"interval"`     //间隔时间，单位秒
}
