package controllers

import (
	"encoding/json"
	"errors"
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego/orm"
)

// HostWatchController operations for HostWatch
type HostWatchController struct {
	BaseController
}

// 关注一个主机
func (c *HostWatchController)AddHostWatch() {
	userId := c.LoginRequired()

	var req models.AddHostWatchReq
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	err = AddHostWatch(userId, req.HostId)
	if err != nil {
		c.ReturnResponse(models.HOST_REWATCH, nil, true)
	}

	c.ReturnResponse(models.SUCCESS, nil, true)
}

func AddHostWatch(userId, hostId int) error {
	userObj, _ := models.GetUserById(userId)
	hostObj, _ := models.GetHostById(hostId)

	o := orm.NewOrm()
	qs := o.QueryTable(new(models.HostWatch))
	cnt, _ := qs.Filter("user_id__exact", userObj).Filter("host_id__exact", hostObj).Count()
	if cnt != 0 {
		return errors.New("重复关注")
	}
	hwObj := models.HostWatch{
		UserId:     userObj,
		HostId:     hostObj,
		Email:      userObj.Email,
	}
	_, err := models.AddHostWatch(&hwObj)
	return err
}
