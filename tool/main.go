package main

import (
	"fmt"
	"os"
	"time"

	"github.com/subranag/remon"
)

type stud struct {
	data int
}

func main() {
	readMem()
}

func readMem() {
	memInfo := new(remon.MemInfo)
	memInfoReader, err := remon.NewMemInfoReader()

	if err != nil {
		fmt.Printf("error reading cpu stats err:%v", err)
		os.Exit(1)
	}
	defer memInfoReader.Close()
	for i := 0; i < 200; i++ {
		memInfoReader.Read(memInfo)
		time.Sleep(500 * time.Millisecond)
	}
}

func readCpu() {
	stats := make(remon.CpuStats)
	prevStats := make(remon.CpuStats)
	cpuStats, err := remon.NewCpuStatsReader()

	if err != nil {
		fmt.Printf("error reading cpu stats err:%v", err)
		os.Exit(1)
	}
	defer cpuStats.Close()

	for i := 0; i < 200; i++ {
		cpuStats.Read(stats)
		if len(prevStats) > 0 {
			for k, v := range stats {
				util := v.Utilization(prevStats[k])
				fmt.Printf("%v:%.2f ", k, util)
			}
			fmt.Println()
		}
		stats.CopyTo(prevStats)
		time.Sleep(500 * time.Millisecond)
	}
}
