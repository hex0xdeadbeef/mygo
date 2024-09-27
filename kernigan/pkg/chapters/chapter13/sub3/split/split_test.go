package split

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {

	const (
		errTemplate = "Split(%q by %q) = %v, want %v"
	)

	var tests = []struct {
		input     string
		separator string
		want      []string
	}{
		{"a:b:c", ":", []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		got := strings.Split(test.input, test.separator)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf(errTemplate, test.input, test.separator, got, test.want)
		}
	}

}
