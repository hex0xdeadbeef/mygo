package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"unsafe"
)

/*
	ATOMICS
1. Package atomics provides low-level atomic memory primitives useful for implementing synchronizations algorithms.
2. Functions of atomics pkg require great care to be used correctly. Except for special, low-level apps, syncronization is better done with
channels or the facilities of the sync pkg.
3. The Swap operation, implemented by SwapT functions is the atomic equivalent of:
	old = *addr
	*addr = newValue
	return old
4. The compare-and-swap operation, implemendted by CompareAndSwapT functions, is the atomic equivalent of:
	if *addr == old {
		*addr = newValue
		return true
	}
	return false
5. The add operation, implemented by the AddT functions, is the atomic equivalent of:
	*addr += delta
	return *addr
6. The load and store operations, implemented by the LoadT and StoreT functions, are the atomic equivalent of:
	1) return *addr
	2) *addr = value
7. In the terms of Go memory model, if the effect of op A is observed by the atomic operation B, then A < B. Additionally, all the atomic operations executed in a program
behave as though executed in some sequentially consistent order.
8. To subtract positive integer value c from x we should use AddUint(32/64)(x, ^unintx(32/64)(c - 1)).
	1) In particular, to decrement value we use: AddUint(32/64)(x, ^uint(32/64)(0))
9. CompareAndSwapT. We check whether the old value is that we expected and if it's true, we swap the value of this addr to new value.
10. atomics must not be copied. We should use it as a pointers or we should store it in a structure with pointer methods.
11. A type Value provides an atomic load and store of a consistently typed value.
	1) The zero value for a Value returns nil from Load() method
	2) Once Store() method has been called, a Value must not be copied.
	3) A Value must not be copied after first use.
12. CAS ABA Problem:
	In multithreaded computing, the ABA proble occurs during synchronization, when a location is read twice, has the same value for both reads, and the read value being the same
	twice is used to conclude that nothing has happened in the interim. However, another thread can execute between the two reads and change the value, do other work, then change
	the value back, thus fooling the first thread into thinking nothing has changed even though the second thread did work that violates that assumption.

The ABA problem occurs when multiple threads accessing shared data interleave. Below is a sequence of events that illustrates the ABA problem:
	1. M1 reads the value A, allowing thread M2 to run
	2. M2 writes value B to the shared memory location
	3. M2 does some dork
	4. M2 writes value A to the shared memory location
	5. M2 is preempted, allowing thread M1 to run
	6. M1 reads the value A from the shared memory location
	7. M1 determines that the shared memory location hasn't been changed and continues
The M1 can fail because of shaded updates/changes made by M2 and the work of an app can be destructed.
*/

func main() {

	ValueUsage()
}

/*
INT32S
*/
func AddInt32Usage() {
	var (
		val int32
		wg  sync.WaitGroup
	)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				// AddInt32 atomically adds delta to *addr and returns the new value.
				// Consider using the more ergonomic and less error-prone [Int32.Add] instead.
				fmt.Println(atomic.AddInt32(&val, 1))
			}
		}()
	}

	wg.Wait()

	fmt.Println("Result:", val)
}

func CompareAndSwapInt32Usage() {
	var (
		val            int32 = 1
		oldVal, newVal int32 = 10, 20
		wg             sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		// CompareAndSwapInt32 executes the compare-and-swap operation for an int32 value.
		// Consider using the more ergonomic and less error-prone [Int32.CompareAndSwap] instead.
		fmt.Printf("old: %d | new: %d | res: %t\n", val, oldVal, atomic.CompareAndSwapInt32(&val, val, oldVal))
	}()

	go func() {
		defer wg.Done()
		// CompareAndSwapInt32 executes the compare-and-swap operation for an int32 value.
		// Consider using the more ergonomic and less error-prone [Int32.CompareAndSwap] instead.
		fmt.Printf("old: %d | new: %d | res: %t\n", oldVal, newVal, atomic.CompareAndSwapInt32(&val, oldVal, newVal))
	}()

	wg.Wait()
}

func LoadInt32Usage() {
	var (
		val int32 = 333
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// LoadInt32 atomically loads *addr.
				// Consider using the more ergonomic and less error-prone [Int32.Load] instead.
				fmt.Println(atomic.LoadInt32(&val))

				if rand.Intn(100)%2 == 0 {
					// AddInt32 atomically adds delta to *addr and returns the new value.
					// Consider using the more ergonomic and less error-prone [Int32.Add] instead.
					atomic.AddInt32(&val, 1)
				}
			}

		}()
	}

	wg.Wait()
}

