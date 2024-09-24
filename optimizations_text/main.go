package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

// https://syslog.ravelin.com/further-dangers-of-large-heaps-in-go-7a267b57d487
/*
	FURTHER DANGERS OF LARGE HEAPS IN GO
1. A large heap size is a problem. The GC needs to look at all 	allocated  memory to see which parts refer to other allocations. It's looking for memory that isn't referred to by any
	piece of memory that we're actually using, because those are the pieces of memory it can free for reuse. And to do that it has to scan through the memory looking for pointers.

2. If we have a large heap, a large amount of allocated memory that you need to keep through the lifetime of a process (for example large lookups tables, or an in-memory database of
	some kind), then to keep the amoun of GC work down you essentially have two choices as follows:
	1. Make sure the memory we allocate contains no pointers. That means:
		- no slices
		- no string
		- no time.Time
		- and definitely no pounters to other allocations.
	If an allocation has no pointers it gets marked as such and the GC doesn't scan it.
	2. Allocate the memory off-heap directly calling the "mmap syscall" by yourself. Then the GC knows nothing about the memory. This has upsides and downsides.
		+ The downside is that this memory can't really be used to reference objects allocated normally, as the GC may think they're no longer in-use and free them.

	If we don't follow either of these practices, and we allocate 50 GB that's kept around for the lifetime of a process, then every GC cycle will scan every bit of that 50 GB. And
	that will take some time. In addition, the GC will set it's memory use target to 100 GB, which may be more memory than we have overall.
*/

/*
		1.
		In this case memory use grows and grows, and in not much more than a minute the process is killed by the OS when memory runs out. Here's what we see if we enable GC trace debug
		output.

		GODEBUG=gctrace=1 ./gcbacklog
		Background GC work generated
		gc 1 @0.804s 21%: 0.012+4528+0.17 ms clock, 0.099+1.4/9054/27147+1.4 ms cpu, 11444->11444->11444 MB, 11445 MB goal, 8 P (forced)
		gc 2 @5.333s 23%: 0.012+6358+0.086 ms clock, 0.099+0/12716/38112+0.68 ms cpu, 11444->11444->11444 MB, 22888 MB goal, 8 P (forced)
		gc 3 @11.764s 31%: 20+53853+1.4 ms clock, 167+37787/107690/0+11 ms cpu, 11505->728829->728783 MB, 22888 MB goal, 8 P
		gc 4 @65.676s 40%: 69+10843+0.036 ms clock, 555+61294/21670/23+0.29 ms cpu, 728844->752155->34785 MB, 1457567 MB goal, 8 P
		Killed: 9

		Once the initial array is allocated the process is using 21% of the available CPU for the GC, and this rises to 40% before it's killed. The GC memory size tartget is quickly 22
		GB (twice our initial allocation), but this rises to an insane 1.4 TB as thins spiral out of control.

	 2. Now, if we change that initial allocation from 1.5 billion
	    pointers to 1.5 billion 8-byte integer things change completely. We use just as much memory, but it doesn't contain pointers. The GC target hits 22 GB, but the GC kicks in more
		frequently and uses less overall CPU, and importantly the target doesn't grow.

	    gc 61 @93.824s 0%: 4.0+4.5+0.075 ms clock, 32+8.9/8.6/4.0+0.60 ms cpu, 22412->22412->11474 MB, 22980 MB goal, 8 P
	    gc 62 @95.290s 0%: 14+4.0+0.085 ms clock, 115+4.3/0.39/0+0.68 ms cpu, 22382->22382->11451 MB, 22949 MB goal, 8 P

	 3. So what are the lessons to learn here?
	    If you're using Go for data processing then we either:
	    - Can't have any long-term large heap allocations
	    - Must ensure that they don't contain any pointers

	    And this means:
	    - no strings
	    - no slices
	    - no time.Time (it contains a pointer to a locale)
	    - no nothing with a pointer hidden pointer in it.
*/
func main() {
	// A huge allocation to give the GC work to do
	lotsOf := make([]*int, 15e8)

	fmt.Println("Background GC work generated")
	// Force a GC to set a baseline we can see if we set CODEBUG=gctrace=1
	runtime.GC()

	var (
		wg sync.WaitGroup
		// The optimal number of goroutines
		workersNumber = runtime.NumCPU()
	)

	for i := 0; i < workersNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work()
		}()
	}

	wg.Wait()

	// Make sure that this memory isn't optimised away
	runtime.KeepAlive(lotsOf)

}

func work() {
	for {
		work := make([]*int, 1e6)
		if f := factorial(20); f != 2432902008176640000 {
			fmt.Println(f)
			os.Exit(1)
		}

		runtime.KeepAlive(work)
	}
}

func factorial(n int) int {
	if n == 1 {
		return 1
	}

	return n * factorial(n-1)
}
