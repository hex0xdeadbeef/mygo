package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

func main() {

	// CloneUsage()

	// CompareUsage()

	// ContainsUsage()

	// ContainsAnyUsage()

	// ContainsFuncUsage()

	// ContainsRuneUsage()

	// CountUsage()

	// CutUsage()

	// CutPrefixUsage()

	// CutSuffixUsage()

	// EqualFoldUsage()

	// FieldsUsage()

	// FieldsFuncUsage()

	// HasPrefixUsage()

	// HasSuffixUsage()

	// IndexUsage()

	// IndexByteUsage()

	// IndexRuneUsage()

	// IndexAnyUsage()

	// JoinUsage()

	// LastIndexUsage()

	// LastIndexAnyUsage()

	// LastIndexByteUsage()

	// LastIndexFuncUsage()

	// MapUsage()

	// RepeatUsage()

	// ReplaceUsage()

	// ReplaceAllUsage()

	// SplitUsage()

	// SplitAfterUsage()

	// SplitAfterNUsage()

	// SplitNUsage()

	// ToLowerUsage()

	// ToLowercaseSpecialUsage()

	// ToTitleUsage()

	// ToTitleSpecialUsage()

	// ToUpperUsage()

	// ToUpperSpecialUsage()

	// ToValidUTF8Usage()

	// TrimUsage()

	// TrimFuncUsage()

	// TrimLeftFuncUsage()

	// TrimPrefixUsage()

	// TrimRightUsage()

	// TrimRightFuncUsage()

	// TrimSpaceUsage()

	// TrimSuffixUsage()

	/*
		strings.Builder()
	*/

	BuilderZeroValueUsage()
	BuilderJoinMultipleByteSlices()
	BuilderJoinStrings()
	BuilderJoinRunes()
	BuilderJoinBytes()
	BuildersComparison()

	BuilderResetUsage([]string{"abc", "def", "qwerty", "zxc"})
	BuildersComparison()
}

func CloneUsage() {
	str := `abcdefg`
	clonedStr := strings.Clone(str)
	fmt.Println(unsafe.Pointer(&str) == unsafe.Pointer(&clonedStr)) // false
	fmt.Println(str == clonedStr)                                   // true

	fmt.Println()
	emptyStr := ""
	// No allocation is made here
	clonedEmptyStr := strings.Clone(emptyStr)
	fmt.Println(unsafe.Pointer(&emptyStr) == unsafe.Pointer(&clonedEmptyStr)) // false
	fmt.Println(emptyStr == clonedEmptyStr)                                   // true
}

func CompareUsage() {

	fmt.Println(strings.Compare("abcdefg", "abc")) // 1
	fmt.Println(strings.Compare("a", "a"))         // 0
	fmt.Println(strings.Compare("a", "b"))         // -1
	fmt.Println(strings.Compare("b", "a"))         // 1

	fmt.Println()

	// The faster ops than strings.Compare()
	fmt.Println("abcdefg" != "abc")
	fmt.Println("a" == "a")
	fmt.Println("b" > "a")
	fmt.Println("a" < "b")

}

func ContainsUsage() {
	fmt.Println(strings.Contains("seafood", "foo")) // true
	fmt.Println(strings.Contains("sea", "foo"))     // false
	fmt.Println(strings.Contains("seafood", ""))    // true
	fmt.Println(strings.Contains("", ""))           // true
}

