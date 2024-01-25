package d2_embeddinganonymous_test

import "fmt"

/* Struct embedding and Anonymous fields */

type point struct {
	X, Y int
}

// NO DECOMPOSED STRUCTURE
/*


type Circle struct {
	X, Y, Radius int
}

type Wheel struct {
	X, Y, Radius, Spokes int
}

func NonEmbeddedUse() {
	leftWheel := Wheel{X: 1, Y: 1, Radius: 100, Spokes: 30}
	fmt.Printf("Wheel parameters: %v", leftWheel)
}


*/

// DECOMPOSED STRUCTURE
/*


type Circle struct {
	center Point
	Radius int
}

type Wheel struct {
	Body   Circle
	Spokes int
}

func DecomposedStructFeatures() {
	leftWheel := Wheel{Body: Circle{Point{X: 0, Y: 0}, 10}, Spokes: 100}
	fmt.Printf("Left wheel parameters.\nType: %[1]T\nValues: %[1]v", leftWheel)
}


*/

// ANONYMOUS FIELD
type Circle struct {
	point
	Radius int
}

type Wheel struct {
	Circle
	Spokes int
}

func AnonymousFields() {
	leftWheel := Wheel{Circle{point{1, 2}, 10}, 15}
	fmt.Printf("\nThe left wheel parameters:\nType: %[1]T\nValues: %#[1]v\n", leftWheel)

	rightWheel := Wheel{
		Circle: Circle{
			point:  point{X: 1, Y: 2},
			Radius: 10,
		},
		Spokes: 100,
	}

	fmt.Printf("\nThe right wheel parameters:\nType: %[1]T\nValues: %#[1]v\n", rightWheel)

}

func SetAnonymousFieldParameter() {
	backLeftValue := Wheel{
		Circle: Circle{
			point:  point{10, 10},
			Radius: 15,
		},
		Spokes: 250,
	}

	fmt.Printf("\nThe backLeftValue parameters BEFORE: \nType: %[1]T\nValues: %[1]v\n", backLeftValue)
	backLeftValue.Radius = 3
	fmt.Printf("\nThe backLeftValue parameters AFTER: \nType: %[1]T\nValues: %[1]v\n", backLeftValue)
}
