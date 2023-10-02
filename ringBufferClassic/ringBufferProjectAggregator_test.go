package ringBufferClassic

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

func TestWriteAndRead(t *testing.T) {
	pb := NewRingBuffer(5)

	data := "test-data"
	pb.Write(data)

	readData, ok := pb.Read()
	if !ok {
		t.Fatalf("Expected data to be read, but got none")
	}

	if readData.(string) != data {
		t.Fatalf("Expected %s, but got %s", data, readData)
	}
}

func TestWriteBeyondCapacity(t *testing.T) {

	pb := NewRingBuffer(3)

	pb.Write("data1")
	pb.Write("data2")
	pb.Write("data3")
	pb.Write("data4") // This should overwrite "data1"

	readData, ok := pb.Read()
	if !ok || readData.(string) != "data2" {
		t.Fatalf("Expected data2, but got %v", readData)
	}

	readData, ok = pb.Read()
	if !ok || readData.(string) != "data3" {
		t.Fatalf("Expected data3, but got %v", readData)
	}

	readData, ok = pb.Read()
	if !ok || readData.(string) != "data4" {
		t.Fatalf("Expected data3, but got %v", readData)
	}

}

func TestWrite(t *testing.T) {

	pb := NewRingBuffer(3)

	pb.Write("data1")
	pb.Write("data2")
	pb.Write("data3")
	pb.Write("data4") // This should overwrite "data1"

	readData, ok := pb.Read()
	if !ok || readData.(string) != "data2" {
		t.Fatalf("Expected data2, but got %v", readData)
	}

	readData, ok = pb.Read()
	if !ok || readData.(string) != "data3" {
		t.Fatalf("Expected data3, but got %v", readData)
	}

	readData, ok = pb.Read()
	if !ok || readData.(string) != "data4" {
		t.Fatalf("Expected data3, but got %v", readData)
	}

	readData, ok = pb.Read()
	if ok {
		t.Fatalf("Expected no data, but got %v", readData)
	}

}

func TestReadFromEmptyBuffer(t *testing.T) {
	pb := NewRingBuffer(3)

	_, ok := pb.Read()
	if ok {
		t.Fatalf("Expected no data from an empty buffer")
	}
}

func BenchmarkRingBufferSlabWriteAndRead(b *testing.B) {
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)
	rb := NewRingBuffer(10)

	data := TestData{ID: 1, Name: "Test"}

	//stats := make(map[string]int64)
	b.ResetTimer()
	startTime := time.Now()
	for i := 0; i < 100; i++ {
		data.ID = i
		rb.Write(data)

	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Elapsed Time: %v\n", elapsedTime)
	b.StopTimer()
	//utility.PrintStats2(rb.arena, stats)

	// Restart the timer for the actual benchmark

	for i := 0; i < 10; i++ {
		_, _ = rb.Read()
		//fmt.Println((*TestData)(unsafe.Pointer(&value[0])))
	}

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calcola e stampa la differenza nell'uso della memoria
	memUsed := memAfter.Alloc - memBefore.Alloc
	mallocs := memAfter.Mallocs - memBefore.Mallocs
	fmt.Printf("Memoria utilizzata: %v bytes\n", memUsed)
	fmt.Printf("Memoria utilizzata: %v mallocs\n", mallocs)

}
