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
1.The length is the part of an array type and it must be computated at compiling time
	1) After declaring an array, it'll filled with zero elements of the corresponding type.

2. We can map int keys (indices) to values of the target type expicitly.
	1) Due to the fact of mapping we can omit following the order of indices

	2) Keys can be enumerator's value with int underlying type.

	3) The max value of indice defines a len of array.

3. If elements of an array are comparable and theirs lengths too (types are different [n]type and [m]type) have the same
value, that's in itself is comparable => we can compare two arrays using == which reports whether all corresponding elements
are equal. The != operator is negation.

4. %x verb allows to print all the elements of an array or slice of bytes in hexadecimal other verbs work as well.

5. When a function is called a copy of each argument value is assigned to the corresponding parameter variable, so the
function receives a copy, not the original array.

6. Passing large arrays into a function is inefficient, because of copying is required. Also no changes will be
reflected on the original array. So that defeat this problem we can pass a pointer to the original array as an argument of a
function.

7. Arrays can't contain itself. Combined declaration for them is forbidden.
*/

func Intoduction() {
	var a [3]int
	fmt.Println(a[0])
	fmt.Println(a[len(a)-1])

	for index, value := range a {
		fmt.Printf("index: %d \t| value: %d\n", index, value)
	}

	for _, value := range a {
		fmt.Printf("value: %d\n", value)
	}
}

func ArrayLiteralInitialization() {
	var q [3]int = [3]int{1, 2, 3}
	var r [3]int = [3]int{1, 2}

	fmt.Printf("q[2]: %d\t|r[2]: %d", q[2], r[2])
}

func EllipsisInitialization() {
	q := [...]int{1, 2, 3}
	fmt.Printf("% d, ", q)
}

func SpecifyPairsArray() {
	type Currency int

	const (
		USD Currency = iota
		EUR
		GBP
		RMB
	)

	symbol := [...]string{USD: string('\U00000024'), EUR: string('\U000020AC'), GBP: string('\U000000A3'), RMB: string('\U000000A5')}
	fmt.Printf("% v \n", symbol)
	fmt.Println(RMB, symbol[RMB])

	r := [...]int{99: -1}
	fmt.Printf("% d", r)

}

func ArrayComparing() {
	a := [2]int{1, 2}
	b := [...]int{1, 2}
	c := [2]int{1, 3}

	// d := [3]int{1, 2}
	fmt.Printf("a == b: %v\n", a == b)
	fmt.Printf("b == c: %v\n", b == c)
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
		[32]uint8
	*/

}

func ZeroArray(pointer *[32]byte) {
	for i := range pointer {
		pointer[i] = 0
	}
}

func ZeroArraySimple(pointer *[32]byte) {
	*pointer = [32]byte{}
}

func CountDifferentBits(firstInput, secondInput string) int {
	firstStringBytes, secondStringBytes := sha256.Sum256([]byte(firstInput)), sha256.Sum256([]byte(secondInput))

	countOfDifferentBits := 0

	for i := 0; i < len(firstStringBytes); i++ {
		countOfDifferentBits += int(
			math.Abs(
				float64(countUsedBits(firstStringBytes[i])) - float64(countUsedBits(secondStringBytes[i]))))
	}

	return countOfDifferentBits
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
	flag.Parse()

	var userInput []string = os.Args[1:] // The first argument of the os.Args[] if an executable file

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
		default:
			fmt.Fprintf(os.Stdout, "The sha256 hash(es) are:\n")
			for _, arg := range userInput[2:] {
				fmt.Fprintf(os.Stdout, "| string: %s\t| hash: %x\n", arg, sha256.Sum256([]byte(arg)))
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "The error occured: no args typed\n")
	}

}
