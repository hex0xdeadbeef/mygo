package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

/*
	STRING INTERNALS

	STRING
1. String is the structure that cannot be changed. It consists of:
	1) A pointer to a sequence of bytes that cannot be changed
	2) An amount of bytes in this sequence

	MISUNDERSTANDING OF RUNES CONSEPT
1. Charset is the set of symbols, for example: the Unicode charset includes 2^21 symbols
2. Encoding is retranslation of a symbols list into a binary code. For example: UTF-8 - the standart of encoding defining the way of encoding all the symnols of Unicode to the varying
amount of bytes (from 1 to 4 bytes).
3. In Unicode we use the concept of a code point to reference to an element. Go's rune - a code pointer of Unicode
4. UTF-8 encodes symbols of an amount from 1 to 4 bytes (up to 32 bits). That's why Go's rune - the alias of int32.
	type rune int32
5. Go doesn't always use the UTF-8 encoding. When we work with a variable, that wasn't initialized with string literal (for example: reading from file system)
6. len(string) returns an amount of bytes in a string, not the amount of runes.
	1) A symbol is not always represented as a single byte. For example '㹤' is 3 bytes in UTF-8
7. We can create a string from its corresponding bytes and vice versa.
8. The source code in Go uses UTF-8 encoding. All the string literals - strings encoded with UTF-8. But if the string has been gotten abroad, there's no guarantee that it's based on UTF-8.

	WRONG ITERATION OVER A STRING
1. An amount of rune in a string can be fetched with utf8.RuneCountInString(string arg)
2. When we convert a string into a slice of runes to range over it, we gets overhead (O(n) - the amount of bytes needed to arrange all the runes in the new slice)
3. Ranging over a string with the for-range loop the first variable is responsible for reflecting a start of the current rune in this string.
4. Working with the string when each rune is represented as a single-byte we can optimize the access to a specific rune:
	We can directly reference to the rune by an index

	STRING SLICING
1. TrimRight and TrimSuffix
	1) TrimRight removes all the trailing elems contained in charset given to the func starting from the right end of given string.
	2) On the other hand the function TrimSuffix removes the first suffix given as the argument. When it's done the processed string is returned. It stops after the first deletion of the tar-
	get suffix.
	The same logic is applied to the TrimPrefix(...) and TrimLeft(...)
	3) Trim is the combination of TrimLeft(...) and TrimRight(...)

	CONCATENATION OPTIMIZATIONS
	The way with premature allocation of memory of the strings beats the way with "+=" 99%, strings.Builder without premature allocation 78%.
	string.Build is the best solution of concatenation when an amount of strings given is bigger than 5.

1. The best way to concatenate two strings is to use the strings.Builder because it manages the underlying builder growing/shrinking it when it's needed.
	1) Since the strings.Builder implements the interface io.StringWriter, there's need of returning an error, but there won't be non-nil error. It's idiomatic to skip this error.
2. strings.Builder includes the byte slice. Every call of WriteString(...) results in append(src, args) incovation applied to this slice.
	1) strings.Builder isn't concurrent safe structure because the multiple appends will cause race conditions.
	2) If the resulting length of all the strings is known before, we should use the strings.Builder.Grow(n int) method because we don't want strings.Builder to make multiple reallocations.
	The (n int) parameter is the amount of bytes.
3. Another good way of concatenation is to take all the strings given, sum theirs lengths, allocate the byte slice of this resulting length and put the bytes of the strings into it.
4. The last way of concatenation is using fmt.Sprintf(...) function


*/

func main() {
	// LenUsage()

	// Utf8CountInString()

	// WrongIteration()

	// ConversionFromStringToRunes()

	// RuneAccessOptimization()

	// TrimRightAndTrimSuffix()
}

func LenUsage() {
	var (
		s                    = "hello"
		chineseSymbol        = "㹤"
		bytesOfChineseSymbol = []byte{0xe3, 0xb9, 0xa4}
	)

	fmt.Println(len(s))             // 5
	fmt.Println(len(chineseSymbol)) // 3

	chineseSymbol = string(bytesOfChineseSymbol)
	fmt.Println(chineseSymbol)
}

