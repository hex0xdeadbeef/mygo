package main

import (
	"fmt"
	"strconv"
)

func main() {
	testB()
}

func testA() {
	// 0x0088 00136 (/Users/dmitriymamykin/Desktop/mygo/concurrencypractice/bin/goroutinestack/main.go:12)     CALL    runtime.morestack_noctxt(SB)
	a := 1
	b := 2

	r := max(a, b)
	fmt.Println(`max: ` + strconv.Itoa(r))
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func testB() {
	var x = make([]int, 10)

	fmt.Println(&x)
	subtestBA(x)
	fmt.Println(&x)
}

func subtestBA(x []int) {
	var y = make([]int, 100)

	fmt.Println(y)

}
