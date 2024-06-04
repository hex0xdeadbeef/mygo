package constructorsandcompositeliterals

import "fmt"

type Person struct {
	id        int
	name      string
	isMarried bool
}

func New(id int, name string) *Person {
	return &Person{id: id, name: name, isMarried: true}
}

func PartiallyNew(id int, name string) *Person {
	return &Person{id: id, name: name}
}

func NoFieldsNewA() *Person {
	return new(Person)
}

func NoFieldsNewB() *Person {
	return &Person{}
}

func MakeUsing() {
	var (
		p *[]int = new([]int)
		v []int  = make([]int, 10, 100)
	)

	fmt.Println(p, v)
}
