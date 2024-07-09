package main

import (
	"fmt"
	"math"
	"math/rand"
)

// https://habr.com/ru/companies/vdsina/articles/515660/
/*
	CACHES
1. Cache - is the tiny, but very efficient and fast memory arranged close to logical segments of CPU

2. The modern processors need only a single tick to sum two 64-bit integers.
	For example: if the processor works with 4GHZ, it needs only 1/4 of nanosecond to do this.
On the other hand HDD needs thousands of nanoseconds, SSD needs hundreds of nanoseconds and we cannot integrate HDD / SSD into CPUs, so there's space between CPU and these ones.

3. From the point 2) we make the conclusion that we need an intermediate memory that is faster than storage device and possible to manipulate plenty of operations of data
transmittions and is closer to the CPU. And there's the solution represented by the RAM.
	1) Most of these storages have the type is called DRAM "dynamic random access memory". They're possible to transmit data much faster than any storage. Inspite of this speed,
	DRAM isn't capable of storing tons of information.
	2) DRAM spends about 100 nanoseconds to find and pass data but it can pass billons bits per second.

4. So far we have CPU, DRAM, Storages (behalf of SSD/HDD) but we need another intermediary. SRAM "Static Random Access memory" takes its turn.
	1) DRAM uses micro-condensators to store data as an electric charge.
	2) SRAM uses transistors that solve the same problem that work with the same speed to the logical parts of CPU (Threads).
		1) SRAM is 10 times faster than DRAM
	3) While arranging the data processing units closer to CPU we loose the capability of storing tons of data. SRAM isn't exception.
	4) SRAM is created with the same technology as CPU, so we can embed the SRAM units into CPU. These units are very close to CPU logical kernels (threads)
Cache - is the set of SRAM units arranged inside the processor.
	1) These SRAM blocks provide max busyness of the CPU thanks to data transmission and data storing with the very fast speeds.

5. CPU internals. Pic: 0_CPU_Internals
	1) The CPU kernel is represented as the area inside the dotted line.
	2) ALU ("Arithmetic Logic Units") - the structures processing math operations.
	3) Registers are ordered into a register file.
		1) Registers are also SRAM but theirs speed is less or equal to ALU they service. Each register stores the only one number. For example 64-bit number. This value can be an element
		of specific data: Code of an instruction or an Address of emory of other data.

6. L1D (data cache) cache stores >=32KB of data and works at the same speed as logical kernels

7. L1I (instructions cache) also stores >=32KB of varying instructions ready to be splitted to smaller ones that must be processed on ALU.
	1) These instructions require its own L0 cache that stores 1,500 operations and is closer to ALU.

8. To get the data needed from the L1 cache we need less than 5 CPU ticks (~1.1 nanosecond).

9. L2 is the unified cache of Instructions & Data.
	0) L2 is also embedded into a physical kernel of the CPU.
	1) The size of L2 cache is >=256KB. This size is picked to provide lower caches with data enough to be processed.
	2) To find and pass the data from the L2 cache we need x2-x4 time relatively to the L1 cache (~3.3 nanoseconds).

10. L3 cache is the mutual cache for all kernels.
	1) Each kernel can freely access the data of another kernel from the L3 cache.
	2) The size of L3 cache can be 2-32 MB
	3) L3 cache has much slower speed relatively to L1 and L2 caches. It spends about CPU 30 ticks to find and pass data further (12,8 ns nanoseconds).

11. Cache makes performance faster, speeding data transmission to logical units and saving a copy of most usable data/instuctions nearly.
	1) Data saved in cache is splitted into 2 parts: data itself and the address of this data in DRAM / Storage (behalf of HDD or SSD)
	2) The address of data stored is called "cache tag"

12. When the CPU performs an operation that needs to read or write from/into memory, it starts from the check of L1 Cache:
	1) If the data has been found ("cache hit" has happend), this access is performed at once
	2) If the tag needed hasn't been found, "cache miss" happened. A new cache tag is created in the L1 cache and other part of the CPU performing search in the other cache levels
	takes its turn (this search can reach up to Storage level).
	3) If all the cache slots are busy (All N slots), to keep the new tag created in L1 cache we need to dump some data to the L2 cache

13. All the dumps result in constant data mixing between all the levels of caches is performed during a few CPU ticks

14. Cache line is 64-byte ordinarily.

15. Set-associative means that a block of system data is bound to the set of lines of a specific set in cache and this block cannot be freely bound to other lines.
	1) "8-way" means that a single block can be bound to 8 rows in cache in the set. The greater associativity, the more chances for data to get into cache during the CPU search and less
	losses happen resulted from "cache misses".
	2) The lacks of this system hide in complexity increasing and energy consumption and performance collapse because we need to procces more rows for a block of data.

16. Incluvisity Policy
	1) "Inclusive cache" means that the data of one cache level can reside on another level. For example: data in L1 cache can reside in L3 cache. It results in faster search during data
	search and absence of the need to traverse higher levels.
	2) L2 cache is "exclusive": all the data resides exclusively for this level and no data is dumped to other levels.
		1) It economizes the place, but results in the cache system needs to traverse L3 cache for the data needed.
	3) "Victim" caches are used to store information that is dumped from the lower cache levels.

17. Write policies.
	1) Most of the modern processors use "write-back" policy. It means that when data is being written into cache levels, there's a delay between a write of the data copy into system memory.
	This pause lasts during the time while data is in cache - RAM gets this information after cache dumping.
*/

