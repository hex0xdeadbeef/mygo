package declarationinitialization_test

import (
	"fmt"
	"io"
	"math"
	"os"
)

func NilInitialValueTypes() {
	var myInterfaceVar io.Reader   // Initial value is nil
	var sliceVar []int             // Initial value is nil
	var pointerVar *int            // Initial value is nil
	var mapVar map[int]int         // Initial value is nil
	var channelVar chan string     // Initial value is nil
	var functionVar func() float64 // Initial value is nil

	fmt.Println(myInterfaceVar, sliceVar, pointerVar, mapVar, channelVar, functionVar)

}

func multipleDeclaration() {
	var i, j, k int                 // All are int
	var b, f, s = true, 2.3, "four" // bool, float64, string
	fmt.Println(i, j, k, b, f, s)
}

func InitializingWithFunctionCall() {
	file, err := os.Open("text.txt")
	if err != nil {
		os.Exit(1)
		return
	}

	io.Copy(os.Stdout, file)
	fmt.Println()

	defer file.Close()
}

func shortVariableDeclaration() {
	i, j := 1, 0
	fmt.Println(i, j)
}

func tupleAssignment() {
	i, j := 1, 0
	i, j = j, i
	fmt.Println(i, j)
}

func shortDeclaringAsAssignment() {
	var value float64 = 100

	doubleSuare := func(value float64) (float64, float64) {
		return math.Pow(float64(value), 2.0), math.Pow(float64(value), 2.0)
	}

	value, theSameValue := doubleSuare(value) // value will be assigned with the squared value
	// Value will be just rewritten

	fmt.Println(value, theSameValue)
}

func shortDeclaringRepeat(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	// file, err := os.Open(fileName) // There must be just assignment

	// Like this
	file, err = os.Open(fileName)
	defer file.Close()
}

func DeclaringInInnerLexicalBlock(fileName string) {
	var file *os.File
	var err error
	fmt.Printf("Old values: %s, %v\n", file, err)
	{
		file, err := os.Open(fileName) // These ones are local variables
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer file.Close()

		fmt.Println("File data: ")
		num, err := io.Copy(os.Stdout, file)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			fmt.Println()
			fmt.Println(num)
		}
	}
	// The upper variables still have initial values (nil, nil)
	fmt.Printf("New values: %s, %v\n", file, err)
}
