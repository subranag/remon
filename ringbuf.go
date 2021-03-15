package remon

import "errors"

//BufferFull error indicates that the buffer is full and no more items can be added
var BufferFull = errors.New("ring buffer is full cannot add items")

// RingBuffer stores entries in a circular buffer of given capacity
// NOTE: RingBuffer is not threadsafe please use it only in a single go routine for now
type RingBuffer interface {
	// Add the given thing into ring buffer
	Add(thing interface{}) error

	// Capacity returns the capacity of the ring buffer
	Capacity() uint

	// Next starts from the head of the buffer and keeps going
	// in a circular manner getting the next element next does not honor size
	// it simply walks the buffer the caller has to decide when to invoke Add
	Next() interface{}
}

type sliceRingBuffer struct {
	elements []interface{}
	capacity uint
	size     uint
	currPos  uint
}

func NewRingBuffer(capacity uint) RingBuffer {
	return &sliceRingBuffer{
		elements: make([]interface{}, capacity),
		capacity: capacity,
		size:     0,
		currPos:  0,
	}
}

func (s *sliceRingBuffer) Add(thing interface{}) error {
	if s.size == s.capacity {
		return BufferFull
	}

	s.elements[s.size] = thing
	s.size++
	return nil
}

func (s *sliceRingBuffer) Capacity() uint {
	return s.capacity
}

func (s *sliceRingBuffer) Next() interface{} {
	element := s.elements[s.currPos]
	s.currPos = (s.currPos + 1) % s.capacity
	return element
}
