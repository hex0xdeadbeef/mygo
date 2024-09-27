package chapter7

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

func InterfaceWriterValues() {

	var w io.Writer // type -> nil, value -> nil
	fmt.Println(w == nil)
	fmt.Printf("%T\n", w)
	/* w.Write([]byte("Hello!")) // "runtime error: invalid memory address or nil pointer dereference" */
	fmt.Println()

	/*
		w = io.Writer(os.Stdout)
	*/
	w = os.Stdout         // type -> *os.File, value -> fd int = 1(stdout)
	fmt.Println(w == nil) // false
	fmt.Printf("%T\n", w)
	w.Write([]byte("Hello!\n"))
	fmt.Println()

	w = new(bytes.Buffer) // type -> *bytes.Buffer, value -> [ data []byte ... ]
	fmt.Println(w == nil)
	fmt.Printf("%T\n", w)
	w.Write([]byte("Hello!\n"))
	fmt.Println()

	w = nil
	fmt.Println(w == nil)
	fmt.Printf("%T\n", w)
	/* w.Write([]byte("Hello!")) // "runtime error: invalid memory address or nil pointer dereference" */
	fmt.Println()

	var instant interface{} = time.Now()
	fmt.Printf("%T\n", instant)
	fmt.Println(instant)
}

func InterfaceComparison() {
	var x interface{} = []int{1, 2, 3}
	fmt.Printf("Type: %T | Value: %v\n", x, x)
	// fmt.Println(x == x) // "runtime error: comparing uncomparable type []int"
}