func ContainsAnyUsage() {
	fmt.Println(strings.ContainsAny("abcdefg", "qwerty")) // true {e}
	fmt.Println(strings.ContainsAny("abcdefg", "qw"))     // false {}

	fmt.Printf("%c\n", 0xfeff)                                             // ""
	fmt.Println(strings.ContainsAny("abcdefg", fmt.Sprintf("%x", 0xfeff))) // true {0xfeff}
	fmt.Println(strings.ContainsAny("abcdefg", ""))                        // false {}

	fmt.Println()
	fmt.Println(strings.ContainsAny("team", "i"))     // false {}
	fmt.Println(strings.ContainsAny("team", "ui"))    // false {}
	fmt.Println(strings.ContainsAny("ure", "ui"))     // true {u}
	fmt.Println(strings.ContainsAny("failure", "ui")) //true {u, i}

	fmt.Printf("%c\n", 0xfeff)                                             // ""
	fmt.Println(strings.ContainsAny("abcdefg", fmt.Sprintf("%x", 0xfeff))) // true {0xfeff}
	fmt.Println(strings.ContainsAny("abcdefg", ""))                        // false {}
	fmt.Println(strings.ContainsAny("", ""))                               // false {}

}

func ContainsFuncUsage() {
	var (
		f = func(r rune) bool {
			return r == 'a' || r == 'e' || r == 'o' || r == 'i' || r == 'u'
		}
	)

	fmt.Println(strings.ContainsFunc("aboba", f)) // true
	fmt.Println(strings.ContainsFunc("qwrt", f))  // false

	f = func(r rune) bool {
		return r != 'a' && r != 'e' && r != 'o' && r != 'i' && r != 'u'
	}
	fmt.Println(strings.ContainsFunc("qwerty", f)) // true
	fmt.Println(strings.ContainsFunc("aeoiu", f))  // false
}

func ContainsRuneUsage() {
	fmt.Println(strings.ContainsRune("abcdefg", 0))  // false
	fmt.Println(strings.ContainsRune("abcdefg", 97)) // true

	fmt.Println(strings.ContainsRune("abcdefg", 'a')) // true
	fmt.Println(strings.ContainsRune("abcdefg", '-')) // false
}

func CountUsage() {
	fmt.Println(strings.Count("aaaaaaaaa", "aaa"))  // "aaa" x3 = 9 symbols -> the result is 3
	fmt.Println(strings.Count("aaaaaaaaa", "aaaa")) // "aaa" x3 = 9 symbols -> the result is 2
	fmt.Println(strings.Count("Привет, мир!", ""))  // The result is 12 (the count of Unicode points) + 1
}

func CutUsage() {
	var (
		before, after string
		found         bool
	)

	before, after, found = strings.Cut("abcDEfg", "DE")
	fmt.Println(fmt.Sprintf("%q, %q, %t", before, after, found))

	before, after, found = strings.Cut("abcDEfg", "XXX")
	fmt.Println(fmt.Sprintf("%q, %q, %t", before, after, found)) // The after is an empty strings, since the sep hasn't been found
}

func CutPrefixUsage() {
	var (
		after string
		found bool
	)

	after, found = strings.CutPrefix("CCCabcdefg", "XXX")
	fmt.Printf("%q, %t\n", after, found)

	after, found = strings.CutPrefix("XXXabcdefg", "XXX")
	fmt.Printf("%q, %t\n", after, found)

	after, found = strings.CutPrefix("XXXabcdefg", "")
	fmt.Printf("%q, %t\n", after, found) // The prefix is an empty strings, so the result is the source string and true
}

func CutSuffixUsage() {
	var (
		before string
		found  bool
	)

	before, found = strings.CutSuffix("abcdefgXXX", "XXX")
	fmt.Printf("%q, %t\n", before, found)

	before, found = strings.CutSuffix("abcdefgXXX", "")
	fmt.Printf("%q, %t\n", before, found) // Since the suffix is an empty strings the result is the source string and true

	before, found = strings.CutSuffix("abcdefgCCC", "XXX")
	fmt.Printf("%q, %t\n", before, found)

}

// Turn strings into the same case and compare them
func EqualFoldUsage() {
	fmt.Println(strings.EqualFold("GO", "go")) // true
	fmt.Println(strings.EqualFold("ab", "Ab")) // true
	fmt.Println(strings.EqualFold("ß", "ss"))  // falses
}

