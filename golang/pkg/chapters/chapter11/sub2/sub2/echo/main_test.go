package main

import (
	"bytes"
	"fmt"
	"testing"
)

var (
	echoTests = []struct {
		newline bool
		sep     string
		args    []string
		want    string
	}{
		{true, "", []string{}, "\n"},
		{false, "", []string{}, ""},
		{true, "\t", []string{"one", "two", "three"}, "one\ttwo\tthree\n"},
		{true, ",", []string{"a", "b", "c"}, "a,b,c\n"},
		{false, ":", []string{"1", "2", "3"}, "1:2:3"},
	}
)

func TestEcho(t *testing.T) {
	const (
		templ = "echo(%v, %q, %q)"
	)

	var (
		description string
	)

	for _, test := range echoTests {

		description = fmt.Sprintf(templ, test.newline, test.sep, test.args)
		out = new(bytes.Buffer)

		if err := echo(test.newline, test.sep, test.args); err != nil {
			t.Errorf("%s failed: %v", description, err)

		}

		if got := out.(*bytes.Buffer).String(); got != test.want {
			t.Errorf("%s = %q, want %q", description, got, test.want)
		}
	}
}