// https://teivah.medium.com/go-and-cpu-caches-af5d32cc5592
/*
1. Modern processors are based on concept of Symmetric MultiProcessing (SMP). In SMP system the processor designed so that two or more cores are connected to a shared memory (also called
main memory or Random Acces Memory "RAM").

2. Also to speed up memory access, the processor has different levels of cache called L1, L2, L3. The most frequent model is to have (L1 & L2) local to a core and L3 shared across all cores.

3. The closer cache is to a CPU core, the smaller is its access latency and size.
	1) L1 has latency ~1,2 ns, CPU cycles ~4, Size 32-512 KB
	2) L2 has latency ~4 ns, CPU cycles ~10, Size 128-24 MB
	3) L3 has latency ~12 ns, CPU cycles ~40, Size 2-32 MB
	4) DRAM has latency ~50-100 ns and it's slower 50 times than L1 cache

4. Locality of reference. When a processor accesses to a particular memory location, it's very likely that:
	1) "Temporary Locality Principle"
	It will access the same location in the near future: this is "Temporary Locality Principle". This is one of the reason to have CPU caches.
	2) "Spatial Locality Principle"
	It will access memory locations nearby: this is "Spatial Locality Principle". Instead of copying a single memory location to the CPU caches, the solution is to copy a "cache line". A
	cache line is contiguous segment of memory.

5. The default size of a cache line is 64-byte. On this macbook (MBP M3 Pro cacheline is 128 byte). It means that instead of copying a single variable to the L1 cache, the processor will
copy a contiguous segment of 64 bytes. For example: instead of copying a single int64 element of a Go's slice, it'll copy this element plus seven elems as well.

6. In memory, all the different rows of a matrix are allocated contiguously. In the MatrixCombination examples the matrix size is a multiple of the cache line size. Hence,
a cache line won't "overtake" on the next row.

7. Consider a situation when the second matrix pointer moves downwards. When the pink pointer accesses to the position (0,4), the processor will cache the whole line.
Therefore, when the pink pointer reaches position (0,5), we may assume that the variable is already present in L1, isn't it?
	1) If the matrix is small enough compared to the size of L1, then yes, the element (0,5) will already be cached.
	2) Otherwise, the cache line will be evicted (dumped) from L1 before the pointer reaches position (0,5). Therefore it'll generate cache miss and the processor will
	have to access the variable differently (using L2 for example).
	3) To optimize matrices optimally we need to arrange two matrices into cachelines optimally: we take the size of cachelines as CacheLineSize, the size of L1D in bytes as
	L1DByteSize and make the following calculations: L1DByteSize / CacheLineSize / 2 -> 256 bytes for each matrix.

8. When a processor needs to access a memory location, there's a translation from the virtual to the physical memory. Using LEA (Load Effective Address) allows computing
a memory address without having to go through this translation. For example: if we manage a slice of int64s elems and that we already have a pointer to the first element
address, we can use LEA to load the second element address by simply shifting the pointer to 8*index bytes. In our example it might be a potential reason why the second test
(when we split the L1 cache for 2 matrices) is faster.

9. MBP Pro M3 Pro cache sizes:
	1) Cache line: "sysctl -a | grep cacheline" 128 Bytes
	2) L1D: "sysctl hw.l1dcachesize" 64 KB for data and 128 KB for instructions. (Each core has this L1D and L1I caches)
	3) L2: "sysctl hw.l2cachesize" 4096 KB -> 4 MB of data (Each core also has the L2 cache)

10. The limiting impacts of cache misses in the case of larger matrix. "Loop Nest Optimization". We have to iterate only within a given block to benefit from cache lines
as much as possible. When we loop over a single row, we dump the previous line into L2 cache and it increases the time, so we need to go through the blocks of specific
size. We upload to the cache only blocks of specific size that can be saved during execution.
	1) We shoud accurately choose the block size to be within L1 cache size boundaries.

11. Outcomes
	1. Because of the Spatial Locality Principle, the processor doesn't put a simple memory address but a cache line.
	2. The L1 cache is local to a given CPU core.

12. Cache coherency (False sharing)
1. Two variables v1 and v2 are stored in DRAM. One thread on a core accesses to v1 whereas another thread on another core accesses v2. Assuming that both are contiguous
(or almost), we end up in a situation where v2 is present in both caches. The question is: what happens if the 1st thread updates its cache line? It could potentially
update any location of its line including v2. Then, when the 2nd thread reads v2, the its value may not be consistent anymore.

2. If two cache lines share some common adresses, the processor will mark  them as "Shared" ones.
	1) If one thread (logical processor) modifies a "Shared" line, it'll mark both as "Modified".
	2) To guarantee caches coherency, it requires coordination between the cores which may significantly degrade the application performance. This problem is called
	"False sharing"

3. "Memory padding". To beat the problem of "False sharing" we should use "Memory Padding". The goal of this method is to make sure that there's enough space between the two
variables (its contents) so that they belong to different cache lines.
	1) We need to discover our L1D cache size
	2) We need to make an appropriate alignment to be confident that two structures is belonged to different cachelines.
	3) It can be done with using the placeholder with the size of alignment needed. We can use int64 or int32 to achieve it.
	4) The lack of the alignment is the bigger amount of memory needed to distribute DRAM areas to different cachelines.
*/