func FieldsUsage() {
	fmt.Println(strings.Fields("a b c d e f g"))
	fmt.Println(strings.Fields("            a b           \tc"))
	fmt.Println(strings.Fields("    ")) // Since the strings consists of only spaces, the result is an empty slice
}

func FieldsFuncUsage() {
	var (
		// We will consider the symbols that are not letters/digits as a separators
		f = func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsNumber(c)
		}
	)

	fmt.Printf("Fields are %q\n", strings.FieldsFunc("             foo1;bar2;baz3...", f)) // Fields are ["foo1" "bar2" "baz3"]
}

func HasPrefixUsage() {
	fmt.Println(strings.HasPrefix("XXXabcdefg", "XX"))
	fmt.Println(strings.HasPrefix("XXabcdefg", "XX"))
	fmt.Println(strings.HasPrefix("abcdefg", ""))
	fmt.Println(strings.HasPrefix("", ""))
	fmt.Println(strings.HasPrefix("abcdefg", "OOO"))
}

func HasSuffixUsage() {
	fmt.Println(strings.HasSuffix("abcdefgXXX", "XX"))
	fmt.Println(strings.HasSuffix("abcdefgXXX", "XX"))
	fmt.Println(strings.HasSuffix("abcdefg", ""))
	fmt.Println(strings.HasSuffix("", ""))
	fmt.Println(strings.HasSuffix("aaaaaaaaa", "qwerty"))
}

func IndexUsage() {
	fmt.Println(strings.Index("sssA", "A"))   // 3
	fmt.Println(strings.Index("XXXaaa", "X")) // 0
	fmt.Println(strings.Index("OOOXXX", "P")) // Since the "P" is not presented in s, the result is -1
	fmt.Println(strings.Index("", ""))        // 0
}

func IndexByteUsage() {
	fmt.Println(strings.IndexByte("AAAA", 'B')) // -1
	fmt.Println(strings.IndexByte("AAAA", 'A')) // 0
	// fmt.Println(strings.IndexByte("abcdefg", '') // illegal rune literal
	fmt.Println(strings.IndexByte("ABCDEFG", 'G')) // 6
}

func IndexRuneUsage() {
	fmt.Println(strings.IndexRune("abcdefg", 'd'))                                                     // 3
	fmt.Println(strings.IndexRune("abcdefg", 'q'))                                                     // -1
	fmt.Println(strings.IndexRune(string([]byte{0xae, 0xed, 0xed, 0xed, 0xff, 0xff}), utf8.RuneError)) // 0
}

func IndexAnyUsage() {
	fmt.Println(strings.IndexAny("chicken", "aeoiu")) // 2
	fmt.Println(strings.IndexAny("qwrty", "aeoui"))   // -1
	fmt.Println(strings.IndexAny("qwerty", "t"))      // 4
}

func JoinUsage() {
	fmt.Println(strings.Join([]string{"", "cat,", "dog,", "bird,", "fish"}, " an animal: "))
	fmt.Println(strings.Join([]string{"A", "B", "C", "D", "E", "F"}, ""))
}

func LastIndexUsage() {
	fmt.Println(strings.LastIndex("gophers are the programmers who use go language", "go")) // 36
	fmt.Println(strings.LastIndex("abcdefg is the eng alphabet", ""))                       // 27
	fmt.Println(strings.LastIndex("abcdefg", "xyz"))                                        // -1
}

func LastIndexAnyUsage() {
	fmt.Println(strings.LastIndexAny("abcdefg abcdefg a a a a a a a a a a a a a a", "a"))    // 42
	fmt.Println(strings.LastIndexAny("abcdefg abcdefg a a a a a a a a a a a a a a", "defg")) // 14
	fmt.Println(strings.LastIndexAny("aaaaaaaaaa", "bcdefg"))                                // -1
}

func LastIndexByteUsage() {
	fmt.Println(strings.LastIndexByte(fmt.Sprintf("%c, %c, %c, abcdefg gdefv", 0xfffd, 0xfffd, 0xfffd), 'g')) // 23
	fmt.Println(strings.LastIndexByte("abcdefgssseq", 'x'))                                                   // -1
}

