package chapter9

import (
	"fmt"
	"runtime"
)

// MPB

func HackerCliché(threadsNumber int) {
	runtime.GOMAXPROCS(threadsNumber)

	for {
		fmt.Print(1)
		go fmt.Print(0)
	}
}
