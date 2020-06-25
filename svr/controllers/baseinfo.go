package controllers

import (
	"encoding/json"
	"github.com/ahojcn/EoA/svr/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type BaseInfoController struct {
	BaseController
}

// 获取主机基础信息
func (c *BaseInfoController)GetBaseInfo()  {
	cpuinfo, _ := cpu.Info()
	diskinfo, _ := disk.Usage("/")
	hostinfo, _ := host.Info()
	meminfo, _ := mem.VirtualMemory()

	d := models.BaseInfo{
		CpuCores:            cpuinfo[0].Cores,
		CpuCacheSize:        cpuinfo[0].CacheSize,
		CpuMhz:              cpuinfo[0].Mhz,
		CpuModelName:        cpuinfo[0].ModelName,
		DiskTotal:           diskinfo.Total / 1024 * 3,
		MemTotal:            meminfo.Total / 1024 * 3,
		HostBootTime:        hostinfo.BootTime,
		HostOS:              hostinfo.OS,
		HostName:            hostinfo.Hostname,
		HostPlatform:        hostinfo.Platform,
		HostPlatformVersion: hostinfo.PlatformVersion,
	}
	b, _ := json.Marshal(d)
	s := string(b)
	c.ReturnResponse(s, true)
}