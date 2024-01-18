package f_constants

import (
	"fmt"
	"math"
	"time"
)

/*
Many computations on constants can be completely evaluated at compile time, reducing the work necessary at runtime
and enabling other compiler optimisations
*/

/* Constants can be uncommited to a particular type. This mechanism lets use arithmetic operations on some values with
much higher precision. Precision can reach up to 256 bits
These types are:
 untyped boolean,
  untyped integer,
   untyped rune,
    untyped floating-poing,
	 untyped complex,
	  untyped string
*/

const (
	e       = 2.718281 // untyped float
	pi      = 3.141592 // untyped float
	IPv4LEn = 4        // For example it'll be used during an array declaring/initializing as a capacity value
)

/* Constants may specify a type as well as a value, but in the absence of an explicit type the type is inferred from
the expression on the right hand */

const (
	noDelay time.Duration = 0               // noDelay is a constant that represents a value 0 of the time.Duration type
	timeout               = 5 * time.Minute // timeout is a constant that represents a value 5 minutes of the time.Minute type
)

func PrintingConstantsValues() {
	fmt.Printf("%T %[1]v\n", noDelay) // There is %T verb allows to print the type of a value
	fmt.Printf("%T %[1]v\n", timeout) // There is %T verb allows to print the type of a value
	fmt.Printf("%T %[1]v\n", time.Minute)
}

/* When a sequence of constants are declared as a group, the right-hand side expression may be ommited*/
const (
	a = 1 // This value descends
	b     // It'll be assigned the value of the upper constant
	c = 2 // This value descends
	d     // It'll be assigned the value of the upper constant
)

func PrintingImplyingConstants() {
	fmt.Printf("a: %d, b: %d, c: %d, d: %d", a, b, c, d)
}

/* Constant declarations may be realized with iota tool. the value of iota begins at zero and increments by one
for each item in the sequence.

iota = 0 ... n
*/

type Weekday int // It implies that the underlying type is int

const (
	Sunday Weekday = iota // 0 <= the iota start value is 0. It will increment itself and be assigned to each
	// underlying variable
	Monday    // 1
	Tuesday   // 2
	Wednesday // 3
	Thursday  // 4
	Friday    // 5
	Saturday  // 6
)

type Flags uint

const (
	FlagUp           = 1 << iota // 000...00001
	FlagBroadcast                // 000...00010
	FlagLoopBack                 // 000...00100
	FlagPointToPoint             // 000...01000
	FlagMulticast                // 000...10000
)

// IsUp(v Flags) returns true if the result of multiplying v by FlagUp has 1 junior bit. It's the check whether v flag
// is FlagUp
func IsUp(v Flags) bool {
	return v&FlagUp == FlagUp
}

// TurnDown(v *Flags) clears the last bit of v value (after operation v isn't FlagUp)
func TurnDown(v *Flags) {
	*v &^= FlagUp
}

// SetBroadcast(v *Flags) sets the third junior bit with 1. It makes v the FlagBroadcast flag
func SetBroadcast(v *Flags) {
	*v |= FlagBroadcast
}

// IsCast(v Flags) looks the values at the 2 and 5 positions and if they're both not zeros returns true
func IsCast(v Flags) bool {
	return v&(FlagBroadcast|FlagMulticast) != 0
}

/* This program computes the powers of 1024 (2^10) */
const (
	_   = 1 << (10 * iota) // This part will be skipped, but iota will continue incrementing itself
	KiB                    // 2^10
	MiB                    // 2^20
	GiB                    // 2^30
	TiB                    // 2^40. It's more than 2^31-1 => There will be the overflwow of int32 type
	PiB                    // 2^50
	EiB                    // 2^60
	ZiB                    // 2^70 /* This constant is too big to be used directly, but it can be used in expressions */
	YiB                    // 2^80 /* This constant is too big to be used directly, but it can be used in expressions */
)

const (
	KB = 10e3
	MB = KB * 10e3
	GB = MB * 10e3
	TB = GB * 10e3
	PB = TB * 10e3
	ET = PB * 10e3
	ZB = ET * 10e3
	YB = ZB * 10e3
)

