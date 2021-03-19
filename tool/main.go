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
	stats := make(remon.CpuStats)
	prevStats := make(remon.CpuStats)
	statsReader, err := remon.NewCpuStatsReader()

	if err != nil {
		fmt.Printf("error reading cpu stats err:%v", err)
		os.Exit(1)
	}
	defer statsReader.Close()

	for i := 0; i < 200; i++ {
		statsReader.Read(stats)
		if len(prevStats) > 0 {
			for k, v := range stats {
				util := v.Utilization(prevStats[k])
				fmt.Printf("%v:%v ", k, util)
			}
			fmt.Println()
		}
		stats.CopyTo(prevStats)
		time.Sleep(500 * time.Millisecond)
	}
}
