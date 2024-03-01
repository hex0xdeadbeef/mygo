package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {
	typeOfTest()
}

func typeOfTest() {
	const (
		three = 3
	)
	var (
		readWriter io.ReadWriter = os.Stdout
	)

	threeType := reflect.TypeOf(three)
	fmt.Println(threeType)

	rdtype := reflect.TypeOf(readWriter)
	fmt.Println(rdtype)

}
