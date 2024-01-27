package b_slices

import (
	"bytes"
	"reflect"
	"testing"
)

type Test struct {
	input    []string
	expected []string
}

type EliminateAdjacentDuplicatesTest struct {
	input    []string
	expected []string
}

var eliminateAdjacentDuplicatesTests = []EliminateAdjacentDuplicatesTest{
	{[]string{}, []string{}},
	{[]string{"apple", "banana", "banana", "orange", "orange", "kiwi"}, []string{"apple", "banana", "orange", "kiwi"}},
	{[]string{"apple", "apple", "apple", "apple"}, []string{"apple"}},
	{[]string{"cat", "dog", "fish", "fish", "cat", "dog", "bird"}, []string{"cat", "dog", "fish", "cat", "dog", "bird"}},
	{[]string{"apple", "orange", "kiwi", "kiwi", "kiwi", "banana"}, []string{"apple", "orange", "kiwi", "banana"}},
	{[]string{"one", "two", "three", "four", "five", "five"}, []string{"one", "two", "three", "four", "five"}},
	{[]string{"hello", "hello", "world", "world", "world", "world", "hello", "hello"}, []string{"hello", "world", "hello"}},
}

func TestEliminateAdjacentDuplicates(t *testing.T) {
	for _, test := range eliminateAdjacentDuplicatesTests {
		result := EliminateAdjacentDuplicates(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("\n| For: %v | expected: %v | got: %v |\n", test.input, test.expected, result)
		}
	}
}

type SquashSpacesTest struct {
	input    []byte
	expected []byte
}

var squashSpacesTests = []SquashSpacesTest{
	{[]byte("   word   "), []byte(" word ")},
	{[]byte("apple  orange"), []byte("apple orange")},
	{[]byte("   leading spaces"), []byte(" leading spaces")},
	{[]byte("trailing spaces   "), []byte("trailing spaces ")},
	{[]byte("   multiple   spaces   in   between   "), []byte(" multiple spaces in between ")},
	{[]byte("   123   456   789   "), []byte(" 123 456 789 ")},
	{[]byte("special   characters   !@#$%^&*()"), []byte("special characters !@#$%^&*()")},
	{[]byte("   leading   and   trailing   spaces   "), []byte(" leading and trailing spaces ")},
	{[]byte("   a   b   c   d   e   "), []byte(" a b c d e ")},
	{[]byte("   a   b   c   "), []byte(" a b c ")},
	{[]byte("  word  word  word  "), []byte(" word word word ")},
	{[]byte("apple  orange  banana  kiwi  grape  watermelon  "), []byte("apple orange banana kiwi grape watermelon ")},
	{[]byte("   cat   dog   fish   bird   "), []byte(" cat dog fish bird ")},
	{[]byte("red   green   blue   yellow   purple   "), []byte("red green blue yellow purple ")},
	{[]byte("  alpha  beta  gamma  delta  epsilon  "), []byte(" alpha beta gamma delta epsilon ")},
	{[]byte("   1   2   3   4   5   6   7   8   9   "), []byte(" 1 2 3 4 5 6 7 8 9 ")},
	{[]byte("   special   characters   !@#$%^&*()   "), []byte(" special characters !@#$%^&*() ")},
	{[]byte("   leading   and   trailing   spaces   "), []byte(" leading and trailing spaces ")},
	{[]byte("   a   b   c   d   e   f   g   h   i   j   k   l   m   n   o   p   q   r   s   t   u   v   w   x   y   z   "), []byte(" a b c d e f g h i j k l m n o p q r s t u v w x y z ")},
	{[]byte("   1   2   3   4   5   6   7   8   9   0   "), []byte(" 1 2 3 4 5 6 7 8 9 0 ")},
}

func TestSquashSpaces(t *testing.T) {
	for _, test := range squashSpacesTests {
		result := SquashSpaces(test.input)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("\nFor:%s\nExpected:%s\nGot:%s\n", string(test.input), string(test.expected), string(result))
		}
	}
}

func TestSquashSpacesSimple(t *testing.T) {
	for _, test := range squashSpacesTests {
		result := SquashSpacesSimple(test.input)
		if !bytes.Equal(result, test.expected) {
			t.Errorf("\nFor:%s\nExpected:%s\nGot:%s\n", string(test.input), string(test.expected), string(result))
		}
	}
}