func StoreInt32Usage() {
	var (
		val int32
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 3; i++ {
				// StoreInt32 atomically stores val into *addr.
				// Consider using the more ergonomic and less error-prone [Int32.Store] instead.
				atomic.StoreInt32(&val, rand.Int31n(128))

				// LoadInt32 atomically loads *addr. Consider using the more ergonomic and less error-prone [Int32.Load] instead.
				fmt.Println(atomic.LoadInt32(&val))
			}
		}()
	}

	wg.Wait()
}

func SwapInt32Usage() {
	var (
		s = make([]int32, 10)

		wg sync.WaitGroup
	)

	for i := 0; i < len(s); i++ {
		s[i] = rand.Int31n(128)
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {

				res := fmt.Sprintln("Before:", s)

				ptr := &s[rand.Intn(len(s))]
				newVal := rand.Intn(128)

				// SwapInt32 atomically stores new into *addr and returns the previous *addr value. Consider using the more ergonomic and less error-prone [Int32.Swap] instead.
				old := atomic.SwapInt32(ptr, int32(newVal))
				res += fmt.Sprintln("Old Val:", old)

				res += fmt.Sprintln("After:", s)

				fmt.Println(res)
			}
		}()
	}

	wg.Wait()
}

/*
INT64S
*/
func AddInt64Usage() {
	var (
		val int64
		wg  sync.WaitGroup
	)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				// AddInt64 atomically adds delta to *addr and returns the new value.
				// Consider using the more ergonomic and less error-prone [Int64.Add] instead (particularly
				// if you target 32-bit platforms; see the bugs section).
				fmt.Println(atomic.AddInt64(&val, 1))
			}
		}()
	}

	wg.Wait()

	fmt.Println("Result:", val)
}

func CompareAndSwapInt64Usage() {
	var (
		val            int64 = 1
		oldVal, newVal int64 = 10, 20
		wg             sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		// CompareAndSwapInt64 executes the compare-and-swap operation for an int64 value.
		// Consider using the more ergonomic and less error-prone [Int64.CompareAndSwap] instead (particularly if you target 32-bit platforms; see the bugs section).
		fmt.Printf("old: %d | new: %d | res: %t\n", val, oldVal, atomic.CompareAndSwapInt64(&val, val, oldVal))
	}()

	go func() {
		defer wg.Done()
		// CompareAndSwapInt64 executes the compare-and-swap operation for an int64 value.
		// Consider using the more ergonomic and less error-prone [Int64.CompareAndSwap] instead (particularly if you target 32-bit platforms; see the bugs section).
		fmt.Printf("old: %d | new: %d | res: %t\n", oldVal, newVal, atomic.CompareAndSwapInt64(&val, oldVal, newVal))
	}()

	wg.Wait()
}

func LoadInt64Usage() {
	var (
		val int64 = 777
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// LoadInt64 atomically loads *addr.
				// Consider using the more ergonomic and less error-prone [Int64.Load] instead (particularly if you target 32-bit platforms; see the bugs section).
				fmt.Println(atomic.LoadInt64(&val))

				if rand.Intn(100)%2 == 0 {
					// AddInt64 atomically adds delta to *addr and returns the new value.
					// Consider using the more ergonomic and less error-prone [Int64.Add] instead (particularly if you target 32-bit platforms; see the bugs section).
					atomic.AddInt64(&val, 1)
				}
			}

		}()
	}

	wg.Wait()
}

func StoreInt64Usage() {
	var (
		val int64
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 3; i++ {
				// StoreInt64 atomically stores val into *addr.
				// Consider using the more ergonomic and less error-prone [Int64.Store] instead (particularly if you target 32-bit platforms; see the bugs section).
				atomic.StoreInt64(&val, rand.Int63n(math.MaxInt64))

				// LoadInt64 atomically loads *addr.
				// Consider using the more ergonomic and less error-prone [Int64.Load] instead (particularly if you target 32-bit platforms; see the bugs section).
				fmt.Println(atomic.LoadInt64(&val))
			}
		}()
	}

	wg.Wait()
}

