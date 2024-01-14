package chapter3

import (
	"fmt"
	"math"
)

func BitsOverflow() {
	var u uint8 = 255
	fmt.Println(u, u+1, u*u) // u+1: 1111 1111 + 0000 0001 = 0 | u*u : 255^2(mod256)

	var i int8 = 127
	fmt.Println(i, i+1, i*i) // i + 1: 1111 1111 + 0000 0001 = 1000 0001 when the senior bit is the sign of the value and
	// other bit are the positive number that is summed to -128 | i*i: -128 + (127*127 mod 128))

}

func incompatibleTypesOperation() {
	var apples uint32 = 1
	var oranges uint16 = 2
	// var compote int = apples + oranges it'll be considere by compiler as error
	var compote int = int(uint(oranges) + uint(apples)) // This is implicit conversion
	fmt.Println(compote)
}

func PrecisionLoseOperations() {
	f := 3.141  // its type is float64
	i := int(f) // Now it has value 3
	fmt.Println(f, i)
	f = 1.99
	fmt.Println(f) // Analogically
}

func OperandIsOutOfRangeForTargetType() {
	f := 1e100 // float64
	i := int(f)
	fmt.Printf("%e := 2pow(%.5g) | int(f) = 2pow(%g)", f, math.Log2(f), math.Log2(float64(i))) /* After conversion
	a float64 value that is biggger than the top boundary of int32 type to int, the result value will have the top
	boundary of int64 typetop boundary value */
}

func HexadecimalOctalDecimalPrinting() {
	oc := 0666
	/*The [1] "adverb" tells the compiler to use the first operand over and over again*/
	fmt.Printf("%d %[1]o %#[1]o\n", oc)
	hex := int64(0xdeadbeef)
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", hex)

}

func RunePrinting() {
	ascii := 'a'
	unicode := '当'
	newline := '\n'
	fmt.Printf("%d %[1]c %#[1]q\n", ascii)
	fmt.Printf("%d %[1]c %#[1]q\n", unicode)
	fmt.Printf("%d %[1]c %#[1]q\n", newline)
}
