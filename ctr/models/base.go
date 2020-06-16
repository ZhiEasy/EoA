package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func Init() {
	_ = orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("sqlconn"))
}