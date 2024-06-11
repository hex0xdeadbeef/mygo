package main

/*
1. init() cannot be invoked directly as a simple function.

*/

import (
	"fmt"

	_ "hundred_issues/internal/incorrectinitfunc/foo"
)

// 1.
var (
	a = func() int {
		fmt.Println("var done")
		return -1
	}()
)

// 2.
func init() {
	fmt.Println("init 2 done")

}

// 3.
func init() {
	fmt.Println("init 1 done")
}

// 4.
func main() {
	fmt.Println("main done")

	// init() // undefined: initcompilerUndeclaredName

}
