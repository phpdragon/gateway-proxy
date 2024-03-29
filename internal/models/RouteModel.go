package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	// import your used driver
	_ "github.com/go-sql-driver/mysql"
)

var routeMap *map[string]Route

func (t *Route) TableName() string {
	return "t_route"
}

func init() {
	// register model
	orm.RegisterModel(new(Route))
}

// QueryAllActiveRoutes ORM 操作说明请查看 https://beego.me/docs/mvc/model/object.md
func QueryAllActiveRoutes() (map[string]Route, error) {
	if nil != routeMap {
		return *routeMap, nil
	}

	dbOrm := orm.NewOrm()
	var routes []Route
	_, err := dbOrm.QueryTable(Route{}).Filter("status", StatusEnable).All(&routes)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	var dataMap = make(map[string]Route)
	for _, route := range routes {
		dataMap[route.UrlPath] = route
	}
	routeMap = &dataMap
	return dataMap, nil
}
