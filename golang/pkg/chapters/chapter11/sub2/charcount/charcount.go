package charcount

import (
	"bytes"
	"io"
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	keys = [...]string{0: "l", 1: "d", 2: "o"}
)

type SymbolsCounter map[string]map[rune]int

// String turns SymbolsCounter instances into string in alphabetic order
func (sc *SymbolsCounter) String() string {
	var (
		curMap   map[rune]int
		keysList []string
		buf      = bytes.NewBuffer(make([]byte, 0))
	)

	for _, key := range keys {
		curMap = (*sc)[key]
		keysList = make([]string, 0)

		for key := range curMap {
			keysList = append(keysList, string(key))
		}

		sort.Strings(keysList)

		var runeKey rune
		for _, key := range keysList {
			buf.WriteString(key)
			runeKey = []rune(key)[0]
			buf.WriteString(strconv.Itoa(curMap[runeKey]))
		}
	}

	return buf.String()
}

// CharCount reads the string by rune and increment the corresponding section by one
func CharCount(inputStr string) SymbolsCounter {
	var (
		quantities             SymbolsCounter = SymbolsCounter(make(map[string]map[rune]int, 0))
		utflen                 [utf8.UTFMax + 1]int
		invalidCharactersCount = 0

		reader = strings.NewReader(inputStr)
	)

	quantities[keys[0]], quantities[keys[1]], quantities[keys[2]] = make(map[rune]int), make(map[rune]int), make(map[rune]int)

	for {
		character, bytesCount, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("reading rune: %v", err)
		}

		if character == unicode.ReplacementChar {
			invalidCharactersCount++
		}

		if unicode.IsLetter(character) {
			quantities[keys[0]][character]++
		} else if unicode.IsDigit(character) {
			quantities[keys[1]][character]++
		} else {
			quantities[keys[2]][character]++
		}

		utflen[bytesCount]++
	}

	return quantities
}
