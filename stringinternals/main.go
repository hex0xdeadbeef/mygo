package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func main() {
	lenUsage()

	fmt.Println()

	concat()

	fmt.Println()

	comparing()

	fmt.Println()
	shareableUnderlyingByteSequence()

	fmt.Println()
	idxClobbering()

	fmt.Println()
	failToCompileIndexChanging()

	fmt.Println()
	immutabilityAndNonaddressabilityOfSingleByte()

	fmt.Println()
	incorrectStringConstant()

	fmt.Println()
	directConversionFromIntegerToString()

	fmt.Println()
	runesLen(nil)

	fmt.Println()
	forRangeRunes()

	fmt.Println()
	forRangeBytes()
}

func lenUsage() {
	hello := `Привет!`

	fmt.Println(len(hello)) // 13 bytes
}

func concat() {
	a := "abcd"
	b := `efg`

	ab := a + b

	fmt.Println(ab)

	c := ""
	c += `!`
	c += "?"
	fmt.Println(c)
}

func comparing() {
	less := `abcd`
	bigger := `abcdefg`

	fmt.Println(less == bigger)
	fmt.Println(bigger > less)
	fmt.Println(less > bigger)
}

func shareableUnderlyingByteSequence() {

	str := "abcdefg"
	// The same underlying byte sequence with str
	subStrA := str
	subStrB := str[1:4]

	fmt.Println(str)
	fmt.Println(subStrA)
	fmt.Println(subStrB)
}

func idxClobbering() {
	hello := `Привет, мир!`

	fmt.Println(hello[:3])
	fmt.Printf("%[1]T %[1]v\n", []byte("abcde")[0])
}

func failToCompileIndexChanging() {
	str := `abcde`
	fmt.Println("Before", str)

	// str[0] = `q` // cannot assign to str[0] (neither addressable nor a map index expression)compilerUnassignableOperand
	strBytes := []byte(str)
	strBytes[0] = 'q'

	str = string(strBytes)
	fmt.Println("After:", str)
}

func immutabilityAndNonaddressabilityOfSingleByte() {
	str := `abcdefg`

	// str[0] = `q` // cannot assign to str[0] (neither addressable nor a map index expression)compilerUnassignableOperand
	// zeroElemPointer := &str[0] // invalid operation: cannot take address of str[0] (value of type byte)compilerUnaddressableOperand

	fmt.Println(str[:0])
}

func incorrectStringConstant() {
	const (
		str = `abc`

		// subStrA = str[5:] // invalid argument: index 5 out of bounds [0:4]compilerInvalidIndex
		// subStrB byte = str[10] // invalid argument: index 10 out of bounds [0:3]compilerInvalidIndex

		// const expr
		strLenA = len(str)

		// not const expr
		// strLenB = len(str[:]) // len(str[:]) (value of type int) is not constantcompilerInvalidConstInit

	)

}

func directConversionFromIntegerToString() {
	num := 100

	fmt.Println(string(num))
	fmt.Println('A')
	fmt.Println(string(-1))
}

func runesLen(runes []rune) int {
	var (
		length int
	)

	for _, v := range runes {
		length += utf8.RuneLen(v)
	}

	fmt.Println(length)
	return length
}

func conversionsBetweenRunesAndBytes(runes []rune) []byte {
	var (
		l = runesLen(runes)

		n, bts = 0, make([]byte, l)
	)

	for _, r := range runes {
		n += utf8.EncodeRune(bts[n:], r)
	}

	fmt.Println(bts)

	return bts
}

func multipleConversions() {
	s := `abcdefg`

	bs := []byte(s)
	s = string(bs)

	rs := []rune(s)
	s = string(rs)

	rs = bytes.Runes(bs)
	bs = conversionsBetweenRunesAndBytes(rs)
}

func compilerOptimizations() {
	var (
		rangingOverByteSliceFromString = func() {
			var str = "abcdef"

			// Here the []byte(str) conversion will not copy the underlying bytes of str.
			for i, b := range []byte(str) {
				fmt.Println(i, " ", b)
			}
		}

		usingStringKeyFromByteSlice = func() {
			var (
				key = []byte{'k', 'e', 'y'}

				m = make(map[string]struct{}, 10)
			)

			// The string([]byte) in this conversion copies the bytes in key
			m[string(key)] = struct{}{}

			// Here, this string([]byte) conversion doesn't copy the bytes in key. The optimization will be still made, even if key is a package-level variable.
			fmt.Println(m[string(key)])
		}

		comparisonOfConversionsByteSliceToString = func() {
			var (
				x = []byte{1023: 'x'}
				y = []byte{1023: 'y'}

				s string
			)

			// None of the four comparisons will make deep copies of the underlying byte slices.
			// But during concatenation the conversions of course will be made.
			if string(x) != string(y) {
				// There will be the only one allocation for the resulting slice
				s = (" " + string(x) + string(y))[1:]
			}

			fmt.Println(s)

		}

		noOptimisation = func() {
			var (
				x = []byte{1023: 'x'}
				y = []byte{1023: 'y'}

				s string
			)

			// During the comparison the copies won't be made
			if string(x) != string(y) {
				// There will be three copies:
				// The copy of "x"
				// The copy of "y"
				// Allcoation for the concatenation
				s = string(x) + string(y)
			}

			fmt.Println(s)
		}
	)

	rangingOverByteSliceFromString()

	usingStringKeyFromByteSlice()

	comparisonOfConversionsByteSliceToString()

	noOptimisation()
}

func forRangeRunes() {
	var (
		s = "éक्षिaπ囧"
	)

	for i, v := range s {
		fmt.Printf("%-3v 0x%x %v \n", i, v, string(v))
	}

	fmt.Println(len(s))
}

func forRangeBytes() {
	var (
		s = "éक्षिaπ囧"
	)

	for i := 0; i < len(s); i++ {
		fmt.Printf("The byte at index %d is 0x%x\n", i, s[i])
	}

	fmt.Println()

	// Here's the bytes aren't copied into a new slice
	for i, b := range []byte(s) {
		fmt.Printf("The byte at index %d is 0x%x\n", i, b)
	}

	fmt.Println(len(s))
}
