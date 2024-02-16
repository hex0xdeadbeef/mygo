package d0_structs

import (
	"fmt"
	"math/rand"
	"time"
)

/*
1. A struct is an aggregate data type that groups ZERO or more named/unnamed values of arbitrary types within a single entity

2. Each value of a struct is called field

3. If the first letter of a field begins with an uppercase letter, the field can be exported, on the other hand if
the field begins with a lowercase letter it's not exported and accessible only within current package.

5. Each field can be accessed using dot notation. Due to this fact we can:
	1) Get a field value

	2) Change the value of a field

	3) Take a pointer of field

6. Every instance of a struct is just a part of memory with a particular size. Due to this fact we can make the
   following conclusions:
		1) When we get an OBJECT POINTER of a struct we get a pointer to the memory area.

		2) When we get an OBJECT FIELD POINTER we get the shifted pointer in the same memory area of this object.

	! FROM THE HOOD THIS PROCESSES TAKE PALCE IN THE COMPILER LEVEL WHILE IT TRANSLATES YOUR CODE TO MACHINE CODE. !

7. If we have a []StructType slice and put an object into it consequently passing the slice into a function that changes
a field of an object and returns a pointer to that object, we won't change the original object, because while appending
or putting with slice[index] the object into []StructType, the slice will have the copy of the object.
	! TO DEFEAT THIS WE SHOULD CREATE AN []*StructType SLICE TO STORE THE LINKS TO THE ORIGINAL OBJECTS. !

8. An aggregate type "S" can't contain the field of type "S", since recursive struct declaration is forbidden.
	! ON THE OTHER HAND IT CAN CONTAIN A FIELD WITH *S TYPE. THIS TRICK ALLOWS US TO CREATE THE DATA TYPES AS LINKED LIST
	AND TREES. !

9. The zero value of a struct is composed of each zero value of its fields.

10. There is the empty struct in Go. It may be useful when we're working with maps and need to emphasize only keys are
relevant

11. There are two flavors of struct literals:
	1) point := Point(1, 2), color.RGBA{red, green, blue, alpha}.
	This kind of literals is used when struct fields are obvious. We should follow the order of fields during initialization.

	2) anim := gif.GIF{LoopCount : nframes}.
	Literals like this are used in the situations when a struct is more sophisticated. In this case we need'n take care of
	the order of fields because names are provided.

	! IF AN FIELD ASSIGNMENT IS OMITTED, THIS FIELD WILL HAVE A CORRESPONDING TO IT ZERO VALUE !

	! THESE FORMS CAN'T BE COMBINED WITHIN A SINGLE LITERAL !

	! WE CAN'T SNEAK AROUND THE RULE THAT UNEXPORTED LITERALS CAN'T BE REFFERED TO FROM ANOTHER PACKAGE USING THE
	SECOND FORM OF A STRUCT LITERAL. UNEXPORTED FIELD WILL BE AVAILABLE TO WORKING WITH ONLY USING EMBEDDED METHODS. !

12. Struct values can be passed to functions and returned from them.
	1) For efficiency programmers pass large structures as a pointer.

	2) When struct object is passed to a function as a pointer, the rutime has a possibility of changing some object
	data.

	3) There are two short forms of declaring struct pointers:
		1) pointPointer := new(Point)
		! IN THIS CASE WE CAN USE THE METHODS OF Point STRUCT WITHOUT INITIALIZING !

		2) pointPointer := &Point{X:1, Y:2} this flavor allows us to use the pointer directly in an expression such
		as function call.

13. If all the fields of a struct are comparable, the struct itself is comparable with == and != operators.
	1) Due to this fact it means comparable structs may be used as keys of a map.
	! From the hood the == operator compares each field of both objects in the up down order. !

14. As the number of structs growth we should factor out their similarities and:
	1) Create new types that consider these features
	&
	2) Create inside the new structure an anonymous field of other structure (it's a field without name)
		1) An anonymous field should be either StructureName or *StructureName !

		2) Thanks to anonymous field we can refer to the names of fields/methods at the leaves of the implicit tree without
		giving the interveining names. For example: rightWheel.X (it's a filed of the senior type "Point")

		3)
		rightWheel := Wheel{Circle{Point{1, 2}, 10}, 15}
		OR
		rightWheel := Wheel{
			Circle: Circle{
				Point:  Point{X: 1, Y: 2},
				Radius: 10,
			},
			Spokes: 100,
		}
		Within a short initialization an object of type that is embedded with anonymous fields we should use the
		descending approach while writing out its values:
			1. Wheel has the "CIRCLE" with 15 spokes

			2. The CIRCLE has "POINT" and 10 radius

			3. The POINT has the coordinates 1, 2

		4) It's forbidden to declare two anonymous fields with the same types, because of names conflict and because the name of
		field is implicitly determined by its type so too is the visibility of field. For example:

			package "p"

			type point struct {
			X, Y int
			}

			"point" may be used as an anonymous type, but:
			1) its long access is allowed only within package level

			2) Out of package level the value of capitalized fields may accessed with shorthand:
				w.X = 100

		5) We may omit any or all the anonymous fields when selecting their subfields.
			rightWheel.X = 100

		6) An anonymous filed may represent a structure doesn't have any field, but contains a methods set. In effect,
		we have an outer object contains an inner object and outer object can use the methods of the inner object.
15. We can print all the fields names with all its values with adverb "#"
*/

/* The classic example of struct type */
/*
type Employee struct {
	ID        int
	Name      string // Can be combined
	Address   string // Can be combined
	DoB       time.Time
	Position  string // Can be combined
	Salary    int
	ManagerId int
}
*/

