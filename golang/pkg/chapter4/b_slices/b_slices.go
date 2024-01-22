package b_slices

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

/*
1. A slice is a lightweight data structure that gives access to a subsequence (or perhaps all) elements of an array,
   which is known as as the slice's underlying array.
2. A slice has three components: a pointer to an underlying array, a length, a capacity.
3. The pointer points to the first element of the array that is reachable through the slice, which is not necessarily
   the array's first element.
4. The length is the number of slice elements. It can't exceed the capacity.
5. The capacity is an amount of elements underlying array can store generally.
6. s[i,j) can be applied to: an array variable, a pointer to an array, another slice
7. Since a slice contains a pointer to an underlying array, passing a slice to a function permits the function to modi-
   fy the underlying array elements
8. Slices aren't comparable, but []byte slices can be compared. We can create our own comparisons, but IT'S FORBIDDEN.
	1) The elemets of slices are indirect (they can contain other types as slices and so on), making it possible for a slice
	to contain itself.
	2) We must recursively compare elements, because slices can contain other slices, strucures, interfaces as well.
	3) Because slice elements are indirect, a fixed slice value may contain different elements at different times as the
	contents of the underlying array are modified.
	DUE TO THE ABOVE POINTS WE MUST NOT COMPARE SLICES
9. The only thing we can compare slice with is nil.
10. The zero value of a slice is nil. When slice is a nil value it hasn't an underlying array and has length and capacity
0.
11. To find out whether a slice is empty we should call the len([]T) function.
12. Functions whose signature includes a slice parameter should contain a logic processing len([T] = 0) behaving.
13. Similar to the previous point functions should treat all zero-length slices the same way whether nil or non-nil
14. Slice can be created with make([]T, len, {capacity}) operator. If capacity is omitted, it's equal to the declared length.
	Under the hood make creates unnamed array variable and returns a slice of it.
15. Built-in append(dst, src) operator appends items to slices.
	1) If a new length is bigger than the available capacity, capacity of the result slice doubles up, elements from the
	source slice copies into the result slice and it's returned.
	2) After appending the amount of elements that is bigger that destiny capacity, the new destiny capacity will be
	equal to (old capacity + amount of added elements)
16. Build-in copy(dst, src) operator copies element from left to right rewriting elements within len(dst)
*/

func SliceOverlapping() {
	// The start 0 indice will be the empty string and won't be printed
	months := [...]string{0: "", 1: "January", 2: "February", 3: "March",
		4: "April", 5: "May", 6: "June",
		7: "July", 8: "August", 9: "September",
		10: "October", 11: "November", 12: "December",
	}

	// Capacity will be len(underlyingArray) - startIndex
	q2 := months[4:7]     // 4: "April", 5: "May", 6: "June", len = 3, cap = 9
	summer := months[6:9] // 6: "June", 7: "July", 8: "August" len = 3, cap = 7
	fmt.Printf("q2: % s\t| summer: % s\n", q2, summer)

	if strings.Contains(strings.Join(q2, ""), "June") == strings.Contains(strings.Join(summer, ""), "June") {
		fmt.Printf("June is containted in both slices\n")
	} else {
		fmt.Printf("June isn't containted in both slices\n")
	}

}

func Slicing() {
	// The start 0 indice will be the empty string and won't be printed
	months := [...]string{0: "", 1: "January", 2: "February", 3: "March",
		4: "April", 5: "May", 6: "June",
		7: "July", 8: "August", 9: "September",
		10: "October", 11: "November", 12: "December",
	}

	// Capacity will be len(underlyingArray) - startIndex
	q2 := months[4:7]     // 4: "April", 5: "May", 6: "June", len = 3, cap = 9
	summer := months[6:9] // 6: "June", 7: "July", 8: "August" len = 3, cap = 7
	fmt.Printf("q2: % s\t| summer: % s\n", q2, summer)

	/* Slicing beyond capacity causes panic error.*/
	// mt.Println(summer[:20])

	// Slicing up beyond len extends slice adding the next (newLen - oldLen) elements of the underlying array
	endlessSummer := summer[:5] // 6: "June", 7: "July", 8: "August", 9: "September", 10: "October"
	fmt.Printf("endlessSummer: % s\n", endlessSummer)
}

