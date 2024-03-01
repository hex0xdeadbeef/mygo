package main

import (
	"fmt"
	"reflect"
)

func main() {
	valueReversing()
}
func valueOfUsing() {
	const (
		three = 3
	)

	threeValue := reflect.ValueOf(three)
	fmt.Println(threeValue.String()) // <int Value>

	strValue := reflect.ValueOf("hello")
	fmt.Println(strValue) // hello
}

func typeOfValue() {
	const (
		number = 100
	)

	val := reflect.ValueOf(number)
	fmt.Println(val.Type())
}

func valueReversing() {
	const (
		number = 100
	)

	value := reflect.ValueOf(number)
	interfaceValue := value.Interface()

	fmt.Println(interfaceValue.(int))
}
