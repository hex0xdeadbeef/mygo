package firtstphonebookimprovement

import (
	"fmt"
	"math/rand"
	"strconv"
)

var (
	data []Entry
)

const (
	MIN = 0
	MAX = 26
)

type Entry struct {
	Name    string
	Surname string
	Tel     string
}

func search(key string) *Entry {
	for i, v := range data {
		if v.Tel == key {
			return &data[i]
		}
	}

	return nil
}

func list() {
	for _, v := range data {
		fmt.Println(v)
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(l int64) string {
	startChar := 'A'
	temp := ""
	var i int64 = 1

	for {
		myRand := random(MIN, MAX)
		newChar := string(byte(startChar) + byte(myRand))
		temp += newChar
		if i == l {
			break
		}

		i++
	}
	return temp
}

func populate(n int, s []Entry) {
	for i := 0; i < n; i++ {
		name := getString(4)
		surname := getString(5)
		n := strconv.Itoa(random(100, 199))
		data = append(data, Entry{Name: name, Surname: surname, Tel: n})
	}
}
