package generics

import (
	"fmt"
	"math"
	"strings"
)

/*go.dev ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

func SumInts(vals map[string]int64) int64 {
	var total int64

	for _, v := range vals {
		total += v
	}

	return total
}

func SumFloats(vals map[string]float64) float64 {
	var total float64

	for _, v := range vals {
		total += v
	}

	return total
}

func SumItsOrFloats[K comparable, V int64 | float64](vals map[K]V) V {
	var total V
	for _, v := range vals {
		total += v
	}

	return total
}

type Numeric interface {
	int64 | float64
}

func SumItsOrFloatsConstrained[K comparable, V Numeric](vals map[K]V) V {
	var total V
	for _, v := range vals {
		total += v
	}

	return total
}

func SimpleUsing() {

	ints := map[string]int64{
		"first":  100,
		"second": 200,
	}

	floats := map[string]float64{
		"first":  100.0,
		"second": 200.0,
	}

	fmt.Printf("Non-generic sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	fmt.Printf("Generic sums: %v and %v\n", SumItsOrFloats(ints), SumItsOrFloats(floats))

	fmt.Printf("Generic with constraint sums: %v and %v\n", SumItsOrFloats(ints), SumItsOrFloatsConstrained(floats))

}

/*Leetware guide ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/* BASICS -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

// Has takes "T" type list and uses the "type inference" mechanism to find out what type it works with.
// Has returns true if the element key is present in the list, otherwise returns false
func Has[T comparable](list []T, key T) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}

	return false
}

func HasUsing() {
	fmt.Printf("List has %q is %t\n", "a", Has([]string{"a", "b", "c"}, "a"))
	fmt.Printf("List has %q is %t\n", "d", Has([]string{"a", "b", "c"}, "d"))
	fmt.Printf("List has \"%d\" is %t\n", 5, Has([]int{1, 2, 3}, 5))
	fmt.Printf("List has \"%d\" is %t\n", 1, Has([]int{1, 2, 3}, 1))
}

// NewEmptyList talkes no parameters of the "T" type and returns the list of elements of type "T", but due to the fact, that we pass no parameters of the specific type, we should explicitly
// point the compiler what type it must use to create and return the list of type "T"
func NewEmptyList[T any]() []T {
	return make([]T, 33)
}

func NewEmptyListUsing() {
	// fmt.Printf("The empty list of %q type list: %#v.", "string", NewEmptyList[]()) // expected operand, found ']'
	fmt.Printf("The empty list of %q type list: %#v.", "int64", NewEmptyList[int64]())

}

// The A, B, C are the type constraints that tell the compiler what types are permitted to be pass into the PrintThings as arguments.
// NOTE: When passing arguments for a1 and a2 we should pass the values of the same type.
// NOTE: A value for c must be the type of "int" or the hand-made type that has the underlying "int" type
func PrintThings[A, B any, C ~int](a1, a2 A, b B, c C) {
	fmt.Printf("%v %v %v %v", a1, a2, b, c)
}

func PrintThingsUsing() {
	// PrintThings(1, "A", 2, 10) // mismatched types untyped int and untyped string (cannot infer A)
	// PrintThings(1, 2, 2, 10.3) // float64 does not satisfy ~int (float64 missing in ~int)
	PrintThings(10, 20, "A string", 1000)
	// PrintThings(10, 20, "A string", int64(1000)) // C (type int64) does not satisfy ~int
	// PrintThings(10, 20, "A string", int32(1000)) // C (type int32) does not satisfy ~int
	PrintThings(10, 20, "A string", int(int64(1000)))
	var f float64 = 10.345
	PrintThings(10, 20, true, int(f))

	type MyInteger int
	mi := MyInteger(333)
	PrintThings(10, 20, byte(127), mi)

}

type Set[C comparable] map[C]struct{}

/*
NOTE: In order to bind the method to the receiver of type "K", we also need to define the generic type in the left declaration as "Set[K]"
*/
func (s Set[C]) Has(key C) bool {
	_, ok := s[key]
	return ok
}

