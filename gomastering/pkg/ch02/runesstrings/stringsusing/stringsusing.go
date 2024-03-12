package stringsusing

import (
	"fmt"
	"strings"
	"unicode"
)

func EqualFoldUsing() {
	fmt.Println(strings.EqualFold("Dmitriy", "DMITriy"))
	fmt.Println(strings.EqualFold("Dmitriy", "DMITri"))
}

func IndexUsing() {
	fmt.Println(strings.Index("Dmitriy", "itr"))

	fmt.Println(strings.Index("Dmitriy", "f"))
}

func HasPrefixSuffixUsing() {
	fmt.Println(strings.HasPrefix("Dmitriy", "Dm"))
	fmt.Println(strings.HasPrefix("Dmitriy", "An"))

	fmt.Println(strings.HasSuffix("Dmitriy", "triy"))
	fmt.Println(strings.HasSuffix("Dmitriy", "kaya"))

}

func FieldsUsing() {
	t := strings.Fields("This is a string!")
	fmt.Printf("Fields: %#v | len: %d\n", t, len(t))

	t = strings.Fields("ThisIs a\tstring!")
	fmt.Printf("Fields: %#v | len: %d\n", t, len(t))
}

func SplitUsing() {
	fmt.Printf("%#v\n", strings.Split("abcd efg", ""))
	fmt.Printf("%s\n", strings.Replace("abcd efg", "", "_", -1))
	fmt.Printf("%s\n", strings.Replace("abcd efg", "", "_", 4))

}

func SplitAfterUsing() {
	fmt.Printf("%#v\n", strings.SplitAfter("123++432++", "++"))

}

func TrimFuncUsing() {
	trimFunc := func(c rune) bool {
		return !unicode.IsLetter(c)
	}

	fmt.Printf("%#v\n", strings.TrimFunc("123 abc ABC \t .", trimFunc))
}
