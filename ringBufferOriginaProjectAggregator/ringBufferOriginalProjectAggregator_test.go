package ringBufferOriginaProjectAggregator

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type TestData struct {
	ID   int
	Name string
}

func BenchmarkRingBufferOriginalWriteAndRead(b *testing.B) {
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	//stats := make(map[string]int64)

	rb := NewPacketBuffer(10)

	data := TestData{}

	startTime := time.Now()
	for i := 0; i < 1000; i++ {
		data.ID = i
		data.Name = "test"
		rb.Write(data)

	}

	for i := 0; i < 1000; i++ {
		_, ok := rb.Read()
		if ok {
			//fmt.Println(value)
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
