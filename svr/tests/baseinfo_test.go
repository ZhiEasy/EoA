package test

import (
	"github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"testing"
)

func TestGetBaseInfo(t *testing.T) {
	cpuinfo, _ := cpu.Info()
	c := cpuinfo[0]
	logrus.Warnln(c.Cores, c.CacheSize, c.Mhz, c.ModelName)

	diskinfo, _ :=disk.Usage("/")
	logrus.Warnln(diskinfo.Total / 1024 / 1024 / 1024)

	hostinfo, _ := host.Info()

	logrus.Warnln(hostinfo.BootTime,
		hostinfo.OS,
		hostinfo.Hostname,
		hostinfo.KernelArch,
		hostinfo.KernelVersion,
		hostinfo.Platform,
		hostinfo.PlatformVersion)

	meminfo, _ := mem.VirtualMemory()
	logrus.Warnln(meminfo.Total/1024/1024/1024)
}
