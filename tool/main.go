package main

import (
	"github.com/subranag/remon"
)

type stud struct {
	data int
}

func main() {
	rb := remon.NewRingBuffer(5)
	counter := 0
	for rb.Next() == nil {
		rb.Add(&stud{data: counter})
		counter += 1
	}
}