func SwapInt64Usage() {
	var (
		s = make([]int64, 10)

		wg sync.WaitGroup
	)

	for i := 0; i < len(s); i++ {
		s[i] = rand.Int63n(128)
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {

				res := fmt.Sprintln("Before:", s)

				ptr := &s[rand.Intn(len(s))]
				newVal := rand.Intn(128)

				// SwapInt64 atomically stores new into *addr and returns the previous *addr value.
				// Consider using the more ergonomic and less error-prone [Int64.Swap] instead (particularly if you target 32-bit platforms; see the bugs section).
				old := atomic.SwapInt64(ptr, int64(newVal))
				res += fmt.Sprintln("Old Val:", old)

				res += fmt.Sprintln("After:", s)

				fmt.Println(res)
			}
		}()
	}

	wg.Wait()
}

/*
UINT32S
*/
func AddUint32Usage() {
	var (
		val uint32
		wg  sync.WaitGroup
	)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				// AddUint32 atomically adds delta to *addr and returns the new value.
				// To subtract a signed positive constant value c from x, do AddUint32(&x, ^uint32(c-1)).
				// In particular, to decrement x, do AddUint32(&x, ^uint32(0)). Consider using the more ergonomic and less error-prone [Uint32.Add] instead.
				fmt.Println(atomic.AddUint32(&val, 1))
			}
		}()
	}

	wg.Wait()

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				// To subtract a signed positive constant value c from x, do AddUint32(&x, ^uint32(c-1)).
				// In particular, to decrement x, do AddUint32(&x, ^uint32(0)). Consider using the more ergonomic and less error-prone [Uint32.Add] instead.
				fmt.Println(atomic.AddUint32(&val, ^uint32(0)))
			}
		}()
	}

	wg.Wait()

	fmt.Println("Result:", val)

}

func CompareAndSwapUint32Usage() {
	var (
		val            uint32 = 1
		oldVal, newVal uint32 = 10, 20
		wg             sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		// CompareAndSwapUint32 executes the compare-and-swap operation for a uint32 value.
		// Consider using the more ergonomic and less error-prone [Uint32.CompareAndSwap] instead.
		fmt.Printf("old: %d | new: %d | res: %t\n", val, oldVal, atomic.CompareAndSwapUint32(&val, val, oldVal))
	}()

	go func() {
		defer wg.Done()
		// CompareAndSwapUint32 executes the compare-and-swap operation for a uint32 value.
		// Consider using the more ergonomic and less error-prone [Uint32.CompareAndSwap] instead.
		fmt.Printf("old: %d | new: %d | res: %t\n", oldVal, newVal, atomic.CompareAndSwapUint32(&val, oldVal, newVal))
	}()

	wg.Wait()

}

func LoadUint32Usage() {
	var (
		val uint32 = 333
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// LoadUint32 atomically loads *addr.
				// Consider using the more ergonomic and less error-prone [Uint32.Load] instead.
				fmt.Println(atomic.LoadUint32(&val))

				if rand.Intn(100)%2 == 0 {
					// AddUint32 atomically adds delta to *addr and returns the new value. To subtract a signed positive constant value c from x, do AddUint32(&x, ^uint32(c-1)).
					// In particular, to decrement x, do AddUint32(&x, ^uint32(0)). Consider using the more ergonomic and less error-prone [Uint32.Add] instead.
					atomic.AddUint32(&val, ^uint32(0))
				}
			}

		}()
	}

	wg.Wait()
}

func StoreUint32Usage() {
	var (
		val uint32
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 3; i++ {
				// StoreUint32 atomically stores val into *addr.
				// Consider using the more ergonomic and less error-prone [Uint32.Store] instead.
				atomic.StoreUint32(&val, rand.Uint32())

				// LoadUint32 atomically loads *addr. Consider using the more ergonomic and less error-prone [Uint32.Load] instead.
				fmt.Println(atomic.LoadUint32(&val))
			}
		}()
	}

	wg.Wait()
}

func SwapUint32Usage() {
	var (
		s = make([]uint32, 10)

		wg sync.WaitGroup
	)

	for i := 0; i < len(s); i++ {
		s[i] = rand.Uint32()
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {

				res := fmt.Sprintln("Before:", s)

				ptr := &s[rand.Intn(len(s))]
				newVal := rand.Uint32()

				// SwapUint32 atomically stores new into *addr and returns the previous *addr value. Consider using the more ergonomic and less error-prone [Uint32.Swap] instead.
				old := atomic.SwapUint32(ptr, newVal)
				res += fmt.Sprintln("Old Val:", old)

				res += fmt.Sprintln("After:", s)

				fmt.Println(res)
			}
		}()
	}

	wg.Wait()
}

