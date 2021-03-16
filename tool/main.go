package main

import (
	"fmt"
	"time"

	"github.com/subranag/remon"
)

type stud struct {
	data int
}

func main() {
	stats := make(remon.CpuStats)
	prevStats := make(remon.CpuStats)

	for i := 0; i < 200; i++ {
		remon.ReadCpuStats(stats)
		if len(prevStats) > 0 {
			fmt.Printf("utilization:%v\n", stats["cpu1"].Utilization(prevStats["cpu1"]))
		}
		time.Sleep(100 * time.Millisecond)
		stats.Copy(prevStats)
	}
}
