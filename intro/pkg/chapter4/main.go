package chapter4

import (
	"fmt"
)

func All_Methods_Using() {
	straightForwardOutput()
	forOutput()
	modifiedForOutput()
	ifStatementPresence()
	caseStatementPresence()
	devisibleByThreeNumbers()
	fizzAndBuzz()
}

func straightForwardOutput() {
	fmt.Println(`1
	2
	3
	4
	5
	6
	7
	8
	9
	10
	`)
}

func forOutput() {
	i := 1
	for i <= 10 {
		fmt.Print(i, " ")
		i++
	}
	fmt.Println()
}

func modifiedForOutput() {
	for i := 1; i <= 10; i++ {
		fmt.Println(i)
	}
	fmt.Println()
}

func ifStatementPresence() {
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			fmt.Print(" ", true)
		} else {
			fmt.Print(" ", false)
		}
	}
	fmt.Println()
}

func caseStatementPresence() {
	for i := 0; i <= 11; i++ {
		switch i {
		case 0:
			fmt.Println("zero")
		case 1:
			fmt.Println("one")
		case 2:
			fmt.Println("two")
		case 3:
			fmt.Println("tree")
		case 4:
			fmt.Println("four")
		case 5:
			fmt.Println("five")
		case 6:
			fmt.Println("six")
		case 7:
			fmt.Println("seven")
		case 8:
			fmt.Println("eight")
		case 9:
			fmt.Println("nine")
		case 10:
			fmt.Println("ten")
		default:
			fmt.Println("Unknown number")
		}
		fmt.Println()
	}

}

func devisibleByThreeNumbers() {
	for i := 3; i < 100; i++ {
		if i%3 == 0 {
			fmt.Print(" ", i)
		}
	}
	fmt.Println()
}

func fizzAndBuzz() {
	for i := 1; i <= 100; i++ {
		if i%15 == 0 {
			fmt.Println("FizzBuzz ")
		} else if i%3 == 0 {
			fmt.Println("Buzz ")
		} else if i%5 == 0 {
			fmt.Println("Fizz ")
		} else {
			fmt.Println(" ", i)
		}
	}
}
