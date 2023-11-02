package overload

import "github.com/astaxie/beego/orm"

// Overload 过载配置表
type Overload struct {
	Id         int    `orm:"column(id);pk" json:"id"`
	AppId      string `orm:"column(app_id)" json:"app_id"`     //应用ID
	UrlPath    string `orm:"column(url_path)" json:"url_path"` //请求路径，建议:/appName/module/action
	Limit      int    `orm:"column(limit)" json:"limit"`       //限制次数
	Interval   int    `orm:"column(interval)" json:"interval"` //间隔时间，单位秒
	Remark     string `orm:"column(remark)" json:"remark"`     //请求路径描述
	State      int    `orm:"column(state)" json:"state"`       //1:启用,0:禁用
	UpdateTime string `orm:"column(update_time)" json:"update_time"`
	CreateTime string `orm:"column(create_time)" json:"create_time"`
}

// EmptyModel 暴露自身，让外部能调用绑定方法
var EmptyModel = &Overload{}

func (t *Overload) TableName() string {
	return "t_overload"
}

func init() {
	// register model
	orm.RegisterModel(EmptyModel)
}
