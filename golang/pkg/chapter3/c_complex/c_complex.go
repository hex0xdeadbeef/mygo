package c_complex

import (
	"fmt"
	"math/cmplx"
)

func Declaration_And_Printing() {
	var z complex128
	var x complex128 = complex(1, 2) // 1 + 2i
	var y complex128 = complex(3, 4) // 3 + 4
	fmt.Println(z)
	fmt.Println(x * y)       // (1, 2) * (3, 4) = ((1*3 - 2*4), (2*3 + 1*4)) = -5 + 10i
	fmt.Println(real(x * y)) // -5
	fmt.Println(imag(x * y)) // 10

	pi := 3.141592i //It's 0 + 3.141592i
	two := 2i       // It's 0 + 2i
	fmt.Println(pi, two)

	fmt.Println(1i * 1i) // It's -1 + 0i
	// math/cmplx using
	fmt.Println(cmplx.Sqrt(-1)) // It's 0 + 1i

	// Arithmetic declaration
	q := 1 + 2i
	w := 3 + 4i
	fmt.Println(q, w)

}

func Comparison() {
	x := 1 + 2i
	y := 3 + 4i
	z := 1 + 2i

	fmt.Println(x == y)
	fmt.Println(x != y)
	fmt.Println(x == z)
}
