package main

import (
	"runtime"
	"time"
)

func main() {

	runtime.GOMAXPROCS(1)

	go A()

	time.Sleep(1 * time.Second)
}

func A() {
	for {
	}
}
