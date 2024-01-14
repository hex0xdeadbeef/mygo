package chapter6

import (
	"fmt"
)

func All_Methods_Using() {

	// Simple functions
	fmt.Println("Simple functions")

	var firstArray []float64 = []float64{98, 93, 77, 82, 83}
	fmt.Println(average(firstArray))
	fmt.Println()

	// Functions that return some values
	fmt.Println("Functions that return some values")

	x, y, err := functionWithMultipleReturnableArguments()
	fmt.Println(x, y, err)
	fmt.Println()

	// Functions that use multiple argument
	fmt.Println("Functions that use multiple argument")

	fmt.Println(functionWithModifiedFormOfMultipleArguments(0))
	fmt.Println(functionWithModifiedFormOfMultipleArguments(1, 2, 3, 4, 5))
	sliceOfIntegers := []int{1, 1, 5}
	anotherSliceOfIntegers := []int{}
	fmt.Println(functionWithModifiedFormOfMultipleArguments(sliceOfIntegers...))
	fmt.Println(functionWithModifiedFormOfMultipleArguments(anotherSliceOfIntegers...))
	fmt.Println()

	// Closures
	fmt.Println("Closures")

	fmt.Println(firstFunctionWithClosure(generateRandomSlice(100)...))
	nextEven := MakeEvenGenerator()
	fmt.Println(nextEven())
	fmt.Println(nextEven())
	fmt.Println(nextEven())
	fmt.Println()

	// Recursion
	fmt.Println("Recursion")

	fmt.Println(recursiveFactorialComputation(3))
	fmt.Println()

	// defer, panic, recover
	fmt.Println("defer, panic, recover")

	deferFunction()
	panicFunction()
	fmt.Println()

	// Pointers
	fmt.Println("Pointers")

	firstPointersPresence()
	creatingAndUsingPointerWithNewOperator()
	fmt.Println()

	// Homework
	fmt.Println("Homework")
	fmt.Println(evenOrOdd(3))
	fmt.Println(evenOrOdd(10))
	fmt.Println(evenOrOdd(4))

	oddFactory := makeOddGenerator()
	fmt.Println(oddFactory())
	fmt.Println(oddFactory())
	fmt.Println(oddFactory())
	fmt.Println(oddFactory())

	fmt.Println(getFibonacci(10))

	solutionOfPanicCase()

	secondArray := []float64{1.0, 2.0}
	secondArray = append(secondArray, 3.0)
	pointerAtFirstElement := getPointerOfFirstNumberOfArray(secondArray...)
	valueOfMemoryAtMemoryPoint := *pointerAtFirstElement
	fmt.Println("Адрес в памяти и значение:", pointerAtFirstElement, "|", valueOfMemoryAtMemoryPoint)

	var randomPointer *int = getRandomPointerToMemoryPoint()
	fmt.Println("Значение по рандомному указателю:", *randomPointer)
	*randomPointer = 10
	fmt.Println("Новое значение по рандомному указателю:", *randomPointer)

	usingPowFunction(1.11111111)

	firstNumberToBeSwapped := 1.11111111111111111111111111111
	usingPowFunction(firstNumberToBeSwapped)
	secondNumberToBeSwapped := 2.2222222222222222222222222222
	usingPowFunction(secondNumberToBeSwapped)

	fmt.Println("Before the swap: ", "1:", "|", firstNumberToBeSwapped, "2:", secondNumberToBeSwapped)
	swapUsingPointers(&firstNumberToBeSwapped, &secondNumberToBeSwapped)
	fmt.Println("After the swap: ", "1:", "|", firstNumberToBeSwapped, "2:", secondNumberToBeSwapped)
}
