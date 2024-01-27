package floats

import (
	"fmt"
	"math"
)

func PrecisionLose() {
	f := 3.141        // its type is float64
	i := int(f)       // Now it has value 3
	fmt.Println(f, i) // 3.141. 3
	f = 1.99
	fmt.Println(f) // 1.99
}

func OperandIsOutOfRangeForTargetType() {
	f := 1e100 // float64
	i := int(f)
	fmt.Printf("%e := 2pow(%.5g) | int(f) = 2pow(%g)", f, math.Log2(f), math.Log2(float64(i))) /* After conversion
	a float64 value that is biggger than the top boundary of int32 type to int, the result value will have the top
	boundary of int64 typetop boundary value */
}

func PrintHexOctDec() {
	oc := 0666 // Oct prefix is 0
	/*The [1] "adverb" tells the compiler to use the first operand over and over again*/
	fmt.Printf("%d %[1]o %#[1]o\n", oc)
	hex := int64(0xdeadbeef)
	fmt.Printf("%d %[1]x %#[1]x %#[1]X\n", hex)
}

func RunePrinting() {
	ascii := 'a'                           // It's an int value that's an equivalent to a corresponding unicode code point
	unicode := '当'                         // It's an int value that's an equivalent to a corresponding unicode code point
	newline := '\n'                        // It's an int value that's an equivalent to a corresponding unicode code point
	fmt.Printf("%d %[1]c %#[1]q\n", ascii) // Verb "q" means that printed value will be surrounded with
	// double quotes
	fmt.Printf("%d %[1]c %#[1]q\n", unicode)
	fmt.Printf("%d %[1]c %#[1]q\n", newline)
}

func MaxMinNonzeroFloatsValues() {
	fmt.Printf("Float32 min value is: 2pow(%.5g) | max: 2pow(%.5g)\n", math.Log2(math.SmallestNonzeroFloat32), math.Log2(math.MaxFloat32))
	fmt.Printf("Float64 min value is: 2pow(%.5g) | max: 2pow(%.5g)\n", math.Log2(math.SmallestNonzeroFloat64), math.Log2(math.MaxFloat64))
}

func Float32ErrorAccumulation() {
	var f float32 = 16777216 // 2pow32 167,772,(16) the part in brackets is unreachable
	var logF = math.Log2(float64(f))
	var fincr = f + 1 // Inspite of incrementing f, it's the same value, because the exponential part of both numbers
	// serves only 6 digits after point
	var logfincr = math.Log2(float64(fincr))

	fmt.Printf("f == fincr: %t | pows of two: %g & %g", f == fincr, logF, logfincr) // true
	// The verb "t" is used for printing bool values
}

func ScientificFloatNotation() {
	const Avogadro = .602214129e24 // Mantissa can be omitted
	const Planck = 0.62606957e-35
}

func FloatVerbsForPrinting() {
	for x := 0; x < 8; x++ {
		fmt.Printf("x = %d, E^x = %8.3g\n", x, math.Exp(float64(x))) // 8 makes the right bound of printing numbers will
		// bind to
	}
}

func InfinitiesAndResultsOfDubiousOperations() {
	var z float64
	fmt.Println(z, -z, 1/z, -1/z, z/z) // 0, -0, +inf, -inf, Nan

	fmt.Println(math.IsNaN(z/z), math.IsInf(1/z, 0), math.IsInf(1/z, -1), math.IsInf(-1/z, 1))
	// The entity is NaN if the mantissa isn't empty and exponential part is filled with 1
	// The entity is (+/-) infinity if the mantissa is empty and exponential part is filled with 1
}

func ComparisonWithNan() {
	var f float64
	var nan float64 = math.NaN()
	// Comparison with NaN always yields false
	fmt.Println(nan == nan, nan > nan, nan < nan) // false, false, false

	// With infinities the same trick works
	fmt.Println(f/f == f/f, -f/f == f/-f) // false, false
}

func NanHandling() {
	if value, ok := getNaN(100); ok {
		fmt.Printf("The value is: %.5g", value)
	} else {
		fmt.Println("An error occured. The value is zero.")
	}

}

func getNaN(value float64) (res float64, ok bool) {
	var f float64
	var result float64 = 1 / f
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, false
	} else {
		return res, true
	}
}
