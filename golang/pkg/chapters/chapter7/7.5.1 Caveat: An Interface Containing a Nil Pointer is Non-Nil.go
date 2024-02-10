package chapter7

import (
	"bytes"
	"fmt"
	"io"
)

const debug = false

func FUsing() {
	/* var buf *bytes.Buffer // it causes the panic */
	var buf io.Writer

	if debug {
		buf = new(bytes.Buffer)
	}
	f(buf) // In case "var buf *bytes.Buffer" we pass into the function a non-nil Writer
	// because it's type descriptor: *bytes.Buffer, value: <nil>
}

func f(out io.Writer) {
	/*
		... do something ...
	*/
	fmt.Printf("Type: %T | Value: %v\n", out, out)

	if out != nil {
		out.Write([]byte("Done!\n"))
	}
}
