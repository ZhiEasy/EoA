package main

import (
	"github.com/ahojcn/EoA/ctr/models"
	_ "github.com/ahojcn/EoA/ctr/routers"
	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	models.Init()
	beego.Run()
}
