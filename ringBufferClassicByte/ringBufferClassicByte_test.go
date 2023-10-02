package ringBufferClassicByte

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

type TestData struct {
	ID   int
	Name string
}

func TestRingBuffer_WriteRead(t *testing.T) {
	rb := NewRingBuffer(5)
	data := []byte("abc")

	rb.Write(data)

	result, ok := rb.Read()

	if !ok || !reflect.DeepEqual(result, []byte("abc")) {
		t.Errorf("Expected 'abc', got %s", *(*[]byte)(unsafe.Pointer(&result[0])))
	}
}

func TestRingBuffer_ReadEmpty(t *testing.T) {
	rb := NewRingBuffer(5)

	_, ok := rb.Read()
	if ok {
		t.Error("Expected no data, but got some")
	}
}

func BenchmarkRingBufferClassicByteWriteAndRead(b *testing.B) {
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	//stats := make(map[string]int64)

	rb := NewRingBuffer(10)

	data := TestData{}
	buf := make([]byte, unsafe.Sizeof(data))

	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		data.ID = i
		data.Name = "test"
		ptrData := (*TestData)(unsafe.Pointer(&buf[0]))
		*ptrData = data
		rb.Write(buf)

	}

	for i := 0; i < 1000; i++ {
		_, ok := rb.Read()
		if ok {
			//fmt.Println(*(*TestData)(unsafe.Pointer(&value[0])))
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Elapsed Time: %v\n", elapsedTime)
	b.StopTimer()

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calcola e stampa la differenza nell'uso della memoria
	memUsed := memAfter.Alloc - memBefore.Alloc
	mallocs := memAfter.Mallocs - memBefore.Mallocs
	fmt.Printf("Memoria utilizzata: %v bytes\n", memUsed)
	fmt.Printf("Memoria utilizzata: %v mallocs\n", mallocs)

}
