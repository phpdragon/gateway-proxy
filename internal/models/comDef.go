package models

const (
	StateDisable = 0
	StateEnable  = 1
)

// Route 定义模型结构体，属性大写开头表示public权限，小写为私有
// 服务路由映射表
type Route struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`         //应用ID
	UrlPath    string `orm:"column(url_path)" json:"url_path"`     //请求路径，建议:/appName/module/action
	ServiceUrl string `orm:"column(service_url)" json:"service"`   //下游Url，支持eureka模式和域名、ip端口模式
	RateLimit  int    `orm:"column(rate_limit)" json:"rate_limit"` //频率限制，每秒次数
	Timeout    int    `orm:"column(timeout)" json:"timeout"`       //超时时间，单位秒
	RspMode    int    `orm:"column(rsp_mode)" json:"rsp_mode"`     //应答模式：0-明文,1-加密
	Remark     string `orm:"column(remark)" json:"remark"`         //请求路径描述
	State      int    `orm:"column(state)" json:"state"`           //1:启用,0:禁用
	UpdateTime string `orm:"column(update_time)" json:"update_time"`
	CreateTime string `orm:"column(create_time)" json:"create_time"`
}
