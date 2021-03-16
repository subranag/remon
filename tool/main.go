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

	for i := 0; i < 200; i++ {
		remon.ReadCpuStats(stats)
		fmt.Printf("%v\n", stats["cpu3"])
		time.Sleep(100 * time.Millisecond)
	}
}
