package chapter6

import (
	"fmt"
	"net/url"
)

type IntList struct {
	Value int
	Tail  *IntList
}

// nil parameter must be documented precisely
func (node *IntList) Sum() int {
	if node == nil {
		return 0
	}

	return node.Value + node.Tail.Sum()
}

func UsingURLValues() {
	m := url.Values{"lang": []string{"eng"}}

	/*
		The change won't be affected on the initial m, because inside the subRoutine() the local link will be redirected
	*/

	subRoutine(m)
	fmt.Println(m)

	m = nil
	firstLanguage := m.Get("lang")
	fmt.Println(firstLanguage) // ""

	// The attempt to append an element into a nil slice
	m.Add("lang", "ru") /* Exception has occurred: panic "assignment to entry in nil map" */

}

func subRoutine(values url.Values) {
	values.Add("lang", "ru")
	// This change doesn't affect on the original object, the local link is just redirected.
	values = nil
}
