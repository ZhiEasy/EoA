package controllers

import (
	"github.com/ahojcn/EoA/svr/models"
	"github.com/astaxie/beego"
)

// BaseController operations for All Controller
type BaseController struct {
	beego.Controller
}


// 权限校验
func (c *BaseController) AuthRequired() {
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