func LastIndexFuncUsage() {
	fmt.Println(strings.LastIndexFunc("abcdefg abcdefg abcdefg abcdefg abcdefg abcdefg", func(r rune) bool { return r > 'e' && r < 'g' })) // 45
	fmt.Println(strings.LastIndexFunc("abcdefg abcdefg abcdefg abcdefg abcdefg abcdefg", func(r rune) bool { return r == 0xfffd }))        // -1
}

func MapUsage() {
	fmt.Println(strings.Map(func(r rune) rune {
		if r == 0xfffd {
			return -1
		}

		if r != 'a' && r != 'e' && r != 'o' && r != 'i' && r != 'u' {
			return '-'
		}

		return r
	},
		"chicken",
	))

	fmt.Println(strings.Map(
		func(r rune) rune {
			if r == 0xfffd || unicode.IsSpace(r) {
				return -1
			}

			return r
		},
		fmt.Sprintf("%[1]cabcedfg, ggggg, %[1]c%[1]c%[1]c%[1]c%[1]c%[1]c%[1]c%[1]c%[1]c", 0xfffd),
	))
}

func RepeatUsage() {
	// fmt.Println(strings.Repeat("a", 2<<61)) // panic "strings: Repeat output length overflow"
	// fmt.Println(strings.Repeat("a", -1)) // panic "strings: negative Repeat count"

	fmt.Println(strings.Repeat("Bebra\n", 16))
}

func ReplaceUsage() {
	fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2))
	fmt.Println(strings.Replace("oink, oink, oink", "oink", "xyz", -1))

	fmt.Println(strings.Replace("oink, oink, oink", "", "|", 3)) // Starting from the 0th index the emptinesses will be replaced with the new
	// N times
	fmt.Println(strings.Replace("oink, oink, oink", "", "|", -1)) // Starting from the 0th index the emptinesses will be replaced with the
	// new up to the end of the strings
}

func ReplaceAllUsage() {
	fmt.Println(strings.ReplaceAll("Привет, Привет, Привет!", "", "|"))
	fmt.Println(strings.ReplaceAll("Привет, Привет, Привет!", "При", "Пока"))
	fmt.Println(strings.ReplaceAll("XYZXYZXYZXYZXYZXYZXYZXYZXYZ", "XYZ", "|"))
	fmt.Println(strings.ReplaceAll("AAAAAAAAAAAA", "AAAAA", "|"))
}

func SplitUsage() {
	fmt.Printf("% q\n", strings.Split("a|b|c|d|e|f|g", "|"))
	fmt.Printf("% q\n", strings.Split("a|b|c|d|e|f|g", ""))
	fmt.Printf("% q\n", strings.Split(" abcd ", ""))
	fmt.Printf("% q\n", strings.Split("", ""))          // The []. Because of emptiness of both s and sep
	fmt.Printf("% q\n", strings.Split("", "badbadbad")) // [""]
}

func SplitNUsage() {
	fmt.Printf("% q\n", strings.SplitN("a|b|c|d|e|f|g", "|", 3))    // ["a" "b" "c|d|e|f|g"]
	fmt.Printf("% q\n", strings.SplitN("a|b|c|d|a|f|g", "|", 0))    // []
	fmt.Printf("% q\n", strings.SplitN("a|b|c|d|a|f|g", "|", 1))    // ["a|b|c|d|a|f|g"]
	fmt.Printf("% q\n", strings.SplitN("a|b|c|d|e|f|g", "xyz", -1)) // ["a|b|c|d|e|f|g"]
	fmt.Printf("% q\n", strings.SplitN("a|b|c|d|e|f|g", "", 5))     // ["a" "|" "b" "|" "c|d|e|f|g"]
	fmt.Printf("% q\n", strings.SplitN("", "", -1))                 // []
	fmt.Printf("% q\n", strings.SplitN("", "", 1))                  // []
}

