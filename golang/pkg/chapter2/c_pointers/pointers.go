package pointers

import (
	"flag"
	"fmt"
	"math"
	"strings"
)

func firstPointer() {
	x := 1              // Short declaring the variable
	p := &x             // Getting a pointer to x
	xValue := *p        // Getting a value of x with its pointer
	fmt.Println(xValue) // Printing current value of x
	*p = 100            // Assignment a new value to x with its pointer
	fmt.Println(x)      // Printing the new value of x
}

func ComparingTwoPointers() {
	var x, y int
	fmt.Println(&x == &x, &x == &y, &x == nil)
}

func ExistenceOfFunctionLocalVariable() {
	pointer1 := functionThatReturnsPointerOfLocalVariable()
	pointer2 := functionThatReturnsPointerOfLocalVariable()
	fmt.Println(pointer1, *pointer1)
	fmt.Println(pointer2, *pointer2)
	fmt.Println(pointer1 == pointer2) // false because each time we create intValue new block of memory is given to it
}

func functionThatReturnsPointerOfLocalVariable() *int {
	intValue := 100 // After fucntion returns, function that's called this one will have access to intValue
	return &intValue
}

func OuterIncrementingWithPointer() {
	value := 1        // This is an alias for pointer
	pointer := &value // This is an alias for value
	incrementingValueWithPointer(pointer)
	fmt.Println(value)
}

func incrementingValueWithPointer(pointer *int) {
	*pointer++
}

func NameAliasingAfterCopyingObjectsOfReferenceType() {

	slice1 := []int{1, 2, 3}
	slice2 := make([]int, 3, 10)
	copy(slice2, slice1) // This statement creates the alias for slice1 (slice1 - pointer)
	fmt.Println("Before changes:\n", slice1, "\n", slice2)

	slice2[0] = 10
	fmt.Println("After changes:\n", slice1, "\n", slice2)

}

// Echo() prints command-line arguments
// We can find out what "n" and "s" do with calling go run main.go -help
// So that set command line arguments we shoul write -n {true/false} or -s {separator}
var newLine = flag.Bool("n", false, "omit trailing newline")
var separator = flag.String("s", " 0===3 ", "separator")

func Echo() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *separator))
	if !*newLine {
		fmt.Println()
	}
}

func CreatingPointerWithNew() {
	pointer1 := createNewIntPointer()
	pointer2 := createNewIntPointer()
	fmt.Println(pointer1 == pointer2) // They're different to each other
}

func createNewIntPointer() *int {
	return new(int)
}

type MyStruct struct {
	array [0]int
}

func TheExceptionToInequalityOfPointers() {

	var int1 = new(int)
	var int2 = new(int)
	fmt.Println("int1 == int2:", int1 == int2) // They are definitely different

	var nanValue1 = math.NaN()
	var nanValue2 = math.NaN()
	fmt.Println("nanValue1 == nanValue2:", nanValue1 == nanValue2) // They are definitely different

	var slice1 = new([]int)
	var slice2 = new([]int)
	fmt.Println("slice1", slice1, "slice2", slice2)
	fmt.Println("slice1 == slice2:", slice1 == slice2) // They are different because memory is allocated for
	// each slice independetly

	/* We can't expand this initialized variables, so compiler creates one pointer to each of them */
	var arr1 = new([0]int)
	var arr2 = new([0]int)
	fmt.Println("arr1", arr1, "arr2", arr2)
	fmt.Println("arr1 == arr2:", arr1 == arr2) // They are equal. Because empty objects (0 elements) are equal

	var object1 = new(MyStruct)
	var object2 = new(MyStruct)
	fmt.Println("object1", object1, "object2", object2)
	fmt.Println("object1 == object2:", object1 == object2) // They are equal

}

func unavailableNewFunctionBecauseOfRedefining(old, new int) int {
	// var value = new(int) // new will be found by compiler and marked as an error
	return old - new
}