/*
UINT64S
*/
func AddUint64Usage() {
	var (
		val uint64
		wg  sync.WaitGroup
	)

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 3; i++ {
				// AddUint64 atomically adds delta to *addr and returns the new value.
				// To subtract a signed positive constant value c from x, do AddUint64(&x, ^uint64(c-1)).
				// In particular, to decrement x, do AddUint64(&x, ^uint64(0)).
				// Consider using the more ergonomic and less error-prone [Uint64.Add] instead (particularly if you target 32-bit platforms; see the bugs section).
				fmt.Println(atomic.AddUint64(&val, 1))
			}
		}()
	}

	wg.Wait()

	fmt.Println("Result:", val)

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				// To subtract a signed positive constant value c from x, do AddUint64(&x, ^uint64(c-1)). In
				// particular, to decrement x, do AddUint64(&x, ^uint64(0)). Consider using the more ergonomic and less error-prone [Uint64.Add] instead (particularly if you
				// target 32-bit platforms; see the bugs section).
				fmt.Println(atomic.AddUint64(&val, ^uint64(3-1)))
			}
		}()
	}

	wg.Wait()

	fmt.Println("Result:", val)
}

func CompareAndSwapUint64Usage() {
	var (
		val            uint64 = 1
		oldVal, newVal uint64 = 10, 20
		wg             sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		// CompareAndSwapUint64 executes the compare-and-swap operation for a uint64 value.
		// Consider using the more ergonomic and less error-prone [Uint64.CompareAndSwap] instead (particularly if you target 32-bit platforms; see the bugs section).
		fmt.Printf("old: %d | new: %d | res: %t\n", val, oldVal, atomic.CompareAndSwapUint64(&val, val, oldVal))
	}()

	go func() {
		defer wg.Done()
		// CompareAndSwapUint64 executes the compare-and-swap operation for a uint64 value.
		// Consider using the more ergonomic and less error-prone [Uint64.CompareAndSwap] instead (particularly if you target 32-bit platforms; see the bugs section).
		fmt.Printf("old: %d | new: %d | res: %t\n", oldVal, newVal, atomic.CompareAndSwapUint64(&val, oldVal, newVal))
	}()

	wg.Wait()

}

func LoadUint64Usage() {
	var (
		val uint64 = 333
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// LoadUint64 atomically loads *addr. Consider using the more ergonomic and less error-prone [Uint64.Load] instead (particularly if you target 32-bit platforms; see
				// the bugs section).
				fmt.Println(atomic.LoadUint64(&val))

				if rand.Intn(100)%2 == 0 {
					// AddUint64 atomically adds delta to *addr and returns the new value. To subtract a signed positive constant value c from x, do AddUint64(&x, ^uint64(c-1)).
					// In particular, to decrement x, do AddUint64(&x, ^uint64(0)). Consider using the more ergonomic and less error-prone [Uint64.Add] instead (particularly if you
					// target 32-bit platforms; see the bugs section).
					atomic.AddUint64(&val, ^uint64(0))
				}
			}

		}()
	}

	wg.Wait()
}

func StoreUint64Usage() {
	var (
		val uint64
		wg  sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 3; i++ {
				// StoreUint64 atomically stores val into *addr.
				// Consider using the more ergonomic and less error-prone [Uint64.Store] instead (particularly if you target 32-bit platforms; see the bugs section).
				atomic.StoreUint64(&val, rand.Uint64())

				// LoadUint64 atomically loads *addr.
				// Consider using the more ergonomic and less error-prone [Uint64.Load] instead (particularly if you target 32-bit platforms; see the bugs section).
				fmt.Println(atomic.LoadUint64(&val))
			}
		}()
	}

	wg.Wait()
}

