package ringBufferSlabInterface

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

func BenchmarkBoh(b *testing.B) {
	rb := NewRingBufferSlab(10, 1024)
	var newValue interface{}
	var ok bool

	go func() {
		for i := 0; i < 10; i++ {
			rb.Write(Person{name: "Alice", age: i})
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			newValue, ok = rb.Read()
			if ok {
				fmt.Printf("Tipo di newValue: %T, valore: %v\n", newValue, newValue)
			}
		}
	}()

}

// Define a test struct to be written to the buffer
type TestData struct {
	ID   int
	Name string
}

type Person struct {
	name string
	age  int
}

// Slab ring write
func BenchmarkRingBufferSlabWrite(b *testing.B) {

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	startTime := time.Now()
	rb := NewRingBufferSlab(10, 100)
	var data TestData
	data.ID = 1
	data.Name = "name"

	for i := 0; i < 1000; i++ {
		//rb.Write(data)
		data.ID = i
		rb.Write(&data)
	}
	for i := 0; i < 1000; i++ {
		value, _ := rb.Read()

		if buffer, ok := value.([]byte); ok {
			var testData TestData
			err := binary.Read(bytes.NewReader(buffer), binary.LittleEndian, testData)
			if err != nil {
				b.Fatalf("Deserialization failed: %v", err)
			}
			fmt.Println(testData)
		} else {
			// gestisci l'errore o ignora se non ti aspetti altri tipi
		}
	}
	duration := time.Since(startTime)
	fmt.Printf("Tempo di esecuzione: %v\n", duration)
	b.ResetTimer()

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	// Calcola e stampa la differenza nell'uso della memoria
	memUsed := memAfter.Mallocs - memBefore.Mallocs
	fmt.Printf("Memoria utilizzata: %v bytes\n", memUsed)

}

// Slab ring read
func BenchmarkRingBufferRead(b *testing.B) {
	data := TestData{ID: 1, Name: "Test"}
	rb := NewRingBufferSlab(10000, int(unsafe.Sizeof(data)))

	// Populate the buffer outside the timer
	for i := 0; i < 1000; i++ {
		rb.Write(&data)
	}
	//b.ResetTimer() // Reset the timer to ignore setup time
	for i := 0; i < 1000; i++ {
		value, _ := rb.Read()
		_ = value
	}

}
