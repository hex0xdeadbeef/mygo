package chapter6

import (
	"fmt"
	"time"
)

/*
Method values
*/
func ExtractMethodValue() {
	point1 := Point{1, 1}

	distanceFromPoint1 := point1.Distance

	fmt.Println(distanceFromPoint1(point1))
	point1.X = 100
	point1.Y = 100

	fmt.Println(distanceFromPoint1(point1))

	point2 := Point{4, 6}
	scalePoint2 := point2.ScaleBy
	scalePoint2(10)
	fmt.Println(point2)
}

type Rocket struct {
}

func (r *Rocket) Launch() {
	fmt.Println("Rocket has been successfully launched!")
}

func RocketStart() {
	rocket := new(Rocket)
	time.AfterFunc(time.Second, func() { rocket.Launch() })
}

/*
Method expressions
*/
func ExtractMethodExpression() {
	point1 := Point{1, 2}
	point2 := Point{3, 4}

	// Extract the method's expressions. Now the distance and scaleBy are usual functions.
	distance := Point.Distance
	scaleBy := (*Point).ScaleBy

	fmt.Println(distance(point1, point2))

	scaleBy(&point1, 100)
	fmt.Println(point1)
}

func (p *Point) Add(q Point) {
	p.X += q.X
	p.Y += q.Y
}

func (p *Point) Sub(q Point) {
	p.X -= q.X
	p.Y -= q.Y
}

func (path Path) TranslateBy(offset Point, add bool) {
	var op func(p *Point, q Point)
	if add {
		op = (*Point).Add
	} else {
		op = (*Point).Sub
	}

	for i := range path {
		op(&path[i], offset)
	}
}
