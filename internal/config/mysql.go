package config

import (
	"fmt"
	"github.com/astaxie/beego/orm"

	// import your used driver
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func NewMySql() {
	dbConfig := &appConfig.Database
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.DbName, dbConfig.Charset)

	// set default database
	if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
		Logger().Error("Init db failed. error: %s", fmt.Sprint(err))
		os.Exit(1)
	}

	Logger().Infof("Init db success. url: %s", fmt.Sprintf("%s:****@tcp(%s)/%s?charset=%s", dbConfig.User, dbConfig.Host, dbConfig.DbName, dbConfig.Charset))
}
