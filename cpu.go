package remon

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

var cpuAgg = regexp.MustCompile(`^cpu  (?P<user>\d+) (?P<nice>\d+) (?P<system>\d+) (?P<idle>\d+) (?P<iowait>\d+) (?P<irq>\d+) (?P<softirq>\d+) (?P<steal>\d+) (?P<guest>\d+) (?P<guest_nice>\d+)`)
var singleCpu = regexp.MustCompile(`^cpu(?P<num>\d+) (?P<user>\d+) (?P<nice>\d+) (?P<system>\d+) (?P<idle>\d+) (?P<iowait>\d+) (?P<irq>\d+) (?P<softirq>\d+) (?P<steal>\d+) (?P<guest>\d+) (?P<guest_nice>\d+)`)

func ReadCpuStats() {
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Printf("error reading /proc/stat:%v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		aggStat := cpuAgg.FindStringSubmatch(line)
		if aggStat != nil {
			fmt.Printf("%v\n", aggStat[0])
		}

		cpuStat := singleCpu.FindStringSubmatch(line)
		if cpuStat != nil {
			fmt.Printf("%v\n", cpuStat[0])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error scanning /proc/stat")
	}
	fmt.Println(cpuAgg.SubexpNames()[1:])
}
