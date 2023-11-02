package dao

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/mysql/consts"
	"github.com/phpdragon/gateway-proxy/internal/mysql/entity"
	"github.com/phpdragon/gateway-proxy/internal/mysql/models/overload"
	"github.com/phpdragon/gateway-proxy/internal/mysql/models/route"
)

var routeConfMap *map[string]entity.RouteConf

// QueryAllActiveRoutes ORM操作
func QueryAllActiveRoutes() (map[string]entity.RouteConf, error) {
	if nil != routeConfMap {
		return *routeConfMap, nil
	}

	dbOrm := orm.NewOrm()
	var routeConfList []entity.RouteConf

	qb, _ := orm.NewQueryBuilder("mysql")
	sql := qb.Select("r.*", "o.*").
		From(route.EmptyModel.TableName() + " as r").
		LeftJoin(overload.EmptyModel.TableName() + " as o").On("r.app_id = o.app_id AND r.url_path = o.url_path").
		Where("r.state = o.state AND r.state = ?").String()

	_, err := dbOrm.Raw(sql, consts.StateEnable).QueryRows(&routeConfList)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	var dataMap = make(map[string]entity.RouteConf)
	for _, routeConf := range routeConfList {
		dataMap[routeConf.UrlPath] = routeConf
	}
	routeConfMap = &dataMap
	return dataMap, nil
}
