package e_strings

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

/* ALL SOURCE FILES ARE CONCIDERED AS UTF-8 CODED TEXT */

func declarationAndInitialization() {
	s := "hello, world"
	fmt.Println(len(s))

	/* When we reference to a string by an index, we require two or more bytes */
	// We can access the elements of string referencing by index and getting a rune value
	fmt.Println(s[0], s[7]) // "104 119" - unicode code points (the initial byte of symbol)
	// There will be a panic error because of referencing to the nonexistent index
	fmt.Println(s[10])
}

func gettingSubstring() {
	str := "hello world"

	/*After slicing the parts of the original strings refer to the original memory points*/
	// The right boundary isn't included
	fmt.Println(str[0:5]) // "hello"
	fmt.Println(str[5:])  // "world"
	fmt.Println(str[:])   // "hello world"

	/*Here might be a panic error as well*/
}

func concatenation() {
	str := "hello world"
	fmt.Println("goodbye, " + str[5:]) // "goodbye, world"
}

func StringComparison() {
	str1 := "abcd"
	str2 := "abe"

	fmt.Println((str1 == str2))
	// This operation runs from left to right and checks the lexographical order
	fmt.Println((str1 > str2)) // false because a == a, b == b, 'e' is located earlier than 'd'
}

func StringsImmutability() {
	s := "left foot"
	/*The string below will be marked as wrong because of stings immutability*/
	// s[0] = 's'
	t := s
	s += ", right foot" // This is just reassignment the new value
	fmt.Println(t)
	fmt.Println(s)
}

func StringLiterals() {
	fmt.Println("\a")
	fmt.Println("abcde\bc")
}

func EscapeSequences() {
	char1 := '\xaa'
	char2 := '\377'
	fmt.Printf("The characters are: %c %c\n", char1, char2) // ª ÿ
	fmt.Printf("The runes are: %d %d\nx", char1, char2)     // 170 255
}

func RawStringLiteral() {
	fmt.Println(`Это сырая строка здесь мы можем делать спуск вниз без \n
	и писать управляющие символы, которые не будут восприняты как урпавляющие.`)
}

func DecodingUTF8() {
	fmt.Println("\xe4\xb8\x96\xe7\x95\x8c") // This is UTF-8 notation
	/*Additionally, \xe4\xb8\x96 is not a legal rune literal */
	fmt.Println("\u4e16\u754c")         // This is Unicode notation (We should convert from hex to binary)
	fmt.Println("\U00004e16\U0000754c") // This is extended Unicode notation (We should convert from hex to binary)

	// We can hold unicode escapes in the rune type directly
	var a rune = '世'
	var b rune = 19990
	fmt.Printf("%c\n", a) // There will be printed: 19990. It's a point of Unicode
	fmt.Printf("%c", b)   // There will be printed 世, the character correspoding to the Unicode point
}

