package chapter2

import (
	"bufio"
	"flag"
	"fmt"
	lc "golang/pkg/chapter2/convertions/lengthconv" // This is alias
	tc "golang/pkg/chapter2/convertions/tempconv"   // This is alias
	wc "golang/pkg/chapter2/convertions/weightconv" // This is alias
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
) // import is executed by compiler before all package

const boilingF = 212.0

var (
	a, b, c int
)

// init() will initialize global variables after import statement finishes
func init() {
	a = b + c
	b = function()
	c = 1
}

func function() int {
	return c + 1
}

func BoilingCelsiusFarenheit() {
	var farenheit = boilingF
	var celsius = (farenheit - 32) * 5 / 9
	fmt.Printf("boiling point = %.2f F* or %.2f C*\n", farenheit, celsius)
}

func BoilingCelsiusFarenheitReturnable(farenheit float64) float64 {
	return (farenheit - 32) * 5 / 9
}

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

// Allocation on heap/stack doesn't depend on using var/new(...)
// It depends whether variable is reachable or not
var global *int

func f() {
	var x int // It's heap-allocated 'cause it's still reachable from the variable global
	x = 1
	global = &x
}

func g() {
	var y = new(int) // It's stack allocated 'case it won't be reacable after g() has returned
	*y = 10
}

func TupleAssignmentOfThreeVariables() {
	var i, j, k = 2, 3, 5

	fmt.Println(i, j, k)
}

func implicitAssignment() {
	medals := []string{"gold", "silver", "platinum"}
	/*medals[0] = "gold"
	medals[1] = "silver"
	medals[2] = "platinum"
	*/
	fmt.Println(medals[0])
}

type Celsius float64 // This is the "Named Type"

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC     Celsius = 0
	BoilingC      Celsius = 100
)

func (c Celsius) String() string {
	return fmt.Sprintf("%g*C", c)
}

type Farenheit float64

func CToF(c Celsius) Farenheit {
	return Farenheit(c*9/5 + 32)
}

func FToC(f Farenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// Named types behave as their underlying types, but operations between two types that have the same underlying
// types are forbidden (We should converse one type to another)
func NamedTypesStatements() {
	fmt.Printf("%g\n", BoilingC-FreezingC) // This operation is equal to float64-float64 (underlying types are the same)
	boilingF := CToF(BoilingC)
	fmt.Printf("%g\n", boilingF-CToF(FreezingC))
	// fmt.Printf("%g\n" boilboilingF - BoilingC) compile error: types mismatch

	var c Celsius
	var f Farenheit
	fmt.Println(c == 0)
	fmt.Println(f >= 0)
	// fmt.Println(c == f) compile error: types mismatch
	fmt.Println(c == Celsius(f)) // true because of conversion f from Farenheit to Celsius

	c = FToC(212.0)

	// When we declare String() function for a named type, there is no need to call it explicitly
	fmt.Println(c.String())
	fmt.Printf("%v\n", c) // no need to call String() explicitly
	fmt.Printf("%s\n", c)
	fmt.Println(c)
	fmt.Println("%g\n", c)
	fmt.Println(float64(c))
}

func Package_Using() {
	fmt.Println("Brrrr! %v\n", tc.AbsoluteZeroC)
	fmt.Println(tc.CToF(tc.BoilingC))
}

func ParseCmdArgConvertToCorF() {
	for _, arg := range os.Args[1:] {
		temperature, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprint(os.Stderr, "cf %v\n", err)
			os.Exit(1)
		}
		f := tc.Farenheit(temperature)
		c := tc.Celsius(temperature)

		fmt.Printf("%s = %s, %s = %s\n", f, tc.FToC(f), c, tc.CToF(c))
	}
}

func GeneralPurposeUnitConversion() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			unit, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "The error is: %v\n", err)
				os.Exit(1)
			}
			unitsConvertation(unit)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		userInput := scanner.Text()
		userInputArgs := strings.Split(userInput, " ")

		for _, arg := range userInputArgs {
			unit, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				errorPrintAndExit(err)
			}
			unitsConvertation(unit)
		}

	}
}

func errorPrintAndExit(err error) {
	fmt.Fprintf(os.Stderr, "The error is: %v\n", err)
	os.Exit(1)
}

func unitsConvertation(unit float64) {
	farenheit := tc.Farenheit(unit)
	celsius := tc.Celsius(unit)

	kilogram := wc.Kilogram(unit)
	pound := wc.Pound(unit)

	meter := lc.Meter(unit)
	foot := lc.Foot(unit)

	fmt.Printf("%s = %s, %s = %s\n", farenheit, tc.FToC(farenheit), celsius, tc.CToF(celsius))
	fmt.Printf("%s = %s, %s = %s\n", kilogram, wc.KGToP(kilogram), pound, wc.PToKG(pound))
	fmt.Printf("%s = %s, %s = %s\n", meter, lc.MToFt(meter), foot, lc.FtToM(foot))
	fmt.Println()
}

func ScopeShadowingFirst() {
	x := "hello" // "function" scope x, it'll be shadowed by for scope x variable
	for i := 0; i < len(x); i++ {
		x := x[i] // "for" scope x, it'll be shadowed by if scope x
		if x != '!' {
			x := x + 'A' - 'a'  // "if" scope x
			fmt.Printf("%c", x) // this statement uses the last mentioned x (x := x + 'A' - 'a' // "if" scope x)
		}
	}
}

func ScopeShadowingSecond() {
	x := "hello"          // "function" scope x, it'll be shadowed by the "for" scape implicit variable
	for _, x := range x { // The value that is fetched from string x is the new variable
		x := x + 'A' - 'a'  // This x is the new variable that is assigned with expression
		fmt.Printf("%c", x) // This statement uses the last mentioned x (x := x + 'A' - 'a' // "if" scope x)
	}
}

func xRet() int {
	if value := rand.Intn(100); value > 50 {
		return 1
	} else {
		return 0
	}
}

func yRet(x int) int {
	if value := rand.Intn(100); value > 50 {
		return 0
	} else {
		return 1
	}
}

func ScopeShadowingThird() {
	if x := xRet(); x == 0 { // 1. Assignment x a value from xRet(). 2. Evaluating the expression
		fmt.Println(x) // Printing
	} else if y := yRet(x); y == x { /* 1. Assignment y a value from yRet(). 2. Evaluating the expression by using
		the evaluated value of x*/
		fmt.Println(x, y)
	} else {
		fmt.Println(x, y) // Printing both variables by by using the evaluated values of x and y
	}
	// fmt.Println(x, y) // Here's no way to reach both x and y variables 'cause they are out of general if statement
}

func ScopeShadowingFourth(fileNames ...string) {
	if file, err := os.Open(fileNames[0]); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		defer file.Close()
		var fileData []byte
		_, err := file.Read(fileData)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		os.Stdout.Write(fileData)
	}
	// defer file.Close() It will not be accessible 'cause file variable scope has ended above
	// f.ReadByte() -||-

}

var cwd string

func DefeatingShadingOfGlobalVariableAssignment() {

	fmt.Println("BEFORE | The address in memory of cwd:", &cwd)
	/* Here's the local cwd shades the upper declared global cwd */
	// cwd, err := os.Getwd()

	/* Here's a method that defeats this problem */
	var err error
	cwd, err = os.Getwd() // cwd has assigned, but hasn't declared

	// cwd, err := os.Getwd() // Here's the local cwd shades the upper declared global cwd
	if err != nil {
		log.Fatalf("os.Getwd failed:", err)
	}

	fmt.Println("AFTER | The address in memory of cwd:", &cwd, "| The working directory is:", cwd)
}
