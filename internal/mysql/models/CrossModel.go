package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/mysql/consts"
)

// CrossDomain 跨域白名单
type CrossDomain struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`     //应用ID
	RouteId    int    `orm:"column(route_id)" json:"route_id"` //路由id：0-全局
	Origin     string `orm:"column(origin)" json:"origin"`     //来源域名链接
	Remark     string `orm:"column(remark)" json:"remark"`     //描述
	State      int    `orm:"column(state)" json:"state"`       //1:启用,0:禁用
	UpdateTime string `orm:"column(update_time)" json:"update_time"`
	CreateTime string `orm:"column(create_time)" json:"create_time"`
}

// CrossDomainModel 暴露自身，让外部能调用绑定方法
var CrossDomainModel = &CrossDomain{}

func init() {
	// register model
	orm.RegisterModel(CrossDomainModel)
}

func (t *CrossDomain) TableName() string {
	return "t_cross_domain"
}

// QueryAllActiveDomains ORM操作
func QueryAllActiveDomains() ([]CrossDomain, error) {
	dbOrm := orm.NewOrm()
	var routes []CrossDomain
	_, err := dbOrm.QueryTable(CrossDomain{}).Filter("state", consts.StateEnable).All(&routes)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}
	return routes, nil
}
