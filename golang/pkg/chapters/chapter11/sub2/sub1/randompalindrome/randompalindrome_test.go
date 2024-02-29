package randompalindrome

import (
	"golang/pkg/chapters/chapter11/sub2/word"
	"math/rand"
	"testing"
	"time"
	"unicode/utf8"
)

func TestRandomPalindrome(t *testing.T) {
	const (
		testsNumber = 1000
		templ       = "IsPalindrome(%s) = false"
	)

	var (
		seed = time.Now().UTC().UnixNano()
		rng  = rand.New(rand.NewSource(seed))

		testcase string
	)

	t.Logf("Random seed: %d", seed)

	for i := 0; i < testsNumber; i++ {
		testcase = randomPalindrome(rng)

		if !word.IsPalindrome(testcase) {
			t.Errorf(templ, testcase)
		}
	}
}

func TestRandomNonPalindrome(t *testing.T) {
	const (
		testsNumber = 1000
		templ       = "!IsPalindrome(%s) = false, runes count = %d"
	)

	var (
		seed = time.Now().UTC().UnixNano()
		rng  = rand.New(rand.NewSource(seed))

		testcase string
	)

	t.Logf("Random seed: %d", seed)

	for i := 0; i < testsNumber; i++ {
		testcase = randomNonPalindrome(rng)

		if word.IsPalindrome(testcase) {
			t.Errorf(templ, testcase, utf8.RuneCountInString(testcase))
		}
	}
}
