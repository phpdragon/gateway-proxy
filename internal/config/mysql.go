package config

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	// import your used driver
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func init() {
	dbConfig := GetDatabaseConfig()
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.DbName, dbConfig.Charset)

	// set default database
	if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
		log.Println("Init db failed. errorCode: ", fmt.Sprint(err))
		os.Exit(1)
	}

	log.Println("Init db success. url:", fmt.Sprintf("%s:****@tcp(%s)/%s?charset=%s", dbConfig.User, dbConfig.Host, dbConfig.DbName, dbConfig.Charset))
}
