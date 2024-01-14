package chapter6

import "fmt"

func deferFunction() {
	{
		defer two()
		first() // Вывод 1
		three() // Вывод 3
	}
	fmt.Println()
	{
		defer first()
		two()   // Вывод 2
		three() // Вывод 3
	}
	fmt.Println()
}

func first() {
	fmt.Println(1)
}

func two() {
	fmt.Println(2)
}

func three() {
	fmt.Println(3)
}

func panicFunction() {
	defer functionThatIncludesDeferToRecoverExecutionAfterPanic()
	panic("PANIC")
}

func functionThatIncludesDeferToRecoverExecutionAfterPanic() {
	str := recover()
	fmt.Println(str)
}
