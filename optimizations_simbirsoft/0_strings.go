package main

import (
	"math/rand"
	"strings"
)

/*
	STRINGS
1. We should use the srings.Builder for the purposes of  concatenations. It's faster than other methods because of dynamic buffer of bytes.
2. The strings.Builder object is efficient, but if we need to construct an object that is small enough and there are many workloads bound to this function, it's better to use this
	function. The performance in this case will be higher.

	The example of this function is:
		func GetFullname(name, secondName, phoneNumber string) string {
			return "Name: " + name +
				"Secondname: " + secondName +
				"Phone number: " + phoneNumber
		}
*/

var (
	tokenSymbolsSet []byte
)

func init() {
	const (
		digitsNumber   = 10
		digitStartByte = '0'

		lettersNumber      = 26
		lowercaseStartByte = 'a'
		uppercaseStartByte = 'A'
	)

	tokenSymbolsSet = make([]byte, 0, digitsNumber+lettersNumber*2)

	for i := 0; i < digitsNumber; i++ {
		tokenSymbolsSet = append(tokenSymbolsSet, digitStartByte+byte(i))
	}

	for i := 0; i < lettersNumber; i++ {
		tokenSymbolsSet = append(tokenSymbolsSet, lowercaseStartByte+byte(i), uppercaseStartByte+byte(i))
	}
}

func getTokenSet(amount int, size int) []string {
	var (
		tokens = make([]string, amount)
	)

	for i := 0; i < amount; i++ {
		tokens[i] = getToken(size)
	}

	return tokens
}

func getToken(tokenSize int) string {
	var (
		buf = make([]byte, 1<<8)
	)

	for i := 0; i < tokenSize; i++ {
		buf[i] = tokenSymbolsSet[rand.Intn(len(tokenSymbolsSet))]
	}

	return (string(buf))
}

func ConcatLinearlyString(strs []string) string {
	var (
		res string
	)

	for i := 0; i < len(strs); i++ {
		res += strs[i]
	}

	return res
}

func ConcatLinearlyBytes(strs []string) string {
	var (
		res = make([]byte, 0, 1<<8)
	)

	for i := 0; i < len(strs); i++ {
		res = append(res, []byte(strs[i])...)
	}

	return string(res)
}

func ConcatBuffer(strs []string) string {
	var (
		bytesCount int
		buf        strings.Builder
	)

	for i := 0; i < len(strs); i++ {
		bytesCount += len(strs[i])
	}

	buf.Grow(bytesCount)

	for i := 0; i < len(strs); i++ {
		buf.WriteString(strs[i])
	}

	return buf.String()
}

func GetFullname(name, secondName, phoneNumber string) string {
	return "Name: " + name +
		"Secondname: " + secondName +
		"Phone number: " + phoneNumber
}
