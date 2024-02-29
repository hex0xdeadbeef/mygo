package word

import (
	"strings"
	"unicode"
)

// IsPalindrome reports whether a string is palindrome. Letter case is ignored, as are non-letters
func IsPalindrome(s string) bool {
	const (
		indexCompensation = 1
	)
	var (
		runes    = formatInput(s)
		boundary = len(runes)/2 + indexCompensation
	)

	if len(runes) <= 1 {
		return true
	}

	for i := 0; i < boundary; i++ {
		if runes[i] != runes[len(runes)-i-indexCompensation] {
			return false
		}
	}

	return true
}

func formatInput(str string) []rune {
	result := make([]rune, 0, len(str))
	str = strings.ToLower(str)

	for _, character := range str {
		if !unicode.IsSpace(character) && !unicode.IsPunct(character) {
			result = append(result, character)
		}
	}

	return result
}