func LenAndRunes() {
	s := "Hello, 世界"
	fmt.Printf("The amount of bytes is: %d\n", len(s))                             // func len() returns the amount of bytes required to store the string
	fmt.Printf("The amount of unicode points is: %d\n", utf8.RuneCountInString(s)) // utf.RuneCountInString(str) returns the amount of unique unicode points

	// This representation isn't valid because of reading each byte partially
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c ", s[i])
	}

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:]) // Function checks the amount of bytes of each symbol, returns the rune
		// that can be represented as (a/an) character/unicode point
		fmt.Printf("%d\t%c\t%d\t%x\t%d\n", i, r, r, r, size)
		i += size
	}
	fmt.Println()

	// Bytes reading occurs implicitly
	for i, unicodePoint := range s {
		fmt.Printf("%d\t", i)
		fmt.Printf("%[1]c\t%[1]d\t%[1]x\t%[1]o\n", unicodePoint)

	}
	fmt.Println()

	// The convenient method to count the number of runes
	var runesNumber int
	for _, _ = range s {
		runesNumber++
	}

	// Another convinient method to count the number of runes
	runesNumber = 0
	for range s {
		runesNumber++
	}

	// MEGAXXL convinient method to count the number of runes
	fmt.Printf("The number of runes is: %d\n", utf8.RuneCountInString(s))

	// If UTF-8 decoder meets an unexpected input byte it prints replacement character '\uFFFD' = �
	fmt.Println("Replacement symbol: '\uFFFD'")

	/* []rune conversion applied to a UTF-8 encoded string returns the sequence of unicode points */
	s = "もちろん、喜んで"
	// "% x" makes the space between each pair of operands
	fmt.Printf("% x\n", s)                  // This statement prints the hexadecimal UTF-8 codes of the string
	stringToRunesArray := []rune(s)         // []runes(str string) returns Unicode code points
	fmt.Printf("%x \n", stringToRunesArray) // This statements prints the hexadecimal Unicode code points

	/* The opposite conversion from []rune(Unicode code points) to string will return UTF-8 codes string */
	fmt.Println(string(stringToRunesArray))

	/* During conversion from int value to string value the first one is perceived as a rune (unicode code point) */
	fmt.Printf("%c\n", 65)     // It prints "A"
	fmt.Printf("%c\n", 0x3eac) // It prints 㺬

	/* Very big codes will produce returning the replacement symbol (�) */
	fmt.Printf("%c", 1234567) // It prints �

}

func BasenameSimple(str string) {
	runes := []rune(str)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == '/' || runes[i] == '\\' {
			runes = runes[i+1:]
			break
		}
	}

	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == '.' {
			runes = runes[:i]
			break
		}
	}

	fmt.Println(string(runes))
}

func BasenameLastIndex(str string) {
	if slashInd, dotInd := strings.LastIndex(str, "/"), strings.LastIndex(str, "."); slashInd != -1 || dotInd != 1 {
		if dotInd != -1 {
			str = str[:dotInd]
		}
		if slashInd != -1 {
			str = str[slashInd+1:]
		}
	}
	fmt.Println(str)
}

func IntegerComma(number string) string {
	curLen := len(number)
	if curLen <= 3 {
		return number
	}
	return IntegerComma(number[:curLen-3]) + "," + number[curLen-3:]
}

func BytesSliceMutability() {
	/* There are three independent blocks of memory. Two of them serves the strings immutability */
	// A string value has an array of bytes, that can't be changed
	str1 := "comma"
	// On the other hand we can create a byte slice of this string we capable to change
	stringBytesSlice := []byte(str1) // []byte(str) statement copies all the bytes from a string to a named variable

	fmt.Printf("BEFORE: The bytes array: %d\n", stringBytesSlice)
	stringBytesSlice[0] = 65 // We can change the values from 0 to 128 (byte)
	fmt.Printf("AFTER: The bytes array: %d\n", stringBytesSlice)

	// We can also do the opposite operation
	str2 := string(stringBytesSlice)
	fmt.Printf("The string from the bytes array: %s\n", str2)

	/*
		IT'S BETTER TO USE THE "bytes" PACKAGE EMBEDED FUNCTIONS(Contains, Count, HasPrefix, Index, Join and so on)
		INSTEAD THEIR "strings" PACKAGE COUNTERPARTS.

		IT PREVENTS SYSTEM FROM EXTRA CONVERSIONS AND UNNECESSARY MEMORY ALLOCATION
	*/
}

func IntsToString(values []int) string {
	var buffer bytes.Buffer
	// buffer.WriteByte('[') // It converts the rune to byte implicitly
	buffer.WriteRune('[') // It's preferable
	for i, value := range values {
		// Initial comma won't be added
		if i > 0 {
			buffer.WriteString(", ") // It converts the string to byte implicitly
		}
		fmt.Fprintf(&buffer, "%d", value) // It converts the string to byte implicitly
	}
	// buffer.WriteByte(']') // It converts the rune to byte implicitly
	buffer.WriteRune(']') // It's preferable

	return buffer.String()
}