func WrongIteration() {
	var (
		a = func() {
			var (
				s = "hêllo"
			)

			// In this case we iterate only over the first runes' positions, not the runes itself.
			// This code has the curious moment:
			// 1) The second rune printed is 'Ã' not the 'ê'.
			// It happens because referencing to the symbol by the exr s[i] we capture only the rune of this byte.
			for i := range s {
				fmt.Printf("Pos %d: %c\n", i, s[i])
			}
			fmt.Println()

			fmt.Printf("len=%d\n", len(s))
		}

		b = func() {
			var (
				s = "hêllo"
			)

			// This will be treated as ranging over all the runes in the string. The i is the start byte of a current Unicode code pointer in the string
			for i, r := range s {
				fmt.Printf("Pos %d: %c\n", i, r)
			}

			fmt.Printf("len=%d\n", len(s))
		}

		c = func() {
			var (
				s = "hêllo"
			)

			// In this case we just retranslate the string as a rune slice
			for i, r := range []rune(s) {
				// The position here means the position in the slice of runes gotten before the execution of the for-range loop
				fmt.Printf("Pos %d: %c\n", i, r)
			}

			fmt.Printf("len=%d\n", len(s))
		}
	)
	a()

	fmt.Println()
	b()

	fmt.Println()
	c()
}

func Utf8CountInString() {
	var (
		s = "hêllo"
	)

	fmt.Println(utf8.RuneCountInString(s))
}

func ConversionFromStringToRunes() {
	var (
		s = "hêllo"
	)

	fmt.Printf("%c\n", []rune(s)[4])
}

func RuneAccessOptimization() {
	var (
		s = "hello"
	)

	fmt.Printf("%c\n", rune(s[1]))
}

func TrimRightAndTrimSuffix() {
	fmt.Println(strings.TrimRight("123ooxo", "xo"))  // 123
	fmt.Println(strings.TrimSuffix("123ooxo", "xo")) // 123oo

	fmt.Println()

	fmt.Println(strings.TrimLeft("xoxoxo123xoxoxo", "xo"))   // 123xoxoxo
	fmt.Println(strings.TrimPrefix("xoxoxo123xoxoxo", "xo")) // xoxo123xoxoxo

	fmt.Println()

	fmt.Println(strings.Trim("xoxoxo123xoxoxo", "ox")) // 123
}

// SillyConcatenation returns a concatenation of strs given
// Since the string isn't changable entity, at each iteration allocation of a new byte sequence happens and it results in low performance and memory overhead
func SillyConcatenation(strs ...string) string {
	var (
		res string
	)

	for i := 0; i < len(strs); i++ {
		res += strs[i]
	}

	return res
}

// OptimizedSillyConcatenation returns a concatenation of given strings
// This function based on the counting all the lengths of a strings given and premature allocating a byte slice of the resulting length
func OptimizedSillyConcatenation(strs ...string) string {
	var (
		resLength int
	)

	for i := 0; i < len(strs); i++ {
		resLength += len(strs[i])
	}

	var (
		res = make([]byte, resLength)
		idx int
	)

	for i := 0; i < len(strs); i++ {
		for j := 0; i < len(strs[j]); j++ {
			res[idx] = strs[i][j]
			idx++
		}
	}

	return string(res)
}

// FastConcatenation returns a concatenation of strs and a slice of errors encountered while concatenating the strs
// Skipping the error of WriteString method is the idiomatic way
func FastConcatenation(strs ...string) string {
	var (
		b *strings.Builder
	)

	for i := 0; i < len(strs); i++ {
		_, _ = b.WriteString(strs[i])
	}

	return b.String()
}

// FastConcatenationOptimized returns a concatenation of strs and a slice of errors encountered while concatenating the strs
// It's optimized concatenation of strs because of premature calculating of resulting bytes count to be used in b.Grow(n int method), when n is the amount of bytes needed to store the
// next string
// Skipping the error of WriteString method is the idiomatic way
func FastConcatenationOptimized(strs ...string) string {
	var (
		b           *strings.Builder
		resBytesCnt int
	)

	for i := 0; i < len(strs); i++ {
		resBytesCnt += len(strs[i])
	}

	b.Grow(resBytesCnt)

	for i := 0; i < len(strs); i++ {
		_, _ = b.WriteString(strs[i])
	}

	return b.String()
}
