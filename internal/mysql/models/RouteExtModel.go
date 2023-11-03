package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/mysql/consts"
)

// RouteConf 服务路由配置
type RouteConf struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`         //应用ID
	UrlPath    string `orm:"column(url_path)" json:"url_path"`     //请求路径，建议:/appName/module/action
	ServiceUrl string `orm:"column(service_url)" json:"service"`   //下游Url，支持eureka模式和域名、ip端口模式
	RateLimit  int    `orm:"column(rate_limit)" json:"rate_limit"` //频率限制，每秒次数
	Timeout    int    `orm:"column(timeout)" json:"timeout"`       //超时时间，单位秒
	RspMode    int    `orm:"column(rsp_mode)" json:"rsp_mode"`     //应答模式：0-明文,1-加密
	CrossMode  int    `orm:"column(cross_mode)" json:"cross_mode"` //跨域模式：0-禁止,1-允许,2-配置
	Limit      int    `orm:"column(limit)" json:"limit"`           //限制次数
	Interval   int    `orm:"column(interval)" json:"interval"`     //间隔时间，单位秒
}

// QueryAllActiveRouteConfMap ORM操作
func QueryAllActiveRouteConfMap() ([]RouteConf, error) {
	dbOrm := orm.NewOrm()
	var records []RouteConf

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("r.*", "o.limit,o.interval").
		From(RouteModel.TableName() + " as r").
		LeftJoin(OverloadModel.TableName() + " as o").
		On("r.app_id = o.app_id AND r.url_path = o.url_path AND r.state = o.state").
		Where("r.state = ?").String()

	_, err := dbOrm.Raw(sql, consts.StateEnable).QueryRows(&records)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	return records, nil
}
