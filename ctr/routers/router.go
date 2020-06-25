package routers

import (
	"github.com/ahojcn/EoA/ctr/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 获取语雀的 OAuth 地址
	beego.Router("/oauth/yuque", &controllers.UserController{}, "get:GetYuQueOAuthPath")
	// 语雀授权回调接口
	beego.Router("/user/oauth/yuque", &controllers.UserController{}, "get:YuQueOAuthRedirect")
	// 用户登录
	beego.Router("/login", &controllers.UserController{}, "post:UserLogin")
	// 退出登录
	beego.Router("/logout", &controllers.UserController{}, "delete:UserLogout")
	// 用户完善信息接口
	beego.Router("/user", &controllers.UserController{}, "post:UpdateUserInfo")
	// 获取当前登录的用户信息
	beego.Router("/user", &controllers.UserController{}, "get:GetUserInfo")
	// 添加主机
	beego.Router("/host", &controllers.HostController{}, "post:AddHost")
	// 获取主机列表
	beego.Router("/host", &controllers.HostController{}, "get:GetHosts")
	// 测试主机连接
	beego.Router("/host/test", &controllers.HostController{}, "post:HostConnectionTest")

	// 关注一个主机
	beego.Router("/host/watch", &controllers.HostWatchController{}, "post:AddHostWatch")
	// 取消关注一个主机
	beego.Router("/host/watch", &controllers.HostWatchController{}, "delete:DeleteHostWatch")
}
