package chapter10

import (
	"fmt"
	"runtime"
)

func PrintOSandArch() {
	fmt.Println(runtime.GOOS, runtime.GOARCH)
}
