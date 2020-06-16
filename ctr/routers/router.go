package routers

import (
	"github.com/ahojcn/EoA/ctr/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/user/oauth", &controllers.UserController{}, "get:OAuth")
}