func Reverse() {
	array := [6]int{12, 15, 16, 1, 13, 24}
	slice := array[3:] // pointer -> [ 1, 13, 24 ]

	fmt.Printf("The array before reversing its slice: % d\n", array)
	swapSlice(slice)
	fmt.Printf("The array after reversing its slice: % d\n", array)

	fmt.Println()

	newArray := [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("The array before reversing its slice: % d\n", newArray)
	swapSlice(newArray[:]) // We must explicitly make a slice of the array
	fmt.Printf("The array after reversing its slice: % d\n", newArray)

}

func Rotation() {
	slice := []int{1, 2, 3, 4, 5}

	fmt.Printf("Array before 3 elements left rotation: %d\n", slice)
	easyRightRotation(3, slice)
	fmt.Printf("Array after 3 elements left rotation: %d\n", slice)
	easyLeftRotation(3, slice)
	fmt.Printf("Array after left 3 elements rotation: %d\n", slice)

}

func easyRightRotation(elementsNumber int, slice []int) {
	if slice == nil {
		fmt.Fprintf(os.Stderr, "The error occured: slice is nil\n")
		os.Exit(1)
	} else {
		swapSlice(slice[:elementsNumber])
		swapSlice(slice[elementsNumber:])
		swapSlice(slice)
	}
}

func easyLeftRotation(elementsNumber int, slice []int) {
	if slice == nil {
		fmt.Fprintf(os.Stderr, "The error occured: slice is nil\n")
		os.Exit(1)
	} else {
		swapSlice(slice)
		swapSlice(slice[:elementsNumber])
		swapSlice(slice[elementsNumber:])
	}
}

// It takes O(n)
func swapSlice(slice []int) {
	l := len(slice)
	for i := 0; i < len(slice)/2; i++ {
		slice[i], slice[l-1-i] = slice[l-1-i], slice[i]
	}
}

func SlicesComparisonInt(sl1, sl2 []int) bool {
	if len(sl1) == len(sl2) {
		return bytes.Equal(getBytes(sl1), getBytes(sl2))
	} else {
		return false
	}
}

func SlicesComparisonString(sl1, sl2 []string) bool {
	if len(sl1) == len(sl2) {
		return bytes.Equal([]byte(strings.Join(sl1, "")), []byte(strings.Join(sl2, "")))
	} else {
		return false
	}
}

func getBytes(slice []int) []byte {
	// Get an empty
	var buf bytes.Buffer

	// Run through the slice
	for _, val := range slice {
		// Write a byte sequence of a number into []byte
		buf.Write([]byte(strconv.Itoa(val)))
	}

	return buf.Bytes()
}

func NilSliceDeclaration() {
	// nil slice declarations
	var slice []int            // len(slice) == 0, s == nil
	slice = nil                // len(slice) == 0, s == nil
	slice = []int{}            // len(slice) == 0, s == nil
	slice = make([]int, 3)[3:] // len(slice) == 0, s == nil
	slice = []int(nil)         // len(slice) == 0, s == nil This statement makes nil slice as such as int(0)
	fmt.Printf("slice = %v, cap(slice) = %d, len(slice) = %d", slice, cap(slice), len(slice))
}
func NilSliceComparison() {
	slice := []int(nil) // len(slice) == 0, slice == nil

	if slice == nil {
		fmt.Fprintf(os.Stderr, "The error occured: slice is nil\n")
		os.Exit(1)
	} else {
		fmt.Println("Everything is a-okay")
	}
}

func NilSliceFunctionArgument() {
	easyLeftRotation(10, nil)
}

func MakeOperator() {
	slice1 := make([]int, 3, 10) // [0, 0, 0], len = 3, cap = 10
	slice2 := make([]string, 5)  // ["", "", "", "", ""], len = 5, cap = 5
	fmt.Printf("%v\n%v\n", slice1, slice2)
}

func Append() {
	str := "Привет, мир!"
	var runes []rune // [], len = 0,

	for _, value := range str {
		fmt.Printf("len: %3d | cap: %3d | underlying array: %v |\n", len(runes), cap(runes), runes)
		runes = append(runes, value)
	}
}

func AppendInt(initialSlice []int, elements ...int) []int {

	if initialSlice == nil {
		initialSlice = make([]int, 0, 1)
	}
	// Create a new slice
	var result []int

	// Set the new incremented length
	resultLength := len(initialSlice) + len(elements)

	// If there is a place in passed original slice
	if resultLength <= cap(initialSlice) {
		// Extend the original slice
		result = initialSlice[:resultLength]
	} else if len(elements) == 1 {
		resultCapacity := 2 * cap(initialSlice)
		result = (make([]int, resultLength, resultCapacity))
		copy(result, initialSlice)
	} else {
		resultCapacity := cap(initialSlice) + len(elements)
		result = make([]int, resultLength, resultCapacity)
		copy(result, initialSlice)
	}
	copy(result[len(initialSlice):], elements)
	return result
}

func Copy() {
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := make([]int, 4, 30) // 0 0 0 0
	copy(slice2, slice1)
	fmt.Println(len(slice2), cap(slice2), slice2) // 4, 30, [ 1 2 3 4 ]
}

func CapacityDoubling() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = append(x, i)
		fmt.Printf("add: %x | len: %3d | cap: %3d | underlying y array: %v\n", unsafe.SliceData(y), len(y), cap(y), y)
		fmt.Printf("add: %x | len: %3d | cap: %3d | underlying x array: %v\n\n", unsafe.SliceData(x), len(x), cap(x), x)
		x = y
	}
}

