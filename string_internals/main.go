package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func main() {
	fmt.Println()
	someFacts()

	fmt.Println()
	detailedFacts()

	fmt.Println()
	forRangeRunesLooping()

	fmt.Println()
	forRangeBytesLooping()

	fmt.Println()
	simpleForIterating()

}

func someFacts() {
	const (
		cnstStrA = `abcdf`
		cnstStrB = "abcdf"
	)

	var (
		a = "a"
		b = `b`
	)

	ab := a + b
	a += b

	fmt.Println(ab == a)

	strA := "abcdef"
	strB := "abc"

	fmt.Println(strB < strA)
}

func detailedFacts() {
	var (
		helloWorld = "hello world!"
		hello      = helloWorld[:5]
	)

	fmt.Println(hello[0])
	fmt.Printf("%T", hello[0])
	// hello[0] = 'q' // cannot assign to hello[0] (neither addressable nor a map index expression)compilerUnassignableOperand
	// pFirstElem = &hello[0] // invalid operation: cannot take address of hello[0] (value of type byte)compilerUnaddressableOperand
}

func constSubtleties() {
	const (
		str = `abc`

		// subStrA byte = str[5] // invalid argument: index 5 out of bounds [0:3]compilerInvalidIndex
		// subStrB byte = str[10:15] // invalid argument: index 10 out of bounds [0:4]compilerInvalidIndex

		// const expr
		lengthA = len(str)
		// not const expr
		// lengthB = len(str[:]) //
	)
}

func RuneToBytes(runes []rune) []byte {
	var (
		n int
	)

	for _, r := range runes {
		n += utf8.RuneLen(r)
	}

	n, resultingBytes := 0, make([]byte, n)

	for _, r := range runes {
		n += utf8.EncodeRune(resultingBytes[n:], r)
	}

	return resultingBytes
}

func runesBytesConversions() {
	s := `abcdefg`

	bs := []byte(s)
	s = string(bs)

	rs := []rune(s)
	s = string(rs)

	rs = bytes.Runes(bs)
	bs = RuneToBytes(rs)
}

func compilerOptimizations() {
	var (
		forRangeLooping = func() {
			var (
				str = `abcdefg`
			)

			// Here, the []byte conversion will
			// not copy the underlying bytes of str
			for i, v := range []byte(str) {
				fmt.Println(i, ":", v)
			}
		}

		mapKeyUsage = func() {
			var (
				strBytes = []byte{'k', 'e', 'y'}
				m        = make(map[string]struct{}, 10)
			)

			// string(str) makes the deep copy of str
			m[string(strBytes)] = struct{}{}

			// The optimization won't make the deep copy of underlying bytes
			fmt.Println(string(strBytes))
		}

		s string
		x = []byte{1023: 'x'}
		y = []byte{1023: 'y'}

		comparing = func() {
			// The two deep copies while comparing won't be made
			// Instead of copying the bytes will be compared one by one.
			if string(x) != string(y) {
				// The underlying bytes will be copied deeply during concatenation
				s = (" " + string(x) + string(y))[1:]
			}

			fmt.Println(s)
		}

		anotherComparing = func() {
			// During comparing the deep copies won't be made. It's not needed
			if string(x) != string(y) {
				// There are the two copies will be made, because there's no non-blank string literal
				s = string(x) + string(y)
			}
		}
	)

	forRangeLooping()

	mapKeyUsage()

	comparing()

	anotherComparing()
}

func forRangeRunesLooping() {
	var (
		str = "éक्षिaπ囧"
	)

	for i, r := range str {
		fmt.Printf("%-4d hex: 0x%x decimal: %d\n", i, r, r)
	}
	fmt.Println(len(str))
}

func forRangeBytesLooping() {
	var (
		str = "éक्षिaπ囧"
	)

	for i, b := range []byte(str) {
		fmt.Printf("Byte №%d | Byte 0x%x | Char %c\n", i, b, b)
	}
	fmt.Println(len(str))

}

func simpleForIterating() {
	var (
		str = "éक्षिaπ囧"
	)

	for i := 0; i < len(str); i++ {
		fmt.Printf("Byte №%d | Byte: 0x%x | Symbol: %c\n", i, str[i], str[i])
	}
	fmt.Println(len(str))
}

func sugaredAppendCopy() {
	hello := []byte("hello")
	world := "world"

	// The usual way is:
	// helloWorld := append(hello, []byte(world)...)

	// Sugared
	helloWorld := append(hello, world...)
	fmt.Println(string(helloWorld))

	helloWorldB := make([]byte, len(hello)+len(world))
	copy(helloWorldB, hello)

	// The usual way
	// copy(helloWorldB[len(hello):], []byte(world))

	// Sugared
	copy(helloWorldB[len(hello):], world)

	fmt.Println(helloWorld)
	fmt.Println(helloWorldB)
}

