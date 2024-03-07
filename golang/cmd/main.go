package main

import (
	"fmt"
)

type Speaker interface {
	Speak()
}

type Person struct {
	firtsname string
	lastname  string
	age       int
}

type Student struct {
	Person

	marks map[string]int
}

func (s *Student) Speak() {
	fmt.Println("Give fucking five!")
}

type Hardworker struct {
	Person
	schedule map[string][]string
	isGay    bool
}

func (c *Hardworker) Speak() {
	fmt.Println("Я должен денег...Много...")
}

func main() {

	george := Student{
		Person{firtsname: "George", lastname: "Odin", age: 22},
		map[string]int{},
	}

	vlad := Hardworker{
		Person{firtsname: "Vlad", lastname: "Riznyk", age: 22},
		make(map[string][]string),
		true,
	}

	speakers := []Speaker{}
	speakers = append(speakers, &george)
	speakers = append(speakers, &vlad)

	for _, speaker := range speakers {
		speaker.Speak()
	}

	var speaker Speaker
	var anotherSpeaker Speaker

	speaker = &george
	anotherSpeaker = &vlad

	if georgeCopy, ok := speaker.(*Student); ok {
		fmt.Println(georgeCopy)
	}

	if vladCopy, ok := anotherSpeaker.(*Hardworker); ok {
		fmt.Println(vladCopy)
	}

}
