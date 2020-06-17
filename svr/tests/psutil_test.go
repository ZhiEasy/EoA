package test

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"testing"
)

func TestGetCPUInfo(t *testing.T) {
	cpuInfo, _ := cpu.Info()
	fmt.Println(cpuInfo)

	times, _ := cpu.Times(true)
	fmt.Println(times)

	per, _ := cpu.Percent(0, true)
	fmt.Println("---",per)
}

func TestGetMemInfo(t *testing.T) {
}