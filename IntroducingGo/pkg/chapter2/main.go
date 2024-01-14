package chapter2

import "fmt"

func All_Methods_Using() {
	fmt.Println("1+1", 1+1)
	fmt.Println("1+1", 1.0+1.0)

	fmt.Println(len("Hello, world"))
	fmt.Println("Hello, world"[1])
	fmt.Println("Hello, " + "world")

	fmt.Println(true && true)
	fmt.Println(true && false)
	fmt.Println(true || true)
	fmt.Println(true || false)
	fmt.Println(!true)

}