func MultipleAppend() {
	var x []int
	x = AppendInt(x, 1)
	fmt.Printf("add: %x | len: %3d | cap: %3d | underlying x array: %v\n\n", unsafe.SliceData(x), len(x), cap(x), x)
	x = AppendInt(x, 2, 3)
	fmt.Printf("add: %x | len: %3d | cap: %3d | underlying x array: %v\n\n", unsafe.SliceData(x), len(x), cap(x), x)
	x = AppendInt(x, 4, 5, 6)
	fmt.Printf("add: %x | len: %3d | cap: %3d | underlying x array: %v\n\n", unsafe.SliceData(x), len(x), cap(x), x)
	x = AppendInt(x, x...) // Slice adding
	fmt.Printf("add: %x | len: %3d | cap: %3d | underlying x array: %v\n\n", unsafe.SliceData(x), len(x), cap(x), x)
	x = AppendInt(x, 100)
	fmt.Printf("add: %x | len: %3d | cap: %3d | underlying x array: %v\n\n", unsafe.SliceData(x), len(x), cap(x), x)
}

func PresenceOfUnusedSliceElements() {
	// Creating a slice that has the underlying array
	sl1 := []int{10, 20, 30, 40}

	// Creating another slice based on the first slice
	sl2 := sl1[1:3] // 20, 30

	fmt.Println("BEFORE ADDING")
	fmt.Println("sl1:")
	fmt.Printf("len: %d | cap: %d | slice: %v\n", len(sl1), cap(sl1), sl1)
	for i := 0; i < len(sl1); i++ {
		fmt.Printf("address: %x | value: %d\n", unsafe.Pointer(&sl1[i]), sl1[i])
	}

	fmt.Println()

	fmt.Println("sl2:")
	fmt.Printf("len: %d | cap: %d | slice: %v\n", len(sl2), cap(sl2), sl2)
	for i := 0; i < len(sl2); i++ {
		fmt.Printf("address: %x | value: %d\n", unsafe.Pointer(&sl2[i]), sl2[i])
	}

	sl1 = append(sl1, 50)
	fmt.Println("\nAFTER ADDING")

	fmt.Println("sl1:")
	fmt.Printf("len: %d | cap: %d | slice: %v\n", len(sl1), cap(sl1), sl1)
	for i := 0; i < len(sl1); i++ {
		fmt.Printf("address: %x | value: %d\n", unsafe.Pointer(&sl1[i]), sl1[i])
	}

	fmt.Println("sl2:")
	fmt.Printf("len: %d | cap: %d | slice: %v\n", len(sl2), cap(sl2), sl2)
	for i := 0; i < len(sl2); i++ {
		fmt.Printf("address: %x | value: %d\n", unsafe.Pointer(&sl2[i]), sl2[i])
	}
}

