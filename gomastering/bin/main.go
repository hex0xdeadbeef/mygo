package main

import (
	"fmt"
	"gomastering/pkg/ch02/randomgen"
)

func main() {
	pass, err := randomgen.SecureRandPass(100)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(pass)
}
