package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

// Report represents a report collected about this system.
type Report struct {
	Hostname       string
	OS             string
	Threads        int
	CPU            string
	CPULoad        float64
	Uptime         uint64
	Virtualization string

	MemoryTotal uint64
	MemoryUsed  uint64

	SwapTotal uint64
	SwapUsed  uint64

	Disk []Disk

	LoadAverage []float64

	Timestamp time.Time
}

// Disk represents a physical disk's information
type Disk struct {
	Mount string
	Used  uint64
	Total uint64
}

// GenerateReport generates a Report
func GenerateReport() (report Report) {
	report = Report{}

	// hostname
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("[warn] os.Hostname() returned error %v\n", err)
	} else {
		report.Hostname = hostname
	}

	// cpu name
	info, err := cpu.Info()
	if err != nil {
		log.Printf("[warn] cpu.Info() returned error %v\n", err)
	} else {
		report.CPU = info[0].ModelName
		report.Threads = len(info)
	}

	// cpu load %
	pct, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Printf("[warn] cpu.Percent() returned error %v\n", err)
	} else {
		report.CPULoad = pct[0] / 100
	}

	// host info
	hostinfo, err := host.Info()
	if err != nil {
		log.Printf("[warn] host.Info() returned error %v\n", err)
	} else {
		report.Uptime = hostinfo.Uptime
		report.OS = fmt.Sprintf("%s (%s, %s)", hostinfo.OS, hostinfo.KernelArch, hostinfo.KernelVersion)
		report.Virtualization = hostinfo.VirtualizationRole + "," + hostinfo.VirtualizationSystem
	}

	// RAM info
	raminfo, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("[warn] mem.VirtualMemory() returned error %v\n", err)
	} else {
		report.MemoryTotal = raminfo.Total
		report.MemoryUsed = raminfo.Used
	}

	// Swap info
	swapinfo, err := mem.SwapMemory()
	if err != nil {
		log.Printf("[warn] mem.SwapMemory() returned error %v\n", err)
	} else {
		report.SwapTotal = swapinfo.Total
		report.SwapUsed = uint64(float64(swapinfo.Total) * (swapinfo.UsedPercent / 100))
	}

	// Disk info
	diskinfo, err := disk.Partitions(false)
	if err != nil {
		log.Printf("[warn] disk.Partitions() returned error %v\n", err)
	} else {
		for _, i := range diskinfo {
			usage, err := disk.Usage(i.Mountpoint)
			if err != nil {
				log.Printf("[warn] disk.Usage(%s) returned error %v\n", i.Mountpoint, err)
			} else {
				report.Disk = append(report.Disk, Disk{
					Mount: i.Mountpoint,
					Used:  usage.Used,
					Total: usage.Total,
				})
			}
		}
	}

	loadinfo, err := load.Avg()
	if err != nil {
		log.Printf("[warn] load.Avg() returned error %v\n", err)
	} else {
		report.LoadAverage = []float64{loadinfo.Load1, loadinfo.Load5, loadinfo.Load15}
	}

	report.Timestamp = time.Now()

	return
}

// ToString turns this Report into a human-readable representation, for debugging purposes
func (report Report) ToString() string {
	virt := ""

	if len(report.Virtualization) > 1 {
		virt = " (" + report.Virtualization + ")"
	}

	disks := ""

	for _, disk := range report.Disk {
		disks += fmt.Sprintf("%s (%d/%d)\n", disk.Mount, disk.Used, disk.Total)
	}

	return fmt.Sprintf("%s%s\nOS: %s\nThreads: %d\nCPU: %s\nLoad average: %f\nCPU load: %f\nUptime: %ds\nRAM: %d/%d\nSwap: %d/%d\n\nDisks:\n%s",
		report.Hostname,
		virt,
		report.OS,
		report.Threads,
		report.CPU,
		report.LoadAverage,
		report.CPULoad,
		report.Uptime,
		report.MemoryUsed,
		report.MemoryTotal,
		report.SwapUsed,
		report.SwapTotal,
		disks)
}