func SwapUint64Usage() {
	var (
		s = make([]uint64, 10)

		wg sync.WaitGroup
	)

	for i := 0; i < len(s); i++ {
		s[i] = rand.Uint64()
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {

				res := fmt.Sprintln("Before:", s)

				ptr := &s[rand.Intn(len(s))]
				newVal := rand.Uint64()

				// SwapUint64 atomically stores new into *addr and returns the previous *addr value. Consider using the more ergonomic and less error-prone [Uint64.Swap] instead (particularly if you target 32-bit platforms; see the bugs section).
				old := atomic.SwapUint64(ptr, newVal)
				res += fmt.Sprintln("Old Val:", old)

				res += fmt.Sprintln("After:", s)

				fmt.Println(res)
			}
		}()
	}

	wg.Wait()
}

/*
POINTERS
*/
func CompareAndSwapPointerUsage() {
	var (
		arr = []int{1, 2, 3}
		ptr = unsafe.Pointer(&arr[0])
	)

	// CompareAndSwapPointer executes the compare-and-swap operation for a unsafe.Pointer value.
	// Consider using the more ergonomic and less error-prone [Pointer.CompareAndSwap] instead.
	fmt.Println(atomic.CompareAndSwapPointer(&ptr, unsafe.Pointer(&arr[0]), unsafe.Pointer(&arr[1])))
	fmt.Println(*(*int)(ptr))
}

func LoadPointer() {
	var (
		val int64 = 777
		ptr       = unsafe.Pointer(&val)

		wg sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// LoadPointer atomically loads *addr. Consider using the more ergonomic and less error-prone [Pointer.Load] instead.
				newPtr := (*int64)(atomic.LoadPointer(&ptr))
				*newPtr += 100

				fmt.Println(*newPtr)
			}

		}()
	}

	wg.Wait()
}

func StorePointerUsage() {
	var (
		s   = make([]int, 10)
		ptr = unsafe.Pointer(&s[0])

		wg sync.WaitGroup
	)

	for i := range s {
		s[i] = rand.Intn(128)
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {
				// StorePointer atomically stores val into *addr.
				// Consider using the more ergonomic and less error-prone [Pointer.Store] instead.
				atomic.StorePointer(&ptr, unsafe.Pointer(&s[rand.Intn(len(s))]))

				// LoadInt64 atomically loads *addr. Consider using the more ergonomic and less error-prone [Int64.Load] instead (particularly if you target 32-bit platforms; see the bugs
				// section).
				fmt.Println(atomic.LoadInt64((*int64)(ptr)))
			}

		}()

	}

	wg.Wait()
}

func SwapPointerUsage() {
	var (
		s   = make([]int32, 10)
		ptr = unsafe.Pointer(&s[0])

		wg sync.WaitGroup
	)

	for i := 0; i < len(s); i++ {
		s[i] = rand.Int31n(128)
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {

				res := fmt.Sprintln("Before:", s)

				// SwapPointer atomically stores new into *addr and returns the previous *addr value. Consider using the more ergonomic and less error-prone [Pointer.Swap] instead.
				new := unsafe.Pointer(&s[rand.Intn(len(s))])
				old := atomic.SwapPointer(&ptr, new)
				res += fmt.Sprintln("Old Val:", old)
				res += fmt.Sprintln("New Val:", new)
				*(*int32)(ptr) = rand.Int31n(128)

				res += fmt.Sprintln("After:", s)

				fmt.Println(res)
			}
		}()
	}

	wg.Wait()
}

/*
UINTPTRS
*/
func AddUintptrUsage() {
	var (
		s   = [10]int{}
		ptr = uintptr(unsafe.Pointer(&s))

		wg sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		s[i] = i * 10
	}

	wg.Add(2)
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			// AddUintptr atomically adds delta to *addr and returns the new value. Consider using the more ergonomic and less error-prone [Uintptr.Add] instead.
			fmt.Println(*(*int)(unsafe.Pointer(
				atomic.AddUintptr(
					&ptr, uintptr(unsafe.Sizeof(new(int))),
				)),
			))
		}
	}()

	go func() {
		defer wg.Done()
		// AddUintptr atomically adds delta to *addr and returns the new value. Consider using the more ergonomic and less error-prone [Uintptr.Add] instead.
		for i := 0; i < 4; i++ {
			fmt.Println(*(*int)(unsafe.Pointer(atomic.AddUintptr(&ptr, uintptr(unsafe.Sizeof(new(int)))))))
		}
	}()

	wg.Wait()

}

