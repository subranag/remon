package remon

import "fmt"

const memInfoPath = "/proc/meminfo"
const MemTotalFormat = "MemTotal: %d kB"
const MemFreeFormat = "MemFree: %d kB"

type MemInfo struct {
	MemTotal     uint64
	MemFree      uint64
	MemAvailable uint64
	Buffers      uint64
	Cached       uint64
	SwapTotal    uint64
	SwapFree     uint64
}

type MemInfoReader interface {

	//Read reads current memory info into memInfo
	//returns an error if there is a problem reading memInfo
	Read(memInfo *MemInfo) error

	//Close releases any resources associated with
	//reading memory info this may vary based in implementation
	Close()
}

type fileMemInfoReader struct {
	fileReader *fileReader
}

func NewMemInfoReader() (MemInfoReader, error) {
	fileReader, err := newReader(memInfoPath)

	if err != nil {
		return nil, err
	}

	return &fileMemInfoReader{fileReader: fileReader}, nil
}

func (m *fileMemInfoReader) Read(memInfo *MemInfo) error {
	return m.fileReader.processLines(func(bytes []byte) {
		var memFree uint64
		n, _ := fmt.Sscanf(string(bytes), MemFreeFormat, &memFree)

		if n > 0 {
			fmt.Println(memFree)
		}
	})
}

func (m *fileMemInfoReader) Close() {
	m.fileReader.close()
}