func SplitAfterUsage() {
	fmt.Printf("% q\n", strings.SplitAfter("a|b|c|d|e|f|g", "|"))   // ["a|" "b|" "c|" "d|" "e|" "f|" "g"]
	fmt.Printf("% q\n", strings.SplitAfter("a|b|c|d|a|f|g", "a"))   // ["a" "|b|c|d|a" "|f|g"]
	fmt.Printf("% q\n", strings.SplitAfter("a|b|c|d|e|f|g", "xyz")) // ["a|b|c|d|e|f|g"]
	fmt.Printf("% q\n", strings.SplitAfter("a|b|c|d|e|f|g", ""))    // ["a" "|" "b" "|" "c" "|" "d" "|" "e" "|" "f" "|" "g"]
	fmt.Printf("% q\n", strings.SplitAfter("", ""))                 // []
}

func SplitAfterNUsage() {
	fmt.Printf("% q\n", strings.SplitAfterN("a|b|c|d|e|f|g", "|", 3))    // ["a|" "b|" "c|d|e|f|g"]
	fmt.Printf("% q\n", strings.SplitAfterN("a|b|c|d|a|f|g", "a", 0))    // []
	fmt.Printf("% q\n", strings.SplitAfterN("a|b|c|d|e|f|g", "xyz", -1)) // ["a|b|c|d|e|f|g"]
	fmt.Printf("% q\n", strings.SplitAfterN("a|b|c|d|e|f|g", "", -1))    // ["a" "|" "b" "|" "c" "|" "d" "|" "e" "|" "f" "|" "g"]
	fmt.Printf("% q\n", strings.SplitAfterN("", "", -1))                 // []
}

func ToLowerUsage() {
	fmt.Println(strings.ToLower("ABCDEFG"))
	fmt.Println(strings.ToLower("12345AbCdE"))
}

func ToLowercaseSpecialUsage() {
	fmt.Println(strings.ToLowerSpecial(unicode.TurkishCase, "Önnek İş")) // önnek iş
}

func ToTitleUsage() {
	fmt.Println(strings.ToTitle("it's titled string")) // IT'S TITLED STRING
}

func ToTitleSpecialUsage() {
	fmt.Println(strings.ToTitleSpecial(unicode.TurkishCase, "dünyanın ilk borsa yapısı Aizonai kabul edilir")) // DÜNYANIN İLK BORSA YAPISI AİZONAİ
	// KABUL EDİLİR

}

func ToUpperUsage() {
	fmt.Println(strings.ToUpper("abcdefg qwerty poiu"))
	fmt.Println(strings.ToUpper("12345AbCdE"))
	fmt.Println(strings.ToUpper("привет, мир!")) // ПРИВЕТ, МИР!
}

func ToUpperSpecialUsage() {
	fmt.Println(strings.ToUpperSpecial(unicode.TurkishCase, "örnek iş")) // ÖRNEK İŞ
}

func ToValidUTF8Usage() {
	fmt.Printf("%s\n", strings.ToValidUTF8("\xff\xff\xffabc\xff\xff\xff", "\uFFFD")) // �abc�
	fmt.Printf("%s\n", strings.ToValidUTF8("a\xffb\xC0\xAFc\xff", "XXX"))            // aXXXbXXXcXXX
	fmt.Printf("%s\n", strings.ToValidUTF8("\xed\xa0\x80", "abc"))                   // abc
}

func TrimUsage() {
	fmt.Println(strings.Trim("ZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCQWERTYZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCZXCZXC", "ZXC")) // QWERTY
	fmt.Println(strings.Trim("QWERTY", ""))                                                                                  // QWERTY
	fmt.Println(strings.Trim("QWERTY", "QWERTY"))                                                                            // ""
}

