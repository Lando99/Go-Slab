package ringBufferClassic

import (
	"sync"
)

type RingBuffer struct {
	data  []interface{}
	read  int
	write int
	full  bool
	mutex sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:  make([]interface{}, size),
		read:  0,
		write: 0,
	}
}

func (rb *RingBuffer) Write(val interface{}) {
	rb.mutex.Lock()

	if rb.full {
		rb.read = (rb.read + 1) % len(rb.data)
	}

	rb.data[rb.write] = val
	rb.write = (rb.write + 1) % len(rb.data)

	rb.full = rb.read == rb.write
	rb.mutex.Unlock()
}

func (rb *RingBuffer) Read() (interface{}, bool) {
	rb.mutex.Lock()
	defer rb.mutex.Unlock()

	if !rb.full && rb.read == rb.write {
		// Buffer vuoto
		return nil, false
	}

	val := rb.data[rb.read]
	rb.read = (rb.read + 1) % len(rb.data)
	rb.full = false

	return val, true
}
