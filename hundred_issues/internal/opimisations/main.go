package main

import "sync"

/*
	OPTIMISATIONS
In most cases readable and understandable code is better than optimized and obfuscated one.
	Make it right, make it clean, make it consice, make it fast - in this order.

	CPU CACHE UNDERSTANDING
1. You needn't be an engineer to be a good racer, it's sufficient to feel the machine.

	CPU ARCHITECTURE
1. The modern CPUs use caching to make memory access faster. In most cases there are three levels of caching:
	- L1 cache: 64KB
	- L2 cache: 256KB
	- L3 cache: 4 MB
2. Logical kernels are called "virtual kernels" or "threads".
3. In intel family splitting a physical kernel to some logical kernels is called "Hyper-Threading"
4. Caching is about saving data and instructions
5. If memory area is closer to logical kernel (thread), an access to it is also faster:
	- L1: 						1 nanosecond 				(arranged on the same crystal as the other part of CPU)
	- L2: ~4*L1 nanosecond 		(~4 nanoseconds on L2) 		(arranged on the same crystal as the other side of CPU)
	- L3: ~2.5*L2 nanosecond 	(~10 nanoseconds on L3) 	(L3 arranged on ther side of CPU)
6. L1 and L2 are called "on-die". It means that they are arranged on the same crystal as the other part of CPU.
7. L3 is called "off-die". It means that there's the bigger delay between data access.
8. RAM speed is less than L1 speed by 50-100 times. While Accessing a single variable in RAM, we can access 100 variables in L1 cache.
9. L3 stores the data repressed from L2 cache.

	CACHE LINE
1. After accessing a specific domain in memory (for example: reading a variable value) the two actions might happen:
	1) This domain can be referenced again by other side.
	2) The referencing to close domains will be done
The first possibility relates to temporary locality and the second to spacial locality. Both possibilities are parts of Reference Locality Principle
(Locality of Reference).
2. Temporary Locality - a reason of need of CPU caches: to make repeatable access to the same variables faster. But because of Spatial Locality, CPU copies
from RAM all the Cache Line, not only single variable.
3. Cache Line - is a continuous segment of memory of fixed size, ordinarily 64 bytes (8 variables of int64). Each time when CPU decides to copy a block of
RAM, it copies this block to a row of cache.
4. Since the RAM is hierarchical, when CPU wants to access a specific memory cell in memory, firstly it checks its presence in L1 cache, then L2, then L3 and
finally in RAM.

	SLICE OF STRUCTURES AND STRUCTURE OF SLICES

	PREDICTION
1. Prediction is the capability of CPU to predict what an application will do to make it faster.
2. In the 64-bit architecure the word is 64-bit, in the 32-bit architecture the word is 32 bit.
3. The concept of striding. There are 3 types of strides:
	1) The stride of size 1 (Unit stride) - all the values we want to access are arranged sequentially, for example: slice of int64. This stride is predictable for
	processor and most efficient.
	2) The stride of constant size (Constant stride) - this stride is also predictable, for example: the looping over each second elem. This stride requires bigger
	amount of cache lines and less efficient than Unit stride.
	3) Unpredictable stride size (Non-unit stride). The stride that is cannot be predicted by CPU, for example: linked list or slice of pointers. Since the CPU
	doesn't know whether the data is allocated sequentially or not, it won't overload the cachelines.

	CACHE ARRANGING STRATEGY
1. When CPU decides to copy a segment of memory and dump it into cache, it must follow the specific strategy.
	1) If cache L1D has the 32KB size and the cache line is 64 bytes. If the segment is arranged in L1D randomly, in the worst case the processor will need to make
	an iteration over 512 cache lines, to calculate a variable value. This type of cache is called "fully associative".
2. The most-usable type of cache is "set-associative" cache. Its principle is based on cache sectioning. In this case the case is separated into sectors. Each
sector has 2 lines\rows. A block of memory can belong to only one segment and its arranging is defined by its address in memory. To understand it we must shrink
the address of the block to three parts:
	1) Block offset. It depends on the size of block. In our case it has the size of 512 bytes and 512 is 2^9 -> the 9 bits of the address represent the offset of
	the block. (block offset = bo)
	2) Sector index. Sector index points to a sector that is related by an address. Since the cache is double-ended set-associative and has 8 lines, there are only
	8/2 = 4sectors. Moreover, 4 is 2^2, so the next 2 bits represent set index (set index = si)
	3) The other part of the address consists of tag bits (tag bits = tb). Since for the simplicity we use the address in 13-bits representation to calculate
	tag bits we use the following: tb = 13 - bo - si -> 2. It means that the other remainig bits represent tag bits.
3. All the modern caches are splitted into sectors. In such cases, based on strides it can turn out that the only one sector will be used. It can result in perform-
ance of application and affect conflictual cache misses. Stride of this type is called "critical". For applications that require an achievement of high performance
we should avoid critical strides in order to get good performance of CPU.

	TO WRITE CODE THAT RESULTS IN FALSE MUTUAL USING
1. A goal of the CPU is providing caches coherence. For example: if a goroutine updates sumA and another one reads sumA (after a synchronization) we wait the
application to observe the latest value.
*/