// The result is: 1024. The type of the result value will be untyped int :)
func UsingBigConstantsInExpressions() {
	fmt.Printf("Some huge untyped constants are computed in expressions (YiB/ZiB).\nThe result of this statement: %[1]v, %[1]T",
		YiB/ZiB)
}

/* Another example of untyped values. The right side statements are assignment a untyped corresponding value */
var (
	x float32    = math.Pi
	y float64    = math.E
	z complex128 = math.Phi
)

/*
When we bind an untyped value to a typed value the following expressions will require an corresponding explicit type
conversion
*/
const Pi64 float64 = math.Pi

var (
	q float32    = float32(Pi64)
	w complex128 = complex128(Pi64)
	r float64    = Pi64 // This statement doesn't require the explicit type conversion
)

/* Literals flavors define what type will be used in an expression. For example: */
const (
	zeroInt     = 0   // It's untyped int zero
	zeroFloat   = 0.0 // It's untyped float zero
	zeroComplex = 0i  // It's untyped complex zero
)

/* The choice of literal may affect the result of a constant devision expression */
func DivisionWithDifferentLiterals() {
	var value float64 = 212
	fmt.Printf("%T\n", (value-32)*5/9)     // 100 It'll be a float64 value
	fmt.Printf("%T\n", 5/9*(value-32))     // 0 It'll be a float64 value
	fmt.Printf("%T\n", 5.0/9.0*(value-32)) // 100 It'll be a float64 value
}

/*
	When an untyped constant occurs on the right hand side of a variable, it's implicitly converted to the target type

(to the left hand variable type)
*/
func ImplicitConversionToTargetType() {
	var f float64 = 3 + 0i // It'll be conversed to float64
	fmt.Printf("%T\n", f)
	f = 2 // It'll be converted to float64
	fmt.Printf("%T\n", f)
	f = 1e123 // It'll be converted to float64
	fmt.Printf("%T\n", f)
	f = 'a' // It'll be also converted to float 64
	fmt.Printf("%T\n", f)
}

// After an explicit conversion the precision may be lost or got the higher value
func RoundingAfterConversion() {
	const deadbeef = 0xdeadbeef // this is untyped int with the high precision 2^31.798819721089007
	fmt.Printf("The original deadbeef value: %v | power of two: %g\n", deadbeef, math.Log2(deadbeef))
	fmt.Printf("The uint32 deadbeef value: %d\n", uint32(deadbeef))   // The original value will be saved
	fmt.Printf("The float32 deadbeef value: %g\n", float32(deadbeef)) /* The 7th digit will get up (float32 provides
	only 6-digit precision) */
	fmt.Printf("The float64 deadbeef value: %g\n", float64(deadbeef)) // The precision will be saved
	// const d = int32(deadbeef) compile error: constant overflows int32
	// const f = float64(1e309) compile error: constant overflows float64
	// e = uint(-1) compile error: constant underflows uint
}

// The integers will have untyped types after assignment but float/complex variables will have particular sized types
// such as float64 or complex128
func ImplicitVariableLiteralTypeDetermination() {
	i := 0 // It'll be untyped integer. Implicit: int(0)
	fmt.Printf("%[1]T %[1]d\n", i)
	r := '\000' // It'll be untyped rune. Implicit: rune('\000') or int32('\000')
	fmt.Printf("%[1]T %[1]c\n", r)
	f := 0.0 // It'll be untyped float. Implicit float64(0.0)
	fmt.Printf("%[1]T %[1]g\n", f)
	c := 0i // It'll be untyped complex. Implicit complex128(0i)
	fmt.Printf("%[1]T %[1]v\n", c)
}

// So that we can explicitly inform the compiler that a varible must have the particular type, we should specify the
// type we need variable to have explicitly
func ExplicitVariableTypeDeclaration() {
	var i1 = int8(0)
	var i2 int8 = 0
	fmt.Printf("%[1]T %[1]d\n", i1) // int8
	fmt.Printf("%[1]T %[1]d\n", i2) // int8
}

// It's important when converting an untyped constant to an interface value since they determine its dynamic type
