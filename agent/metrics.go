package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetOs() string {
	os := runtime.GOOS
	return os
}

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

// Block of temerature code
func GetCPUTemperature() (float64, error) {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0, err
	}

	for _, t := range temps {
		if t.SensorKey == "Package id 0" || t.SensorKey == "CPU" {
			return t.Temperature, nil
		}
	}

	if len(temps) > 0 {
		return temps[0].Temperature, nil
	}

	return 0, fmt.Errorf("CPU temperature not found")
}

func GetTemperature() float64 {
	temp, err := GetCPUTemperature()
	if err != nil {
		fmt.Printf("Get error: %s", err)
	}
	return temp
}

//
