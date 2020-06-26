package models

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type HostInfo struct {
	BootTime        uint64 `json:"boot_time"`        // 启动时间
	OS              string `json:"os"`               // 操作系统
	Name            string `json:"name"`             // 主机名
	Platform        string `json:"platform"`         // 操作平台
	PlatformVersion string `json:"platform_version"` // 操作平台版本
}

func GetHostInfo() *HostInfo {
	i, _ := host.Info()
	return &HostInfo{
		BootTime:        i.BootTime,
		OS:              i.OS,
		Name:            i.Hostname,
		Platform:        i.Platform,
		PlatformVersion: i.PlatformVersion,
	}
}

type DiskInfo struct {
	Total       uint64  `json:"total"`        // 磁盘总大小
	Used        uint64  `json:"used"`         // 使用量
	UsedPercent float64 `json:"used_percent"` // 使用率
}

func GetDiskInfo() *DiskInfo {
	i, _ := disk.Usage("/")
	return &DiskInfo{
		Total:       i.Total / 1024 / 1024 / 1024,
		Used:        i.Used / 1024 / 1024 / 1024,
		UsedPercent: i.UsedPercent,
	}
}

type MemInfo struct {
	Total       uint64  `json:"total"`        // 内存总量
	Used        uint64  `json:"used"`         // 内存使用量
	UsedPercent float64 `json:"used_percent"` // 内存使用率
}

func GetMemInfo() *MemInfo {
	i, _ := mem.VirtualMemory()
	return &MemInfo{
		Total:       i.Total / 1024 / 1024 / 1024,
		Used:        i.Used / 1024 / 1024 / 1024,
		UsedPercent: i.UsedPercent,
	}
}

type CpuInfo struct {
	Cores        int32     `json:"cores"`         // CPU 核心数
	CacheSize    int32     `json:"cache_size"`    // CPU 缓存大小
	Mhz          float64   `json:"mhz"`           // CPU 频率
	ModelName    string    `json:"model_name"`    // CPU 型号
	PercentTotal []float64 `json:"percent_total"` // CPU 总使用率
}

func GetCpuInfo() *CpuInfo {
	i, _ := cpu.Info()
	p, _ := cpu.Percent(time.Duration(0), false)
	return &CpuInfo{
		Cores:        i[0].Cores,
		CacheSize:    i[0].CacheSize,
		Mhz:          i[0].Mhz,
		ModelName:     i[0].ModelName,
		PercentTotal: p,
	}
}

type BaseInfo struct {
	HostInfo *HostInfo `json:"host_info"`
	DiskInfo *DiskInfo  `json:"disk_info"`
	MemInfo  *MemInfo   `json:"mem_info"`
	CpuInfo  *CpuInfo   `json:"cpu_info"`
}

func GetBaseInfo() *BaseInfo {
	return &BaseInfo{
		HostInfo: GetHostInfo(),
		DiskInfo: GetDiskInfo(),
		MemInfo:  GetMemInfo(),
		CpuInfo:  GetCpuInfo(),
	}
}
