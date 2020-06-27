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
	// 删除主机
	beego.Router("/host", &controllers.HostController{}, "delete:DeleteHost")
	// 获取主机列表
	beego.Router("/host", &controllers.HostController{}, "get:GetHosts")
	// 测试主机ssh连接
	beego.Router("/host/ssh", &controllers.HostController{}, "post:HostConnectionTest")
	// 测试主机svr连接
	beego.Router("/host/svr", &controllers.HostController{}, "post:HostConnectionSvr")

	// svr 启动后的回调接口
	beego.Router("/cb/host/baseinfo", &controllers.HostController{}, "post:BaseInfoCallBack")

	// 关注一个主机
	beego.Router("/host/watch", &controllers.HostWatchController{}, "post:AddHostWatch")
	// 取消关注一个主机
	beego.Router("/host/watch", &controllers.HostWatchController{}, "delete:DeleteHostWatch")

	// 创建主机监控
	beego.Router("/task/host", &controllers.TaskController{}, "post:AddHostInfoTask")

	// 增加一个负责人邮件
	beego.Router("/host/blame", &controllers.HostBlameEmailController{}, "post:AddHostBlameEmail")
	// 删除一个负责人邮件
	beego.Router("/host/blame", &controllers.HostBlameEmailController{}, "delete:DeleteHostBlameEmail")
}