func main() {

}

/*
1. The Principle of Temporary Localization is applied to these variables:
  - i
  - length
  - total

At each iteration we reference to these variables.
2. The Principle of Spatial Localization is applied to these variables:
  - Code instructions
  - Slice s. Since behind this slice the underlying slice stands and all the elements of it are arranged at the same memory area and close to each other, the
    access to s[0] means the access to s[1], s[2] and other elems.

3. When sumA references to s[0] this adress hasn't been arranged in cache. If the CPU decides to cache this variable, it'll copy the whole segment of memory (64
bytes)
4. Firstly, referencing to s[0] results in cache miss because the memory address hasn't been arranged in cache yet. It's called "compulsory miss". But if the CPU
gets the access to 0x000 memory block (referencing to elems from 1st to 7th), it'll result in from cache reading (dumping to the cache). The same logic is applied
when sumA references to s[8]. As a result looping over 16 elements has resulted in: 2 compulsory misses and 2 forced cache dumps.
*/
func sumA(s []int64) int64 {
	var (
		total  int64
		length = len(s)
	)

	for i := 0; i < length; i++ {
		// When there will be the first referencing to s[0] this address hasn't been arranged in cache.
		// If the CPU decides to cache this variable, it'll cache the whole segment of memory (64 bytes) from RAM.
		total += s[i]
	}

	return total
}

type (
	Foo struct {
		a int64
		b int64
	}

	Bar struct {
		a []int64
		b []int64
	}
)

/*
In this case. If the foos will have 16 elems and bar will have the same amount of elems a and b, the first function will use 4 cache lines and the second only 2.
It results in the speed performance more than 20% percent finally,

To optimize an application we need to organize our data in the way when we get maximum efficiency from each cache line.
*/
func sumFoo(foos []Foo) int64 {
	var total int64

	for i := 0; i < len(foos); i++ {
		total += foos[i].a
	}

	return total
}

func SumBar(bar Bar) int64 {
	var total int64

	for i := 0; i < len(bar.a); i++ {
		total += bar.a[i]
	}

	return total
}

type node struct {
	val  int64
	next *node
}

/*
In both cases we have the same density of elements arranging. But the execution of sum2 is much faster (~70%) and that's why:
*/
func linkedListSum(root *node) int64 {
	var total int64

	for root != nil {
		total += root.val
		root = root.next
	}

	return total
}

func sumB(s []int64) int64 {
	var total int64

	for i := 0; i < len(s); i += 2 {
		total += s[i]
	}

	return total
}

func calculateSum512(s [][512]int64) int64 {
	var sum int64

	for i := 0; i < len(s); i++ {
		for j := 0; j < 8; j++ {
			sum += s[i][j]
		}
	}

	return sum
}

type (
	Input struct {
		a int64
		b int64
	}

	Result struct {
		sumA int64
		sumB int64
	}
)

/*
In this case we have 2 kernels (P) that has 2 OS threads and 2 goroutines working on incrementing the variables of res.
Both cache lines are replicated because L1D relates to each kernel.

A goal of the CPU is providing caches coherence. For example: if a goroutine updates sumA and another one reads sumA (after a synchronization) we wait the
application to observe the latest value.

In our case it doesn't happen. Both goroutines reference to theirs own variables, not to mutual. The processor doesn't know that there's the conflict. 
*/
func count(inputs []Input) Result {
	var (
		wg  = &sync.WaitGroup{}
		res = Result{}
	)

	wg.Add(2)
	go func() {
		defer wg.Done()

		for i := 0; i < len(inputs); i++ {
			res.sumA += inputs[i].a
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < len(inputs); i++ {
			res.sumB += inputs[i].b
		}
	}()

	wg.Wait()

	return res
}
