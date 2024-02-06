package chapter6

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

// The meethod of type Point
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// The simple function that isn't bound to a particular type
func DistanceFunc(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

/*
Path is named slice type
*/
type Path []Point

func (path Path) Distance() float64 {
	var distance float64
	for i := 0; i < len(path)-1; i++ {
		// path[i] uses its "DistanceFunc(p, q Point) float64" own method
		distance += path[i].Distance(path[i+1])
	}

	return distance
}

func GetPerimeter() {
	shapeCircuit := Path{
		{1, 1},
		{5, 1},
		{5, 4},
		{1, 1},
	}

	// shapeCircuit calls to the its own "(path Path) Distance() float64" method
	fmt.Println(shapeCircuit.Distance())
}
