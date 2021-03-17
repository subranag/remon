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
	statsReader, err := remon.NewCpuStatsReader()

	if err != nil {
		fmt.Printf("error reading cpu stats err:%v", err)
	}
	defer statsReader.Close()

	statsReader.ReadStats(stats)
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
	statsReader.ReadStats(stats)
	fmt.Printf("%v\n", stats["cpu"])
}
