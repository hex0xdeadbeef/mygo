package integers

import "fmt"

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