package chapter7

import (
	"fmt"
	"math"
)

type MultipleShape struct {
	shapes []Shape // This using of Shape means, that we can use an interface type as a field
}

func (multiShape *MultipleShape) area() float64 {
	total := 0.0
	for _, shape := range multiShape.shapes {
		total += shape.area()
	}
	return total
}

func (multiShape *MultipleShape) perimeter() float64 {
	total := 0.0
	for _, shape := range multiShape.shapes {
		total += shape.perimeter()
	}
	return total
}

func (multipleShape *MultipleShape) getShapes() []Shape {
	return multipleShape.shapes
}

func (multipleShape *MultipleShape) addShape(shapes ...Shape) {
	multipleShape.shapes = append(multipleShape.shapes, shapes...)
}

type Shape interface {
	area() float64 // The function must be implemented by types that implenet the respective interface
	perimeter() float64
}

type Circle struct {
	x, y, r float64
}

func (circle *Circle) area() float64 {
	return math.Pow(circle.r, 2) * math.Pi
}

func (circle *Circle) perimeter() float64 {
	return math.Pi * 2
}

func changeRField(circle *Circle, floatNumber float64) {
	circle.r = floatNumber
}

type Rectangle struct {
	x1, x2, y1, y2 float64
}

func (rectangle *Rectangle) area() float64 {
	return math.Sqrt(
		math.Pow(rectangle.x2-rectangle.x1, 2)) +
		math.Sqrt(math.Pow(rectangle.y2-rectangle.y1, 2))
}

func (rectangle *Rectangle) perimeter() float64 {
	return math.Abs(rectangle.x1-rectangle.x2) + math.Abs(rectangle.y2-rectangle.y1)
}

func circleDeclaration() {
	var firstCircle Circle      // Значение
	secondCircle := new(Circle) // Указатель

	fmt.Print("Параметры первой окружности после объявления: ")
	fmt.Println(firstCircle, firstCircle, firstCircle.y, firstCircle.r)
	fmt.Print("Параметры второй окружности после объявления: ")
	fmt.Println(secondCircle, secondCircle.x, secondCircle.y, secondCircle.r)
}

func circleInitialization() {
	firstCircle := Circle{x: 1.0, y: 2.0, r: 10.0} // Значение
	secondCircle := Circle{3.1, 2.1, 15.3}         // Значение
	thirdCircle := &Circle{1, 2, 3}                // Указатель

	fmt.Print("Параметры первой окружности после инициализации: ")
	fmt.Println(firstCircle, firstCircle.x, firstCircle.y, firstCircle.r)
	fmt.Print("Параметры второй окружности после инициализации: ")
	fmt.Println(secondCircle, secondCircle.x, secondCircle.y, secondCircle.r)
	fmt.Print("Параметры третьей окружности после инициализации: ")
	fmt.Println(thirdCircle, thirdCircle.x, thirdCircle.y, thirdCircle.r)
}

func circleArea(circle Circle) float64 {
	return math.Pi * math.Pow(circle.r, 2)
}

func usingCircleAreaFunctions() {

	firstCircle := Circle{4.15, 6.17, 10.0}
	fmt.Println("The circle area before any changes: ", firstCircle)
	fmt.Println("Area of the circle: ", circleArea(firstCircle))

	changeRField(&firstCircle, 21.1)
	fmt.Println("The circle fields after the X value change: ", firstCircle)
	fmt.Println("Area of the circle: ", circleArea(firstCircle))
}

func circleMethodUsing() {
	circle := Circle{1, 44.0, 15.12}
	fmt.Println(circle.area())
}

func getShapesAreaSum(shapes ...Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.area()
	}
	return total
}

func usingGetShapesAreaSum() {
	rectangle := Rectangle{1, 1, 2, 2}
	circle := Circle{1.0, 1.0, 15}

	var shapes []Shape
	shapes = append(shapes, &rectangle, &circle)

	fmt.Println("Sum of shapes areas is:", getShapesAreaSum(shapes...))
}

func modifiedUsingGetShapesAreaSum() {

	var rectangle Shape = &Rectangle{1, 1, 2, 2}
	circle := &Circle{1.0, 1.0, 200}

	var firstMultipleShape MultipleShape
	/* When Using a MultipleShape object !!! Slice of interface can hold only pointers to objects that
	implement this interface */
	firstMultipleShape.addShape(rectangle, circle)

	fmt.Println("Sum of shapes areas is:", firstMultipleShape.area())
}
