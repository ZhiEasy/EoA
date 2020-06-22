package models

type BaseInfo struct {
	CpuCores int `json:"cpu_cores"`  // 核心数 个
	CpuCacheSize int `json:"cpu_cache_size"`  // 缓存大小 MB
	CpuMhz int `json:"cpu_mhz"`  // 2200Mhz
	CpuModelName string `json:"cpu_model_name"`  // 型号
	DiskTotal int `json:"disk_total"`  // 磁盘大小 GB
	HostBootTime uint64 `json:"host_boot_time"`  // 启动时间
	HostOS string `json:"host_os"`  // 操作系统
	HostName string `json:"host_name"`  // 主机名
	HostPlatform string `json:"host_platform"`  // 操作平台
	HostPlatformVersion string `json:"host_platform_version"`  // 操作系统版本
}
