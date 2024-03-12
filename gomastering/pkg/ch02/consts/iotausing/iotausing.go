package iotausing

import "fmt"

const (
	p2_0 = 1 << iota
	_
	p2_2
	_
	p2_4
	_
	p2_6
)

func UnskippedValuesUsing() {
	fmt.Println(p2_0)
	fmt.Println(p2_2)
	fmt.Println(p2_4)
	fmt.Println(p2_6)
}