// https://dev.to/ashevelyov/understanding-cpu-cache-and-prefetching-in-go-44bk
/*
1. Why Cache Matters?
The CPU cache acts as a ready stash of pre-tied memory, ensuring the the sprinter (CPU) keeps racing ahead without unnecessary stops.

Arrays in go, due to their continuous memory layout, play a pivotal role in this caching dance. When we access an array element, it's not just that specific element
that's whisked into the cache; a chunk of nearby elems comes along for the ride.

2. A prefetcher is like that assistant on the sidelines who notices the sprinter's predictable pattern and gets those shoes ready even before they're needed. If we're
looping over an array in Go, the prefetcher catches on quickly. It anticipates future data needs and ensures the cache is primed and ready.

3. The concert line-up (Array vs. Linked List)
Imagine a concert where attendees (data points) stand in:
	1) A line (array). The security checks each attemdee one by one. This is efficient since everyone is in order.
	2) attendees are scattered throughout the venue in random seats (linked list) and the security has to dart from one place to another to check each person.
The difference between looping over this structure the following: to traverse linked list we spend x4 time of traversing the slice.
*/

// https://mecha-mind.medium.com/demystifying-cpu-caches-with-examples-810534628d71
/*

1. Some notable factors can affect the performance of a computer program utilizing the CPU cache are as follows:
	1) Impact of cache lines
	2) Cache associativity
	3) False Cache Line Sharing

2. Important pointers regarding CPU caches:
	1) Reading / Writing to main memory (RAM (DRAM)) is expensive as compared to CPU caches (50-100 nanosecs).
	2) The access times for RAM (DRAM) is 62.9 ns whereas it's 12,8 ns for L3 cache, 3,3 ns for L2 cache, 1,1 ns for L1 cache
	3) During each read query, it's first checked in CPU's L1 cache, if it's found then returned, else check L2 cache and so on till L3 cache or L4 cache (if present)
	4) If nothing was found in CPU caches, fetch from main memory (DRAM) and add to L1, L2, L3 caches
	5) When a byte is read not only the byte is returned, but a block of bytes is returned known as "cache line"
	6) Generally the "cache line" is 64-bytes. Any read for a byte in the same cache line is read from L1, L2 or L3 cache because the cache line of 64 bytes is already
	cache in the CPU.
	7) For e.g if the cache lines are aligned s.t. 0-63 bytes is one cache line, 64-127 is the next one and so on. Then if byte 10 is read, all bytes from 0 to 63 is
	transfered from RAM (DRAM) to the CPU and cached. Next if request for byte 57 is made, it's already present in the CPU's cache since 0 <= 57 <= 63, hence read from
	cache.
	8) A cache line read from main memory can occupy one of N cache slots i.e. if an existing cache line already occupies one of the N slots, then the incoming cache line
	can only occupy one of the remaining N-1 slots.
	9) If all the N slots are busy then we have to evict one of the N slots to make room for the new cache line, LRU policy is used to evict cache slots.

3. Finding L1, L2 and L3 caches' sizes is: sysctl -a | grep cache
	hw.cachelinesize: 128 Bytes
	hw.l1icachesize: 131072 is equal to 128KB
	hw.l1dcachesize: 65536 is equal to 64KB
	hw.l2cachesize: 4194304 is equal to 4MB

	hw.l1icachesize: 131072 is equal to 128KB
	hw.l1dcachesize: 65536 is equal to 64KB

4. L1 cache resides in each core of the CPU. For example: if a CPU is 12-core, then each core has 128 KB of L1I (Instructions) and 64KB of L1D (Data). The sum is 192KB. So
in the case of my macbook the total size of L1 caches is 2304 KB

5. Also L2 cache is also present in each core of the CPU. In the case of my macbook the size of L2 cache is 4MB, so the sum is 48 MB

6. In the case of L3 cache it's shared among all the physical kernels.
	1) Each byte present in L3 cache is also present in L2 and L1 caches.

7. Cache line size i.e. number of bytes prefetched for each byte read from main memory is 128 bytes.

8. L1 data associativity is 12 i.e. Each cache line can occupy one of the 12 slots in the L1 cache. Similarly for L2 cache, each line can occupy one of the 20 slots and
for L3 cache can occupy one of the 12 slots.

9. To calculate how many cache lines we need use to iterate over slice with the slice N we should use the following logic:
	1) Take count elements count as N, multiply it on the size of the element type (size of aggregated elems): N * M (bytes) = S
	2) Divide S by the size of cache line Q: S / Q = an amount of cachelines needed to iterate over the size of N elements with elements of size M.

10. Each cache line has an address in cache, that is bound to the address in DRAM. When the CPU references to a particular address in memory, it checks presence of this cache
line in cache.
	1) "Cache hit". If the cache line is present in cache, the CPU can get the data needed from this cache without referencing to DRAM. It makes this operation faster.
	2) "Cache miss". If the cache line corresponding to the address needed isn't present in cache, the CPU is responsible for loading this cache line from DRAM. It requires
	additional time and slows down the operation of data accessing. As a result, the load of a full cache line happens even if we need to access the only one elem.

11. Memory managing.
	1) "Replacement Policy". If the cache is full and a new "Cache miss" happens, a cache line is picked to be evicted from the cache to release the place for the new cache line
	2) "Prefetching". To decrease cache misses' delays the CPU can download data in advance. This data is assumed to be required by the CPU to be used in the close future. It
	can be done based on an analysis of data access sequence (Spatial locality).

12. "False sharing".
Assume that we have 4 threads: A, B, C, D and we repeatedly update integers at indices i, j, k, l. Now if i, j, k, l belong to the same cache line for e.g 0 -> i, j -> 1, k -> 2,
l = 3, then whenever a thread updates an integer, the entire cache line is invalidated and the entire cache line is again fetched from DRAM.
	1) Where each core updates its own L1 or L2 cache independently, the values need to be in sync across all cores. Syncing takes time. For example: if thread A updates the
	integer at index i = 0, then the thread B cannot update the integer at index j = 1 in its own L1/L2 cache concurrently because the entire cache line is invalidated and until
	the update is synced across all cores, the next update cannot happen in the same cache line.
	2) But if i, j, k, l belong to different cache lines for e.g. i = 0, j = 16, k = 32 and l = 48, then an update made by a thread doesn't invalidate the other cache lines and
	hence threads updating the other cache lines will update the integers in its own core's cache concurrently without having to sync across all cores.

13. Key Takeaway: Care must be taken while using multi-threading with CPU bound tasks such as above as the additional costs due to "context switching", creation and deletion of
threads and false cache line sharing can significantly reduce the gains.

14. Cache Associativity.
If we look at only the highest cache i.e. L2 cache, it has 4 MB cache size and 12-way associativity.
	1) 12-way associativity implies that each cache line can be placed in one of 12 slots in the L3 cache.
	2) Number of cache lines that can be placed in L2 cache = 4 * 1024 * 1024 (L2 size in Bytes) / 128 (Cache line size in Bytes) = 32,768 cache lines. The number 32,768 means
	the total number of slots available in the L2 cache.
	3) Since L2 is 12-ways associative, number of different 12-slots in the L2 cache is 32,768 / 12 = 2730 groups.
	4) But how to determine which cache lines go into which group?
	Since 0 ... 2,730 can be represented using 12 bits, thus if we index each cache line as 0, 1, 2 ... and so (for representing all the 32,768 cachelines we spend 15 bits) on
	then all cache lines whose lower 12 bits are the same, can be put in the same group in L2 cache. If there are more than 12 cache lines with 12 lower bits same, then there
	will be cache evictions in L3 cache.
Key Takeaway: When accessing elements in an array in a stepped fashion, try not to use a step size that is power of two.
*/

