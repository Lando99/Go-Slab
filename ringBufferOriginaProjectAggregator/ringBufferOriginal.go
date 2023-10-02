package ringBufferOriginaProjectAggregator

import (
	"sync"
)

type PacketBuffer struct {
	buffer     []interface{}
	BufferSize int
	readIndex  int
	writeIndex int
	mutex      sync.Mutex
}

func NewPacketBuffer(bufferSize int) *PacketBuffer {
	return &PacketBuffer{
		buffer:     make([]interface{}, bufferSize),
		BufferSize: bufferSize,
	}
}

func (pb *PacketBuffer) Write(b interface{}) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	nextIndex := (pb.writeIndex + 1) % pb.BufferSize
	if nextIndex == pb.readIndex {
		// Buffer is full, discard the oldest packet
		pb.readIndex = (pb.readIndex + 1) % pb.BufferSize
	}

	pb.buffer[pb.writeIndex] = b
	pb.writeIndex = nextIndex
}

func (pb *PacketBuffer) Read() (interface{}, bool) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	if pb.readIndex == pb.writeIndex {
		return nil, false
	}

	out := pb.buffer[pb.readIndex]
	pb.readIndex = (pb.readIndex + 1) % pb.BufferSize

	return out, true
}

func (pb *PacketBuffer) ReadSlice(d int) ([]interface{}, bool) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()

	if pb.readIndex == pb.writeIndex {
		return nil, false
	}

	if (pb.readIndex+d > pb.BufferSize) || (pb.readIndex+d > pb.writeIndex) {
		return nil, false
	}

	out := pb.buffer[pb.readIndex : pb.readIndex+d]
	pb.readIndex = (pb.readIndex + d) % pb.BufferSize

	return out, true
}
