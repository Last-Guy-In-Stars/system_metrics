package main

import (
	"fmt"
	"os"
	"time"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetOsName() string {
	os_name, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting os name")
		return "No name on host"
	}
	return os_name
}

func GetCPU() int32 {
	cpu, err := cpu.Percent(time.Second, false) // Параметр false берет среднюю по всем процессорам
	if err != nil {
		fmt.Println("Error getting cpu")
		return 0
	}
	if len(cpu) > 0 {
		return int32(cpu[0])
	}
	return 0
}

func GetMemory() int32 {
	mem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error getting memory")
		return 0
	}
	return int32(mem.Used / (1024 * 1024)) // МБ
}