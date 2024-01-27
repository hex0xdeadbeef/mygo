package main

import (
	"fmt"
	a "golang/pkg/chapter4/b_slices"
)

func main() {

	str := a.SquashSpacesSimple([]byte("\t\t\t\t    aaa bbb www wdqdwqdq   \t\n"))
	fmt.Println(str, len(str), cap(str))
}
