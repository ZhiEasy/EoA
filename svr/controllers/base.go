package controllers

import (
	"crypto/md5"
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
func (c *BaseController)LoginRequired() {
	ak := c.GetSession("access_key")
	md5.Sum("")
	if ak == nil {
		c.StopRun()
	}
	//if ak !=
}

// 返回JSON Response信息
func (c *BaseController)ReturnResponse(data interface{}, stopRun bool) {
	c.Data["json"] = data
	c.ServeJSON()
	if stopRun {
		c.StopRun()
	}
	return
}

