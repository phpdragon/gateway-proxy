package application

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/mysql/consts"
)

// Application 过载配置表
type Application struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"` //应用ID
	Name       string `orm:"column(name)" json:"name"`     //应用名称
	Remark     string `orm:"column(remark)" json:"remark"` //应用描述
	State      int    `orm:"column(state)" json:"state"`   //1:启用,0:禁用
	UpdateTime string `orm:"column(update_time)" json:"update_time"`
	CreateTime string `orm:"column(create_time)" json:"create_time"`
}

// EmptyModel 暴露自身，让外部能调用绑定方法
var EmptyModel = &Application{}
var allActiveApps *map[string]Application

func (t *Application) TableName() string {
	return "t_application"
}

func init() {
	// register model
	orm.RegisterModel(EmptyModel)
}

func IsExist(appId string) bool {
	//dbOrm := orm.NewOrm()
	//return dbOrm.QueryTable(*EmptyModel).Filter("state", consts.StateEnable).Exist()
	apps, _ := QueryAllActiveApp()
	_, ok := apps[appId]
	return ok
}

// QueryAllActiveApp ORM操作
func QueryAllActiveApp() (map[string]Application, error) {
	if nil != allActiveApps {
		return *allActiveApps, nil
	}

	dbOrm := orm.NewOrm()
	var routes []Application
	_, err := dbOrm.QueryTable(*EmptyModel).Filter("state", consts.StateEnable).All(&routes)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}

	var dataMap = make(map[string]Application)
	for _, app := range routes {
		dataMap[app.AppId] = app
	}
	allActiveApps = &dataMap
	return dataMap, nil
}
