package database

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)



var Dao orm.Ormer

func Dbsetting() {

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("dbstring"),30)

	Dao=orm.NewOrm()
}

