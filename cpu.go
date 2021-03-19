package remon

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

var cpuAgg = regexp.MustCompile(`^cpu  (?P<user>\d+) (?P<nice>\d+) (?P<system>\d+) (?P<idle>\d+) (?P<iowait>\d+) (?P<irq>\d+) (?P<softirq>\d+) (?P<steal>\d+) (?P<guest>\d+) (?P<guest_nice>\d+)`)
var singleCpu = regexp.MustCompile(`^cpu(?P<num>\d+) (?P<user>\d+) (?P<nice>\d+) (?P<system>\d+) (?P<idle>\d+) (?P<iowait>\d+) (?P<irq>\d+) (?P<softirq>\d+) (?P<steal>\d+) (?P<guest>\d+) (?P<guest_nice>\d+)`)

//CpuUsage is the representation of CPU usage in a system
type CpuUsage struct {
	//Name is the name of the cpu like cpu0, cpu1 ... cpuN
	Name string
	//Idle the idle time of the CPU
	Idle uint64
	//Total the total time of the CPU Used = Total - Idle
	Total uint64
}

//Utilization provides the CPU utilization as compared the CpuUsage provided
func (c *CpuUsage) Utilization(p *CpuUsage) float64 {
	pUsage := p.Total - p.Idle
	cUsage := c.Total - c.Idle

	dp := cUsage - pUsage
	dt := c.Total - p.Total
	return (float64(dp) / float64(dt)) * 100.0
}

//CpuStats is the map representing CPU usage with the key being CPU name
//and the value being CpuUsage of that CPU
type CpuStats map[string]*CpuUsage

//CopyTo copies the stats from this stats to t
func (s CpuStats) CopyTo(t CpuStats) {
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

//CpuStatsReader can be used to read the CPU stats, please note to call the Close
//function to release resources otherwise there might be resource leakage
type CpuStatsReader interface {

	//Read reads the CpuStats for all CPUs in the host
	//aggregate stats of the CPU are present in the "cpu" entry in CpuStats
	//if the CPU stats cannot be read the function returns an error
	Read(stats CpuStats) error

	// Close closes the reader and releases any resources
	// acquired for reading CPU status
	Close()
}

//fileCpuStatsReader is the Linux flavor implementation of the CpuStatsReader
type fileCpuStatsReader struct {
	cpuStatsFile *os.File
	reader       *bufio.Reader
}

//NewCpuStatsReader returns a new CpuStatsReader error otherwise
func NewCpuStatsReader() (CpuStatsReader, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Printf("error reading /proc/stat:%v\n", err)
		return nil, err
	}

	fileStatsReader := &fileCpuStatsReader{cpuStatsFile: file, reader: bufio.NewReader(file)}
	return fileStatsReader, nil
}

func (s *fileCpuStatsReader) Read(stats CpuStats) error {
	s.cpuStatsFile.Seek(0, io.SeekStart)
	s.reader.Reset(s.cpuStatsFile)
	for {
		bytes, err := s.reader.ReadBytes('\n')

		if len(bytes) > 0 {
			agg := cpuAgg.FindSubmatch(bytes)
			if len(agg) > 0 {
				readCpuStat(agg, true, stats)
			}

			cpu := singleCpu.FindSubmatch(bytes)
			if len(cpu) > 0 {
				readCpuStat(cpu, false, stats)
			}
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			// if the error is EOF we simply break
			break
		}
	}
	return nil
}

func (s *fileCpuStatsReader) Close() {
	if s.cpuStatsFile != nil {
		s.cpuStatsFile.Close()
	}
}

func readCpuStat(rawStat [][]byte, agg bool, stats CpuStats) {
	cpuTag := "cpu"
	values := rawStat[1:]
	if !agg {
		cpuTag += string(rawStat[1])
		values = rawStat[2:]
	}
	user := byteToStat(values[0])
	nice := byteToStat(values[1])
	system := byteToStat(values[2])
	idle := byteToStat(values[3])
	iowait := byteToStat(values[4])
	irq := byteToStat(values[5])
	softirq := byteToStat(values[6])
	steal := byteToStat(values[7])

	idleTotal := idle + iowait
	total := user + nice + system + idle + iowait + irq + softirq + steal

	_, present := stats[cpuTag]

	if !present {
		stats[cpuTag] = &CpuUsage{Name: cpuTag, Idle: idleTotal, Total: total}
		return
	}

	stats[cpuTag].Name = cpuTag
	stats[cpuTag].Idle = idleTotal
	stats[cpuTag].Total = total
}

func byteToStat(stat []byte) uint64 {
	val, err := strconv.ParseUint(string(stat), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot have cpu usage value that is not a number found:%v", stat))
	}
	return val
}
