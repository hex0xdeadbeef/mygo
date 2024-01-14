package chapter1

import (
	"fmt"
)

type Person struct {
	Lastname  string
	Firstname string
}

// This is a comment

func Hello_Go() {
	fmt.Println("Hello, Dmitriy")
	person := Person{Lastname: "Mamykin", Firstname: "Dmitriy"}
	fmt.Println(person.Firstname + " " + person.Lastname)

}
