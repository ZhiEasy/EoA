package controllers

import (
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// BaseController operations for All Controller
type BaseController struct {
	beego.Controller
	o orm.Ormer
}


// 登录校验
// 未登录，返回给前端 JSON 信息
// 已登录
// 	1. 信息完善过，函数返回 user_id
//  2. 信息未完善，返回状态
func (c *BaseController)LoginRequired() int {
	var resp models.Response
	userId := c.GetSession("user_id")
	if userId == nil {
		c.ReturnResponse(models.AUTH_ERROR, resp, true)
	}
	//user, err := models.GetUserById(userId.(int))
	//if err != nil {
	//	c.ReturnResponse(models.AUTH_ERROR, nil, true)
	//	return userId.(int)  // 不让下面 user 报警告
	//}
	//var yuque models.YuQueUserInfo
	//_ = json.Unmarshal([]byte(user.YuqueInfo), &yuque)
	//userInfo := models.UserProfile {
	//	Id:         user.Id,
	//	CreateTime: user.CreateTime,
	//	Name:       user.Name,
	//	Email:      user.Email,
	//	AvatarUrl:  yuque.Data.AvatarURL,
	//}
	//// 判断这个已经授权过的用户是否完善了信息
	//if user.Pwd == "" || user.Name == "" || user.Email == "" {
	//	// 如果没有完善信息
	//	c.ReturnResponse(models.NEED_UPDATE_INFO, userInfo, true)
	//}
	return userId.(int)
}

// 返回JSON Response信息
func (c *BaseController)ReturnResponse(code int, data interface{}, stopRun bool) {
	var resp models.Response
	resp.Status = code
	resp.Msg = models.ResponseText(code)
	resp.Data = data
	c.Data["json"] = resp
	c.ServeJSON()
	if stopRun {
		c.StopRun()
	}
	return
}
