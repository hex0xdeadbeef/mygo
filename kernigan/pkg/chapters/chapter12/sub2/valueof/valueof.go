package main

import (
	"fmt"
	"reflect"
)

const (
	hundred = 100
)

func main() {
	valueOfUsing()
	fmt.Println()
	valueReversing()
}

func valueOfUsing() {
	hundredVal := reflect.ValueOf(hundred)
	fmt.Println(hundredVal.String()) // <int Value>

	nilVal := reflect.ValueOf(nil)
	fmt.Println(nilVal.String())

	strValue := reflect.ValueOf("hello")
	fmt.Println(strValue) // hello
}

func valueReversing() {
	hundredVal := reflect.ValueOf(hundred)
	fmt.Println(hundredVal.String())

	interfaceValue := hundredVal.Interface()

	fmt.Println(interfaceValue.(int))
}