func RewritingElementsOfSlice() {
	slice1 := make([]int, 5, 8) // { | 0, 0, 0, 0, 0 | { _, _, _ } }
	var slice2 []int            // {  }
	fmt.Println(unsafe.SliceData(slice1), slice1, len(slice1), cap(slice1))
	fmt.Println(unsafe.SliceData(slice2), slice2, len(slice2), cap(slice2))
	fmt.Println()

	slice2 = append(slice1, 1, 1, 1) // { | 0, 0, 0, 0, 0, 1, 1, 1 }
	slice1[0] = 10                   // sl1: { | 10, 0, 0, 0, 0 | { _, _, _ } }, sl2: { | 10, 0, 0, 0, 0, 1, 1, 1 }
	fmt.Println(unsafe.SliceData(slice1), slice1, len(slice1), cap(slice1))
	fmt.Println(unsafe.SliceData(slice2), slice2, len(slice2), cap(slice2))
	fmt.Println()

	slice1 = append(slice1, 7) // sl1: { | 10, 0, 0, 0, 0, 7 | { _, _ } }, sl2: { | 10, 0, 0, 0, 0, 7, 1, 1 }
	fmt.Println(unsafe.SliceData(slice1), slice1, len(slice1), cap(slice1))
	fmt.Println(unsafe.SliceData(slice2), slice2, len(slice2), cap(slice2))
	fmt.Println()
}

// After copying the second slice will have a new memory block different from the second
func CopyingSlices() {

	slice1 := []int{1, 2, 3}
	slice2 := make([]int, 3, 10)
	copy(slice2, slice1) //
	fmt.Println("Before changes:\n", slice1, "\n", slice2)

	slice2[0] = 10 // This change won't affect both slices
	fmt.Println("After changes:\n", slice1, "\n", slice2)
}

func NonEmptySimple(strings []string) []string {
	i := 0
	for _, str := range strings {
		if str != "" {
			strings[i] = str
			i++
		}
	}
	return strings[:i]
}

func NonEmptyPlumpy(strings []string) []string {
	s := strings[:0]
	i := 0
	for _, str := range strings {
		if str != "" {
			s = append(s, str)
			i++
		}
	}
	return s[:i]
}

func NonEmptyFat(strings []string) []string {
	var i int
	for ; i < len(strings); i++ {
		if strings[i] == "" {
			strings = append(strings[:i], strings[i+1:]...)
			i--
		}
	}
	return strings
}

type Stack struct {
	body []int
}

type StackIsEmpty interface {
	Error() error
}

type StackIsFilled interface {
	Error() error
}

func (stack *Stack) Init(values ...int) (int, StackIsFilled) {
	if len(stack.body) == 0 {
		stack.body = append(stack.body, values...)
		return len(stack.body), nil
	} else {
		var err StackIsFilled
		return 0, err
	}

}

func (stack *Stack) String() (string, StackIsEmpty) {
	if len(stack.body) != 0 {
		strBody := ""
		for _, value := range stack.body {
			strBody += string(rune(value))
		}
		return strBody, nil
	} else {
		var err StackIsEmpty
		return "", err
	}

}

func (stack *Stack) Push(value int) {
	if stack.body == nil {
		stack.body = make([]int, 0, 1)
	} else {
		stack.body = append(stack.body, value)
	}
}

func (stack *Stack) GetTop() (int, StackIsEmpty) {
	if len(stack.body) != 0 {
		return stack.body[len(stack.body)-1], nil
	} else {
		var err StackIsEmpty
		return 0, err
	}

}

func (stack *Stack) Pop() (int, StackIsEmpty) {
	if len(stack.body) != 0 {
		topValue := stack.body[len(stack.body)-1]
		stack.body = stack.body[:len(stack.body)-1]
		return topValue, nil
	} else {
		var err StackIsEmpty
		return 0, err
	}
}

func (stack *Stack) RemoveAtIndex(index int) (int, StackIsEmpty) {
	if len(stack.body) != 0 && index >= 0 && index < len(stack.body) {
		value := stack.body[index]
		if index != 0 {
			copy(stack.body[:index], stack.body[index+1:])
		} else {
			stack.body = stack.body[1:]
		}
		return value, nil
	} else {
		var err StackIsEmpty
		return 0, err
	}
}

