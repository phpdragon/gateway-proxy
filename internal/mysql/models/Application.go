package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/phpdragon/gateway-proxy/internal/mysql/consts"
)

// Application 过载配置表
type Application struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`         //应用ID
	Name       string `orm:"column(name)" json:"name"`             //应用名称
	Remark     string `orm:"column(remark)" json:"remark"`         //应用描述
	CrossMode  int    `orm:"column(cross_mode)" json:"cross_mode"` //跨域模式：0-禁止,1-允许,2-配置
	AuthMode   int    `orm:"column(auth_mode)" json:"auth_mode"`   //鉴权模式：0-不鉴权,1-报头,2-URL
	AuthCode   string `orm:"column(auth_code)" json:"auth_code"`   //鉴权码
	State      int    `orm:"column(state)" json:"state"`           //1:启用,0:禁用
	UpdateTime string `orm:"column(update_time)" json:"update_time"`
	CreateTime string `orm:"column(create_time)" json:"create_time"`
}

// AppModel 暴露自身，让外部能调用绑定方法
var AppModel = &Application{}

func init() {
	// register model
	orm.RegisterModel(AppModel)
}

func (t *Application) TableName() string {
	return "t_application"
}

// QueryAllActiveApp ORM操作
func (t *Application) QueryAllActiveApp() ([]Application, error) {
	return QueryAllActiveApp()
}

// QueryAllActiveApp ORM操作
func QueryAllActiveApp() ([]Application, error) {
	dbOrm := orm.NewOrm()
	var records []Application
	_, err := dbOrm.QueryTable(*AppModel).Filter("state", consts.StateEnable).All(&records)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		return nil, err
	} else if errors.Is(err, orm.ErrNoRows) {
		return nil, err
	}
	return records, nil
}

func (t *Application) IsExistApp(appId string) bool {
	dbOrm := orm.NewOrm()
	return dbOrm.QueryTable(*AppModel).Filter("state", consts.StateEnable).Exist()
}
