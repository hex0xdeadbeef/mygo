package main

import (
	"fmt"
	"golang/pkg/chapters/chapter9"
)

func main() {
	chapter9.DepositSafe(200)
	fmt.Println(chapter9.Withdraw(100))
}
