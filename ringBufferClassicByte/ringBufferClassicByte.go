package ringBufferClassicByte

import (
	"sync"
)

type RingBuffer struct {
	data  [][]byte
	read  int
	write int
	full  bool
	mutex sync.Mutex
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:  make([][]byte, size),
		read:  0,
		write: 0,
	}
}

func (rb *RingBuffer) Write(val []byte) {
	rb.mutex.Lock()

	if rb.full {
		rb.read = (rb.read + 1) % len(rb.data)
	}

	bufferCopy := make([]byte, len(val))
	copy(bufferCopy, val)
	rb.data[rb.write] = bufferCopy

	rb.write = (rb.write + 1) % len(rb.data)

	rb.full = rb.read == rb.write
	rb.mutex.Unlock()
}

func (rb *RingBuffer) Read() ([]byte, bool) {
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
