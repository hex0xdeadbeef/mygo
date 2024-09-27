package chapter5

import (
	"fmt"

	"golang.org/x/net/html"
	. "golang.org/x/net/html"
)

func soleTitle(doc *Node) (title string, err error) {
	type bailout struct{}

	// The defer function that contains the recover logic
	defer func() {
		switch p := recover(); p {
		case nil:
		case bailout{}:
			err = fmt.Errorf("multiple title elements")
		default:
			panic(p)
		}

	}()

	forEachNode(doc,
		func(node *Node) {
			if node.Type == html.ElementNode &&
				node.Data == "title" &&
				// If a child has "title" data -> the title isn't single
				node.FirstChild != nil {
				if title != "" {
					panic(bailout{})
				}
				title = node.FirstChild.Data
			}

		},
		nil)

	if title == "" {
		return "", fmt.Errorf("no title element")
	}

	return title, nil
}

func NonReturnStatement() (name string) {
	defer func() {
		if p := recover(); p != nil {
			name = "John Doe"
		}
	}()

	panic("Harry Potter")
}
