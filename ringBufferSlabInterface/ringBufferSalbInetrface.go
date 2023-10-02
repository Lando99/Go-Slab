// Package ringBufferSlabInterface provides an implementation of a ring buffer using the slab memory allocator.
package ringBufferSlabInterface

import (
	"github.com/couchbase/go-slab"
	"sync"
	"unsafe"
)

// RingBufferSlab RingBuffer represents a circular buffer (or ring buffer) where data can be written to and read from.
type RingBufferSlab struct {
	buffer     []slab.Loc // The buffer that will store the data locations.
	bufferSize int        // The size of the ring buffer.
	readIndex  int        // The index from where we will read the next piece of data.
	writeIndex int        // The index where we will write the next piece of data.
	full       bool
	arena      *slab.Arena // The memory arena used for slab allocation.
	mutex      sync.Mutex  // Mutex to ensure thread safety.

	buf     []byte
	dataPtr *interface{}
}

// NewRingBufferSlab NewRingBuffer initializes and returns a new RingBuffer with the specified buffer and slab sizes.
func NewRingBufferSlab(bufferSize, slabSize int) *RingBufferSlab {
	return &RingBufferSlab{
		buffer:     make([]slab.Loc, bufferSize),
		bufferSize: bufferSize,
		arena:      slab.NewArena(1, slabSize, 2, nil),
		readIndex:  0, // The index from where we will read the next piece of data.
		writeIndex: 0, // The index where we will write the next piece of data.
		full:       false,
	}

}

func (rb *RingBufferSlab) Write(b interface{}) {
	//stats := make(map[string]int64)
	rb.mutex.Lock()
	if rb.full {
		rb.buf = rb.arena.LocToBuf(rb.buffer[rb.readIndex])
		rb.arena.DecRef(rb.buf)
		rb.readIndex = (rb.readIndex + 1) % rb.bufferSize
	}

	rb.buf = rb.arena.Alloc(int(unsafe.Sizeof(b)))

	rb.dataPtr = (*interface{})(unsafe.Pointer(&rb.buf[0]))
	*rb.dataPtr = b
	rb.buffer[rb.writeIndex] = rb.arena.BufToLoc(rb.buf)

	rb.writeIndex = (rb.writeIndex + 1) % rb.bufferSize

	rb.full = rb.readIndex == rb.writeIndex

	rb.mutex.Unlock()
	//PrintStats(rb.arena, stats)
}

func (rb *RingBufferSlab) Read() (interface{}, bool) {
	rb.mutex.Lock()

	if !rb.full && rb.readIndex == rb.writeIndex {
		// Buffer vuoto
		rb.mutex.Unlock()
		return nil, false
	}

	rb.buf = rb.arena.LocToBuf(rb.buffer[rb.readIndex])
	rb.dataPtr = (*interface{})(unsafe.Pointer(&rb.buf[0]))
	rb.arena.DecRef(rb.buf)

	rb.readIndex = (rb.readIndex + 1) % rb.bufferSize
	rb.full = false

	rb.mutex.Unlock()
	return *rb.dataPtr, true
}
