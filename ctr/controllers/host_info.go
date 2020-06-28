package controllers

import (
	"encoding/json"
	"github.com/ahojcn/EoA/ctr/models"
	"github.com/astaxie/beego/orm"
	"math"
	"time"
)

// HostInfoController operations for HostInfo
type HostInfoController struct {
	BaseController
}

// 获取主机资源监控记录
func (c *HostInfoController) GetHostInfo() {
	_ = c.LoginRequired(true)

	hostId, err := c.GetInt("host_id")
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	start := c.GetString("start")
	end := c.GetString("end")
	if start == "" || end == "" {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	_, err = time.ParseInLocation("2006-01-02", start, time.Local)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}
	_, err = time.ParseInLocation("2006-01-02", end, time.Local)
	if err != nil {
		c.ReturnResponse(models.REQUEST_DATA_ERROR, nil, true)
	}

	var his []models.HostInfo
	c.o = orm.NewOrm()
	_, err = c.o.QueryTable(new(models.HostInfo)).Filter("host_id", hostId).Filter("create_time__gte", start).Filter("create_time__lte", end).All(&his)
	if err != nil {
		c.ReturnResponse(models.SERVER_ERROR, nil, true)
	}

	data := make(map[string]interface{})
	cpu := make([]float64, 0)
	mem := make([]float64, 0)
	disk := make([]float64, 0)
	t := make([]string, 0)
	for _, hi := range his {
		var info models.BaseInfo
		err = json.Unmarshal([]byte(hi.Info), &info)
		if err != nil {
			continue
		}
		cpu = append(cpu, math.Trunc(info.CpuInfo.PercentTotal[0]*1e2+0.5) * 1e-2)
		mem = append(mem, math.Trunc(info.MemInfo.UsedPercent*1e2+0.5) * 1e-2)
		disk = append(disk, math.Trunc(info.DiskInfo.UsedPercent*1e2+0.5) * 1e-2)
		t = append(t, hi.CreateTime.Format("2006-01-02 15:04:05"))
	}
	data["cpu"] = cpu
	data["mem"] = mem
	data["disk"] = disk
	data["time"] = t

	c.ReturnResponse(models.SUCCESS, data, true)
}
