package remon

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var cpuAgg = regexp.MustCompile(`^cpu  (?P<user>\d+) (?P<nice>\d+) (?P<system>\d+) (?P<idle>\d+) (?P<iowait>\d+) (?P<irq>\d+) (?P<softirq>\d+) (?P<steal>\d+) (?P<guest>\d+) (?P<guest_nice>\d+)`)
var singleCpu = regexp.MustCompile(`^cpu(?P<num>\d+) (?P<user>\d+) (?P<nice>\d+) (?P<system>\d+) (?P<idle>\d+) (?P<iowait>\d+) (?P<irq>\d+) (?P<softirq>\d+) (?P<steal>\d+) (?P<guest>\d+) (?P<guest_nice>\d+)`)

type CpuUsage struct {
	Name  string
	Idle  uint64
	Total uint64
}

func (c *CpuUsage) Utilization(p *CpuUsage) float64 {
	pUsage := p.Total - p.Idle
	cUsage := c.Total - c.Idle

	dp := cUsage - pUsage
	dt := c.Total - p.Total
	return (float64(dp) / float64(dt)) * 100.0
}

type CpuStats map[string]*CpuUsage

func (s CpuStats) Copy(t CpuStats) {
	for k, v := range s {
		_, ok := t[k]
		if !ok {
			t[k] = &CpuUsage{}
		}
		t[k].Name = v.Name
		t[k].Idle = v.Idle
		t[k].Total = v.Total
	}
}

func ReadCpuStats(stats CpuStats) {
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
			stats["cpu"] = readCpuStat(aggStat, true)
		}

		cpuStat := singleCpu.FindStringSubmatch(line)
		if cpuStat != nil {
			stats["cpu"+cpuStat[1]] = readCpuStat(cpuStat, false)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error scanning /proc/stat")
	}
}

func readCpuStat(rawStat []string, agg bool) *CpuUsage {
	cpuTag := "cpu"
	values := rawStat[1:]
	if !agg {
		cpuTag += rawStat[1]
		values = rawStat[2:]
	}
	user := strToStat(values[0])
	nice := strToStat(values[1])
	system := strToStat(values[2])
	idle := strToStat(values[3])
	iowait := strToStat(values[4])
	irq := strToStat(values[5])
	softirq := strToStat(values[6])
	steal := strToStat(values[7])

	idleTotal := idle + iowait
	total := user + nice + system + idle + iowait + irq + softirq + steal
	return &CpuUsage{Name: cpuTag, Idle: idleTotal, Total: total}
}

func strToStat(stat string) uint64 {
	val, err := strconv.ParseUint(stat, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot have cpu usage value that is not a number found:%v", stat))
	}
	return val
}
