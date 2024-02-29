package assertions

import (
	"fmt"
	"strings"
	"testing"
)

// A poor assertion function
func assertEqual(x, y int) {
	if x != y {
		panic(fmt.Sprintf("%d != %d", x, y))
	}
}

var (
	tests = []struct {
		input     string
		separator string
		want      int
	}{
		{"a:b:c", ":", 3},
		{"1,2,3,4", ",", 4},
	}
)

func TestSplitWithAssertion(t *testing.T) {
	for _, test := range tests {
		assertEqual(len(strings.Split(test.input, test.separator)), test.want)
	}
}

func TestSplit(t *testing.T) {
	const (
		template = "Split(%q %q) returned %d words, want %d"
	)

	for _, test := range tests {
		if got := len(strings.Split(test.input, test.separator)); got != test.want {
			t.Errorf(template, test.input, test.separator, got, test.want)
		}
	}
}
