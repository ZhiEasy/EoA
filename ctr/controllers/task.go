package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"time"
)

// TaskController operations for Task
type TaskController struct {
	BaseController
}

// 创建主机监控报警
// TODO ctr 系统重新启动时候将以前的任务重建
func (c *TaskController) AddHostInfoTask() {
	userId := c.LoginRequired(true)
	var req models.AddHostInfoTask
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &req)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	userObj, err := models.GetUserById(userId)
	if err != nil {
		c.ReturnResponse(models.AUTH_ERROR, nil, true)
	}
	hostObj, err := models.GetHostById(req.HostId)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
		return
	}

	// 更新主机信息
	dls, err := json.Marshal(req.DiskLine)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	hostObj.DiskLine = string(dls)
	mls, err := json.Marshal(req.MemLine)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	hostObj.MemLine = string(mls)
	cls, err := json.Marshal(req.CpuLine)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	hostObj.CpuLine = string(cls)
	_ = models.UpdateHostById(hostObj)

	var taskObj models.Task
	taskObj.UserId = userObj
	taskObj.HostId = hostObj
	taskObj.Spec = req.Spec
	taskObj.Name = fmt.Sprintf("%v", time.Now().Unix())
	taskObj.Description = req.Description
	taskObj.Type = 0 // 0 主机监控
	_, _ = models.AddTask(&taskObj)

	// 创建 Task
	t := toolbox.NewTask(taskObj.Name, req.Spec, func() error {
		resp, err := SvrTest(hostObj.Ip)
		if err != nil {
			logrus.Warnf("监控：获取主机信息失败 %v", err)
			return nil
		}
		var data models.BaseInfo
		tmp, _ := json.Marshal(resp.Data)
		_ = json.Unmarshal([]byte(string(tmp)), &data)

		// 整理邮件列表
		// 整理邮件列表放在这里是因为可能会删除一些邮件负责人，当放在外面时候就获取不到最新的邮件列表了
		var emailList []string
		var hws []models.HostWatch
		c.o = orm.NewOrm()
		_, _ = c.o.QueryTable(new(models.HostWatch)).Filter("user_id", userObj.Id).Filter("host_id", hostObj.Id).All(&hws)
		for _, e := range hws {
			emailList = append(emailList, e.Email)
		}
		var hbe []models.HostBlameEmail
		_, _ = c.o.QueryTable(new(models.HostBlameEmail)).Filter("host_id", hostObj.Id).All(&hbe)
		for _, e := range hbe {
			emailList = append(emailList, e.Email)
		}

		if data.MemInfo.UsedPercent > req.MemLine[1] || data.MemInfo.UsedPercent < req.MemLine[0] {
			logrus.Infof("内存报警：已使用%v，限制条件%v\n", data.MemInfo.UsedPercent, req.MemLine)
			context := fmt.Sprintf("任务Id：%s<br/>"+
				"主机IP：%s<br/>"+
				"报警原因：当前内存使用率超出范围<br/>"+
				"设定范围：%s<br/>"+
				"<p style=\"color: red;\">当前：%v</p>", taskObj.Name, hostObj.Ip, hostObj.MemLine, data.MemInfo.UsedPercent)
			SendMail(emailList, "EoA主机监控报警", context)
		}
		if data.CpuInfo.PercentTotal[0] > req.CpuLine[1] || data.CpuInfo.PercentTotal[0] < req.CpuLine[0] {
			logrus.Infof("CPU报警：已使用%v，限制条件%v\n", data.CpuInfo.PercentTotal[0], req.CpuLine)
			context := fmt.Sprintf("任务Id：%s<br/>"+
				"主机IP：%s<br/>"+
				"报警原因：当前【CPU使用率】超出范围<br/>"+
				"设定范围：%s<br/>"+
				"<p style=\"color: red;\">当前：%v</p>", taskObj.Name, hostObj.Ip, hostObj.CpuLine, data.CpuInfo.PercentTotal)
			SendMail(emailList, "EoA主机监控报警", context)
		}
		if data.DiskInfo.UsedPercent > req.DiskLine[1] || data.DiskInfo.UsedPercent < req.DiskLine[0] {
			logrus.Infof("CPU报警：已使用%v，限制条件%v\n", data.DiskInfo.UsedPercent, req.DiskLine)
			context := fmt.Sprintf("任务Id：%s<br/>"+
				"主机IP：%s<br/>"+
				"报警原因：当前【磁盘使用率】超出范围<br/>"+
				"设定范围：%s<br/>"+
				"<p style=\"color: red;\">当前：%v</p>", taskObj.Name, hostObj.Ip, hostObj.DiskLine, data.DiskInfo.UsedPercent)
			SendMail(emailList, "EoA主机监控报警", context)
		}

		// 保存监控信息
		hostinfo := models.HostInfo{
			HostId:     hostObj,
			Info:       string(tmp),
		}
		_, _ = models.AddHostInfo(&hostinfo)

		return nil
	})
	toolbox.AddTask(taskObj.Name, t)
	toolbox.StartTask()

	c.ReturnResponse(models.SUCCESS, taskObj.Name, true)
}

// 取消主机监控任务
func (c *TaskController) DeleteHostInfoTask() {
	userId := c.LoginRequired(true)

	taskId := c.GetString("task_id")
	var taskObj models.Task
	c.o = orm.NewOrm()
	err := c.o.QueryTable(new(models.Task)).Filter("id", taskId).One(&taskObj)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	// 非本人创建不得删除
	if taskObj.UserId.Id != userId {
		c.ReturnResponse(models.AUTH_ERROR, nil, true)
	}

	// 删除 task 信息
	err = models.DeleteTask(taskObj.Id)
	if err != nil {
		logrus.Warnf("删除任务失败\n用户：%v\n参数taskId：%v\n原因：%v", userId, taskId, err)
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	c.ReturnResponse(models.SUCCESS, taskObj, true)
}