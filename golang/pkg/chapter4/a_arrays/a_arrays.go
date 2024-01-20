package a_arrays

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"math"
	"os"
)

/*
1.The length is the part of an array type and it must be computable at compiling time

2. We can match int keys (indices) with target type expicitly. Keys can be enumerator's value with int underlying type.
The max value of indice defines a len of array.

3.If elements of an array are comparable and theirs lengths (types are different [n]type and [m]type) have the same
value, that's in itself is comparable => we can compare two arrays using == which reports whether all corresponding
elements are equal. The != operator is negation.

4. %x verb allows to print all the elements of an array or slice of bytes in hexadecimal

5. When a function is called a copy of each argument value is assigned to the corresponding parameter variable, so the
function receives a copy, not the original.

6. Passing large arrays into a function is inefficient, because of copying is required. Also no changes will reflected
on an original array. So that defeat this problem we can pass a pointer to an original array as a argument of a function.
*/

func Intoduction() {
	var a [3]int             // [ 0, 0 , 0]
	fmt.Println(a[0])        // 0
	fmt.Println(a[len(a)-1]) // 0

	// Prints the indices and values
	for index, value := range a {
		fmt.Printf("index: %d \t| value: %d\n", index, value)
	}

	// Prints the values only
	for _, value := range a {
		fmt.Printf("value: %d\n", value)
	}
}

// Initially arrays are filled with zero-values of a base type
func ArrayLiteralInitialization() {
	var q [3]int = [3]int{1, 2, 3} // [1, 2, 3]
	var r [3]int = [3]int{1, 2}    // [1, 2, 0]

	fmt.Printf("q[2]: %d\t|r[2]: %d", q[2], r[2]) // 3 | 0
}

func EllipsisInitialization() {
	q := [...]int{1, 2, 3} // The length will be recongized based on an amount of literals
	fmt.Printf("% d, ", q)
}

func SpecifyPairsArray() {
	type Currency int

	// It can be used as an index because of the underlying int type
	const (
		USD Currency = iota
		EUR
		GBP
		RMB
	)

	// Each index is assigned to a string value
	symbol := [...]string{USD: string('\U00000024'), EUR: string('\U000020AC'), GBP: string('\U000000A3'), RMB: string('\U000000A5')}
	fmt.Printf("% v \n", symbol)  // Only symbols will be printed
	fmt.Println(RMB, symbol[RMB]) // 3 ¥

	r := [...]int{99: -1} // Array that has the length 99 and the last index has -1 value
	fmt.Printf("% d", r)

}

// So that we can compare two arrays they must have: the same types, the same lengths
func ArrayComparing() {
	a := [2]int{1, 2}
	b := [...]int{1, 2}
	c := [2]int{1, 3}

	// d := [3]int{1, 2}
	fmt.Printf("a == b: %v\n", a == b) // true
	fmt.Printf("b == c: %v\n", b == c) // false
	// The following array has the different length unlike others ones => comparison is impossible
	// fmt.Printf("a == d: %v", a == d) invalid operation: a == d (mismatched types [2]int and [3]int)
}

// If encoded input strings digests are equal, it's extremely likely that the to messages are the same
func DigestsComparison() {
	// The types of c1 and c2 are [32]byte, because of the size of the crypto digest (256 bits => 32 bytes)
	c1 := sha256.Sum256([]byte("x"))
	c2 := sha256.Sum256([]byte("X"))

	fmt.Printf("%x\n%x\n%t\n%T\n", c1, c2, c1 == c2, c1)
	/*
		2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881 the digest(hash) of c1
		4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015 the digest(hash) of c2
		false // the result of comparison of two arrays
		[32]uint8 // Inspite of the initial type, the result type is [32]uint8 (byte is the alias of uint)
	*/

}

func ZeroArray(pointer *[32]byte) {
	// The pointer is just the beginning of an array
	for i := range pointer {
		pointer[i] = 0
	}
}

func ZeroArraySimple(pointer *[32]byte) {
	*pointer = [32]byte{} // We just assign the {0, ..., 0} array to *pointer
}

func CountDifferentBits(firstInput, secondInput string) int {
	// Get bytes for each input string
	firstStringBytes, secondStringBytes := sha256.Sum256([]byte(firstInput)), sha256.Sum256([]byte(secondInput))

	// Declare and initialize the result value of the different bytes amount
	total := 0

	// Run through the arrays and count a current bytes difference amount summing total up
	for i := 0; i < len(firstStringBytes); i++ {
		total += int(
			math.Abs(
				float64(countUsedBits(firstStringBytes[i])) - float64(countUsedBits(secondStringBytes[i]))))
	}

	return total
}

func countUsedBits(b byte) int {
	total := 0
	for b != 0 {
		b &= (b - 1)
		total++
	}
	return total
}

var shaPointer = flag.String("sha", "256", "choose hash flavor")

func PrintDigest() {
	// Read all flags
	flag.Parse()

	// Read user input
	var userInput []string = os.Args[1:]

	if len(userInput) != 0 {
		switch *shaPointer {
		case "384":
			fmt.Fprintf(os.Stdout, "The sha384 hash(es) are:\n")
			for _, arg := range userInput[2:] {
				fmt.Fprintf(os.Stdout, "| string: %s\t| hash: %x\n", arg, sha512.Sum384([]byte(arg)))
			}
		case "512":
			fmt.Fprintf(os.Stdout, "The sha512 hash(es) are:\n")
			for _, arg := range userInput[2:] {
				fmt.Fprintf(os.Stdout, "| string: %s\t| hash: %x\n", arg, sha512.Sum512([]byte(arg)))
			}
		case "256":
			fmt.Fprintf(os.Stdout, "The sha256 hash(es) are:\n")
			for _, arg := range userInput[2:] {
				fmt.Fprintf(os.Stdout, "| string: %s\t| hash: %x\n", arg, sha256.Sum256([]byte(arg)))
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "The error occured: no args typed\n")
	}

}
