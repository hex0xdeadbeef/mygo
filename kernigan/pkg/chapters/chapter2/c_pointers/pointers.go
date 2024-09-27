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
	var x, y int                               // They're both nil initially
	fmt.Println(&x == &x, &x == &y, &x == nil) // true false false
}

func ExistenceOfFunctionLocalVariable() {
	pointer1 := getLocalVarPointer()
	pointer2 := getLocalVarPointer()
	fmt.Println(pointer1, *pointer1)  // First memory address
	fmt.Println(pointer2, *pointer2)  // Second memory address
	fmt.Println(pointer1 == pointer2) // false because each time a new address is created
}

func getLocalVarPointer() *int {
	intValue := 100 // After fucntion returns, function that's called this one will have access to intValue
	return &intValue
}

func OuterIncrementing() {
	value := 1        // This is an alias for pointer
	pointer := &value // This is an alias for value
	incrementWithPointer(pointer)
	fmt.Println(value) // 2
}

func incrementWithPointer(pointer *int) {
	*pointer++
}

// Echo() prints command-line arguments
// We can find out what "n" and "s" do with calling go run main.go -help
// So that set command line arguments we should write -n={true/false} or -s={separator}
var newLine = flag.Bool("n", false, "omit trailing newline") // the type is *bool
var separator = flag.String("s", " 0===3 ", "separator")     // the type is *string

func Echo() {
	flag.Parse()                                     // Reads all flags
	fmt.Print(strings.Join(flag.Args(), *separator)) // If flag present, it'll be applied
	if !*newLine {                                   // If flag present, it'll be applied
		fmt.Println()
	}
}

func CreatingPointerWithNew() {
	pointer1 := createPointerWithNew()
	pointer2 := createPointerWithNew()
	fmt.Println(pointer1 == pointer2) // They're different to each other
}

func createPointerWithNew() *int {
	return new(int) // Returns a unique pointer to int value
}

func EmptyObjectsComparison() {

	type MyStruct struct {
		array [0]int
	}

	var (
		int1 = new(int)
		int2 = new(int)

		nanValue1 = math.NaN()
		nanValue2 = math.NaN()

		z float64 // This variable is used for infinities creation

		slice1 = new([]int)
		slice2 = new([]int)

		arr1 = new([0]int)
		arr2 = new([0]int)

		object1 = new(MyStruct)
		object2 = new(MyStruct)
	)

	fmt.Println("int1 == int2:", int1 == int2) // false

	fmt.Println("nanValue1 == nanValue2:", nanValue1 == nanValue2) // false

	fmt.Println(z)
	fmt.Println("inf == inf:", 1.0/z == 1.0/z) // true

	fmt.Println("slice1", slice1, "slice2", slice2)
	fmt.Println("slice1 == slice2:", slice1 == slice2) // false. They are different because memory is allocated for
	// each slice independetly

	/* We can't expand this initialized variables, so compiler creates one pointer to each of them */
	fmt.Println("arr1", arr1, "arr2", arr2)
	fmt.Println("arr1 == arr2:", arr1 == arr2) // true They are equal. Because empty objects (0 elements) are equal

	fmt.Println("object1", object1, "object2", object2)
	fmt.Println("object1 == object2:", object1 == object2) // true They are equal
}

func unavailableNew(old, new int) int {
	// var value = new(int) // new will be found by compiler and marked as an error
	return old - new
}
