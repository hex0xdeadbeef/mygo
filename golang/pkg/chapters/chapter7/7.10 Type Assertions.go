package chapter7

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

/*
Basics
*/
func ConcreteTypeAssertionUsing() {
	var w io.Writer
	w = os.Stdout

	f := w.(*os.File)      // success f == os.Stdout
	c := w.(*bytes.Buffer) // panic: "interface conversion: io.Writer is *os.File, not *bytes.Buffer"

	fmt.Println(f, c)

}

func InterfaceAssertionUsing() {
	var w io.Writer
	w = os.Stdout
	rw := w.(io.ReadWriter) // success: *os.File has both Read and Write methods

	w = new(ByteCounter)
	rw = w.(io.ReadWriter) // panic: interface conversion: *chapter7.ByteCounter is not io.ReadWriter: mi...

	var vault []byte
	rw.Read(vault)
}

func NilInterfaceAssertion() {
	var w io.Writer // dynamic type -> nil, dynamic value -> nil

	rw := w.(io.ReadWriter) // panic: "interface conversion: interface is nil, not io.ReadWriter"

	var vault []byte
	rw.Read(vault)
}

func DowncastInterfaceType() {
	var rw io.ReadWriter
	rw = os.Stdout

	w := rw.(io.Writer) // fails only if rw is nil

	w.Write([]byte("Hello!"))
}

func TypeAssertionValidityCheck() {
	var w io.Writer = os.Stdout

	f, ok := w.(io.ReadWriteCloser)
	if ok {
		fmt.Printf("The variable \"f\" is now of %T type.", f)
	}

	b, ok := w.(*bytes.Buffer) // The dynamic type of the w is *os.File, and it doesn't match with *bytes.Buffer
	if ok {
		fmt.Printf("The variable \"b\" is now of %T type.", b)
	}
}

func ShadowingVariableInAssertionCheck() {
	var w io.Writer = os.Stdout

	fmt.Printf("The address of the original \"%v\"\n", &w)

	if w, ok := w.(*os.File); ok {
		fmt.Printf("The address of the shadowing variable \"%v\"\n", &w)
		w.Write([]byte("Hello!"))
	}
}
