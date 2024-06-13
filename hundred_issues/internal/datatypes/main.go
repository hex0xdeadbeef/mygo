package main

/*
	DATA TYPES TRAPS

	OCTAL NUMBERS
1. The octal numbers are useful in different scenarios. For example: opening a file with os.OpenFile
2. We should use the explicit prefix of octal numbers: "0o"

	INTEGERS OVERFLOWS
1. If we exceed the top/low boundary of an integer type, the overflow happens. It results in balderdash in the answer.
2. If the number is signed and it's been overflowed, the sign changes in reverse.
3. Intergers' overflow that can be found during compilation results in the compile-error.
4. During execution overflow/underflow doesn't result in panic. It can result in unexpected errors.
5. We need control overflow/underflow.
6. So that control the overflow we should use these techniques:
	1) While incrementring we should compare the value being incremented with the corresponding math constant of the given type
	2) Working with the untyped integers we should use the corresponding vals.

	FLOATS PROBLEMS
1. The arithmetic of floating-point numbers is approximation of the real arithmetic.
2. The float32 keeps the accuracy up to 6 digits within the all number
3. The float64 keeps the accuracy up to 15 digits within the all number
4. Both float32 and float64 follow IEEE-754 convention. Some bits are mantissa and others is the power part.
	1) float32: 23 bits - mantissa, 8 bits - exponential part, 1 bit - sign
	2) float 64: 52 bits - mantissa, 11 bits - exponential part, 1 bit - sign
5. To represent a floating-point number into a decimal one the formula:
	sign * 2^exponent * mantissa
	is used
6. The outcomes of approximations instead of precise vals:
	1) The "==" operator cannot be used in comparisons -> we should use the epsilon value to find out whether the numbers are equal with the
	specific precision or not.
	2) The result of operation over two floating-point numbers depends on a specific architecture of the processor this operations performed upon.
	Many processors have the module responsible for computations over floating-point numbers (it's called FPU)
	3) There's no guarantee that the result of computation on the first computer will be the same as on the other one.
7. The errors can occure when:
	1) Casting types from integers to floats
	2) Sequential ops over floats
8. The order of ops applied. To get more precise result we apply the ops in the following way:
		1) Multiply | Divide
		2) Summarize | Subtraction
	The decomposed expression will work much more time, but will be more precise. There's a choice of precision/execution time.
9. When summarizing / subtraction we should group the elems of the close precision together to achieve bigger presicion
*/

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"unsafe"
)

func main() {

	// OctalLiterals()

	// IntegersOverflowA()

	// FloatingPointA()

	// FloatingPointB()

	// InfinitiesAndNan()

	// FloatingPointC()

	// FloatingPointD()

	// AppendElemIntoOneSourceArr()

	// EmptyNilSlices()

	// SliceCopyOps()

	// AppendUnwantedOutcomes()

	// FullSliceExpression()

	// SliceCapLeak()

	// SliceAndPointers()
}

func OctalLiterals() {
	fmt.Println(100 + 010) // 100 in decimal + 10 in octet system (0*8 + 1*8) -> the result is 108
	fmt.Println(100 + 0o010)
}

