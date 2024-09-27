package word

import (
	"testing"
)

var (
	tests = []struct {
		input    string
		expected bool
	}{
		{"aba", true},
		{"aaaa", true},
		{"racecar", true},
		{"level", true},
		{"madam", true},
		{"rotator", true},
		{"deed", true},
		{"civic", true},
		{"radar", true},
		{"able was I ere I saw elba", true},
		{"12321", true},
		{"Able was I ere I saw Elba", true},
		{"A man, a plan, a canal: Panama", true},
		{"abb", false},
		{"abcd", false},
		{"hello", false},
		{"world", false},
		{"golang", false},
		{"1234", false},
		{"abcdefg", false},
		{"hijklmn", false},
		{"opqrstu", false},
		{"vwxyz", false},
	}
)

func TestIsPalindrome(t *testing.T) {
	const (
		errMessageTempl = "!IsPalindrome(\"%s\") = false"
	)

	for _, test := range tests {
		if got := IsPalindrome(test.input); got != test.expected {
			t.Errorf(errMessageTempl, test.input)
		}
	}
}
