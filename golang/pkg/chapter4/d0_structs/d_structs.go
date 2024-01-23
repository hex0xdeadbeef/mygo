package d0_structs

import (
	"fmt"
	"math/rand"
	"time"
)

/*
1. A struct is an aggregate data type that groups ZERO or more named values of arbitrary types as a single entity

2. Each value of a struct is called field

3. Each fild can be accessed using dot notation. Due to this fact we can:
	1) Get a field value
	2) Change the value of a field
	3) Take a pointer of field

4. Every object of a struct is just a part of memory with a particular size. Due to this fact we can make the
   following conclusions:
		1) When we get an OBJECT POINTER of a struct we get a pointer to the memory area.
		2) When we get an OBJECT FIELD POINTER we get the a shifted pointer in the same memory area of this obejct.
	From the hood this processes take palce in the compiler level while it translates your code to machine code.

5. If we have a []StructType slice and put an object into it subsequently passing the slice into a function that re-
turns a pointer to StructType and change its field, we won't change the original object, because []StructType stotes
the copies of the initial object. To defeat this we should create an []*StructType slice to store the links to the
original objects.

6. If the first letter of a field begins with an uppercase letter, the field can be exported, on the other hand if
the field begins with a lowercase letter it's not exportable.

7. An aggregate type S can't contain the field of type S, since recursive struct declaration is forbidden.
	! On the other hand it can contain a field with *S type. This trick allows us to create the data types as linked
	list and trees. !

8. The zero value of a struct is composed of each zero value of its fields.
9. There is the empty struct in Go. It may be useful when we're working with maps and need to emphasize, that only
keys are relevant
10. There are two flavors of struct literals:
	1) point := Point(1, 2), color.RGBA{red, green, blue, alpha}. This kind of literals is used when struct fields
	are obvious. We should follow the order of fields.
	2) anim := gif.GIF{LoopCount : nframes}. Literals like this are used in the situations when a struct is more so-
	phisticated. In this case we need'n take care of the order of fields because names are provided.
	! If an field assignment is omitted, this field will have a corresponding to it zero value !
	! These forms can't be combined within a single literal !
	! We can't sneak around the rule that unexported literals can't be reffered to from another package using the
	second form of a struct literal. Unexported field will be available to working with only using embeded methods.
11. Struct values can be passed to functions and returned from them.
	1) For efficiency programmers pass large structures as a pointer.
	2) When struct object is passed to a function as a pointer, the rutine has a possibility of changing some object
	fields
	3) There are two short forms of declaring struct pointers:
		1) pointPointer := new(Point)
		2) *pointPointer = &Point{X:1, Y:2} this flavor allows us to use the pointer directly in an expression such
		as function call.

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
func ChangeFieldWithPointer(parameter *int) {
	*parameter = rand.Intn(10e5)
}

func GetStructObjectPointer() {
	var Bob Employee
	var bobPointer *Employee
	bobPointer = &Bob

	fmt.Printf("Bob's fields before changing them with a pointer: %v\n", Bob)
	changeEmployeesParametersWithPointer(bobPointer)
	fmt.Printf("Bob's fields before changing them with a pointer: %v\n", Bob)
}
func changeEmployeesParametersWithPointer(employee *Employee) {
	(*employee).Position = "Senior"
	employee.ManagerId = rand.Intn(10e5)
}

func FindObjectChangeItsField() {
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
			if (*employees[index]).ManagerId == requiredManagerID {
				return &(*employees[index])
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
		root = new(Tree)
		root.value = value
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
	return Point{X: point.X * factor, Y: point.Y * factor}
}

func Bonus(employee *Employee, percent int) int {
	return employee.Salary * percent
}

// This function will work with a copy of a pointer -> it can modify the field of an object
func AwardAnnualRaise(employee *Employee) {
	employee.Salary *= 105 / 100
}
