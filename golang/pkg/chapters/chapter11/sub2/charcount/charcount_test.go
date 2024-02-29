package charcount

import "testing"

var tests = []struct {
	input    string
	expected string
}{
	{"abcd", "a1b1c1d1"},
	{"a777mp", "a1m1p173"},
	{"a777mp!!!", "a1m1p173!3"},
	{"", ""},
}

func TestCharCount(t *testing.T) {
	const (
		templ = "CharCount(%s) = %s, want %s"
	)
	for _, test := range tests {
		preGot := CharCount(test.input)
		if got := preGot.String(); got != test.expected {
			t.Errorf(templ, test.input, got, test.expected)
		}
	}
}
