package chapter8_9

import (
	"fmt"
	"strings"
)

// "strings" package
func containsUsing(source string, subsequence string) bool {
	return strings.Contains(source, subsequence)
}

func countUsing(source string, subsequence string) int {
	return strings.Count(source, subsequence)
}

func hasPrefixUsing(source string, prefix string) bool {
	return strings.HasPrefix(source, prefix)
}

func hasSuffixUsing(source string, prefix string) bool {
	return strings.HasSuffix(source, prefix)
}

func indexUsing(source string, subsequence string) int {
	return strings.Index(source, subsequence) // Index may return -1 if substring you want to find isn't occured in this string
}

// The important thing is that the followed functions return a new string instead of changing a source string
func joinUsing(stringsSlice []string, separator string) string {
	return strings.Join(stringsSlice, separator)
}

func repeatUsing(initial string, timesToRepeat int) string {
	return strings.Repeat(initial, timesToRepeat)
}

func replaceUsing(initial string, sequenceToReplace string, sequenceToBeReplaced string, timesToReplace int) string {
	return strings.Replace(initial, sequenceToReplace, sequenceToBeReplaced, timesToReplace) // You can passed in the value -1 to replace  all substrings string contains
}

func splitUsing(source string, spliterator string) []string {
	return strings.Split(source, spliterator) // If you want get an array wich contains all symbols you should use the "" spliterator
}

func toLowerUsing(source string) string {
	return strings.ToLower(source)
}

func toUpperUsing(source string) string {
	return strings.ToUpper(source)
}

// Converting a string into []byte and vice versa
func convertingStringToByteSlice(source string) []byte {
	return []byte(source)
}

func convertingByteArrayToString(bytes []byte) string {
	return string(bytes)
}

func String_All_Nethods_Using() {
	// "strings" package
	fmt.Println("strings package")
	fmt.Println("{")

	// Methods that return bool values
	fmt.Println(containsUsing("test", "es"))
	fmt.Println()
	fmt.Println(countUsing("test", "t"))
	fmt.Println()
	fmt.Println(hasPrefixUsing("test", "tes"))
	fmt.Println()
	fmt.Println(hasSuffixUsing("test", "st"))
	fmt.Println()
	fmt.Println(indexUsing("test", "t"))
	fmt.Println()

	// Methods that returns a new string instead of modifying a source one
	fmt.Println(joinUsing([]string{"a", "b", "c", "d"}, "+"))
	fmt.Println()
	fmt.Println(repeatUsing("bebra ", 10))
	fmt.Println()
	fmt.Println(replaceUsing("abobaabobaabobaabobaabobaaboba", "bo", "BO", -1))
	fmt.Println()
	fmt.Println(splitUsing("ABCDEFG", ""))
	fmt.Println()
	fmt.Println(toLowerUsing("CAPITALIZED MESSAGE"))
	fmt.Println()
	fmt.Println(toUpperUsing("capitalized message"))
	fmt.Println()

	// Converting a string to abyte array string and vice versa
	fmt.Println(convertingStringToByteSlice("abcdefg")) // We'll get a byte array
	fmt.Println()
	fmt.Println(convertingByteArrayToString([]byte{97, 98, 99, 100, 101, 102, 103})) // We'll get a string
	fmt.Println()
	fmt.Println("}")
	fmt.Println()
}
