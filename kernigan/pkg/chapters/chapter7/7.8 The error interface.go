package chapter7

import (
	"errors"
	"fmt"
	"io"
	"syscall"
)

func ErrorsComparison() {
	fmt.Println(errors.New("EOF") == io.EOF) // false
}

func SyscallPackageUsing() {
	var err error = syscall.Errno(2)
	fmt.Println(err.Error()) // no such file or directory
	fmt.Println(err)         // no such file or directory
}
