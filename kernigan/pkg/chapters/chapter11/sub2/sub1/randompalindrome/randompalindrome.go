package randompalindrome

import (
	"math/rand"
	"unicode"
)

const (
	maxLength         = 25
	unicodeSet        = 0x1000
	indexCompensation = 1
)

// randomPalindrome takes rng argument in order to have possibility to reproduce testcase, generate palindrome string, and returns it.
func randomPalindrome(rng *rand.Rand) string {
	var (
		length      = rng.Intn(maxLength)
		symbols     = make([]rune, length)
		symbolToAdd rune
	)

	for i := 0; i < (length+1)/2; i++ {
		symbolToAdd = rune(rng.Intn(unicodeSet))
		symbols[i] = symbolToAdd
		symbols[length-indexCompensation-i] = symbolToAdd
	}

	return string(symbols)
}

// randomNonPalindrome takes rng argument in order to have possibility to reproduce testcase, generate non palindrome string, and returns it.
func randomNonPalindrome(rng *rand.Rand) string {
	const (
		minLength         = 2
		unicodeshiftRight = 1
	)

	var (
		length      = minLength + rng.Intn(maxLength)
		symbols     = make([]rune, length)
		symbolToAdd rune
	)

	for i := 0; i < length; i++ {
		symbolToAdd = rune(rng.Intn(unicodeSet))
		symbols[i] = symbolToAdd
	}

	for _, symbol := range symbols {
		if unicode.IsSpace(symbol) && unicode.IsPunct(symbol) {
			length--
		}
	}

	if length <= 0 {
		symbols[0] = 'A'
	}

	return string(symbols)
}