func TrimFuncUsage() {
	fmt.Println(strings.TrimFunc("ZXCZXCZXCZXCZXCZXCabcdefgZXCZXCZXCZXCZXCZXC", func(r rune) bool {
		if r == 'Z' || r == 'X' || r == 'C' {
			return true
		}

		return false
	}))
}

func TrimLeftFuncUsage() {
	fmt.Println(strings.TrimLeftFunc("xxxQWERTYxxx", func(r rune) bool {
		return r == 'x'
	})) // QWERTYxxx
}

func TrimPrefixUsage() {
	fmt.Println(strings.TrimPrefix("xxxQWERTYxxx", "xxx")) // QWERTYxxx
	fmt.Println(strings.TrimPrefix("zzzQWERTYxxx", "xxx"))
}

func TrimRightUsage() {
	fmt.Println(strings.TrimRight("abcdefgXXX", "X"))      // abcdefg
	fmt.Println(strings.TrimRight("abcdefgZXCZXC", "ZXC")) // abcdefg
	fmt.Println(strings.TrimRight("abcdefgZXCZXC", "ZC"))  // abcdefgZXCZX
}

func TrimRightFuncUsage() {
	fmt.Println(strings.TrimRightFunc("abcde111fgs12417916429124412", func(r rune) bool {
		return unicode.IsDigit(r)
	})) // abcde111fgs
}

func TrimSpaceUsage() {
	fmt.Println(strings.TrimSpace("            \t\n\rwdoqjdwijodqwijoqdwjiodqwjidwq\t         \n")) // wdoqjdwijodqwijoqdwjiodqwjidwq
}

func TrimSuffixUsage() {
	fmt.Println(strings.TrimSuffix("abcdefgXXX", "XXX")) // abcdefg
	fmt.Println(strings.TrimSuffix("abcdefgXXX", ""))    // abcdefgXXX
}

func BuilderZeroValueUsage() {
	var (
		builder strings.Builder
	)

	for i := 0; i < 11; i++ {
		fmt.Fprintf(&builder, "%d...\n", i)
	}

	builder.WriteString("Launching!")

	fmt.Println(builder.String())
}

func BuilderJoinMultipleByteSlices(byteSlices ...[]byte) string {
	var (
		strBuilder strings.Builder
	)

	for _, v := range byteSlices {
		strBuilder.Write(v)
	}

	return strBuilder.String()

}

func BuilderJoinStrings(strs ...string) string {
	var (
		strBuilder strings.Builder
	)

	for _, v := range strs {
		strBuilder.WriteString(v)
	}

	return strBuilder.String()
}

func BuilderJoinBytes(bts ...byte) string {
	var (
		strBuilder strings.Builder
	)

	for _, v := range bts {
		strBuilder.WriteByte(v)
	}

	return strBuilder.String()
}

func BuilderJoinRunes(runes ...rune) string {
	var (
		strBuilder strings.Builder
	)

	for _, r := range runes {
		strBuilder.WriteRune(r)
	}

	return strBuilder.String()
}

func BuilderResetUsage(strs []string) {
	var (
		sb strings.Builder
	)

	for _, s := range strs {
		fmt.Printf("B: cap(): %d, len(): %d\t", sb.Cap(), sb.Len())
		sb.WriteString(s)
		fmt.Printf("A: cap(): %d, len(): %d\n", sb.Cap(), sb.Len())
	}
	resultingStr := sb.String()
	fmt.Println(resultingStr)

	sb.Reset()
	fmt.Printf("RESET: cap(): %d, len(): %d\t", sb.Cap(), sb.Len())

	strToBeAppended := `abcdefg`
	for _, r := range strToBeAppended {
		sb.WriteRune(r)
	}

	fmt.Println(strToBeAppended)
}

func BuildersComparison() {
	var (
		sb1, sb2 strings.Builder
	)

	sb1.WriteString("abcdefg")
	sb2.WriteString("abcdefg")

	fmt.Println(unsafe.Pointer(&sb1) == unsafe.Pointer(&sb2)) // false
}
