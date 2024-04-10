package chapter6

import (
	"fmt"
	"math/rand"
)

func evenOrOdd(number int) (int, bool) {
	number /= 2
	return number, number%2 != 0
}

func greatestNumberInListOfNumbers(arguments ...int) int {
	greatestNumber := arguments[0]
	for i := 0; i < len(arguments); i++ {
		if greatestNumber < arguments[i] {
			greatestNumber = arguments[1]
		}
	}
	return greatestNumber
}

func makeOddGenerator() func() int {
	nextReturningValue := 1
	return func() (ret int) {
		returnableNumber := nextReturningValue
		nextReturningValue += 2
		return returnableNumber
	}
}

func getFibonacci(number int) int {
	if number < 0 {
		return 0
	} else if number == 1 {
		return 1
	} else {
		return getFibonacci(number-1) + getFibonacci(number-2)
	}
}

func solutionOfPanicCase() {

	defer panicHandler()

	var roundOfNumberTen []int = make([]int, 10, 15)

	for i := 0; i <= cap(roundOfNumberTen); i++ {

		numberFactory := func() (int, bool) {
			generatedNumber := rand.Intn(10)
			if generatedNumber == 0 {
				return generatedNumber, true
			} else {
				return generatedNumber, false
			}
		}
		var numberToAppend int
		var errorPointer bool
		numberToAppend, errorPointer = numberFactory()

		if errorPointer {
			panic("попытка добавить нуль в слайс.")
		} else {
			roundOfNumberTen = append(roundOfNumberTen, numberToAppend)
		}
	}
}

func panicHandler() {
	fmt.Println()
	str := recover()
	fmt.Println("При выполнении программы произошла ошибка:", str)
}

func getPointerOfFirstNumberOfArray(arguments ...float64) *float64 {
	return &arguments[0]
}

func getRandomPointerToMemoryPoint() *int {
	return new(int)
}

func squareNumber(numberToBePowed *float64) {
	*numberToBePowed *= *numberToBePowed
}

func usingPowFunction(floatNumber float64) {
	squareNumber(&floatNumber)
}

func swapUsingPointers(firstNumber *float64, secondNumber *float64) {
	var bucket float64 = *firstNumber
	*firstNumber = *secondNumber
	*secondNumber = bucket
}