/* type Employee struct {
	ID                      int
	Name, Address, Position string
	DoB                     time.Time
	Salary                  int
	ManagerId               int
	employee                Employee // the compiler error: recursive declaring is forbidden
} */

type Employee struct {
	ID                      int
	Name, Address, Position string
	DoB                     time.Time
	Salary                  int
	ManagerId               int
}

type Person struct {
	Firstname, Lastname              string
	Age, PassportID, DriverLicenseID uint32
	Gender                           string
	partner                          *Person // This is the link to another person
}

type EmptyStruct struct{}

func CreateStructObject() {
	var Bob Employee
	fmt.Println(Bob)
}

func GetField() {
	var Bob Employee
	fmt.Printf("%d, %s, %s, %v, %s, %d, %d\n",
		Bob.ID,
		Bob.Name,
		Bob.Address,
		Bob.DoB,
		Bob.Position,
		Bob.Salary,
		Bob.ManagerId)
}

func ChangeFieldValue() {
	var Bob Employee
	fmt.Printf("Bob's salary before demoting: %d\n", Bob.Salary)
	Bob.Salary -= 5000
	fmt.Printf("Bob's salary after demoting: %d\n", Bob.Salary)
}

func GetFieldPointer() {
	var Bob Employee
	managerIdPointer := &Bob.ManagerId
	fmt.Printf("Bob's ManagerId field before changing it with a pointer: %d\n", Bob.ManagerId)
	ChangeFieldWithPointer(managerIdPointer)
	fmt.Printf("Bob's ManagerId field after changing it with a pointer: %d\n", Bob.ManagerId)
}
func ChangeFieldWithPointer(field *int) {
	*field = rand.Intn(10e5)
}

func GetStructObjectPointer() {
	var Bob Employee
	var bobPointer *Employee
	bobPointer = &Bob

	fmt.Printf("Bob's fields before changing them with a pointer: %v\n", Bob)
	changeEmployeeParameterWithPointer(bobPointer)
	fmt.Printf("Bob's fields before changing them with a pointer: %v\n", Bob)
}
func changeEmployeeParameterWithPointer(employee *Employee) {
	(*employee).Position = "Senior"
	employee.ManagerId = rand.Intn(10e5)
}

func FindObjectChangeField() {
	var Bob Employee
	Bob.Salary = 7000
	Bob.ManagerId = 10450100
	fmt.Printf("Bob's fields BEFORE changing them with a pointer: %v\n", Bob)

	var employees []*Employee
	employees = append(employees, &Bob)

	pointerToBob := getEmployeeByID(10450100, employees)
	pointerToBob.Salary = 100
	fmt.Printf("Bob's fields AFTER changing them with a pointer: %v\n", Bob)
}
func getEmployeeByID(requiredManagerID int, employees []*Employee) *Employee {
	if len(employees) != 0 {
		for index, _ := range employees {
			if employees[index].ManagerId == requiredManagerID {
				return employees[index]
			}
		}
	}
	return nil
}

type Tree struct {
	value       int
	left, right *Tree
}

func Sort(values []int) []int {
	var root *Tree
	for _, value := range values {
		root = add(root, value)
	}
	return appendInts(root, values[:0])
}
func add(root *Tree, value int) *Tree {
	if root == nil {
		root := &Tree{value: value, left: nil, right: nil}
		return root
	}

	if value < root.value {
		root.left = add(root.left, value)
	} else {
		root.right = add(root.right, value)
	}
	return root
}
func appendInts(root *Tree, values []int) []int {
	if root != nil {
		values = appendInts(root.left, values)
		values = append(values, root.value)
		values = appendInts(root.right, values)
	}
	return values
}

func MapWithEmptyStructValues() {
	emptyValuesMap := make(map[string]EmptyStruct)
	emptyValuesMap["Rustam"] = EmptyStruct{}
	emptyValuesMap["Dmitriy"] = EmptyStruct{}
	emptyValuesMap["Denis"] = EmptyStruct{}

	for key, _ := range emptyValuesMap {
		fmt.Printf(" %s", key)
	}
	fmt.Println()
}

/* Struct literals */

type Point struct {
	X, Y int
}

func StructLiteral() {
	point := Point{X: 0, Y: 0}
	fmt.Printf("%[1]T. %[1]v", point)
}

// This function will have a copy of the passed argument
func Scale(point Point, factor int) Point {
	point.X *= factor
	point.Y *= factor
	return point
}

func Bonus(employee *Employee, percent int) int {
	return employee.Salary*percent - employee.Salary
}

// This function will work with a copy of a pointer -> it can modify the field of an object
func AwardAnnualRaise(employee *Employee) {
	employee.Salary *= 105 / 100
}

/* Structs comparison */
func StructComparisonFat(point1, point2 Point) bool {
	return point1.X == point2.X && point1.Y == point1.Y
}

func StructComparisonLight(point1, point2 Point) bool {
	return point1 == point2
}

func ComparePoints(point1, point2 Point) bool {
	return StructComparisonFat(point1, point2) == StructComparisonLight(point1, point2)
}

type Address struct {
	hostname string
	port     int
}

func StructKeyMap() {
	address1 := Address{"google.com", 8080}
	address2 := Address{"microsoft.com", 8000}

	addressesMap := make(map[Address]bool)
	addressesMap[address1] = true
	addressesMap[address2] = false

	fmt.Println("\tAddresses status:")
	for address, isConnected := range addressesMap {
		fmt.Printf("add: %5v <-> %t\n", address, isConnected)
	}
}
