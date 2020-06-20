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
// 已登录，函数返回 user_id
func (c *BaseController)LoginRequired() int {
	var resp models.Response
	userId := c.GetSession("user_id")
	if userId == nil {
		resp.Status = -1
		resp.Data = nil
		resp.Msg = "认证失败"
		c.Data["json"] = resp
		c.ServeJSON()
		c.StopRun()
	}
	return userId.(int)
}

// 返回JSON Response信息
func (c *BaseController)ReturnResponse(status int, msg string, data interface{}, stopRun bool) {
	var resp models.Response
	resp.Status = status
	resp.Msg = msg
	resp.Data = data
	c.Data["json"] = resp
	c.ServeJSON()
	if stopRun {
		c.StopRun()
	}
	return
}
