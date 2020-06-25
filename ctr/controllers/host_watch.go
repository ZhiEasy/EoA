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
func (c *HostWatchController) AddHostWatch() {
	userId := c.LoginRequired(true)

	var req models.AddHostWatchReq
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	err = AddHostWatch(userId, req.HostId)
	if err != nil {
		c.ReturnResponse(models.HOST_REWATCH, nil, true)
	}

	c.ReturnResponse(models.SUCCESS, nil, true)
}

// 取消关注一个主机
func (c *HostWatchController) DeleteHostWatch() {
	userId := c.LoginRequired(true)

	hostId := c.GetString("host_id")

	var hw models.HostWatch
	c.o = orm.NewOrm()
	qs := c.o.QueryTable(new(models.HostWatch))
	err := qs.Filter("user_id", userId).Filter("host_id", hostId).One(&hw)
	if err != nil {
		c.ReturnResponse(models.REQUEST_ERROR, nil, true)
	}

	err = models.DeleteHostWatch(hw.Id)
	if err != nil {
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
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
		UserId: userObj,
		HostId: hostObj,
		Email:  userObj.Email,
	}
	_, err := models.AddHostWatch(&hwObj)
	return err
}