// https://en.wikipedia.org/wiki/Cache_replacement_policies
/*
	CACHE REPLACEMENT POLICY (LRU)
1. The average memory reference time is:
	T = Tm * m + (Th + E) where:
	1) Tm - time to make main memory (DRAM) access when there's a miss
	2) m - miss ratio = 1 - (hit ratio)
	3) Th - latency: time to reference the cache
	4) E - secondary effects
	We get:
	T = T[main-memory access] * (1 - (hit ratio)) + (T[cache-access] + E)
2. A cache has two primary figures of merit:
	1) latency
	2) hit ratio
A number of secondary factors also affect cache performance.
3. Least Recently Used (LRU) Policy.
	1) Discards least recently used items first.
	2) This algorithm requires keeping track of what was used and when, which is cumbersome. It requires "age bits" for cache lines and tracks the least recently used cache line
	based on them.
	3) When a cache line is used, the age of the other cache lines changes.
	4) The global counter set to 0. Each time we use cache line, the global counter is incremented and the cache block gets this counter value and is put to cache.
	5) If the cache is full, we pick the element with the smallest mark value and evict it from the cache, put the cacheline needed into the cache and assign incremented counter
	to cacheline wes used.
For example: the acces sequence is A B C D E F and cache size is 4 cache lines.
	1) After installation A B C D in the left-to-right order the marks will be the following: A(0) B(1) C(2) D(3)
	2) After the access of cache line E we pick the element with the smallest mark (A(0)) value and evict it from the cache, assign to this cell a new cache line (E) and assign
	the incremented counter to this cacheline (E(4))
	_E(4) B(1) C(2) D(3) | A
	3) After using a cacheline in cache we assign incremented global counter to cacheline used and assign it to the used cache line. For example we use B:
	E(4) B(5) C(2) D(3) | A
	4) Accessing F we evict from the cache value the element with the smallest mark and replace it to the element with the smallest mark and assign to this cache line incremented
	global counter:
	_E(4) _B(5) F(6) D(3) | A C
*/