// So that get an array value at a particular indices with a pointer we should use the following form: (*name)[index]
func swapArray(array *[]int) {
	l := len(*array)
	for i := 0; i < l/2; i++ {
		(*array)[i], (*array)[l-i-1] = (*array)[l-i-1], (*array)[i]
	}
}

func RotationLeftSingle(slice []int, n int) []int {
	l := len(slice)
	if l != 0 {
		if n > 0 && n != l {
			result := make([]int, 0)
			result = append(result, slice[n:]...)
			result = append(result, slice[:n]...)
			return result
		} else if n < 0 {
			return RotationRightSingle(slice, int(math.Abs(float64(n))))
		} else {
			return slice
		}
	} else {
		return nil
	}
}

func RotationRightSingle(slice []int, n int) []int {
	l := len(slice)
	if l != 0 {
		if n > 0 && n != l {
			n %= l
			result := make([]int, 0)
			result = append(result, slice[l-n:]...)
			result = append(result, slice[:l-n]...)
			return result
		} else if n < 0 {
			return RotationLeftSingle(slice, int(math.Abs(float64(n))))
		} else {
			return slice
		}

	} else {
		return nil
	}
}

func EliminateAdjacentDuplicates(strings []string) []string {
	if len(strings) > 1 {
		// Managing index of a slice
		index := 0
		// Running through the slice and matching adjacent strings
		for ; index < len(strings)-1; index++ {
			if bytes.Equal([]byte(strings[index]), []byte(strings[index+1])) {
				/*
					Reslice the original slice according to the following logic:
					Append all the data before the index of repeated element after this indice.
				*/
				strings = append(strings[:index], strings[index+1:]...)
				// Get started from the index you removed
				index--
			}
		}
	}
	return strings
}

func SquashSpaces(sequenceUTF8 []byte) []byte {
	left := 0
	right := 0

	// Running through a bytes array
	for right < len(sequenceUTF8) {

		outerRune, outerSize := utf8.DecodeRune(sequenceUTF8[right:])

		if !unicode.IsSpace(outerRune) {
			for i := 0; i < outerSize; i++ {

				sequenceUTF8[left] = sequenceUTF8[right]
				left++
				right++
			}
			continue
		}
		sequenceUTF8[left] = ' '
		left++

		for right < len(sequenceUTF8) {
			innerRune, innerSize := utf8.DecodeRune(sequenceUTF8[right:])
			if !unicode.IsSpace(innerRune) {
				break
			} else {
				right += innerSize
			}
		}

	}

	sequenceUTF8 = sequenceUTF8[:left]
	return sequenceUTF8
}

func SquashSpacesSimple(sequenceUTF8 []byte) []byte {
	// Get a new pointer that is out of memory of the underlying array
	sl := sequenceUTF8[:0]
	// Get a string representation of a bytes array
	str := string(sequenceUTF8)

	// Create the counter of the elements what will be added
	counter := 0
	// Create a flag that will be restrict repeated space symbols adding
	isSpace := false

	for _, run := range str {
		if !unicode.IsSpace(run) {
			// Get some rune bytes
			runeBytes := []byte(string(run))
			// Rewrite the bytes of the underlying array
			sl = append(sl, runeBytes...)
			// Increment the counter
			counter += len(runeBytes)

			// Reverse the flag
			isSpace = false
		} else if !isSpace {
			// Append the byte os the space
			sl = append(sl, ' ')
			// Increment the counter
			counter++
			// Reverse the flag
			isSpace = true
		}
	}

	// Cut up the result so that it won't have balderdash
	sl = sl[:counter]
	return sl
}

func ReverseUTF8(sequenceUTF8 []byte) []byte {
	runes := []rune(string(sequenceUTF8))
	return swapRunes(runes)
}

func swapRunes(runesSlice []rune) []byte {
	l := len(runesSlice)
	for i := 0; i < l/2; i++ {
		runesSlice[i], runesSlice[l-1-i] = runesSlice[l-1-i], runesSlice[i]
	}
	return []byte(string(runesSlice))
}
