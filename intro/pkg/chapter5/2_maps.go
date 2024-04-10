package chapter5

import (
	"fmt"
)

func creatingAMap() {
	x := make(map[string]int)
	x["key"] = 10
	fmt.Println(x["key"])

	y := make(map[int]float64)
	y[1] = 1.11111
	fmt.Println(y[1])
}

func deletingMapElements() {
	a := make(map[int]float64)
	a[1] = 1.111

	delete(a, 1)
}

func poppingAnUnexistingElementFromAMap() {
	a := make(map[int]string)
	a[0] = "a"
	a[1] = "b"
	a[2] = "c"

	fmt.Println(a[1])
}

func preventingPoppingBasicElementForATypeFromMap() {
	a := make(map[int]string)
	a[0] = "a"
	a[1] = "b"
	a[2] = "c"

	// if ok != false {
	// 	fmt.Println(number)
	// } else {
	// 	fmt.Println("Значения в мапе не существует.")
	// }

	// Более простая запись
	number, ok := a[0]
	if ok {
		fmt.Println(number)
	}
}

func mapOfMaps() {
	chemicalElements := make(map[string]map[string]string)
	chemicalElements["H"] = map[string]string{"name": "Hydrogen", "state": "gas"}
	chemicalElements["He"] = map[string]string{"name": "Helium", "State": "gas"}

	fmt.Println(chemicalElements["H"])
	elementState, elementStateExist := chemicalElements["H"]["state"]
	if elementStateExist {
		fmt.Println(elementState)
	}

}
