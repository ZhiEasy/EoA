package controllers

import (
	"github.com/astaxie/beego"
)

// BaseController operations for All Controller
type BaseController struct {
	beego.Controller
}


// 权限校验
func (c *BaseController)LoginRequired() {
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