/*
	MECHANICAL SYMPATHY
	 1. Data-Layout.
	    Organize your data structures in a way that minimizes cache-misses. This typically involves keeping data in contiguous memory blocks and using cache-friendly data
		structures like arrays or slices instead of linked lists.
	 2. Access patterns.
	    Access data in sequential or predictable pattern, allowing the CPU's cache prefetching to work efficiently. This is reffered to as a stride.

		Having a predictable stride pattern like the unit stride gives the CPU more predictability to import a cache line because it's well aware of the fact that the next elem
		will be accessed.

		Having a bad access pattern decays predictability for the CPU and causes internal latency.
*/

/*
	POINTERS PROBLEM
1. Assume that our structure fits the size of L1 cache and we've decided to use the pointer to this structure instead of using this structure itself. In this case instead of
copying the structure itself we copy the address of this structure and we have to go to the DRAM to get the contents of it.
2. The conclusion is: it's better to pass the structure itself instead of passing the pointer to it because of the larger time needed to go to DRAM.
3. The takeaway is: if we work with the small instances we can remove pointers at all. In this case we ease the work with CPU caches forcing it to cache elems directly and reduce
work of GC.
*/
func main() {
	fmt.Println(math.Log2(2730))
}

func createMatrix(size int) [][]int64 {
	var (
		matrix = make([][]int64, size)
	)
	for i := range matrix {
		matrix[i] = make([]int64, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			matrix[i][j] = int64(rand.Intn(128))
		}
	}

	return matrix
}