/*
NOTE: to define a function that will work with a constraining types we should:
 1. Declare the constraining set of types we will use in the function body
 2. Define the parameters that will be of a type in the constraining set
 3. Optionally: provide the function signature with return values of a type in the constraining type set.
*/
func NewSet[C comparable](values ...C) Set[C] {
	set := make(map[C]struct{}, len(values))

	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}

func NewSetUsing() {
	fmt.Printf("The type of %q is %v\n", "int", NewSet([]int{1, 2, 3, 4, 5, 5, 6, 6, 7, 7, 7, 10, -1}...))
	fmt.Printf("The type of %q is %v\n", "string", NewSet("apple", "car", "desk"))
	fmt.Printf("The type of %q is %v\n", "bool", NewSet(true, true, true, false, true, false))
}

func SetHasUsing() {

	intSet := NewSet(10, -1, 100, 500, 1000, 200)

	fmt.Printf("Set: %v has a value %v is: %t\n", intSet, 300, intSet.Has(300))
	fmt.Printf("Set: %v has a value %v is: %t\n", intSet, 1000, intSet.Has(1000))

	fmt.Println()

	strSet := NewSet("apple", "one", "wolf", "desk", "phone")
	fmt.Printf("Set: %v has a value %v is: %t\n", strSet, "wolf", strSet.Has("green"))
	fmt.Printf("Set: %v has a value %v is: %t\n", strSet, "wolf", strSet.Has("wolf"))
}

/* CONSTRAINTS-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

// [T any]
func HasWithAnyConstraint[T any](key T, isEqual func(a, b T) bool, values ...T) bool {
	for _, v := range values {
		if isEqual(key, v) {
			return true
		}
	}

	return false
}

// NOTE: So that we can work with the constraint "any" we should provide the "HasWithAnyConstraint" function with the subroutine that will define the behavior of working with the values of
// type "T"
func HasWithAnyConstraintUsing() {
	fmt.Println(HasWithAnyConstraint(
		10,
		func(a, b int) bool {
			return a == b
		},
		-1, 1, 2, 3, 10, 4, 5,
	))

	fmt.Println(HasWithAnyConstraint(
		"apple",
		func(a, b string) bool {
			const prefix = "app"
			return strings.HasPrefix(a, prefix) && strings.HasPrefix(b, prefix)
		},
		"original", "airplane", "intilligent", "desk", "tablet",
	))

}

// [C comparable]
func HasWithComparableConstraint[C comparable](key C, values ...C) bool {
	for _, v := range values {
		if key == v {
			return true
		}
	}

	return false
}

// NOTE: the argumets we pass into the function with the [C comparable] contraint must implement the logic of working with the "==" and "!=" operators, so the types of arguments should be the
// one of the following one: numeric, "string", struct with all the comparable fields.
func HasWithComparableConstraintUsing() {

	fmt.Println(HasWithComparableConstraint(-15, 10, -15, 23, 33, -15))
	fmt.Println(HasWithComparableConstraint("original", "airplane", "intilligent", "desk", "tablet"))

	type Pair struct {
		name string
		id   int
	}

	fmt.Println(HasWithComparableConstraint(Pair{"Dmitriy", 1}, Pair{"Oleg", 3}, Pair{"Dmitriy", 1}, Pair{"George", 1}))
}

// [M MyInterface]
type Equalizer[T any] interface {
	IsEqual(other T) bool
}

type Passport struct {
	Fullname string
	ID       int
}

func (pi *Passport) IsEqual(other *Passport) bool {
	return *pi == *other
}

// NOTE: In this case when we put the interface name into the constraint set we should explicitly write out the preceding type parameter we declare in advance after the interface name.
func HasWithInterfaceConstraint[E Equalizer[E]](key E, values ...E) bool {
	for _, v := range values {
		if key.IsEqual(v) {
			return true
		}
	}

	return false
}

func HasWithInterfaceConstraintUsing() {
	passport1 := &Passport{"Dmitriy Mamykin", 333}
	passport2 := &Passport{"George Odin", 777}
	passport3 := &Passport{"Alexandr Zherenovskiy", 111}
	passport4 := &Passport{"Dmitriy Mamykin", 333}

	fmt.Println(HasWithInterfaceConstraint(passport1, passport2, passport3, passport4))
}

// [U ~underlyingType]

type ProductName string

func HasWithUnderlyingTypeConstraint[U ~string](key U, values ...U) bool {
	for _, v := range values {
		if key == v {
			return true
		}
	}

	return false
}

func HasWithUnderlyingTypeConstraintUsing() {
	pn1 := ProductName("Superb desk")
	pn2 := ProductName("The best kettle")
	pn3 := ProductName("Squared spoon")
	pn4 := ProductName("Soft refrigerator")
	pn5 := ProductName("Warm bread")

	fmt.Println(HasWithUnderlyingTypeConstraint(ProductName("Warm bread"), pn1, pn2, pn3, pn4, pn5))
}

// [U type1|...|~typeN]

// NOTE: In this case we must check whether all the types support the operations we  want to ptocess in the function body. If their do we won't get the compiler error, otherwise we will.

/*
func Compare[T int | bool](a, b T) int {
	if a == b {
		return 0
	}

	if a < b { // invalid operation: a < b (type parameter T is not comparable with <)
		return -1
	}

	return 1
}
*/

