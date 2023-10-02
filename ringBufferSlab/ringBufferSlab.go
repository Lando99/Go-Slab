package ringBufferSlab

import (
	"github.com/couchbase/go-slab"
	"sync"
)

type RingBufferSlab struct {
	buffer     []slab.Loc
	slabSize   int
	bufferSize int
	readIndex  int
	writeIndex int
	full       bool
	arena      *slab.Arena
	mutex      sync.Mutex
}

// NewRingBufferSlab creates and initializes a new RingBufferSlab with the
// specified buffer size and slab size. It returns a pointer to the newly
// created RingBufferSlab.
func NewRingBufferSlab(bufferSize, slabSize int) *RingBufferSlab {

	// Initialize a new RingBufferSlab instance with the provided sizes and
	// a new memory arena for slab allocation.
	ring := &RingBufferSlab{
		// Array to store the locations of the allocated slabs.
		buffer: make([]slab.Loc, bufferSize),

		// Set the buffer and slab sizes as provided.
		bufferSize: bufferSize,
		slabSize:   slabSize,

		// Create a new memory arena for slab allocation.
		// The parameters define the min and max sizes of the slabs and
		// the factor by which the slab size will grow.
		arena: slab.NewArena(1, slabSize, 2, nil),

		// Initial read and write indices are set to 0.
		readIndex:  0,
		writeIndex: 0,

		// Initially, the buffer is not full.
		full: false,
	}

	// Initialize the buffer by allocating slabs in the arena for each slot.
	ring.initBuffer()

	// Return the initialized RingBufferSlab.
	return ring
}

// initBuffer initializes the RingBufferSlab by allocating slabs of memory within the arena.
// For each slot in the buffer, it allocates a slab of the specified size (rb.slabSize)
// and stores its location in the buffer.
func (rb *RingBufferSlab) initBuffer() {
	for i := 0; i < rb.bufferSize; i++ {
		// Allocate a slab of memory within the arena of size 'rb.slabSize'.
		buf := rb.arena.Alloc(rb.slabSize)

		// Store the location of the allocated slab in the buffer.
		rb.buffer[i] = rb.arena.BufToLoc(buf)
	}
}

// Write writes the provided byte slice 'b' into the RingBufferSlab.
func (rb *RingBufferSlab) Write(b []byte) {
	rb.mutex.Lock() // Lock for concurrency.

	if rb.full { // Move read index if buffer is full.
		rb.readIndex = (rb.readIndex + 1) % rb.bufferSize
	}

	buf := rb.arena.LocToBuf(rb.buffer[rb.writeIndex]) // Get buffer location.
	copy(buf, b)                                       // Copy data.

	// Update buffer and move write index.
	//rb.buffer[rb.writeIndex] = rb.arena.BufToLoc(buf)
	rb.writeIndex = (rb.writeIndex + 1) % rb.bufferSize

	rb.full = rb.readIndex == rb.writeIndex // Update 'full' status.

	rb.mutex.Unlock() // Release the lock.
}

// WriteDeprecated writes the provided byte slice 'b' into the RingBufferSlab.
// This version of the function is less efficient in terms of CPU computing
// compared to the more recent implementations.
func (rb *RingBufferSlab) WriteDeprecated(b []byte) {
	rb.mutex.Lock() // Lock for concurrency.

	// If the buffer is full, it retrieves the buffer at the read index
	// and decreases its reference count before advancing the read index.
	// This extra step can be computationally expensive.
	if rb.full {
		buf := rb.arena.LocToBuf(rb.buffer[rb.readIndex])
		rb.arena.DecRef(buf)
		rb.readIndex = (rb.readIndex + 1) % rb.bufferSize
	}

	// Memory is allocated in the slab for the incoming data.
	buf := rb.arena.Alloc(len(b))
	copy(buf, b) // Copy the provided data into the buffer.

	// Update buffer location and advance the write index.
	rb.buffer[rb.writeIndex] = rb.arena.BufToLoc(buf)
	rb.writeIndex = (rb.writeIndex + 1) % rb.bufferSize

	rb.full = rb.readIndex == rb.writeIndex // Check if buffer is full.

	rb.mutex.Unlock() // Release the lock.
}

func (rb *RingBufferSlab) Read() ([]byte, bool) {
	rb.mutex.Lock()

	if !rb.full && rb.readIndex == rb.writeIndex {
		rb.mutex.Unlock()
		return nil, false
	}

	buf := rb.arena.LocToBuf(rb.buffer[rb.readIndex])

	//rb.arena.DecRef(buf)
	rb.readIndex = (rb.readIndex + 1) % rb.bufferSize
	rb.full = false

	rb.mutex.Unlock()
	return buf, true
}
