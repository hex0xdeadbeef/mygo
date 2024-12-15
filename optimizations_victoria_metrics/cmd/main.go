package main

import (
	"fmt"
	"runtime"
)

func main() {
	printAlloc()
	var s []byte

	printAlloc()

	s = []byte{1_000_000: 0}
	s = s[:0]

	printAlloc()

	runtime.GC()

	runtime.KeepAlive(&s)

}

func printAlloc() {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", m.Alloc/1<<20)
}
