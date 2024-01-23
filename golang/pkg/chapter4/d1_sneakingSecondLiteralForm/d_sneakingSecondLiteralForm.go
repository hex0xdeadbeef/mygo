package d1_sneakingsecondliteralform

import (
	"fmt"
	a "golang/pkg/chapter4/d0_structs"
)

func AttemptToSneakAroundRule() {
	/* point := a.Point{x: 1, Y: 100} // unknown field x in struct literal of type d_structs.Point */
	person := a.Person{
		Firstname:       "Anya",
		Lastname:        "Zolotushkina",
		Age:             20,
		PassportID:      351911,
		DriverLicenseID: 191344,
		/*
		 partner:         &a.Person{Firstname: "Ivan", Lastname: "Ivanov"},
		 unknown field partner in struct literal of type d0_structs.Person
		*/
	}
	fmt.Printf("%[1]T, %[1]v", person)
}
