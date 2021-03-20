package remon

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type fileReader struct {
	file   *os.File
	reader *bufio.Reader
}

func newReader(path string) (*fileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading %v:%w\n", path, err)
	}
	return &fileReader{file: file, reader: bufio.NewReader(file)}, nil
}

func (f *fileReader) processLines(callback func([]byte)) error {
	f.file.Seek(0, io.SeekStart)
	f.reader.Reset(f.file)
	for {
		bytes, err := f.reader.ReadBytes('\n')

		if len(bytes) > 0 {
			callback(bytes)
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

func (f *fileReader) close() {
	if f.file != nil {
		f.file.Close()
	}
}