func CompareAndSwapUintptrUsage() {
	var (
		arr = []int{1, 2, 3}
		ptr = uintptr(unsafe.Pointer(&arr[0]))
	)

	// CompareAndSwapPointer executes the compare-and-swap operation for a unsafe.Pointer value.
	// Consider using the more ergonomic and less error-prone [Pointer.CompareAndSwap] instead.
	fmt.Println(atomic.CompareAndSwapUintptr(&ptr, uintptr(unsafe.Pointer(&arr[0])), uintptr(unsafe.Pointer(&arr[1]))))
	fmt.Println(*(*int)(unsafe.Pointer(ptr)))
}

func LoadUintptrUsage() {
	var (
		val uint64 = 333
		ptr        = uintptr(unsafe.Pointer(&val))

		wg sync.WaitGroup
	)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// LoadUintptr atomically loads *addr. Consider using the more ergonomic and less error-prone [Uintptr.Load] instead.
				newPtr := (*uint64)(unsafe.Pointer(atomic.LoadUintptr(&ptr)))
				*newPtr += 100
			}

		}()
	}

	wg.Wait()

	fmt.Println(*(*uint64)(unsafe.Pointer(atomic.LoadUintptr(&ptr))))
}

func StoreUintptrUsage() {
	var (
		s   = make([]int, 10)
		ptr = uintptr(unsafe.Pointer(&s[0]))

		wg sync.WaitGroup
	)

	for i := range s {
		s[i] = rand.Intn(128)
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {
				// StoreUintptr atomically stores val into *addr. Consider using the more ergonomic and less error-prone [Uintptr.Store] instead.
				atomic.StoreUintptr(&ptr, uintptr(unsafe.Pointer(&s[rand.Intn(len(s))])))

				// LoadInt64 atomically loads *addr. Consider using the more ergonomic and less error-prone [Int64.Load] instead (particularly if you target 32-bit platforms; see the bugs
				// section).
				fmt.Print(atomic.LoadInt64((*int64)(unsafe.Pointer(ptr))), " ")
			}

		}()

	}

	wg.Wait()
}

func SwapUintptrUsage() {
	var (
		s   = make([]int32, 10)
		ptr = uintptr(unsafe.Pointer(&s[0]))

		wg sync.WaitGroup
	)

	for i := 0; i < len(s); i++ {
		s[i] = rand.Int31n(128)
	}
	fmt.Println(s)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 5; i++ {

				res := fmt.Sprintln("Before:", s)

				// SwapUintptr atomically stores new into *addr and returns the previous *addr value. Consider using the more ergonomic and less error-prone [Uintptr.Swap] instead.
				new := uintptr(unsafe.Pointer(&s[rand.Intn(len(s))]))
				old := atomic.SwapUintptr(&ptr, uintptr(new))
				res += fmt.Sprintln("Old Val:", old)
				res += fmt.Sprintln("New Val:", new)
				*(*int32)(unsafe.Pointer(ptr)) = rand.Int31n(128)

				res += fmt.Sprintln("After:", s)

				fmt.Println(res)
			}
		}()
	}

	wg.Wait()
}

/*
	TYPES
*/

/*
	atomic.Bool
*/

func BoolUsage() {
	var (
		flag = &atomic.Bool{}
		wg   sync.WaitGroup
	)

	fmt.Print("CAS:")

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("false -> true: ", flag.CompareAndSwap(false, true), ", ")
	}()

	go func() {
		defer wg.Done()
		fmt.Print("true -> false: ", flag.CompareAndSwap(true, false), ", ")
	}()

	wg.Wait()
	fmt.Println()

	fmt.Println("Load:", flag.Load())

	flag.Store(false)
	fmt.Println("Store:", flag.Load())

	fmt.Println("Swap:", flag.Swap(true), flag.Load())

}

/*
	atomic.Int32
*/

func Int32Usage() {
	var (
		val = &atomic.Int32{}
		wg  sync.WaitGroup
	)

	fmt.Print("CAS: ")

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("0 -> 10: ", val.CompareAndSwap(0, 10), ", ")
	}()

	go func() {
		defer wg.Done()
		fmt.Print("10 -> 20: ", val.CompareAndSwap(10, 20), ", ")
	}()

	wg.Wait()
	fmt.Println()

	fmt.Print("Add: ")
	wg = sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		val.Add(100)
	}()

	go func() {
		defer wg.Done()
		val.Add(200)
	}()
	wg.Wait()

	fmt.Println(val.Load())

	fmt.Println("Load:", val.Load())

	val.Store(30)
	fmt.Println("Store:", val.Load())

	fmt.Println("Swap:", val.Swap(40), val.Load())

}

