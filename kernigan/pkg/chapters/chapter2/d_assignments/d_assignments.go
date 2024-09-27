package d_assignments

import "fmt"

func TupleAssignmentOfThreeVariables() {
	var i, j, k = 2, 3, 5

	fmt.Println(i, j, k)
}

func implicitAssignment() {
	medals := []string{"gold", "silver", "platinum"}
	/*medals[0] = "gold"
	medals[1] = "silver"
	medals[2] = "platinum"
	*/
	fmt.Println(medals[0])
}