package main

import (
	"encoding/binary"
	"fmt"
	"github.com/couchbase/go-slab"
	"log"
	"main/utility"
	"math"
	"net/http"
	_ "net/http/pprof"
	"unsafe"
)

type Person struct {
	Name    [25]byte
	Surname [25]byte
	Age     int
}
type Person2 struct {
	Name    string
	Surname string
	Age     int
}
type Session struct {
	ID      string
	UserID  string
	Numbers []int
	// time
	// other fields to manage for the session
}
type Session2 struct {
	ID     int
	Name   [10]byte // We use a byte array instead of a string to simplify the example
	Active bool
}

func main() {

	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
	}()

	// Create Arena
	arena := slab.NewArena(128, 1024, 2, nil)

	// Allocate memory for an integer
	var intBuf []byte
	intBuf = arena.Alloc(8)                      // 8 bytes for an int64
	binary.LittleEndian.PutUint64(intBuf, 57555) // insert value into the buffer

	fmt.Println("Allocated memory for int:", intBuf)
	fmt.Println("Integer stored in buffer:", binary.LittleEndian.Uint64(intBuf))

	// Allocate memory for a float
	var floatBuf []byte
	floatBuf = arena.Alloc(8) // 8 bytes for a float64

	bits := math.Float64bits(123.456)
	binary.LittleEndian.PutUint64(floatBuf, bits)

	fmt.Println("Allocated memory for float:", floatBuf) // print the bytes saved in the buffer
	fmt.Println("Float stored in buffer:", math.Float64frombits(binary.LittleEndian.Uint64(floatBuf)))

	// Allocate memory for a string
	var stringBuf []byte
	stringBuf = arena.Alloc(128) // Bytes needed for the string
	writeStringToBuffer(stringBuf, "Hello, World!")
	fmt.Println("Allocated memory for string:", stringBuf)
	fmt.Println("String stored in buffer:", string(stringBuf))

	// Allocate memory for a PERSON variable
	personBuf := arena.Alloc(int(unsafe.Sizeof(Person{})))

	// Insert values into the buffer
	writeStringToBuffer(personBuf[:25], "John")
	writeStringToBuffer(personBuf[25:50], "Doe")
	binary.LittleEndian.PutUint32(personBuf[50:], uint32(30))

	fmt.Println("Allocated memory for person:", personBuf)
	fmt.Println("String stored in buffer:", string(personBuf[0:25]), string(personBuf[25:50]), binary.LittleEndian.Uint64(personBuf[50:]))

	// Create an array of INTEGER ARRAYS
	nums := []int64{10, 20, 30, 40, 50}
	// Calculate how many bytes are needed for the array
	size := len(nums) * binary.Size(nums[0])
	// Allocate memory for the array
	intArrayBuf := arena.Alloc(size)

	// Write the array to the buffer
	for i, num := range nums {
		binary.LittleEndian.PutUint64(intArrayBuf[i*8:], uint64(num))
	}
	fmt.Println("Allocated memory for int:", intArrayBuf)

	for i := 0; i < len(nums); i++ {
		fmt.Println("Integer stored in buffer:", i, "=", binary.LittleEndian.Uint64(intArrayBuf[(i*int(unsafe.Sizeof(int(0)))):((i+1)*int(unsafe.Sizeof(int(0))))]))
	}

	// Remember to release memory when finished
	defer arena.DecRef(intBuf)
	defer arena.DecRef(floatBuf)
	defer arena.DecRef(stringBuf)
	defer arena.DecRef(personBuf)

	// Try the session of a SERVER CONNECTION

	// Allocate memory for the session
	sessionBuf := arena.Alloc(128)

	// Use the buffer as a Session
	//session := (*Session)(unsafe.Pointer(&sessionBuf[0]))

	// Initialize the session
	(*Session)(unsafe.Pointer(&sessionBuf[0])).ID = "fghjkjhgfdsdfghj"
	(*Session)(unsafe.Pointer(&sessionBuf[0])).UserID = "session"
	(*Session)(unsafe.Pointer(&sessionBuf[0])).Numbers = append((*Session)(unsafe.Pointer(&sessionBuf[0])).Numbers, 10)
	// Manage the connection
	fmt.Println("Allocated memory for person:", sessionBuf)
	fmt.Println((*Session)(unsafe.Pointer(&sessionBuf[0])).UserID)
	fmt.Println((*Session)(unsafe.Pointer(&sessionBuf[0])).ID)
	fmt.Println((*Session)(unsafe.Pointer(&sessionBuf[0])).Numbers)
	// Finally, deallocate memory for the session
	arena.DecRef(sessionBuf)

	// Try PERSON2
	// Allocate memory for the session
	person2Buf := arena.Alloc(int(unsafe.Sizeof(Person2{})))

	// Use the buffer as a Session
	//session := (*Session)(unsafe.Pointer(&sessionBuf[0]))
	// Initialize the session
	(*Person2)(unsafe.Pointer(&person2Buf[0])).Name = "Matteo"
	(*Person2)(unsafe.Pointer(&person2Buf[0])).Surname = "Lando"
	(*Person2)(unsafe.Pointer(&person2Buf[0])).Age = 23
	// Manage the connection
	fmt.Println("Allocated memory for person2:", person2Buf)
	fmt.Println((*Person2)(unsafe.Pointer(&person2Buf[0])).Name)
	fmt.Println((*Person2)(unsafe.Pointer(&person2Buf[0])).Surname)
	fmt.Println((*Person2)(unsafe.Pointer(&person2Buf[0])).Age)
	// Finally, deallocate memory for the session
	arena.DecRef(person2Buf)

	// ARRAY SESSIONS
	// Create a slice of pointers to Session2
	numSessions := 10
	sessions := make([]*Session2, numSessions)

	// Allocate each Session individually in the slab allocator
	for i := 0; i < numSessions; i++ {
		// Allocate memory for a Session
		sessionBuf := arena.Alloc(int(unsafe.Sizeof(Session2{})))

		// Get a pointer to the newly allocated Session
		session := (*Session2)(unsafe.Pointer(&sessionBuf[0]))

		// Set some sample values
		session.ID = i
		copy(session.Name[:], fmt.Sprintf("----------%d", i))
		session.Active = i%2 == 0 // Activate for even indices

		// Save the pointer to the session in the slice
		sessions[i] = session
	}

	// Print the sessions
	for i, session := range sessions {
		fmt.Printf("Session %d: ID = %d, Name = %s, Active = %v\n",
			i, session.ID, session.Name, session.Active)
	}

	// NEW
	// Create a slice of Loc's to track the locations of the buffers.
	locSlice := make([]slab.Loc, 10)

	// Allocate each Data instance and track its location in locSlice.
	for i := range locSlice {
		buf := arena.Alloc(int(unsafe.Sizeof(Data{})))

		dataBuffer := (*Data)(unsafe.Pointer(&buf[0]))
		dataBuffer.name = fmt.Sprintf("Name%d", i)
		dataBuffer.age = i
		locSlice[i] = arena.BufToLoc(buf)
	}

	for _, loc := range locSlice {
		data := (*Data)(unsafe.Pointer(&arena.LocToBuf(loc)[0]))
		fmt.Println(data.name, data.age)
	}

	// Now DecRef each element in the slice using the Loc's we saved earlier.
	for _, loc := range locSlice {
		buf := arena.LocToBuf(loc)
		arena.DecRef(buf)
	}

	// TEST concatenation of buffers of different sizes
	// 1. Creation of an arena with an initial chunk size of 10,
	// slab size of 1000, and a growth factor of 2.
	stats := make(map[string]int64)
	arenaConcatTest := slab.NewArena(10, 100, 2.0, nil)

	// Allocation of two buffers.
	buf1 := arenaConcatTest.Alloc(10) // Allocate a buffer of size 10.
	buf2 := arenaConcatTest.Alloc(15) // Allocate a buffer of size 15.

	copy(buf1, "Hello")
	copy(buf2, " World!")

	// 2. Connect buf1 to buf2.
	arenaConcatTest.SetNext(buf1, buf2)

	// 3. Retrieve the buffer following buf1.
	nextBuf := arenaConcatTest.GetNext(buf1)
	fmt.Println(string(buf1) + string(nextBuf)) // print "Hello World!"

	utility.PrintStats2(arenaConcatTest, stats)
	// Don't forget to free the memory when you're done.
	arenaConcatTest.DecRef(buf1)
	arenaConcatTest.DecRef(buf2)
	arenaConcatTest.DecRef(nextBuf) // Since GetNext increases the reference count.

}

type Data struct {
	name string
	age  int
}

// Writes string characters to the byte buffer
func writeStringToBuffer(buffer []byte, str string) {
	for i := 0; i < len(str); i++ {
		buffer[i] = str[i]
	}
}
