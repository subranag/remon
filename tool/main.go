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
	statsReader, err := remon.NewCpuStatsReader()

	if err != nil {
		fmt.Printf("error reading cpu stats err:%v", err)
		os.Exit(1)
	}
	defer statsReader.Close()

	statsReader.Read(stats)
	fmt.Printf("%v\n", stats["cpu0"])
	time.Sleep(500 * time.Millisecond)
	fmt.Println()
	statsReader.Read(stats)
	fmt.Printf("%v\n", stats["cpu0"])
}