type SimpleNumeric interface {
	~int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | ~float64
}

func HasWithTypeUnionConstraint[SN SimpleNumeric](a, b SN) int {
	if a == b {
		return 0
	}

	if a < b {
		return -1
	}

	return 1
}

func HasWithTypeUnionConstraintUsing() {
	fmt.Println(HasWithTypeUnionConstraint(-1.0, 1.0))
}

/* COMBINING CONSTRAINTS -------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

type Currency interface {
	~int | ~int64
	ISO4127Code() string
	Decimal() int
}

type USD int64

func (usd USD) ISO4127Code() string {
	return "$USD"
}

func (b USD) Decimal() int {
	return 2
}

func PrintBalance[T Currency](b T) {
	balance := float64(b) / math.Pow10(b.Decimal())
	fmt.Printf("%.*f %s\n", b.Decimal(), balance, b.ISO4127Code())
}

func PrintBalanceUsing() {
	PrintBalance(USD(333))
}

/* TYPE PARAMETERS PROPOSAL -----------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

func Print2[T1, T2 any](s1 []T1, s2 []T2) {
	for _, v := range s1 {
		fmt.Println(v)
	}

	for _, v := range s2 {
		fmt.Println(v)
	}
}

func Print2Same[T any](s1, s2 []T) {
	for _, v := range s1 {
		fmt.Println(v)
	}

	for _, v := range s2 {
		fmt.Println(v)
	}
}

// Generic type
type Vector[T any] []T

func (v *Vector[T]) Push(val T) {
	*v = append(*v, val)
}

var v Vector[int]

// Referencing itself must be done with the same order of type parameters
type Node[T any] struct {
	next *Node[T]
	val  T
}

type P[T1, T2 any] struct {
	next *P[T2, T1] // INVALID; MUST BE [T1, T2]
}

// This restriction applies to both direct and indirect references
type ListHead[T any] struct {
	head *ListElement[T]
}

type ListElement[T any] struct {
	next *ListElement[T]
	val  T
	head *ListHead[T]
}

// The type parameter of a generic type may have constraints other than "any"
type StringableVector[T fmt.Stringer] []T

func (s StringableVector[T]) String() string {
	var sb strings.Builder

	for i, v := range s {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(v.String())
	}

	return sb.String()
}
