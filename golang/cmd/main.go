package main

import (
	"fmt"
	c3f "golang/pkg/chapter3/f_constants"
	"unsafe"
)

func PresenceOfUnusedSliceElements() {

	// Creating a slice that has the underlying array
	sl1 := make([]int, 0, 4)
	sl1 = append(sl1, 10, 20, 30, 40)

	// Creating another slice based on the first slice
	sl2 := sl1[1:3]

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

	fmt.Println()

	fmt.Println("sl2:")
	fmt.Printf("len: %d | cap: %d | slice: %v\n", len(sl2), cap(sl2), sl2)
	for i := 0; i < len(sl2); i++ {
		fmt.Printf("address: %x | value: %d\n", unsafe.Pointer(&sl2[i]), sl2[i])
	}
}

func RewritingElementsOfSlice() {
	slice1 := make([]int, 5, 8)
	slice2 := []int{1}
	fmt.Println(unsafe.SliceData(slice1), slice1, len(slice1), cap(slice1))
	fmt.Println(unsafe.SliceData(slice2), slice2, len(slice2), cap(slice2))
	fmt.Println()

	slice2 = append(slice1, 1, 1, 1)
	slice1[0] = 10
	fmt.Println(unsafe.SliceData(slice1), slice1, len(slice1), cap(slice1))
	fmt.Println(unsafe.SliceData(slice2), slice2, len(slice2), cap(slice2))
	fmt.Println()

	slice1 = append(slice1, 7)
	fmt.Println(unsafe.SliceData(slice1), slice1, len(slice1), cap(slice1))
	fmt.Println(unsafe.SliceData(slice2), slice2, len(slice2), cap(slice2))
	fmt.Println()
}

func main() {
	// c3bth.LocalServer()
	// c3c.Declaration_And_Printing()
	// c3cm.Server()
	// c3e.StringComparison()
	// c3e.StringsImmutability()
	// c3e.StringLiterals()
	// c3e.EscapeSequences()
	// c3e.DecodingUTF8()
	// c3e.LenAndRunes()
	// c3e.BasenameSimple("もちろん/もちろん/もちろん/もちろん.com")
	// c3e.BasenameLastIndex("もちろん/もちろん/もちろん/もちろん.com")
	// fmt.Println(c3e.IntegerComma("12345"))
	// c3e.BytesSliceMutability()
	// fmt.Println(c3e.IntsToString([]int{1, 2, 3}))
	// c3e.LenAndRunes()
	// c3e.IntegerCommaBuffer("5")
	// c3e.FloatCommaBuffer("1111.e10")
	// fmt.Println(c3e.IsAnagramTo("abcd", "dbca"))
	// c3e.IntegerToStringAndBack()
	// SliceTrick()
	// PresenceOfUnusedSliceElements()
	// c3f.UsingBigConstantsInExpressions()
	// c3f.DivisionWithDifferentLiterals()
	// c3f.ImplicitConversionToTargetType()
	// c3f.RoundingAfterConversion()
	// c3f.ImplicitVariableLiteralTypeDetermination()
	c3f.ExplicitVariableTypeDeclaration()
}
