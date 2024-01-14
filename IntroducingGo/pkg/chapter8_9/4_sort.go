package chapter8_9

import (
	"fmt"
	"math/rand"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (people ByAge) Len() int {
	return len(people)
}

func (person ByAge) Less(i, j int) bool {
	return person[i].Age < person[j].Age
}

func (person ByAge) Swap(i, j int) {
	person[i], person[j] = person[j], person[i]
}

func sortingPeople() {

	rustamsFactory := func(rustamsCount int) []Person {
		rustamsSlice := make([]Person, 0, rustamsCount)
		for i := 0; i < rustamsCount; i++ {
			rustamsSlice = append(rustamsSlice, Person{"Rustam", rand.Intn(46)})
		}
		return rustamsSlice
	}

	people := make([]Person, 0, 0)
	people = append(people, rustamsFactory(100)...)
	sort.Sort(ByAge(people))

	fmt.Println(people)
}

func Sort_All_Nethods_Using() {
	fmt.Println("{")
	sortingPeople()
	fmt.Println("}")
	fmt.Println()
}
