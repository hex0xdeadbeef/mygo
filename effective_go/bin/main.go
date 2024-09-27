package main

import (
	"effectivego/cmd/forloop"
	"effectivego/cmd/gettersandsetters"
	"effectivego/cmd/redeclarationreassignment"
	"fmt"

	"effectivego/cmd/deferstmt"
)

func main() {

	getterUsing()

	redeclarationreassignment.ShortRepetitions()

	forloop.ForString()

	deferstmt.DeferUsing()
	deferstmt.B()
}

func getterUsing() {
	var (
		obj = gettersandsetters.New("")
	)
	obj.Owner()

	if obj.Owner() == "" {
		obj.SetOwner("Dmitriy")
		fmt.Println(obj.Owner())
	}

}
