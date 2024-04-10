package chapter8_9

import (
	"container/list"
	"fmt"
)

func List_All_Nethods_Using() {

	creatingAList()
	fmt.Println()
}

func creatingAList() {
	var firstList list.List
	firstList.PushBack("бебра")
	firstList.PushBack(-1)
	firstList.PushFront(1)
	firstList.PushBack("xyz")

	for element := firstList.Front(); element != nil; element = element.Next() {
		fmt.Print(element.Value, " ")
	}

	fmt.Println()
}
