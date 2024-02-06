package e_functions

import (
	"fmt"
	"os"
	"runtime"
)

/*
Defer function stack:
0. defer 1
1. defer 2
2. defer 3
*/
func PressF(x int) {
	fmt.Printf("f(%d)\n", x+0/x)
	defer fmt.Printf("defer %d\n", x)
	PressF(x - 1)
}

/*
The function holds the full information about called functions. They're piled in the reverse order. And the last
function call will be at the top of the stack.
*/
func PrintStack() {
	var buf []byte
	// false - provides the result with information only about only the current goroutine trace
	// true - provides extra information about all other goroutines after the trace of the current one
	numberOfBytesWritten := runtime.Stack(buf, true)
	os.Stderr.Write(buf[:numberOfBytesWritten])
}
