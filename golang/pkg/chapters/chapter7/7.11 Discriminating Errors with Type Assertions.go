package chapter7

import (
	"errors"
	"fmt"
	"os"
)

func UsingErrorTypeAssertion() {
	_, err := os.Open("/no/such/file")
	fmt.Println(err) // open /no/such/file: no such file or directory

	fmt.Printf("%#v\n", err) // &fs.PathError{Op:"open", Path:"/no/such/file", Err:0x2}
}

func DistinguishingErrorsWithAssertion() {
	var ErrNotExist = errors.New("file doesn't exist")

	/*
		IsNotExist returns a boolen indicating whether the error is known to report that a file or directory doen't exist.
	*/

	// var IsNotExist func(err error) bool
	// isNotExist = func(err error) bool {
	// 	if pe, ok := err.(*os.PathError); ok {
	// 		err = pe.Err
	// 	}
	// 	return err == syscall.ENOENT
	// }

	_, err := os.Open("no/such/file")

	fmt.Println(os.IsNotExist(ErrNotExist)) // false
	fmt.Println(os.IsNotExist(err))         // true
}
