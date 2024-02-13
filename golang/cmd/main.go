package main

import (
	"fmt"
	. "golang/pkg/chapters/chapter7"
)

func main() {
	filename := "file.xml"

	root, err := ConstructTree("file.xml")
	if err != nil {
		fmt.Printf("constructing tree from file: %s; %s", filename, err)
	}

	PrintTree(root)
}