func OctalFileFlag() error {
	// 0644 - reading for all the users and writing only for the current user
	file, err := os.OpenFile("foo", os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	fmt.Println(file)
	return nil
}

func UnderscoreDelimiterUsage() {
	fmt.Println(1_000_000_000)
	fmt.Println(0b000_000_111)
	fmt.Println(0x123_435)
}

func IntegersOverflowA() {
	var (
		counter = math.MaxInt64
	)

	counter++

	// fmt.Println(counter) // -9223372036854775808 - the var has been overflowed
}

func IntegersOverflowB() {
	var (
	// i int64= math.MaxFloat64 + 1 //  cannot use math.MaxFloat64 + 1 (untyped float constant 1.79769e+308) as int64 value in variable
	// declaration (truncated)compilerTruncatedFloat
	)
}

func Increment32Control(counter int32) (int32, error) {
	if counter == math.MaxInt32 {
		return 0, fmt.Errorf("int32 overflow")
	}

	counter++
	return counter, nil
}

func IncrementUntypedIntControl(counter int) (int, error) {
	if counter == math.MaxInt {
		return 0, fmt.Errorf("untyped int overflow")
	}

	counter++
	return counter, nil
}

func IncrementUntypedUintControl(counter uint) (uint, error) {
	if counter == math.MaxUint {
		return 0, fmt.Errorf("untyped uint overflow")
	}

	counter++
	return counter, nil
}

func SumOverflowControl(a, b int) (int, error) {
	// The check won't work. Because the result of summarizing can be an overflowed val
	// if a+b > math.MaxInt {
	// 	return 0, fmt.Errorf("the target sum overflows untyped int")
	// }

	if a > math.MaxInt-b {
		return 0, fmt.Errorf("the target sum overflows untyped int")
	}

	return a + b, nil
}

func MultiplyingOverflowControl(a, b int) (int, error) {
	if a == 0 || b == 0 {
		return 0, nil
	}

	res := a * b
	if a == 1 || b == 1 {
		return res, nil
	}

	// If we multiply the math.MinInt val the product will overflow untyped int
	if a == math.MinInt || b == math.MinInt {
		return 0, fmt.Errorf("the target product overflows untyped int")
	}

	// If the equality isn't completed, the overflow has happened
	if res/a != b {
		return 0, fmt.Errorf("the target product overflows untyped int")
	}

	return 0, nil
}

func FloatingPointA() {
	var (
		num float32 = 1.0001
	)

	// Our expectations of getting the 1.00020001 result will fail
	fmt.Println(num * num)
}

func FloatingPointB() {
	fmt.Println(math.SmallestNonzeroFloat64) // 5e-324
	fmt.Println(math.MaxFloat64)             // 1.7976931348623157e+308
}

func InfinitiesAndNan() {
	var (
		zero        float64
		positiveInf float64 = 1 / zero
		negativeInf float64 = -1 / zero
		nan                 = zero / zero
	)

	fmt.Println(positiveInf)
	fmt.Println(negativeInf)
	fmt.Println(nan)

	fmt.Println(positiveInf == positiveInf, negativeInf == negativeInf, positiveInf > negativeInf, negativeInf < positiveInf)
	fmt.Println(math.IsNaN(nan))
	fmt.Println(math.IsInf(negativeInf, -1))
	fmt.Println(nan == nan, nan > nan, nan < nan)
}

func FloatingPointC() {
	var (
		errorGatheringA = func(n int) float64 {
			result := 10_000.

			for i := 0; i < n; i++ {
				result += 1.0001
			}

			return result
		}

		errorsGatheringB = func(n int) float64 {
			result := 0.

			for i := 0; i < n; i++ {
				result += 1.0001
			}

			return result + 10_000
		}
	)

	fmt.Println(errorGatheringA(10))  // 10010.000999999993
	fmt.Println(errorsGatheringB(10)) // 10010.001

	fmt.Println(errorGatheringA(1_000))  // 11000.099999999293
	fmt.Println(errorsGatheringB(1_000)) // 11000.099999999982
}

func FloatingPointD() {
	a := 100_000.001
	b := 1.0001
	c := 1.0002

	fmt.Println(a * (b + c)) // 200030.00200030004
	fmt.Println(a*b + a*c)   // 200030.0020003
}

func AppendElemIntoOneSourceArr() {
	s1 := make([]int, 3, 6)
	s1[1] = 1     // s1 = [0, 1, 0]
	s2 := s1[1:3] // s2 = [1, 0]

	s2 = append(s2, 2) // s2 = [1, 0, 2]

	fmt.Println(s1) // s1 = [0, 1, 0]
	fmt.Println(s2) // s2 = [1, 0, 2]

	s2 = append(s2, 3, 4, 5)
	p1, p2 := unsafe.Pointer(&s1), unsafe.Pointer(&s2)
	fmt.Println(p1 == p2)
	fmt.Println(s1)
	fmt.Println(s2)
}

type (
	CustomInt int
)

func ConvertSlow(ints []int) []CustomInt {
	customInts := make([]CustomInt, 0)

	for _, num := range ints {
		customInts = append(customInts, CustomInt(num))
	}

	return customInts
}

func ConverWithReuseLengthAppend(ints []int) []CustomInt {
	customInts := make([]CustomInt, 0, len(ints))

	for _, v := range ints {
		customInts = append(customInts, CustomInt(v))
	}

	return customInts
}
func ConverWithReuseLengthInit(ints []int) []CustomInt {
	customInts := make([]CustomInt, len(ints))

	for i, v := range ints {
		customInts[i] = CustomInt(v)
	}

	return customInts
}

func EmptyNilSlices() {
	var (
		nilA []int
		nilB = []int(nil)

		emptyA = []string{}
		emptyB = make([]string, 0)
	)

	fmt.Println(nilA == nil, nilB == nil)
	fmt.Println(emptyA == nil, emptyB == nil)
}

func EmptinessSliceTest[T any](s []T) any {
	if len(s) == 0 {
		return false
	}

	return s
}

func SliceCopyOps() {
	var (
		wrongA = func() {
			var (
				src = []int{1, 2, 3}
				dst []int
			)
			copy(dst, src)

			fmt.Println(dst)
		}

		rightA = func() {
			var (
				src = []int{1, 2, 3}
				dst = make([]int, len(src))
			)

			copy(dst, src)
			fmt.Println(dst)

			copy(dst[1:], src)
			fmt.Println(dst)
		}

		alternativeCopyA = func() {
			var (
				src = []int{1, 2, 3}
				dst = append([]int(nil), src...)
			)

			fmt.Println(dst)
		}
	)

	fmt.Println()
	wrongA()
	rightA()

	fmt.Println()
	alternativeCopyA()
}

func AppendUnwantedOutcomes() {
	var (
		a = func() {
			var (
				s1 = []int{1, 2, 3}
			)
			fmt.Println(s1, len(s1), cap(s1))

			s2 := s1[1:2]
			fmt.Println(s2, len(s2), cap(s2))

			fmt.Println()

			s3 := append(s2, 10)
			fmt.Println(s2, len(s2), cap(s2))
			fmt.Println(s3, len(s3), cap(s3))
		}

		f = func(s []int) {
			_ = append(s, rand.Intn(100))
		}

		fModified = func(s []int) []int {
			return append(s, rand.Intn(100))
		}

		changesVisibility = func() {
			var (
				s = []int{1, 2, 3}
			)

			f(s[:2])

			fmt.Println(s)
		}

		safeContentsChanging = func() {
			var (
				s     = []int{1, 2, 3}
				sCopy = make([]int, 2)
			)
			copy(sCopy, s)

			sCopy = fModified(sCopy)

			res := append(sCopy, s[2])
			fmt.Println(res)
		}
	)

	fmt.Println()
	a()

	fmt.Println()
	changesVisibility()

	fmt.Println()
	safeContentsChanging()
}

func FullSliceExpression() {
	var (
		s = []int{1, 2, 3}
		// s[low:high:max] and capacity is max-low
		// fullSliceExpr = s[:2:2] // -> the capacity of s is 2 - 0 = 2

		f = func(s []int) []int {
			s = append(s, rand.Intn(100))
			return s
		}
	)

	fmt.Println((f(s[:2])))
	fmt.Println((f(s[:2:2])))
}

func printAlloc() {
	const (
		mapToKB = 1024
	)

	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	fmt.Printf("Bytes in Heap: %-10d KB | Live objects: %-10d |  Freed: %-10d\n", m.Alloc/mapToKB, m.Mallocs, m.Frees)
}

func SliceCapLeak() {
	var (
		receiveMsg = func() []byte {
			const (
				msgSize int = 1_000_000
			)

			var (
				msg = make([]byte, 0, msgSize)
			)

			for range msgSize {
				msg = append(msg, byte(rand.Intn(128)))
			}

			return msg

		}

		storeMsgTypeCopy = func(storage [][]byte, msg []byte) [][]byte {
			const (
				typeRange = 5
			)

			var (
				msgType = make([]byte, typeRange)
			)
			copy(msgType, msg)

			// As a result we will store the source slices of the capacity 1_000_000 bytes.
			// return append(storage, msg[:typeRange])

			// In this case we don't store the link to the firs elem of msg, so the source slice won't be supported in future
			return append(storage, msgType)
		}

		_ = func(storage [][]byte, msg []byte) [][]byte {
			const (
				typeRange = 5
			)

			var (
				msgType = make([]byte, typeRange)
			)
			copy(msgType, msg)

			// As a result we will store the source slices of the capacity 1_000_000 bytes.
			// return append(storage, msg[:typeRange])

			// In this case we don't store the link to the firs elem of msg, so the source slice won't be supported in future
			return append(storage, msgType)
		}

		consumeMsgs = func() {
			const (
				msgsNumber = 5
			)

			var (
				msgTypeStorage = make([][]byte, 0, 1_000)
			)
			printAlloc()
			fmt.Println()

			for range msgsNumber {
				printAlloc()
				msgTypeStorage = storeMsgTypeCopy(msgTypeStorage, receiveMsg())
				// storeMsgTypeFullSliceExpr(msgTypeStorage, receiveMsg())
				printAlloc()

				fmt.Println()
			}

			for _, msgTypeBytes := range msgTypeStorage {
				fmt.Println(string(msgTypeBytes))
			}
		}
	)

	consumeMsgs()
}

/*
1) Allocates a slice of 1000 Storage elemes
2) Iterates over every Storage elem and allocates for Storage.body a 1 MB slice
3) Invokes keepFirstTwoElemsOnly that returns only the first two elems using reslicing and then calls runtime.GC
*/
func SliceAndPointers() {
	type (
		Storage struct {
			body []byte
		}
	)

	const (
		storagesCount           = 1_000
		fromBytesToMBMultiplier = 1 << 20
	)

	var (
		storages = make([]Storage, storagesCount)

		keepFirstTwoElemsOnly = func(storages []Storage) []Storage {
			// In this case the block of memory initially allocated for 1_000 storages is stored even after reslicing
			/*
				Before storages allocation
				Bytes in Heap: 138        KB | Live objects: 168        |  Freed: 5

				After storages allocation
				Bytes in Heap: 1024148    KB | Live objects: 1206       |  Freed: 29

				After reslicing and GC invocation
				Bytes in Heap: 1024150    KB | Live objects: 1212       |  Freed: 32
			*/
			// return storages[:2]

			// In this case we create a new pointer, so the excess of memory won't be stored after reassigning
			/*
				Before storages allocation
				Bytes in Heap: 138        KB | Live objects: 167        |  Freed: 5

				After storages allocation
				Bytes in Heap: 1024151    KB | Live objects: 1210       |  Freed: 29

				After reslicing and GC invocation
				Bytes in Heap: 2201       KB | Live objects: 1219       |  Freed: 1030
			*/
			// var (
			// 	storagesCopy = make([]Storage, 2)
			// )
			// copy(storagesCopy, storages)
			// return storagesCopy

			// In this case we want to save the initial capacity of storages and nil the rest of the storages slice
			/*
				Before storages allocation
				Bytes in Heap: 138        KB | Live objects: 168        |  Freed: 5

				After storages allocation
				Bytes in Heap: 1024151    KB | Live objects: 1209       |  Freed: 29

				After reslicing and GC invocation
				Bytes in Heap: 2201       KB | Live objects: 1218       |  Freed: 1030
			*/
			for i := 2; i < len(storages); i++ {
				storages[i].body = nil
			}

			return storages
		}
	)

	fmt.Println("Before storages allocation")
	printAlloc()
	fmt.Println()
	for i := 0; i < storagesCount; i++ {
		storages[i].body = make([]byte, fromBytesToMBMultiplier)

	}
	fmt.Println("After storages allocation")
	printAlloc()
	fmt.Println()

	// The block of memory initially allocated for 1_000 storages is stored even after reslicing
	firstTwoStorages := keepFirstTwoElemsOnly(storages)

	runtime.GC()
	fmt.Println("After reslicing and GC invocation")
	printAlloc()
	fmt.Println()

	runtime.KeepAlive(firstTwoStorages)
}
