package deferstmt

import "fmt"

func DeferUsing() {
	for i := 0; i < 5; i++ {
		defer fmt.Println(i)
	}
}

func trace(s string) string {
	fmt.Println("entering:", s)
	return s
}

func untrace(s string) {
	fmt.Println("leaving:", s)
}

func a() {
	defer untrace(trace("a"))
	fmt.Println("in a")

}

func B() {
	defer untrace(trace("b"))
	fmt.Println("in b")
	a()
}