func MatrixCombination(a, b [][]int64) ([][]int64, error) {

	if len(a) != len(b) || len(a[0]) != (len(b[0])) {
		return nil, fmt.Errorf("sizes mismatch")
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a); j++ {
			// a[i][j] += b[i][j]
			a[i][j] += b[j][i]
		}
	}

	return a, nil
}
func MatrixCombinationLoopBlocking(a, b [][]int64) ([][]int64, error) {
	const (
		blockSize = 64

		rowSize = 1 << 14
	)

	if len(a) != len(b) || len(a[0]) != (len(b[0])) {
		return nil, fmt.Errorf("sizes mismatch")
	}

	for i := 0; i < len(a); i += blockSize {
		for j := 0; j < len(a); j += blockSize {

			for ii := i; ii < i+blockSize; ii++ {
				for jj := j; jj < j+blockSize; jj++ {
					a[ii][jj] += b[jj][ii]
					// a[ii][jj] += b[ii][jj]
				}
			}

		}
	}

	return a, nil
}

const (
	defaultSize = 1e6
)

type Node struct {
	val  int64
	next *Node
}

func createArray() []int64 {
	var (
		s = make([]int64, defaultSize)
	)

	for i := 0; i < defaultSize; i++ {
		s[i] = int64(rand.Intn(128))
	}

	return s
}

func createLinkedList() *Node {
	var (
		cur  = &Node{val: rand.Int63(), next: nil}
		root = cur
	)

	for i := 0; i < defaultSize; i++ {
		newNode := &Node{val: rand.Int63(), next: nil}
		cur.next = newNode
		cur = cur.next
	}

	return root
}

// Cache-friendly code using an array (or slice)
func CacheFriendlyStruct() {
	type Point struct {
		x, y float64
	}

	_ = make([]float64, 1e3)
}

func CacheFriendlyRanging() {
	array := make([]int, 1e3)
	for i := 0; i < 1e3; i++ {
		array[i] = rand.Intn(128)
	}

	// Bad access pattern: access every other element in the array
	// Constant stride
	/*
		In this case we refer to the old elems an using this way force processor to refer to the DRAM often.
		The DRAM accessing is 50-100 times slower than accessing L1D cache.
	*/
	badArray := make([]int, 1e3/2)
	for i := 0; i < 1e3/2; i += 2 {
		badArray[i/2] = badArray[i]
	}

	// Good access pattern: accessing neighboring elements in the array
	// Unit stride
	goodArray := make([]int, 1e3/2)
	for i := 0; i < 1e3/2; i += 2 {
		goodArray[i/2] = goodArray[i] + goodArray[i+1]
	}

}


