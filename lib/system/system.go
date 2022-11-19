package system

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"time"
)

type cpuInfo struct {
	PhysicalCnt     int     `json:"physicalCnt"`     // CPU物理核数
	LogicalCnt      int     `json:"logicalCnt"`      // CPU逻辑核数
	CpuTotalPercent float64 `json:"cpuTotalPercent"` // 1s内负载
}

type memInfo struct {
	Total           uint64  `json:"total"`           // 内存总大小
	Available       uint64  `json:"available"`       // 可用内存大小
	MemTotalPercent float64 `json:"memTotalPercent"` // 内存使用率
}
type diskInfo struct {
	Path        string  `json:"path"`        // 分区名
	Total       uint64  `json:"total"`       // 分区总大小
	Used        uint64  `json:"used"`        // 分区已使用容量
	Free        uint64  `json:"free"`        // 分区空闲容量
	UsedPercent float64 `json:"usedPercent"` // 分区使用百分比
}
type systemInfo struct {
	KernelVersion string `json:"kernelVersion"` // 内核版本
	Platform      string `json:"platform"`      // 平台
	OsFamily      string `json:"osFamily"`      // 操作系统家族
	OsVersion     string `json:"osVersion"`     // 操作系统版本
}

type systemStatus struct {
	Cpu        cpuInfo    `json:"cpu"`
	Mem        memInfo    `json:"mem"`
	Disk       diskInfo   `json:"disk"`
	SystemInfo systemInfo `json:"systemInfo"`
}

func GetInfo() systemStatus {
	cpuStatus, err := getCpuInfo()
	if err != nil {
		//
	}
	memStatus, err := getMemInfo()
	if err != nil {
		//
	}
	diskStatus, err := getDiskInfo()
	if err != nil {
		//
	}
	system, err := getSystemInfo()
	if err != nil {
		//
	}
	status := systemStatus{cpuStatus, memStatus, diskStatus, system}

	return status
}

func getCpuInfo() (cpuInfo, error) {
	physicalCnt, err := cpu.Counts(false) //物理核数
	if err != nil {
		return cpuInfo{}, err
	}
	logicalCnt, err := cpu.Counts(true) // 逻辑核数
	if err != nil {
		return cpuInfo{}, err
	}

	percent, err := cpu.Percent(1*time.Second, false) // 1s内负载
	if err != nil {
		return cpuInfo{}, err
	}
	totalPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", percent[0]), 64)

	return cpuInfo{physicalCnt, logicalCnt, totalPercent}, nil

}

func getMemInfo() (memInfo, error) {
	info, err := mem.VirtualMemory()
	if err != nil {
		return memInfo{}, err
	}
	total := info.Total                                                             // 总内存
	availableMem := info.Available                                                  // 可用内存
	usedPercent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info.UsedPercent), 64) // 使用率

	return memInfo{total, availableMem, usedPercent}, nil
}

func getDiskInfo() (diskInfo, error) {
	info, err := disk.Usage("/") // 需要监控的磁盘分区， 项目中改为从环境变量获取
	if err != nil {
		return diskInfo{}, err
	}
	path := info.Path                                                           // 分区名
	total := info.Total                                                         // 分区总容量
	used := info.Used                                                           // 分区已使用的容量
	free := info.Free                                                           // 分区空闲的容量
	percent, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", info.UsedPercent), 64) // 分区使用百分比

	return diskInfo{path, total, used, free, percent}, nil
}

func getSystemInfo() (systemInfo, error) {
	kernelVersion, err := host.KernelVersion() //内核版本
	if err != nil {
		return systemInfo{}, err
	}
	platform, family, version, err := host.PlatformInformation() // 平台、操作系统、版本
	if err != nil {
		return systemInfo{}, err
	}

	return systemInfo{kernelVersion, platform, family, version}, nil
}
