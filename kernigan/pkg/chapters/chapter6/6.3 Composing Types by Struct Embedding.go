package chapter6

import (
	"fmt"
	"image/color"
	"sync"
)

type ColoredPoint struct {
	*Point
	Color color.RGBA
}

// "Point" Distance() wrapper
func (from *ColoredPoint) Distance(to *ColoredPoint) float64 {
	return from.Point.Distance(*to.Point)
}

// "Point" ScaleBy() wrapper
func (c *ColoredPoint) ScaleBy(scale float64) {
	c.Point.ScaleBy(scale)
}

func SelectFieldAndMethods() {
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}

	colPoint1 := &ColoredPoint{&Point{X: 1, Y: 1}, red}
	colPoint2 := &ColoredPoint{&Point{X: 10, Y: 10}, blue}

	/*
		Internal fields passing
	*/
	fmt.Println(colPoint1.Point.Distance(*colPoint2.Point))
	colPoint2.Point.ScaleBy(31)
	fmt.Println(colPoint1.Point.Distance(*colPoint2.Point))

	fmt.Println()
	/*
		Wrappers using
	*/
	fmt.Println(colPoint1.Distance(colPoint2))
	colPoint1.ScaleBy(333)
	fmt.Println(colPoint1.Distance(colPoint2))

	fmt.Println()
	/* Redirection the Point field of the colPoint1 */
	colPoint1.Point = colPoint2.Point
	fmt.Println(colPoint1.Distance(colPoint2))
	fmt.Println(colPoint1.Color)
	fmt.Println(colPoint2.Color)
}

type MultipleAnonymousFields struct {
	Point
	color.RGBA
}

func (maf *MultipleAnonymousFields) PrintSun() {
	fmt.Println("Hello, Sun!")
}

type Z struct {
	X
	C
}
type X struct {
}

func (x *X) f() {}

type C struct {
}

func (c *C) f() {}

func MethodsAmbiguity() {
	zInstance := &Z{}

	/*
		zInstance.f() // ambiguous selector zInstance.f
	*/

	fmt.Println(zInstance)
}

// The cache has all the methods of sync.Mutex
var cache = struct {
	sync.Mutex
	mapping map[string]string
}{
	mapping: make(map[string]string),
}

func Lookup(key string) string {
	cache.Lock()
	v := cache.mapping[key]
	cache.Unlock()
	return v
}
