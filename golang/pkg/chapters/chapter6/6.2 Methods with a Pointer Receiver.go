package chapter6

import "fmt"

/*
Scales the object fields.
*/
func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

func (p *Point) String() string {
	return fmt.Sprintf("(X: %g, Y: %g)", p.X, p.Y)
}

// Attempts to attach method to a pointer/interface type.
type A *int

/*
func (A) f() {} // invalid receiver type A (POINTER or interface type)
*/

type B interface{}

/*
func (B) f() {} // invalid receiver type B (pointer or INTERFACE type)
*/

func ScaleByUsing() {
	// Use pointer object
	pointPointer := &Point{X: 1, Y: 10}
	pointPointer.ScaleBy(100)
	fmt.Println(pointPointer)

	// Create a pointer explicitly
	point1 := Point{X: 2, Y: 3}
	pointPtr := &point1
	pointPtr.ScaleBy(333)
	fmt.Println(pointPtr)

	// Use place-in pointer
	point2 := Point{X: 7, Y: 17}
	(&point2).ScaleBy(777)
	fmt.Println(&point2)

	/*
		Use a non-pointer object directly.
		The pointer will be taken by the compiler implicitly.
	*/
	point3 := Point{X: 13, Y: 13}
	point3.ScaleBy(13)
	fmt.Println(point3)

	/*
		Attempt to apply implicit pointer conversion on a literal value
	*/
	/* Point{-1, -1}.ScaleBy(100) // cannot call pointer method ScaleBy on Point */

	/*
		Point's literal calls the method of "Point" type that takes the receiver of type Point
	*/
	// Point{X: 10, Y: 10}.Distance(point3)
}

func PointerObjectCallNonPointerReceiverMethod() {
	point := Point{1, 1}
	pointPointer := &point

	// (p Point) Distance(q Point) float64, but receiver now has the *Point type.
	pointPointer.Distance(Point{3, 3})
	(&point).Distance(Point{3, 3})
}