/*
	atomic.Int64
*/

func Int64Usage() {
	var (
		val = &atomic.Int64{}
		wg  sync.WaitGroup
	)

	fmt.Print("CAS: ")

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("0 -> 10: ", val.CompareAndSwap(0, 10), ", ")
	}()

	go func() {
		defer wg.Done()
		fmt.Print("10 -> 20: ", val.CompareAndSwap(10, 20), ", ")
	}()

	wg.Wait()
	fmt.Println()

	fmt.Print("Add: ")
	wg = sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		val.Add(100)
	}()

	go func() {
		defer wg.Done()
		val.Add(200)
	}()
	wg.Wait()

	fmt.Println(val.Load())

	fmt.Println("Load:", val.Load())

	val.Store(30)
	fmt.Println("Store:", val.Load())

	fmt.Println("Swap:", val.Swap(40), val.Load())

}

/*
	atomic.Pointer
*/

func PointerUsage() {
	var (
		a   int = 0
		val     = &atomic.Pointer[int]{}
		wg  sync.WaitGroup
	)

	val.Store(&a)

	fmt.Print("CAS: ")

	b, c := 10, 20

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("0 -> 10: ", val.CompareAndSwap(&a, &b), ", ")
	}()

	go func() {
		defer wg.Done()
		fmt.Print("10 -> 20: ", val.CompareAndSwap(&b, &c), ", ")
	}()

	wg.Wait()
	fmt.Println()

	fmt.Println("Load:", val.Load())

	d := 30
	val.Store(&d)
	fmt.Println("Store:", val.Load())

	e := 40
	fmt.Println("Swap:", val.Swap(&e), val.Load())

}

/*
atomic.Uint32
*/
func Uint32Usage() {
	var (
		val = &atomic.Uint32{}
		wg  sync.WaitGroup
	)

	fmt.Print("CAS: ")

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("0 -> 10: ", val.CompareAndSwap(0, 10), ", ")
	}()

	go func() {
		defer wg.Done()
		fmt.Print("10 -> 20: ", val.CompareAndSwap(10, 20), ", ")
	}()

	wg.Wait()
	fmt.Println()

	fmt.Print("Add: ")
	wg = sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		val.Add(100)
	}()

	go func() {
		defer wg.Done()
		val.Add(200)
	}()
	wg.Wait()

	fmt.Println(val.Load())

	fmt.Println("Load:", val.Load())

	val.Store(30)
	fmt.Println("Store:", val.Load())

	fmt.Println("Swap:", val.Swap(40), val.Load())

}

/*
atomic.Uint64
*/
func Uint64Usage() {
	var (
		val = &atomic.Uint64{}
		wg  sync.WaitGroup
	)

	fmt.Print("CAS: ")

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("0 -> 10: ", val.CompareAndSwap(0, 10), ", ")
	}()

	go func() {
		defer wg.Done()
		fmt.Print("10 -> 20: ", val.CompareAndSwap(10, 20), ", ")
	}()

	wg.Wait()
	fmt.Println()

	fmt.Print("Add: ")
	wg = sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		val.Add(100)
	}()

	go func() {
		defer wg.Done()
		val.Add(200)
	}()
	wg.Wait()

	fmt.Println(val.Load())

	fmt.Println("Load:", val.Load())

	val.Store(30)
	fmt.Println("Store:", val.Load())

	fmt.Println("Swap:", val.Swap(40), val.Load())

}

/*
atomic.Value
*/
func ValueUsage() {
	var (
		val = &atomic.Value{}

		wg sync.WaitGroup
	)
	val.Store("qwerty")

	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Print("Store: ", "abcdefg\n")
		val.Store("abcdefg")
	}()

	go func() {
		defer wg.Done()
		fmt.Printf("CAS: %q -> %q: ", "abcdefg", "poiuyt")
		fmt.Println(val.CompareAndSwap("abcdefg", "poiuyt"))
		fmt.Println()
	}()

	wg.Wait()

	fmt.Println("Swap to  \"zxcvbn\":", val.Swap("zxcvbn"))
	fmt.Println("Load:", val.Load())
}
