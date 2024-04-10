package chapter7

import (
	"fmt"
)

type Person struct {
	name string
}

func (person *Person) talk() {
	fmt.Println("Hi, I'm", person.name)
}

/*
type Android struct {
	person Person
	name string
}*/

/*
В данном случае некорректно говорить, что Android содержит в себе человека. Скорее Андроид это человек.
В таких случаях необходимо вкладывать анонимное поле типа Person:

Посредством обращения через имя типа Person мы можем использовать методы типа Person
*/

type Android struct {
	Person
	model string
}

func (android *Android) setModel(model string) {
	android.model = model
}

func (android *Android) sayItsModel() {
	fmt.Println("Я являюсь моделью с номером:", android.model)
}

func usingTypeMethodsWithAnonymousFieldOrDirectly() {
	firstAndroid := new(Android)
	firstAndroid.setModel("13.1")

	firstAndroid.Person.talk() // call method with anonymous field
	firstAndroid.sayItsModel()

	firstAndroid.talk() // call method directly
	firstAndroid.sayItsModel()
}
