package routers

import (
	"github.com/ahojcn/EoA/ctr/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 语雀授权回调接口
	beego.Router("/user/oauth", &controllers.UserController{}, "get:OAuth")
	// 用户完善信息接口
	beego.Router("/user", &controllers.UserController{}, "post:UpdateUserInfo")
}
