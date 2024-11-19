package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/mysql/consts"
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
	CrossMode  int    `orm:"column(cross_mode)" json:"cross_mode"` //跨域模式：0-禁止,1-允许,2-配置
	Remark     string `orm:"column(remark)" json:"remark"`         //请求路径描述
	State      int    `orm:"column(state)" json:"state"`           //1:启用,0:禁用
	UpdateTime string `orm:"column(update_time)" json:"update_time"`
	CreateTime string `orm:"column(create_time)" json:"create_time"`
}

// RouteModel 暴露自身，让外部能调用绑定方法
var RouteModel = &Route{}

func init() {
	// register model
	orm.RegisterModel(RouteModel)
}

func (t *Route) TableName() string {
	return "t_route"
}

// QueryAllActiveRoutes ORM操作
func QueryAllActiveRoutes() (map[string]Route, error) {
	dbOrm := orm.NewOrm()
	var routes []Route
	_, err := dbOrm.QueryTable(Route{}).Filter("state", consts.StateEnable).All(&routes)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	var dataMap = make(map[string]Route)
	for _, route := range routes {
		dataMap[route.UrlPath] = route
	}
	return dataMap, nil
}
