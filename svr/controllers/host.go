package controllers

import (
	"github.com/Sirupsen/logrus"
	"github.com/ahojcn/EoA/svr/models"
)

type HostController struct {
	BaseController
}

// 获取主机信息
func (c *HostController) GetHostInfo() {
	logrus.Warnln(111)
	info := models.GetBaseInfo()
	c.ReturnResponse(models.SUCCESS, info, true)
}
