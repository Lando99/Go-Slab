package main

import (
	"arena"
	"fmt"
	"strconv"
)

type T struct {
	val int
}

func main() {
	// Create an arena in the beginning of the function.
	mem := arena.NewArena()
	// Free the arena in the end.
	defer mem.Free()

	// Allocate a bunch of objects from the arena.
	for i := 0; i < 10; i++ {
		obj := arena.New[T](mem)
		obj.val = 1
		fmt.Println(strconv.Itoa(obj.val))
	}

	fmt.Println("hello")
	// Or a slice with length and capacity.
	slice := arena.MakeSlice[T](mem, 100, 200)
	_ = slice
}
