package controllers

import (
	"encoding/json"
	"github.com/ahojcn/EoA/ctr/models"
)

// HostBlameEmailController operations for HostBlameEmail
type HostBlameEmailController struct {
	BaseController
}

// 增加一个邮件到负责人表
func (c *HostBlameEmailController) AddHostBlameEmail() {
	_ = c.LoginRequired(true)

	var req models.AddHostBlameEmailReq
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	hostObj, err := models.GetHostById(req.HostId)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	for _, e := range req.Email {
		data := models.HostBlameEmail{
			HostId: hostObj,
			Email:  e,
		}
		_, _ = models.AddHostBlameEmail(&data)
	}

	c.ReturnResponse(models.SUCCESS, nil, true)
}

// 删除一个负责人邮件
func (c *HostBlameEmailController) DeleteHostBlameEmail() {
	_ = c.LoginRequired(true)

	id, err := c.GetInt("id")
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	err = models.DeleteHostBlameEmail(id)
	if err != nil {
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	c.ReturnResponse(models.SUCCESS, nil, true)
}
