package ringBufferSlab

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

// TestData is a sample struct to be stored within the buffer.
type TestData struct {
	ID   int
	Name string
}

func TestRingBuffer_WriteRead(t *testing.T) {
	rb := NewRingBufferSlab(5, 100)
	data := []byte("abc")
	buf := make([]byte, unsafe.Sizeof(data))

	ptrData := (*[]byte)(unsafe.Pointer(&buf[0]))
	*ptrData = data
	rb.Write(buf)

	result, ok := rb.Read()

	if !ok || !reflect.DeepEqual(*(*[]byte)(unsafe.Pointer(&result[0])), []byte("abc")) {
		t.Errorf("Expected 'abc', got %s", *(*[]byte)(unsafe.Pointer(&result[0])))
	}
}

func TestRingBuffer_ReadEmpty(t *testing.T) {
	rb := NewRingBufferSlab(10, 100)

	_, ok := rb.Read()
	if ok {
		t.Error("Expected no data, but got some")
	}
}

// BenchmarkRingBufferSlabWriteAndRead measures the performance of writing and reading from the Slab-based ring buffer.
func BenchmarkRingBufferSlabWriteAndRead(b *testing.B) {
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	data := TestData{Name: "test", ID: 1}
	rb := NewRingBufferSlab(10, int(unsafe.Sizeof(data)))

	buf := make([]byte, unsafe.Sizeof(data))

	// Record the start time for the benchmark.
	startTime := time.Now()

	// Write data to the buffer.
	for i := 0; i < 1000; i++ {
		data.ID = i
		data.Name = "test"
		ptrData := (*TestData)(unsafe.Pointer(&buf[0]))
		*ptrData = data
		rb.Write(buf)
	}

	// Read data from the buffer.
	for i := 0; i < 1000; i++ {
		_, ok := rb.Read()
		if ok {
			// Uncomment the below line to print the data read from the buffer.
			//fmt.Println((*TestData)(unsafe.Pointer(&value[0])))
		}
	}

	elapsedTime := time.Since(startTime)

	// Display the elapsed time.
	fmt.Printf("Elapsed Time: %v\n", elapsedTime)
	b.StopTimer()

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calculate and print the difference in memory usage.
	memUsed := memAfter.Alloc - memBefore.Alloc
	mallocs := memAfter.Mallocs - memBefore.Mallocs
	fmt.Printf("Memory used: %v bytes\n", memUsed)
	fmt.Printf("Memory allocations: %v mallocs\n", mallocs)
}
