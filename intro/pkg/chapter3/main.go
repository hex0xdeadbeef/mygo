package chapter3

import (
	"fmt"
)

var a string = "Hi!"

func All_Methods_Using() {

	var x string
	x = "Hello, world"
	fmt.Println(x)
	x = "second"
	// x += "second"
	fmt.Println(x)

	var y string = "hello"
	var z string = "world"
	// var z string = "hello"
	fmt.Println(y == z)

	// var a string = "Hello world"
	b := 5
	fmt.Println(b)

	fmt.Println(a)
	f()

	const c string = "Hello, world!"
	fmt.Println(c)

	var (
		q int     = 5
		w int     = 10
		e float32 = 15
	)

	var (
		r = 20
		t = 25
		u = 30
	)

	fmt.Println(q, w, e)
	fmt.Println(r, t, u)

	userInput()

	farenheitToCelsius()
	feetToMeters()
}

func f() {
	fmt.Println(a)
}

func userInput() {
	fmt.Println("Enter a number")
	var input float64
	fmt.Scanf("%f", &input)

	output := input * 2

	fmt.Println(output)
}

func farenheitToCelsius() {
	fmt.Println("Enter a farenheit: ")
	var farenheotDegreeInputByAUser float64
	fmt.Scanf("%f", &farenheotDegreeInputByAUser)

	celsius := (farenheotDegreeInputByAUser - 32.0) * 5 / 9

	fmt.Println(celsius)
}

func feetToMeters() {
	fmt.Println("Enter the meters value you want: ")
	var metersInput float64
	fmt.Scanf("%f", &metersInput)

	feet := metersInput * 0.3048

	fmt.Println("Result is: ", feet)
}
