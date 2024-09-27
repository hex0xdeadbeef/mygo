package randompalindrome_test

import (
	"golang/pkg/chapters/chapter11/sub2/word"

	"testing"
)

func BenchmarkIsPalindtome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		word.IsPalindrome("A man, a plan, a canal: Panama")
	}
}

